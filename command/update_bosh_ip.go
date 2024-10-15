package command

import (
	"fmt"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/om"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/vm"
)

type UpdateBoshIP struct {
	omClient *om.Client
	Options  struct {
		OldDirectorIP string `long:"old-director-ip"   description:"previous ip of bosh director" required:"true"`
		NewDirectorIP string `long:"new-director-ip"   description:"new ip of bosh director" required:"true"`
	}
}

func NewUpdateBoshIPCommand(omClient *om.Client) *UpdateBoshIP {
	return &UpdateBoshIP{
		omClient: omClient,
	}
}

// Execute - runs the migration
func (u *UpdateBoshIP) Execute([]string) error {
	commands := []string{
		fmt.Sprintf("sudo sed -i 's/%s/%s/g' /var/vcap/bosh/settings.json", u.Options.OldDirectorIP, u.Options.NewDirectorIP),
		fmt.Sprintf("sudo sed -i 's/%s/%s/g' /var/vcap/bosh/spec.json", u.Options.OldDirectorIP, u.Options.NewDirectorIP),
		fmt.Sprintf("sudo sed -i 's/%s/%s/g' /var/vcap/bosh/etc/blobstore-dav.json", u.Options.OldDirectorIP, u.Options.NewDirectorIP),
		"sudo killall -9 bosh-agent",
	}

	runner := vm.NewRunner(u.omClient)
	return runner.Execute(commands, vm.Filter{})
}
