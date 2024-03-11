package main

import (
	"fmt"
	"os"

	goflags "github.com/jessevdk/go-flags"
	"github.com/vmware-tanzu-labs/bosh-agent-ip-changer/command"
)

type CommandHolder struct {
	Version  command.VersionCommand `command:"version" description:"Print version information and exit"`
	ChangeIP command.ChangeIP       `command:"changeIP" description:"Changes the bosh agent configuration when bosh IP address changes"`
}

var Command CommandHolder

func main() {
	parser := goflags.NewParser(&Command, goflags.HelpFlag)
	parser.NamespaceDelimiter = "-"

	_, err := parser.Parse()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
