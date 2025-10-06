// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package observe

const (
	// ---
	// Generate Root:

	ObservationTypePKIGenerateRoot = "pki/root/generate"

	// ---
	// Generate Intermediate:

	ObservationTypePKIGenerateIntermediate = "pki/intermediate/generate"

	// ---
	// Issue + Sign

	// ObservationTypePKIIssue observations will be emitted for both the issue (pki/issue/:name) and
	// issue-and-sign endpoints (pki/sign/:name). Observations for isssued-and-signed certs will
	// contain "signed" = true, and unsigned certs will contain "signed" = false.
	ObservationTypePKIIssue = "pki/issue"

	// ObservationTypePKICIEPSIssue observations will be emitted for both the CIEPS issue and
	// issue-and-sign endpoints. Observations for issued-and-signed certs will
	// contain "signed" = true, and unsigned certs will contain "signed" = false.
	ObservationTypePKICIEPSIssue = "pki/cieps/issue"

	// ---
	// Issuer Related Observations

	ObservationTypePKIIssuerRead       = "pki/issuer/read"
	ObservationTypePKIIssuerWrite      = "pki/issuer/write"
	ObservationTypePKIIssuerPatch      = "pki/issuer/patch"
	ObservationTypePKIIssuerDelete     = "pki/issuer/delete"
	ObservationTypePKIIssuerResignCRLs = "pki/issuer/resign-crls"
	// ObservationTypePKIIssuersImport is emitted when an import happens for issuers.
	// This can happen via /pki/config/ca, /pki/issuers/import/bundle, /pki/intermediate/set-signed,
	// and others.
	ObservationTypePKIIssuersImport = "pki/issuer/import"
	// ObservationTypePKIConfigIssuersWrite is emitted both for writes to /pki/config/issuers
	// and /pki/root/replace, as they have the same code path.
	ObservationTypePKIConfigIssuersWrite = "pki/config/issuers/write"
	ObservationTypePKIConfigIssuersRead  = "pki/config/issuers/read"

	// ObservationTypePKIReadIssuerCertificate is issued when the issuer's certificate is read,
	// i.e. the following:
	// https://developer.hashicorp.com/vault/api-docs/secret/pki#read-issuer-certificate
	ObservationTypePKIReadIssuerCertificate = "pki/issuer/certificate/read"

	// ---
	// Role related observations

	ObservationTypePKIRoleRead   = "pki/role/read"
	ObservationTypePKIRoleWrite  = "pki/role/write"
	ObservationTypePKIRolePatch  = "pki/role/patch"
	ObservationTypePKIRoleDelete = "pki/role/delete"

	// ---
	// Cert metadata

	// ObservationTypePKIReadCertificateMetadata is emitted when /pki/cert-metadata/:serial is called.
	ObservationTypePKIReadCertificateMetadata = "pki/certificate-metadata/read"

	// ---
	// Tidy

	// ObservationTypePKITidy is emitted when a tidy operation is accepted, not completed.
	ObservationTypePKITidy = "pki/tidy"

	// ---
	// Revoke

	ObservationTypePKIRevoke = "pki/revoke"

	// ---
	// Rotate CRLs

	// ObservationTypePKIRotateCRL is emitted when pki/crl/rotate is called, which forces a rotation of all issuers' CRLs.
	ObservationTypePKIRotateCRL = "pki/crl/rotate"
	// ObservationTypePKIRotateDeltaCRL is emitted when pki/crl/rotate-delta is called, which forces a rotation of all issuers' delta CRLs.
	ObservationTypePKIRotateDeltaCRL = "pki/crl/rotate-delta"

	// ---
	// Key Related Observations

	ObservationTypePKIKeysGenerate    = "pki/keys/generate"
	ObservationTypePKIKeysImport      = "pki/keys/import"
	ObservationTypePKIConfigKeysWrite = "pki/config/keys/write"
	ObservationTypePKIConfigKeysRead  = "pki/config/keys/read"
	ObservationTypePKIKeyRead         = "pki/key/read"
	ObservationTypePKIKeyWrite        = "pki/key/write"
	ObservationTypePKIKeyDelete       = "pki/key/delete"

	// ---
	// OCSP Related Observations
	// Note that statuses are kept to their values as per https://datatracker.ietf.org/doc/html/rfc6960 and
	// are not translated to be 'human-readable'. This observation covers both pki/ocsp and pki/unified-ocsp
	// endpoints, returning a "unified" boolean in the body.

	ObservationTypePKIOCSP = "pki/ocsp"

	// ---
	// Config Related Observations

	// ObservationTypePKIConfigClusterRead will be emitted on a read to
	// pki/config/cluster
	ObservationTypePKIConfigClusterRead = "pki/config/integrations/cluster/read"
	// ObservationTypePKIConfigClusterWrite will be emitted on a write to
	// pki/config/cluster.
	ObservationTypePKIConfigClusterWrite = "pki/config/integrations/cluster/write"

	// ObservationTypePKIConfigIntegrationsGardiumRead will be emitted on a read to
	// pki/config/integrations/gardium. It will not include any user-specified URLs.
	ObservationTypePKIConfigIntegrationsGardiumRead = "pki/config/integrations/gardium/read"
	// ObservationTypePKIConfigIntegrationsGardiumWrite will be emitted on a write to
	// pki/config/integrations/gardium. It will not include any user-specified URLs.
	ObservationTypePKIConfigIntegrationsGardiumWrite = "pki/config/integrations/gardium/write"

	// ObservationTypePKIConfigURLsRead will be emitted on a read to
	// pki/config/urls. It will not include any user-specified URLs.
	ObservationTypePKIConfigURLsRead = "pki/config/urls/read"
	// ObservationTypePKIConfigURLsWrite will be emitted on a write to
	// pki/config/urls. It will not include any user-specified URLs.
	ObservationTypePKIConfigURLsWrite = "pki/config/urls/write"

	// ObservationTypePKIConfigExternalPolicyRead is emitted when a read call goes to
	// pki/config/external-policy (CIEPS).
	ObservationTypePKIConfigExternalPolicyRead = "pki/config/external-policy/read"
	// ObservationTypePKIConfigExternalPolicyWrite is emitted when a write call goes to
	// pki/config/external-policy (CIEPS). Note that any sensitive information, like
	// certificates or URLs.
	ObservationTypePKIConfigExternalPolicyWrite = "pki/config/external-policy/write"

	ObservationTypePKIConfigCRLRead  = "pki/config/crl/read"
	ObservationTypePKIConfigCRLWrite = "pki/config/crl/write"

	// ---
	// ACME Related Observations

	ObservationTypePKIConfigACMERead  = "pki/config/acme/read"
	ObservationTypePKIConfigACMEWrite = "pki/config/acme/write"

	ObservationTypePKIAcmeRevoke         = "pki/acme/revoke"
	ObservationTypePKIAcmeNewOrder       = "pki/acme/order/new-order"
	ObservationTypePKIAcmeListOrders     = "pki/acme/order/list-orders"
	ObservationTypePKIAcmeGetOrder       = "pki/acme/order/get-order"
	ObservationTypePKIAcmeFinalizeOrder  = "pki/acme/order/finalize-order"
	ObservationTypePKIAcmeFetchOrderCert = "pki/acme/order/fetch-order-cert"
	ObservationTypePKIAcmeNewAccount     = "pki/acme/account/new-account"
	ObservationTypePKIAcmeUpdateAccount  = "pki/acme/account/update-account"
	ObservationTypePKIAcmeChallenge      = "pki/acme/challenge"
	ObservationTypePKIAcmeAuthorization  = "pki/acme/authorization"
	ObservationTypePKIAcmeNewEab         = "pki/acme/new-eab"

	// ---
	// EST Related Observations

	ObservationTypePKIConfigESTRead  = "pki/config/est/read"
	ObservationTypePKIConfigESTWrite = "pki/config/est/write"

	ObservationTypePKIESTCACerts = "pki/est/ca-certs"

	ObservationTypePKIESTEnroll   = "pki/est/enroll"
	ObservationTypePKIESTReEnroll = "pki/est/re-enroll"

	// ---
	// CMPv2 Related Observations

	ObservationTypePKIConfigCMPv2Read  = "pki/config/cmpv2/read"
	ObservationTypePKIConfigCMPv2Write = "pki/config/cmpv2/write"

	ObservationTypePKICMPCertRequest = "pki/cmpv2/cert-request"

	// ---
	// SCEP Related Observations

	ObservationTypePKIConfigSCEPRead  = "pki/config/scep/read"
	ObservationTypePKIConfigSCEPWrite = "pki/config/scep/write"

	ObservationTypePKISCEPPKIOperation = "pki/scep/operation/pki"
)
