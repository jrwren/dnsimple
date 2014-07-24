package dnsimple

import (
	"fmt"
	"strconv"
)

type RecordResponse struct {
	Record Record `json:"record"`
}

// Record is used to represent a retrieved Record. All properties
// are set as strings.
type Record struct {
	Name       string `json:"name"`
	Content    string `json:"content"`
	DomainId   int64  `json:"domain_id"`
	Id         int64  `json:"id"`
	Prio       int64  `json:"prio"`
	RecordType string `json:"record_type"`
	Ttl        int64  `json:"ttl"`
}

// Returns the domain id
func (r *Record) StringDomainId() string {
	return strconv.FormatInt(r.DomainId, 10)
}

// Returns the id
func (r *Record) StringId() string {
	return strconv.FormatInt(r.Id, 10)
}

// Returns the string for prio
func (r *Record) StringPrio() string {
	return strconv.FormatInt(r.Prio, 10)
}

// Returns the string for Locked
func (r *Record) StringTtl() string {
	return strconv.FormatInt(r.Ttl, 10)
}

// CreateRecord contains the request parameters to create a new
// record.
type CreateRecord struct {
	Name  string // name of the record
	Value string // where the record points
	Type  string // type, i.e a, mx
	Ttl   string // TTL of record
}

// CreateRecord creates a record from the parameters specified and
// returns an error if it fails. If no error and an ID is returned,
// the Record was succesfully created.
func (c *Client) CreateRecord(domain string, opts *CreateRecord) (string, error) {
	// Make the request parameters
	params := make(map[string]interface{})

	params["name"] = opts.Name
	params["record_type"] = opts.Type
	params["content"] = opts.Value

	if opts.Ttl != "" {
		ttl, err := strconv.ParseInt(opts.Ttl, 0, 0)
		if err != nil {
			return "", nil
		}
		params["ttl"] = ttl
	}

	endpoint := fmt.Sprintf("/domains/%s/records", domain)

	req, err := c.NewRequest(params, "POST", endpoint)
	if err != nil {
		return "", err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return "", fmt.Errorf("Error creating record: %s", err)
	}

	record := new(RecordResponse)

	err = decodeBody(resp, &record)

	if err != nil {
		return "", fmt.Errorf("Error parsing record response: %s", err)
	}

	// The request was successful
	return record.Record.StringId(), nil
}

// DestroyRecord destroys a record by the ID specified and
// returns an error if it fails. If no error is returned,
// the Record was succesfully destroyed.
func (c *Client) DestroyRecord(id string) error {
	var body map[string]interface{}
	req, err := c.NewRequest(body, "DELETE", fmt.Sprintf("/records/%s", id))

	if err != nil {
		return err
	}

	_, err = checkResp(c.Http.Do(req))

	if err != nil {
		return fmt.Errorf("Error destroying record: %s", err)
	}

	// The request was successful
	return nil
}

// RetrieveRecord gets  a record by the ID specified and
// returns a Record and an error. An error will be returned for failed
// requests with a nil Record.
func (c *Client) RetrieveRecord(id string) (Record, error) {
	var body map[string]interface{}
	req, err := c.NewRequest(body, "GET", fmt.Sprintf("/records/%s", id))

	if err != nil {
		return Record{}, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return Record{}, fmt.Errorf("Error destroying record: %s", err)
	}

	record := new(RecordResponse)

	err = decodeBody(resp, record)

	if err != nil {
		return Record{}, fmt.Errorf("Error decoding record response: %s", err)
	}

	// The request was successful
	return record.Record, nil
}