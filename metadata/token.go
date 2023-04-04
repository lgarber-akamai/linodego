package metadata

import (
	"context"
	"strconv"
)

type GenerateTokenOptions struct {
	ExpirySeconds int
}

func (c *Client) GenerateToken(ctx context.Context, opts GenerateTokenOptions) (string, error) {
	req := c.R(ctx)

	tokenExpirySeconds := 3600
	if opts.ExpirySeconds != 0 {
		tokenExpirySeconds = opts.ExpirySeconds
	}

	req.SetHeader("X-Metadata-Token-Expiry-Seconds", strconv.Itoa(tokenExpirySeconds))

	resp, err := req.Put("token")
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
