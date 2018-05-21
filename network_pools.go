package golinode

import (
	"fmt"

	"github.com/go-resty/resty"
)

// IPv6PoolsPagedResponse represents a paginated IPv6Pool API response
type IPv6PoolsPagedResponse struct {
	*PageOptions
	Data []*IPv6Range
}

// Endpoint gets the endpoint URL for IPv6Pool
func (IPv6PoolsPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.IPv6Pools.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends IPv6Pools when processing paginated IPv6Pool responses
func (resp *IPv6PoolsPagedResponse) AppendData(r *IPv6PoolsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of IPv6Pool
func (IPv6PoolsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(IPv6PoolsPagedResponse{})
}

// ListIPv6Pools lists IPv6Pools
func (c *Client) ListIPv6Pools(opts *ListOptions) ([]*IPv6Range, error) {
	response := IPv6PoolsPagedResponse{}
	err := c.ListHelper(&response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetIPv6Pool gets the template with the provided ID
func (c *Client) GetIPv6Pool(id string) (*IPv6Range, error) {
	e, err := c.IPv6Pools.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().SetResult(&IPv6Range{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*IPv6Range), nil
}
