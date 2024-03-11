package command

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu-labs/bosh-agent-ip-changer/pkg/om"
	"github.com/vmware-tanzu-labs/bosh-agent-ip-changer/ssh_client"
)

type ChangeIP struct {
	CACert            string `long:"ca-cert" env:"OM_CA_CERT" description:"OpsManager CA certificate path or value"`
	ClientID          string `short:"c"  long:"client-id"             env:"OM_CLIENT_ID"                           description:"Client ID for the Ops Manager VM"`
	ClientSecret      string `short:"s"  long:"client-secret"         env:"OM_CLIENT_SECRET"                       description:"Client Secret for the Ops Manager VM"`
	Password          string `short:"p"  long:"password"              env:"OM_PASSWORD"                            description:"admin password for the Ops Manager VM"`
	SkipSSLValidation bool   `short:"k"  long:"skip-ssl-validation"   env:"OM_SKIP_SSL_VALIDATION"                 description:"skip ssl certificate validation during http requests"`
	Username          string `short:"u"  long:"username"              env:"OM_USERNAME"                            description:"admin username for the Ops Manager VM"`
	Target            string `short:"t"  long:"target"                env:"OM_TARGET"                              description:"location of the Ops Manager VM" required:"true"`
	Debug             bool   `long:"debug"  description:"sets log level to debug"`
	OldDirectorIP     string `long:"old-director-ip"   description:"previous ip of bosh director" required:"true"`
	NewDirectorIP     string `long:"new-director-ip"   description:"new ip of bosh director" required:"true"`
}

func initLogging(debug bool) {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// Execute - runs the migration
func (c *ChangeIP) Execute([]string) error {
	initLogging(c.Debug)
	client := om.New(c.Target, c.Username, c.Password, c.ClientID, c.ClientSecret, c.SkipSSLValidation, c.CACert)
	boshVMs, err := client.VMs()
	if err != nil {
		return err
	}
	for deployment, vms := range boshVMs {
		for _, vm := range vms {
			if strings.EqualFold("running", vm.JobState) {
				log.Infof("skipping vm %s of deployment %s as already healthy bosh agent", vm.InstanceName, deployment)
				continue
			}
			log.Infof("Processing vm %s of deployment %s", vm.InstanceName, deployment)
			password, err := client.VMCredentials(deployment, vm.JobName)
			if err != nil {
				return err
			}

			conn, err := ssh_client.Connect(fmt.Sprintf("%s:22", vm.IPs[0]), "vcap", password)
			if err != nil {
				return err
			}
			_, err = conn.SendCommands(
				fmt.Sprintf("sudo sed -i 's/%s/%s/g' /var/vcap/bosh/settings.json", c.OldDirectorIP, c.NewDirectorIP),
				fmt.Sprintf("sudo sed -i 's/%s/%s/g' /var/vcap/bosh/spec.json", c.OldDirectorIP, c.NewDirectorIP),
				fmt.Sprintf("sudo sed -i 's/%s/%s/g' /var/vcap/bosh/etc/blobstore-dav.json", c.OldDirectorIP, c.NewDirectorIP),
				"sudo killall -9 bosh-agent")
			if err != nil {
				return err
			}

		}
	}
	//ctx := context.Background()

	// c, err := m.combinedConfig()
	// if err != nil {
	// 	return err
	// }

	// return migrate.RunFoundationMigrationWithConfig(c, ctx)
	return nil
}
