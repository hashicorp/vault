package okta

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

const (
	profileEmailFilter       = "profile.email"
	profileLoginFilter       = "profile.login"
	profileStatusFilter      = "status"
	profileIDFilter          = "id"
	profileFirstNameFilter   = "profile.firstName"
	profileLastNameFilter    = "profile.lastName"
	profileLastUpdatedFilter = "lastUpdated"
	// UserStatusActive is a  constant to represent OKTA User State returned by the API
	UserStatusActive = "ACTIVE"
	// UserStatusStaged is a  constant to represent OKTA User State returned by the API
	UserStatusStaged = "STAGED"
	// UserStatusProvisioned is a  constant to represent OKTA User State returned by the API
	UserStatusProvisioned = "PROVISIONED"
	// UserStatusRecovery is a  constant to represent OKTA User State returned by the API
	UserStatusRecovery = "RECOVERY"
	// UserStatusLockedOut is a  constant to represent OKTA User State returned by the API
	UserStatusLockedOut = "LOCKED_OUT"
	// UserStatusPasswordExpired is a  constant to represent OKTA User State returned by the API
	UserStatusPasswordExpired = "PASSWORD_EXPIRED"
	// UserStatusSuspended is a  constant to represent OKTA User State returned by the API
	UserStatusSuspended = "SUSPENDED"
	// UserStatusDeprovisioned is a  constant to represent OKTA User State returned by the API
	UserStatusDeprovisioned = "DEPROVISIONED"

	oktaFilterTimeFormat = "2006-01-02T15:05:05.000Z"
)

// UsersService handles communication with the User data related
// methods of the OKTA API.
type UsersService service

// ActivationResponse - Response coming back from a user activation
type activationResponse struct {
	ActivationURL string `json:"activationUrl"`
}

type provider struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type recoveryQuestion struct {
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

type passwordValue struct {
	Value string `json:"value,omitempty"`
}
type credentials struct {
	Password         *passwordValue    `json:"password,omitempty"`
	Provider         *provider         `json:"provider,omitempty"`
	RecoveryQuestion *recoveryQuestion `json:"recovery_question,omitempty"`
}

type userProfile struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Login       string `json:"login"`
	MobilePhone string `json:"mobilePhone,omitempty"`
	SecondEmail string `json:"secondEmail,omitempty"`
	PsEmplid    string `json:"psEmplid,omitempty"`
	NickName    string `json:"nickname,omitempty"`
	DisplayName string `json:"displayName,omitempty"`

	ProfileURL        string `json:"profileUrl,omitempty"`
	PreferredLanguage string `json:"preferredLanguage,omitempty"`
	UserType          string `json:"userType,omitempty"`
	Organization      string `json:"organization,omitempty"`
	Title             string `json:"title,omitempty"`
	Division          string `json:"division,omitempty"`
	Department        string `json:"department,omitempty"`
	CostCenter        string `json:"costCenter,omitempty"`
	EmployeeNumber    string `json:"employeeNumber,omitempty"`
	PrimaryPhone      string `json:"primaryPhone,omitempty"`
	StreetAddress     string `json:"streetAddress,omitempty"`
	City              string `json:"city,omitempty"`
	State             string `json:"state,omitempty"`
	ZipCode           string `json:"zipCode,omitempty"`
	CountryCode       string `json:"countryCode,omitempty"`
}

type userLinks struct {
	ChangePassword struct {
		Href string `json:"href"`
	} `json:"changePassword"`
	ChangeRecoveryQuestion struct {
		Href string `json:"href"`
	} `json:"changeRecoveryQuestion"`
	Deactivate struct {
		Href string `json:"href"`
	} `json:"deactivate"`
	ExpirePassword struct {
		Href string `json:"href"`
	} `json:"expirePassword"`
	ForgotPassword struct {
		Href string `json:"href"`
	} `json:"forgotPassword"`
	ResetFactors struct {
		Href string `json:"href"`
	} `json:"resetFactors"`
	ResetPassword struct {
		Href string `json:"href"`
	} `json:"resetPassword"`
}

// User is a struct that represents a user object from OKTA.
type User struct {
	Activated       string          `json:"activated,omitempty"`
	Created         string          `json:"created,omitempty"`
	Credentials     credentials     `json:"credentials,omitempty"`
	ID              string          `json:"id,omitempty"`
	LastLogin       string          `json:"lastLogin,omitempty"`
	LastUpdated     string          `json:"lastUpdated,omitempty"`
	PasswordChanged string          `json:"passwordChanged,omitempty"`
	Profile         userProfile     `json:"profile"`
	Status          string          `json:"status,omitempty"`
	StatusChanged   string          `json:"statusChanged,omitempty"`
	Links           userLinks       `json:"_links,omitempty"`
	MFAFactors      []userMFAFactor `json:"-,omitempty"`
	Groups          []Group         `json:"-,omitempty"`
}

type userMFAFactor struct {
	ID          string    `json:"id,omitempty"`
	FactorType  string    `json:"factorType,omitempty"`
	Provider    string    `json:"provider,omitempty"`
	VendorName  string    `json:"vendorName,omitempty"`
	Status      string    `json:"status,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	Profile     struct {
		CredentialID string `json:"credentialId,omitempty"`
	} `json:"profile,omitempty"`
}

// NewUser object to create user objects in OKTA
type NewUser struct {
	Profile     userProfile  `json:"profile"`
	Credentials *credentials `json:"credentials,omitempty"`
}

type newPasswordSet struct {
	Credentials credentials `json:"credentials"`
}

type resetPasswordResponse struct {
	ResetPasswordURL string `json:"resetPasswordUrl"`
}

// NewUser - Returns a new user object. This is used to create users in OKTA. It only has the properties that
// OKTA will take as input. The "User" object has more feilds that are OKTA returned like the ID, etc
func (s *UsersService) NewUser() NewUser {
	return NewUser{}
}

// SetPassword Adds a specified password to the new User
func (u *NewUser) SetPassword(passwordIn string) {

	if passwordIn != "" {

		pass := new(passwordValue)
		pass.Value = passwordIn

		var cred *credentials
		if u.Credentials == nil {
			cred = new(credentials)
		} else {
			cred = u.Credentials
		}

		cred.Password = pass
		u.Credentials = cred

	}
}

// SetRecoveryQuestion - Sets a custom security question and answer on a user object
func (u *NewUser) SetRecoveryQuestion(questionIn string, answerIn string) {

	if questionIn != "" && answerIn != "" {
		recovery := new(recoveryQuestion)

		recovery.Question = questionIn
		recovery.Answer = answerIn

		var cred *credentials
		if u.Credentials == nil {
			cred = new(credentials)
		} else {
			cred = u.Credentials
		}
		cred.RecoveryQuestion = recovery
		u.Credentials = cred

	}
}

func (u User) String() string {
	return stringify(u)
	// return fmt.Sprintf("ID: %v \tLogin: %v", u.ID, u.Profile.Login)
}

// GetByID returns a user object for a specific OKTA ID.
// Generally the id input string is the cryptic OKTA key value from User.ID. However, the OKTA API may accept other values like "me", or login shortname
func (s *UsersService) GetByID(id string) (*User, *Response, error) {
	u := fmt.Sprintf("users/%v", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

// UserListFilterOptions is a struct that you can populate which will "filter" user searches
// the exported struct fields should allow you to do different filters based on what is allowed in the OKTA API.
//  The filter OKTA API is limited in the fields it can search
//  NOTE: In the current form you can't add parenthesis and ordering
// OKTA API Supports only a limited number of properties:
// status, lastUpdated, id, profile.login, profile.email, profile.firstName, and profile.lastName.
// http://developer.okta.com/docs/api/resources/users.html#list-users-with-a-filter
type UserListFilterOptions struct {
	Limit         int    `url:"limit,omitempty"`
	EmailEqualTo  string `url:"-"`
	LoginEqualTo  string `url:"-"`
	StatusEqualTo string `url:"-"`
	IDEqualTo     string `url:"-"`

	FirstNameEqualTo string `url:"-"`
	LastNameEqualTo  string `url:"-"`
	//  API documenation says you can search with "starts with" but these don't work

	// FirstNameStartsWith    string    `url:"-"`
	// LastNameStartsWith     string    `url:"-"`

	// This will be built by internal - may not need to export
	FilterString  string     `url:"filter,omitempty"`
	NextURL       *url.URL   `url:"-"`
	GetAllPages   bool       `url:"-"`
	NumberOfPages int        `url:"-"`
	LastUpdated   dateFilter `url:"-"`
}

// PopulateGroups will populate the groups a user is a member of. You pass in a pointer to an existing users
func (s *UsersService) PopulateGroups(user *User) (*Response, error) {
	u := fmt.Sprintf("users/%v/groups", user.ID)
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, err
	}
	// TODO: If user has more than 200 groups this will only return those first 200
	resp, err := s.client.Do(req, &user.Groups)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// PopulateEnrolledFactors will populate the Enrolled MFA Factors a user is a member of.
// You pass in a pointer to an existing users
// http://developer.okta.com/docs/api/resources/factors.html#list-enrolled-factors
func (s *UsersService) PopulateEnrolledFactors(user *User) (*Response, error) {
	u := fmt.Sprintf("users/%v/factors", user.ID)
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, err
	}
	// TODO: If user has more than 200 groups this will only return those first 200
	resp, err := s.client.Do(req, &user.MFAFactors)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// List users with status of LOCKED_OUT
// filter=status eq "LOCKED_OUT"
// List users updated after 06/01/2013 but before 01/01/2014
// filter=lastUpdated gt "2013-06-01T00:00:00.000Z" and lastUpdated lt "2014-01-01T00:00:00.000Z"
// List users updated after 06/01/2013 but before 01/01/2014 with a status of ACTIVE
// filter=lastUpdated gt "2013-06-01T00:00:00.000Z" and lastUpdated lt "2014-01-01T00:00:00.000Z" and status eq "ACTIVE"
// TODO - Currently no way to do parenthesis
// List users updated after 06/01/2013 but with a status of LOCKED_OUT or RECOVERY
// filter=lastUpdated gt "2013-06-01T00:00:00.000Z" and (status eq "LOCKED_OUT" or status eq "RECOVERY")

// OTKA API docs: http://developer.okta.com/docs/api/resources/users.html#list-users-with-a-filter

func appendToFilterString(currFilterString string, appendFilterKey string, appendFilterOperator string, appendFilterValue string) (rs string) {
	if currFilterString != "" {
		rs = fmt.Sprintf("%v and %v %v \"%v\"", currFilterString, appendFilterKey, appendFilterOperator, appendFilterValue)
	} else {
		rs = fmt.Sprintf("%v %v \"%v\"", appendFilterKey, appendFilterOperator, appendFilterValue)
	}

	return rs
}

// ListWithFilter will use the input UserListFilterOptions to find users and return a paged result set
func (s *UsersService) ListWithFilter(opt *UserListFilterOptions) ([]User, *Response, error) {
	var u string
	var err error

	pagesRetreived := 0

	if opt.NextURL != nil {
		u = opt.NextURL.String()
	} else {
		if opt.EmailEqualTo != "" {
			opt.FilterString = appendToFilterString(opt.FilterString, profileEmailFilter, FilterEqualOperator, opt.EmailEqualTo)
		}
		if opt.LoginEqualTo != "" {
			opt.FilterString = appendToFilterString(opt.FilterString, profileLoginFilter, FilterEqualOperator, opt.LoginEqualTo)
		}

		if opt.StatusEqualTo != "" {
			opt.FilterString = appendToFilterString(opt.FilterString, profileStatusFilter, FilterEqualOperator, opt.StatusEqualTo)
		}

		if opt.IDEqualTo != "" {
			opt.FilterString = appendToFilterString(opt.FilterString, profileIDFilter, FilterEqualOperator, opt.IDEqualTo)
		}

		if opt.FirstNameEqualTo != "" {
			opt.FilterString = appendToFilterString(opt.FilterString, profileFirstNameFilter, FilterEqualOperator, opt.FirstNameEqualTo)
		}

		if opt.LastNameEqualTo != "" {
			opt.FilterString = appendToFilterString(opt.FilterString, profileLastNameFilter, FilterEqualOperator, opt.LastNameEqualTo)
		}

		//  API documenation says you can search with "starts with" but these don't work
		// if opt.FirstNameStartsWith != "" {
		// 	opt.FilterString = appendToFilterString(opt.FilterString, profileFirstNameFilter, filterStartsWithOperator, opt.FirstNameStartsWith)
		// }

		// if opt.LastNameStartsWith != "" {
		// 	opt.FilterString = appendToFilterString(opt.FilterString, profileLastNameFilter, filterStartsWithOperator, opt.LastNameStartsWith)
		// }

		if !opt.LastUpdated.Value.IsZero() {
			opt.FilterString = appendToFilterString(opt.FilterString, profileLastUpdatedFilter, opt.LastUpdated.Operator, opt.LastUpdated.Value.UTC().Format(oktaFilterTimeFormat))
		}

		if opt.Limit == 0 {
			opt.Limit = defaultLimit
		}

		u, err = addOptions("users", opt)

	}

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	users := make([]User, 1)
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	pagesRetreived++

	if (opt.NumberOfPages > 0 && pagesRetreived < opt.NumberOfPages) || opt.GetAllPages {

		for {

			if pagesRetreived == opt.NumberOfPages {
				break
			}
			if resp.NextURL != nil {
				var userPage []User
				pageOption := new(UserListFilterOptions)
				pageOption.NextURL = resp.NextURL
				pageOption.NumberOfPages = 1
				pageOption.Limit = opt.Limit

				userPage, resp, err = s.ListWithFilter(pageOption)
				if err != nil {
					return users, resp, err
				} else {
					users = append(users, userPage...)
					pagesRetreived++
				}
			} else {
				break
			}
		}
	}
	return users, resp, err
}

// Create - Creates a new user. You must pass in a "newUser" object created from Users.NewUser()
// There are many differnt reasons that OKTA may reject the request so you have to check the error messages
func (s *UsersService) Create(userIn NewUser, createAsActive bool) (*User, *Response, error) {

	u := fmt.Sprintf("users?activate=%v", createAsActive)

	req, err := s.client.NewRequest("POST", u, userIn)

	if err != nil {
		return nil, nil, err
	}

	newUser := new(User)
	resp, err := s.client.Do(req, newUser)
	if err != nil {
		return nil, resp, err
	}

	return newUser, resp, err
}

// Activate Activates a user. You can have OKTA send an email by including a "sendEmail=true"
// If you pass in sendEmail=false, then activationResponse.ActivationURL will have a string URL that
// can be sent to the end user. You can discard response if sendEmail=true
func (s *UsersService) Activate(id string, sendEmail bool) (*activationResponse, *Response, error) {
	u := fmt.Sprintf("users/%v/lifecycle/activate?sendEmail=%v", id, sendEmail)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	activationInfo := new(activationResponse)
	resp, err := s.client.Do(req, activationInfo)

	if err != nil {
		return nil, resp, err
	}

	return activationInfo, resp, err
}

// Deactivate - Deactivates a user
func (s *UsersService) Deactivate(id string) (*Response, error) {
	u := fmt.Sprintf("users/%v/lifecycle/deactivate", id)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req, nil)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// Suspend - Suspends a user - If user is NOT active an Error will come back based on OKTA API:
// http://developer.okta.com/docs/api/resources/users.html#suspend-user
func (s *UsersService) Suspend(id string) (*Response, error) {
	u := fmt.Sprintf("users/%v/lifecycle/suspend", id)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req, nil)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// Unsuspend - Unsuspends a user - If user is NOT SUSPENDED, an Error will come back based on OKTA API:
// http://developer.okta.com/docs/api/resources/users.html#unsuspend-user
func (s *UsersService) Unsuspend(id string) (*Response, error) {
	u := fmt.Sprintf("users/%v/lifecycle/unsuspend", id)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req, nil)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// Unlock - Unlocks a user - Per docs, only for OKTA Mastered Account
// http://developer.okta.com/docs/api/resources/users.html#unlock-user
func (s *UsersService) Unlock(id string) (*Response, error) {
	u := fmt.Sprintf("users/%v/lifecycle/unlock", id)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req, nil)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// SetPassword - Sets a user password to an Admin provided String
func (s *UsersService) SetPassword(id string, newPassword string) (*User, *Response, error) {

	if id == "" || newPassword == "" {
		return nil, nil, errors.New("please provide a User ID and Password")
	}

	passwordUpdate := new(newPasswordSet)

	pass := new(passwordValue)
	pass.Value = newPassword

	passwordUpdate.Credentials.Password = pass

	u := fmt.Sprintf("users/%v", id)
	req, err := s.client.NewRequest("POST", u, passwordUpdate)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

// ResetPassword - Generates a one-time token (OTT) that can be used to reset a user’s password.
// The OTT link can be automatically emailed to the user or returned to the API caller and distributed using a custom flow.
// http://developer.okta.com/docs/api/resources/users.html#reset-password
// If you pass in sendEmail=false, then resetPasswordResponse.resetPasswordUrl will have a string URL that
// can be sent to the end user. You can discard response if sendEmail=true
func (s *UsersService) ResetPassword(id string, sendEmail bool) (*resetPasswordResponse, *Response, error) {
	u := fmt.Sprintf("users/%v/lifecycle/reset_password?sendEmail=%v", id, sendEmail)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	resetInfo := new(resetPasswordResponse)
	resp, err := s.client.Do(req, resetInfo)

	if err != nil {
		return nil, resp, err
	}

	return resetInfo, resp, err
}

// PopulateMFAFactors will populate the MFA Factors a user is a member of. You pass in a pointer to an existing users
func (s *UsersService) PopulateMFAFactors(user *User) (*Response, error) {
	u := fmt.Sprintf("users/%v/factors", user.ID)

	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, &user.MFAFactors)
	if err != nil {
		return resp, err
	}

	return resp, err
}
