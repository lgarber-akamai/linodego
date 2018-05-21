package golinode

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

// Image represents a deployable Image object for use with Linode Instances
type Image struct {
	CreatedStr  string `json:"created"`
	UpdatedStr  string `json:"updated"`
	ID          string
	Label       string
	Description string
	Type        string
	IsPublic    bool
	Size        int
	Vendor      string
	Deprecated  bool

	CreatedBy string     `json:"created_by"`
	Created   *time.Time `json:"-"`
	Updated   *time.Time `json:"-"`
}

func (l *Image) fixDates() *Image {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
}

// ImagesPagedResponse represents a linode API response for listing of images
type ImagesPagedResponse struct {
	*PageOptions
	Data []*Image
}

func (ImagesPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.Images.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *ImagesPagedResponse) AppendData(r *ImagesPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

func (ImagesPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(ImagesPagedResponse{})
}

// ListImages lists Images
func (c *Client) ListImages(opts *ListOptions) ([]*Image, error) {
	response := ImagesPagedResponse{}
	err := c.ListHelper(&response, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	if err != nil {
		return nil, err
	}
	return response.Data, nil

}

// GetImage gets the Image with the provided ID
func (c *Client) GetImage(id string) (*Image, error) {
	e, err := c.Images.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := coupleAPIErrors(c.Images.R().Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*Image), nil
}
