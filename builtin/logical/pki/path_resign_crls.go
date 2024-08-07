// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/pki_backend"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	crlNumberParam          = "crl_number"
	deltaCrlBaseNumberParam = "delta_crl_base_number"
	nextUpdateParam         = "next_update"
	crlsParam               = "crls"
	formatParam             = "format"
)

var (
	akOid       = asn1.ObjectIdentifier{2, 5, 29, 35}
	crlNumOid   = asn1.ObjectIdentifier{2, 5, 29, 20}
	deltaCrlOid = asn1.ObjectIdentifier{2, 5, 29, 27}
)

func pathResignCrls(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issuer/" + framework.GenericNameRegex(issuerRefParam) + "/resign-crls",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKIIssuer,
			OperationVerb:   "resign",
			OperationSuffix: "crls",
		},

		Fields: map[string]*framework.FieldSchema{
			issuerRefParam: {
				Type: framework.TypeString,
				Description: `Reference to a existing issuer; either "default"
for the configured default issuer, an identifier or the name assigned
to the issuer.`,
				Default: defaultRef,
			},
			crlNumberParam: {
				Type:        framework.TypeInt,
				Description: `The sequence number to be written within the CRL Number extension.`,
			},
			deltaCrlBaseNumberParam: {
				Type: framework.TypeInt,
				Description: `Using a zero or greater value specifies the base CRL revision number to encode within
 a Delta CRL indicator extension, otherwise the extension will not be added.`,
				Default: -1,
			},
			nextUpdateParam: {
				Type: framework.TypeString,
				Description: `The amount of time the generated CRL should be
valid; defaults to 72 hours.`,
				Default: pki_backend.DefaultCrlConfig.Expiry,
			},
			crlsParam: {
				Type:        framework.TypeStringSlice,
				Description: `A list of PEM encoded CRLs to combine, originally signed by the requested issuer.`,
			},
			formatParam: {
				Type: framework.TypeString,
				Description: `The format of the combined CRL, can be "pem" or "der". If "der", the value will be
base64 encoded. Defaults to "pem".`,
				Default: "pem",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathUpdateResignCrlsHandler,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"crl": {
								Type:        framework.TypeString,
								Description: `CRL`,
								Required:    true,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis: `Combine and sign with the provided issuer different CRLs`,
		HelpDescription: `Provide two or more PEM encoded CRLs signed by the issuer,
 normally from separate Vault clusters to be combined and signed.`,
	}
}

func pathSignRevocationList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issuer/" + framework.GenericNameRegex(issuerRefParam) + "/sign-revocation-list",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKIIssuer,
			OperationVerb:   "sign",
			OperationSuffix: "revocation-list",
		},

		Fields: map[string]*framework.FieldSchema{
			issuerRefParam: {
				Type: framework.TypeString,
				Description: `Reference to a existing issuer; either "default"
for the configured default issuer, an identifier or the name assigned
to the issuer.`,
				Default: defaultRef,
			},
			crlNumberParam: {
				Type:        framework.TypeInt,
				Description: `The sequence number to be written within the CRL Number extension.`,
			},
			deltaCrlBaseNumberParam: {
				Type: framework.TypeInt,
				Description: `Using a zero or greater value specifies the base CRL revision number to encode within
 a Delta CRL indicator extension, otherwise the extension will not be added.`,
				Default: -1,
			},
			nextUpdateParam: {
				Type: framework.TypeString,
				Description: `The amount of time the generated CRL should be
valid; defaults to 72 hours.`,
				Default: pki_backend.DefaultCrlConfig.Expiry,
			},
			formatParam: {
				Type: framework.TypeString,
				Description: `The format of the combined CRL, can be "pem" or "der". If "der", the value will be
base64 encoded. Defaults to "pem".`,
				Default: "pem",
			},
			"revoked_certs": {
				Type: framework.TypeSlice,
				Description: `A list of maps containing the keys serial_number (string), revocation_time (string), 
and extensions (map with keys id (string), critical (bool), value (string))`,
			},
			"extensions": {
				Type: framework.TypeSlice,
				Description: `A list of maps containing extensions with keys id (string), critical (bool), 
value (string)`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathUpdateSignRevocationListHandler,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"crl": {
								Type:        framework.TypeString,
								Description: `CRL`,
								Required:    true,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis: `Generate and sign a CRL based on the provided parameters.`,
		HelpDescription: `Given a list of revoked certificates and other parameters, 
return a signed CRL based on the parameter values.`,
	}
}

func (b *backend) pathUpdateResignCrlsHandler(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if b.UseLegacyBundleCaStorage() {
		return logical.ErrorResponse("This API cannot be used until the migration has completed"), nil
	}

	issuerRef := GetIssuerRef(data)
	crlNumber := data.Get(crlNumberParam).(int)
	deltaCrlBaseNumber := data.Get(deltaCrlBaseNumberParam).(int)
	nextUpdateStr := data.Get(nextUpdateParam).(string)
	rawCrls := data.Get(crlsParam).([]string)

	format, err := parseCrlFormat(data.Get(formatParam).(string))
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	nextUpdateOffset, err := parseutil.ParseDurationSecond(nextUpdateStr)
	if err != nil {
		return logical.ErrorResponse("invalid value for %s: %v", nextUpdateParam, err), nil
	}

	if nextUpdateOffset <= 0 {
		return logical.ErrorResponse("%s parameter must be greater than 0", nextUpdateParam), nil
	}

	if crlNumber < 0 {
		return logical.ErrorResponse("%s parameter must be 0 or greater", crlNumberParam), nil
	}
	if deltaCrlBaseNumber < -1 {
		return logical.ErrorResponse("%s parameter must be -1 or greater", deltaCrlBaseNumberParam), nil
	}

	if issuerRef == "" {
		return logical.ErrorResponse("%s parameter cannot be blank", issuerRefParam), nil
	}

	providedCrls, err := decodePemCrls(rawCrls)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	sc := b.makeStorageContext(ctx, request.Storage)
	caBundle, err := getCaBundle(sc, issuerRef)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if err := verifyCrlsAreFromIssuersKey(caBundle.Certificate, providedCrls); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	revokedCerts, warnings, err := getAllRevokedCertsFromPem(providedCrls)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	now := time.Now()
	template := &x509.RevocationList{
		SignatureAlgorithm:  caBundle.RevocationSigAlg,
		RevokedCertificates: revokedCerts,
		Number:              big.NewInt(int64(crlNumber)),
		ThisUpdate:          now,
		NextUpdate:          now.Add(nextUpdateOffset),
	}

	if deltaCrlBaseNumber > -1 {
		ext, err := certutil.CreateDeltaCRLIndicatorExt(int64(deltaCrlBaseNumber))
		if err != nil {
			return nil, fmt.Errorf("could not create crl delta indicator extension: %w", err)
		}
		template.ExtraExtensions = []pkix.Extension{ext}
	}

	crlBytes, err := x509.CreateRevocationList(rand.Reader, template, caBundle.Certificate, caBundle.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating new CRL: %w", err)
	}

	body := encodeResponse(crlBytes, format == "der")

	return &logical.Response{
		Warnings: warnings,
		Data: map[string]interface{}{
			"crl": body,
		},
	}, nil
}

func (b *backend) pathUpdateSignRevocationListHandler(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if b.UseLegacyBundleCaStorage() {
		return logical.ErrorResponse("This API cannot be used until the migration has completed"), nil
	}

	issuerRef := GetIssuerRef(data)
	crlNumber := data.Get(crlNumberParam).(int)
	deltaCrlBaseNumber := data.Get(deltaCrlBaseNumberParam).(int)
	nextUpdateStr := data.Get(nextUpdateParam).(string)
	nextUpdateOffset, err := parseutil.ParseDurationSecond(nextUpdateStr)
	if err != nil {
		return logical.ErrorResponse("invalid value for %s: %v", nextUpdateParam, err), nil
	}

	if nextUpdateOffset <= 0 {
		return logical.ErrorResponse("%s parameter must be greater than 0", nextUpdateParam), nil
	}

	if crlNumber < 0 {
		return logical.ErrorResponse("%s parameter must be 0 or greater", crlNumberParam), nil
	}
	if deltaCrlBaseNumber < -1 {
		return logical.ErrorResponse("%s parameter must be -1 or greater", deltaCrlBaseNumberParam), nil
	}

	if issuerRef == "" {
		return logical.ErrorResponse("%s parameter cannot be blank", issuerRefParam), nil
	}

	format, err := parseCrlFormat(data.Get(formatParam).(string))
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	revokedCerts, err := parseRevokedCertsParam(data.Get("revoked_certs").([]interface{}))
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	crlExtensions, err := parseExtensionsParam(data.Get("extensions").([]interface{}))
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	sc := b.makeStorageContext(ctx, request.Storage)
	caBundle, err := getCaBundle(sc, issuerRef)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if deltaCrlBaseNumber > -1 {
		ext, err := certutil.CreateDeltaCRLIndicatorExt(int64(deltaCrlBaseNumber))
		if err != nil {
			return nil, fmt.Errorf("could not create crl delta indicator extension: %w", err)
		}
		crlExtensions = append(crlExtensions, ext)
	}

	now := time.Now()
	template := &x509.RevocationList{
		SignatureAlgorithm:  caBundle.RevocationSigAlg,
		RevokedCertificates: revokedCerts,
		Number:              big.NewInt(int64(crlNumber)),
		ThisUpdate:          now,
		NextUpdate:          now.Add(nextUpdateOffset),
		ExtraExtensions:     crlExtensions,
	}

	crlBytes, err := x509.CreateRevocationList(rand.Reader, template, caBundle.Certificate, caBundle.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating new CRL: %w", err)
	}

	body := encodeResponse(crlBytes, format == "der")

	return &logical.Response{
		Data: map[string]interface{}{
			"crl": body,
		},
	}, nil
}

func parseRevokedCertsParam(revokedCerts []interface{}) ([]pkix.RevokedCertificate, error) {
	var parsedCerts []pkix.RevokedCertificate
	seenSerials := make(map[*big.Int]int)
	for i, entry := range revokedCerts {
		if revokedCert, ok := entry.(map[string]interface{}); ok {
			serialNum, err := parseSerialNum(revokedCert)
			if err != nil {
				return nil, fmt.Errorf("failed parsing serial_number from entry %d: %w", i, err)
			}

			if origEntry, exists := seenSerials[serialNum]; exists {
				serialNumStr := revokedCert["serial_number"]
				return nil, fmt.Errorf("duplicate serial number: %s, original entry %d and %d", serialNumStr, origEntry, i)
			}

			seenSerials[serialNum] = i

			revocationTime, err := parseRevocationTime(revokedCert)
			if err != nil {
				return nil, fmt.Errorf("failed parsing revocation_time from entry %d: %w", i, err)
			}

			extensions, err := parseCertExtensions(revokedCert)
			if err != nil {
				return nil, fmt.Errorf("failed parsing extensions from entry %d: %w", i, err)
			}

			parsedCerts = append(parsedCerts, pkix.RevokedCertificate{
				SerialNumber:   serialNum,
				RevocationTime: revocationTime,
				Extensions:     extensions,
			})
		}
	}

	return parsedCerts, nil
}

func parseCertExtensions(cert map[string]interface{}) ([]pkix.Extension, error) {
	extRaw, exists := cert["extensions"]
	if !exists || extRaw == nil || extRaw == "" {
		// We don't require extensions to be populated
		return []pkix.Extension{}, nil
	}

	extListRaw, ok := extRaw.([]interface{})
	if !ok {
		return nil, errors.New("'extensions' field did not contain a slice")
	}

	return parseExtensionsParam(extListRaw)
}

func parseExtensionsParam(extRawList []interface{}) ([]pkix.Extension, error) {
	var extensions []pkix.Extension
	seenOid := make(map[string]struct{})
	for i, entryRaw := range extRawList {
		entry, ok := entryRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("extension entry %d not a map", i)
		}
		extension, err := parseExtension(entry)
		if err != nil {
			return nil, fmt.Errorf("failed parsing extension entry %d: %w", i, err)
		}

		parsedIdStr := extension.Id.String()
		if _, exists := seenOid[parsedIdStr]; exists {
			return nil, fmt.Errorf("duplicate extension id: %s", parsedIdStr)
		}

		seenOid[parsedIdStr] = struct{}{}

		extensions = append(extensions, extension)
	}

	return extensions, nil
}

func parseExtension(entry map[string]interface{}) (pkix.Extension, error) {
	asnObjectId, err := parseExtAsn1ObjectId(entry)
	if err != nil {
		return pkix.Extension{}, err
	}

	if asnObjectId.Equal(akOid) {
		return pkix.Extension{}, fmt.Errorf("authority key object identifier (%s) is reserved", akOid.String())
	}

	if asnObjectId.Equal(crlNumOid) {
		return pkix.Extension{}, fmt.Errorf("crl number object identifier (%s) is reserved", crlNumOid.String())
	}

	if asnObjectId.Equal(deltaCrlOid) {
		return pkix.Extension{}, fmt.Errorf("delta crl object identifier (%s) is reserved", deltaCrlOid.String())
	}

	critical, err := parseExtCritical(entry)
	if err != nil {
		return pkix.Extension{}, err
	}

	extVal, err := parseExtValue(entry)
	if err != nil {
		return pkix.Extension{}, err
	}

	return pkix.Extension{
		Id:       asnObjectId,
		Critical: critical,
		Value:    extVal,
	}, nil
}

func parseExtValue(entry map[string]interface{}) ([]byte, error) {
	valRaw, exists := entry["value"]
	if !exists {
		return nil, errors.New("missing 'value' field")
	}

	valStr, err := parseutil.ParseString(valRaw)
	if err != nil {
		return nil, fmt.Errorf("'value' field value was not a string: %w", err)
	}

	if len(valStr) == 0 {
		return []byte{}, nil
	}

	decodeString, err := base64.StdEncoding.DecodeString(valStr)
	if err != nil {
		return nil, fmt.Errorf("failed base64 decoding 'value' field: %w", err)
	}
	return decodeString, nil
}

func parseExtCritical(entry map[string]interface{}) (bool, error) {
	critRaw, exists := entry["critical"]
	if !exists || critRaw == nil || critRaw == "" {
		// Optional field, so just return as if they provided the value false.
		return false, nil
	}

	myBool, err := parseutil.ParseBool(critRaw)
	if err != nil {
		return false, fmt.Errorf("critical field value failed to be parsed: %w", err)
	}

	return myBool, nil
}

func parseExtAsn1ObjectId(entry map[string]interface{}) (asn1.ObjectIdentifier, error) {
	idRaw, idExists := entry["id"]
	if !idExists {
		return asn1.ObjectIdentifier{}, errors.New("missing id field")
	}

	oidStr, err := parseutil.ParseString(idRaw)
	if err != nil {
		return nil, fmt.Errorf("'id' field value was not a string: %w", err)
	}

	if len(oidStr) == 0 {
		return asn1.ObjectIdentifier{}, errors.New("zero length object identifier")
	}

	// Parse out dot notation
	oidParts := strings.Split(oidStr, ".")
	oid := make(asn1.ObjectIdentifier, len(oidParts), len(oidParts))
	for i := range oidParts {
		oidIntVal, err := strconv.Atoi(oidParts[i])
		if err != nil {
			return nil, fmt.Errorf("failed parsing asn1 index element %d value %s: %w", i, oidParts[i], err)
		}
		oid[i] = oidIntVal
	}
	return oid, nil
}

func parseRevocationTime(cert map[string]interface{}) (time.Time, error) {
	var revTime time.Time
	revTimeRaw, exists := cert["revocation_time"]
	if !exists {
		return revTime, errors.New("missing 'revocation_time' field")
	}
	revTime, err := parseutil.ParseAbsoluteTime(revTimeRaw)
	if err != nil {
		return revTime, fmt.Errorf("failed parsing time %v: %w", revTimeRaw, err)
	}
	return revTime, nil
}

func parseSerialNum(cert map[string]interface{}) (*big.Int, error) {
	serialNumRaw, serialExists := cert["serial_number"]
	if !serialExists {
		return nil, errors.New("missing 'serial_number' field")
	}
	return parseSerialNumStr(serialNumRaw)
}

func parseSerialNumStr(serialNumRaw interface{}) (*big.Int, error) {
	serialNumStr, err := parseutil.ParseString(serialNumRaw)
	if err != nil {
		return nil, fmt.Errorf("'serial_number' field value was not a string: %w", err)
	}
	// Clean up any provided serials to decoder
	for _, separator := range []string{":", ".", "-", " "} {
		serialNumStr = strings.ReplaceAll(serialNumStr, separator, "")
	}
	// Prefer hex.DecodeString over certutil.ParseHexFormatted as we don't need a separator
	serialBytes, err := hex.DecodeString(serialNumStr)
	if err != nil {
		return nil, fmt.Errorf("'serial_number' failed converting to bytes: %w", err)
	}

	bigIntSerial := big.Int{}
	bigIntSerial.SetBytes(serialBytes)
	return &bigIntSerial, nil
}

func parseCrlFormat(requestedValue string) (string, error) {
	format := strings.ToLower(requestedValue)
	switch format {
	case "pem", "der":
		return format, nil
	default:
		return "", fmt.Errorf("unknown format value of %s", requestedValue)
	}
}

func verifyCrlsAreFromIssuersKey(caCert *x509.Certificate, crls []*x509.RevocationList) error {
	for i, crl := range crls {
		// At this point we assume if the issuer's key signed the CRL that is a good enough check
		// to validate that we owned/generated the provided CRL.
		if err := crl.CheckSignatureFrom(caCert); err != nil {
			return fmt.Errorf("CRL index: %d was not signed by requested issuer", i)
		}
	}

	return nil
}

func encodeResponse(crlBytes []byte, derFormatRequested bool) string {
	if derFormatRequested {
		return base64.StdEncoding.EncodeToString(crlBytes)
	}

	block := pem.Block{
		Type:  "X509 CRL",
		Bytes: crlBytes,
	}
	return string(pem.EncodeToMemory(&block))
}

func getAllRevokedCertsFromPem(crls []*x509.RevocationList) ([]pkix.RevokedCertificate, []string, error) {
	uniqueCert := map[string]pkix.RevokedCertificate{}
	var warnings []string
	for _, crl := range crls {
		for _, curCert := range crl.RevokedCertificates {
			serial := serialFromBigInt(curCert.SerialNumber)
			// Get rid of any extensions the existing certificate might have had.
			curCert.Extensions = []pkix.Extension{}

			existingCert, exists := uniqueCert[serial]
			if !exists {
				// First time we see the revoked cert
				uniqueCert[serial] = curCert
				continue
			}

			if existingCert.RevocationTime.Equal(curCert.RevocationTime) {
				// Same revocation times, just skip it
				continue
			}

			warn := fmt.Sprintf("Duplicate serial %s with different revocation "+
				"times detected, using oldest revocation time", serial)
			warnings = append(warnings, warn)

			if existingCert.RevocationTime.After(curCert.RevocationTime) {
				uniqueCert[serial] = curCert
			}
		}
	}

	var revokedCerts []pkix.RevokedCertificate
	for _, cert := range uniqueCert {
		revokedCerts = append(revokedCerts, cert)
	}

	return revokedCerts, warnings, nil
}

func getCaBundle(sc *storageContext, issuerRef string) (*certutil.CAInfoBundle, error) {
	issuerId, err := sc.resolveIssuerReference(issuerRef)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve issuer %s: %w", issuerRefParam, err)
	}

	return sc.fetchCAInfoByIssuerId(issuerId, issuing.CRLSigningUsage)
}

func decodePemCrls(rawCrls []string) ([]*x509.RevocationList, error) {
	var crls []*x509.RevocationList
	for i, rawCrl := range rawCrls {
		crl, err := decodePemCrl(rawCrl)
		if err != nil {
			return nil, fmt.Errorf("failed decoding crl %d: %w", i, err)
		}
		crls = append(crls, crl)
	}

	return crls, nil
}

func decodePemCrl(crl string) (*x509.RevocationList, error) {
	block, rest := pem.Decode([]byte(crl))
	if len(rest) != 0 {
		return nil, errors.New("invalid crl; should be one PEM block only")
	}

	return x509.ParseRevocationList(block.Bytes)
}
