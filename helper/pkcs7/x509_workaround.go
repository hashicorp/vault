// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pkcs7

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"

	ctasn "github.com/google/certificate-transparency-go/asn1"
	ctx509 "github.com/google/certificate-transparency-go/x509"
	ctpkix "github.com/google/certificate-transparency-go/x509/pkix"
)

func CustomParseX509Certificates(data []byte) ([]*x509.Certificate, error) {
	// The Intune client generate a self-signed certificate with a critical authority key extension,
	// let's try to work-around that
	certificates, err := ctx509.ParseCertificates(data)
	if err != nil {
		return nil, err
	}

	converted := make([]*x509.Certificate, len(certificates))
	for i, ctCert := range certificates {
		converted[i] = convertCtCert(ctCert)
	}

	return converted, nil
}

func convertCtCert(ctxCert *ctx509.Certificate) *x509.Certificate {
	return &x509.Certificate{
		Raw:                         ctxCert.Raw,
		RawTBSCertificate:           ctxCert.RawTBSCertificate,
		RawSubjectPublicKeyInfo:     ctxCert.RawSubjectPublicKeyInfo,
		RawSubject:                  ctxCert.RawSubject,
		RawIssuer:                   ctxCert.RawIssuer,
		Signature:                   ctxCert.Signature,
		SignatureAlgorithm:          x509.SignatureAlgorithm(ctxCert.SignatureAlgorithm),
		PublicKeyAlgorithm:          x509.PublicKeyAlgorithm(ctxCert.PublicKeyAlgorithm),
		PublicKey:                   ctxCert.PublicKey,
		Version:                     ctxCert.Version,
		SerialNumber:                ctxCert.SerialNumber,
		Issuer:                      convertPkixName(ctxCert.Issuer),
		Subject:                     convertPkixName(ctxCert.Subject),
		NotBefore:                   ctxCert.NotBefore,
		NotAfter:                    ctxCert.NotAfter,
		KeyUsage:                    x509.KeyUsage(ctxCert.KeyUsage),
		Extensions:                  convertExtensions(ctxCert.Extensions),
		ExtraExtensions:             convertExtensions(ctxCert.ExtraExtensions),
		UnhandledCriticalExtensions: convertAsn1ObjectIdentifiers(ctxCert.UnhandledCriticalExtensions),
		ExtKeyUsage:                 convertExtKeyUsages(ctxCert.ExtKeyUsage),
		UnknownExtKeyUsage:          nil,
		BasicConstraintsValid:       ctxCert.BasicConstraintsValid,
		IsCA:                        ctxCert.IsCA,
		MaxPathLen:                  ctxCert.MaxPathLen,
		MaxPathLenZero:              ctxCert.MaxPathLenZero,
		SubjectKeyId:                ctxCert.SubjectKeyId,
		AuthorityKeyId:              ctxCert.AuthorityKeyId,
		OCSPServer:                  ctxCert.OCSPServer,
		IssuingCertificateURL:       ctxCert.IssuingCertificateURL,
		DNSNames:                    ctxCert.DNSNames,
		EmailAddresses:              ctxCert.EmailAddresses,
		IPAddresses:                 ctxCert.IPAddresses,
		URIs:                        ctxCert.URIs,
		PermittedDNSDomainsCritical: ctxCert.PermittedDNSDomainsCritical,
		PermittedDNSDomains:         ctxCert.PermittedDNSDomains,
		ExcludedDNSDomains:          ctxCert.ExcludedDNSDomains,
		PermittedIPRanges:           ctxCert.PermittedIPRanges,
		ExcludedIPRanges:            ctxCert.ExcludedIPRanges,
		PermittedEmailAddresses:     ctxCert.PermittedEmailAddresses,
		ExcludedEmailAddresses:      ctxCert.ExcludedEmailAddresses,
		PermittedURIDomains:         ctxCert.PermittedURIDomains,
		ExcludedURIDomains:          ctxCert.ExcludedURIDomains,
		CRLDistributionPoints:       ctxCert.CRLDistributionPoints,
		PolicyIdentifiers:           convertAsn1ObjectIdentifiers(ctxCert.PolicyIdentifiers),
		// Newer fields that are not in Google Transparency
		Policies:                  nil,
		InhibitAnyPolicy:          0,
		InhibitAnyPolicyZero:      false,
		InhibitPolicyMapping:      0,
		InhibitPolicyMappingZero:  false,
		RequireExplicitPolicy:     0,
		RequireExplicitPolicyZero: false,
		PolicyMappings:            nil,
	}
}

func convertExtKeyUsages(usages []ctx509.ExtKeyUsage) []x509.ExtKeyUsage {
	converted := make([]x509.ExtKeyUsage, len(usages))
	for i, extKeyUsage := range usages {
		converted[i] = x509.ExtKeyUsage(extKeyUsage)
	}
	return converted
}

func convertAsn1ObjectIdentifiers(extensions []ctasn.ObjectIdentifier) []asn1.ObjectIdentifier {
	converted := make([]asn1.ObjectIdentifier, len(extensions))
	for i, extension := range extensions {
		converted[i] = asn1.ObjectIdentifier(extension)
	}
	return converted
}

func convertExtensions(extensions []ctpkix.Extension) []pkix.Extension {
	converted := make([]pkix.Extension, len(extensions))
	for i, ext := range extensions {
		converted[i] = pkix.Extension{
			Id:       asn1.ObjectIdentifier(ext.Id),
			Critical: ext.Critical,
			Value:    ext.Value,
		}
	}
	return converted
}

func convertPkixName(name ctpkix.Name) pkix.Name {
	return pkix.Name{
		Country:            name.Country,
		Organization:       name.Organization,
		OrganizationalUnit: name.OrganizationalUnit,
		Locality:           name.Locality,
		Province:           name.Province,
		StreetAddress:      name.StreetAddress,
		PostalCode:         name.PostalCode,
		SerialNumber:       name.SerialNumber,
		CommonName:         name.CommonName,
		Names:              convertAttributeTypeAndValue(name.Names),
		ExtraNames:         convertAttributeTypeAndValue(name.ExtraNames),
	}
}

func convertAttributeTypeAndValue(names []ctpkix.AttributeTypeAndValue) []pkix.AttributeTypeAndValue {
	converted := make([]pkix.AttributeTypeAndValue, len(names))
	for i, name := range names {
		converted[i] = pkix.AttributeTypeAndValue{
			Type:  asn1.ObjectIdentifier(name.Type),
			Value: name.Value,
		}
	}
	return converted
}
