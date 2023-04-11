package metadata

import (
	"context"
)

type IPv4Data struct {
	Public  []IPNet `json:"public"`
	Private []IPNet `json:"private"`
	Elastic []IPNet `json:"elastic"`
}

type IPv6Data struct {
	Ranges        []IPNet `json:"ranges"`
	LinkLocal     IPNet   `json:"link-local"`
	ElasticRanges []IPNet `json:"elastic-ranges"`
}

type NetworkData struct {
	VLANID int      `json:"vlan-id"`
	IPv4   IPv4Data `json:"ipv4"`
	IPv6   IPv6Data `json:"ipv6"`
}

func (c *Client) GetNetwork(ctx context.Context) (*NetworkData, error) {
	req := c.R(ctx).SetResult(&NetworkData{})

	resp, err := req.Get("network")
	if err != nil {
		return nil, err
	}

	return resp.Result().(*NetworkData), nil
}
