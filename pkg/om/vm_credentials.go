package om

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type vmCredential struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	Password string `json:"password"`
}

func (c *Client) VMCredentials(deployment, instanceGroup string) (string, error) {
	if credentials, ok := c.vmCredentials[deployment]; ok {
		if password, ok := credentials[instanceGroup]; ok {
			return password, nil
		} else {
			return "", fmt.Errorf("unable to find credential %s for deployment %s", instanceGroup, deployment)
		}

	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v0/deployed/products/%s/vm_credentials", deployment), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	output := []vmCredential{}
	err = json.Unmarshal(respBody, &output)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal response: %s", err)
	}
	credMap := make(map[string]string)
	for _, credential := range output {
		name := credential.Name
		credMap[strings.Split(name, "-")[0]] = credential.Password
	}
	c.vmCredentials[deployment] = credMap

	if credentials, ok := c.vmCredentials[deployment]; ok {
		if password, ok := credentials[instanceGroup]; ok {
			return password, nil
		} else {
			return "", fmt.Errorf("unable to find credential %s for deployment %s", instanceGroup, deployment)
		}
	} else {
		return "", fmt.Errorf("unable to find credential %s for deployment %s", instanceGroup, deployment)
	}
}
