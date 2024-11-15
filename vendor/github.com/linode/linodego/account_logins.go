package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type Login struct {
	ID         int        `json:"id"`
	Datetime   *time.Time `json:"datetime"`
	IP         string     `json:"ip"`
	Restricted bool       `json:"restricted"`
	Username   string     `json:"username"`
	Status     string     `json:"status"`
}

func (c *Client) ListLogins(ctx context.Context, opts *ListOptions) ([]Login, error) {
	response, err := getPaginatedResults[Login](ctx, c, "account/logins", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Login) UnmarshalJSON(b []byte) error {
	type Mask Login

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

func (c *Client) GetLogin(ctx context.Context, loginID int) (*Login, error) {
	e := formatAPIPath("account/logins/%d", loginID)

	response, err := doGETRequest[Login](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}
