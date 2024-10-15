package vm

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/om"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/ssh"
)

type Runner struct {
	omClient *om.Client
}

func NewRunner(omClient *om.Client) *Runner {
	return &Runner{
		omClient: omClient,
	}
}

// Execute runs the commands on all the VMs that match the specified filters
func (r Runner) Execute(commands []string, filter Filter) error {
	if commands == nil || len(commands) == 0 {
		return errors.New("no commands to execute")
	}

	boshVMs, err := r.omClient.VMs()
	if err != nil {
		return err
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
				return err
			}

			conn, err := ssh.Connect(fmt.Sprintf("%s:22", vm.IPs[0]), "vcap", password)
			if err != nil {
				return err
			}
			_, err = conn.SendCommands(commands...)
			if err != nil {
				return err
			}
			count++
		}
	}

	log.Infof("Processed %d vms", count)

	return nil
}
