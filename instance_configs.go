package golinode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

type InstanceConfig struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID          int
	Label       string                   `json:"label"`
	Comments    string                   `json:"comments"`
	Devices     *InstanceConfigDeviceMap `json:"devices"`
	Helpers     *InstanceConfigHelpers   `json:"helpers"`
	MemoryLimit int                      `json:"memory_limit"`
	Kernel      string                   `json:"kernel"`
	InitRD      int                      `json:"init_rd"`
	RootDevice  string                   `json:"root_device"`
	RunLevel    string                   `json:"run_level"`
	VirtMode    string                   `json:"virt_mode"`
	Created     *time.Time               `json:"-"`
	Updated     *time.Time               `json:"-"`
}

type InstanceConfigDevice struct {
	DiskID   int `json:"disk_id"`
	VolumeID int `json:"volume_id"`
}

type InstanceConfigDeviceMap struct {
	SDA *InstanceConfigDevice `json:"sda"`
	SDB *InstanceConfigDevice `json:"sdb"`
	SDC *InstanceConfigDevice `json:"sdc"`
	SDD *InstanceConfigDevice `json:"sdd"`
	SDE *InstanceConfigDevice `json:"sde"`
	SDF *InstanceConfigDevice `json:"sdf"`
	SDG *InstanceConfigDevice `json:"sdg"`
	SDH *InstanceConfigDevice `json:"sdh"`
}

type InstanceConfigHelpers struct {
	UpdateDBDisabled  bool `json:"updatedb_disabled"`
	Distro            bool `json:"distro"`
	ModulesDep        bool `json:"modules_dep"`
	Network           bool `json:"network"`
	DevTmpFsAutomount bool `json:"devtmpfs_automount"`
}

// InstanceConfigsPagedResponse represents a paginated InstanceConfig API response
type InstanceConfigsPagedResponse struct {
	*PageOptions
	Data []*InstanceConfig
}

// InstanceConfigCreateOptions are InstanceConfig settings that can be used at creation
type InstanceConfigCreateOptions struct {
	Label       string                   `json:"label"`
	Comments    string                   `json:"comments"`
	Devices     *InstanceConfigDeviceMap `json:"devices"`
	Helpers     *InstanceConfigHelpers   `json:"helpers"`
	MemoryLimit int                      `json:"memory_limit"`
	Kernel      string                   `json:"kernel"`
	InitRD      int                      `json:"init_rd"`
	RootDevice  string                   `json:"root_device"`
	RunLevel    string                   `json:"run_level"`
	VirtMode    string                   `json:"virt_mode"`
}

// InstanceConfigUpdateOptions are InstanceConfig settings that can be used in updates
type InstanceConfigUpdateOptions InstanceConfigCreateOptions

func (i InstanceConfig) getCreateOptions() InstanceConfigCreateOptions {
	return InstanceConfigCreateOptions{
		Label:       i.Label,
		Comments:    i.Comments,
		Devices:     i.Devices,
		Helpers:     i.Helpers,
		MemoryLimit: i.MemoryLimit,
		Kernel:      i.Kernel,
		InitRD:      i.InitRD,
		RootDevice:  i.RootDevice,
		RunLevel:    i.RunLevel,
		VirtMode:    i.VirtMode,
	}
}

func (i InstanceConfig) getUpdateOptions() InstanceConfigUpdateOptions {
	return InstanceConfigUpdateOptions{
		Label:       i.Label,
		Comments:    i.Comments,
		Devices:     i.Devices,
		Helpers:     i.Helpers,
		MemoryLimit: i.MemoryLimit,
		Kernel:      i.Kernel,
		InitRD:      i.InitRD,
		RootDevice:  i.RootDevice,
		RunLevel:    i.RunLevel,
		VirtMode:    i.VirtMode,
	}
}

// EndpointWithID gets the endpoint URL for InstanceConfigs of a given Instance
func (InstanceConfigsPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.InstanceConfigs.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends InstanceConfigs when processing paginated InstanceConfig responses
func (resp *InstanceConfigsPagedResponse) AppendData(r *InstanceConfigsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of InstanceConfig
func (InstanceConfigsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(InstanceConfigsPagedResponse{})
}

// ListInstanceConfigs lists InstanceConfigs
func (c *Client) ListInstanceConfigs(linodeID int, opts *ListOptions) ([]*InstanceConfig, error) {
	response := InstanceConfigsPagedResponse{}
	err := c.ListHelperWithID(response, linodeID, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *InstanceConfig) fixDates() *InstanceConfig {
	v.Created, _ = parseDates(v.CreatedStr)
	v.Updated, _ = parseDates(v.UpdatedStr)
	return v
}

// GetInstanceConfig gets the template with the provided ID
func (c *Client) GetInstanceConfig(linodeID int, configID int) (*InstanceConfig, error) {
	e, err := c.InstanceConfigs.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, configID)
	r, err := c.R().SetResult(&InstanceConfig{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceConfig).fixDates(), nil
}

// CreateInstanceConfig creates a new InstanceConfig for the given Instance
func (c *Client) CreateInstanceConfig(linodeID int, createOpts InstanceConfigCreateOptions) (*InstanceConfig, error) {
	var body string
	e, err := c.InstanceConfigs.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}

	req := c.R().SetResult(&InstanceConfig{})

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

	return r.Result().(*InstanceConfig), nil
}

// UpdateInstanceConfig update an InstanceConfig for the given Instance
func (c *Client) UpdateInstanceConfig(linodeID int, configID int, updateOpts InstanceConfigUpdateOptions) (*InstanceConfig, error) {
	var body string
	e, err := c.InstanceConfigs.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, configID)
	req := c.R().SetResult(&InstanceConfig{})

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

	return r.Result().(*InstanceConfig), nil
}

// RenameInstanceConfig renames an InstanceConfig
func (c *Client) RenameInstanceConfig(linodeID int, configID int, label string) (*InstanceConfig, error) {
	return c.UpdateInstanceConfig(linodeID, configID, InstanceConfigUpdateOptions{Label: label})
}

// DeleteInstanceConfig deletes a Linode InstanceConfig
func (c *Client) DeleteInstanceConfig(id int) error {
	e, err := c.InstanceConfigs.EndpointWithID(id)
	if err != nil {
		return err
	}

	if _, err = c.R().Delete(e); err != nil {
		return err
	}

	return nil
}
