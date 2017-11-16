// Automatically generated by pegomock. DO NOT EDIT!
// Source: github.com/hootsuite/atlantis/server/events (interfaces: GitlabMergeRequestGetter)

package mocks

import (
	"reflect"

	go_gitlab "github.com/lkysow/go-gitlab"
	pegomock "github.com/petergtz/pegomock"
)

type MockGitlabMergeRequestGetter struct {
	fail func(message string, callerSkip ...int)
}

func NewMockGitlabMergeRequestGetter() *MockGitlabMergeRequestGetter {
	return &MockGitlabMergeRequestGetter{fail: pegomock.GlobalFailHandler}
}

func (mock *MockGitlabMergeRequestGetter) GetMergeRequest(repoFullName string, pullNum int) (*go_gitlab.MergeRequest, error) {
	params := []pegomock.Param{repoFullName, pullNum}
	result := pegomock.GetGenericMockFrom(mock).Invoke("GetMergeRequest", params, []reflect.Type{reflect.TypeOf((**go_gitlab.MergeRequest)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 *go_gitlab.MergeRequest
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(*go_gitlab.MergeRequest)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockGitlabMergeRequestGetter) VerifyWasCalledOnce() *VerifierGitlabMergeRequestGetter {
	return &VerifierGitlabMergeRequestGetter{mock, pegomock.Times(1), nil}
}

func (mock *MockGitlabMergeRequestGetter) VerifyWasCalled(invocationCountMatcher pegomock.Matcher) *VerifierGitlabMergeRequestGetter {
	return &VerifierGitlabMergeRequestGetter{mock, invocationCountMatcher, nil}
}

func (mock *MockGitlabMergeRequestGetter) VerifyWasCalledInOrder(invocationCountMatcher pegomock.Matcher, inOrderContext *pegomock.InOrderContext) *VerifierGitlabMergeRequestGetter {
	return &VerifierGitlabMergeRequestGetter{mock, invocationCountMatcher, inOrderContext}
}

type VerifierGitlabMergeRequestGetter struct {
	mock                   *MockGitlabMergeRequestGetter
	invocationCountMatcher pegomock.Matcher
	inOrderContext         *pegomock.InOrderContext
}

func (verifier *VerifierGitlabMergeRequestGetter) GetMergeRequest(repoFullName string, pullNum int) *GitlabMergeRequestGetter_GetMergeRequest_OngoingVerification {
	params := []pegomock.Param{repoFullName, pullNum}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "GetMergeRequest", params)
	return &GitlabMergeRequestGetter_GetMergeRequest_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type GitlabMergeRequestGetter_GetMergeRequest_OngoingVerification struct {
	mock              *MockGitlabMergeRequestGetter
	methodInvocations []pegomock.MethodInvocation
}

func (c *GitlabMergeRequestGetter_GetMergeRequest_OngoingVerification) GetCapturedArguments() (string, int) {
	repoFullName, pullNum := c.GetAllCapturedArguments()
	return repoFullName[len(repoFullName)-1], pullNum[len(pullNum)-1]
}

func (c *GitlabMergeRequestGetter_GetMergeRequest_OngoingVerification) GetAllCapturedArguments() (_param0 []string, _param1 []int) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]string, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(string)
		}
		_param1 = make([]int, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(int)
		}
	}
	return
}
