package main

import (
	"fmt"
	"os"

	goflags "github.com/jessevdk/go-flags"
	"github.com/vmware-tanzu-labs/opsman-utils/command"
)

type CommandHolder struct {
	Version      command.VersionCommand `command:"version" description:"Print version information and exit"`
	UpdateBoshIP command.UpdateBoshIP   `command:"update-bosh-ip" description:"Updates every bosh agent to use the specified BOSH Director IP"`
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
