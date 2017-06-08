package hana

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCredsCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredsCreateRead,
		},

		HelpSynopsis:    pathCredsCreateHelpSyn,
		HelpDescription: pathCredsCreateHelpDesc,
	}
}

func (b *backend) pathCredsCreateRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Get the role
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	// Determine if we have a lease configuration
	leaseConfig, err := b.LeaseConfig(req.Storage)
	if err != nil {
		return nil, err
	}
	if leaseConfig == nil {
		leaseConfig = &configLease{}
	}

	// Generate username and password for new user
	displayName := req.DisplayName
	if len(displayName) > 32 {
		displayName = displayName[:32]
	}
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	username := fmt.Sprintf("%s_%s", displayName, userUUID)
	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// HANA does not allow hyphens in usernames
	username = strings.Replace(username, "-", "_", -1)

	// Tokens have a token prefix. Change this to Vault for HANA user management clarity
	username = strings.Replace(username, "token", "vault", 1)
	username = strings.ToUpper(username)

	// Most HANA configurations have password constraints.
	// Prefix with A1a to cover the base case (user must change password upon login anyway)
	password = strings.Replace(password, "-", "_", -1)
	password = "A1a" + password

	// Get our handle
	db, err := b.DB(req.Storage)
	if err != nil {
		return nil, err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Request HANA server time plus lease duration.
	// This ensures the created account is deactivated server-side upon lease revocation
	var validUntil string
	timeQuery := fmt.Sprintf("SELECT TO_NVARCHAR(add_seconds(CURRENT_TIMESTAMP,"+
		"%f), 'YYYY-MM-DD HH24:MI:SS') FROM DUMMY", (leaseConfig.TTL).Seconds())
	err = db.QueryRow(timeQuery).Scan(&validUntil)
	if err != nil {
		return nil, err
	}

	// Execute each query
	for _, query := range strutil.ParseArbitraryStringSlice(role.SQL, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		stmt, err := tx.Prepare(Query(query, map[string]string{
			"name":        username,
			"password":    password,
			"valid_until": validUntil,
		}))
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		if _, err := stmt.Exec(); err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the secret
	resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"username":    username,
		"password":    password,
		"valid_until": validUntil,
	}, map[string]interface{}{
		"username": username,
	})

	resp.Secret.TTL = leaseConfig.TTL
	return resp, nil
}

const pathCredsCreateHelpSyn = `
Request database credentials for a certain role.
`

const pathCredsCreateHelpDesc = `
This path reads database credentials for a certain role. The
database credentials will be generated on demand.
`
