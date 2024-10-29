package ssh_test

import (
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu-labs/opsman-utils/pkg/ssh"
	"os"
	"strings"
	"testing"
)

func TestBOSHAllProxyConnectionIntegration(t *testing.T) {
	boshVMAddress := os.Getenv("BOSH_VM_ADDRESS")
	boshVMPassword := os.Getenv("BOSH_VM_PASSWORD")
	boshAllProxy := os.Getenv("BOSH_ALL_PROXY")

	if boshVMAddress == "" || boshVMPassword == "" || boshAllProxy == "" {
		t.Skip("SSH integration tests require BOSH_VM_ADDRESS, BOSH_VM_PASSWORD, and BOSH_ALL_PROXY env vars to run")
	}

	if !strings.Contains(boshVMAddress, ":") {
		boshVMAddress += ":22"
	}

	conn := ssh.New()
	out, err := conn.Execute(boshVMAddress, "vcap", boshVMPassword, "echo 'hello world'")
	require.NoError(t, err)
	require.Equal(t, "hello world\r\n", string(out))
}
