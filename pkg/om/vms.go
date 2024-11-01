package om

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vmware-tanzu-labs/opsman-utils/pkg/bosh"
)

type credential struct {
	Credential string `json:"credential"`
}

func (c *Client) VMs() (map[string][]bosh.VM, error) {
	result := make(map[string][]bosh.VM)

	req, err := http.NewRequest("GET", "/api/v0/deployed/director/credentials/bosh_commandline_credentials", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := credential{}
	err = json.Unmarshal(respBody, &output)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal director credentials response: %s", err)
	}
	keyValues := parseKeyValues(output.Credential)

	clientID := keyValues["BOSH_CLIENT"]
	clientSecret := keyValues["BOSH_CLIENT_SECRET"]
	environment := keyValues["BOSH_ENVIRONMENT"]
	boshClient := bosh.New(environment, clientID, clientSecret)
	vms, err := boshClient.GetAllVMs()
	if err != nil {
		return nil, err
	}

	for _, vm := range vms {
		if list, ok := result[vm.Deployment]; ok {
			list = append(list, vm)
			result[vm.Deployment] = list
		} else {
			result[vm.Deployment] = []bosh.VM{vm}
		}
	}
	return result, nil

}

func parseKeyValues(credentials string) map[string]string {
	values := make(map[string]string)
	kvs := strings.Split(credentials, " ")
	for _, kv := range kvs {
		if strings.Contains(kv, "=") {
			k := strings.Split(kv, "=")[0]
			v := strings.Split(kv, "=")[1]
			values[k] = v
		}
	}
	return values
}
