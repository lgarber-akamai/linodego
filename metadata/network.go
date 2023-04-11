package metadata

import (
	"context"
	"net/netip"
)

type IPv4Data struct {
	Public  []netip.Prefix `json:"public"`
	Private []netip.Prefix `json:"private"`
	Elastic []netip.Prefix `json:"elastic"`
}

type IPv6Data struct {
	Ranges        []netip.Prefix `json:"ranges"`
	LinkLocal     netip.Prefix   `json:"link-local"`
	ElasticRanges []netip.Prefix `json:"elastic-ranges"`
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
