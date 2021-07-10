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
}

type ClientConfig struct {
	TargetAddress string
}

func NewClient(config *ClientConfig) (*Client, error) {
	return &Client{
		httpClient:    &http.Client{},
		targetAddress: config.TargetAddress,
	}, nil
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
