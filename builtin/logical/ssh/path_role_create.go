package ssh

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoleCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/(?P<name>\\w+)",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "name of the policy",
			},
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "username in target",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IP of the target machine",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathRoleCreateWrite,
		},
		HelpSynopsis:    pathRoleCreateHelpSyn,
		HelpDescription: pathRoleCreateHelpDesc,
	}
}

func (b *backend) pathRoleCreateWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)
	username := d.Get("username").(string)
	ipRaw := d.Get("ip").(string)
	if roleName == "" {
		return logical.ErrorResponse("Invalid 'name'"), nil
	}
	if ipRaw == "" {
		return logical.ErrorResponse("Invalid 'ip'"), nil
	}

	//find the role to be used for installing dynamic key
	roleEntry, err := req.Storage.Get(fmt.Sprintf("policy/%s", roleName))
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %s", err)
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("Role '%s' not found", roleName)), nil
	}
	var role sshRole
	if err := roleEntry.DecodeJSON(&role); err != nil {
		return nil, err
	}

	if username == "" {
		username = role.DefaultUser
	}

	//validate the IP address
	ipAddr := net.ParseIP(ipRaw)
	if ipAddr == nil {
		return logical.ErrorResponse(fmt.Sprintf("Invalid IP '%s'", ipRaw)), nil
	}
	ip := ipAddr.String()

	ipMatched := false
	for _, item := range strings.Split(role.CIDR, ",") {
		_, cidrIPNet, err := net.ParseCIDR(item)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid cidr entry '%s'", item)), nil
		}
		ipMatched = cidrIPNet.Contains(ipAddr)
		if ipMatched {
			break
		}
	}
	if !ipMatched {
		return logical.ErrorResponse(fmt.Sprintf("IP[%s] does not belong to role[%s]", ip, roleName)), nil
	}

	//fetch the host key to be used for installation
	keyEntry, err := req.Storage.Get(fmt.Sprintf("keys/%s", role.KeyName))
	if err != nil {
		return nil, fmt.Errorf("Key '%s' not found error:%s", role.KeyName, err)
	}
	var hostKey sshHostKey
	if err := keyEntry.DecodeJSON(&hostKey); err != nil {
		return nil, fmt.Errorf("Error reading the host key: %s", err)
	}

	//store the host key to file. Use it as parameter for scp command
	hostKeyFileName := fmt.Sprintf("./vault_ssh_%s_%s_shared.pem", username, ip)
	err = ioutil.WriteFile(hostKeyFileName, []byte(hostKey.Key), 0600)

	dynamicPrivateKeyFileName := fmt.Sprintf("vault_ssh_%s_%s_otk.pem", username, ip)
	dynamicPublicKeyFileName := fmt.Sprintf("vault_ssh_%s_%s_otk.pem.pub", username, ip)

	//delete the temporary files if they are already present
	err = removeFile(dynamicPrivateKeyFileName)
	if err != nil {
		return nil, fmt.Errorf("Error removing dynamic private key file: '%s'", err)
	}
	err = removeFile(dynamicPublicKeyFileName)
	if err != nil {
		return nil, fmt.Errorf("Error removing dynamic private key file: '%s'", err)
	}

	//generate RSA key pair
	dynamicPublicKey, dynamicPrivateKey, _ := generateRSAKeys()

	//save the public key pair to a file
	ioutil.WriteFile(dynamicPublicKeyFileName, []byte(dynamicPublicKey), 0644)

	//send the public key to target machine
	err = uploadFileScp(dynamicPublicKeyFileName, username, ip, hostKey.Key)
	if err != nil {
		return nil, err
	}

	//connect to target machine
	session, err := createSSHPublicKeysSession(username, ip, hostKey.Key)
	if err != nil {
		return nil, fmt.Errorf("Unable to create SSH Session using public keys: %s", err)
	}
	if session == nil {
		return nil, fmt.Errorf("Invalid session object")
	}

	authKeysFileName := fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)
	tempKeysFileName := fmt.Sprintf("/home/%s/temp_authorized_keys", username)

	//commands to be run on target machine
	grepCmd := fmt.Sprintf("grep -vFf %s %s > %s", dynamicPublicKeyFileName, authKeysFileName, tempKeysFileName)
	catCmdRemoveDuplicate := fmt.Sprintf("cat %s > %s", tempKeysFileName, authKeysFileName)
	catCmdAppendNew := fmt.Sprintf("cat %s >> %s", dynamicPublicKeyFileName, authKeysFileName)
	removeCmd := fmt.Sprintf("rm -f %s %s", tempKeysFileName, dynamicPublicKeyFileName)

	targetCmd := fmt.Sprintf("%s;%s;%s;%s", grepCmd, catCmdRemoveDuplicate, catCmdAppendNew, removeCmd)

	//run the commands on target machine
	if err := session.Run(targetCmd); err != nil {
		return nil, err
	}
	session.Close()

	result := b.Secret(SecretOneTimeKeyType).Response(map[string]interface{}{
		"key": dynamicPrivateKey,
	}, map[string]interface{}{
		"username":           username,
		"ip":                 ip,
		"host_key_name":      role.KeyName,
		"dynamic_public_key": dynamicPublicKey,
	})
	lease, _ := b.Lease(req.Storage)
	if lease != nil {
		result.Secret.Lease = lease.Lease
		result.Secret.LeaseGracePeriod = lease.LeaseMax
	}
	return result, nil
}

type sshCIDR struct {
	CIDR []string
}

const pathRoleCreateHelpSyn = `
Creates a dynamic key for the target machine.
`

const pathRoleCreateHelpDesc = `
This path will generates a new key for establishing SSH session with
target host. Previously registered shared key belonging to target
infrastructure will be used to install the new key at the target. If
this backend is mounted at 'ssh', then "ssh/creds/role" would generate
a dynamic key for 'web' role.

The dynamic keys will have a lease associated with them. The access
keys can be revoked by using the lease ID.
`
