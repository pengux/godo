package godo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Status is the response status from API after each request
type Status string

const (
	// StatusOK indicates that the request could be executed successfully
	StatusOK Status = "OK"
	// StatusError indicates that there was an error while processing the request, more information about the error should be available in the "message" field of the response
	StatusError Status = "ERROR"

	// APIURL is the URL for Digitalocean's API
	APIURL = "https://api.digitalocean.com/v1"
)

// Client represents a new client which sends request to the API
type Client struct {
	ClientID string
	APIKey   string
}

// Event represents a event at DigitalOcean
type Event struct {
	ID           string  `json:"id"`
	ActionStatus string  `json:"action_status"`
	DropletID    int     `json:"droplet_id"`
	EventTypeID  int     `json:"event_type_id"`
	Percentage   float64 `json:"percentage"`
}

// Region represent available regions within DigitalOcean cloud
type Region struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Size represents a droplet size
type Size struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Memory       int     `json:"memory"`
	CPU          int     `json:"cpu"`
	Disk         int     `json:"disk"`
	CostPerHour  float64 `json:"cost_per_hour"`
	CostPerMonth string  `json:"cost_per_month"`
}

// NewClient returns a new Client struct
func NewClient(clientID string, apiKey string) *Client {
	return &Client{
		clientID,
		apiKey,
	}
}

// GetEventByID returns information about an event by its ID
func (c *Client) GetEventByID(ID int) (*Event, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		Event   Event  `json:"event"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/events/%d", ID), &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get event with ID %d: %v", ID, DOResp.Message)
	}

	return &DOResp.Event, nil
}

// GetAllRegions returns all available regions
func (c *Client) GetAllRegions() ([]Region, error) {
	var DOResp struct {
		Status  Status   `json:"status"`
		Regions []Region `json:"regions"`
		Message string   `json:"message"`
	}

	err := c.doGet("/regions", &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get regions: %v", DOResp.Message)
	}

	return DOResp.Regions, nil
}

// GetAllSizes returns all available sizes for a droplet
func (c *Client) GetAllSizes() ([]Size, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		Sizes   []Size `json:"sizes"`
		Message string `json:"message"`
	}

	err := c.doGet("/sizes", &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get sizes: %v", DOResp.Message)
	}

	return DOResp.Sizes, nil
}

func (c *Client) doGet(endpoint string, i interface{}) error {
	url := fmt.Sprintf("%s%s", APIURL, endpoint)

	if !strings.Contains(url, "?") {
		url += "?"
	} else {
		url += "&"
	}
	url += fmt.Sprintf("client_id=%s&api_key=%s", c.ClientID, c.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		panic(err)
	}

	return nil
}
