package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// The details and enrollment information of a Beta program that an account is enrolled in.
type AccountBetaProgram struct {
	Label       string     `json:"label"`
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Started     *time.Time `json:"-"`
	Ended       *time.Time `json:"-"`

	// Date the account was enrolled in the beta program
	Enrolled *time.Time `json:"-"`
}

// AccountBetaProgramCreateOpts fields are those accepted by JoinBetaProgram
type AccountBetaProgramCreateOpts struct {
	ID string `json:"id"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (cBeta *AccountBetaProgram) UnmarshalJSON(b []byte) error {
	type Mask AccountBetaProgram

	p := struct {
		*Mask
		Started  *parseabletime.ParseableTime `json:"started"`
		Ended    *parseabletime.ParseableTime `json:"ended"`
		Enrolled *parseabletime.ParseableTime `json:"enrolled"`
	}{
		Mask: (*Mask)(cBeta),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	cBeta.Started = (*time.Time)(p.Started)
	cBeta.Ended = (*time.Time)(p.Ended)
	cBeta.Enrolled = (*time.Time)(p.Enrolled)

	return nil
}

// ListAccountBetaPrograms lists all beta programs an account is enrolled in.
func (c *Client) ListAccountBetaPrograms(ctx context.Context, opts *ListOptions) ([]AccountBetaProgram, error) {
	response, err := getPaginatedResults[AccountBetaProgram](ctx, c, "/account/betas", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetAccountBetaProgram gets the details of a beta program an account is enrolled in.
func (c *Client) GetAccountBetaProgram(ctx context.Context, betaID string) (*AccountBetaProgram, error) {
	b := formatAPIPath("/account/betas/%s", betaID)

	response, err := doGETRequest[AccountBetaProgram](ctx, c, b)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// JoinBetaProgram enrolls an account into a beta program.
func (c *Client) JoinBetaProgram(ctx context.Context, opts AccountBetaProgramCreateOpts) (*AccountBetaProgram, error) {
	e := "account/betas"
	response, err := doPOSTRequest[AccountBetaProgram](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}
