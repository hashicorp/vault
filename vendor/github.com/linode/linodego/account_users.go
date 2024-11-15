package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type UserType string

const (
	UserTypeProxy   UserType = "proxy"
	UserTypeParent  UserType = "parent"
	UserTypeChild   UserType = "child"
	UserTypeDefault UserType = "default"
)

// LastLogin represents a LastLogin object
type LastLogin struct {
	LoginDatetime *time.Time `json:"-"`
	Status        string     `json:"status"`
}

// User represents a User object
type User struct {
	Username            string     `json:"username"`
	Email               string     `json:"email"`
	LastLogin           *LastLogin `json:"last_login"`
	UserType            UserType   `json:"user_type"`
	Restricted          bool       `json:"restricted"`
	TFAEnabled          bool       `json:"tfa_enabled"`
	SSHKeys             []string   `json:"ssh_keys"`
	PasswordCreated     *time.Time `json:"-"`
	VerifiedPhoneNumber *string    `json:"verified_phone_number"`
}

// UserCreateOptions fields are those accepted by CreateUser
type UserCreateOptions struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Restricted bool   `json:"restricted"`
}

// UserUpdateOptions fields are those accepted by UpdateUser
type UserUpdateOptions struct {
	Username   string `json:"username,omitempty"`
	Restricted *bool  `json:"restricted,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ll *LastLogin) UnmarshalJSON(b []byte) error {
	type Mask LastLogin

	p := struct {
		*Mask
		LoginDatetime *parseabletime.ParseableTime `json:"login_datetime"`
	}{
		Mask: (*Mask)(ll),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	ll.LoginDatetime = (*time.Time)(p.LoginDatetime)

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *User) UnmarshalJSON(b []byte) error {
	type Mask User

	p := struct {
		*Mask
		PasswordCreated *parseabletime.ParseableTime `json:"password_created"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.PasswordCreated = (*time.Time)(p.PasswordCreated)

	return nil
}

// GetCreateOptions converts a User to UserCreateOptions for use in CreateUser
func (i User) GetCreateOptions() (o UserCreateOptions) {
	o.Username = i.Username
	o.Email = i.Email
	o.Restricted = i.Restricted

	return
}

// GetUpdateOptions converts a User to UserUpdateOptions for use in UpdateUser
func (i User) GetUpdateOptions() (o UserUpdateOptions) {
	o.Username = i.Username
	o.Restricted = copyBool(&i.Restricted)

	return
}

// ListUsers lists Users on the account
func (c *Client) ListUsers(ctx context.Context, opts *ListOptions) ([]User, error) {
	response, err := getPaginatedResults[User](ctx, c, "account/users", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetUser gets the user with the provided ID
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	e := formatAPIPath("account/users/%s", userID)
	response, err := doGETRequest[User](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateUser creates a User.  The email address must be confirmed before the
// User account can be accessed.
func (c *Client) CreateUser(ctx context.Context, opts UserCreateOptions) (*User, error) {
	e := "account/users"
	response, err := doPOSTRequest[User](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateUser updates the User with the specified id
func (c *Client) UpdateUser(ctx context.Context, userID string, opts UserUpdateOptions) (*User, error) {
	e := formatAPIPath("account/users/%s", userID)
	response, err := doPUTRequest[User](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteUser deletes the User with the specified id
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	e := formatAPIPath("account/users/%s", userID)
	err := doDELETERequest(ctx, c, e)
	return err
}
