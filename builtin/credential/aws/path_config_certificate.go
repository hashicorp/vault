package aws

import (
	"crypto"
	"crypto/dsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// dsaSignature represents the contents of the signature of a signed
// content using digital signature algorithm.
type dsaSignature struct {
	R, S *big.Int
}

// As per AWS documentation, this public key is valid for US East (N. Virginia),
// US West (Oregon), US West (N. California), EU (Ireland), EU (Frankfurt),
// Asia Pacific (Tokyo), Asia Pacific (Seoul), Asia Pacific (Singapore),
// Asia Pacific (Sydney), and South America (Sao Paulo)
const defaultAWSPublicCert = `
-----BEGIN CERTIFICATE-----
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

func pathConfigCertificate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/certificate$",
		Fields: map[string]*framework.FieldSchema{
			"aws_public_cert": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     defaultAWSPublicCert,
				Description: "AWS Public key required to verify PKCS7 signature.",
			},
		},

		ExistenceCheck: b.pathConfigCertificateExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigCertificateCreateUpdate,
			logical.UpdateOperation: b.pathConfigCertificateCreateUpdate,
			logical.ReadOperation:   b.pathConfigCertificateRead,
		},

		HelpSynopsis:    pathConfigCertificateSyn,
		HelpDescription: pathConfigCertificateDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigCertificateExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := awsPublicCertificateEntry(req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// Decodes the PEM encoded certiticate and parses it into a x509 cert.
func decodePEMAndParseCertificate(certificate string) (*x509.Certificate, error) {
	// Decode the PEM block and error out if a block is not detected in the first attempt.
	decodedPublicCert, rest := pem.Decode([]byte(certificate))
	if len(rest) != 0 {
		return nil, fmt.Errorf("invalid certificate; failed to decode certificate")
	}

	// Check if the certificate can be parsed.
	publicCert, err := x509.ParseCertificate(decodedPublicCert.Bytes)
	if err != nil {
		return nil, err
	}
	if publicCert == nil {
		return nil, fmt.Errorf("invalid certificate; failed to parse certificate")
	}
	return publicCert, nil
}

// awsPublicCertificateParsed will fetch the storage entry for the certificate,
// decodes it and returns the parsed certificate.
func awsPublicCertificateParsed(s logical.Storage) (*x509.Certificate, error) {
	certEntry, err := awsPublicCertificateEntry(s)
	if err != nil {
		return nil, err
	}
	if certEntry == nil {
		return decodePEMAndParseCertificate(defaultAWSPublicCert)
	}
	return decodePEMAndParseCertificate(certEntry.AWSPublicCert)
}

// awsPublicCertificate is used to get the configured AWS Public Key that is used
// to verify the PKCS#7 signature of the instance identity document.
func awsPublicCertificateEntry(s logical.Storage) (*awsPublicCert, error) {
	entry, err := s.Get("config/certificate")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		// Existence check depends on this being nil when the storage entry is not present.
		return nil, nil
	}

	var result awsPublicCert
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// pathConfigCertificateRead is used to view the configured AWS Public Key that is
// used to verify the PKCS#7 signature of the instance identity document.
func (b *backend) pathConfigCertificateRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	certificateEntry, err := awsPublicCertificateEntry(req.Storage)
	if err != nil {
		return nil, err
	}
	if certificateEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"aws_public_cert": certificateEntry.AWSPublicCert,
		},
	}, nil
}

// pathConfigCertificateCreateUpdate is used to register an AWS Public Key that is
// used to verify the PKCS#7 signature of the instance identity document.
func (b *backend) pathConfigCertificateCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	// Check if there is already a certificate entry registered.
	certEntry, err := awsPublicCertificateEntry(req.Storage)
	if err != nil {
		return nil, err
	}
	if certEntry == nil {
		certEntry = &awsPublicCert{}
	}

	// Check if the value is provided by the client.
	certStrB64, ok := data.GetOk("aws_public_cert")
	if ok {
		certBytes, err := base64.StdEncoding.DecodeString(certStrB64.(string))
		if err != nil {
			return nil, err
		}

		certEntry.AWSPublicCert = string(certBytes)
	} else if req.Operation == logical.CreateOperation {
		certEntry.AWSPublicCert = data.Get("aws_public_cert").(string)
	}

	// If explicitly set to empty string, error out.
	if certEntry.AWSPublicCert == "" {
		return logical.ErrorResponse("missing aws_public_cert"), nil
	}

	// Verify the certificate by decoding it and parsing it.
	publicCert, err := decodePEMAndParseCertificate(certEntry.AWSPublicCert)
	if err != nil {
		return nil, err
	}
	if publicCert == nil {
		return logical.ErrorResponse("invalid certificate; failed to decode and parse certificate"), nil
	}

	// Before trusting the signature provided, validate its signature.

	// Extract the signature of the certificate.
	dsaSig := &dsaSignature{}
	dsaSigRest, err := asn1.Unmarshal(publicCert.Signature, dsaSig)
	if err != nil {
		return nil, err
	}
	if len(dsaSigRest) != 0 {
		return nil, fmt.Errorf("failed to unmarshal certificate's signature")
	}

	certHashFunc := crypto.SHA1.New()

	// RawTBSCertificate will contain the information in the certificate that is signed.
	certHashFunc.Write(publicCert.RawTBSCertificate)

	// Verify the signature using the public key present in the certificate.
	if !dsa.Verify(publicCert.PublicKey.(*dsa.PublicKey), certHashFunc.Sum(nil), dsaSig.R, dsaSig.S) {
		return logical.ErrorResponse("invalid certificate; failed to verify certificate's signature"), nil
	}

	// If none of the checks fail, save the provided certificate.
	entry, err := logical.StorageEntryJSON("config/certificate", certEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// Struct awsPublicCert holds the AWS Public Key that is used to verify the PKCS#7 signature
// of the instnace identity document.
type awsPublicCert struct {
	AWSPublicCert string `json:"aws_public_cert" structs:"aws_public_cert" mapstructure:"aws_public_cert"`
}

const pathConfigCertificateSyn = `
Configure the AWS Public Key that is used to verify the PKCS#7 signature of the identidy document.
`

const pathConfigCertificateDesc = `
AWS Public Key used to verify the PKCS#7 signature of the identity document
varies by region. It can be found in AWS's documentation. The default key that
is used to verify the signature is the one that is applicable for following regions:
US East (N. Virginia), US West (Oregon), US West (N. California), EU (Ireland),
EU (Frankfurt), Asia Pacific (Tokyo), Asia Pacific (Seoul), Asia Pacific (Singapore),
Asia Pacific (Sydney), and South America (Sao Paulo).
`
