// Package server is the main package for Atlantis. It handles the web server
// and executing commands that come in via pull request comments.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"flag"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/hootsuite/atlantis/server/events"
	"github.com/hootsuite/atlantis/server/events/locking"
	"github.com/hootsuite/atlantis/server/events/locking/boltdb"
	"github.com/hootsuite/atlantis/server/events/models"
	"github.com/hootsuite/atlantis/server/events/run"
	"github.com/hootsuite/atlantis/server/events/terraform"
	"github.com/hootsuite/atlantis/server/events/vcs"
	"github.com/hootsuite/atlantis/server/events/webhooks"
	"github.com/hootsuite/atlantis/server/logging"
	"github.com/hootsuite/atlantis/server/static"
	"github.com/lkysow/go-gitlab"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

const LockRouteName = "lock-detail"

// Server runs the Atlantis web server. It's used for webhook requests and the
// Atlantis UI.
type Server struct {
	Router             *mux.Router
	Port               int
	CommandHandler     *events.CommandHandler
	Logger             *logging.SimpleLogger
	Locker             locking.Locker
	AtlantisURL        string
	EventsController   *EventsController
	IndexTemplate      TemplateWriter
	LockDetailTemplate TemplateWriter
	Planner            *events.PlanExecutor
	GHHostname         string
	GHToken            string
	GHUser             string
}

// Config configures Server.
// The mapstructure tags correspond to flags in cmd/server.go and are used when
// the config is parsed from a YAML file.
type Config struct {
	AtlantisURL         string          `mapstructure:"atlantis-url"`
	DataDir             string          `mapstructure:"data-dir"`
	GithubHostname      string          `mapstructure:"gh-hostname"`
	GithubToken         string          `mapstructure:"gh-token"`
	GithubUser          string          `mapstructure:"gh-user"`
	GithubWebHookSecret string          `mapstructure:"gh-webhook-secret"`
	GitlabHostname      string          `mapstructure:"gitlab-hostname"`
	GitlabToken         string          `mapstructure:"gitlab-token"`
	GitlabUser          string          `mapstructure:"gitlab-user"`
	GitlabWebHookSecret string          `mapstructure:"gitlab-webhook-secret"`
	LogLevel            string          `mapstructure:"log-level"`
	Port                int             `mapstructure:"port"`
	RequireApproval     bool            `mapstructure:"require-approval"`
	SlackToken          string          `mapstructure:"slack-token"`
	Webhooks            []WebhookConfig `mapstructure:"webhooks"`
}

type WebhookConfig struct {
	Event          string `mapstructure:"event"`
	WorkspaceRegex string `mapstructure:"workspace-regex"`
	Kind           string `mapstructure:"kind"`
	// Slack specific
	Channel string `mapstructure:"channel"`
}

func NewServer(config Config) (*Server, error) {
	var supportedVCSHosts []vcs.Host
	var githubClient *vcs.GithubClient
	var gitlabClient *vcs.GitlabClient
	if config.GithubUser != "" {
		supportedVCSHosts = append(supportedVCSHosts, vcs.Github)
		var err error
		githubClient, err = vcs.NewGithubClient(config.GithubHostname, config.GithubUser, config.GithubToken)
		if err != nil {
			return nil, err
		}
	}
	if config.GitlabUser != "" {
		supportedVCSHosts = append(supportedVCSHosts, vcs.Gitlab)
		gitlabClient = &vcs.GitlabClient{
			Client: gitlab.NewClient(nil, config.GitlabToken),
		}
	}
	var webhooksConfig []webhooks.Config
	for _, c := range config.Webhooks {
		config := webhooks.Config{
			Channel:        c.Channel,
			Event:          c.Event,
			Kind:           c.Kind,
			WorkspaceRegex: c.WorkspaceRegex,
		}
		webhooksConfig = append(webhooksConfig, config)
	}
	webhooksManager, err := webhooks.NewMultiWebhookSender(webhooksConfig, webhooks.NewSlackClient(config.SlackToken))
	if err != nil {
		return nil, errors.Wrap(err, "initializing webhooks")
	}
	vcsClient := vcs.NewDefaultClientProxy(githubClient, gitlabClient)
	commitStatusUpdater := &events.DefaultCommitStatusUpdater{Client: vcsClient}
	terraformClient, err := terraform.NewClient()
	// The flag.Lookup call is to detect if we're running in a unit test. If we
	// are, then we don't error out because we don't have/want terraform
	// installed on our CI system where the unit tests run.
	if err != nil && flag.Lookup("test.v") == nil {
		return nil, errors.Wrap(err, "initializing terraform")
	}
	markdownRenderer := &events.MarkdownRenderer{}
	boltdb, err := boltdb.New(config.DataDir)
	if err != nil {
		return nil, err
	}
	lockingClient := locking.NewClient(boltdb)
	run := &run.Run{}
	configReader := &events.ProjectConfigManager{}
	concurrentRunLocker := events.NewEnvLock()
	workspace := &events.FileWorkspace{
		DataDir: config.DataDir,
	}
	projectPreExecute := &events.ProjectPreExecute{
		Locker:       lockingClient,
		Run:          run,
		ConfigReader: configReader,
		Terraform:    terraformClient,
	}
	applyExecutor := &events.ApplyExecutor{
		VCSClient:         vcsClient,
		Terraform:         terraformClient,
		RequireApproval:   config.RequireApproval,
		Run:               run,
		Workspace:         workspace,
		ProjectPreExecute: projectPreExecute,
		Webhooks:          webhooksManager,
	}
	planExecutor := &events.PlanExecutor{
		VCSClient:         vcsClient,
		Terraform:         terraformClient,
		Run:               run,
		Workspace:         workspace,
		ProjectPreExecute: projectPreExecute,
		Locker:            lockingClient,
		ProjectFinder:     &events.ProjectFinder{},
	}
	helpExecutor := &events.HelpExecutor{}
	pullClosedExecutor := &events.PullClosedExecutor{
		VCSClient: vcsClient,
		Locker:    lockingClient,
		Workspace: workspace,
	}
	logger := logging.NewSimpleLogger("server", nil, false, logging.ToLogLevel(config.LogLevel))
	eventParser := &events.EventParser{
		GithubUser:  config.GithubUser,
		GithubToken: config.GithubToken,
		GitlabUser:  config.GitlabUser,
		GitlabToken: config.GitlabToken,
	}
	commandHandler := &events.CommandHandler{
		ApplyExecutor:            applyExecutor,
		PlanExecutor:             planExecutor,
		HelpExecutor:             helpExecutor,
		LockURLGenerator:         planExecutor,
		EventParser:              eventParser,
		VCSClient:                vcsClient,
		GithubPullGetter:         githubClient,
		GitlabMergeRequestGetter: gitlabClient,
		CommitStatusUpdater:      commitStatusUpdater,
		EnvLocker:                concurrentRunLocker,
		MarkdownRenderer:         markdownRenderer,
		Logger:                   logger,
	}
	eventsController := &EventsController{
		CommandRunner:          commandHandler,
		PullCleaner:            pullClosedExecutor,
		Parser:                 eventParser,
		Logger:                 logger,
		GithubWebHookSecret:    []byte(config.GithubWebHookSecret),
		GithubRequestValidator: &DefaultGithubRequestValidator{},
		GitlabRequestParser:    &DefaultGitlabRequestParser{},
		GitlabWebHookSecret:    []byte(config.GitlabWebHookSecret),
		SupportedVCSHosts:      supportedVCSHosts,
	}
	router := mux.NewRouter()
	return &Server{
		Router:             router,
		Port:               config.Port,
		CommandHandler:     commandHandler,
		Logger:             logger,
		Locker:             lockingClient,
		AtlantisURL:        config.AtlantisURL,
		EventsController:   eventsController,
		IndexTemplate:      indexTemplate,
		LockDetailTemplate: lockTemplate,
		Planner:            planExecutor,
		GHHostname:         config.GithubHostname,
		GHToken:            config.GithubToken,
		GHUser:             config.GithubUser,
	}, nil
}

func (s *Server) Start() error {
	s.Router.HandleFunc("/", s.Index).Methods("GET").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.URL.Path == "/" || r.URL.Path == "/index.html"
	})
	s.Router.PathPrefix("/static/").Handler(http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, AssetInfo: static.AssetInfo}))
	s.Router.HandleFunc("/events", s.postEvents).Methods("POST")
	s.Router.HandleFunc("/locks", s.DeleteLockRoute).Methods("DELETE").Queries("id", "{id:.*}")
	lockRoute := s.Router.HandleFunc("/lock", s.GetLockRoute).Methods("GET").Queries("id", "{id}").Name(LockRouteName)

	s.Router.HandleFunc("/plans", s.plan).Methods("POST")

	// function that planExecutor can use to construct detail view url
	// injecting this here because this is the earliest routes are created
	s.CommandHandler.SetLockURL(func(lockID string) string {
		// ignoring error since guaranteed to succeed if "id" is specified
		u, _ := lockRoute.URL("id", url.QueryEscape(lockID))
		return s.AtlantisURL + u.RequestURI()
	})
	n := negroni.New(&negroni.Recovery{
		Logger:     log.New(os.Stdout, "", log.LstdFlags),
		PrintStack: false,
		StackAll:   false,
		StackSize:  1024 * 8,
	}, NewRequestLogger(s.Logger))
	n.UseHandler(s.Router)
	s.Logger.Warn("Atlantis started - listening on port %v", s.Port)
	return cli.NewExitError(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), n), 1)
}

type PlanRequestBody struct {
	Repo      string
	Branch    string
	Path      string
	Workspace string
	Version   string
}

type PlanResponseBody struct {
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Changes bool   `json:"changes"`
}

func (s *Server) plan(w http.ResponseWriter, r *http.Request) {
	// Deal with the request body.
	d := json.NewDecoder(r.Body)
	var reqBody PlanRequestBody
	err := d.Decode(&reqBody)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// Set up Plan.
	repoFullName := reqBody.Repo
	owner := strings.Split(repoFullName, "/")[0]
	name := strings.Split(repoFullName, "/")[1]
	cloneURL := fmt.Sprintf("https://%s:%s@%s/%s.git", s.GHUser, s.GHToken, s.GHHostname, repoFullName)
	sanitizedCloneURL := fmt.Sprintf("https://%s/%s.git", s.GHHostname, repoFullName)

	repo := models.Repo{
		FullName:          repoFullName,
		Owner:             owner,
		Name:              name,
		CloneURL:          cloneURL,
		SanitizedCloneURL: sanitizedCloneURL,
	}
	// Defaults.
	branch := "master"
	path := "."
	if reqBody.Branch != "" {
		branch = reqBody.Branch
	}
	if reqBody.Path != "" {
		path = reqBody.Path
	}
	pull := models.PullRequest{
		Num:    0,
		Branch: branch,
	}
	user := models.User{
		Username: "git-atlantis-user",
	}
	cmd := &events.Command{
		Name:        events.Plan,
		Environment: reqBody.Workspace,
		Verbose:     true,
		Flags:       []string{},
	}
	ctx := events.CommandContext{
		BaseRepo: repo,
		HeadRepo: repo,
		Pull:     pull,
		User:     user,
		Command:  cmd,
		Log:      s.Logger,
	}
	projects := []models.Project{
		models.Project{
			RepoFullName: repoFullName,
			Path:         path,
		},
	}

	// Run plan and respond with result.
	cr := s.Planner.RunPlan(&ctx, projects)
	if len(cr.ProjectResults) == 0 {
		s.respond(w, logging.Error, http.StatusInternalServerError, "No project results from plan")
		return
	}
	res := cr.ProjectResults[0]
	respStruct := PlanResponseBody{}
	switch res.Status() {
	case vcs.Success:
		respStruct.Success = true
		respStruct.Output = res.PlanSuccess.TerraformOutput
		respStruct.Changes = res.PlanSuccess.Changes
	case vcs.Failed:
		respStruct.Success = false
		if err := res.Error; err != nil {
			respStruct.Output = res.Error.Error()
		} else {
			respStruct.Output = res.Failure
		}
	}
	respBody, err := json.Marshal(respStruct)
	if err != nil {
		panic("oh no, can't json.Marshal(...)!")
	}

	w.Header().Set("Content-Type", "application/json")
	// s.respond(w, logging.Info, http.StatusOK, "Plan executed")
	w.Write(respBody)
}

func (s *Server) Index(w http.ResponseWriter, _ *http.Request) {
	locks, err := s.Locker.List()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Could not retrieve locks: %s", err)
		return
	}

	var results []LockIndexData
	for id, v := range locks {
		lockURL, _ := s.Router.Get(LockRouteName).URL("id", url.QueryEscape(id))
		results = append(results, LockIndexData{
			LockURL:      lockURL.String(),
			RepoFullName: v.Project.RepoFullName,
			PullNum:      v.Pull.Num,
			Time:         v.Time,
		})
	}
	s.IndexTemplate.Execute(w, results) // nolint: errcheck
}

func (s *Server) GetLockRoute(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No lock id in request")
		return
	}
	s.GetLock(w, r, id)
}

// GetLock handles a lock detail page view. getLockRoute is expected to
// be called before. This function was extracted to make it testable.
func (s *Server) GetLock(w http.ResponseWriter, _ *http.Request, id string) {
	// get details for lock id
	idUnencoded, err := url.QueryUnescape(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid lock id")
		return
	}

	// for the given lock key get lock data
	lock, err := s.Locker.GetLock(idUnencoded)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	if lock == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "No lock found at that id")
		return
	}

	// extract the repo owner and repo name
	repo := strings.Split(lock.Project.RepoFullName, "/")

	l := LockDetailData{
		LockKeyEncoded:  id,
		LockKey:         idUnencoded,
		RepoOwner:       repo[0],
		RepoName:        repo[1],
		PullRequestLink: lock.Pull.URL,
		LockedBy:        lock.Pull.Author,
		Environment:     lock.Env,
	}

	s.LockDetailTemplate.Execute(w, l) // nolint: errcheck
}

func (s *Server) DeleteLockRoute(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		s.respond(w, logging.Warn, http.StatusBadRequest, "No lock id in request")
		return
	}
	s.DeleteLock(w, r, id)
}

func (s *Server) DeleteLock(w http.ResponseWriter, _ *http.Request, id string) {
	idUnencoded, err := url.PathUnescape(id)
	if err != nil {
		s.respond(w, logging.Warn, http.StatusBadRequest, "Invalid lock id: %s", err)
		return
	}
	lock, err := s.Locker.Unlock(idUnencoded)
	if err != nil {
		s.respond(w, logging.Error, http.StatusInternalServerError, "Failed to delete lock %s: %s", idUnencoded, err)
		return
	}
	if lock == nil {
		s.respond(w, logging.Warn, http.StatusNotFound, "No lock found at that id", idUnencoded)
		return
	}
	s.respond(w, logging.Info, http.StatusOK, "Deleted lock id %s", idUnencoded)
}

// postEvents handles POST requests to our /events endpoint. These should be
// VCS webhook requests.
func (s *Server) postEvents(w http.ResponseWriter, r *http.Request) {
	s.EventsController.Post(w, r)
}
func (s *Server) respond(w http.ResponseWriter, lvl logging.LogLevel, code int, format string, args ...interface{}) {
	response := fmt.Sprintf(format, args...)
	s.Logger.Log(lvl, response)
	w.WriteHeader(code)
	fmt.Fprintln(w, response)
}
