package ssh

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLookup(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.pathLookup\n")
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
	ip := d.Get("ip").(string)
	//ip := "127.0.0.1"
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return logical.ErrorResponse(fmt.Sprintf("Invalid IP '%s'", ip)), nil
	}
	keys, _ := req.Storage.List("policy/")
	var matchingRoles []string
	for _, item := range keys {
		if contains, _ := containsIP(req.Storage, item, ip); contains {
			matchingRoles = append(matchingRoles, item)
		}
	}
	log.Printf("Vishal: req.Path: %#v \n Keys:%#v\n", req.Path, keys)
	return &logical.Response{
		Data: map[string]interface{}{
			"roles": matchingRoles,
		},
	}, nil
}

func containsIP(s logical.Storage, roleName string, ip string) (bool, error) {
	roleEntry, err := s.Get("policy/" + roleName)
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
		log.Println(item)
		_, cidrIPNet, _ := net.ParseCIDR(item)
		ipMatched = cidrIPNet.Contains(net.ParseIP(ip))
		if ipMatched {
			break
		}
	}
	return ipMatched, nil
}

const pathLookupSyn = `
pathLookupSyn
`

const pathLookupDesc = `
pathLoookupDesc
`
