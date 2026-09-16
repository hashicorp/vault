// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"time"

	"github.com/cloudflare/circl/sign/mldsa/mldsa65"
	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
	"github.com/hashicorp/vault/sdk/helper/errutil"
)

// OIDs for ML-DSA algorithms as defined in NIST FIPS 204.
// Per FIPS 204, the key OID and signature OID are identical for ML-DSA:
// the same OID identifies both the public key algorithm in
// SubjectPublicKeyInfo and the signature algorithm in
// signatureAlgorithm / TBS signatureAlgorithm fields.
var (
	oidMLDSA65 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 3, 18}
	oidMLDSA87 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 3, 19}
)

// x509Certificate is the ASN.1 structure of an X.509 certificate.
type x509Certificate struct {
	TBSCertificate     asn1.RawValue
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

// isMLDSAKey returns true if the signer uses an ML-DSA key type.
func isMLDSAKey(signer crypto.Signer) bool {
	keyType := GetPrivateKeyTypeFromSigner(signer)
	return keyType == MLDSA65PrivateKey || keyType == MLDSA87PrivateKey
}

// isMLDSAKeyType returns true if the key type string is an ML-DSA type.
func isMLDSAKeyType(keyType string) bool {
	return keyType == "ml-dsa-65" || keyType == "ml-dsa-87"
}

// mldsaSignatureAlgorithmOID returns the OID for the given ML-DSA key type.
func mldsaSignatureAlgorithmOID(signer crypto.Signer) (asn1.ObjectIdentifier, error) {
	keyType := GetPrivateKeyTypeFromSigner(signer)
	switch keyType {
	case MLDSA65PrivateKey:
		return oidMLDSA65, nil
	case MLDSA87PrivateKey:
		return oidMLDSA87, nil
	default:
		return nil, fmt.Errorf("not an ML-DSA key type: %s", keyType)
	}
}

// createMLDSACertificate creates an X.509 certificate signed with an ML-DSA key.
// Go's standard x509.CreateCertificate does not yet support ML-DSA, so we
// construct the certificate manually using ASN.1.
//
// The approach:
// 1. Build the TBS certificate structure directly using ASN.1,
//    embedding the correct ML-DSA signature algorithm OID.
// 2. Marshal the TBS to DER.
// 3. Sign the TBS DER bytes with the ML-DSA signer.
// 4. Assemble the final certificate DER with signature.
// 5. Verify the signature to ensure correctness.
func createMLDSACertificate(template, parent *x509.Certificate, pub crypto.PublicKey, signer crypto.Signer) ([]byte, error) {
	sigAlgOID, err := mldsaSignatureAlgorithmOID(signer)
	if err != nil {
		return nil, err
	}

	// Build the TBS certificate from scratch because Go's
	// x509.CreateCertificate does not support ML-DSA public keys.
	tbsCert, err := buildTBSCertificate(template, parent, pub, sigAlgOID)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error building TBS certificate: %v", err)}
	}

	// Step 2: Marshal TBS to DER
	tbsDER, err := asn1.Marshal(tbsCert)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error marshaling TBS certificate: %v", err)}
	}

	// Step 3: Sign the TBS DER bytes
	sig, err := signer.Sign(nil, tbsDER, crypto.Hash(0))
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error signing certificate with ML-DSA: %v", err)}
	}

	// Step 4: Assemble the final certificate
	sigAlg := pkix.AlgorithmIdentifier{Algorithm: sigAlgOID}
	cert := x509Certificate{
		TBSCertificate:     asn1.RawValue{FullBytes: tbsDER},
		SignatureAlgorithm: sigAlg,
		SignatureValue:     asn1.BitString{Bytes: sig, BitLength: len(sig) * 8},
	}

	certDER, err := asn1.Marshal(cert)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error marshaling certificate: %v", err)}
	}

	// Step 5: Verify the signature to ensure correctness
	if err := verifyMLDSASignature(signer.Public(), tbsDER, sig); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("post-creation signature verification failed: %v", err)}
	}

	return certDER, nil
}

// tbsCertificate is a simplified ASN.1 TBS certificate structure.
type tbsCertificate struct {
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           validity
	Subject            asn1.RawValue
	PublicKeyInfo      asn1.RawValue
	Extensions         []pkix.Extension `asn1:"optional,explicit,tag:3"`
}

type validity struct {
	NotBefore, NotAfter asn1.RawValue
}

// marshalPublicKeyInfo marshals the ML-DSA public key into a
// SubjectPublicKeyInfo structure.
func marshalPublicKeyInfo(pub crypto.PublicKey, algOID asn1.ObjectIdentifier) (asn1.RawValue, error) {
	type publicKeyInfo struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}

	type binaryMarshaler interface {
		MarshalBinary() ([]byte, error)
	}

	bm, ok := pub.(binaryMarshaler)
	if !ok {
		return asn1.RawValue{}, fmt.Errorf("public key does not implement MarshalBinary")
	}

	pubBytes, err := bm.MarshalBinary()
	if err != nil {
		return asn1.RawValue{}, fmt.Errorf("error marshaling public key: %v", err)
	}

	pki := publicKeyInfo{
		Algorithm: pkix.AlgorithmIdentifier{Algorithm: algOID},
		PublicKey: asn1.BitString{Bytes: pubBytes, BitLength: len(pubBytes) * 8},
	}

	pkiDER, err := asn1.Marshal(pki)
	if err != nil {
		return asn1.RawValue{}, fmt.Errorf("error marshaling SubjectPublicKeyInfo: %v", err)
	}

	return asn1.RawValue{FullBytes: pkiDER}, nil
}

// buildTBSCertificate constructs a TBS certificate structure from the template.
func buildTBSCertificate(template, parent *x509.Certificate, pub crypto.PublicKey, sigAlgOID asn1.ObjectIdentifier) (tbsCertificate, error) {
	// Marshal subject and issuer names
	issuerRDN, err := asn1.Marshal(parent.Subject.ToRDNSequence())
	if err != nil {
		return tbsCertificate{}, fmt.Errorf("error marshaling issuer: %v", err)
	}

	subjectRDN, err := asn1.Marshal(template.Subject.ToRDNSequence())
	if err != nil {
		return tbsCertificate{}, fmt.Errorf("error marshaling subject: %v", err)
	}

	// Marshal validity times
	notBeforeDER, err := asn1.Marshal(template.NotBefore)
	if err != nil {
		return tbsCertificate{}, fmt.Errorf("error marshaling NotBefore: %v", err)
	}
	notAfterDER, err := asn1.Marshal(template.NotAfter)
	if err != nil {
		return tbsCertificate{}, fmt.Errorf("error marshaling NotAfter: %v", err)
	}

	// Marshal public key info
	pkInfo, err := marshalPublicKeyInfo(pub, sigAlgOID)
	if err != nil {
		return tbsCertificate{}, err
	}

	// Build extensions from the template
	extensions, err := buildExtensions(template)
	if err != nil {
		return tbsCertificate{}, err
	}

	tbs := tbsCertificate{
		Version:            2, // v3
		SerialNumber:       template.SerialNumber,
		SignatureAlgorithm: pkix.AlgorithmIdentifier{Algorithm: sigAlgOID},
		Issuer:             asn1.RawValue{FullBytes: issuerRDN},
		Validity: validity{
			NotBefore: asn1.RawValue{FullBytes: notBeforeDER},
			NotAfter:  asn1.RawValue{FullBytes: notAfterDER},
		},
		Subject:       asn1.RawValue{FullBytes: subjectRDN},
		PublicKeyInfo: pkInfo,
		Extensions:    extensions,
	}

	return tbs, nil
}

// buildExtensions creates the certificate extensions from the template.
func buildExtensions(template *x509.Certificate) ([]pkix.Extension, error) {
	var extensions []pkix.Extension

	// Subject Key Identifier
	if len(template.SubjectKeyId) > 0 {
		skidDER, err := asn1.Marshal(template.SubjectKeyId)
		if err != nil {
			return nil, fmt.Errorf("error marshaling SKID: %v", err)
		}
		extensions = append(extensions, pkix.Extension{
			Id:    asn1.ObjectIdentifier{2, 5, 29, 14},
			Value: skidDER,
		})
	}

	// Authority Key Identifier
	if len(template.AuthorityKeyId) > 0 {
		type authKeyId struct {
			KeyIdentifier []byte `asn1:"optional,tag:0"`
		}
		akidDER, err := asn1.Marshal(authKeyId{KeyIdentifier: template.AuthorityKeyId})
		if err != nil {
			return nil, fmt.Errorf("error marshaling AKID: %v", err)
		}
		extensions = append(extensions, pkix.Extension{
			Id:    asn1.ObjectIdentifier{2, 5, 29, 35},
			Value: akidDER,
		})
	}

	// Basic Constraints
	if template.BasicConstraintsValid || template.IsCA {
		type basicConstraints struct {
			IsCA       bool `asn1:"optional"`
			MaxPathLen int  `asn1:"optional,default:-1"`
		}
		bc := basicConstraints{IsCA: template.IsCA}
		if template.MaxPathLen > 0 || template.MaxPathLenZero {
			bc.MaxPathLen = template.MaxPathLen
		} else {
			bc.MaxPathLen = -1
		}
		bcDER, err := asn1.Marshal(bc)
		if err != nil {
			return nil, fmt.Errorf("error marshaling basic constraints: %v", err)
		}
		extensions = append(extensions, pkix.Extension{
			Id:       asn1.ObjectIdentifier{2, 5, 29, 19},
			Critical: true,
			Value:    bcDER,
		})
	}

	// Key Usage
	if template.KeyUsage != 0 {
		ku, err := marshalKeyUsage(template.KeyUsage)
		if err != nil {
			return nil, fmt.Errorf("error marshaling key usage: %v", err)
		}
		extensions = append(extensions, ku)
	}

	// Extended Key Usage
	if len(template.ExtKeyUsage) > 0 || len(template.UnknownExtKeyUsage) > 0 {
		var oids []asn1.ObjectIdentifier
		for _, eku := range template.ExtKeyUsage {
			oid, ok := ekuToOID(eku)
			if ok {
				oids = append(oids, oid)
			}
		}
		oids = append(oids, template.UnknownExtKeyUsage...)

		ekuDER, err := asn1.Marshal(oids)
		if err != nil {
			return nil, fmt.Errorf("error marshaling EKU: %v", err)
		}
		extensions = append(extensions, pkix.Extension{
			Id:    asn1.ObjectIdentifier{2, 5, 29, 37},
			Value: ekuDER,
		})
	}

	// SAN (Subject Alternative Name)
	sanDER, err := marshalSANExtension(template)
	if err != nil {
		return nil, err
	}
	if sanDER != nil {
		extensions = append(extensions, pkix.Extension{
			Id:    asn1.ObjectIdentifier{2, 5, 29, 17},
			Value: sanDER,
		})
	}

	// Include any extra extensions from the template
	extensions = append(extensions, template.ExtraExtensions...)

	return extensions, nil
}

// ekuToOID maps an x509.ExtKeyUsage to its OID.
func ekuToOID(eku x509.ExtKeyUsage) (asn1.ObjectIdentifier, bool) {
	switch eku {
	case x509.ExtKeyUsageAny:
		return asn1.ObjectIdentifier{2, 5, 29, 37, 0}, true
	case x509.ExtKeyUsageServerAuth:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 1}, true
	case x509.ExtKeyUsageClientAuth:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 2}, true
	case x509.ExtKeyUsageCodeSigning:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 3}, true
	case x509.ExtKeyUsageEmailProtection:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 4}, true
	case x509.ExtKeyUsageIPSECEndSystem:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 5}, true
	case x509.ExtKeyUsageIPSECTunnel:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 6}, true
	case x509.ExtKeyUsageIPSECUser:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 7}, true
	case x509.ExtKeyUsageTimeStamping:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 8}, true
	case x509.ExtKeyUsageOCSPSigning:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 9}, true
	case x509.ExtKeyUsageMicrosoftServerGatedCrypto:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 10, 3, 3}, true
	case x509.ExtKeyUsageNetscapeServerGatedCrypto:
		return asn1.ObjectIdentifier{2, 16, 840, 1, 113730, 4, 1}, true
	case x509.ExtKeyUsageMicrosoftCommercialCodeSigning:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 22}, true
	case x509.ExtKeyUsageMicrosoftKernelCodeSigning:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 61, 1, 1}, true
	default:
		return nil, false
	}
}

// marshalSANExtension builds the SAN extension DER bytes from the template.
func marshalSANExtension(template *x509.Certificate) ([]byte, error) {
	if len(template.DNSNames) == 0 && len(template.EmailAddresses) == 0 &&
		len(template.IPAddresses) == 0 && len(template.URIs) == 0 {
		return nil, nil
	}

	var rawValues []asn1.RawValue

	for _, name := range template.DNSNames {
		rawValues = append(rawValues, asn1.RawValue{
			Tag:   2, // dNSName
			Class: asn1.ClassContextSpecific,
			Bytes: []byte(name),
		})
	}

	for _, email := range template.EmailAddresses {
		rawValues = append(rawValues, asn1.RawValue{
			Tag:   1, // rfc822Name
			Class: asn1.ClassContextSpecific,
			Bytes: []byte(email),
		})
	}

	for _, ip := range template.IPAddresses {
		rawIP := ip.To4()
		if rawIP == nil {
			rawIP = ip.To16()
		}
		rawValues = append(rawValues, asn1.RawValue{
			Tag:   7, // iPAddress
			Class: asn1.ClassContextSpecific,
			Bytes: rawIP,
		})
	}

	for _, uri := range template.URIs {
		rawValues = append(rawValues, asn1.RawValue{
			Tag:   6, // uniformResourceIdentifier
			Class: asn1.ClassContextSpecific,
			Bytes: []byte(uri.String()),
		})
	}

	sanDER, err := asn1.Marshal(rawValues)
	if err != nil {
		return nil, fmt.Errorf("error marshaling SAN extension: %v", err)
	}

	return sanDER, nil
}

// createMLDSACSR creates a certificate signing request signed with an ML-DSA key.
func createMLDSACSR(template *x509.CertificateRequest, signer crypto.Signer) ([]byte, error) {
	sigAlgOID, err := mldsaSignatureAlgorithmOID(signer)
	if err != nil {
		return nil, err
	}

	// Marshal subject
	subjectRDN, err := asn1.Marshal(template.Subject.ToRDNSequence())
	if err != nil {
		return nil, fmt.Errorf("error marshaling CSR subject: %v", err)
	}

	// Marshal public key info
	pkInfo, err := marshalPublicKeyInfo(signer.Public(), sigAlgOID)
	if err != nil {
		return nil, err
	}

	// Build CSR attributes (extensions)
	var attributes []asn1.RawValue
	var extensions []pkix.Extension

	// Add SANs
	sanDER, err := marshalCSRSANExtension(template)
	if err != nil {
		return nil, err
	}
	if sanDER != nil {
		extensions = append(extensions, pkix.Extension{
			Id:    asn1.ObjectIdentifier{2, 5, 29, 17},
			Value: sanDER,
		})
	}

	// Add extra extensions from template
	extensions = append(extensions, template.ExtraExtensions...)

	if len(extensions) > 0 {
		extDER, err := asn1.Marshal(extensions)
		if err != nil {
			return nil, fmt.Errorf("error marshaling CSR extensions: %v", err)
		}
		// extensionRequest attribute OID: 1.2.840.113549.1.9.14
		attrType, err := asn1.Marshal(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 14})
		if err != nil {
			return nil, fmt.Errorf("error marshaling attr type: %v", err)
		}
		// Wrap in SET
		attrValues, err := asn1.Marshal(asn1.RawValue{
			Class:      asn1.ClassUniversal,
			Tag:        asn1.TagSet,
			IsCompound: true,
			Bytes:      extDER,
		})
		if err != nil {
			return nil, fmt.Errorf("error marshaling attr values: %v", err)
		}

		attrBytes := append(attrType, attrValues...)
		attributes = append(attributes, asn1.RawValue{
			Class:      asn1.ClassUniversal,
			Tag:        asn1.TagSequence,
			IsCompound: true,
			Bytes:      attrBytes,
		})
	}

	// Build CertificationRequestInfo
	type certificationRequestInfo struct {
		Version       int
		Subject       asn1.RawValue
		PublicKeyInfo asn1.RawValue
		Attributes    asn1.RawValue `asn1:"tag:0"`
	}

	var attrDER []byte
	if len(attributes) > 0 {
		for _, attr := range attributes {
			b, err := asn1.Marshal(attr)
			if err != nil {
				return nil, fmt.Errorf("error marshaling attribute: %v", err)
			}
			attrDER = append(attrDER, b...)
		}
	}

	cri := certificationRequestInfo{
		Version:       0,
		Subject:       asn1.RawValue{FullBytes: subjectRDN},
		PublicKeyInfo: pkInfo,
		Attributes:    asn1.RawValue{Class: asn1.ClassContextSpecific, Tag: 0, IsCompound: true, Bytes: attrDER},
	}

	criDER, err := asn1.Marshal(cri)
	if err != nil {
		return nil, fmt.Errorf("error marshaling CertificationRequestInfo: %v", err)
	}

	// Sign
	sig, err := signer.Sign(nil, criDER, crypto.Hash(0))
	if err != nil {
		return nil, fmt.Errorf("error signing CSR with ML-DSA: %v", err)
	}

	// Assemble CSR
	type certificationRequest struct {
		CertificationRequestInfo asn1.RawValue
		SignatureAlgorithm       pkix.AlgorithmIdentifier
		Signature                asn1.BitString
	}

	csr := certificationRequest{
		CertificationRequestInfo: asn1.RawValue{FullBytes: criDER},
		SignatureAlgorithm:       pkix.AlgorithmIdentifier{Algorithm: sigAlgOID},
		Signature:                asn1.BitString{Bytes: sig, BitLength: len(sig) * 8},
	}

	// Verify the signature to ensure correctness
	if err := verifyMLDSASignature(signer.Public(), criDER, sig); err != nil {
		return nil, fmt.Errorf("post-creation CSR signature verification failed: %v", err)
	}

	return asn1.Marshal(csr)
}

// verifyMLDSASignature verifies an ML-DSA signature using the appropriate
// circl verification function based on the public key type.
func verifyMLDSASignature(pub crypto.PublicKey, msg, sig []byte) error {
	switch pk := pub.(type) {
	case *mldsa65.PublicKey:
		if !mldsa65.Verify(pk, msg, nil, sig) {
			return fmt.Errorf("ML-DSA-65 signature verification failed")
		}
	case *mldsa87.PublicKey:
		if !mldsa87.Verify(pk, msg, nil, sig) {
			return fmt.Errorf("ML-DSA-87 signature verification failed")
		}
	default:
		return fmt.Errorf("unsupported public key type for ML-DSA verification: %T", pub)
	}
	return nil
}

// parseMLDSACertificate parses an ML-DSA certificate from DER bytes.
// Go's x509.ParseCertificate does not recognize ML-DSA signature algorithms,
// so we parse the ASN.1 structure manually and populate an x509.Certificate
// with the fields we can extract.
func parseMLDSACertificate(der []byte) (*x509.Certificate, error) {
	var cert struct {
		TBS struct {
			Version      int `asn1:"optional,explicit,default:0,tag:0"`
			SerialNumber *big.Int
			SigAlg       pkix.AlgorithmIdentifier
			Issuer       asn1.RawValue
			Validity     struct {
				NotBefore time.Time
				NotAfter  time.Time
			}
			Subject    asn1.RawValue
			PublicKey  asn1.RawValue
			Extensions []pkix.Extension `asn1:"optional,explicit,tag:3"`
		} `asn1:"sequence"`
		SigAlg    pkix.AlgorithmIdentifier
		Signature asn1.BitString
	}

	rest, err := asn1.Unmarshal(der, &cert)
	if err != nil {
		return nil, fmt.Errorf("error parsing ML-DSA certificate ASN.1: %v", err)
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("trailing data after certificate")
	}

	// Parse Subject and Issuer RDN sequences
	var issuerRDN pkix.RDNSequence
	if _, err := asn1.Unmarshal(cert.TBS.Issuer.FullBytes, &issuerRDN); err != nil {
		return nil, fmt.Errorf("error parsing issuer: %v", err)
	}
	var subjectRDN pkix.RDNSequence
	if _, err := asn1.Unmarshal(cert.TBS.Subject.FullBytes, &subjectRDN); err != nil {
		return nil, fmt.Errorf("error parsing subject: %v", err)
	}

	var issuerName pkix.Name
	issuerName.FillFromRDNSequence(&issuerRDN)
	var subjectName pkix.Name
	subjectName.FillFromRDNSequence(&subjectRDN)

	result := &x509.Certificate{
		Raw:                der,
		SerialNumber:       cert.TBS.SerialNumber,
		Issuer:             issuerName,
		Subject:            subjectName,
		NotBefore:          cert.TBS.Validity.NotBefore,
		NotAfter:           cert.TBS.Validity.NotAfter,
		RawIssuer:          cert.TBS.Issuer.FullBytes,
		RawSubject:         cert.TBS.Subject.FullBytes,
		Signature:          cert.Signature.Bytes,
		SignatureAlgorithm: x509.UnknownSignatureAlgorithm,
	}

	// Parse extensions
	for _, ext := range cert.TBS.Extensions {
		result.Extensions = append(result.Extensions, ext)

		switch {
		case ext.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 19}): // Basic Constraints
			var bc struct {
				IsCA       bool `asn1:"optional"`
				MaxPathLen int  `asn1:"optional,default:-1"`
			}
			if _, err := asn1.Unmarshal(ext.Value, &bc); err == nil {
				result.BasicConstraintsValid = true
				result.IsCA = bc.IsCA
				if bc.MaxPathLen >= 0 {
					result.MaxPathLen = bc.MaxPathLen
					if bc.MaxPathLen == 0 {
						result.MaxPathLenZero = true
					}
				}
			}

		case ext.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 14}): // Subject Key Identifier
			var skid []byte
			if _, err := asn1.Unmarshal(ext.Value, &skid); err == nil {
				result.SubjectKeyId = skid
			}

		case ext.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 35}): // Authority Key Identifier
			var akid struct {
				KeyIdentifier []byte `asn1:"optional,tag:0"`
			}
			if _, err := asn1.Unmarshal(ext.Value, &akid); err == nil {
				result.AuthorityKeyId = akid.KeyIdentifier
			}

		case ext.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 15}): // Key Usage
			var bitString asn1.BitString
			if _, err := asn1.Unmarshal(ext.Value, &bitString); err == nil {
				var usage int
				for i, b := range bitString.Bytes {
					for bit := 0; bit < 8; bit++ {
						if (b>>uint(7-bit))&1 != 0 {
							usage |= 1 << uint(i*8+bit)
						}
					}
				}
				result.KeyUsage = x509.KeyUsage(usage)
			}

		case ext.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 17}): // Subject Alternative Names
			parseSANExtension(ext.Value, result)

		case ext.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 37}): // Extended Key Usage
			var oids []asn1.ObjectIdentifier
			if _, err := asn1.Unmarshal(ext.Value, &oids); err == nil {
				for _, oid := range oids {
					if eku, ok := oidToExtKeyUsage(oid); ok {
						result.ExtKeyUsage = append(result.ExtKeyUsage, eku)
					} else {
						result.UnknownExtKeyUsage = append(result.UnknownExtKeyUsage, oid)
					}
				}
			}

		default:
			result.ExtraExtensions = append(result.ExtraExtensions, ext)
		}
	}

	return result, nil
}

// parseSANExtension parses the Subject Alternative Name extension value
// and populates the certificate's DNSNames, EmailAddresses, IPAddresses,
// and URIs fields.
func parseSANExtension(value []byte, cert *x509.Certificate) {
	var rawValues []asn1.RawValue
	if _, err := asn1.Unmarshal(value, &rawValues); err != nil {
		return
	}
	for _, v := range rawValues {
		switch v.Tag {
		case 1: // rfc822Name (email)
			cert.EmailAddresses = append(cert.EmailAddresses, string(v.Bytes))
		case 2: // dNSName
			cert.DNSNames = append(cert.DNSNames, string(v.Bytes))
		case 6: // uniformResourceIdentifier
			if u, err := url.Parse(string(v.Bytes)); err == nil {
				cert.URIs = append(cert.URIs, u)
			}
		case 7: // iPAddress
			cert.IPAddresses = append(cert.IPAddresses, net.IP(v.Bytes))
		}
	}
}

// oidToExtKeyUsage maps a known EKU OID to x509.ExtKeyUsage.
func oidToExtKeyUsage(oid asn1.ObjectIdentifier) (x509.ExtKeyUsage, bool) {
	switch {
	case oid.Equal(asn1.ObjectIdentifier{2, 5, 29, 37, 0}):
		return x509.ExtKeyUsageAny, true
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 1}):
		return x509.ExtKeyUsageServerAuth, true
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 2}):
		return x509.ExtKeyUsageClientAuth, true
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 3}):
		return x509.ExtKeyUsageCodeSigning, true
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 4}):
		return x509.ExtKeyUsageEmailProtection, true
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 8}):
		return x509.ExtKeyUsageTimeStamping, true
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 9}):
		return x509.ExtKeyUsageOCSPSigning, true
	default:
		return 0, false
	}
}

// marshalCSRSANExtension builds the SAN extension DER for a CSR template.
func marshalCSRSANExtension(template *x509.CertificateRequest) ([]byte, error) {
	if len(template.DNSNames) == 0 && len(template.EmailAddresses) == 0 &&
		len(template.IPAddresses) == 0 && len(template.URIs) == 0 {
		return nil, nil
	}

	var rawValues []asn1.RawValue

	for _, name := range template.DNSNames {
		rawValues = append(rawValues, asn1.RawValue{
			Tag:   2,
			Class: asn1.ClassContextSpecific,
			Bytes: []byte(name),
		})
	}

	for _, email := range template.EmailAddresses {
		rawValues = append(rawValues, asn1.RawValue{
			Tag:   1,
			Class: asn1.ClassContextSpecific,
			Bytes: []byte(email),
		})
	}

	for _, ip := range template.IPAddresses {
		rawIP := ip.To4()
		if rawIP == nil {
			rawIP = ip.To16()
		}
		rawValues = append(rawValues, asn1.RawValue{
			Tag:   7,
			Class: asn1.ClassContextSpecific,
			Bytes: rawIP,
		})
	}

	for _, uri := range template.URIs {
		rawValues = append(rawValues, asn1.RawValue{
			Tag:   6,
			Class: asn1.ClassContextSpecific,
			Bytes: []byte(uri.String()),
		})
	}

	return asn1.Marshal(rawValues)
}
