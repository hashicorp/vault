package radius

import (
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/vault/helper/mfa"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"layeh.com/radius"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: mfa.MFARootPaths(),

			Unauthenticated: []string{
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			//	pathGroups(&b),
			//	pathGroupsList(&b),
			pathUsers(&b),
			pathUsersList(&b),
		},
			mfa.MFAPaths(b.Backend, pathLogin(&b))...,
		),

		AuthRenew: b.pathLoginRenew,
	}

	return &b
}

type backend struct {
	*framework.Backend
}

func (b *backend) Login(req *logical.Request, username string, password string) ([]string, *logical.Response, error) {

	cfg, err := b.Config(req)
	if err != nil {
		return nil, nil, err
	}
	if cfg == nil {
		return nil, logical.ErrorResponse("radius backend not configured"), nil
	}

	hostport := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	packet := radius.New(radius.CodeAccessRequest, []byte(cfg.Secret))
	packet.Add("User-Name", username)
	packet.Add("User-Password", password)
	packet.Add("NAS-Port", uint32(cfg.NasPort))

	//dial_timeout := time.Duration(10) * time.Second
	//raad_timeout := time.Duration(10) * time.Second

	client := radius.Client{
		DialTimeout: time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
	}

	received, err := client.Exchange(packet, hostport)

	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}

	if received.Code != radius.CodeAccessAccept {
		return nil, logical.ErrorResponse("access denied by the authentication server"), nil
	}

	radiusResponse := &logical.Response{
		Data: map[string]interface{}{},
	}

	// Retrieve policies
	var policies []string
	user, err := b.user(req.Storage, username)
	if err == nil && user != nil {
		policies = append(policies, user.Policies...)
	}

	// Policies from each group may overlap
	policies = policyutil.SanitizePolicies(policies, cfg.AllowUnknownUsers)

	if len(policies) == 0 {
		return nil, logical.ErrorResponse("user has no associated policies"), nil
	}

	return policies, radiusResponse, nil
}

const backendHelp = `
The "radius" credential provider allows authentication against
a RADIUS server, checking username and associating users
to set of policies.

Configuration of the server is done through the "config" and "users"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
