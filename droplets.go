package godo

import (
	"fmt"
	"strings"
	"time"
)

const (
	// EndpointDroplets is the endpoint string for droplets
	EndpointDroplets = "/droplets"
)

// Droplet maps to the droplet(s) field in the response
type Droplet struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	ImageID          int       `json:"image_id"`
	SizeID           int       `json:"size_id"`
	RegionID         int       `json:"region_id"`
	BackupsActive    bool      `json:"backups_active"`
	IPAdress         string    `json:"ip_address"`
	PrivateIPAddress string    `json:"private_ip_address"`
	Locked           bool      `json:"locked"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

// NewDroplet maps to the data that is required to create a new droplet
type NewDroplet struct {
	// Name is required
	Name string

	// Either SizeID or SizeSlug must be set
	SizeID   int
	SizeSlug string

	// Either IamgeID or ImageSlug must be set
	ImageID   int
	ImageSlug string

	// Either RegionID or RegionSlug must be set
	RegionID   int
	RegionSlug string

	SSHKeyIDs         []string
	PrivateNetworking bool
	BackupsEnabled    bool
}

// PartialDroplet maps to the partial droplet data in the response when a new droplet is created successfully
type PartialDroplet struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	ImageID int    `json:"image_id"`
	SizeID  int    `json:"size_id"`
	EventID int    `json:"event_id"`
}

// CreateDroplet creates a new droplet
func (c *Client) CreateDroplet(n NewDroplet) (*PartialDroplet, error) {
	// Validate
	if n.SizeID == 0 && n.SizeSlug == "" {
		return nil, fmt.Errorf("size ID or slug must be set")
	}

	if n.ImageID == 0 && n.ImageSlug == "" {
		return nil, fmt.Errorf("image ID or slug must be set")
	}

	if n.RegionID == 0 && n.RegionSlug == "" {
		return nil, fmt.Errorf("region ID or slug must be set")
	}

	s := fmt.Sprintf("/droplets/new?name=%s", n.Name)

	if n.SizeID != 0 {
		s += fmt.Sprintf("&size_id=%d", n.SizeID)
	} else {
		s += "&size_slug=" + n.SizeSlug
	}

	if n.ImageID != 0 {
		s += fmt.Sprintf("&image_id=%d", n.ImageID)
	} else {
		s += "&image_slug=" + n.ImageSlug
	}

	if n.RegionID != 0 {
		s += fmt.Sprintf("&region_id=%d", n.RegionID)
	} else {
		s += "&region_slug=" + n.RegionSlug
	}

	if len(n.SSHKeyIDs) > 0 {
		s += "&ssh_key_ids=" + strings.Join(n.SSHKeyIDs, ",")
	}

	if n.PrivateNetworking {
		s += "&private_networking=true"
	}

	if n.BackupsEnabled {
		s += "&backups_enabled=true"
	}

	var DOResp struct {
		Status  Status         `json:"status"`
		Droplet PartialDroplet `json:"droplet"`
		Message string         `json:"message"`
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not create droplet: %v", DOResp.Message)
	}

	return &DOResp.Droplet, nil
}

// DeleteDropletByID returns a domain by its ID. Returns an event ID on success
func (c *Client) DeleteDropletByID(ID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/destroy", ID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not delete droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// GetAllDroplets returns all active droplets
func (c *Client) GetAllDroplets() ([]Droplet, error) {
	var DOResp struct {
		Status   Status    `json:"status"`
		Droplets []Droplet `json:"droplets"`
		Message  string    `json:"message"`
	}

	err := c.doGet("/droplets", &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get droplets: %v", DOResp.Message)
	}

	return DOResp.Droplets, nil
}

// GetDropletByID returns a domain by its ID
func (c *Client) GetDropletByID(ID int) (*Droplet, error) {
	var DOResp struct {
		Status  Status  `json:"status"`
		Droplet Droplet `json:"droplet"`
		Message string  `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d", ID), &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get droplet with ID %d: %v", ID, DOResp.Message)
	}

	return &DOResp.Droplet, nil
}

// RebootDroplet reboot a droplet. This is the preferred method to use if a server is not responding. Returns an event ID on success.
func (c *Client) RebootDroplet(ID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/reboot", ID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not reboot droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// PowerCycleDroplet power cycle a droplet. This will turn off the droplet and then turn it back on. Returns an event ID on success.
func (c *Client) PowerCycleDroplet(ID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/power_cycle", ID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not reboot droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// ShutDownDroplet shut down a running droplet. This will turn off the droplet but it will remain in client's account. Returns an event ID on success.
func (c *Client) ShutDownDroplet(ID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/shutdown", ID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not shut down droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// PowerOffDroplet power off a running droplet. The droplet will remain in client's account. Returns an event ID on success.
func (c *Client) PowerOffDroplet(ID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/power_off", ID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not power off droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// PowerOnDroplet power on a powered off droplet. Returns an event ID on success.
func (c *Client) PowerOnDroplet(ID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/power_on", ID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not power on droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// ResetRootPassDroplet reset root's password for a droplet. Please be aware that this will reboot the droplet to allow resetting the password. Returns an event ID on success.
func (c *Client) ResetRootPassDroplet(ID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/password_reset", ID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not reset root's password for droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// ResizeDroplet resizes a droplet to a different size. The size param can be either string or integer. Returns an event ID on success.
func (c *Client) ResizeDroplet(ID int, size interface{}) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	s := fmt.Sprintf("/droplets/%d/resize", ID)

	switch size.(type) {
	case string:
		s += fmt.Sprintf("&size_slug=%s", size)
	case int:
		s += fmt.Sprintf("&size_id=%d", size)
	default:
		return 0, fmt.Errorf("size must be either a string or integer")
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not resize the droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// TakeSnapshotOnDroplet takes a snapshot of the droplet once it has been powered off, which can later be restored or used to create a new droplet from the same image. Please be aware this may cause a reboot. If name is an empty string, it will default to date/time. Returns an event ID on success.
func (c *Client) TakeSnapshotOnDroplet(ID int, name string) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	s := fmt.Sprintf("/droplets/%d/snapshot", ID)

	if name != "" {
		s += fmt.Sprintf("&name=%s", name)
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not take snapshot of droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// RestoreDroplet restores a droplet from a previous image or snapshot. This will be a mirror copy of the image or snapshot to the droplet. Returns an event ID on success.
func (c *Client) RestoreDroplet(ID, imageID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/restore?image_id=%d", ID, imageID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not restore droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// RebuildDroplet reinstalls a droplet with a default image. This is useful if you want to start again but retain the same IP address for your droplet. Returns an event ID on success.
func (c *Client) RebuildDroplet(ID, imageID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/rebuild?image_id=%d", ID, imageID), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not rebuild droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}

// RenameDroplet renames a droplet. Returns an event ID on success.
func (c *Client) RenameDroplet(ID int, name string) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/droplets/%d/rename?name=%s", ID, name), &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not rename droplet with ID %d: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}
