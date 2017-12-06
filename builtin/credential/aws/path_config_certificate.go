package awsauth

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// dsaSignature represents the contents of the signature of a signed
// content using digital signature algorithm.
type dsaSignature struct {
	R, S *big.Int
}

// This certificate is used to verify the PKCS#7 signature of the instance
// identity document. As per AWS documentation, this public key is valid for
// US East (N. Virginia), US West (Oregon), US West (N. California), EU
// (Ireland), EU (Frankfurt), Asia Pacific (Tokyo), Asia Pacific (Seoul), Asia
// Pacific (Singapore), Asia Pacific (Sydney), and South America (Sao Paulo).
//
// It's also the same certificate, but for some reason listed separately, for
// GovCloud (US)
const genericAWSPublicCertificatePkcs7 = `-----BEGIN CERTIFICATE-----
MIIC7TCCAq0CCQCWukjZ5V4aZzAJBgcqhkjOOAQDMFwxCzAJBgNVBAYTAlVTMRkw
FwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYD
VQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQzAeFw0xMjAxMDUxMjU2MTJaFw0z
ODAxMDUxMjU2MTJaMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9u
IFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNl
cnZpY2VzIExMQzCCAbcwggEsBgcqhkjOOAQBMIIBHwKBgQCjkvcS2bb1VQ4yt/5e
ih5OO6kK/n1Lzllr7D8ZwtQP8fOEpp5E2ng+D6Ud1Z1gYipr58Kj3nssSNpI6bX3
VyIQzK7wLclnd/YozqNNmgIyZecN7EglK9ITHJLP+x8FtUpt3QbyYXJdmVMegN6P
hviYt5JH/nYl4hh3Pa1HJdskgQIVALVJ3ER11+Ko4tP6nwvHwh6+ERYRAoGBAI1j
k+tkqMVHuAFcvAGKocTgsjJem6/5qomzJuKDmbJNu9Qxw3rAotXau8Qe+MBcJl/U
hhy1KHVpCGl9fueQ2s6IL0CaO/buycU1CiYQk40KNHCcHfNiZbdlx1E9rpUp7bnF
lRa2v1ntMX3caRVDdbtPEWmdxSCYsYFDk4mZrOLBA4GEAAKBgEbmeve5f8LIE/Gf
MNmP9CM5eovQOGx5ho8WqD+aTebs+k2tn92BBPqeZqpWRa5P/+jrdKml1qx4llHW
MXrs3IgIb6+hUIB+S8dz8/mmO0bpr76RoZVCXYab2CZedFut7qc3WUH9+EUAH5mw
vSeDCOUMYQR7R9LINYwouHIziqQYMAkGByqGSM44BAMDLwAwLAIUWXBlk40xTwSw
7HX32MxXYruse9ACFBNGmdX2ZBrVNGrN9N2f6ROk0k9K
-----END CERTIFICATE-----
`

// This certificate is used to verify the instance identity document using the
// RSA digest of the same
const genericAWSPublicCertificateIdentity = `-----BEGIN CERTIFICATE-----
MIIDIjCCAougAwIBAgIJAKnL4UEDMN/FMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNV
BAYTAlVTMRMwEQYDVQQIEwpXYXNoaW5ndG9uMRAwDgYDVQQHEwdTZWF0dGxlMRgw
FgYDVQQKEw9BbWF6b24uY29tIEluYy4xGjAYBgNVBAMTEWVjMi5hbWF6b25hd3Mu
Y29tMB4XDTE0MDYwNTE0MjgwMloXDTI0MDYwNTE0MjgwMlowajELMAkGA1UEBhMC
VVMxEzARBgNVBAgTCldhc2hpbmd0b24xEDAOBgNVBAcTB1NlYXR0bGUxGDAWBgNV
BAoTD0FtYXpvbi5jb20gSW5jLjEaMBgGA1UEAxMRZWMyLmFtYXpvbmF3cy5jb20w
gZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAIe9GN//SRK2knbjySG0ho3yqQM3
e2TDhWO8D2e8+XZqck754gFSo99AbT2RmXClambI7xsYHZFapbELC4H91ycihvrD
jbST1ZjkLQgga0NE1q43eS68ZeTDccScXQSNivSlzJZS8HJZjgqzBlXjZftjtdJL
XeE4hwvo0sD4f3j9AgMBAAGjgc8wgcwwHQYDVR0OBBYEFCXWzAgVyrbwnFncFFIs
77VBdlE4MIGcBgNVHSMEgZQwgZGAFCXWzAgVyrbwnFncFFIs77VBdlE4oW6kbDBq
MQswCQYDVQQGEwJVUzETMBEGA1UECBMKV2FzaGluZ3RvbjEQMA4GA1UEBxMHU2Vh
dHRsZTEYMBYGA1UEChMPQW1hem9uLmNvbSBJbmMuMRowGAYDVQQDExFlYzIuYW1h
em9uYXdzLmNvbYIJAKnL4UEDMN/FMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEF
BQADgYEAFYcz1OgEhQBXIwIdsgCOS8vEtiJYF+j9uO6jz7VOmJqO+pRlAbRlvY8T
C1haGgSI/A1uZUKs/Zfnph0oEI0/hu1IIJ/SKBDtN5lvmZ/IzbOPIJWirlsllQIQ
7zvWbGd9c9+Rm3p04oTvhup99la7kZqevJK0QRdD/6NpCKsqP/0=
-----END CERTIFICATE-----`

// pathListCertificates creates a path that enables listing of all
// the AWS public certificates registered with Vault.
func pathListCertificates(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/certificates/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathCertificatesList,
		},

		HelpSynopsis:    pathListCertificatesHelpSyn,
		HelpDescription: pathListCertificatesHelpDesc,
	}
}

func pathConfigCertificate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/certificate/" + framework.GenericNameRegex("cert_name"),
		Fields: map[string]*framework.FieldSchema{
			"cert_name": {
				Type:        framework.TypeString,
				Description: "Name of the certificate.",
			},
			"aws_public_cert": {
				Type:        framework.TypeString,
				Description: "Base64 encoded AWS Public cert required to verify PKCS7 signature of the EC2 instance metadata.",
			},
			"type": {
				Type:    framework.TypeString,
				Default: "pkcs7",
				Description: `
Takes the value of either "pkcs7" or "identity", indicating the type of
document which can be verified using the given certificate. The reason is that
the PKCS#7 document will have a DSA digest and the identity signature will have
an RSA signature, and accordingly the public certificates to verify those also
vary. Defaults to "pkcs7".`,
			},
		},

		ExistenceCheck: b.pathConfigCertificateExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigCertificateCreateUpdate,
			logical.UpdateOperation: b.pathConfigCertificateCreateUpdate,
			logical.ReadOperation:   b.pathConfigCertificateRead,
			logical.DeleteOperation: b.pathConfigCertificateDelete,
		},

		HelpSynopsis:    pathConfigCertificateSyn,
		HelpDescription: pathConfigCertificateDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigCertificateExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	certName := data.Get("cert_name").(string)
	if certName == "" {
		return false, fmt.Errorf("missing cert_name")
	}

	entry, err := b.lockedAWSPublicCertificateEntry(req.Storage, certName)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// pathCertificatesList is used to list all the AWS public certificates registered with Vault
func (b *backend) pathCertificatesList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	certs, err := req.Storage.List("config/certificate/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(certs), nil
}

// Decodes the PEM encoded certiticate and parses it into a x509 cert
func decodePEMAndParseCertificate(certificate string) (*x509.Certificate, error) {
	// Decode the PEM block and error out if a block is not detected in the first attempt
	decodedPublicCert, rest := pem.Decode([]byte(certificate))
	if len(rest) != 0 {
		return nil, fmt.Errorf("invalid certificate; should be one PEM block only")
	}

	// Check if the certificate can be parsed
	publicCert, err := x509.ParseCertificate(decodedPublicCert.Bytes)
	if err != nil {
		return nil, err
	}
	if publicCert == nil {
		return nil, fmt.Errorf("invalid certificate; failed to parse certificate")
	}
	return publicCert, nil
}

// awsPublicCertificates returns a slice of all the parsed AWS public
// certificates, which are used to verify either the SHA256 RSA signature, or
// the PKCS7 signatures of the instance identity documents. This method will
// append the certificates registered using `config/certificate/<cert_name>`
// endpoint, along with the default certificate in the backend.
func (b *backend) awsPublicCertificates(s logical.Storage, isPkcs bool) ([]*x509.Certificate, error) {
	// Lock at beginning and use internal method so that we are consistent as
	// we iterate through
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	var certs []*x509.Certificate

	defaultCert := genericAWSPublicCertificateIdentity
	if isPkcs {
		defaultCert = genericAWSPublicCertificatePkcs7
	}

	// Append the generic certificate provided in the AWS EC2 instance metadata documentation
	decodedCert, err := decodePEMAndParseCertificate(defaultCert)
	if err != nil {
		return nil, err
	}
	certs = append(certs, decodedCert)

	// Get the list of all the registered certificates
	registeredCerts, err := s.List("config/certificate/")
	if err != nil {
		return nil, err
	}

	// Iterate through each certificate, parse and append it to a slice
	for _, cert := range registeredCerts {
		certEntry, err := b.nonLockedAWSPublicCertificateEntry(s, cert)
		if err != nil {
			return nil, err
		}
		if certEntry == nil {
			return nil, fmt.Errorf("certificate storage has a nil entry under the name:%s\n", cert)
		}
		// Append relevant certificates only
		if (isPkcs && certEntry.Type == "pkcs7") ||
			(!isPkcs && certEntry.Type == "identity") {
			decodedCert, err := decodePEMAndParseCertificate(certEntry.AWSPublicCert)
			if err != nil {
				return nil, err
			}
			certs = append(certs, decodedCert)
		}
	}

	return certs, nil
}

// lockedSetAWSPublicCertificateEntry is used to store the AWS public key in
// the storage. This method acquires lock before creating or updating a storage
// entry.
func (b *backend) lockedSetAWSPublicCertificateEntry(s logical.Storage, certName string, certEntry *awsPublicCert) error {
	if certName == "" {
		return fmt.Errorf("missing certificate name")
	}

	if certEntry == nil {
		return fmt.Errorf("nil AWS public key certificate")
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return b.nonLockedSetAWSPublicCertificateEntry(s, certName, certEntry)
}

// nonLockedSetAWSPublicCertificateEntry is used to store the AWS public key in
// the storage. This method does not acquire lock before reading the storage.
// If locking is desired, use lockedSetAWSPublicCertificateEntry instead.
func (b *backend) nonLockedSetAWSPublicCertificateEntry(s logical.Storage, certName string, certEntry *awsPublicCert) error {
	if certName == "" {
		return fmt.Errorf("missing certificate name")
	}

	if certEntry == nil {
		return fmt.Errorf("nil AWS public key certificate")
	}

	entry, err := logical.StorageEntryJSON("config/certificate/"+certName, certEntry)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("failed to create storage entry for AWS public key certificate")
	}

	return s.Put(entry)
}

// lockedAWSPublicCertificateEntry is used to get the configured AWS Public Key
// that is used to verify the PKCS#7 signature of the instance identity
// document.
func (b *backend) lockedAWSPublicCertificateEntry(s logical.Storage, certName string) (*awsPublicCert, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedAWSPublicCertificateEntry(s, certName)
}

// nonLockedAWSPublicCertificateEntry reads the certificate information from
// the storage. This method does not acquire lock before reading the storage.
// If locking is desired, use lockedAWSPublicCertificateEntry instead.
func (b *backend) nonLockedAWSPublicCertificateEntry(s logical.Storage, certName string) (*awsPublicCert, error) {
	entry, err := s.Get("config/certificate/" + certName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var certEntry awsPublicCert
	if err := entry.DecodeJSON(&certEntry); err != nil {
		return nil, err
	}

	// Handle upgrade for certificate type
	persistNeeded := false
	if certEntry.Type == "" {
		certEntry.Type = "pkcs7"
		persistNeeded = true
	}

	if persistNeeded {
		if err := b.nonLockedSetAWSPublicCertificateEntry(s, certName, &certEntry); err != nil {
			return nil, err
		}
	}

	return &certEntry, nil
}

// pathConfigCertificateDelete is used to delete the previously configured AWS
// Public Key that is used to verify the PKCS#7 signature of the instance
// identity document.
func (b *backend) pathConfigCertificateDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	certName := data.Get("cert_name").(string)
	if certName == "" {
		return logical.ErrorResponse("missing cert_name"), nil
	}

	return nil, req.Storage.Delete("config/certificate/" + certName)
}

// pathConfigCertificateRead is used to view the configured AWS Public Key that
// is used to verify the PKCS#7 signature of the instance identity document.
func (b *backend) pathConfigCertificateRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	certName := data.Get("cert_name").(string)
	if certName == "" {
		return logical.ErrorResponse("missing cert_name"), nil
	}

	certificateEntry, err := b.lockedAWSPublicCertificateEntry(req.Storage, certName)
	if err != nil {
		return nil, err
	}
	if certificateEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(certificateEntry).Map(),
	}, nil
}

// pathConfigCertificateCreateUpdate is used to register an AWS Public Key that
// is used to verify the PKCS#7 signature of the instance identity document.
func (b *backend) pathConfigCertificateCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	certName := data.Get("cert_name").(string)
	if certName == "" {
		return logical.ErrorResponse("missing certificate name"), nil
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// Check if there is already a certificate entry registered
	certEntry, err := b.nonLockedAWSPublicCertificateEntry(req.Storage, certName)
	if err != nil {
		return nil, err
	}
	if certEntry == nil {
		certEntry = &awsPublicCert{}
	}

	// Check if type information is provided
	certTypeRaw, ok := data.GetOk("type")
	if ok {
		certEntry.Type = strings.ToLower(certTypeRaw.(string))
	} else if req.Operation == logical.CreateOperation {
		certEntry.Type = data.Get("type").(string)
	}

	switch certEntry.Type {
	case "pkcs7":
	case "identity":
	default:
		return logical.ErrorResponse(fmt.Sprintf("invalid certificate type %q", certEntry.Type)), nil
	}

	// Check if the value is provided by the client
	certStrData, ok := data.GetOk("aws_public_cert")
	if ok {
		if certBytes, err := base64.StdEncoding.DecodeString(certStrData.(string)); err == nil {
			certEntry.AWSPublicCert = string(certBytes)
		} else {
			certEntry.AWSPublicCert = certStrData.(string)
		}
	} else {
		// aws_public_cert should be supplied for both create and update operations.
		// If it is not provided, throw an error.
		return logical.ErrorResponse("missing aws_public_cert"), nil
	}

	// If explicitly set to empty string, error out
	if certEntry.AWSPublicCert == "" {
		return logical.ErrorResponse("invalid aws_public_cert"), nil
	}

	// Verify the certificate by decoding it and parsing it
	publicCert, err := decodePEMAndParseCertificate(certEntry.AWSPublicCert)
	if err != nil {
		return nil, err
	}
	if publicCert == nil {
		return logical.ErrorResponse("invalid certificate; failed to decode and parse certificate"), nil
	}

	// If none of the checks fail, save the provided certificate
	if err := b.nonLockedSetAWSPublicCertificateEntry(req.Storage, certName, certEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

// Struct awsPublicCert holds the AWS Public Key that is used to verify the PKCS#7 signature
// of the instnace identity document.
type awsPublicCert struct {
	AWSPublicCert string `json:"aws_public_cert" structs:"aws_public_cert" mapstructure:"aws_public_cert"`
	Type          string `json:"type" structs:"type" mapstructure:"type"`
}

const pathConfigCertificateSyn = `
Adds the AWS Public Key that is used to verify the PKCS#7 signature of the identidy document.
`

const pathConfigCertificateDesc = `
AWS Public Key which is used to verify the PKCS#7 signature of the identity document,
varies by region. The public key(s) can be found in AWS EC2 instance metadata documentation.
The default key that is used to verify the signature is the one that is applicable for
following regions: US East (N. Virginia), US West (Oregon), US West (N. California),
EU (Ireland), EU (Frankfurt), Asia Pacific (Tokyo), Asia Pacific (Seoul), Asia Pacific (Singapore),
Asia Pacific (Sydney), and South America (Sao Paulo).

If the instances belongs to region other than the above, the public key(s) for the
corresponding regions should be registered using this endpoint. PKCS#7 is verified
using a collection of certificates containing the default certificate and all the
certificates that are registered using this endpoint.
`
const pathListCertificatesHelpSyn = `
Lists all the AWS public certificates that are registered with the backend.
`
const pathListCertificatesHelpDesc = `
Certificates will be listed by their respective names that were used during registration.
`
