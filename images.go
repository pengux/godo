package godo

import "fmt"

// Image represents a Digitalocean image.
type Image struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Distribution string   `json:"distribution"`
	Slug         string   `json:"slug"`
	Public       bool     `json:"public"`
	RegionIDs    []int    `json:"regions"`
	RegionSlugs  []string `json:"region_slugs"`
}

// DeleteImage deletes an image. There is no way to restore a deleted image so be careful and ensure any data is properly backed up.
func (c *Client) DeleteImage(ID interface{}) error {
	var DOResp struct {
		Status  Status `json:"status"`
		Message string `json:"message"`
	}

	var s string
	switch ID.(type) {
	case string, int:
		s = fmt.Sprintf("/images/%v/destroy", ID)
	default:
		return fmt.Errorf("ID must be either a string or integer")
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return err
	}

	if DOResp.Status == StatusError {
		return fmt.Errorf("could not delete image with ID %v: %v", ID, DOResp.Message)
	}

	return nil
}

// GetAllImages returns all available images for the client ID.
func (c *Client) GetAllImages() ([]Image, error) {
	var DOResp struct {
		Status  Status  `json:"status"`
		Images  []Image `json:"images"`
		Message string  `json:"message"`
	}

	err := c.doGet("/images", &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get images: %v", DOResp.Message)
	}

	return DOResp.Images, nil
}

// GetImageByID returns information about an image by its ID, which can be either integer or string
func (c *Client) GetImageByID(ID interface{}) (*Image, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		Image   Image  `json:"image"`
		Message string `json:"message"`
	}

	var s string
	switch ID.(type) {
	case string, int:
		s = fmt.Sprintf("/images/%v", ID)
	default:
		return nil, fmt.Errorf("ID must be either a string or integer")
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get image with ID %v: %v", ID, DOResp.Message)
	}

	return &DOResp.Image, nil
}

// TransferImage transfers an image to a specified region. Returns an event ID on success.
func (c *Client) TransferImage(ID interface{}, regionID int) (int, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		EventID int    `json:"event_id"`
		Message string `json:"message"`
	}

	var s string
	switch ID.(type) {
	case string, int:
		s = fmt.Sprintf("/images/%v/transfer?region_id=%d", ID, regionID)
	default:
		return 0, fmt.Errorf("ID must be either a string or integer")
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return 0, err
	}

	if DOResp.Status == StatusError {
		return 0, fmt.Errorf("could not transfer image with ID %v: %v", ID, DOResp.Message)
	}

	return DOResp.EventID, nil
}
