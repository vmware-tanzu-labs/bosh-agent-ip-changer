// Code generated by counterfeiter. DO NOT EDIT.
package vmfakes

import (
	"sync"

	"github.com/vmware-tanzu-labs/opsman-utils/pkg/vm"
)

type FakeSSHRunner struct {
	ExecuteStub        func(string, string, string, ...string) ([]byte, error)
	executeMutex       sync.RWMutex
	executeArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 []string
	}
	executeReturns struct {
		result1 []byte
		result2 error
	}
	executeReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSSHRunner) Execute(arg1 string, arg2 string, arg3 string, arg4 ...string) ([]byte, error) {
	fake.executeMutex.Lock()
	ret, specificReturn := fake.executeReturnsOnCall[len(fake.executeArgsForCall)]
	fake.executeArgsForCall = append(fake.executeArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 []string
	}{arg1, arg2, arg3, arg4})
	stub := fake.ExecuteStub
	fakeReturns := fake.executeReturns
	fake.recordInvocation("Execute", []interface{}{arg1, arg2, arg3, arg4})
	fake.executeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSSHRunner) ExecuteCallCount() int {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return len(fake.executeArgsForCall)
}

func (fake *FakeSSHRunner) ExecuteCalls(stub func(string, string, string, ...string) ([]byte, error)) {
	fake.executeMutex.Lock()
	defer fake.executeMutex.Unlock()
	fake.ExecuteStub = stub
}

func (fake *FakeSSHRunner) ExecuteArgsForCall(i int) (string, string, string, []string) {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	argsForCall := fake.executeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeSSHRunner) ExecuteReturns(result1 []byte, result2 error) {
	fake.executeMutex.Lock()
	defer fake.executeMutex.Unlock()
	fake.ExecuteStub = nil
	fake.executeReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeSSHRunner) ExecuteReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.executeMutex.Lock()
	defer fake.executeMutex.Unlock()
	fake.ExecuteStub = nil
	if fake.executeReturnsOnCall == nil {
		fake.executeReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.executeReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeSSHRunner) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSSHRunner) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ vm.SSHRunner = new(FakeSSHRunner)
