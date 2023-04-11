package linodego

import (
	"context"
	"encoding/json"
)

// LishAuthMethod constants start with AuthMethod and include Linode API Lish Authentication Methods
type LishAuthMethod string

// LishAuthMethod constants are the methods of authentication allowed when connecting via Lish
const (
	AuthMethodPasswordKeys LishAuthMethod = "password_keys"
	AuthMethodKeysOnly     LishAuthMethod = "keys_only"
	AuthMethodDisabled     LishAuthMethod = "disabled"
)

// ProfileReferrals represent a User's status in the Referral Program
type ProfileReferrals struct {
	Total     int     `json:"total"`
	Completed int     `json:"completed"`
	Pending   int     `json:"pending"`
	Credit    float64 `json:"credit"`
	Code      string  `json:"code"`
	URL       string  `json:"url"`
}

// Profile represents a Profile object
type Profile struct {
	UID                int              `json:"uid"`
	Username           string           `json:"username"`
	Email              string           `json:"email"`
	Timezone           string           `json:"timezone"`
	EmailNotifications bool             `json:"email_notifications"`
	IPWhitelistEnabled bool             `json:"ip_whitelist_enabled"`
	TwoFactorAuth      bool             `json:"two_factor_auth"`
	Restricted         bool             `json:"restricted"`
	LishAuthMethod     LishAuthMethod   `json:"lish_auth_method"`
	Referrals          ProfileReferrals `json:"referrals"`
	AuthorizedKeys     []string         `json:"authorized_keys"`
}

// ProfileUpdateOptions fields are those accepted by UpdateProfile
type ProfileUpdateOptions struct {
	Email              string         `json:"email,omitempty"`
	Timezone           string         `json:"timezone,omitempty"`
	EmailNotifications *bool          `json:"email_notifications,omitempty"`
	IPWhitelistEnabled *bool          `json:"ip_whitelist_enabled,omitempty"`
	LishAuthMethod     LishAuthMethod `json:"lish_auth_method,omitempty"`
	AuthorizedKeys     *[]string      `json:"authorized_keys,omitempty"`
	TwoFactorAuth      *bool          `json:"two_factor_auth,omitempty"`
	Restricted         *bool          `json:"restricted,omitempty"`
}

// GetUpdateOptions converts a Profile to ProfileUpdateOptions for use in UpdateProfile
func (i Profile) GetUpdateOptions() (o ProfileUpdateOptions) {
	o.Email = i.Email
	o.Timezone = i.Timezone
	o.EmailNotifications = copyBool(&i.EmailNotifications)
	o.IPWhitelistEnabled = copyBool(&i.IPWhitelistEnabled)
	o.LishAuthMethod = i.LishAuthMethod
	authorizedKeys := make([]string, len(i.AuthorizedKeys))
	copy(authorizedKeys, i.AuthorizedKeys)
	o.AuthorizedKeys = &authorizedKeys
	o.TwoFactorAuth = copyBool(&i.TwoFactorAuth)
	o.Restricted = copyBool(&i.Restricted)

	return
}

// GetProfile returns the Profile of the authenticated user
func (c *Client) GetProfile(ctx context.Context) (*Profile, error) {
	e := "profile"
	req := c.R(ctx).SetResult(&Profile{})
	r, err := CoupleAPIErrors(req.Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*Profile), nil
}

// UpdateProfile updates the Profile with the specified id
func (c *Client) UpdateProfile(ctx context.Context, opts ProfileUpdateOptions) (*Profile, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	e := "profile"
	req := c.R(ctx).SetResult(&Profile{}).SetBody(string(body))
	r, err := CoupleAPIErrors(req.Put(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*Profile), nil
}
