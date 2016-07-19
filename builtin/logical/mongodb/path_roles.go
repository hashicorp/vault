package mongodb

import (
	"encoding/json"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"db": {
				Type:        framework.TypeString,
				Description: "Name of the authentication database for users generated for this role.",
			},
			"roles": {
				Type:        framework.TypeString,
				Description: "MongoDB roles to assign to the users generated for this role.",
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

func (b *backend) Role(s logical.Storage, n string) (*roleStorageEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleStorageEntry
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
	role, err := b.Role(req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	rolesJsonBytes, err := json.Marshal(role.MongoDBRoles.toStandardRolesArray())
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"db":    role.DB,
			"roles": string(rolesJsonBytes),
		},
	}, nil
}

func (b *backend) pathRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("Missing name"), nil
	}

	roleDB := data.Get("db").(string)
	if roleDB == "" {
		return logical.ErrorResponse("db parameter is required"), nil
	}

	// Example roles JSON:
	//
	// [ "readWrite", { "role": "readWrite", "db": "test" } ]
	//
	// For storage, we convert such an array into a homogeneous array of role documents like:
	//
	// [ { "role": "readWrite" }, { "role": "readWrite", "db": "test" } ]
	//
	var roles []mongodbRole
	rolesJson := []byte(data.Get("roles").(string))
	if len(rolesJson) > 0 {
		var rolesArray []interface{}
		err := json.Unmarshal(rolesJson, &rolesArray)
		if err != nil {
			return nil, err
		}
		for _, rawRole := range rolesArray {
			switch role := rawRole.(type) {
			case string:
				roles = append(roles, mongodbRole{Role: role})
			case map[string]interface{}:
				if db, ok := role["db"].(string); ok {
					if roleName, ok := role["role"].(string); ok {
						roles = append(roles, mongodbRole{Role: roleName, DB: db})
					}
				}
			}
		}
	}

	// Store it
	entry, err := logical.StorageEntryJSON("role/"+name, &roleStorageEntry{
		DB:           roleDB,
		MongoDBRoles: roles,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (roles mongodbRoles) toStandardRolesArray() []interface{} {
	// Convert array of role documents like:
	//
	// [ { "role": "readWrite" }, { "role": "readWrite", "db": "test" } ]
	//
	// into a "standard" MongoDB roles array containing both strings and role documents:
	//
	// [ "readWrite", { "role": "readWrite", "db": "test" } ]
	//
	// MongoDB's createUser command accepts the latter.
	//
	var standardRolesArray []interface{}
	for _, role := range roles {
		if role.DB == "" {
			standardRolesArray = append(standardRolesArray, role.Role)
		} else {
			standardRolesArray = append(standardRolesArray, role)
		}
	}
	return standardRolesArray
}

type roleStorageEntry struct {
	DB           string       `json:"db"`
	MongoDBRoles mongodbRoles `json:"roles"`
}

type mongodbRole struct {
	Role string `json:"role" bson:"role"`
	DB   string `json:"db"   bson:"db"`
}

type mongodbRoles []mongodbRole

const pathRoleHelpSyn = `
Manage the roles used to generate MongoDB credentials.
`

const pathRoleHelpDesc = `
This path lets you manage the roles used to generate MongoDB credentials.

The "db" parameter specifies the authentication database for users
generated for a given role.

The "roles" parameter specifies the MongoDB roles that should be assigned
to users created for a given role. Just like when creating a user directly
using db.createUser, the roles JSON array can specify both built-in roles
and user-defined roles for both the database the user is created in and
for other databases.

For example, the following roles JSON array grants the "readWrite"
permission on both the user's authentication database and the "test"
database:

[ "readWrite", { "role": "readWrite", "db": "test" } ]

Please consult the MongoDB documentation for more
details on Role-Based Access Control in MongoDB.
`
