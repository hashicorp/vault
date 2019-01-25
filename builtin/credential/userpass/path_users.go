package userpass

import (
	"context"
	"fmt"
	"strings"
	"time"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/tokenhelper"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathUsersList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathUserList,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func pathUsers(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "users/" + framework.GenericNameRegex("username"),
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username for this user.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathUserDelete,
			logical.ReadOperation:   b.pathUserRead,
			logical.UpdateOperation: b.pathUserWrite,
			logical.CreateOperation: b.pathUserWrite,
		},

		ExistenceCheck: b.userExistenceCheck,

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}

	tokenhelper.AddTokenFields(p.Fields)

	return p
}

func (b *backend) userExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	userEntry, err := b.user(ctx, req.Storage, data.Get("username").(string))
	if err != nil {
		return false, err
	}

	return userEntry != nil, nil
}

func (b *backend) user(ctx context.Context, s logical.Storage, username string) (*UserEntry, error) {
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	entry, err := s.Get(ctx, "user/"+strings.ToLower(username))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result UserEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	var needsUpgrade bool
	if result.OldTTL != 0 {
		needsUpgrade = true
		result.TTL = result.OldTTL
		result.OldTTL = 0
	}
	if result.OldMaxTTL != 0 {
		needsUpgrade = true
		result.MaxTTL = result.OldMaxTTL
		result.OldMaxTTL = 0
	}
	if len(result.OldPolicies) != 0 {
		needsUpgrade = true
		result.Policies = result.OldPolicies
		result.OldPolicies = nil
	}
	if result.OldBoundCIDRs != nil {
		needsUpgrade = true
		result.BoundCIDRs = result.OldBoundCIDRs
		result.OldBoundCIDRs = nil
	}
	if needsUpgrade && (b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby)) {
		if err := b.setUser(ctx, s, strings.ToLower(username), &result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (b *backend) setUser(ctx context.Context, s logical.Storage, username string, userEntry *UserEntry) error {
	entry, err := logical.StorageEntryJSON("user/"+username, userEntry)
	if err != nil {
		return err
	}

	return s.Put(ctx, entry)
}

func (b *backend) pathUserList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	users, err := req.Storage.List(ctx, "user/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(users), nil
}

func (b *backend) pathUserDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "user/"+strings.ToLower(d.Get("username").(string)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathUserRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	user, err := b.user(ctx, req.Storage, strings.ToLower(d.Get("username").(string)))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	user.PopulateTokenData(data)
	return &logical.Response{
		Data: data,
	}, nil
}

func (b *backend) userCreateUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("username").(string))
	userEntry, err := b.user(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}
	// Due to existence check, user will only be nil if it's a create operation
	if userEntry == nil {
		userEntry = &UserEntry{}
	}

	if _, ok := d.GetOk("password"); ok {
		userErr, intErr := b.updateUserPassword(req, d, userEntry)
		if intErr != nil {
			return nil, err
		}
		if userErr != nil {
			return logical.ErrorResponse(userErr.Error()), logical.ErrInvalidRequest
		}
	}

	if err := userEntry.ParseTokenFields(req, d); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	return nil, b.setUser(ctx, req.Storage, username, userEntry)
}

func (b *backend) pathUserWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	password := d.Get("password").(string)
	if req.Operation == logical.CreateOperation && password == "" {
		return logical.ErrorResponse("missing password"), logical.ErrInvalidRequest
	}
	return b.userCreateUpdate(ctx, req, d)
}

type UserEntry struct {
	tokenhelper.TokenParams

	// Password is deprecated in Vault 0.2 in favor of
	// PasswordHash, but is retained for backwards compatibility.
	Password string

	// PasswordHash is a bcrypt hash of the password. This is
	// used instead of the actual password in Vault 0.2+.
	PasswordHash []byte

	// These token-related fields have been moved to the embedded tokenhelper.TokenParams struct
	OldPolicies []string `json:"Policies"`

	// Duration after which the user will be revoked unless renewed
	OldTTL time.Duration `json:"TTL"`

	// Maximum duration for which user can be valid
	OldMaxTTL time.Duration `json:"MaxTTL"`

	OldBoundCIDRs []*sockaddr.SockAddrMarshaler `json:"BoundCIDRs"`
}

const pathUserHelpSyn = `
Manage users allowed to authenticate.
`

const pathUserHelpDesc = `
This endpoint allows you to create, read, update, and delete users
that are allowed to authenticate.

Deleting a user will not revoke auth for prior authenticated users
with that name. To do this, do a revoke on "login/<username>" for
the username you want revoked. If you don't need to revoke login immediately,
then the next renew will cause the lease to expire.
`
