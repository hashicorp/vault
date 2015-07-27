package ssh

import (
	"fmt"
	"net"

	"github.com/hashicorp/vault/helper/uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCredsCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/(?P<name>[-\\w]+)",
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
			logical.WriteOperation: b.pathCredsCreateWrite,
		},
		HelpSynopsis:    pathCredsCreateHelpSyn,
		HelpDescription: pathCredsCreateHelpDesc,
	}
}

func (b *backend) pathCredsCreateWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("Missing name"), nil
	}

	username := d.Get("username").(string)

	ipRaw := d.Get("ip").(string)
	if ipRaw == "" {
		return logical.ErrorResponse("Missing ip"), nil
	}

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

	// Set the default username
	if username == "" {
		if role.DefaultUser == "" {
			return logical.ErrorResponse("No default username registered. Use 'username' option"), nil
		}
		username = role.DefaultUser
	}

	// Validate the IP address
	ipAddr := net.ParseIP(ipRaw)
	if ipAddr == nil {
		return logical.ErrorResponse(fmt.Sprintf("Invalid IP '%s'", ipRaw)), nil
	}
	ip := ipAddr.String()
	ipMatched, err := cidrContainsIP(ip, role.CIDR)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("Error validating IP: %s", err)), nil
	}
	if !ipMatched {
		return logical.ErrorResponse(fmt.Sprintf("IP[%s] does not belong to role[%s]", ip, roleName)), nil
	}

	var result *logical.Response
	if role.KeyType == KeyTypeOTP {
		// Generate salted OTP
		otp := uuid.GenerateUUID()
		otpSalted := b.salt.SaltID(otp)
		entry, err := req.Storage.Get("otp/" + otpSalted)
		// Make sure that new OTP is not replacing an existing one
		for err == nil && entry != nil {
			otp := uuid.GenerateUUID()
			otpSalted := b.salt.SaltID(otp)
			entry, err = req.Storage.Get("otp/" + otpSalted)
		}
		entry, err = logical.StorageEntryJSON("otp/"+otpSalted, sshOTP{
			Username: username,
			IP:       ip,
		})
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(entry); err != nil {
			return nil, err
		}
		result = b.Secret(SecretOTPType).Response(map[string]interface{}{
			"key_type": role.KeyType,
			"key":      otp,
		}, map[string]interface{}{
			"otp": otp,
		})
	} else if role.KeyType == KeyTypeDynamic {
		// Fetch the host key to be used for dynamic key installation
		keyEntry, err := req.Storage.Get(fmt.Sprintf("keys/%s", role.KeyName))
		if err != nil {
			return nil, fmt.Errorf("key '%s' not found error:%s", role.KeyName, err)
		}

		if keyEntry == nil {
			return nil, fmt.Errorf("key '%s' not found", role.KeyName, err)
		}

		var hostKey sshHostKey
		if err := keyEntry.DecodeJSON(&hostKey); err != nil {
			return nil, fmt.Errorf("error reading the host key: %s", err)
		}

		// Generate RSA key pair
		dynamicPublicKey, dynamicPrivateKey, err := generateRSAKeys()
		if err != nil {
			return nil, fmt.Errorf("error generating key: %s", err)
		}

		// Transfer the public key to target machine
		err = uploadPublicKeyScp(dynamicPublicKey, username, ip, role.Port, hostKey.Key)
		if err != nil {
			return nil, err
		}

		// Add the public key to authorized_keys file in target machine
		err = installPublicKeyInTarget(role.AdminUser, username, ip, role.Port, hostKey.Key)
		if err != nil {
			return nil, fmt.Errorf("error adding public key to authorized_keys file in target")
		}

		result = b.Secret(SecretDynamicKeyType).Response(map[string]interface{}{
			"key":      dynamicPrivateKey,
			"key_type": role.KeyType,
		}, map[string]interface{}{
			"admin_user":         role.AdminUser,
			"username":           username,
			"ip":                 ip,
			"host_key_name":      role.KeyName,
			"dynamic_public_key": dynamicPublicKey,
			"port":               role.Port,
		})
	} else {
		return nil, fmt.Errorf("key type unknown")
	}

	// Change the lease information to reflect user's choice
	lease, _ := b.Lease(req.Storage)
	if lease != nil {
		result.Secret.Lease = lease.Lease
		result.Secret.LeaseGracePeriod = lease.LeaseMax
	}
	return result, nil
}

type sshOTP struct {
	Username string `json:"username"`
	IP       string `json:"ip"`
}

type sshCIDR struct {
	CIDR []string
}

const pathCredsCreateHelpSyn = `
Creates a dynamic key for the target machine.
`

const pathCredsCreateHelpDesc = `
This path will generate a new key for establishing SSH session with
target host. The key can either be a long lived dynamic key or a One
Time Password (OTP), using 'key_type' parameter being 'dynamic' or 
'otp' respectively. For dynamic keys, a named key should be supplied.
Create named key using the 'keys/' endpoint, and this represents the
shared SSH key of target host. If this backend is mounted at 'ssh',
then "ssh/creds/web" would generate a key for 'web' role.

Keys will have a lease associated with them. The access keys can be
revoked by using the lease ID.
`
