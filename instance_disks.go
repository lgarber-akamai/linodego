package golinode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

type InstanceDisk struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID         int
	Label      string
	Status     string
	Size       int
	Filesystem string
	Created    *time.Time `json:"-"`
	Updated    *time.Time `json:"-"`
}

// InstanceDisksPagedResponse represents a paginated InstanceDisk API response
type InstanceDisksPagedResponse struct {
	*PageOptions
	Data []*InstanceDisk
}

// InstanceDiskCreateOptions are InstanceDisk settings that can be used at creation
type InstanceDiskCreateOptions struct {
	Label string `json:"label"`
	Size  int    `json:"size"`

	// Image is optional, but requires RootPass if provided
	Image    string `json:"image,omitempty"`
	RootPass string `json:"root_pass,omitempty"`

	Filesystem      string            `json:"filesystem,omitempty"`
	AuthorizedKeys  []string          `json:"authorized_keys,omitempty"`
	ReadOnly        bool              `json:"read_only,omitempty"`
	StackscriptID   int               `json:"stackscript_id,omitempty"`
	StackscriptData map[string]string `json:"stackscript_data,omitempty"`
}

// InstanceDiskUpdateOptions are InstanceDisk settings that can be used in updates
type InstanceDiskUpdateOptions struct {
	Label    string `json:"label"`
	ReadOnly bool   `json:"read_only,omitempty"`
}

// EndpointWithID gets the endpoint URL for InstanceDisks of a given Instance
func (InstanceDisksPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.InstanceDisks.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends InstanceDisks when processing paginated InstanceDisk responses
func (resp *InstanceDisksPagedResponse) AppendData(r *InstanceDisksPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of InstanceDisk
func (InstanceDisksPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(InstanceDisksPagedResponse{})
}

// ListInstanceDisks lists InstanceDisks
func (c *Client) ListInstanceDisks(linodeID int, opts *ListOptions) ([]*InstanceDisk, error) {
	response := InstanceDisksPagedResponse{}
	err := c.ListHelperWithID(response, linodeID, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *InstanceDisk) fixDates() *InstanceDisk {
	v.Created, _ = parseDates(v.CreatedStr)
	v.Updated, _ = parseDates(v.UpdatedStr)
	return v
}

// GetInstanceDisk gets the template with the provided ID
func (c *Client) GetInstanceDisk(linodeID int, configID int) (*InstanceDisk, error) {
	e, err := c.InstanceDisks.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, configID)
	r, err := c.R().SetResult(&InstanceDisk{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceDisk).fixDates(), nil
}

// CreateInstanceDisk creates a new InstanceDisk for the given Instance
func (c *Client) CreateInstanceDisk(linodeID int, createOpts InstanceDiskCreateOptions) (*InstanceDisk, error) {
	var body string
	e, err := c.InstanceDisks.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}

	req := c.R().SetResult(&InstanceDisk{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, err
	}

	r, err := req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	if err != nil {
		return nil, err
	}

	return r.Result().(*InstanceDisk), nil
}

// UpdateInstanceDisk creates a new InstanceDisk for the given Instance
func (c *Client) UpdateInstanceDisk(linodeID int, diskID int, updateOpts InstanceDiskUpdateOptions) (*InstanceDisk, error) {
	var body string
	e, err := c.InstanceDisks.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, diskID)

	req := c.R().SetResult(&InstanceDisk{})

	if bodyData, err := json.Marshal(updateOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, err
	}

	r, err := req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(e)

	if err != nil {
		return nil, err
	}

	return r.Result().(*InstanceDisk), nil
}

// RenameInstanceDisk renames an InstanceDisk
func (c *Client) RenameInstanceDisk(linodeID int, diskID int, label string) (*InstanceDisk, error) {
	return c.UpdateInstanceDisk(linodeID, diskID, InstanceDiskUpdateOptions{Label: label})
}

// DeleteInstanceDisk deletes a Linode InstanceDisk
func (c *Client) DeleteInstanceDisk(id int) error {
	e, err := c.InstanceDisks.EndpointWithID(id)
	if err != nil {
		return err
	}

	if _, err = c.R().Delete(e); err != nil {
		return err
	}

	return nil
}
