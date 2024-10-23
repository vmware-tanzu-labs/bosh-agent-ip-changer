package vm

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/bosh"
)

const boshUser = "vcap"

//counterfeiter:generate . OpsMan
type OpsMan interface {
	VMs() (map[string][]bosh.VM, error)
	VMCredentials(deployment, instanceGroup string) (string, error)
}

//counterfeiter:generate . SSHRunner
type SSHRunner interface {
	Execute(address, user, password string, commands ...string) error
}

type Runner struct {
	omClient  OpsMan
	sshClient SSHRunner
}

func NewRunner(omClient OpsMan, sshClient SSHRunner) *Runner {
	return &Runner{
		omClient:  omClient,
		sshClient: sshClient,
	}
}

// Execute runs the commands on all the VMs that match the specified filters
func (r Runner) Execute(commands []string, filter Filter) (int, error) {
	if commands == nil || len(commands) == 0 {
		return 0, errors.New("no commands to execute")
	}

	boshVMs, err := r.omClient.VMs()
	if err != nil {
		return 0, err
	}

	count := 0
	for deployment, vms := range boshVMs {
		if !filter.ShouldProcessDeployment(deployment) {
			continue
		}
		for _, vm := range vms {
			if !filter.ShouldProcessInstanceGroup(vm.InstanceGroup) {
				continue
			}
			log.Infof("Processing vm %s of deployment %s", vm.InstanceName, deployment)
			password, err := r.omClient.VMCredentials(deployment, vm.InstanceGroup)
			if err != nil {
				return count, err
			}

			addr := fmt.Sprintf("%s:22", vm.IPs[0])
			err = r.sshClient.Execute(addr, boshUser, password, commands...)
			if err != nil {
				return count, err
			}
			count++
		}
	}

	log.Infof("Processed %d vms", count)

	return count, nil
}
