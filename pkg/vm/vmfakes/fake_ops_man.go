// Code generated by counterfeiter. DO NOT EDIT.
package vmfakes

import (
	"sync"

	"github.com/vmware-tanzu-labs/opsman-utils/pkg/bosh"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/vm"
)

type FakeOpsMan struct {
	VMCredentialsStub        func(string, string) (string, error)
	vMCredentialsMutex       sync.RWMutex
	vMCredentialsArgsForCall []struct {
		arg1 string
		arg2 string
	}
	vMCredentialsReturns struct {
		result1 string
		result2 error
	}
	vMCredentialsReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	VMsStub        func() (map[string][]bosh.VM, error)
	vMsMutex       sync.RWMutex
	vMsArgsForCall []struct {
	}
	vMsReturns struct {
		result1 map[string][]bosh.VM
		result2 error
	}
	vMsReturnsOnCall map[int]struct {
		result1 map[string][]bosh.VM
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOpsMan) VMCredentials(arg1 string, arg2 string) (string, error) {
	fake.vMCredentialsMutex.Lock()
	ret, specificReturn := fake.vMCredentialsReturnsOnCall[len(fake.vMCredentialsArgsForCall)]
	fake.vMCredentialsArgsForCall = append(fake.vMCredentialsArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.VMCredentialsStub
	fakeReturns := fake.vMCredentialsReturns
	fake.recordInvocation("VMCredentials", []interface{}{arg1, arg2})
	fake.vMCredentialsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeOpsMan) VMCredentialsCallCount() int {
	fake.vMCredentialsMutex.RLock()
	defer fake.vMCredentialsMutex.RUnlock()
	return len(fake.vMCredentialsArgsForCall)
}

func (fake *FakeOpsMan) VMCredentialsCalls(stub func(string, string) (string, error)) {
	fake.vMCredentialsMutex.Lock()
	defer fake.vMCredentialsMutex.Unlock()
	fake.VMCredentialsStub = stub
}

func (fake *FakeOpsMan) VMCredentialsArgsForCall(i int) (string, string) {
	fake.vMCredentialsMutex.RLock()
	defer fake.vMCredentialsMutex.RUnlock()
	argsForCall := fake.vMCredentialsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeOpsMan) VMCredentialsReturns(result1 string, result2 error) {
	fake.vMCredentialsMutex.Lock()
	defer fake.vMCredentialsMutex.Unlock()
	fake.VMCredentialsStub = nil
	fake.vMCredentialsReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeOpsMan) VMCredentialsReturnsOnCall(i int, result1 string, result2 error) {
	fake.vMCredentialsMutex.Lock()
	defer fake.vMCredentialsMutex.Unlock()
	fake.VMCredentialsStub = nil
	if fake.vMCredentialsReturnsOnCall == nil {
		fake.vMCredentialsReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.vMCredentialsReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeOpsMan) VMs() (map[string][]bosh.VM, error) {
	fake.vMsMutex.Lock()
	ret, specificReturn := fake.vMsReturnsOnCall[len(fake.vMsArgsForCall)]
	fake.vMsArgsForCall = append(fake.vMsArgsForCall, struct {
	}{})
	stub := fake.VMsStub
	fakeReturns := fake.vMsReturns
	fake.recordInvocation("VMs", []interface{}{})
	fake.vMsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeOpsMan) VMsCallCount() int {
	fake.vMsMutex.RLock()
	defer fake.vMsMutex.RUnlock()
	return len(fake.vMsArgsForCall)
}

func (fake *FakeOpsMan) VMsCalls(stub func() (map[string][]bosh.VM, error)) {
	fake.vMsMutex.Lock()
	defer fake.vMsMutex.Unlock()
	fake.VMsStub = stub
}

func (fake *FakeOpsMan) VMsReturns(result1 map[string][]bosh.VM, result2 error) {
	fake.vMsMutex.Lock()
	defer fake.vMsMutex.Unlock()
	fake.VMsStub = nil
	fake.vMsReturns = struct {
		result1 map[string][]bosh.VM
		result2 error
	}{result1, result2}
}

func (fake *FakeOpsMan) VMsReturnsOnCall(i int, result1 map[string][]bosh.VM, result2 error) {
	fake.vMsMutex.Lock()
	defer fake.vMsMutex.Unlock()
	fake.VMsStub = nil
	if fake.vMsReturnsOnCall == nil {
		fake.vMsReturnsOnCall = make(map[int]struct {
			result1 map[string][]bosh.VM
			result2 error
		})
	}
	fake.vMsReturnsOnCall[i] = struct {
		result1 map[string][]bosh.VM
		result2 error
	}{result1, result2}
}

func (fake *FakeOpsMan) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.vMCredentialsMutex.RLock()
	defer fake.vMCredentialsMutex.RUnlock()
	fake.vMsMutex.RLock()
	defer fake.vMsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOpsMan) recordInvocation(key string, args []interface{}) {
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

var _ vm.OpsMan = new(FakeOpsMan)
