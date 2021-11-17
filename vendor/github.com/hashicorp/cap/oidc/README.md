# oidc
[![Go Reference](https://pkg.go.dev/badge/github.com/hashicorp/cap/oidc.svg)](https://pkg.go.dev/github.com/hashicorp/cap/oidc)

oidc is a package for writing clients that integrate with OIDC Providers using
OIDC flows.

Primary types provided by the package:

* `Request`: represents one OIDC authentication flow for a user.  It contains the
  data needed to uniquely represent that one-time flow across the multiple
  interactions needed to complete the OIDC flow the user is attempting.  All
  Requests contain an expiration for the user's OIDC flow.

* `Token`: represents an OIDC id_token, as well as an Oauth2 access_token and
  refresh_token (including the the access_token expiry)

* `Config`: provides the configuration for a typical 3-legged OIDC
  authorization code flow (for example: client ID/Secret, redirectURL, supported
  signing algorithms, additional scopes requested, etc)

* `Provider`: provides integration with an OIDC provider. 
  The provider provides capabilities like: generating an auth URL, exchanging
  codes for tokens, verifying tokens, making user info requests, etc.

* `Alg`: represents asymmetric signing algorithms

* `Error`: provides an error and provides the ability to specify an error code,
  operation that raised the error, the kind of error, and any wrapped error

#### [oidc.callback](callback/)
[![Go Reference](https://pkg.go.dev/badge/github.com/hashicorp/cap/oidc/callback.svg)](https://pkg.go.dev/github.com/hashicorp/cap/oidc/callback)
 
The callback package includes handlers (http.HandlerFunc) which can be used
for the callback leg an OIDC flow. Callback handlers for both the authorization
code flow (with optional PKCE) and the implicit flow are provided.

<hr>

### Examples:

* [CLI example](examples/cli/) which implements an OIDC
  user authentication CLI.  

* [SPA example](examples/spa) which implements an OIDC user
  authentication SPA (single page app). 

<hr>

Example of a provider using an authorization code flow:

```go
// Create a new provider config
pc, err := oidc.NewConfig(
  "http://your-issuer.com/",
  "your_client_id",
  "your_client_secret",
  []oidc.Alg{oidc.RS256},
  []string{"http://your_redirect_url"},
)
if err != nil {
  // handle error
}

// Create a provider
p, err := oidc.NewProvider(pc)
if err != nil {
  // handle error
}
defer p.Done()

	
// Create a Request for a user's authentication attempt that will use the
// authorization code flow.  (See NewRequest(...) using the WithPKCE and
// WithImplicit options for creating a Request that uses those flows.)	
oidcRequest, err := oidc.NewRequest(2 * time.Minute, "http://your_redirect_url/callback")
if err != nil {
  // handle error
}

// Create an auth URL
authURL, err := p.AuthURL(context.Background(), oidcRequest)
if err != nil {
  // handle error
}
fmt.Println("open url to kick-off authentication: ", authURL)
```

Create a http.Handler for OIDC authentication response redirects.

```go
func NewHandler(ctx context.Context, p *oidc.Provider, rw callback.RequestReader) (http.HandlerFunc, error)
  if p == nil { 
    // handle error
  }
  if rw == nil {
    // handle error
  }
  return func(w http.ResponseWriter, r *http.Request) {
    oidcRequest, err := rw.Read(ctx, req.FormValue("state"))
    if err != nil {
      // handle error
    }
    // Exchange(...) will verify the tokens before returning. 
    token, err := p.Exchange(ctx, oidcRequest, req.FormValue("state"), req.FormValue("code"))
    if err != nil {
      // handle error
    }
    var claims map[string]interface{}
    if err := t.IDToken().Claims(&claims); err != nil {
      // handle error
    }

    // Get the user's claims via the provider's UserInfo endpoint
    var infoClaims map[string]interface{}
    err = p.UserInfo(ctx, token.StaticTokenSource(), claims["sub"].(string), &infoClaims)
    if err != nil {
      // handle error
    }
    resp := struct {
	  IDTokenClaims  map[string]interface{}
	  UserInfoClaims map[string]interface{}
    }{claims, infoClaims}
    enc := json.NewEncoder(w)
    if err := enc.Encode(resp); err != nil {
	    // handle error
    }
  }
}
```
  
 
