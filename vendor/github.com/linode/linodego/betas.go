package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// Beta Program is a new product or service that is not generally available to all Akamai customers.
// Users must enroll into a beta in order to access the functionality.
type BetaProgram struct {
	Label       string `json:"label"`
	ID          string `json:"id"`
	Description string `json:"description"`

	// Start date of the beta program.
	Started *time.Time `json:"-"`

	// End date of the beta program.
	Ended *time.Time `json:"-"`

	// Greenlight is a program that allows customers to gain access to
	// certain beta programs and to collect direct feedback from those customers.
	GreenlightOnly bool `json:"greenlight_only"`

	// Link to product marketing page for the beta program.
	MoreInfo string `json:"more_info"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (beta *BetaProgram) UnmarshalJSON(b []byte) error {
	type Mask BetaProgram

	p := struct {
		*Mask
		Started *parseabletime.ParseableTime `json:"started"`
		Ended   *parseabletime.ParseableTime `json:"ended"`
	}{
		Mask: (*Mask)(beta),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	beta.Started = (*time.Time)(p.Started)
	beta.Ended = (*time.Time)(p.Ended)

	return nil
}

// ListBetaPrograms lists active beta programs
func (c *Client) ListBetaPrograms(ctx context.Context, opts *ListOptions) ([]BetaProgram, error) {
	response, err := getPaginatedResults[BetaProgram](ctx, c, "/betas", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetBetaProgram gets the beta program's detail with the ID
func (c *Client) GetBetaProgram(ctx context.Context, betaID string) (*BetaProgram, error) {
	e := formatAPIPath("betas/%s", betaID)
	response, err := doGETRequest[BetaProgram](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}
