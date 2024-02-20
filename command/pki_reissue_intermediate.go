// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/posener/complete"
)

type PKIReIssueCACommand struct {
	*BaseCommand

	flagConfig          string
	flagReturnIndicator string
	flagDefaultDisabled bool
	flagList            bool

	flagKeyStorageSource string
	flagNewIssuerName    string
}

func (c *PKIReIssueCACommand) Synopsis() string {
	return "Uses a parent certificate and a template certificate to create a new issuer on a child mount"
}

func (c *PKIReIssueCACommand) Help() string {
	helpText := `
Usage: vault pki reissue PARENT TEMPLATE CHILD_MOUNT options
`
	return strings.TrimSpace(helpText)
}

func (c *PKIReIssueCACommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "type",
		Target:     &c.flagKeyStorageSource,
		Default:    "internal",
		EnvVar:     "",
		Usage:      `Options are “existing” - to use an existing key inside vault, “internal” - to generate a new key inside vault, or “kms” - to link to an external key.  Exported keys are not available through this API.`,
		Completion: complete.PredictSet("internal", "existing", "kms"),
	})

	f.StringVar(&StringVar{
		Name:    "issuer_name",
		Target:  &c.flagNewIssuerName,
		Default: "",
		EnvVar:  "",
		Usage:   `If present, the newly created issuer will be given this name`,
	})

	return set
}

func (c *PKIReIssueCACommand) Run(args []string) int {
	// Parse Args
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	args = f.Args()

	if len(args) < 3 {
		c.UI.Error("Not enough arguments: expected parent issuer and child-mount location and some key_value argument")
		return 1
	}

	stdin := (io.Reader)(os.Stdin)
	userData, err := parseArgsData(stdin, args[3:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	// Check We Have a Client
	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to obtain client: %v", err))
		return 1
	}

	parentIssuer := sanitizePath(args[0]) // /pki/issuer/default
	templateIssuer := sanitizePath(args[1])
	intermediateMount := sanitizePath(args[2])

	templateIssuerBundle, err := readIssuer(client, templateIssuer)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error fetching template certificate %v : %v", templateIssuer, err))
		return 1
	}
	certificate := templateIssuerBundle.certificate

	useExistingKey := c.flagKeyStorageSource == "existing"
	keyRef := ""
	if useExistingKey {
		keyRef = templateIssuerBundle.keyId

		if keyRef == "" {
			c.UI.Error(fmt.Sprintf("Template issuer %s did not have a key id field set in response which is required", templateIssuer))
			return 1
		}
	}

	templateData, err := parseTemplateCertificate(*certificate, useExistingKey, keyRef)
	data := updateTemplateWithData(templateData, userData)

	return pkiIssue(c.BaseCommand, parentIssuer, intermediateMount, c.flagNewIssuerName, c.flagKeyStorageSource, data)
}

func updateTemplateWithData(template map[string]interface{}, changes map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{}

	for key, value := range template {
		data[key] = value
	}

	// ttl and not_after set the same thing.  Delete template ttl if using not_after:
	if _, ok := changes["not_after"]; ok {
		delete(data, "ttl")
	}

	// If we are updating the key_type, do not set key_bits
	if _, ok := changes["key_type"]; ok && changes["key_type"] != template["key_type"] {
		delete(data, "key_bits")
	}

	for key, value := range changes {
		data[key] = value
	}

	return data
}

func parseTemplateCertificate(certificate x509.Certificate, useExistingKey bool, keyRef string) (templateData map[string]interface{}, err error) {
	// Generate Certificate Signing Parameters
	templateData = map[string]interface{}{
		"common_name": certificate.Subject.CommonName,
		"alt_names":   certutil.MakeAltNamesCommaSeparatedString(certificate.DNSNames, certificate.EmailAddresses),
		"ip_sans":     certutil.MakeIpAddressCommaSeparatedString(certificate.IPAddresses),
		"uri_sans":    certutil.MakeUriCommaSeparatedString(certificate.URIs),
		// other_sans (string: "") - Specifies custom OID/UTF8-string SANs. These must match values specified on the role in allowed_other_sans (see role creation for allowed_other_sans globbing rules). The format is the same as OpenSSL: <oid>;<type>:<value> where the only current valid type is UTF8. This can be a comma-delimited list or a JSON string slice.
		// Punting on Other_SANs, shouldn't really be on CAs
		"signature_bits":        certutil.FindSignatureBits(certificate.SignatureAlgorithm),
		"exclude_cn_from_sans":  certutil.DetermineExcludeCnFromCertSans(certificate),
		"ou":                    certificate.Subject.OrganizationalUnit,
		"organization":          certificate.Subject.Organization,
		"country":               certificate.Subject.Country,
		"locality":              certificate.Subject.Locality,
		"province":              certificate.Subject.Province,
		"street_address":        certificate.Subject.StreetAddress,
		"postal_code":           certificate.Subject.PostalCode,
		"serial_number":         certificate.Subject.SerialNumber,
		"ttl":                   (certificate.NotAfter.Sub(certificate.NotBefore)).String(),
		"max_path_length":       certificate.MaxPathLen,
		"permitted_dns_domains": strings.Join(certificate.PermittedDNSDomains, ","),
		"use_pss":               certutil.IsPSS(certificate.SignatureAlgorithm),
	}

	if useExistingKey {
		templateData["skid"] = hex.EncodeToString(certificate.SubjectKeyId) // TODO: Double Check this with someone
		if keyRef == "" {
			return nil, fmt.Errorf("unable to create certificate template for existing key without a key_id")
		}
		templateData["key_ref"] = keyRef
	} else {
		templateData["key_type"] = certutil.GetKeyType(certificate.PublicKeyAlgorithm.String())
		templateData["key_bits"] = certutil.FindBitLength(certificate.PublicKey)
	}

	return templateData, nil
}
