package bosh

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/gogobosh"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	ClientID     string
	ClientSecret string
	Environment  string
	client       *gogobosh.Client
}

func New(environment, clientID, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Environment:  environment,
	}
}

func (c *Client) GetAllVMs() ([]VM, error) {
	log.Debug("Getting all BOSH managed VMs")

	client, err := c.getOrCreateUnderlyingClient()
	if err != nil {
		return []VM{}, err
	}

	deployments, err := client.GetDeployments()
	if err != nil {
		return []VM{}, fmt.Errorf(
			"failed to get bosh deployments, this can happen because of incorrect login details: %w", err)
	}

	// TODO: This gives back a confusing error when given a bad client secret
	var r []VM
	for _, d := range deployments {
		log.Infof("Found deployment %s", d.Name)
		vms, err := client.GetDeploymentVMs(d.Name)
		if err != nil {
			return []VM{}, err
		}

		log.Infof("With %d BOSH managed VMs", len(vms))

		for _, vm := range vms {
			instanceName := vm.JobName + "/" + vm.ID
			r = append(r, VM{
				InstanceName:  instanceName,
				InstanceGroup: vm.JobName,
				Deployment:    d.Name,
				IPs:           vm.IPs,
				JobState:      vm.JobState,
			})
		}
	}
	return r, nil
}

func (c *Client) getOrCreateUnderlyingClient() (*gogobosh.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	config := &gogobosh.Config{
		BOSHAddress:       fmt.Sprintf("https://%s:25555", c.Environment),
		ClientID:          c.ClientID,
		ClientSecret:      c.ClientSecret,
		HttpClient:        http.DefaultClient,
		SkipSslValidation: true,
	}

	log.Debugf("Creating bosh client to connect to %s", config.BOSHAddress)
	client, err := gogobosh.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create bosh client: %w", err)
	}

	c.client = client
	return client, nil
}
