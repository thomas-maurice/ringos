package ringo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Client struct {
	httpClient    *http.Client
	targetAddress string
	debug         bool
	username      string
	password      string
	authEnabled   bool
}

type ClientConfig struct {
	TargetAddress string
	Debug         bool
	Username      string
	Password      string
	AuthEnabled   bool
}

func NewClient(config *ClientConfig) (*Client, error) {
	return &Client{
		httpClient:    &http.Client{},
		targetAddress: config.TargetAddress,
		debug:         config.Debug,
		username:      config.Username,
		password:      config.Password,
		authEnabled:   config.AuthEnabled,
	}, nil
}

func (c *Client) URL(pathes ...string) string {
	return c.buildURL(pathes...)
}

func (c *Client) buildURL(pathes ...string) string {
	return fmt.Sprintf(
		c.targetAddress+"/%s",
		strings.TrimLeft(
			strings.Join(pathes, "/"), "/"),
	)
}

func (c *Client) Get(url string, ret interface{}) (int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	if c.debug {
		bytes, _ := httputil.DumpRequestOut(req, true)
		fmt.Println(string(bytes))
	}

	if c.authEnabled {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if c.debug {
		rbytes, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(rbytes))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if ret != nil {
		err := json.Unmarshal(bytes, ret)
		if err != nil {
			return 0, err
		}
	}

	return resp.StatusCode, nil
}

func (c *Client) Post(url string, data interface{}, ret interface{}) (int, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(dataBytes))
	if err != nil {
		return 0, err
	}

	req.Header.Add("Content-Type", "application/json")

	if c.debug {
		bytes, _ := httputil.DumpRequestOut(req, true)
		fmt.Println(string(bytes))
	}

	if c.authEnabled {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if c.debug {
		rbytes, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(rbytes))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if ret != nil {
		err := json.Unmarshal(bytes, ret)
		if err != nil {
			return 0, err
		}
	}

	return resp.StatusCode, nil
}

func (c *Client) Status() (*DeviceStatus, error) {
	var status DeviceStatus
	_, err := c.Get(c.buildURL("api", "status"), &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (c *Client) Colour() (*ColourResponse, error) {
	var colour ColourResponse
	_, err := c.Get(c.buildURL("api", "colour"), &colour)
	if err != nil {
		return nil, err
	}
	return &colour, nil
}

func (c *Client) GetConfig() (*Config, error) {
	var config Config
	_, err := c.Get(c.buildURL("api", "config"), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Client) Chase() (*ChaseResponse, error) {
	var chase ChaseResponse
	_, err := c.Get(c.buildURL("api", "chase"), &chase)
	if err != nil {
		return nil, err
	}
	return &chase, nil
}

func (c *Client) Breathing() (*BreathingResponse, error) {
	var breathing BreathingResponse
	_, err := c.Get(c.buildURL("api", "breathing"), &breathing)
	if err != nil {
		return nil, err
	}
	return &breathing, nil
}

func (c *Client) SetColour(colour *ColourRequest) (*GenericAnswer, error) {
	var response GenericAnswer
	_, err := c.Post(c.buildURL("api", "colour"), colour, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) SetChase(chase *ChaseRequest) (*GenericAnswer, error) {
	var response GenericAnswer
	_, err := c.Post(c.buildURL("api", "chase"), chase, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) SetBreathing(breathing *BreathingRequest) (*GenericAnswer, error) {
	var response GenericAnswer
	_, err := c.Post(c.buildURL("api", "breathing"), breathing, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) KnownNetworks() (KnownNetworkList, error) {
	var knownNetworks KnownNetworkList
	_, err := c.Get(c.buildURL("api", "net", "list"), &knownNetworks)
	if err != nil {
		return nil, err
	}
	return knownNetworks, nil
}

func (c *Client) ScanNetworks() (ScannedNetworks, error) {
	var scannedNetworks ScannedNetworks
	_, err := c.Get(c.buildURL("api", "net", "scan"), &scannedNetworks)
	if err != nil {
		return nil, err
	}
	return scannedNetworks, nil
}

func (c *Client) DeleteNetwork(req DeleteNetworkRequest) (*GenericAnswer, error) {
	var resp GenericAnswer
	_, err := c.Post(c.buildURL("api", "net", "del"), &req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) AddNetwork(req AddNetworkRequest) (*GenericAnswer, error) {
	var resp GenericAnswer
	_, err := c.Post(c.buildURL("api", "net", "add"), &req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateConfig(req Config) (*GenericAnswer, error) {
	var resp GenericAnswer
	_, err := c.Post(c.buildURL("api", "config"), &req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
