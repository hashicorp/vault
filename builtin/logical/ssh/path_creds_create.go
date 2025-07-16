// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"fmt"
	"net"
	"strings"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type sshOTP struct {
	Username string `json:"username" structs:"username" mapstructure:"username"`
	IP       string `json:"ip" structs:"ip" mapstructure:"ip"`
	RoleName string `json:"role_name" structs:"role_name" mapstructure:"role_name"`
}

func pathCredsCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameWithAtRegex("role"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixSSH,
			OperationVerb:   "generate",
			OperationSuffix: "credentials",
		},

		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: "[Required] Name of the role",
			},
			"username": {
				Type:        framework.TypeString,
				Description: "[Optional] Username in remote host",
			},
			"ip": {
				Type:        framework.TypeString,
				Description: "[Required] IP of the remote host",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathCredsCreateWrite,
		},

		HelpSynopsis:    pathCredsCreateHelpSyn,
		HelpDescription: pathCredsCreateHelpDesc,
	}
}

func (b *backend) pathCredsCreateWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("role").(string)
	if roleName == "" {
		return logical.ErrorResponse("Missing role"), nil
	}

	ipRaw := d.Get("ip").(string)
	if ipRaw == "" {
		return logical.ErrorResponse("Missing ip"), nil
	}

	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Role %q not found", roleName)), nil
	}

	// username is an optional parameter.
	username := d.Get("username").(string)

	// Set the default username
	if username == "" {
		if role.DefaultUser == "" {
			return logical.ErrorResponse("No default username registered. Use 'username' option"), nil
		}
		username = role.DefaultUser
	}

	if role.AllowedUsers != "" {
		// Check if the username is present in allowed users list.
		err := validateUsername(username, role.AllowedUsers)

		// If username is not present in allowed users list, check if it
		// is the default username in the role. If neither is true, then
		// that username is not allowed to generate a credential.
		if err != nil && username != role.DefaultUser {
			return logical.ErrorResponse("Username is not present is allowed users list"), nil
		}
	} else if username != role.DefaultUser {
		return logical.ErrorResponse("Username has to be either in allowed users list or has to be a default username"), nil
	}

	// Validate the IP address
	ipAddr := net.ParseIP(ipRaw)
	if ipAddr == nil {
		return logical.ErrorResponse(fmt.Sprintf("Invalid IP %q", ipRaw)), nil
	}

	// Check if the IP belongs to the registered list of CIDR blocks under the role
	ip := ipAddr.String()

	zeroAddressEntry, err := b.getZeroAddressRoles(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error retrieving zero-address roles: %w", err)
	}
	var zeroAddressRoles []string
	if zeroAddressEntry != nil {
		zeroAddressRoles = zeroAddressEntry.Roles
	}

	err = validateIP(ip, roleName, role.CIDRList, role.ExcludeCIDRList, zeroAddressRoles)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("Error validating IP: %v", err)), nil
	}

	var result *logical.Response
	if role.KeyType == KeyTypeOTP {
		// Generate an OTP
		otp, err := b.GenerateOTPCredential(ctx, req, &sshOTP{
			Username: username,
			IP:       ip,
			RoleName: roleName,
		})
		if err != nil {
			return nil, err
		}

		// Return the information relevant to user of OTP type and save
		// the data required for later use in the internal section of secret.
		// In this case, saving just the OTP is sufficient since there is
		// no need to establish connection with the remote host.
		result = b.Secret(SecretOTPType).Response(map[string]interface{}{
			"key_type": role.KeyType,
			"key":      otp,
			"username": username,
			"ip":       ip,
			"port":     role.Port,
		}, map[string]interface{}{
			"otp": otp,
		})
	} else if role.KeyType == KeyTypeDynamic {
		return nil, fmt.Errorf("dynamic key types have been removed")
	} else {
		return nil, fmt.Errorf("key type unknown")
	}

	return result, nil
}

// Generates a UUID OTP and its salted value based on the salt of the backend.
func (b *backend) GenerateSaltedOTP(ctx context.Context) (string, string, error) {
	str, err := uuid.GenerateUUID()
	if err != nil {
		return "", "", err
	}
	salt, err := b.Salt(ctx)
	if err != nil {
		return "", "", err
	}

	return str, salt.SaltID(str), nil
}

// Generates an UUID OTP and creates an entry for the same in storage backend with its salted string.
func (b *backend) GenerateOTPCredential(ctx context.Context, req *logical.Request, sshOTPEntry *sshOTP) (string, error) {
	otp, otpSalted, err := b.GenerateSaltedOTP(ctx)
	if err != nil {
		return "", err
	}

	// Check if there is an entry already created for the newly generated OTP.
	entry, err := b.getOTP(ctx, req.Storage, otpSalted)

	// If entry already exists for the OTP, make sure that new OTP is not
	// replacing an existing one by recreating new ones until an unused
	// OTP is generated. It is very unlikely that this is the case and this
	// code is just for safety.
	for err == nil && entry != nil {
		otp, otpSalted, err = b.GenerateSaltedOTP(ctx)
		if err != nil {
			return "", err
		}
		entry, err = b.getOTP(ctx, req.Storage, otpSalted)
		if err != nil {
			return "", err
		}
	}

	// Store an entry for the salt of OTP.
	newEntry, err := logical.StorageEntryJSON("otp/"+otpSalted, sshOTPEntry)
	if err != nil {
		return "", err
	}
	if err := req.Storage.Put(ctx, newEntry); err != nil {
		return "", err
	}
	return otp, nil
}

// ValidateIP first checks if the role belongs to the list of privileged
// roles that could allow any IP address and if there is a match, IP is
// accepted immediately. If not, IP is searched in the allowed CIDR blocks
// registered with the role. If there is a match, then it is searched in the
// excluded CIDR blocks and if IP is found there as well, an error is returned.
// IP is valid only if it is encompassed by allowed CIDR blocks and not by
// excluded CIDR blocks.
func validateIP(ip, roleName, cidrList, excludeCidrList string, zeroAddressRoles []string) error {
	// Search IP in the zero-address list
	for _, role := range zeroAddressRoles {
		if roleName == role {
			return nil
		}
	}

	// Search IP in allowed CIDR blocks
	ipMatched, err := cidrListContainsIP(ip, cidrList)
	if err != nil {
		return err
	}
	if !ipMatched {
		return fmt.Errorf("IP does not belong to role")
	}

	if len(excludeCidrList) == 0 {
		return nil
	}

	// Search IP in exclude list
	ipMatched, err = cidrListContainsIP(ip, excludeCidrList)
	if err != nil {
		return err
	}
	if ipMatched {
		return fmt.Errorf("IP does not belong to role")
	}

	return nil
}

// Checks if the username supplied by the user is present in the list of
// allowed users registered which creation of role.
func validateUsername(username, allowedUsers string) error {
	if allowedUsers == "" {
		return fmt.Errorf("username not in allowed users list")
	}

	// Role was explicitly configured to allow any username.
	if allowedUsers == "*" {
		return nil
	}

	userList := strings.Split(allowedUsers, ",")
	for _, user := range userList {
		if strings.TrimSpace(user) == username {
			return nil
		}
	}

	return fmt.Errorf("username not in allowed users list")
}

const pathCredsCreateHelpSyn = `
Creates a credential for establishing SSH connection with the remote host.
`

const pathCredsCreateHelpDesc = `
This path will generate a new key for establishing SSH session with
target host. The key can be a One Time Password (OTP) using 'key_type'
being 'otp'.

Keys will have a lease associated with them. The access keys can be
revoked by using the lease ID.
`
