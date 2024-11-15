package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// Stackscript represents a Linode StackScript
type Stackscript struct {
	ID                int               `json:"id"`
	Username          string            `json:"username"`
	Label             string            `json:"label"`
	Description       string            `json:"description"`
	Ordinal           int               `json:"ordinal"`
	LogoURL           string            `json:"logo_url"`
	Images            []string          `json:"images"`
	DeploymentsTotal  int               `json:"deployments_total"`
	DeploymentsActive int               `json:"deployments_active"`
	IsPublic          bool              `json:"is_public"`
	Mine              bool              `json:"mine"`
	Created           *time.Time        `json:"-"`
	Updated           *time.Time        `json:"-"`
	RevNote           string            `json:"rev_note"`
	Script            string            `json:"script"`
	UserDefinedFields *[]StackscriptUDF `json:"user_defined_fields"`
	UserGravatarID    string            `json:"user_gravatar_id"`
}

// StackscriptUDF define a single variable that is accepted by a Stackscript
type StackscriptUDF struct {
	// A human-readable label for the field that will serve as the input prompt for entering the value during deployment.
	Label string `json:"label"`

	// The name of the field.
	Name string `json:"name"`

	// An example value for the field.
	Example string `json:"example"`

	// A list of acceptable single values for the field.
	OneOf string `json:"oneOf,omitempty"`

	// A list of acceptable values for the field in any quantity, combination or order.
	ManyOf string `json:"manyOf,omitempty"`

	// The default value. If not specified, this value will be used.
	Default string `json:"default,omitempty"`
}

// StackscriptCreateOptions fields are those accepted by CreateStackscript
type StackscriptCreateOptions struct {
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	IsPublic    bool     `json:"is_public"`
	RevNote     string   `json:"rev_note"`
	Script      string   `json:"script"`
}

// StackscriptUpdateOptions fields are those accepted by UpdateStackscript
type StackscriptUpdateOptions StackscriptCreateOptions

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Stackscript) UnmarshalJSON(b []byte) error {
	type Mask Stackscript

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)

	return nil
}

// GetCreateOptions converts a Stackscript to StackscriptCreateOptions for use in CreateStackscript
func (i Stackscript) GetCreateOptions() StackscriptCreateOptions {
	return StackscriptCreateOptions{
		Label:       i.Label,
		Description: i.Description,
		Images:      i.Images,
		IsPublic:    i.IsPublic,
		RevNote:     i.RevNote,
		Script:      i.Script,
	}
}

// GetUpdateOptions converts a Stackscript to StackscriptUpdateOptions for use in UpdateStackscript
func (i Stackscript) GetUpdateOptions() StackscriptUpdateOptions {
	return StackscriptUpdateOptions{
		Label:       i.Label,
		Description: i.Description,
		Images:      i.Images,
		IsPublic:    i.IsPublic,
		RevNote:     i.RevNote,
		Script:      i.Script,
	}
}

// ListStackscripts lists Stackscripts
func (c *Client) ListStackscripts(ctx context.Context, opts *ListOptions) ([]Stackscript, error) {
	response, err := getPaginatedResults[Stackscript](ctx, c, "linode/stackscripts", opts)
	return response, err
}

// GetStackscript gets the Stackscript with the provided ID
func (c *Client) GetStackscript(ctx context.Context, scriptID int) (*Stackscript, error) {
	e := formatAPIPath("linode/stackscripts/%d", scriptID)
	response, err := doGETRequest[Stackscript](ctx, c, e)
	return response, err
}

// CreateStackscript creates a StackScript
func (c *Client) CreateStackscript(ctx context.Context, opts StackscriptCreateOptions) (*Stackscript, error) {
	e := "linode/stackscripts"
	response, err := doPOSTRequest[Stackscript](ctx, c, e, opts)
	return response, err
}

// UpdateStackscript updates the StackScript with the specified id
func (c *Client) UpdateStackscript(ctx context.Context, scriptID int, opts StackscriptUpdateOptions) (*Stackscript, error) {
	e := formatAPIPath("linode/stackscripts/%d", scriptID)
	response, err := doPUTRequest[Stackscript](ctx, c, e, opts)
	return response, err
}

// DeleteStackscript deletes the StackScript with the specified id
func (c *Client) DeleteStackscript(ctx context.Context, scriptID int) error {
	e := formatAPIPath("linode/stackscripts/%d", scriptID)
	err := doDELETERequest(ctx, c, e)
	return err
}
