// Automatically generated by pegomock. DO NOT EDIT!
// Source: github.com/hootsuite/atlantis/server/events/webhooks (interfaces: WebhookSender)

package mocks

import (
	webhooks "github.com/hootsuite/atlantis/server/events/webhooks"
	pegomock "github.com/petergtz/pegomock"
	"reflect"
)

type MockWebhookSender struct {
	fail func(message string, callerSkip ...int)
}

func NewMockWebhookSender() *MockWebhookSender {
	return &MockWebhookSender{fail: pegomock.GlobalFailHandler}
}

func (mock *MockWebhookSender) Send(_param0 webhooks.ApplyResult) error {
	params := []pegomock.Param{_param0}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Send", params, []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(error)
		}
	}
	return ret0
}

func (mock *MockWebhookSender) VerifyWasCalledOnce() *VerifierWebhookSender {
	return &VerifierWebhookSender{mock, pegomock.Times(1), nil}
}

func (mock *MockWebhookSender) VerifyWasCalled(invocationCountMatcher pegomock.Matcher) *VerifierWebhookSender {
	return &VerifierWebhookSender{mock, invocationCountMatcher, nil}
}

func (mock *MockWebhookSender) VerifyWasCalledInOrder(invocationCountMatcher pegomock.Matcher, inOrderContext *pegomock.InOrderContext) *VerifierWebhookSender {
	return &VerifierWebhookSender{mock, invocationCountMatcher, inOrderContext}
}

type VerifierWebhookSender struct {
	mock                   *MockWebhookSender
	invocationCountMatcher pegomock.Matcher
	inOrderContext         *pegomock.InOrderContext
}

func (verifier *VerifierWebhookSender) Send(_param0 webhooks.ApplyResult) *WebhookSender_Send_OngoingVerification {
	params := []pegomock.Param{_param0}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Send", params)
	return &WebhookSender_Send_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type WebhookSender_Send_OngoingVerification struct {
	mock              *MockWebhookSender
	methodInvocations []pegomock.MethodInvocation
}

func (c *WebhookSender_Send_OngoingVerification) GetCapturedArguments() webhooks.ApplyResult {
	_param0 := c.GetAllCapturedArguments()
	return _param0[len(_param0)-1]
}

func (c *WebhookSender_Send_OngoingVerification) GetAllCapturedArguments() (_param0 []webhooks.ApplyResult) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]webhooks.ApplyResult, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(webhooks.ApplyResult)
		}
	}
	return
}