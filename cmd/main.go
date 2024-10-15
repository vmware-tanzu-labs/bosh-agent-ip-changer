package main

import (
	"fmt"
	"github.com/vmware-tanzu-labs/opsman-utils/command"
	"os"
)

func main() {
	err := command.Main(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
