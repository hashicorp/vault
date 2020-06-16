package golangsdk

/*
AuthOptions stores information needed to authenticate to an OpenStack Cloud.
You can populate one manually, or use a provider's AuthOptionsFromEnv() function
to read relevant information from the standard environment variables. Pass one
to a provider's AuthenticatedClient function to authenticate and obtain a
ProviderClient representing an active session on that provider.

Its fields are the union of those recognized by each identity implementation and
provider.

An example of manually providing authentication information:

  opts := golangsdk.AuthOptions{
    IdentityEndpoint: "https://openstack.example.com:5000/v2.0",
    Username: "{username}",
    Password: "{password}",
    TenantID: "{tenant_id}",
  }

  provider, err := openstack.AuthenticatedClient(opts)

An example of using AuthOptionsFromEnv(), where the environment variables can
be read from a file, such as a standard openrc file:

  opts, err := openstack.AuthOptionsFromEnv()
  provider, err := openstack.AuthenticatedClient(opts)
*/
type AuthOptions struct {
	// IdentityEndpoint specifies the HTTP endpoint that is required to work with
	// the Identity API of the appropriate version. While it's ultimately needed by
	// all of the identity services, it will often be populated by a provider-level
	// function.
	//
	// The IdentityEndpoint is typically referred to as the "auth_url" or
	// "OS_AUTH_URL" in the information provided by the cloud operator.
	IdentityEndpoint string `json:"-"`

	// Username is required if using Identity V2 API. Consult with your provider's
	// control panel to discover your account's username. In Identity V3, either
	// UserID or a combination of Username and DomainID or DomainName are needed.
	Username string `json:"username,omitempty"`
	UserID   string `json:"-"`

	Password string `json:"password,omitempty"`

	// At most one of DomainID and DomainName must be provided if using Username
	// with Identity V3. Otherwise, either are optional.
	DomainID   string `json:"-"`
	DomainName string `json:"name,omitempty"`

	// The TenantID and TenantName fields are optional for the Identity V2 API.
	// The same fields are known as project_id and project_name in the Identity
	// V3 API, but are collected as TenantID and TenantName here in both cases.
	// Some providers allow you to specify a TenantName instead of the TenantId.
	// Some require both. Your provider's authentication policies will determine
	// how these fields influence authentication.
	// If DomainID or DomainName are provided, they will also apply to TenantName.
	// It is not currently possible to authenticate with Username and a Domain
	// and scope to a Project in a different Domain by using TenantName. To
	// accomplish that, the ProjectID will need to be provided as the TenantID
	// option.
	TenantID   string `json:"tenantId,omitempty"`
	TenantName string `json:"tenantName,omitempty"`

	// AllowReauth should be set to true if you grant permission for Gophercloud to
	// cache your credentials in memory, and to allow Gophercloud to attempt to
	// re-authenticate automatically if/when your token expires.  If you set it to
	// false, it will not cache these settings, but re-authentication will not be
	// possible.  This setting defaults to false.
	//
	// NOTE: The reauth function will try to re-authenticate endlessly if left
	// unchecked. The way to limit the number of attempts is to provide a custom
	// HTTP client to the provider client and provide a transport that implements
	// the RoundTripper interface and stores the number of failed retries. For an
	// example of this, see here:
	// https://github.com/rackspace/rack/blob/1.0.0/auth/clients.go#L311
	AllowReauth bool `json:"-"`

	// TokenID allows users to authenticate (possibly as another user) with an
	// authentication token ID.
	TokenID string `json:"-"`

	// AgencyNmae is the name of agnecy
	AgencyName string `json:"-"`

	// AgencyDomainName is the name of domain who created the agency
	AgencyDomainName string `json:"-"`

	// DelegatedProject is the name of delegated project
	DelegatedProject string `json:"-"`
}

// ToTokenV2CreateMap allows AuthOptions to satisfy the AuthOptionsBuilder
// interface in the v2 tokens package
func (opts AuthOptions) ToTokenV2CreateMap() (map[string]interface{}, error) {
	// Populate the request map.
	authMap := make(map[string]interface{})

	if opts.Username != "" {
		if opts.Password != "" {
			authMap["passwordCredentials"] = map[string]interface{}{
				"username": opts.Username,
				"password": opts.Password,
			}
		} else {
			return nil, ErrMissingInput{Argument: "Password"}
		}
	} else if opts.TokenID != "" {
		authMap["token"] = map[string]interface{}{
			"id": opts.TokenID,
		}
	} else {
		return nil, ErrMissingInput{Argument: "Username"}
	}

	if opts.TenantID != "" {
		authMap["tenantId"] = opts.TenantID
	}
	if opts.TenantName != "" {
		authMap["tenantName"] = opts.TenantName
	}

	return map[string]interface{}{"auth": authMap}, nil
}

func (opts *AuthOptions) ToTokenV3CreateMap(scope map[string]interface{}) (map[string]interface{}, error) {
	type domainReq struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	type projectReq struct {
		Domain *domainReq `json:"domain,omitempty"`
		Name   *string    `json:"name,omitempty"`
		ID     *string    `json:"id,omitempty"`
	}

	type userReq struct {
		ID       *string    `json:"id,omitempty"`
		Name     *string    `json:"name,omitempty"`
		Password string     `json:"password"`
		Domain   *domainReq `json:"domain,omitempty"`
	}

	type passwordReq struct {
		User userReq `json:"user"`
	}

	type tokenReq struct {
		ID string `json:"id"`
	}

	type identityReq struct {
		Methods  []string     `json:"methods"`
		Password *passwordReq `json:"password,omitempty"`
		Token    *tokenReq    `json:"token,omitempty"`
	}

	type authReq struct {
		Identity identityReq `json:"identity"`
	}

	type request struct {
		Auth authReq `json:"auth"`
	}

	// Populate the request structure based on the provided arguments. Create and return an error
	// if insufficient or incompatible information is present.
	var req request

	if opts.Password == "" {
		if opts.TokenID != "" {
			// Because we aren't using password authentication, it's an error to also provide any of the user-based authentication
			// parameters.
			if opts.Username != "" {
				return nil, ErrUsernameWithToken{}
			}
			if opts.UserID != "" {
				return nil, ErrUserIDWithToken{}
			}

			// Configure the request for Token authentication.
			req.Auth.Identity.Methods = []string{"token"}
			req.Auth.Identity.Token = &tokenReq{
				ID: opts.TokenID,
			}
		} else {
			// If no password or token ID are available, authentication can't continue.
			return nil, ErrMissingPassword{}
		}
	} else {
		// Password authentication.
		req.Auth.Identity.Methods = []string{"password"}

		// At least one of Username and UserID must be specified.
		if opts.Username == "" && opts.UserID == "" {
			return nil, ErrUsernameOrUserID{}
		}

		if opts.Username != "" {
			// If Username is provided, UserID may not be provided.
			if opts.UserID != "" {
				return nil, ErrUsernameOrUserID{}
			}

			// Either DomainID or DomainName must also be specified.
			if opts.DomainID == "" && opts.DomainName == "" {
				return nil, ErrDomainIDOrDomainName{}
			}

			if opts.DomainID != "" {
				if opts.DomainName != "" {
					return nil, ErrDomainIDOrDomainName{}
				}

				// Configure the request for Username and Password authentication with a DomainID.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{
						Name:     &opts.Username,
						Password: opts.Password,
						Domain:   &domainReq{ID: &opts.DomainID},
					},
				}
			}

			if opts.DomainName != "" {
				// Configure the request for Username and Password authentication with a DomainName.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{
						Name:     &opts.Username,
						Password: opts.Password,
						Domain:   &domainReq{Name: &opts.DomainName},
					},
				}
			}
		}

		if opts.UserID != "" {
			// If UserID is specified, neither DomainID nor DomainName may be.
			if opts.DomainID != "" {
				return nil, ErrDomainIDWithUserID{}
			}
			if opts.DomainName != "" {
				return nil, ErrDomainNameWithUserID{}
			}

			// Configure the request for UserID and Password authentication.
			req.Auth.Identity.Password = &passwordReq{
				User: userReq{ID: &opts.UserID, Password: opts.Password},
			}
		}
	}

	b, err := BuildRequestBody(req, "")
	if err != nil {
		return nil, err
	}

	if len(scope) != 0 {
		b["auth"].(map[string]interface{})["scope"] = scope
	}

	return b, nil
}

func (opts *AuthOptions) ToTokenV3ScopeMap() (map[string]interface{}, error) {
	var scope scopeInfo

	if opts.TenantID != "" {
		scope.ProjectID = opts.TenantID
	} else {
		if opts.TenantName != "" {
			scope.ProjectName = opts.TenantName
			scope.DomainID = opts.DomainID
			scope.DomainName = opts.DomainName
		} else {
			// support scoping to domain
			scope.DomainID = opts.DomainID
			scope.DomainName = opts.DomainName
		}
	}
	return scope.BuildTokenV3ScopeMap()
}

func (opts *AuthOptions) CanReauth() bool {
	return opts.AllowReauth
}

func (opts *AuthOptions) AuthTokenID() string {
	return ""
}

func (opts *AuthOptions) AuthHeaderDomainID() string {
	return ""
}

// Implements the method of AuthOptionsProvider
func (opts AuthOptions) GetIdentityEndpoint() string {
	return opts.IdentityEndpoint
}

type scopeInfo struct {
	ProjectID   string
	ProjectName string
	DomainID    string
	DomainName  string
}

func (scope *scopeInfo) BuildTokenV3ScopeMap() (map[string]interface{}, error) {
	if scope.ProjectName != "" {
		// ProjectName provided: either DomainID or DomainName must also be supplied.
		// ProjectID may not be supplied.
		if scope.DomainID == "" && scope.DomainName == "" {
			return nil, ErrScopeDomainIDOrDomainName{}
		}
		if scope.ProjectID != "" {
			return nil, ErrScopeProjectIDOrProjectName{}
		}

		if scope.DomainID != "" {
			// ProjectName + DomainID
			return map[string]interface{}{
				"project": map[string]interface{}{
					"name":   &scope.ProjectName,
					"domain": map[string]interface{}{"id": &scope.DomainID},
				},
			}, nil
		}

		if scope.DomainName != "" {
			// ProjectName + DomainName
			return map[string]interface{}{
				"project": map[string]interface{}{
					"name":   &scope.ProjectName,
					"domain": map[string]interface{}{"name": &scope.DomainName},
				},
			}, nil
		}
	} else if scope.ProjectID != "" {
		// ProjectID provided. ProjectName, DomainID, and DomainName may not be provided.
		if scope.DomainID != "" {
			return nil, ErrScopeProjectIDAlone{}
		}
		if scope.DomainName != "" {
			return nil, ErrScopeProjectIDAlone{}
		}

		// ProjectID
		return map[string]interface{}{
			"project": map[string]interface{}{
				"id": &scope.ProjectID,
			},
		}, nil
	} else if scope.DomainID != "" {
		// DomainID provided. ProjectID, ProjectName, and DomainName may not be provided.
		if scope.DomainName != "" {
			return nil, ErrScopeDomainIDOrDomainName{}
		}

		// DomainID
		return map[string]interface{}{
			"domain": map[string]interface{}{
				"id": &scope.DomainID,
			},
		}, nil
	} else if scope.DomainName != "" {
		// DomainName
		return map[string]interface{}{
			"domain": map[string]interface{}{
				"name": &scope.DomainName,
			},
		}, nil
	}

	return nil, nil
}

type AgencyAuthOptions struct {
	TokenID          string
	DomainID         string
	AgencyName       string
	AgencyDomainName string
	DelegatedProject string
}

func (opts *AgencyAuthOptions) CanReauth() bool {
	return false
}

func (opts *AgencyAuthOptions) AuthTokenID() string {
	return opts.TokenID
}

func (opts *AgencyAuthOptions) AuthHeaderDomainID() string {
	return opts.DomainID
}

func (opts *AgencyAuthOptions) ToTokenV3ScopeMap() (map[string]interface{}, error) {
	scope := scopeInfo{
		ProjectName: opts.DelegatedProject,
		DomainName:  opts.AgencyDomainName,
	}

	return scope.BuildTokenV3ScopeMap()
}

func (opts *AgencyAuthOptions) ToTokenV3CreateMap(scope map[string]interface{}) (map[string]interface{}, error) {
	type assumeRoleReq struct {
		DomainName string `json:"domain_name"`
		AgencyName string `json:"xrole_name"`
	}

	type identityReq struct {
		Methods    []string      `json:"methods"`
		AssumeRole assumeRoleReq `json:"assume_role"`
	}

	type authReq struct {
		Identity identityReq `json:"identity"`
	}

	var req authReq
	req.Identity.Methods = []string{"assume_role"}
	req.Identity.AssumeRole = assumeRoleReq{
		DomainName: opts.AgencyDomainName,
		AgencyName: opts.AgencyName,
	}
	r, err := BuildRequestBody(req, "auth")
	if err != nil {
		return r, err
	}

	if len(scope) != 0 {
		r["auth"].(map[string]interface{})["scope"] = scope
	}
	return r, nil
}
