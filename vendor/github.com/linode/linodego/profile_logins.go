package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// Profile represents a Profile object
type ProfileLogin struct {
	Datetime   *time.Time `json:"datetime"`
	ID         int        `json:"id"`
	IP         string     `json:"ip"`
	Restricted bool       `json:"restricted"`
	Status     string     `json:"status"`
	Username   string     `json:"username"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *ProfileLogin) UnmarshalJSON(b []byte) error {
	type Mask ProfileLogin

	l := struct {
		*Mask
		Datetime *parseabletime.ParseableTime `json:"datetime"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &l); err != nil {
		return err
	}

	i.Datetime = (*time.Time)(l.Datetime)

	return nil
}

// GetProfileLogin returns the Profile Login of the authenticated user
func (c *Client) GetProfileLogin(ctx context.Context, id int) (*ProfileLogin, error) {
	e := formatAPIPath("profile/logins/%d", id)
	response, err := doGETRequest[ProfileLogin](ctx, c, e)
	return response, err
}

// ListProfileLogins lists Profile Logins of the authenticated user
func (c *Client) ListProfileLogins(ctx context.Context, opts *ListOptions) ([]ProfileLogin, error) {
	response, err := getPaginatedResults[ProfileLogin](ctx, c, "profile/logins", opts)
	return response, err
}
