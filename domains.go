package godo

import (
	"fmt"
	"net"
)

// Domain maps to the domain(s) field in the response
type Domain struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	TTL               int    `json:"ttl"`
	LiveZoneFile      string `json:"live_zone_file"`
	Error             string `json:"error"`
	ZoneFileWithError string `json:"zone_file_with_error"`
}

// PartialDomain maps to the partial domain data in the response when a new domain is created successfully
type PartialDomain struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DomainRecord maps to the domain record response
type DomainRecord struct {
	ID         int    `json:"id"`
	DomainID   int    `json:"domain_id"`
	RecordType string `json:"record_type"`
	Name       string `json:"name"`
	Data       string `json:"data"`
	Priority   int    `json:"priority"`
	Port       int    `json:"port"`
	Weight     int    `json:"weight"`
}

// CreateDomain creates a new domain
func (c *Client) CreateDomain(name string, IP net.IP) (*PartialDomain, error) {
	// Validate
	if name == "" {
		return nil, fmt.Errorf("name must be set")
	}

	if len(IP) == 0 {
		return nil, fmt.Errorf("IP address must be set and valid")
	}

	s := fmt.Sprintf("/domains/new?name=%s&ip_address=%s", name, IP)

	var DOResp struct {
		Status  Status        `json:"status"`
		Domain  PartialDomain `json:"domain"`
		Message string        `json:"message"`
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not create domain: %v", DOResp.Message)
	}

	return &DOResp.Domain, nil
}

// DeleteDomainByID returns a domain by its ID
func (c *Client) DeleteDomainByID(ID interface{}) error {
	var DOResp struct {
		Status  Status `json:"status"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/domains/%v/destroy", ID), &DOResp)
	if err != nil {
		return err
	}

	if DOResp.Status == StatusError {
		return fmt.Errorf("could not delete domain with ID %v: %v", ID, DOResp.Message)
	}

	return nil
}

// GetAllDomains returns all current domain
func (c *Client) GetAllDomains() ([]Domain, error) {
	var DOResp struct {
		Status  Status   `json:"status"`
		Domains []Domain `json:"domains"`
		Message string   `json:"message"`
	}

	err := c.doGet("/domains", &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get domains: %v", DOResp.Message)
	}

	return DOResp.Domains, nil
}

// GetDomainByID returns a domain by its ID
func (c *Client) GetDomainByID(ID int) (*Domain, error) {
	var DOResp struct {
		Status  Status `json:"status"`
		Domain  Domain `json:"domain"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/domains/%d", ID), &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get domain with ID %d: %v", ID, DOResp.Message)
	}

	return &DOResp.Domain, nil
}

// CreateDomainRecord creates a record for a domain by ID, if sucessfully it will returns a new DomainRecord
func (c *Client) CreateDomainRecord(ID interface{}, r DomainRecord) (*DomainRecord, error) {
	// Validate
	if r.RecordType == "" {
		return nil, fmt.Errorf("record type must be set")
	}

	if r.Data == "" {
		return nil, fmt.Errorf("data value must be set")
	}

	s := fmt.Sprintf("/domains/%v/records/new?record_type=%s&data=%s", ID, r.RecordType, r.Data)

	if r.Name != "" {
		s += fmt.Sprintf("&name=%s", r.Name)
	}

	if r.Priority != 0 {
		s += fmt.Sprintf("&priority=%d", r.Priority)
	}

	if r.Port != 0 {
		s += fmt.Sprintf("&port=%d", r.Port)
	}

	if r.Weight != 0 {
		s += fmt.Sprintf("&weight=%d", r.Weight)
	}

	var DOResp struct {
		Status  Status       `json:"status"`
		Record  DomainRecord `json:"record"`
		Message string       `json:"message"`
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not create record for domain %v: %v", ID, DOResp.Message)
	}

	return &DOResp.Record, nil
}

// GetAllRecordsByDomain returns all current domain records for a specific domain. The domainID can be integer or string
func (c *Client) GetAllRecordsByDomain(domainID interface{}) ([]DomainRecord, error) {
	var DOResp struct {
		Status  Status         `json:"status"`
		Records []DomainRecord `json:"records"`
		Message string         `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/domains/%v/records", domainID), &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get records for domain %v: %v", domainID, DOResp.Message)
	}

	return DOResp.Records, nil
}

// GetRecordByDomain return a domain record by domain ID and record ID. domainID can be integer or string
func (c *Client) GetRecordByDomain(domainID interface{}, ID int) (*DomainRecord, error) {
	var DOResp struct {
		Status  Status       `json:"status"`
		Record  DomainRecord `json:"record"`
		Message string       `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/domains/%v/records/%d", domainID, ID), &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not get record for domain %v with ID %d: %v", domainID, ID, DOResp.Message)
	}

	return &DOResp.Record, nil
}

// UpdateRecordByDomain updates a domain record by domain ID and record ID. domainID can be integer or string
func (c *Client) UpdateRecordByDomain(domainID interface{}, r DomainRecord) (*DomainRecord, error) {
	// Validate
	if r.ID == 0 {
		return nil, fmt.Errorf("record ID must be set")
	}

	if r.RecordType == "" {
		return nil, fmt.Errorf("record type must be set")
	}

	if r.Data == "" {
		return nil, fmt.Errorf("data value must be set")
	}

	s := fmt.Sprintf("/domains/%v/records/new?record_type=%s&data=%s", domainID, r.ID, r.RecordType, r.Data)

	if r.Name != "" {
		s += fmt.Sprintf("&name=%s", r.Name)
	}

	if r.Priority != 0 {
		s += fmt.Sprintf("&priority=%d", r.Priority)
	}

	if r.Port != 0 {
		s += fmt.Sprintf("&port=%d", r.Port)
	}

	if r.Weight != 0 {
		s += fmt.Sprintf("&weight=%d", r.Weight)
	}

	var DOResp struct {
		Status  Status       `json:"status"`
		Record  DomainRecord `json:"record"`
		Message string       `json:"message"`
	}

	err := c.doGet(s, &DOResp)
	if err != nil {
		return nil, err
	}

	if DOResp.Status == StatusError {
		return nil, fmt.Errorf("could not create record %d for domain %v: %v", r.ID, domainID, DOResp.Message)
	}

	return &DOResp.Record, nil
}

// DeleteRecordByDomain delete a domain record
func (c *Client) DeleteRecordByDomain(domainID interface{}, ID int) error {
	var DOResp struct {
		Status  Status `json:"status"`
		Message string `json:"message"`
	}

	err := c.doGet(fmt.Sprintf("/domains/%v/records/%d/destroy", domainID, ID), &DOResp)
	if err != nil {
		return err
	}

	if DOResp.Status == StatusError {
		return fmt.Errorf("could not delete record %d for domain with ID %v: %v", domainID, ID, DOResp.Message)
	}

	return nil
}
