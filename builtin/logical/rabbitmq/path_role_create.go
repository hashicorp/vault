package rabbitmq

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	rabbithole "github.com/michaelklishin/rabbit-hole"
)

func pathCreds(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredsRead,
		},

		HelpSynopsis:    pathRoleCreateReadHelpSyn,
		HelpDescription: pathRoleCreateReadHelpDesc,
	}
}

// Issues the credential based on the role name
func (b *backend) pathCredsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	// Get the role
	role, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	// Ensure username is unique
	uuidVal, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	username := fmt.Sprintf("%s-%s", req.DisplayName, uuidVal)

	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// Get the client configuration
	client, err := b.Client(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return logical.ErrorResponse("failed to get the client"), nil
	}

	// Register the generated credentials in the backend, with the RabbitMQ server
	resp, err := client.PutUser(username, rabbithole.UserSettings{
		Password: password,
		Tags:     role.Tags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create a new user with the generated credentials")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.Logger().Error(fmt.Sprintf("unable to close response body: %s", err))
		}
	}()
	if !isIn200s(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error creating user %s - %d: %s", username, resp.StatusCode, body)
	}

	success := false
	defer func() {
		if success {
			return
		}
		// Delete the user because it's in an unknown state.
		resp, err := client.DeleteUser(username)
		if err != nil {
			b.Logger().Error(fmt.Sprintf("failed to delete %s: %s", username, err))
		}
		if !isIn200s(resp.StatusCode) {
			body, _ := ioutil.ReadAll(resp.Body)
			b.Logger().Error(fmt.Sprintf("error deleting %s - %d: %s", username, resp.StatusCode, body))
		}
	}()

	// If the role had vhost permissions specified, assign those permissions
	// to the created username for respective vhosts.
	for vhost, permission := range role.VHosts {
		if err := func() error {
			resp, err := client.UpdatePermissionsIn(vhost, username, rabbithole.Permissions{
				Configure: permission.Configure,
				Write:     permission.Write,
				Read:      permission.Read,
			})
			if err != nil {
				return err
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					b.Logger().Error(fmt.Sprintf("unable to close response body: %s", err))
				}
			}()
			if !isIn200s(resp.StatusCode) {
				body, _ := ioutil.ReadAll(resp.Body)
				return fmt.Errorf("error updating vhost permissions for %s - %d: %s", vhost, resp.StatusCode, body)
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}

	// If the role had vhost topic permissions specified, assign those permissions
	// to the created username for respective vhosts and exchange.
	for vhost, permissions := range role.VHostTopics {
		for exchange, permission := range permissions {
			if err := func() error {
				resp, err := client.UpdateTopicPermissionsIn(vhost, username, rabbithole.TopicPermissions{
					Exchange: exchange,
					Write:    permission.Write,
					Read:     permission.Read,
				})
				if err != nil {
					return err
				}
				defer func() {
					if err := resp.Body.Close(); err != nil {
						b.Logger().Error(fmt.Sprintf("unable to close response body: %s", err))
					}
				}()
				if !isIn200s(resp.StatusCode) {
					body, _ := ioutil.ReadAll(resp.Body)
					return fmt.Errorf("error updating vhost permissions for %s - %d: %s", vhost, resp.StatusCode, body)
				}
				return nil
			}(); err != nil {
				return nil, err
			}
		}
	}
	success = true

	// Return the secret
	response := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"username": username,
		"password": password,
	}, map[string]interface{}{
		"username": username,
	})

	// Determine if we have a lease
	lease, err := b.Lease(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if lease != nil {
		response.Secret.TTL = lease.TTL
		response.Secret.MaxTTL = lease.MaxTTL
	}

	return response, nil
}

func isIn200s(respStatus int) bool {
	return respStatus >= 200 && respStatus < 300
}

const pathRoleCreateReadHelpSyn = `
Request RabbitMQ credentials for a certain role.
`

const pathRoleCreateReadHelpDesc = `
This path reads RabbitMQ credentials for a certain role. The
RabbitMQ credentials will be generated on demand and will be automatically
revoked when the lease is up.
`
