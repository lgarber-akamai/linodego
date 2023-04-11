package linodego

import (
	"context"
)

type GrantsListResponse = UserGrants

func (c *Client) GrantsList(ctx context.Context) (*GrantsListResponse, error) {
	e := "profile/grants"
	r, err := CoupleAPIErrors(c.R(ctx).SetResult(GrantsListResponse{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*GrantsListResponse), err
}
