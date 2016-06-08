package azuresql

import (
	"fmt"

	"strings"

	"github.com/Azure/azure-sdk-for-go/management/sql"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCredsCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("ip"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IP parameter to create a firewall rule.",
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
	ip := data.Get("ip").(string)

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

	// Generate our username and password
	displayName := req.DisplayName
	if len(displayName) > 10 {
		displayName = displayName[:10]
	}
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	username := fmt.Sprintf("%s-%s", displayName, userUUID)
	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// If there is an IP specified, generate a firewall rule for that IP address
	fwrule := ip
	if len(ip) > 0 {
		client, err := b.AzureClient(req.Storage)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Azure Subscription error: %s", err)), nil
		}
		fwrule = strings.Join([]string{userUUID, ip}, "-")
		err = client.CreateFirewallRule(b.server, sql.FirewallRuleCreateParams{Name: fwrule, StartIPAddress: ip, EndIPAddress: ip})
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Unable to create firewall rule: %s", err)), nil
		}
	}

	// Get our connection
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

	// Execute each query
	for _, query := range SplitSQL(role.SQL) {
		stmt, err := db.Prepare(Query(query, map[string]string{
			"name":     username,
			"password": password,
		}))
		if err != nil {
			return nil, err
		}
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
		"username": username,
		"password": password,
		"fwrule":   fwrule,
	}, map[string]interface{}{
		"username": username,
		"fwrule":   fwrule,
	})
	resp.Secret.TTL = leaseConfig.TTL
	return resp, nil
}

const pathCredsCreateHelpSyn = `
Request database credentials for a certain role.
`

const pathCredsCreateHelpDesc = `
This path reads database credentials for a certain role. The
database credentials will be generated on demand and will be automatically
revoked when the lease is up. 

If role is generated with /creds/<role-name>/<ip-address>, an Azure firewall 
allow rule will be created for the given IP address (config/subscription needs 
to be configured).
`
