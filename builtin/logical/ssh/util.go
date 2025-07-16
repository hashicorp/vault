// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ssh"
)

// Creates a new RSA key pair with the given key length. The private key will be
// of pem format and the public key will be of OpenSSH format.
func generateRSAKeys(keyBits int) (publicKeyRsa string, privateKeyRsa string, err error) {
	privateKey, err := cryptoutil.GenerateRSAKey(rand.Reader, keyBits)
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key-pair: %w", err)
	}

	privateKeyRsa = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))

	sshPublicKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key-pair: %w", err)
	}
	publicKeyRsa = "ssh-rsa " + base64.StdEncoding.EncodeToString(sshPublicKey.Marshal())
	return
}

// Takes an IP address and role name and checks if the IP is part
// of CIDR blocks belonging to the role.
func roleContainsIP(ctx context.Context, s logical.Storage, roleName string, ip string) (bool, error) {
	if roleName == "" {
		return false, fmt.Errorf("missing role name")
	}

	if ip == "" {
		return false, fmt.Errorf("missing ip")
	}

	roleEntry, err := s.Get(ctx, fmt.Sprintf("roles/%s", roleName))
	if err != nil {
		return false, fmt.Errorf("error retrieving role %w", err)
	}
	if roleEntry == nil {
		return false, fmt.Errorf("role %q not found", roleName)
	}

	var role sshRole
	if err := roleEntry.DecodeJSON(&role); err != nil {
		return false, fmt.Errorf("error decoding role %q", roleName)
	}

	if matched, err := cidrListContainsIP(ip, role.CIDRList); err != nil {
		return false, err
	} else {
		return matched, nil
	}
}

// Returns true if the IP supplied by the user is part of the comma
// separated CIDR blocks
func cidrListContainsIP(ip, cidrList string) (bool, error) {
	if len(cidrList) == 0 {
		return false, fmt.Errorf("IP does not belong to role")
	}
	for _, item := range strings.Split(cidrList, ",") {
		_, cidrIPNet, err := net.ParseCIDR(item)
		if err != nil {
			return false, fmt.Errorf("invalid CIDR entry %q", item)
		}
		if cidrIPNet.Contains(net.ParseIP(ip)) {
			return true, nil
		}
	}
	return false, nil
}

func parsePublicSSHKey(key string) (ssh.PublicKey, error) {
	keyParts := strings.Split(key, " ")
	if len(keyParts) > 1 {
		// Someone has sent the 'full' public key rather than just the base64 encoded part that the ssh library wants
		key = keyParts[1]
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePublicKey([]byte(decodedKey))
}

func convertMapToStringValue(initial map[string]interface{}) map[string]string {
	result := map[string]string{}
	for key, value := range initial {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}

func convertMapToIntSlice(initial map[string]interface{}) (map[string][]int, error) {
	var err error
	result := map[string][]int{}

	for key, value := range initial {
		result[key], err = parseutil.SafeParseIntSlice(value, 0 /* no upper bound on number of keys lengths per key type */)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Serve a template processor for custom format inputs
func substQuery(tpl string, data map[string]string) string {
	for k, v := range data {
		tpl = strings.ReplaceAll(tpl, fmt.Sprintf("{{%s}}", k), v)
	}

	return tpl
}
