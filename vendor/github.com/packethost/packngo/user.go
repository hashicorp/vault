package packngo

import (
	"path"
)

const (
	usersBasePath = "/users"
	userBasePath  = "/user"
)

// UserService interface defines available user methods
type UserService interface {
	Create(*UserCreateRequest) (*User, *Response, error)
	Current() (*User, *Response, error)
	List(*ListOptions) ([]User, *Response, error)
	Get(string, *GetOptions) (*User, *Response, error)
	Update(*UserUpdateRequest) (*User, *Response, error)
}

type usersRoot struct {
	Users []User `json:"users"`
	Meta  meta   `json:"meta"`
}

// SocialAccounts are social usernames or urls
type SocialAccounts struct {
	GitHub   string `json:"github,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
	Twitter  string `json:"twitter,omitempty"`
	Gravatar string `json:"gravatar,omitempty"`
	Keybase  string `json:"keybase,omitempty"`
}

// User represents an Equinix Metal user
type User struct {
	ID               string          `json:"id"`
	ShortID          string          `json:"short_id"`
	FirstName        string          `json:"first_name,omitempty"`
	LastName         string          `json:"last_name,omitempty"`
	FullName         string          `json:"full_name,omitempty"`
	Email            string          `json:"email,omitempty"`
	SocialAccounts   *SocialAccounts `json:"social_accounts,omitempty"`
	CustomData       interface{}     `json:"customdata,omitempty"`
	OptIn            *bool           `json:"opt_in,omitempty"`
	OptInUpdatedAt   string          `json:"opt_in_updated_at,omitempty"`
	DefaultProjectID *string         `json:"default_project_id,omitempty"`
	NumberOfSSHKeys  int             `json:"number_of_ssh_keys,omitempty"`
	Language         *string         `json:"language,omitempty"`
	// MailingAddress TODO: format
	VerificationStage string `json:"verification_stage,omitempty"`
	MaxProjects       *int   `json:"max_projects,omitempty"`
	LastLogin         string `json:"last_login_at,omitempty"`

	// Features effect the behavior of the API and UI
	Features []string `json:"features,omitempty"`

	// TwoFactorAuth is the form of two factor auth, "app" or "sms"
	// Renamed from TwoFactor in packngo v0.16 to match API
	TwoFactorAuth         string  `json:"two_factor_auth,omitempty"`
	DefaultOrganizationID string  `json:"default_organization_id,omitempty"`
	AvatarURL             string  `json:"avatar_url,omitempty"`
	AvatarThumbURL        string  `json:"avatar_thumb_url,omitempty"`
	Created               string  `json:"created_at,omitempty"`
	Updated               string  `json:"updated_at,omitempty"`
	TimeZone              string  `json:"timezone,omitempty"`
	Emails                []Email `json:"emails,omitempty"`
	PhoneNumber           string  `json:"phone_number,omitempty"`
	URL                   string  `json:"href,omitempty"`
	Restricted            bool    `json:"restricted,omitempty"`
	Staff                 bool    `json:"staff,omitempty"`
}

// UserLite is an abbreviated listing of an Equinix Metal user
type UserLite struct {
	*Href          `json:",inline"`
	ID             string     `json:"id"`
	ShortID        string     `json:"short_id"`
	FirstName      string     `json:"first_name,omitempty"`
	LastName       string     `json:"last_name,omitempty"`
	FullName       string     `json:"full_name,omitempty"`
	Email          string     `json:"email,omitempty"`
	CreatedAt      *Timestamp `json:"created_at,omitempty"`
	UpdatedAt      *Timestamp `json:"updated_at,omitempty"`
	Level          string     `json:"level,omitempty"`
	AvatarThumbURL string     `json:"avatar_thumb_url,omitempty"`
}

// UserCreateRequest struct for UserService.Create
type UserCreateRequest struct {
	InvitationID string         `json:"invitation_id,omitempty"`
	Nonce        string         `json:"nonce,omitempty"`
	FirstName    string         `json:"first_name,omitempty"`
	LastName     string         `json:"last_name,omitempty"`
	Password     string         `json:"password,omitempty"`
	Customdata   *interface{}   `json:"customdata,omitempty"`
	Emails       []EmailRequest `json:"emails,omitempty"`
}

// UserUpdateRequest struct for UserService.Update
type UserUpdateRequest struct {
	FirstName   *string      `json:"first_name,omitempty"`
	LastName    *string      `json:"last_name,omitempty"`
	PhoneNumber *string      `json:"phone_number,omitempty"`
	Timezone    *string      `json:"timezone,omitempty"`
	Password    *string      `json:"password,omitempty"`
	Avatar      *string      `json:"avatar,omitempty"`
	Customdata  *interface{} `json:"customdata,omitempty"`
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

// Create a User with the given UserCreateRequest. New user VerificationStage
// will be AccountCreated, unless UserCreateRequest contains an valid
// InvitationID and Nonce in which case the VerificationStage will be Verified.
func (s *UserServiceOp) Create(createRequest *UserCreateRequest) (*User, *Response, error) {
	opts := &GetOptions{}
	endpointPath := path.Join(usersBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	user := new(User)

	resp, err := s.client.DoRequest("POST", apiPathQuery, createRequest, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
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
	if validateErr := ValidateUUID(userID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(usersBasePath, userID)
	apiPathQuery := opts.WithQuery(endpointPath)
	user := new(User)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

// Update updates the current user
func (s *UserServiceOp) Update(updateRequest *UserUpdateRequest) (*User, *Response, error) {
	opts := &GetOptions{}
	endpointPath := path.Join(userBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	user := new(User)

	resp, err := s.client.DoRequest("PUT", apiPathQuery, updateRequest, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}
