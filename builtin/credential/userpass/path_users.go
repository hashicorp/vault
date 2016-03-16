package userpass

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/crypto/bcrypt"
)

func pathUsers(b *backend) *framework.Path {
	return &framework.Path{
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

			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Comma-separated list of policies",
			},
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "The lease duration which decides login expiration",
			},
			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "Maximum duration after which login should expire",
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
}

func (b *backend) userExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	username := data.Get("username").(string)
	if username == "" {
		return false, fmt.Errorf("missing username")
	}

	userEntry, err := b.user(req.Storage, username)
	if err != nil {
		return false, err
	}

	return userEntry != nil, nil
}

func (b *backend) user(s logical.Storage, n string) (*UserEntry, error) {
	entry, err := s.Get("user/" + strings.ToLower(n))
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

	return &result, nil
}

func (b *backend) setUser(s logical.Storage, username string, userEntry *UserEntry) error {
	entry, err := logical.StorageEntryJSON("user/"+username, userEntry)
	if err != nil {
		return err
	}

	return s.Put(entry)
}

func (b *backend) pathUserDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("user/" + strings.ToLower(d.Get("username").(string)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathUserRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	user, err := b.user(req.Storage, strings.ToLower(d.Get("username").(string)))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policies": strings.Join(user.Policies, ","),
		},
	}, nil
}

func (b *backend) userCreateUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("username").(string))
	userEntry, err := b.user(req.Storage, username)
	if err != nil {
		return nil, err
	}
	// Due to existence check, user will only be nil if it's a create operation
	if userEntry == nil {
		userEntry = &UserEntry{}
	}

	// Set/update the values of UserEntry only if fields are supplied
	if passwordRaw, ok := d.GetOk("password"); ok {
		// Generate a hash of the password
		hash, err := bcrypt.GenerateFromPassword([]byte(passwordRaw.(string)), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		userEntry.PasswordHash = hash
	}

	if policiesRaw, ok := d.GetOk("policies"); ok {
		policies := strings.Split(policiesRaw.(string), ",")
		for i, p := range policies {
			policies[i] = strings.TrimSpace(p)
		}
		userEntry.Policies = policies
	}

	_, ttlSet := d.GetOk("ttl")
	_, maxTTLSet := d.GetOk("max_ttl")
	if ttlSet || maxTTLSet {
		ttlStr := d.Get("ttl").(string)
		maxTTLStr := d.Get("max_ttl").(string)
		userEntry.TTL, userEntry.MaxTTL, err = b.SanitizeTTL(ttlStr, maxTTLStr)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("err: %s", err)), nil
		}
	}

	return nil, b.setUser(req.Storage, username, userEntry)
}

func (b *backend) pathUserWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	password := d.Get("password").(string)
	if password == "" {
		return nil, fmt.Errorf("missing password")
	}
	return b.userCreateUpdate(req, d)
}

type UserEntry struct {
	// Password is deprecated in Vault 0.2 in favor of
	// PasswordHash, but is retained for backwards compatibilty.
	Password string

	// PasswordHash is a bcrypt hash of the password. This is
	// used instead of the actual password in Vault 0.2+.
	PasswordHash []byte

	Policies []string

	// Duration after which the user will be revoked unless renewed
	TTL time.Duration

	// Maximum duration for which user can be valid
	MaxTTL time.Duration
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
