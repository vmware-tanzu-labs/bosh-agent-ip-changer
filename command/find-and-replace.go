package command

import (
	"fmt"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/om"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/vm"
)

type FindAndReplace struct {
	omClient *om.Client
	Options  struct {
		FilePath          string   `long:"file-path"          description:"full file path to perform the find and replace operation on" required:"true"`
		ReplaceExpression string   `long:"replace-expression" description:"find and replace expression using sed syntax"                required:"true"`
		Deployment        string   `long:"deployment"         description:"optional deployment name to filter by"                       required:"false"`
		InstanceGroup     []string `long:"instance-group"     description:"optional instance group name(s) to filter by"                required:"false"`
		RestartJobName    string   `long:"restart-job"        description:"optional instance job to restart"                            required:"false"`
	}
}

func NewFindAndReplaceCommand(omClient *om.Client) *FindAndReplace {
	return &FindAndReplace{
		omClient: omClient,
	}
}

// Execute runs the find and replace operation
func (f *FindAndReplace) Execute([]string) error {
	commands := []string{
		fmt.Sprintf("sudo sed -i '%s' '%s'", f.Options.ReplaceExpression, f.Options.FilePath),
	}
	if len(f.Options.RestartJobName) > 0 {
		commands = append(commands, "sudo killall -9 "+f.Options.RestartJobName)
	}

	filter := vm.Filter{
		Deployment:     f.Options.Deployment,
		InstanceGroups: f.Options.InstanceGroup,
	}

	runner := vm.NewRunner(f.omClient)
	return runner.Execute(commands, filter)
}
