package pki

import (
	"encoding/pem"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `issue/(?P<role>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The desired role with configuration for this request",
			},
			"common_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The requested common name; if you want more than one, specify the alternative names in the alt_names map",
			},
			"alt_names": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The requested Subject Alternative Names, if any, in a comma-delimited list",
			},
			"ip_sans": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The requested IP SANs, if any, in a common-delimited list",
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

	rawSigningBundle, caCert, err := fetchCAInfo(req)
	if err != nil {
		return logical.ErrorResponse("Could not fetch the CA certificate; has it been set?"), nil
	}

	if time.Now().Add(lease).After(caCert.NotAfter) {
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
		RawSigningBundle: rawSigningBundle,
		CACert:           caCert,
		CommonNames:      commonNames,
		IPSANs:           ipSANs,
		KeyType:          role.KeyType,
		KeyBits:          role.KeyBits,
		Lease:            lease,
		Usage:            usage,
	}

	rawBundle, userErr, intErr := createCertificate(creationBundle)
	switch {
	case userErr != nil:
		return logical.ErrorResponse(userErr.Error()), nil
	case intErr != nil:
		return nil, intErr
	}

	serial := strings.ToLower(getOctalFormatted(rawBundle.SerialNumber.Bytes(), ":"))

	resp := b.Secret(SecretCertsType).Response(map[string]interface{}{
		"serial": serial,
	}, map[string]interface{}{
		"serial": serial,
	})

	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rawBundle.CertificateBytes,
	}
	certificateString := string(pem.EncodeToMemory(&block))
	resp.Data["certificate"] = certificateString

	block.Bytes = rawSigningBundle.CertificateBytes
	caString := string(pem.EncodeToMemory(&block))
	resp.Data["issuing_ca"] = caString

	block.Bytes = rawBundle.PrivateKeyBytes
	switch rawBundle.PrivateKeyType {
	case RSAPrivateKeyType:
		block.Type = "RSA PRIVATE KEY"
	case ECPrivateKeyType:
		block.Type = "EC PRIVATE KEY"
	default:
		return nil, fmt.Errorf("Could not determine private key type when creating block")
	}
	resp.Data["private_key"] = string(pem.EncodeToMemory(&block))

	resp.Secret.Lease = lease

	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "certs/" + serial,
		Value: rawBundle.CertificateBytes,
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
