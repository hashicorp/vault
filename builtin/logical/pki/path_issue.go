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
		Pattern: `issue/(?P<role>\w[\w-]+\w)`,
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
				Description: "The requested lease",
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

	leaseField := data.Get("lease").(string)
	if len(leaseField) == 0 {
		leaseField = role.Lease
	}

	lease, err := time.ParseDuration(leaseField)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Invalid requested lease: %s", err)), nil
	}
	leaseMax, err := time.ParseDuration(role.LeaseMax)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Invalid lease: %s", err)), nil
	}

	if lease > leaseMax {
		return logical.ErrorResponse("Lease expires after maximum allowed by this role"), nil
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

	if time.Now().Add(lease).After(signingBundle.Certificate.NotAfter) {
		return logical.ErrorResponse(fmt.Sprintf("Cannot satisfy request, as maximum lease is beyond the expiration of the CA certificate")), nil
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
		Lease:         lease,
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

	resp.Secret.Lease = lease

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
