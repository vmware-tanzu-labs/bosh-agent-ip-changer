package command

import (
	"errors"
	goflags "github.com/jessevdk/go-flags"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/om"
	"os"
)

type globalOptions struct {
	CACert            string `long:"ca-cert" env:"OM_CA_CERT" description:"OpsManager CA certificate path or value"`
	ClientID          string `short:"c"  long:"client-id"             env:"OM_CLIENT_ID"                           description:"Client ID for the Ops Manager VM"`
	ClientSecret      string `short:"s"  long:"client-secret"         env:"OM_CLIENT_SECRET"                       description:"Client Secret for the Ops Manager VM"`
	Password          string `short:"p"  long:"password"              env:"OM_PASSWORD"                            description:"admin password for the Ops Manager VM"`
	SkipSSLValidation bool   `short:"k"  long:"skip-ssl-validation"   env:"OM_SKIP_SSL_VALIDATION"                 description:"skip ssl certificate validation during http requests"`
	Username          string `short:"u"  long:"username"              env:"OM_USERNAME"                            description:"admin username for the Ops Manager VM"`
	Target            string `short:"t"  long:"target"                env:"OM_TARGET"                              description:"location of the Ops Manager VM" required:"true"`

	Debug   bool `long:"debug" description:"sets log level to debug"`
	Version bool `short:"v"  long:"version" description:"prints the release version"`
}

func Main(args []string) error {
	var global globalOptions
	parser := goflags.NewParser(&global, goflags.PassDoubleDash|goflags.PassAfterNonOption)
	parser.Name = "opsmanutils"

	// allow version command instead of only version flag
	if len(args) > 1 && args[1] == "version" {
		args[1] = "--version"
	}

	// parse only the global options at this point
	args, _ = parser.ParseArgs(args[1:])

	// if --version was specified print and exit immediately
	if global.Version {
		return NewVersionCommand().Execute(nil)
	}

	initLogging(global.Debug)

	omClient := om.New(
		global.Target,
		global.Username,
		global.Password,
		global.ClientID,
		global.ClientSecret,
		global.SkipSSLValidation,
		global.CACert)

	_, err := parser.AddCommand("update-bosh-ip",
		"Updates every bosh agent to use the specified BOSH Director IP",
		"Updates every bosh agent to use the specified BOSH Director IP",
		NewUpdateBoshIPCommand(omClient))
	if err != nil {
		return err
	}

	_, err = parser.AddCommand("find-and-replace",
		"Performs a text find and replace operation against the specified instances",
		"Performs a text find and replace operation against the specified instances",
		NewFindAndReplaceCommand(omClient))
	if err != nil {
		return err
	}

	// allow help command in addition to help flag
	if len(args) > 0 && args[0] == "help" {
		args[0] = "--help"
	}

	// if parsing of the command sub-options fails dump out the CLI help message
	parser.Options |= goflags.HelpFlag
	_, err = parser.ParseArgs(args)
	if isUsageErr(err) {
		parser.WriteHelp(os.Stdout)
		return nil
	}

	return err
}

func isUsageErr(err error) bool {
	if err != nil {
		var e *goflags.Error
		if errors.As(err, &e) {
			switch {
			case errors.Is(e.Type, goflags.ErrHelp), errors.Is(e.Type, goflags.ErrCommandRequired):
				return true
			}
		}
	}
	return false
}
