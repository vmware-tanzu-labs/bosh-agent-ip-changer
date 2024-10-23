package command

import (
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/om"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/ssh"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/vm"
)

type Execute struct {
	omClient      *om.Client
	sshConnection *ssh.Connection
	Options       struct {
		Commands       []string `long:"command"        description:"command(s) to execute"                        required:"true"`
		Deployment     string   `long:"deployment"     description:"optional deployment name to filter by"        required:"false"`
		InstanceGroup  []string `long:"instance-group" description:"optional instance group name(s) to filter by" required:"false"`
		RestartJobName string   `long:"restart-job"    description:"optional instance job to restart"             required:"false"`
	}
}

func NewExecuteCommand(omClient *om.Client, sshConnection *ssh.Connection) *Execute {
	return &Execute{
		omClient:      omClient,
		sshConnection: sshConnection,
	}
}

// Execute runs the commands on the specified instances
func (f *Execute) Execute([]string) error {
	commands := f.Options.Commands
	if len(f.Options.RestartJobName) > 0 {
		commands = append(commands, "sudo killall -9 "+f.Options.RestartJobName)
	}

	filter := vm.Filter{
		Deployment:     f.Options.Deployment,
		InstanceGroups: f.Options.InstanceGroup,
	}

	runner := vm.NewRunner(f.omClient, f.sshConnection)
	_, err := runner.Execute(commands, filter)
	return err
}
