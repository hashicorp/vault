package cassandra

import (
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	defaultCreationCQL = `CREATE USER '{{username}}' WITH PASSWORD '{{password}}' NOSUPERUSER;`
	defaultRollbackCQL = `DROP USER '{{username}}';`
)

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},

			"creation_cql": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: defaultCreationCQL,
				Description: `CQL to create a user and optionally grant
authorization. If not supplied, a default that
creates non-superuser accounts with the built-in
password authenticator will be used; no
authorization grants will be configured. Separate
statements by semicolons; use @file to load from a
file. Valid template values are '{{username}}' and
'{{password}}' -- the single quotes are important!`,
			},

			"rollback_cql": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: defaultRollbackCQL,
				Description: `CQL to roll back an account operation. This will
be used if there is an error during execution of a
statement passed in via the "creation_cql" parameter
parameter. The default simply drops the user, which
should generally be sufficient. Separate statements
by semicolons; use @file to load from a file. Valid
template values are '{{username}}' and
'{{password}}' -- the single quotes are important!`,
			},

			"lease": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "4h",
				Description: "The lease length; defaults to 4 hours",
			},

			"consistency": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "Quorum",
				Description: "The consistency level for the operations; defaults to Quorum.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.UpdateOperation: b.pathRoleCreate,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func getRole(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("role/" + data.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := getRole(req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(role).Map(),
	}, nil
}

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	creationCQL := data.Get("creation_cql").(string)

	rollbackCQL := data.Get("rollback_cql").(string)

	leaseRaw := data.Get("lease").(string)
	lease, err := time.ParseDuration(leaseRaw)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error parsing lease value of %s: %s", leaseRaw, err)), nil
	}

	consistencyStr := data.Get("consistency").(string)
	_, err = gocql.ParseConsistencyWrapper(consistencyStr)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error parsing consistency value of %q: %v", consistencyStr, err)), nil
	}

	entry := &roleEntry{
		Lease:       lease,
		CreationCQL: creationCQL,
		RollbackCQL: rollbackCQL,
		Consistency: consistencyStr,
	}

	// Store it
	entryJSON, err := logical.StorageEntryJSON("role/"+name, entry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entryJSON); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleEntry struct {
	CreationCQL string        `json:"creation_cql" structs:"creation_cql"`
	Lease       time.Duration `json:"lease" structs:"lease"`
	RollbackCQL string        `json:"rollback_cql" structs:"rollback_cql"`
	Consistency string        `json:"consistency" structs:"consistency"`
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.

The "creation_cql" parameter customizes the CQL string used to create users
and assign them grants. This can be a sequence of CQL queries separated by
semicolons. Some substitution will be done to the CQL string for certain keys.
The names of the variables must be surrounded by '{{' and '}}' to be replaced.
Note that it is important that single quotes are used, not double quotes.

  * "username" - The random username generated for the DB user.

  * "password" - The random password generated for the DB user.

If no "creation_cql" parameter is given, a default will be used:

` + defaultCreationCQL + `

This default should be suitable for Cassandra installations using the password
authenticator but not configured to use authorization.

Similarly, the "rollback_cql" is used if user creation fails, in the absense of
Cassandra transactions. The default should be suitable for almost any
instance of Cassandra:

` + defaultRollbackCQL + `

"lease" the lease time; if not set the mount/system defaults are used.
`
