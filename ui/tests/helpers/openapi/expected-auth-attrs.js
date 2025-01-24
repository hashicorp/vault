/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// The constants within this file represent the expected model attributes as parsed from OpenAPI
// if changes are made to the OpenAPI spec, that may result in changes that must be reflected
// here AND ensured to not cause breaking changes within the UI.

const userpass = {
  user: {
    username: {
      editType: 'string',
      helpText: 'Username for this user.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Username',
      type: 'string',
    },
    password: {
      editType: 'string',
      helpText: 'Password for this user.',
      fieldGroup: 'default',
      sensitive: true,
      type: 'string',
    },
    tokenBoundCidrs: {
      editType: 'stringArray',
      helpText:
        'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Bound CIDRs",
    },
    tokenExplicitMaxTtl: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Explicit Maximum TTL",
    },
    tokenMaxTtl: {
      editType: 'ttl',
      helpText: 'The maximum lifetime of the generated token',
      fieldGroup: 'Tokens',
      label: "Generated Token's Maximum TTL",
    },
    tokenNoDefaultPolicy: {
      editType: 'boolean',
      helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
      fieldGroup: 'Tokens',
      label: "Do Not Attach 'default' Policy To Generated Tokens",
      type: 'boolean',
    },
    tokenNumUses: {
      editType: 'number',
      helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
      fieldGroup: 'Tokens',
      label: 'Maximum Uses of Generated Tokens',
      type: 'number',
    },
    tokenPeriod: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
      fieldGroup: 'Tokens',
      label: "Generated Token's Period",
    },
    tokenPolicies: {
      editType: 'stringArray',
      helpText: 'A list of policies that will apply to the generated token for this user.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Policies",
    },
    tokenTtl: {
      editType: 'ttl',
      helpText: 'The initial ttl of the token to generate',
      fieldGroup: 'Tokens',
      label: "Generated Token's Initial TTL",
    },
    tokenType: {
      editType: 'string',
      helpText: 'The type of token to generate, service or batch',
      fieldGroup: 'Tokens',
      label: "Generated Token's Type",
      type: 'string',
    },
  },
};

const azure = {
  'auth-config/azure': {
    clientId: {
      editType: 'string',
      fieldGroup: 'default',
      helpText:
        'The OAuth2 client id to connection to Azure. This value can also be provided with the AZURE_CLIENT_ID environment variable.',
      label: 'Client ID',
      type: 'string',
    },
    clientSecret: {
      editType: 'string',
      fieldGroup: 'default',
      helpText:
        'The OAuth2 client secret to connection to Azure. This value can also be provided with the AZURE_CLIENT_SECRET environment variable.',
      type: 'string',
    },
    environment: {
      editType: 'string',
      fieldGroup: 'default',
      helpText:
        'The Azure environment name. If not provided, AzurePublicCloud is used. This value can also be provided with the AZURE_ENVIRONMENT environment variable.',
      type: 'string',
    },
    maxRetries: {
      editType: 'number',
      fieldGroup: 'default',
      helpText:
        'The maximum number of attempts a failed operation will be retried before producing an error.',
      type: 'number',
    },
    maxRetryDelay: {
      editType: 'ttl',
      fieldGroup: 'default',
      helpText: 'The maximum delay allowed before retrying an operation.',
    },
    resource: {
      editType: 'string',
      fieldGroup: 'default',
      helpText:
        'The resource URL for the vault application in Azure Active Directory. This value can also be provided with the AZURE_AD_RESOURCE environment variable.',
      type: 'string',
    },
    retryDelay: {
      editType: 'ttl',
      fieldGroup: 'default',
      helpText: 'The initial amount of delay to use before retrying an operation, increasing exponentially.',
    },
    rootPasswordTtl: {
      editType: 'ttl',
      fieldGroup: 'default',
      helpText:
        'The TTL of the root password in Azure. This can be either a number of seconds or a time formatted duration (ex: 24h, 48ds)',
    },
    tenantId: {
      editType: 'string',
      fieldGroup: 'default',
      helpText:
        'The tenant id for the Azure Active Directory. This is sometimes referred to as Directory ID in AD. This value can also be provided with the AZURE_TENANT_ID environment variable.',
      label: 'Tenant ID',
      type: 'string',
    },
    identityTokenAudience: {
      editType: 'string',
      fieldGroup: 'default',
      helpText: 'Audience of plugin identity tokens',
      type: 'string',
    },
    identityTokenTtl: {
      editType: 'ttl',
      fieldGroup: 'default',
      helpText: 'Time-to-live of plugin identity tokens',
    },
  },
};

const cert = {
  'auth-config/cert': {
    disableBinding: {
      editType: 'boolean',
      helpText:
        'If set, during renewal, skips the matching of presented client identity with the client identity used during login. Defaults to false.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    enableIdentityAliasMetadata: {
      editType: 'boolean',
      helpText:
        'If set, metadata of the certificate including the metadata corresponding to allowed_metadata_extensions will be stored in the alias. Defaults to false.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    enableMetadataOnFailures: {
      editType: 'boolean',
      fieldGroup: 'default',
      helpText: 'If set, metadata of the client certificate will be returned on authentication failures.',
      type: 'boolean',
    },
    ocspCacheSize: {
      editType: 'number',
      helpText: 'The size of the in memory OCSP response cache, shared by all configured certs',
      fieldGroup: 'default',
      type: 'number',
    },
    roleCacheSize: {
      editType: 'number',
      fieldGroup: 'default',
      helpText: 'The size of the in memory role cache',
      type: 'number',
    },
  },
  cert: {
    name: {
      editType: 'string',
      helpText: 'The name of the certificate',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Name',
      type: 'string',
    },
    allowedCommonNames: {
      editType: 'stringArray',
      helpText: 'A list of names. At least one must exist in the Common Name. Supports globbing.',
      fieldGroup: 'Constraints',
    },
    allowedDnsSans: {
      editType: 'stringArray',
      helpText: 'A list of DNS names. At least one must exist in the SANs. Supports globbing.',
      fieldGroup: 'Constraints',
      label: 'Allowed DNS SANs',
    },
    allowedEmailSans: {
      editType: 'stringArray',
      helpText: 'A list of Email Addresses. At least one must exist in the SANs. Supports globbing.',
      fieldGroup: 'Constraints',
      label: 'Allowed Email SANs',
    },
    allowedMetadataExtensions: {
      editType: 'stringArray',
      helpText:
        'A list of OID extensions. Upon successful authentication, these extensions will be added as metadata if they are present in the certificate. The metadata key will be the string consisting of the OID numbers separated by a dash (-) instead of a dot (.) to allow usage in ACL templates.',
      fieldGroup: 'default',
    },
    allowedNames: {
      editType: 'stringArray',
      helpText:
        'A list of names. At least one must exist in either the Common Name or SANs. Supports globbing. This parameter is deprecated, please use allowed_common_names, allowed_dns_sans, allowed_email_sans, allowed_uri_sans.',
      fieldGroup: 'Constraints',
    },
    allowedOrganizationalUnits: {
      editType: 'stringArray',
      helpText: 'A list of Organizational Units names. At least one must exist in the OU field.',
      fieldGroup: 'Constraints',
    },
    allowedUriSans: {
      editType: 'stringArray',
      helpText: 'A list of URIs. At least one must exist in the SANs. Supports globbing.',
      fieldGroup: 'Constraints',
      label: 'Allowed URI SANs',
    },
    certificate: {
      editType: 'file',
      helpText: 'The public certificate that should be trusted. Must be x509 PEM encoded.',
      fieldGroup: 'default',
      type: 'string',
    },
    displayName: {
      editType: 'string',
      helpText: 'The display name to use for clients using this certificate.',
      fieldGroup: 'default',
      type: 'string',
    },
    ocspCaCertificates: {
      editType: 'file',
      helpText: 'Any additional CA certificates needed to communicate with OCSP servers',
      fieldGroup: 'default',
      type: 'string',
    },
    ocspEnabled: {
      editType: 'boolean',
      helpText: 'Whether to attempt OCSP verification of certificates at login',
      fieldGroup: 'default',
      type: 'boolean',
    },
    ocspFailOpen: {
      editType: 'boolean',
      helpText:
        'If set to true, if an OCSP revocation cannot be made successfully, login will proceed rather than failing. If false, failing to get an OCSP status fails the request.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    ocspQueryAllServers: {
      editType: 'boolean',
      helpText:
        'If set to true, rather than accepting the first successful OCSP response, query all servers and consider the certificate valid only if all servers agree.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    ocspServersOverride: {
      editType: 'stringArray',
      helpText:
        'A list of OCSP server addresses. If unset, the OCSP server is determined from the AuthorityInformationAccess extension on the certificate being inspected.',
      fieldGroup: 'default',
    },
    requiredExtensions: {
      editType: 'stringArray',
      helpText:
        "A list of extensions formatted as 'oid:value'. Expects the extension value to be some type of ASN1 encoded string. All values much match. Supports globbing on 'value'.",
      fieldGroup: 'default',
    },
    tokenBoundCidrs: {
      editType: 'stringArray',
      helpText:
        'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Bound CIDRs",
    },
    tokenExplicitMaxTtl: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Explicit Maximum TTL",
    },
    tokenMaxTtl: {
      editType: 'ttl',
      helpText: 'The maximum lifetime of the generated token',
      fieldGroup: 'Tokens',
      label: "Generated Token's Maximum TTL",
    },
    tokenNoDefaultPolicy: {
      editType: 'boolean',
      helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
      fieldGroup: 'Tokens',
      label: "Do Not Attach 'default' Policy To Generated Tokens",
      type: 'boolean',
    },
    tokenNumUses: {
      editType: 'number',
      helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
      fieldGroup: 'Tokens',
      label: 'Maximum Uses of Generated Tokens',
      type: 'number',
    },
    tokenPeriod: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
      fieldGroup: 'Tokens',
      label: "Generated Token's Period",
    },
    tokenPolicies: {
      editType: 'stringArray',
      helpText: 'A list of policies that will apply to the generated token for this user.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Policies",
    },
    tokenTtl: {
      editType: 'ttl',
      helpText: 'The initial ttl of the token to generate',
      fieldGroup: 'Tokens',
      label: "Generated Token's Initial TTL",
    },
    tokenType: {
      editType: 'string',
      helpText: 'The type of token to generate, service or batch',
      fieldGroup: 'Tokens',
      label: "Generated Token's Type",
      type: 'string',
    },
  },
};

const gcp = {
  'auth-config/gcp': {
    credentials: {
      editType: 'string',
      helpText:
        'Google credentials JSON that Vault will use to verify users against GCP APIs. If not specified, will use application default credentials',
      fieldGroup: 'default',
      label: 'Credentials',
      type: 'string',
    },
    customEndpoint: {
      editType: 'object',
      helpText: 'Specifies overrides for various Google API Service Endpoints used in requests.',
      fieldGroup: 'default',
      type: 'object',
    },
    gceAlias: {
      editType: 'string',
      helpText: 'Indicates what value to use when generating an alias for GCE authentications.',
      fieldGroup: 'default',
      type: 'string',
    },
    gceMetadata: {
      editType: 'stringArray',
      helpText:
        "The metadata to include on the aliases and audit logs generated by this plugin. When set to 'default', includes: instance_creation_timestamp, instance_id, instance_name, project_id, project_number, role, service_account_id, service_account_email, zone. Not editing this field means the 'default' fields are included. Explicitly setting this field to empty overrides the 'default' and means no metadata will be included. If not using 'default', explicit fields must be sent like: 'field1,field2'.",
      fieldGroup: 'default',
      defaultValue: 'field1,field2',
      label: 'gce_metadata',
    },
    iamAlias: {
      editType: 'string',
      helpText: 'Indicates what value to use when generating an alias for IAM authentications.',
      fieldGroup: 'default',
      type: 'string',
    },
    iamMetadata: {
      editType: 'stringArray',
      helpText:
        "The metadata to include on the aliases and audit logs generated by this plugin. When set to 'default', includes: project_id, role, service_account_id, service_account_email. Not editing this field means the 'default' fields are included. Explicitly setting this field to empty overrides the 'default' and means no metadata will be included. If not using 'default', explicit fields must be sent like: 'field1,field2'.",
      fieldGroup: 'default',
      defaultValue: 'field1,field2',
      label: 'iam_metadata',
    },
    identityTokenAudience: {
      editType: 'string',
      fieldGroup: 'default',
      helpText: 'Audience of plugin identity tokens',
      type: 'string',
    },
    identityTokenTtl: {
      editType: 'ttl',
      fieldGroup: 'default',
      helpText: 'Time-to-live of plugin identity tokens',
    },
    serviceAccountEmail: {
      editType: 'string',
      fieldGroup: 'default',
      helpText: 'Email ID for the Service Account to impersonate for Workload Identity Federation.',
      type: 'string',
    },
  },
};

const github = {
  'auth-config/github': {
    baseUrl: {
      editType: 'string',
      helpText:
        'The API endpoint to use. Useful if you are running GitHub Enterprise or an API-compatible authentication server.',
      fieldGroup: 'GitHub Options',
      label: 'Base URL',
      type: 'string',
    },
    organization: {
      editType: 'string',
      helpText: 'The organization users must be part of',
      fieldGroup: 'default',
      type: 'string',
    },
    organizationId: {
      editType: 'number',
      helpText: 'The ID of the organization users must be part of',
      fieldGroup: 'default',
      type: 'number',
    },
    tokenBoundCidrs: {
      editType: 'stringArray',
      helpText:
        'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Bound CIDRs",
    },
    tokenExplicitMaxTtl: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Explicit Maximum TTL",
    },
    tokenMaxTtl: {
      editType: 'ttl',
      helpText: 'The maximum lifetime of the generated token',
      fieldGroup: 'Tokens',
      label: "Generated Token's Maximum TTL",
    },
    tokenNoDefaultPolicy: {
      editType: 'boolean',
      helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
      fieldGroup: 'Tokens',
      label: "Do Not Attach 'default' Policy To Generated Tokens",
      type: 'boolean',
    },
    tokenNumUses: {
      editType: 'number',
      helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
      fieldGroup: 'Tokens',
      label: 'Maximum Uses of Generated Tokens',
      type: 'number',
    },
    tokenPeriod: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
      fieldGroup: 'Tokens',
      label: "Generated Token's Period",
    },
    tokenPolicies: {
      editType: 'stringArray',
      helpText: 'A list of policies that will apply to the generated token for this user.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Policies",
    },
    tokenTtl: {
      editType: 'ttl',
      helpText: 'The initial ttl of the token to generate',
      fieldGroup: 'Tokens',
      label: "Generated Token's Initial TTL",
    },
    tokenType: {
      editType: 'string',
      helpText: 'The type of token to generate, service or batch',
      fieldGroup: 'Tokens',
      label: "Generated Token's Type",
      type: 'string',
    },
  },
};

const jwt = {
  'auth-config/jwt': {
    boundIssuer: {
      editType: 'string',
      helpText: "The value against which to match the 'iss' claim in a JWT. Optional.",
      fieldGroup: 'default',
      type: 'string',
    },
    defaultRole: {
      editType: 'string',
      helpText:
        'The default role to use if none is provided during login. If not set, a role is required during login.',
      fieldGroup: 'default',
      type: 'string',
    },
    jwksCaPem: {
      editType: 'string',
      helpText:
        'The CA certificate or chain of certificates, in PEM format, to use to validate connections to the JWKS URL. If not set, system certificates are used.',
      fieldGroup: 'default',
      type: 'string',
    },
    jwksPairs: {
      editType: 'objectArray',
      fieldGroup: 'default',
      helpText:
        'Set of JWKS Url and CA certificate (or chain of certificates) pairs. CA certificates must be in PEM format. Cannot be used with "jwks_url" or "jwks_ca_pem".',
    },
    jwksUrl: {
      editType: 'string',
      helpText:
        'JWKS URL to use to authenticate signatures. Cannot be used with "oidc_discovery_url" or "jwt_validation_pubkeys".',
      fieldGroup: 'default',
      type: 'string',
    },
    jwtSupportedAlgs: {
      editType: 'stringArray',
      helpText: 'A list of supported signing algorithms. Defaults to RS256.',
      fieldGroup: 'default',
    },
    jwtValidationPubkeys: {
      editType: 'stringArray',
      helpText:
        'A list of PEM-encoded public keys to use to authenticate signatures locally. Cannot be used with "jwks_url" or "oidc_discovery_url".',
      fieldGroup: 'default',
    },
    namespaceInState: {
      editType: 'boolean',
      helpText:
        'Pass namespace in the OIDC state parameter instead of as a separate query parameter. With this setting, the allowed redirect URL(s) in Vault and on the provider side should not contain a namespace query parameter. This means only one redirect URL entry needs to be maintained on the provider side for all vault namespaces that will be authenticating against it. Defaults to true for new configs.',
      fieldGroup: 'default',
      defaultValue: true,
      label: 'Namespace in OIDC state',
      type: 'boolean',
    },
    oidcClientId: {
      editType: 'string',
      helpText: 'The OAuth Client ID configured with your OIDC provider.',
      fieldGroup: 'default',
      type: 'string',
    },
    oidcClientSecret: {
      editType: 'string',
      helpText: 'The OAuth Client Secret configured with your OIDC provider.',
      fieldGroup: 'default',
      sensitive: true,
      type: 'string',
    },
    oidcDiscoveryCaPem: {
      editType: 'string',
      helpText:
        'The CA certificate or chain of certificates, in PEM format, to use to validate connections to the OIDC Discovery URL. If not set, system certificates are used.',
      fieldGroup: 'default',
      type: 'string',
    },
    oidcDiscoveryUrl: {
      editType: 'string',
      helpText:
        'OIDC Discovery URL, without any .well-known component (base path). Cannot be used with "jwks_url" or "jwt_validation_pubkeys".',
      fieldGroup: 'default',
      type: 'string',
    },
    oidcResponseMode: {
      editType: 'string',
      helpText:
        "The response mode to be used in the OAuth2 request. Allowed values are 'query' and 'form_post'.",
      fieldGroup: 'default',
      type: 'string',
    },
    oidcResponseTypes: {
      editType: 'stringArray',
      helpText:
        "The response types to request. Allowed values are 'code' and 'id_token'. Defaults to 'code'.",
      fieldGroup: 'default',
    },
    providerConfig: {
      editType: 'object',
      helpText: 'Provider-specific configuration. Optional.',
      fieldGroup: 'default',
      label: 'Provider Config',
      type: 'object',
    },
    unsupportedCriticalCertExtensions: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'A list of ASN1 OIDs of certificate extensions marked Critical that are unsupported by Vault and should be ignored. This option should very rarely be needed except in specialized PKI environments.',
    },
  },
};

const kubernetes = {
  'auth-config/kubernetes': {
    disableLocalCaJwt: {
      editType: 'boolean',
      helpText:
        'Disable defaulting to the local CA cert and service account JWT when running in a Kubernetes pod',
      fieldGroup: 'default',
      label: 'Disable use of local CA and service account JWT',
      type: 'boolean',
    },
    kubernetesCaCert: {
      editType: 'string',
      helpText:
        "Optional PEM encoded CA cert for use by the TLS client used to talk with the API. If it is not set and disable_local_ca_jwt is true, the system's trusted CA certificate pool will be used.",
      fieldGroup: 'default',
      label: 'Kubernetes CA Certificate',
      type: 'string',
    },
    kubernetesHost: {
      editType: 'string',
      helpText:
        'Host must be a host string, a host:port pair, or a URL to the base of the Kubernetes API server.',
      fieldGroup: 'default',
      type: 'string',
    },
    pemKeys: {
      editType: 'stringArray',
      helpText:
        'Optional list of PEM-formated public keys or certificates used to verify the signatures of kubernetes service account JWTs. If a certificate is given, its public key will be extracted. Not every installation of Kubernetes exposes these keys.',
      fieldGroup: 'default',
      label: 'Service account verification keys',
    },
    tokenReviewerJwt: {
      editType: 'string',
      helpText:
        'A service account JWT (or other token) used as a bearer token to access the TokenReview API to validate other JWTs during login. If not set the JWT used for login will be used to access the API.',
      fieldGroup: 'default',
      label: 'Token Reviewer JWT',
      type: 'string',
    },
    useAnnotationsAsAliasMetadata: {
      editType: 'boolean',
      fieldGroup: 'default',
      helpText:
        'Use annotations from the client token\'s associated service account as alias metadata for the Vault entity. Only annotations with the prefix "vault.hashicorp.com/alias-metadata-" will be used. Note that Vault will need permission to read service accounts from the Kubernetes API.',
      label: 'Use annotations of JWT service account as alias metadata',
      type: 'boolean',
    },
  },
  role: {
    name: {
      editType: 'string',
      helpText: 'Name of the role.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Name',
      type: 'string',
    },
    aliasNameSource: {
      editType: 'string',
      helpText:
        'Source to use when deriving the Alias name. valid choices: "serviceaccount_uid" : <token.uid> e.g. 474b11b5-0f20-4f9d-8ca5-65715ab325e0 (most secure choice) "serviceaccount_name" : <namespace>/<serviceaccount> e.g. vault/vault-agent default: "serviceaccount_uid"',
      fieldGroup: 'default',
      type: 'string',
    },
    audience: {
      editType: 'string',
      helpText: 'Optional Audience claim to verify in the jwt.',
      fieldGroup: 'default',
      type: 'string',
    },
    boundServiceAccountNames: {
      editType: 'stringArray',
      helpText:
        'List of service account names able to access this role. If set to "*" all names are allowed.',
      fieldGroup: 'default',
    },
    boundServiceAccountNamespaces: {
      editType: 'stringArray',
      helpText: 'List of namespaces allowed to access this role. If set to "*" all namespaces are allowed.',
      fieldGroup: 'default',
    },
    tokenBoundCidrs: {
      editType: 'stringArray',
      helpText:
        'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Bound CIDRs",
    },
    tokenExplicitMaxTtl: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Explicit Maximum TTL",
    },
    tokenMaxTtl: {
      editType: 'ttl',
      helpText: 'The maximum lifetime of the generated token',
      fieldGroup: 'Tokens',
      label: "Generated Token's Maximum TTL",
    },
    tokenNoDefaultPolicy: {
      editType: 'boolean',
      helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
      fieldGroup: 'Tokens',
      label: "Do Not Attach 'default' Policy To Generated Tokens",
      type: 'boolean',
    },
    tokenNumUses: {
      editType: 'number',
      helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
      fieldGroup: 'Tokens',
      label: 'Maximum Uses of Generated Tokens',
      type: 'number',
    },
    tokenPeriod: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
      fieldGroup: 'Tokens',
      label: "Generated Token's Period",
    },
    tokenPolicies: {
      editType: 'stringArray',
      helpText: 'A list of policies that will apply to the generated token for this user.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Policies",
    },
    tokenTtl: {
      editType: 'ttl',
      helpText: 'The initial ttl of the token to generate',
      fieldGroup: 'Tokens',
      label: "Generated Token's Initial TTL",
    },
    tokenType: {
      editType: 'string',
      helpText: 'The type of token to generate, service or batch',
      fieldGroup: 'Tokens',
      label: "Generated Token's Type",
      type: 'string',
    },
  },
};

const ldap = {
  'auth-config/ldap': {
    anonymousGroupSearch: {
      editType: 'boolean',
      helpText:
        'Use anonymous binds when performing LDAP group searches (if true the initial credentials will still be used for the initial connection test).',
      fieldGroup: 'default',
      label: 'Anonymous group search',
      type: 'boolean',
    },
    binddn: {
      editType: 'string',
      helpText: 'LDAP DN for searching for the user DN (optional)',
      fieldGroup: 'default',
      label: 'Name of Object to bind (binddn)',
      type: 'string',
    },
    bindpass: {
      editType: 'string',
      helpText: 'LDAP password for searching for the user DN (optional)',
      fieldGroup: 'default',
      sensitive: true,
      type: 'string',
    },
    caseSensitiveNames: {
      editType: 'boolean',
      helpText:
        'If true, case sensitivity will be used when comparing usernames and groups for matching policies.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    certificate: {
      editType: 'file',
      helpText:
        'CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded (optional)',
      fieldGroup: 'default',
      label: 'CA certificate',
      type: 'string',
    },
    clientTlsCert: {
      editType: 'file',
      helpText: 'Client certificate to provide to the LDAP server, must be x509 PEM encoded (optional)',
      fieldGroup: 'default',
      label: 'Client certificate',
      type: 'string',
    },
    clientTlsKey: {
      editType: 'file',
      helpText: 'Client certificate key to provide to the LDAP server, must be x509 PEM encoded (optional)',
      fieldGroup: 'default',
      label: 'Client key',
      type: 'string',
    },
    connectionTimeout: {
      editType: 'ttl',
      helpText:
        'Timeout, in seconds, when attempting to connect to the LDAP server before trying the next URL in the configuration.',
      fieldGroup: 'default',
    },
    denyNullBind: {
      editType: 'boolean',
      helpText:
        "Denies an unauthenticated LDAP bind request if the user's password is empty; defaults to true",
      fieldGroup: 'default',
      type: 'boolean',
    },
    dereferenceAliases: {
      editType: 'string',
      helpText:
        "When aliases should be dereferenced on search operations. Accepted values are 'never', 'finding', 'searching', 'always'. Defaults to 'never'.",
      possibleValues: ['never', 'finding', 'searching', 'always'],
      fieldGroup: 'default',
      type: 'string',
    },
    discoverdn: {
      editType: 'boolean',
      helpText: 'Use anonymous bind to discover the bind DN of a user (optional)',
      fieldGroup: 'default',
      label: 'Discover DN',
      type: 'boolean',
    },
    enableSamaccountnameLogin: {
      editType: 'boolean',
      fieldGroup: 'default',
      helpText:
        'If true, matching sAMAccountName attribute values will be allowed to login when upndomain is defined.',
      type: 'boolean',
    },
    groupattr: {
      editType: 'string',
      helpText:
        'LDAP attribute to follow on objects returned by <groupfilter> in order to enumerate user group membership. Examples: "cn" or "memberOf", etc. Default: cn',
      fieldGroup: 'default',
      defaultValue: 'cn',
      label: 'Group Attribute',
      type: 'string',
    },
    groupdn: {
      editType: 'string',
      helpText: 'LDAP search base to use for group membership search (eg: ou=Groups,dc=example,dc=org)',
      fieldGroup: 'default',
      label: 'Group DN',
      type: 'string',
    },
    groupfilter: {
      editType: 'string',
      helpText:
        'Go template for querying group membership of user (optional) The template can access the following context variables: UserDN, Username Example: (&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}})) Default: (|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))',
      fieldGroup: 'default',
      label: 'Group Filter',
      type: 'string',
    },
    insecureTls: {
      editType: 'boolean',
      helpText: 'Skip LDAP server SSL Certificate verification - VERY insecure (optional)',
      fieldGroup: 'default',
      label: 'Insecure TLS',
      type: 'boolean',
    },
    maxPageSize: {
      editType: 'number',
      helpText:
        "If set to a value greater than 0, the LDAP backend will use the LDAP server's paged search control to request pages of up to the given size. This can be used to avoid hitting the LDAP server's maximum result size limit. Otherwise, the LDAP backend will not use the paged search control.",
      fieldGroup: 'default',
      type: 'number',
    },
    passwordPolicy: {
      editType: 'string',
      fieldGroup: 'default',
      helpText: 'Password policy to use to rotate the root password',
      type: 'string',
    },
    requestTimeout: {
      editType: 'ttl',
      helpText:
        'Timeout, in seconds, for the connection when making requests against the server before returning back an error.',
      fieldGroup: 'default',
    },
    starttls: {
      editType: 'boolean',
      helpText: 'Issue a StartTLS command after establishing unencrypted connection (optional)',
      fieldGroup: 'default',
      label: 'Issue StartTLS',
      type: 'boolean',
    },
    tlsMaxVersion: {
      editType: 'string',
      helpText:
        "Maximum TLS version to use. Accepted values are 'tls10', 'tls11', 'tls12' or 'tls13'. Defaults to 'tls12'",
      possibleValues: ['tls10', 'tls11', 'tls12', 'tls13'],
      fieldGroup: 'default',
      label: 'Maximum TLS Version',
      type: 'string',
    },
    tlsMinVersion: {
      editType: 'string',
      helpText:
        "Minimum TLS version to use. Accepted values are 'tls10', 'tls11', 'tls12' or 'tls13'. Defaults to 'tls12'",
      possibleValues: ['tls10', 'tls11', 'tls12', 'tls13'],
      fieldGroup: 'default',
      label: 'Minimum TLS Version',
      type: 'string',
    },
    tokenBoundCidrs: {
      editType: 'stringArray',
      helpText:
        'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Bound CIDRs",
    },
    tokenExplicitMaxTtl: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Explicit Maximum TTL",
    },
    tokenMaxTtl: {
      editType: 'ttl',
      helpText: 'The maximum lifetime of the generated token',
      fieldGroup: 'Tokens',
      label: "Generated Token's Maximum TTL",
    },
    tokenNoDefaultPolicy: {
      editType: 'boolean',
      helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
      fieldGroup: 'Tokens',
      label: "Do Not Attach 'default' Policy To Generated Tokens",
      type: 'boolean',
    },
    tokenNumUses: {
      editType: 'number',
      helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
      fieldGroup: 'Tokens',
      label: 'Maximum Uses of Generated Tokens',
      type: 'number',
    },
    tokenPeriod: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
      fieldGroup: 'Tokens',
      label: "Generated Token's Period",
    },
    tokenPolicies: {
      editType: 'stringArray',
      helpText: 'A list of policies that will apply to the generated token for this user.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Policies",
    },
    tokenTtl: {
      editType: 'ttl',
      helpText: 'The initial ttl of the token to generate',
      fieldGroup: 'Tokens',
      label: "Generated Token's Initial TTL",
    },
    tokenType: {
      editType: 'string',
      helpText: 'The type of token to generate, service or batch',
      fieldGroup: 'Tokens',
      label: "Generated Token's Type",
      type: 'string',
    },
    upndomain: {
      editType: 'string',
      helpText: 'Enables userPrincipalDomain login with [username]@UPNDomain (optional)',
      fieldGroup: 'default',
      label: 'User Principal (UPN) Domain',
      type: 'string',
    },
    url: {
      editType: 'string',
      helpText:
        'LDAP URL to connect to (default: ldap://127.0.0.1). Multiple URLs can be specified by concatenating them with commas; they will be tried in-order.',
      fieldGroup: 'default',
      label: 'URL',
      type: 'string',
    },
    usePre111GroupCnBehavior: {
      editType: 'boolean',
      helpText:
        'In Vault 1.1.1 a fix for handling group CN values of different cases unfortunately introduced a regression that could cause previously defined groups to not be found due to a change in the resulting name. If set true, the pre-1.1.1 behavior for matching group CNs will be used. This is only needed in some upgrade scenarios for backwards compatibility. It is enabled by default if the config is upgraded but disabled by default on new configurations.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    useTokenGroups: {
      editType: 'boolean',
      helpText:
        'If true, use the Active Directory tokenGroups constructed attribute of the user to find the group memberships. This will find all security groups including nested ones.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    userattr: {
      editType: 'string',
      helpText: 'Attribute used for users (default: cn)',
      fieldGroup: 'default',
      defaultValue: 'cn',
      label: 'User Attribute',
      type: 'string',
    },
    userdn: {
      editType: 'string',
      helpText: 'LDAP domain to use for users (eg: ou=People,dc=example,dc=org)',
      fieldGroup: 'default',
      label: 'User DN',
      type: 'string',
    },
    userfilter: {
      editType: 'string',
      helpText:
        'Go template for LDAP user search filer (optional) The template can access the following context variables: UserAttr, Username Default: ({{.UserAttr}}={{.Username}})',
      fieldGroup: 'default',
      label: 'User Search Filter',
      type: 'string',
    },
    usernameAsAlias: {
      editType: 'boolean',
      helpText: 'If true, sets the alias name to the username',
      fieldGroup: 'default',
      type: 'boolean',
    },
  },
  group: {
    name: {
      editType: 'string',
      helpText: 'Name of the LDAP group.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Name',
      type: 'string',
    },
    policies: {
      editType: 'stringArray',
      helpText: 'A list of policies associated to the group.',
      fieldGroup: 'default',
    },
  },
  user: {
    name: {
      editType: 'string',
      helpText: 'Name of the LDAP user.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Name',
      type: 'string',
    },
    groups: {
      editType: 'stringArray',
      helpText: 'A list of additional groups associated with the user.',
      fieldGroup: 'default',
    },
    policies: {
      editType: 'stringArray',
      helpText: 'A list of policies associated with the user.',
      fieldGroup: 'default',
    },
  },
};

const okta = {
  'auth-config/okta': {
    apiToken: {
      editType: 'string',
      helpText: 'Okta API key.',
      fieldGroup: 'default',
      label: 'API Token',
      type: 'string',
    },
    baseUrl: {
      editType: 'string',
      helpText:
        'The base domain to use for the Okta API. When not specified in the configuration, "okta.com" is used.',
      fieldGroup: 'default',
      label: 'Base URL',
      type: 'string',
    },
    bypassOktaMfa: {
      editType: 'boolean',
      helpText:
        'When set true, requests by Okta for a MFA check will be bypassed. This also disallows certain status checks on the account, such as whether the password is expired.',
      fieldGroup: 'default',
      label: 'Bypass Okta MFA',
      type: 'boolean',
    },
    orgName: {
      editType: 'string',
      helpText: 'Name of the organization to be used in the Okta API.',
      fieldGroup: 'default',
      label: 'Organization Name',
      type: 'string',
    },
    tokenBoundCidrs: {
      editType: 'stringArray',
      helpText:
        'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Bound CIDRs",
    },
    tokenExplicitMaxTtl: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Explicit Maximum TTL",
    },
    tokenMaxTtl: {
      editType: 'ttl',
      helpText: 'The maximum lifetime of the generated token',
      fieldGroup: 'Tokens',
      label: "Generated Token's Maximum TTL",
    },
    tokenNoDefaultPolicy: {
      editType: 'boolean',
      helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
      fieldGroup: 'Tokens',
      label: "Do Not Attach 'default' Policy To Generated Tokens",
      type: 'boolean',
    },
    tokenNumUses: {
      editType: 'number',
      helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
      fieldGroup: 'Tokens',
      label: 'Maximum Uses of Generated Tokens',
      type: 'number',
    },
    tokenPeriod: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
      fieldGroup: 'Tokens',
      label: "Generated Token's Period",
    },
    tokenPolicies: {
      editType: 'stringArray',
      helpText: 'A list of policies that will apply to the generated token for this user.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Policies",
    },
    tokenTtl: {
      editType: 'ttl',
      helpText: 'The initial ttl of the token to generate',
      fieldGroup: 'Tokens',
      label: "Generated Token's Initial TTL",
    },
    tokenType: {
      editType: 'string',
      helpText: 'The type of token to generate, service or batch',
      fieldGroup: 'Tokens',
      label: "Generated Token's Type",
      type: 'string',
    },
  },
  group: {
    name: {
      editType: 'string',
      helpText: 'Name of the Okta group.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Name',
      type: 'string',
    },
    policies: {
      editType: 'stringArray',
      helpText: 'A list of policies associated to the group.',
      fieldGroup: 'default',
    },
  },
  user: {
    name: {
      editType: 'string',
      helpText: 'Name of the user.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Name',
      type: 'string',
    },
    groups: {
      editType: 'stringArray',
      helpText: 'List of groups associated with the user.',
      fieldGroup: 'default',
    },
    policies: {
      editType: 'stringArray',
      helpText: 'List of policies associated with the user.',
      fieldGroup: 'default',
    },
  },
};

const radius = {
  'auth-config/radius': {
    dialTimeout: {
      editType: 'ttl',
      helpText: 'Number of seconds before connect times out (default: 10)',
      fieldGroup: 'default',
      defaultValue: 10,
    },
    host: {
      editType: 'string',
      helpText: 'RADIUS server host',
      fieldGroup: 'default',
      label: 'Host',
      type: 'string',
    },
    nasIdentifier: {
      editType: 'string',
      helpText: 'RADIUS NAS Identifier field (optional)',
      fieldGroup: 'default',
      label: 'NAS Identifier',
      type: 'string',
    },
    nasPort: {
      editType: 'number',
      helpText: 'RADIUS NAS port field (default: 10)',
      fieldGroup: 'default',
      defaultValue: 10,
      label: 'NAS Port',
      type: 'number',
    },
    port: {
      editType: 'number',
      helpText: 'RADIUS server port (default: 1812)',
      fieldGroup: 'default',
      defaultValue: 1812,
      type: 'number',
    },
    readTimeout: {
      editType: 'ttl',
      helpText: 'Number of seconds before response times out (default: 10)',
      fieldGroup: 'default',
      defaultValue: 10,
    },
    secret: {
      editType: 'string',
      helpText: 'Secret shared with the RADIUS server',
      fieldGroup: 'default',
      type: 'string',
    },
    tokenBoundCidrs: {
      editType: 'stringArray',
      helpText:
        'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Bound CIDRs",
    },
    tokenExplicitMaxTtl: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Explicit Maximum TTL",
    },
    tokenMaxTtl: {
      editType: 'ttl',
      helpText: 'The maximum lifetime of the generated token',
      fieldGroup: 'Tokens',
      label: "Generated Token's Maximum TTL",
    },
    tokenNoDefaultPolicy: {
      editType: 'boolean',
      helpText: "If true, the 'default' policy will not automatically be added to generated tokens",
      fieldGroup: 'Tokens',
      label: "Do Not Attach 'default' Policy To Generated Tokens",
      type: 'boolean',
    },
    tokenNumUses: {
      editType: 'number',
      helpText: 'The maximum number of times a token may be used, a value of zero means unlimited',
      fieldGroup: 'Tokens',
      label: 'Maximum Uses of Generated Tokens',
      type: 'number',
    },
    tokenPeriod: {
      editType: 'ttl',
      helpText:
        'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
      fieldGroup: 'Tokens',
      label: "Generated Token's Period",
    },
    tokenPolicies: {
      editType: 'stringArray',
      helpText: 'A list of policies that will apply to the generated token for this user.',
      fieldGroup: 'Tokens',
      label: "Generated Token's Policies",
    },
    tokenTtl: {
      editType: 'ttl',
      helpText: 'The initial ttl of the token to generate',
      fieldGroup: 'Tokens',
      label: "Generated Token's Initial TTL",
    },
    tokenType: {
      editType: 'string',
      helpText: 'The type of token to generate, service or batch',
      fieldGroup: 'Tokens',
      label: "Generated Token's Type",
      type: 'string',
    },
    unregisteredUserPolicies: {
      editType: 'string',
      helpText:
        'List of policies to grant upon successful RADIUS authentication of an unregistered user (default: empty)',
      fieldGroup: 'default',
      label: 'Policies for unregistered users',
      type: 'string',
    },
  },
  user: {
    name: {
      editType: 'string',
      helpText: 'Name of the RADIUS user.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Name',
      type: 'string',
    },
    policies: {
      editType: 'stringArray',
      helpText: 'A list of policies associated to the user.',
      fieldGroup: 'default',
    },
  },
};

export default {
  azure,
  userpass,
  cert,
  gcp,
  github,
  jwt,
  kubernetes,
  ldap,
  okta,
  radius,
  // aws is the only method that doesn't leverage OpenApi in practice
};
