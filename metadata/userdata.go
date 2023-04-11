package metadata

import (
	"context"
	"github.com/linode/linodego"
)

func (c *Client) GetUserData(ctx context.Context) (string, error) {
	// Getting user-data requires the text/plain content type
	req := c.R(ctx).
		ExpectContentType("text/plain").
		SetHeader("Content-Type", "text/plain")

	resp, err := linodego.CoupleAPIErrors(req.Get("user-data"))
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
