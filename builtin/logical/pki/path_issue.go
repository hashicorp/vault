package pki

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The desired role with configuration for this
request`,
			},
			"common_name": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The requested common name; if you want more than
one, specify the alternative names in the
alt_names map`,
			},
			"alt_names": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The requested Subject Alternative Names, if any,
in a comma-delimited list`,
			},
			"ip_sans": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The requested IP SANs, if any, in a
common-delimited list`,
			},
			"lease": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `The requested lease. DEPRECATED: use "ttl" instead.`,
			},
			"ttl": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The requested Time To Live for the certificate;
sets the expiration date. If not specified
the role default TTL it used. Cannot be larer
than the role max TTL.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathIssueCert,
		},

		HelpSynopsis:    pathIssueCertHelpSyn,
		HelpDescription: pathIssueCertHelpDesc,
	}
}

func (b *backend) pathIssueCert(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the common name(s)
	var commonNames []string
	cn := data.Get("common_name").(string)
	if len(cn) == 0 {
		return logical.ErrorResponse("The common_name field is required"), nil
	}
	commonNames = []string{cn}

	cnAlt := data.Get("alt_names").(string)
	if len(cnAlt) != 0 {
		for _, v := range strings.Split(cnAlt, ",") {
			commonNames = append(commonNames, v)
		}
	}

	// Get the role
	role, err := b.getRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", roleName)), nil
	}

	// Get any IP SANs
	ipSANs := []net.IP{}

	ipAlt := data.Get("ip_sans").(string)
	if len(ipAlt) != 0 {
		if !role.AllowIPSANs {
			return logical.ErrorResponse(fmt.Sprintf("IP Subject Alternative Names are not allowed in this role, but was provided %s", ipAlt)), nil
		}
		for _, v := range strings.Split(ipAlt, ",") {
			parsedIP := net.ParseIP(v)
			if parsedIP == nil {
				return logical.ErrorResponse(fmt.Sprintf("The value '%s' is not a valid IP address", v)), nil
			}
			ipSANs = append(ipSANs, parsedIP)
		}
	}

	ttlField := data.Get("ttl").(string)
	if len(ttlField) == 0 {
		ttlField = data.Get("lease").(string)
		if len(ttlField) == 0 {
			ttlField = role.TTL
		}
	}

	var ttl time.Duration
	if len(ttlField) == 0 {
		ttl, err = b.System().DefaultLeaseTTL()
		if err != nil {
			return nil, fmt.Errorf("Error fetching default TTL: %s", err)
		}
	} else {
		ttl, err = time.ParseDuration(ttlField)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Invalid requested ttl: %s", err)), nil
		}
	}

	var maxTTL time.Duration
	if len(role.MaxTTL) == 0 {
		maxTTL, err = b.System().MaxLeaseTTL()
		if err != nil {
			return nil, fmt.Errorf("Error fetching max TTL: %s", err)
		}
	} else {
		maxTTL, err = time.ParseDuration(role.MaxTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Invalid ttl: %s", err)), nil
		}
	}

	if ttl > maxTTL {
		// Don't error if they were using system defaults, only error if
		// they specifically chose a bad TTL
		if len(ttlField) == 0 {
			ttl = maxTTL
		} else {
			return logical.ErrorResponse("TTL is larger than maximum allowed by this role"), nil
		}
	}

	badName, err := validateCommonNames(req, commonNames, role)
	if len(badName) != 0 {
		return logical.ErrorResponse(fmt.Sprintf("Name %s not allowed by this role", badName)), nil
	} else if err != nil {
		return nil, fmt.Errorf("Error validating name %s: %s", badName, err)
	}

	signingBundle, caErr := fetchCAInfo(req)
	switch caErr.(type) {
	case certutil.UserError:
		return logical.ErrorResponse(fmt.Sprintf("Could not fetch the CA certificate: %s", caErr)), nil
	case certutil.InternalError:
		return nil, fmt.Errorf("Error fetching CA certificate: %s", caErr)
	}

	if time.Now().Add(ttl).After(signingBundle.Certificate.NotAfter) {
		return logical.ErrorResponse(fmt.Sprintf("Cannot satisfy request, as TTL is beyond the expiration of the CA certificate")), nil
	}

	var usage certUsage
	if role.ServerFlag {
		usage = usage | serverUsage
	}
	if role.ClientFlag {
		usage = usage | clientUsage
	}
	if role.CodeSigningFlag {
		usage = usage | codeSigningUsage
	}

	creationBundle := &certCreationBundle{
		SigningBundle: signingBundle,
		CACert:        signingBundle.Certificate,
		CommonNames:   commonNames,
		IPSANs:        ipSANs,
		KeyType:       role.KeyType,
		KeyBits:       role.KeyBits,
		TTL:           ttl,
		Usage:         usage,
	}

	parsedBundle, err := createCertificate(creationBundle)
	switch err.(type) {
	case certutil.UserError:
		return logical.ErrorResponse(err.Error()), nil
	case certutil.InternalError:
		return nil, err
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("Error converting raw cert bundle to cert bundle: %s", err)
	}

	resp := b.Secret(SecretCertsType).Response(
		structs.New(cb).Map(),
		map[string]interface{}{
			"serial_number": cb.SerialNumber,
		})

	resp.Secret.TTL = ttl

	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "certs/" + cb.SerialNumber,
		Value: parsedBundle.CertificateBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to store certificate locally")
	}

	return resp, nil
}

const pathIssueCertHelpSyn = `
Request certificates using a certain role with the provided common name.
`

const pathIssueCertHelpDesc = `
This path allows requesting certificates to be issued according to the
policy of the given role. The certificate will only be issued if the
requested common name is allowed by the role policy.
`
