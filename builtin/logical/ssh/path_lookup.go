package ssh

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLookup(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "lookup",
		Fields: map[string]*framework.FieldSchema{
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "[Required] IP address of remote host",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLookupWrite,
		},
		HelpSynopsis:    pathLookupSyn,
		HelpDescription: pathLookupDesc,
	}
}

func (b *backend) pathLookupWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ipAddr := d.Get("ip").(string)
	if ipAddr == "" {
		return logical.ErrorResponse("Missing ip"), nil
	}
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return logical.ErrorResponse(fmt.Sprintf("Invalid IP %q", ip.String())), nil
	}

	// Get all the roles created in the backend.
	keys, err := req.Storage.List(ctx, "roles/")
	if err != nil {
		return nil, err
	}

	// Look for roles which has CIDR blocks that encompasses the given IP
	// and create a list out of it.
	var matchingRoles []string
	for _, role := range keys {
		if contains, _ := roleContainsIP(ctx, req.Storage, role, ip.String()); contains {
			matchingRoles = append(matchingRoles, role)
		}
	}

	// Add roles that are allowed to accept any IP address.
	zeroAddressEntry, err := b.getZeroAddressRoles(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if zeroAddressEntry != nil {
		matchingRoles = append(matchingRoles, zeroAddressEntry.Roles...)
	}

	// This list may potentially reveal more information than it is supposed to.
	// The roles for which the client is not authorized to will also be displayed.
	// However, if the client tries to use the role for which the client is not
	// authenticated, it will fail. It is not a problem. In a way this can be
	// viewed as a feature. The client can ask for permissions to be given for
	// a specific role if things are not working!
	//
	// Ideally, the role names should be filtered and only the roles which
	// the client is authorized to see, should be returned.
	return &logical.Response{
		Data: map[string]interface{}{
			"roles": matchingRoles,
		},
	}, nil
}

const pathLookupSyn = `
List all the roles associated with the given IP address.
`

const pathLookupDesc = `
The IP address for which the key is requested, is searched in the CIDR blocks
registered with vault using the 'roles' endpoint. Keys can be generated only by
specifying the 'role' name. The roles that can be used to generate the key for
a particular IP, are listed via this endpoint. For example, if this backend is
mounted at "ssh", then "ssh/lookup" lists the roles associated with keys can be
generated for a target IP, if the CIDR block encompassing the IP is registered
with vault.
`
