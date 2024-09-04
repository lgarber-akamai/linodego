package linodego

import (
	"context"
)

// NetworkTransferPrice represents a single valid network transfer price.
type NetworkTransferPrice struct {
	baseType
}

// ListNetworkTransferPrices lists network transfer prices. This endpoint is cached by default.
func (c *Client) ListNetworkTransferPrices(ctx context.Context, opts *ListOptions) ([]NetworkTransferPrice, error) {
	e := "network-transfer/prices"

	endpoint, err := generateListCacheURL(e, opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]NetworkTransferPrice), nil
	}

	response, err := getPaginatedResults[NetworkTransferPrice](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}
