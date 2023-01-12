package ers

import (
	"fmt"
	"net/http"
)

type Config struct {
	Hostname string `json:"hostname,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Client struct {
	config Config
}

func (cli *Client) GetEP(epName string) (*EP, error) {
	uri := fmt.Sprintf("https://%v:%v/ers/config/endpoint/name/%v", cli.config.Hostname, cli.config.Port, epName)

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.SetBasicAuth(cli.config.Username, cli.config.Password)
	var ep *EP
	return ep, nil
}

func (cli *Client) GetEPGroup(name string) (*EPGroup, error) {
	var epg *EPGroup
	return epg, nil
}

func NewClient(config Config) *Client {
	return &Client{config: config}
}
