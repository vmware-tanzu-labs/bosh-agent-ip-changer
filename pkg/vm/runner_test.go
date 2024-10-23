package vm_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/bosh"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/vm"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/vm/vmfakes"
	"testing"
)

func TestGetVMsReturnsErr(t *testing.T) {
	ssh := &vmfakes.FakeSSHRunner{}
	om := &vmfakes.FakeOpsMan{}
	om.VMsReturns(nil, errors.New("can't connect to Opsman"))

	runner := vm.NewRunner(om, ssh)
	processed, err := runner.Execute([]string{"echo 'hello world'"}, vm.Filter{})
	require.Error(t, err)
	require.Equal(t, 0, processed)
}

func TestNoVMsReturned(t *testing.T) {
	vms := map[string][]bosh.VM{}

	opsmanClient := &vmfakes.FakeOpsMan{}
	opsmanClient.VMsReturns(vms, nil)

	sshClient := &vmfakes.FakeSSHRunner{}
	runner := vm.NewRunner(opsmanClient, sshClient)
	processed, err := runner.Execute([]string{"echo 'hello world'"}, vm.Filter{})
	require.NoError(t, err)
	require.Equal(t, 0, processed)
}

func TestGetVMPasswordReturnsErr(t *testing.T) {
	vms := map[string][]bosh.VM{
		"cf": {
			{
				InstanceName:  "vm1",
				Deployment:    "cf",
				InstanceGroup: "group1",
				JobState:      "running",
				IPs:           []string{"192.168.0.1"},
			},
		},
	}

	ssh := &vmfakes.FakeSSHRunner{}
	om := &vmfakes.FakeOpsMan{}
	om.VMsReturns(vms, nil)
	om.VMCredentialsReturnsOnCall(0, "", errors.New("can't get VM password"))

	runner := vm.NewRunner(om, ssh)
	processed, err := runner.Execute([]string{"echo 'hello world'"}, vm.Filter{})
	require.Error(t, err)
	require.Equal(t, 0, processed)
}

func TestOneVMProcessed(t *testing.T) {
	vms := map[string][]bosh.VM{
		"cf": {
			{
				InstanceName:  "vm1",
				Deployment:    "cf",
				InstanceGroup: "group1",
				JobState:      "running",
				IPs:           []string{"192.168.0.1"},
			},
		},
	}

	ssh := &vmfakes.FakeSSHRunner{}
	om := &vmfakes.FakeOpsMan{}
	om.VMsReturns(vms, nil)
	om.VMCredentialsReturnsOnCall(0, "vm1password", nil)

	runner := vm.NewRunner(om, ssh)
	processed, err := runner.Execute([]string{"echo 'hello world'"}, vm.Filter{})
	require.NoError(t, err)
	require.Equal(t, 1, processed)

	// validate the right SSH commands were executed
	hostAndPort, user, password, commands := ssh.ExecuteArgsForCall(0)
	require.Equal(t, "192.168.0.1:22", hostAndPort)
	require.Equal(t, "vcap", user)
	require.Equal(t, "vm1password", password)
	require.Equal(t, []string{"echo 'hello world'"}, commands)
}

func TestInstanceGroupFilters(t *testing.T) {
	vms := map[string][]bosh.VM{
		"cf": {
			{
				InstanceName:  "vm1",
				Deployment:    "cf",
				InstanceGroup: "diego_cell",
				JobState:      "running",
				IPs:           []string{"192.168.0.1"},
			},
			{
				InstanceName:  "vm2",
				Deployment:    "cf",
				InstanceGroup: "diego_cell",
				JobState:      "running",
				IPs:           []string{"192.168.0.2"},
			},
			{
				InstanceName:  "vm3",
				Deployment:    "cf",
				InstanceGroup: "router",
				JobState:      "running",
				IPs:           []string{"192.168.0.3"},
			},
			{
				InstanceName:  "vm4",
				Deployment:    "cf",
				InstanceGroup: "mysql",
				JobState:      "running",
				IPs:           []string{"192.168.0.4"},
			},
		},
		"mysql": {
			{
				InstanceName:  "vm5",
				Deployment:    "mysql",
				InstanceGroup: "mysql",
				JobState:      "running",
				IPs:           []string{"192.168.0.5"},
			},
			{
				InstanceName:  "vm6",
				Deployment:    "mysql",
				InstanceGroup: "mysql_proxy",
				JobState:      "running",
				IPs:           []string{"192.168.0.6"},
			},
		},
	}

	ssh := &vmfakes.FakeSSHRunner{}
	om := &vmfakes.FakeOpsMan{}
	om.VMsReturns(vms, nil)
	om.VMCredentialsReturns("password", nil)

	// filter by instance groups only
	runner := vm.NewRunner(om, ssh)
	processed, err := runner.Execute([]string{"echo 'hello world'"}, vm.Filter{
		InstanceGroups: []string{"diego_cell"},
	})
	require.NoError(t, err)
	require.Equal(t, 2, processed)

	hostAndPort, user, password, commands := ssh.ExecuteArgsForCall(0)
	require.Equal(t, "192.168.0.1:22", hostAndPort)
	require.Equal(t, "vcap", user)
	require.Equal(t, "password", password)
	require.Equal(t, []string{"echo 'hello world'"}, commands)

	hostAndPort, user, password, commands = ssh.ExecuteArgsForCall(1)
	require.Equal(t, "192.168.0.2:22", hostAndPort)

	// filter by deployment only
	ssh = &vmfakes.FakeSSHRunner{}
	runner = vm.NewRunner(om, ssh)
	processed, err = runner.Execute([]string{"echo 'hello world'"}, vm.Filter{
		Deployment: "mysql",
	})
	require.NoError(t, err)
	require.Equal(t, 2, processed)

	hostAndPort, user, password, commands = ssh.ExecuteArgsForCall(0)
	require.Equal(t, "192.168.0.5:22", hostAndPort)

	hostAndPort, user, password, commands = ssh.ExecuteArgsForCall(1)
	require.Equal(t, "192.168.0.6:22", hostAndPort)

	// filter by deployment and instance group
	ssh = &vmfakes.FakeSSHRunner{}
	runner = vm.NewRunner(om, ssh)
	processed, err = runner.Execute([]string{"echo 'hello world'"}, vm.Filter{
		Deployment:     "cf",
		InstanceGroups: []string{"mysql"},
	})
	require.NoError(t, err)
	require.Equal(t, 1, processed)

	hostAndPort, user, password, commands = ssh.ExecuteArgsForCall(0)
	require.Equal(t, "192.168.0.4:22", hostAndPort)
}
