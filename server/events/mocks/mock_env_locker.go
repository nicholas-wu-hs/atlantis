// Automatically generated by pegomock. DO NOT EDIT!
// Source: github.com/hootsuite/atlantis/server/events (interfaces: EnvLocker)

package mocks

import (
	"reflect"

	pegomock "github.com/petergtz/pegomock"
)

type MockEnvLocker struct {
	fail func(message string, callerSkip ...int)
}

func NewMockEnvLocker() *MockEnvLocker {
	return &MockEnvLocker{fail: pegomock.GlobalFailHandler}
}

func (mock *MockEnvLocker) TryLock(repoFullName string, env string, pullNum int) bool {
	params := []pegomock.Param{repoFullName, env, pullNum}
	result := pegomock.GetGenericMockFrom(mock).Invoke("TryLock", params, []reflect.Type{reflect.TypeOf((*bool)(nil)).Elem()})
	var ret0 bool
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(bool)
		}
	}
	return ret0
}

func (mock *MockEnvLocker) Unlock(repoFullName string, env string, pullNum int) {
	params := []pegomock.Param{repoFullName, env, pullNum}
	pegomock.GetGenericMockFrom(mock).Invoke("Unlock", params, []reflect.Type{})
}

func (mock *MockEnvLocker) VerifyWasCalledOnce() *VerifierEnvLocker {
	return &VerifierEnvLocker{mock, pegomock.Times(1), nil}
}

func (mock *MockEnvLocker) VerifyWasCalled(invocationCountMatcher pegomock.Matcher) *VerifierEnvLocker {
	return &VerifierEnvLocker{mock, invocationCountMatcher, nil}
}

func (mock *MockEnvLocker) VerifyWasCalledInOrder(invocationCountMatcher pegomock.Matcher, inOrderContext *pegomock.InOrderContext) *VerifierEnvLocker {
	return &VerifierEnvLocker{mock, invocationCountMatcher, inOrderContext}
}

type VerifierEnvLocker struct {
	mock                   *MockEnvLocker
	invocationCountMatcher pegomock.Matcher
	inOrderContext         *pegomock.InOrderContext
}

func (verifier *VerifierEnvLocker) TryLock(repoFullName string, env string, pullNum int) *EnvLocker_TryLock_OngoingVerification {
	params := []pegomock.Param{repoFullName, env, pullNum}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "TryLock", params)
	return &EnvLocker_TryLock_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type EnvLocker_TryLock_OngoingVerification struct {
	mock              *MockEnvLocker
	methodInvocations []pegomock.MethodInvocation
}

func (c *EnvLocker_TryLock_OngoingVerification) GetCapturedArguments() (string, string, int) {
	repoFullName, env, pullNum := c.GetAllCapturedArguments()
	return repoFullName[len(repoFullName)-1], env[len(env)-1], pullNum[len(pullNum)-1]
}

func (c *EnvLocker_TryLock_OngoingVerification) GetAllCapturedArguments() (_param0 []string, _param1 []string, _param2 []int) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]string, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(string)
		}
		_param1 = make([]string, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(string)
		}
		_param2 = make([]int, len(params[2]))
		for u, param := range params[2] {
			_param2[u] = param.(int)
		}
	}
	return
}

func (verifier *VerifierEnvLocker) Unlock(repoFullName string, env string, pullNum int) *EnvLocker_Unlock_OngoingVerification {
	params := []pegomock.Param{repoFullName, env, pullNum}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Unlock", params)
	return &EnvLocker_Unlock_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type EnvLocker_Unlock_OngoingVerification struct {
	mock              *MockEnvLocker
	methodInvocations []pegomock.MethodInvocation
}

func (c *EnvLocker_Unlock_OngoingVerification) GetCapturedArguments() (string, string, int) {
	repoFullName, env, pullNum := c.GetAllCapturedArguments()
	return repoFullName[len(repoFullName)-1], env[len(env)-1], pullNum[len(pullNum)-1]
}

func (c *EnvLocker_Unlock_OngoingVerification) GetAllCapturedArguments() (_param0 []string, _param1 []string, _param2 []int) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]string, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(string)
		}
		_param1 = make([]string, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(string)
		}
		_param2 = make([]int, len(params[2]))
		for u, param := range params[2] {
			_param2[u] = param.(int)
		}
	}
	return
}
