package packngo

import "path"

const usersBasePath = "/users"
const userBasePath = "/user"

// UserService interface defines available user methods
type UserService interface {
	List(*ListOptions) ([]User, *Response, error)
	Get(string, *GetOptions) (*User, *Response, error)
	Current() (*User, *Response, error)
}

type usersRoot struct {
	Users []User `json:"users"`
	Meta  meta   `json:"meta"`
}

// User represents an Equinix Metal user
type User struct {
	ID                    string  `json:"id"`
	FirstName             string  `json:"first_name,omitempty"`
	LastName              string  `json:"last_name,omitempty"`
	FullName              string  `json:"full_name,omitempty"`
	Email                 string  `json:"email,omitempty"`
	TwoFactor             string  `json:"two_factor_auth,omitempty"`
	DefaultOrganizationID string  `json:"default_organization_id,omitempty"`
	AvatarURL             string  `json:"avatar_url,omitempty"`
	Facebook              string  `json:"twitter,omitempty"`
	Twitter               string  `json:"facebook,omitempty"`
	LinkedIn              string  `json:"linkedin,omitempty"`
	Created               string  `json:"created_at,omitempty"`
	Updated               string  `json:"updated_at,omitempty"`
	TimeZone              string  `json:"timezone,omitempty"`
	Emails                []Email `json:"emails,omitempty"`
	PhoneNumber           string  `json:"phone_number,omitempty"`
	URL                   string  `json:"href,omitempty"`
	VPN                   bool    `json:"vpn"`
}

func (u User) String() string {
	return Stringify(u)
}

// UserServiceOp implements UserService
type UserServiceOp struct {
	client *Client
}

// Get method gets a user by userID
func (s *UserServiceOp) List(opts *ListOptions) (users []User, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(usersBasePath)

	for {
		subset := new(usersRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		users = append(users, subset.Users...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Returns the user object for the currently logged-in user.
func (s *UserServiceOp) Current() (*User, *Response, error) {
	user := new(User)

	resp, err := s.client.DoRequest("GET", userBasePath, nil, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

func (s *UserServiceOp) Get(userID string, opts *GetOptions) (*User, *Response, error) {
	endpointPath := path.Join(usersBasePath, userID)
	apiPathQuery := opts.WithQuery(endpointPath)
	user := new(User)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}
