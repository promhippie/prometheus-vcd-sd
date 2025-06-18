package client

import (
	"net/url"

	"github.com/vmware/go-vcloud-director/v3/govcd"
)

// New initializes a new client.
func New(endpoint *url.URL, insecure bool, username, password, org, vdc string) *Client {
	return &Client{
		Upstream: govcd.NewVCDClient(
			*endpoint,
			insecure,
		),
		Endpoint:     endpoint,
		Insecure:     insecure,
		Username:     username,
		Password:     password,
		Organization: org,
		Datacenter:   vdc,
	}
}

// Client abstracts some cloud provider client handling.
type Client struct {
	Upstream *govcd.VCDClient

	Endpoint     *url.URL
	Insecure     bool
	Username     string
	Password     string
	Organization string
	Datacenter   string
}

// Authenticate wraps the auth for the cloud provider.
func (c *Client) Authenticate() error {
	return c.Upstream.Authenticate(
		c.Username,
		c.Password,
		c.Organization,
	)
}

// Disconnect wraps the logout for the cloud provider.
func (c *Client) Disconnect() error {
	return c.Upstream.Disconnect()
}
