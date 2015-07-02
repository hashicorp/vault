package ssh

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLookup(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "lookup",
		Fields: map[string]*framework.FieldSchema{
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IP address of target",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathLookupWrite,
		},
		HelpSynopsis:    pathLookupSyn,
		HelpDescription: pathLookupDesc,
	}
}

func (b *backend) pathLookupWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ipAddr := d.Get("ip").(string)
	if ipAddr == "" {
		return logical.ErrorResponse("Invalid 'ip'"), nil
	}
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return logical.ErrorResponse(fmt.Sprintf("Invalid IP '%s'", ip.String())), nil
	}

	keys, err := req.Storage.List("policy/")
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return &logical.Response{
			Data: map[string]interface{}{
				"roles": nil,
			},
		}, nil
	}

	var matchingRoles []string
	for _, item := range keys {
		if contains, _ := containsIP(req.Storage, item, ip.String()); contains {
			matchingRoles = append(matchingRoles, item)
		}
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"roles": matchingRoles,
		},
	}, nil
}

func containsIP(s logical.Storage, roleName string, ip string) (bool, error) {
	if roleName == "" || ip == "" {
		return false, fmt.Errorf("invalid parameters")
	}
	roleEntry, err := s.Get(fmt.Sprintf("policy/%s", roleName))
	if err != nil {
		return false, fmt.Errorf("error retrieving role '%s'", err)
	}
	if roleEntry == nil {
		return false, fmt.Errorf("role '%s' not found", roleName)
	}
	var role sshRole
	if err := roleEntry.DecodeJSON(&role); err != nil {
		return false, fmt.Errorf("error decoding role '%s'", roleName)
	}
	ipMatched := false
	for _, item := range strings.Split(role.CIDR, ",") {
		_, cidrIPNet, err := net.ParseCIDR(item)
		if err != nil {
			return false, fmt.Errorf("invalid cidr entry '%s'", item)
		}
		ipMatched = cidrIPNet.Contains(net.ParseIP(ip))
		if ipMatched {
			break
		}
	}
	return ipMatched, nil
}

const pathLookupSyn = `
Lists 'roles' that can be used to create a dynamic key.
`

const pathLookupDesc = `
The IP address to which the key is requested will be searched in
the CIDR blocks registered with vault using the 'roles' endpoint.
Key can be generated only by specifying the 'role' name. The roles
that can be used to generate the key for a particular IP will be
listed through this endpoint. For example, if this backend is mounted
at "ssh" then "ssh/lookup" lists the roles associated with 
Keys can be generated for a target IP if the CIDR block encompassing
the IP is registered with vault.
`
