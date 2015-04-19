package postgresql

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	_ "github.com/lib/pq"
)

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/(?P<name>\\w+)",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},

			"sql": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "SQL string to create a user. See help for more info.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathRoleCreate,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	sql := data.Get("sql").(string)

	// Get our connection
	db, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	// Test the query by trying to prepare it
	stmt, err := db.Prepare(Query(sql, map[string]string{
		"name":       "foo",
		"password":   "bar",
		"expiration": "",
	}))
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error testing query: %s", err)), nil
	}
	stmt.Close()

	// Store it
	entry, err := logical.StorageEntryJSON("role/"+name, map[string]interface{}{
		"sql": sql,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.

The "sql" parameter customizes the SQL string used to create the role.
This can only be a single SQL query. Some substitution will be done to the
SQL string for certain keys. The names of the variables must be surrounded
by "{{" and "}}" to be replaced.

  * "name" - The random username generated for the DB user.

  * "password" - The random password generated for the DB user.

  * "expiration" - The timestamp when this user will expire.

Example of a decent SQL query to use:

  CREATE ROLE "{{name}}" WITH
    LOGIN
    PASSWORD '{{password}}'
    VALID UNTIL '{{expiration}}';

Note the above user wouldn't be able to access anything. To give a user access
to resources, create roles manually in PostgreSQL, then use the "IN ROLE"
clause for CREATE ROLE to add the user to more roles.
`
