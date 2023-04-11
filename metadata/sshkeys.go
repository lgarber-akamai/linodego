package metadata

import (
	"context"
	"github.com/linode/linodego"
)

type SSHKeysUserData struct {
	Root []string `json:"root"`
}

type SSHKeysData struct {
	Users SSHKeysUserData `json:"users"`
}

func (c *Client) GetSSHKeys(ctx context.Context) (*SSHKeysData, error) {
	req := c.R(ctx).SetResult(&SSHKeysData{})

	resp, err := linodego.CoupleAPIErrors(req.Get("ssh-keys"))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*SSHKeysData), nil
}
