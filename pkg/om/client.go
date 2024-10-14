package om

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cloudfoundry-community/go-uaa"
	"golang.org/x/oauth2"
)

type Client struct {
	caCert             string
	clientID           string
	clientSecret       string
	insecureSkipVerify bool
	password           string
	target             string
	token              *oauth2.Token
	username           string
	connectTimeout     time.Duration
	requestTimeout     time.Duration
	vmCredentials      map[string]map[string]string
}

func New(target, username, password string,
	clientID, clientSecret string,
	insecureSkipVerify bool,
	caCert string) *Client {
	return &Client{
		caCert:             caCert,
		clientID:           clientID,
		clientSecret:       clientSecret,
		insecureSkipVerify: insecureSkipVerify,
		password:           password,
		target:             target,
		username:           username,
		connectTimeout:     time.Duration(10) * time.Second,
		requestTimeout:     time.Duration(1800) * time.Second,
		vmCredentials:      make(map[string]map[string]string),
	}
}

func (c *Client) do(request *http.Request) (*http.Response, error) {
	token := c.token
	target := c.target

	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "https://" + target
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("could not parse target url: %s", err)
	}

	targetURL.Path = "/uaa"

	request.URL.Scheme = targetURL.Scheme
	request.URL.Host = targetURL.Host

	client, err := newHTTPClient(
		c.insecureSkipVerify,
		c.caCert,
		c.requestTimeout,
		c.connectTimeout,
	)

	if err != nil {
		return nil, err
	}

	if token != nil && token.Valid() {
		request.Header.Set(
			"Authorization",
			fmt.Sprintf("Bearer %s", token.AccessToken),
		)
		return client.Do(request)
	}

	options := []uaa.Option{
		uaa.WithSkipSSLValidation(c.insecureSkipVerify),
		uaa.WithClient(client),
	}

	var authOption uaa.AuthenticationOption

	if c.username != "" && c.password != "" {
		authOption = uaa.WithPasswordCredentials(
			"opsman",
			"",
			c.username,
			c.password,
			uaa.JSONWebToken,
		)
	} else {
		authOption = uaa.WithClientCredentials(
			c.clientID,
			c.clientSecret,
			uaa.JSONWebToken,
		)
	}

	api, err := uaa.New(
		targetURL.String(),
		authOption,
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("could not init UAA client: %w", err)
	}

	for i := 0; i <= 2; i++ {
		token, err = api.Token(request.Context())
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("token could not be retrieved from target url: %w", err)
	}

	request.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", token.AccessToken),
	)

	c.token = token

	return client.Do(request)
}

func newHTTPClient(insecureSkipVerify bool, caCert string, requestTimeout time.Duration, connectTimeout time.Duration) (*http.Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
	}
	err := setCACert(caCert, tlsConfig)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: tlsConfig,
			Dial: (&net.Dialer{
				Timeout:   connectTimeout,
				KeepAlive: 30 * time.Second,
			}).Dial,
		},
		Timeout: requestTimeout,
	}, nil
}

func setCACert(caCert string, tlsConfig *tls.Config) error {
	if caCert == "" {
		return nil
	}

	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		caCertPool = x509.NewCertPool()
	}
	if !strings.Contains(caCert, "BEGIN") {
		contents, err := os.ReadFile(caCert)
		if err != nil {
			return fmt.Errorf("could not load ca cert from file: %s", err)
		}
		caCert = string(contents)
	}
	if ok := caCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
		return errors.New("could not use ca cert")
	}

	tlsConfig.RootCAs = caCertPool
	return nil
}
