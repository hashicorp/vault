/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// The constants within this file represent the expected model attributes as parsed from OpenAPI
// if changes are made to the OpenAPI spec, that may result in changes that must be reflected
// here AND ensured to not cause breaking changes within the UI.

const ssh = {
  'role-ssh': {
    role: {
      editType: 'string',
      helpText: '[Required for all types] Name of the role being created.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Role',
      type: 'string',
    },
    algorithmSigner: {
      editType: 'string',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] When supplied, this value specifies a signing algorithm for the key. Possible values: ssh-rsa, rsa-sha2-256, rsa-sha2-512, default, or the empty string.',
      possibleValues: ['', 'default', 'ssh-rsa', 'rsa-sha2-256', 'rsa-sha2-512'],
      fieldGroup: 'default',
      label: 'Signing Algorithm',
      type: 'string',
    },
    allowBareDomains: {
      editType: 'boolean',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If set, host certificates that are requested are allowed to use the base domains listed in "allowed_domains", e.g. "example.com". This is a separate option as in some cases this can be considered a security threat.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowEmptyPrincipals: {
      editType: 'boolean',
      fieldGroup: 'default',
      helpText:
        'Whether to allow issuing certificates with no valid principals (meaning any valid principal). Exists for backwards compatibility only, the default of false is highly recommended.',
      type: 'boolean',
    },
    allowHostCertificates: {
      editType: 'boolean',
      helpText:
        "[Not applicable for OTP type] [Optional for CA type] If set, certificates are allowed to be signed for use as a 'host'.",
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowSubdomains: {
      editType: 'boolean',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If set, host certificates that are requested are allowed to use subdomains of those listed in "allowed_domains".',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowUserCertificates: {
      editType: 'boolean',
      helpText:
        "[Not applicable for OTP type] [Optional for CA type] If set, certificates are allowed to be signed for use as a 'user'.",
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowUserKeyIds: {
      editType: 'boolean',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If true, users can override the key ID for a signed certificate with the "key_id" field. When false, the key ID will always be the token display name. The key ID is logged by the SSH server and can be useful for auditing.',
      fieldGroup: 'default',
      label: 'Allow User Key IDs',
      type: 'boolean',
    },
    allowedCriticalOptions: {
      editType: 'string',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] A comma-separated list of critical options that certificates can have when signed. To allow any critical options, set this to an empty string.',
      fieldGroup: 'default',
      type: 'string',
    },
    allowedDomains: {
      editType: 'string',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If this option is not specified, client can request for a signed certificate for any valid host. If only certain domains are allowed, then this list enforces it.',
      fieldGroup: 'default',
      type: 'string',
    },
    allowedDomainsTemplate: {
      editType: 'boolean',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If set, Allowed domains can be specified using identity template policies. Non-templated domains are also permitted.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowedExtensions: {
      editType: 'string',
      helpText:
        "[Not applicable for OTP type] [Optional for CA type] A comma-separated list of extensions that certificates can have when signed. An empty list means that no extension overrides are allowed by an end-user; explicitly specify '*' to allow any extensions to be set.",
      fieldGroup: 'default',
      type: 'string',
    },
    allowedUserKeyLengths: {
      editType: 'object',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If set, allows the enforcement of key types and minimum key sizes to be signed.',
      fieldGroup: 'default',
      type: 'object',
    },
    allowedUsers: {
      editType: 'string',
      helpText:
        "[Optional for all types] [Works differently for CA type] If this option is not specified, or is '*', client can request a credential for any valid user at the remote host, including the admin user. If only certain usernames are to be allowed, then this list enforces it. If this field is set, then credentials can only be created for default_user and usernames present in this list. Setting this option will enable all the users with access to this role to fetch credentials for all other usernames in this list. Use with caution. N.B.: with the CA type, an empty list means that no users are allowed; explicitly specify '*' to allow any user.",
      fieldGroup: 'default',
      type: 'string',
    },
    allowedUsersTemplate: {
      editType: 'boolean',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If set, Allowed users can be specified using identity template policies. Non-templated users are also permitted.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    cidrList: {
      editType: 'string',
      helpText:
        '[Optional for OTP type] [Not applicable for CA type] Comma separated list of CIDR blocks for which the role is applicable for. CIDR blocks can belong to more than one role.',
      fieldGroup: 'default',
      label: 'CIDR List',
      type: 'string',
    },
    defaultCriticalOptions: {
      editType: 'object',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] Critical options certificates should have if none are provided when signing. This field takes in key value pairs in JSON format. Note that these are not restricted by "allowed_critical_options". Defaults to none.',
      fieldGroup: 'default',
      type: 'object',
    },
    defaultExtensions: {
      editType: 'object',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] Extensions certificates should have if none are provided when signing. This field takes in key value pairs in JSON format. Note that these are not restricted by "allowed_extensions". Defaults to none.',
      fieldGroup: 'default',
      type: 'object',
    },
    defaultExtensionsTemplate: {
      editType: 'boolean',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If set, Default extension values can be specified using identity template policies. Non-templated extension values are also permitted.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    defaultUser: {
      editType: 'string',
      helpText:
        "[Required for OTP type] [Optional for CA type] Default username for which a credential will be generated. When the endpoint 'creds/' is used without a username, this value will be used as default username.",
      fieldGroup: 'default',
      label: 'Default Username',
      type: 'string',
    },
    defaultUserTemplate: {
      editType: 'boolean',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] If set, Default user can be specified using identity template policies. Non-templated users are also permitted.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    excludeCidrList: {
      editType: 'string',
      helpText:
        '[Optional for OTP type] [Not applicable for CA type] Comma separated list of CIDR blocks. IP addresses belonging to these blocks are not accepted by the role. This is particularly useful when big CIDR blocks are being used by the role and certain parts of it needs to be kept out.',
      fieldGroup: 'default',
      label: 'Exclude CIDR List',
      type: 'string',
    },
    keyIdFormat: {
      editType: 'string',
      helpText:
        "[Not applicable for OTP type] [Optional for CA type] When supplied, this value specifies a custom format for the key id of a signed certificate. The following variables are available for use: '{{token_display_name}}' - The display name of the token used to make the request. '{{role_name}}' - The name of the role signing the request. '{{public_key_hash}}' - A SHA256 checksum of the public key that is being signed.",
      fieldGroup: 'default',
      label: 'Key ID Format',
      type: 'string',
    },
    keyType: {
      editType: 'string',
      helpText:
        "[Required for all types] Type of key used to login to hosts. It can be either 'otp' or 'ca'. 'otp' type requires agent to be installed in remote hosts.",
      possibleValues: ['otp', 'ca'],
      fieldGroup: 'default',
      defaultValue: 'ca',
      type: 'string',
    },
    maxTtl: {
      editType: 'ttl',
      helpText: '[Not applicable for OTP type] [Optional for CA type] The maximum allowed lease duration',
      fieldGroup: 'default',
      label: 'Max TTL',
    },
    notBeforeDuration: {
      editType: 'ttl',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] The duration that the SSH certificate should be backdated by at issuance.',
      fieldGroup: 'default',
      defaultValue: 30,
      label: 'Not before duration',
    },
    port: {
      editType: 'number',
      helpText:
        "[Optional for OTP type] [Not applicable for CA type] Port number for SSH connection. Default is '22'. Port number does not play any role in creation of OTP. For 'otp' type, this is just a way to inform client about the port number to use. Port number will be returned to client by Vault server along with OTP.",
      fieldGroup: 'default',
      defaultValue: 22,
      type: 'number',
    },
    ttl: {
      editType: 'ttl',
      helpText:
        '[Not applicable for OTP type] [Optional for CA type] The lease duration if no specific lease duration is requested. The lease duration controls the expiration of certificates issued by this backend. Defaults to the value of max_ttl.',
      fieldGroup: 'default',
      label: 'TTL',
    },
  },
};

const kmip = {
  'kmip/config': {
    defaultTlsClientKeyBits: {
      editType: 'number',
      helpText: 'Client certificate key bits, valid values depend on key type',
      fieldGroup: 'default',
      defaultValue: 256,
      label: 'Default TLS Client Key bits',
      type: 'number',
    },
    defaultTlsClientKeyType: {
      editType: 'string',
      helpText: 'Client certificate key type, rsa or ec',
      possibleValues: ['rsa', 'ec'],
      fieldGroup: 'default',
      defaultValue: 'ec',
      label: 'Default TLS Client Key type',
      type: 'string',
    },
    defaultTlsClientTtl: {
      editType: 'ttl',
      helpText:
        'Client certificate TTL in either an integer number of seconds (3600) or an integer time unit (1h)',
      fieldGroup: 'default',
      defaultValue: '336h',
      label: 'Default TLS Client TTL',
    },
    listenAddrs: {
      editType: 'stringArray',
      helpText:
        'A list of address:port to listen on. A bare address without port may be provided, in which case port 5696 is assumed.',
      fieldGroup: 'default',
      defaultValue: '127.0.0.1:5696',
    },
    serverHostnames: {
      editType: 'stringArray',
      helpText:
        "A list of hostnames to include in the server's TLS certificate as SAN DNS names. The first will be used as the common name (CN).",
      fieldGroup: 'default',
    },
    serverIps: {
      editType: 'stringArray',
      helpText: "A list of IP to include in the server's TLS certificate as SAN IP addresses.",
      fieldGroup: 'default',
    },
    tlsCaKeyBits: {
      editType: 'number',
      helpText: 'CA key bits, valid values depend on key type',
      fieldGroup: 'default',
      defaultValue: 256,
      label: 'TLS CA Key bits',
      type: 'number',
    },
    tlsCaKeyType: {
      editType: 'string',
      helpText: 'CA key type, rsa or ec',
      possibleValues: ['rsa', 'ec'],
      fieldGroup: 'default',
      defaultValue: 'ec',
      label: 'TLS CA Key type',
      type: 'string',
    },
    tlsMinVersion: {
      editType: 'string',
      helpText: 'Min TLS version',
      fieldGroup: 'default',
      defaultValue: 'tls12',
      label: 'Minimum TLS Version',
      type: 'string',
    },
  },
  'kmip/role': {
    role: {
      editType: 'string',
      helpText: 'Name of the role.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Role',
      type: 'string',
    },
    operationActivate: {
      editType: 'boolean',
      helpText: 'Allow the "Activate" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Activate',
      type: 'boolean',
    },
    operationAddAttribute: {
      editType: 'boolean',
      helpText: 'Allow the "Add Attribute" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Add Attribute',
      type: 'boolean',
    },
    operationAll: {
      editType: 'boolean',
      helpText:
        'Allow ALL operations to be performed by this role. This can be overridden if other allowed operations are set to false within the same request.',
      fieldGroup: 'default',
      label: 'All',
      type: 'boolean',
    },
    operationCreate: {
      editType: 'boolean',
      helpText: 'Allow the "Create" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Create',
      type: 'boolean',
    },
    operationCreateKeyPair: {
      editType: 'boolean',
      helpText: 'Allow the "Create Key Pair" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Create Key Pair',
      type: 'boolean',
    },
    operationDecrypt: {
      editType: 'boolean',
      helpText: 'Allow the "Decrypt" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Decrypt',
      type: 'boolean',
    },
    operationDeleteAttribute: {
      editType: 'boolean',
      helpText: 'Allow the "Delete Attribute" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Delete Attribute',
      type: 'boolean',
    },
    operationDestroy: {
      editType: 'boolean',
      helpText: 'Allow the "Destroy" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Destroy',
      type: 'boolean',
    },
    operationDiscoverVersions: {
      editType: 'boolean',
      helpText: 'Allow the "Discover Versions" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Discover Versions',
      type: 'boolean',
    },
    operationEncrypt: {
      editType: 'boolean',
      helpText: 'Allow the "Encrypt" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Encrypt',
      type: 'boolean',
    },
    operationGet: {
      editType: 'boolean',
      helpText: 'Allow the "Get" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Get',
      type: 'boolean',
    },
    operationGetAttributeList: {
      editType: 'boolean',
      helpText: 'Allow the "Get Attribute List" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Get Attribute List',
      type: 'boolean',
    },
    operationGetAttributes: {
      editType: 'boolean',
      helpText: 'Allow the "Get Attributes" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Get Attributes',
      type: 'boolean',
    },
    operationImport: {
      editType: 'boolean',
      helpText: 'Allow the "Import" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Import',
      type: 'boolean',
    },
    operationLocate: {
      editType: 'boolean',
      helpText: 'Allow the "Locate" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Locate',
      type: 'boolean',
    },
    operationMac: {
      editType: 'boolean',
      helpText: 'Allow the "Mac" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Mac',
      type: 'boolean',
    },
    operationMacVerify: {
      editType: 'boolean',
      helpText: 'Allow the "Mac Verify" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Mac Verify',
      type: 'boolean',
    },
    operationModifyAttribute: {
      editType: 'boolean',
      helpText: 'Allow the "Modify Attribute" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Modify Attribute',
      type: 'boolean',
    },
    operationNone: {
      editType: 'boolean',
      helpText:
        'Allow NO operations to be performed by this role. This can be overridden if other allowed operations are set to true within the same request.',
      fieldGroup: 'default',
      label: 'None',
      type: 'boolean',
    },
    operationQuery: {
      editType: 'boolean',
      helpText: 'Allow the "Query" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Query',
      type: 'boolean',
    },
    operationRegister: {
      editType: 'boolean',
      helpText: 'Allow the "Register" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Register',
      type: 'boolean',
    },
    operationRekey: {
      editType: 'boolean',
      helpText: 'Allow the "Rekey" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Rekey',
      type: 'boolean',
    },
    operationRekeyKeyPair: {
      editType: 'boolean',
      helpText: 'Allow the "Rekey Key Pair" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Rekey Key Pair',
      type: 'boolean',
    },
    operationRevoke: {
      editType: 'boolean',
      helpText: 'Allow the "Revoke" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Revoke',
      type: 'boolean',
    },
    operationRngRetrieve: {
      editType: 'boolean',
      helpText: 'Allow the "Rng Retrieve" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Rng Retrieve',
      type: 'boolean',
    },
    operationRngSeed: {
      editType: 'boolean',
      helpText: 'Allow the "Rng Seed" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Rng Seed',
      type: 'boolean',
    },
    operationSign: {
      editType: 'boolean',
      helpText: 'Allow the "Sign" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Sign',
      type: 'boolean',
    },
    operationSignatureVerify: {
      editType: 'boolean',
      helpText: 'Allow the "Signature Verify" operation to be performed by this role',
      fieldGroup: 'default',
      label: 'Signature Verify',
      type: 'boolean',
    },
    tlsClientKeyBits: {
      editType: 'number',
      helpText: 'Client certificate key bits, valid values depend on key type',
      fieldGroup: 'default',
      defaultValue: 521,
      label: 'TLS Client Key bits',
      type: 'number',
    },
    tlsClientKeyType: {
      editType: 'string',
      helpText: 'Client certificate key type, rsa or ec',
      possibleValues: ['rsa', 'ec'],
      fieldGroup: 'default',
      defaultValue: 'ec',
      label: 'TLS Client Key type',
      type: 'string',
    },
    tlsClientTtl: {
      editType: 'ttl',
      helpText:
        'Client certificate TTL in either an integer number of seconds (10) or an integer time unit (10s)',
      fieldGroup: 'default',
      defaultValue: '86400',
      label: 'TLS Client TTL',
    },
  },
};

const pki = {
  'pki/config/acme': {
    allowRoleExtKeyUsage: {
      editType: 'boolean',
      helpText:
        'whether the ExtKeyUsage field from a role is used, defaults to false meaning that certificate will be signed with ServerAuth.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowedIssuers: {
      editType: 'stringArray',
      helpText:
        'which issuers are allowed for use with ACME; by default, this will only be the primary (default) issuer',
      fieldGroup: 'default',
    },
    allowedRoles: {
      editType: 'stringArray',
      helpText:
        "which roles are allowed for use with ACME; by default via '*', these will be all roles including sign-verbatim; when concrete role names are specified, any default_directory_policy role must be included to allow usage of the default acme directories under /pki/acme/directory and /pki/issuer/:issuer_id/acme/directory.",
      fieldGroup: 'default',
    },
    defaultDirectoryPolicy: {
      editType: 'string',
      helpText:
        'the policy to be used for non-role-qualified ACME requests; by default ACME issuance will be otherwise unrestricted, equivalent to the sign-verbatim endpoint; one may also specify a role to use as this policy, as "role:<role_name>", the specified role must be allowed by allowed_roles',
      fieldGroup: 'default',
      type: 'string',
    },
    dnsResolver: {
      editType: 'string',
      helpText:
        'DNS resolver to use for domain resolution on this mount. Defaults to using the default system resolver. Must be in the format <host>:<port>, with both parts mandatory.',
      fieldGroup: 'default',
      type: 'string',
    },
    eabPolicy: {
      editType: 'string',
      helpText:
        "Specify the policy to use for external account binding behaviour, 'not-required', 'new-account-required' or 'always-required'",
      fieldGroup: 'default',
      type: 'string',
    },
    enabled: {
      editType: 'boolean',
      helpText:
        'whether ACME is enabled, defaults to false meaning that clusters will by default not get ACME support',
      fieldGroup: 'default',
      type: 'boolean',
    },
    maxTtl: {
      editType: 'ttl',
      fieldGroup: 'default',
      helpText:
        'Specify the maximum TTL for ACME certificates. Role TTL values will be limited to this value',
    },
  },
  'pki/certificate/generate': {
    role: {
      editType: 'string',
      helpText: 'The desired role with configuration for this request',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Role',
      type: 'string',
    },
    altNames: {
      editType: 'string',
      helpText:
        'The requested Subject Alternative Names, if any, in a comma-delimited list. If email protection is enabled for the role, this may contain email addresses.',
      fieldGroup: 'default',
      label: 'DNS/Email Subject Alternative Names (SANs)',
      type: 'string',
    },
    certMetadata: {
      editType: 'string',
      fieldGroup: 'default',
      helpText:
        "User supplied metadata to store associated with this certificate's serial number, base64 encoded",
      label: 'Certificate Metadata',
      type: 'string',
    },
    commonName: {
      editType: 'string',
      helpText:
        'The requested common name; if you want more than one, specify the alternative names in the alt_names map. If email protection is enabled in the role, this may be an email address.',
      fieldGroup: 'default',
      type: 'string',
    },
    excludeCnFromSans: {
      editType: 'boolean',
      helpText:
        'If true, the Common Name will not be included in DNS or Email Subject Alternate Names. Defaults to false (CN is included).',
      fieldGroup: 'default',
      label: 'Exclude Common Name from Subject Alternative Names (SANs)',
      type: 'boolean',
    },
    format: {
      editType: 'string',
      helpText:
        'Format for returned data. Can be "pem", "der", or "pem_bundle". If "pem_bundle", any private key and issuing cert will be appended to the certificate pem. If "der", the value will be base64 encoded. Defaults to "pem".',
      possibleValues: ['pem', 'der', 'pem_bundle'],
      fieldGroup: 'default',
      defaultValue: 'pem',
      type: 'string',
    },
    ipSans: {
      editType: 'stringArray',
      helpText: 'The requested IP SANs, if any, in a comma-delimited list',
      fieldGroup: 'default',
      label: 'IP Subject Alternative Names (SANs)',
    },
    issuerRef: {
      editType: 'string',
      helpText:
        'Reference to a existing issuer; either "default" for the configured default issuer, an identifier or the name assigned to the issuer.',
      fieldGroup: 'default',
      type: 'string',
    },
    notAfter: {
      editType: 'string',
      helpText:
        'Set the not after field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ',
      fieldGroup: 'default',
      type: 'string',
    },
    otherSans: {
      editType: 'stringArray',
      helpText:
        'Requested other SANs, in an array with the format <oid>;UTF8:<utf8 string value> for each entry.',
      fieldGroup: 'default',
      label: 'Other SANs',
    },
    privateKeyFormat: {
      editType: 'string',
      helpText:
        'Format for the returned private key. Generally the default will be controlled by the "format" parameter as either base64-encoded DER or PEM-encoded DER. However, this can be set to "pkcs8" to have the returned private key contain base64-encoded pkcs8 or PEM-encoded pkcs8 instead. Defaults to "der".',
      possibleValues: ['', 'der', 'pem', 'pkcs8'],
      fieldGroup: 'default',
      defaultValue: 'der',
      type: 'string',
    },
    removeRootsFromChain: {
      editType: 'boolean',
      helpText: 'Whether or not to remove self-signed CA certificates in the output of the ca_chain field.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    serialNumber: {
      editType: 'string',
      helpText:
        "The Subject's requested serial number, if any. See RFC 4519 Section 2.31 'serialNumber' for a description of this field. If you want more than one, specify alternative names in the alt_names map using OID 2.5.4.5. This has no impact on the final certificate's Serial Number field.",
      fieldGroup: 'default',
      type: 'string',
    },
    ttl: {
      editType: 'ttl',
      helpText:
        'The requested Time To Live for the certificate; sets the expiration date. If not specified the role default, backend default, or system default TTL is used, in that order. Cannot be larger than the role max TTL.',
      fieldGroup: 'default',
      label: 'TTL',
    },
    uriSans: {
      editType: 'stringArray',
      helpText: 'The requested URI SANs, if any, in a comma-delimited list.',
      fieldGroup: 'default',
      label: 'URI Subject Alternative Names (SANs)',
    },
    userIds: {
      editType: 'stringArray',
      helpText:
        'The requested user_ids value to place in the subject, if any, in a comma-delimited list. Restricted by allowed_user_ids. Any values are added with OID 0.9.2342.19200300.100.1.1.',
      fieldGroup: 'default',
      label: 'User ID(s)',
    },
  },
  'pki/certificate/sign': {
    role: {
      editType: 'string',
      helpText: 'The desired role with configuration for this request',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Role',
      type: 'string',
    },
    altNames: {
      editType: 'string',
      helpText:
        'The requested Subject Alternative Names, if any, in a comma-delimited list. If email protection is enabled for the role, this may contain email addresses.',
      fieldGroup: 'default',
      label: 'DNS/Email Subject Alternative Names (SANs)',
      type: 'string',
    },
    certMetadata: {
      editType: 'string',
      fieldGroup: 'default',
      helpText:
        "User supplied metadata to store associated with this certificate's serial number, base64 encoded",
      label: 'Certificate Metadata',
      type: 'string',
    },
    commonName: {
      editType: 'string',
      helpText:
        'The requested common name; if you want more than one, specify the alternative names in the alt_names map. If email protection is enabled in the role, this may be an email address.',
      fieldGroup: 'default',
      type: 'string',
    },
    csr: {
      editType: 'string',
      helpText: 'PEM-format CSR to be signed.',
      fieldGroup: 'default',
      type: 'string',
    },
    excludeCnFromSans: {
      editType: 'boolean',
      helpText:
        'If true, the Common Name will not be included in DNS or Email Subject Alternate Names. Defaults to false (CN is included).',
      fieldGroup: 'default',
      label: 'Exclude Common Name from Subject Alternative Names (SANs)',
      type: 'boolean',
    },
    format: {
      editType: 'string',
      helpText:
        'Format for returned data. Can be "pem", "der", or "pem_bundle". If "pem_bundle", any private key and issuing cert will be appended to the certificate pem. If "der", the value will be base64 encoded. Defaults to "pem".',
      possibleValues: ['pem', 'der', 'pem_bundle'],
      fieldGroup: 'default',
      defaultValue: 'pem',
      type: 'string',
    },
    ipSans: {
      editType: 'stringArray',
      helpText: 'The requested IP SANs, if any, in a comma-delimited list',
      fieldGroup: 'default',
      label: 'IP Subject Alternative Names (SANs)',
    },
    issuerRef: {
      editType: 'string',
      helpText:
        'Reference to a existing issuer; either "default" for the configured default issuer, an identifier or the name assigned to the issuer.',
      fieldGroup: 'default',
      type: 'string',
    },
    notAfter: {
      editType: 'string',
      helpText:
        'Set the not after field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ',
      fieldGroup: 'default',
      type: 'string',
    },
    otherSans: {
      editType: 'stringArray',
      helpText:
        'Requested other SANs, in an array with the format <oid>;UTF8:<utf8 string value> for each entry.',
      fieldGroup: 'default',
      label: 'Other SANs',
    },
    privateKeyFormat: {
      editType: 'string',
      helpText:
        'Format for the returned private key. Generally the default will be controlled by the "format" parameter as either base64-encoded DER or PEM-encoded DER. However, this can be set to "pkcs8" to have the returned private key contain base64-encoded pkcs8 or PEM-encoded pkcs8 instead. Defaults to "der".',
      possibleValues: ['', 'der', 'pem', 'pkcs8'],
      fieldGroup: 'default',
      defaultValue: 'der',
      type: 'string',
    },
    removeRootsFromChain: {
      editType: 'boolean',
      helpText: 'Whether or not to remove self-signed CA certificates in the output of the ca_chain field.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    serialNumber: {
      editType: 'string',
      helpText:
        "The Subject's requested serial number, if any. See RFC 4519 Section 2.31 'serialNumber' for a description of this field. If you want more than one, specify alternative names in the alt_names map using OID 2.5.4.5. This has no impact on the final certificate's Serial Number field.",
      fieldGroup: 'default',
      type: 'string',
    },
    ttl: {
      editType: 'ttl',
      helpText:
        'The requested Time To Live for the certificate; sets the expiration date. If not specified the role default, backend default, or system default TTL is used, in that order. Cannot be larger than the role max TTL.',
      fieldGroup: 'default',
      label: 'TTL',
    },
    uriSans: {
      editType: 'stringArray',
      helpText: 'The requested URI SANs, if any, in a comma-delimited list.',
      fieldGroup: 'default',
      label: 'URI Subject Alternative Names (SANs)',
    },
    userIds: {
      editType: 'stringArray',
      helpText:
        'The requested user_ids value to place in the subject, if any, in a comma-delimited list. Restricted by allowed_user_ids. Any values are added with OID 0.9.2342.19200300.100.1.1.',
      fieldGroup: 'default',
      label: 'User ID(s)',
    },
  },
  'pki/config/cluster': {
    aiaPath: {
      editType: 'string',
      helpText:
        "Optional URI to this mount's AIA distribution point; may refer to an external non-Vault responder. This is for resolving AIA URLs and providing the {{cluster_aia_path}} template parameter and will not be used for other purposes. As such, unlike path above, this could safely be an insecure transit mechanism (like HTTP without TLS). For example: http://cdn.example.com/pr1/pki",
      fieldGroup: 'default',
      type: 'string',
    },
    path: {
      editType: 'string',
      helpText:
        "Canonical URI to this mount on this performance replication cluster's external address. This is for resolving AIA URLs and providing the {{cluster_path}} template parameter but might be used for other purposes in the future. This should only point back to this particular PR replica and should not ever point to another PR cluster. It may point to any node in the PR replica, including standby nodes, and need not always point to the active node. For example: https://pr1.vault.example.com:8200/v1/pki",
      fieldGroup: 'default',
      type: 'string',
    },
  },
  'pki/role': {
    name: {
      editType: 'string',
      helpText: 'Name of the role',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Name',
      type: 'string',
    },
    allowAnyName: {
      editType: 'boolean',
      helpText:
        'If set, clients can request certificates for any domain, regardless of allowed_domains restrictions. See the documentation for more information.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowBareDomains: {
      editType: 'boolean',
      helpText:
        'If set, clients can request certificates for the base domains themselves, e.g. "example.com" of domains listed in allowed_domains. This is a separate option as in some cases this can be considered a security threat. See the documentation for more information.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowGlobDomains: {
      editType: 'boolean',
      helpText:
        'If set, domains specified in allowed_domains can include shell-style glob patterns, e.g. "ftp*.example.com". See the documentation for more information.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowIpSans: {
      editType: 'boolean',
      helpText:
        'If set, IP Subject Alternative Names are allowed. Any valid IP is accepted and No authorization checking is performed.',
      fieldGroup: 'default',
      defaultValue: true,
      label: 'Allow IP Subject Alternative Names',
      type: 'boolean',
    },
    allowLocalhost: {
      editType: 'boolean',
      helpText:
        'Whether to allow "localhost" and "localdomain" as a valid common name in a request, independent of allowed_domains value.',
      fieldGroup: 'default',
      defaultValue: true,
      type: 'boolean',
    },
    allowSubdomains: {
      editType: 'boolean',
      helpText:
        'If set, clients can request certificates for subdomains of domains listed in allowed_domains, including wildcard subdomains. See the documentation for more information.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowWildcardCertificates: {
      editType: 'boolean',
      helpText:
        'If set, allows certificates with wildcards in the common name to be issued, conforming to RFC 6125\'s Section 6.4.3; e.g., "*.example.net" or "b*z.example.net". See the documentation for more information.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowedDomains: {
      editType: 'stringArray',
      helpText:
        'Specifies the domains this role is allowed to issue certificates for. This is used with the allow_bare_domains, allow_subdomains, and allow_glob_domains to determine matches for the common name, DNS-typed SAN entries, and Email-typed SAN entries of certificates. See the documentation for more information. This parameter accepts a comma-separated string or list of domains.',
      fieldGroup: 'default',
    },
    allowedDomainsTemplate: {
      editType: 'boolean',
      helpText:
        'If set, Allowed domains can be specified using identity template policies. Non-templated domains are also permitted.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowedOtherSans: {
      editType: 'stringArray',
      helpText:
        'If set, an array of allowed other names to put in SANs. These values support globbing and must be in the format <oid>;<type>:<value>. Currently only "utf8" is a valid type. All values, including globbing values, must use this syntax, with the exception being a single "*" which allows any OID and any value (but type must still be utf8).',
      fieldGroup: 'default',
      label: 'Allowed Other Subject Alternative Names',
    },
    allowedSerialNumbers: {
      editType: 'stringArray',
      helpText:
        'If set, an array of allowed serial numbers to put in Subject. These values support globbing.',
      fieldGroup: 'default',
    },
    allowedUriSans: {
      editType: 'stringArray',
      helpText:
        'If set, an array of allowed URIs for URI Subject Alternative Names. Any valid URI is accepted, these values support globbing.',
      fieldGroup: 'default',
      label: 'Allowed URI Subject Alternative Names',
    },
    allowedUriSansTemplate: {
      editType: 'boolean',
      helpText:
        'If set, Allowed URI SANs can be specified using identity template policies. Non-templated URI SANs are also permitted.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    allowedUserIds: {
      editType: 'stringArray',
      helpText:
        'If set, an array of allowed user-ids to put in user system login name specified here: https://www.rfc-editor.org/rfc/rfc1274#section-9.3.1',
      fieldGroup: 'default',
    },
    backend: {
      editType: 'string',
      helpText: 'Backend Type',
      fieldGroup: 'default',
      type: 'string',
    },
    basicConstraintsValidForNonCa: {
      editType: 'boolean',
      helpText: 'Mark Basic Constraints valid when issuing non-CA certificates.',
      fieldGroup: 'default',
      label: 'Basic Constraints Valid for Non-CA',
      type: 'boolean',
    },
    clientFlag: {
      editType: 'boolean',
      helpText:
        'If set, certificates are flagged for client auth use. Defaults to true. See also RFC 5280 Section 4.2.1.12.',
      fieldGroup: 'default',
      defaultValue: true,
      type: 'boolean',
    },
    cnValidations: {
      editType: 'stringArray',
      helpText:
        "List of allowed validations to run against the Common Name field. Values can include 'email' to validate the CN is a email address, 'hostname' to validate the CN is a valid hostname (potentially including wildcards). When multiple validations are specified, these take OR semantics (either email OR hostname are allowed). The special value 'disabled' allows disabling all CN name validations, allowing for arbitrary non-Hostname, non-Email address CNs.",
      fieldGroup: 'default',
      label: 'Common Name Validations',
    },
    codeSigningFlag: {
      editType: 'boolean',
      helpText:
        'If set, certificates are flagged for code signing use. Defaults to false. See also RFC 5280 Section 4.2.1.12.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    country: {
      editType: 'stringArray',
      helpText: 'If set, Country will be set to this value in certificates issued by this role.',
      fieldGroup: 'default',
    },
    emailProtectionFlag: {
      editType: 'boolean',
      helpText:
        'If set, certificates are flagged for email protection use. Defaults to false. See also RFC 5280 Section 4.2.1.12.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    enforceHostnames: {
      editType: 'boolean',
      helpText:
        'If set, only valid host names are allowed for CN and DNS SANs, and the host part of email addresses. Defaults to true.',
      fieldGroup: 'default',
      defaultValue: true,
      type: 'boolean',
    },
    extKeyUsage: {
      editType: 'stringArray',
      helpText:
        'A comma-separated string or list of extended key usages. Valid values can be found at https://golang.org/pkg/crypto/x509/#ExtKeyUsage -- simply drop the "ExtKeyUsage" part of the name. To remove all key usages from being set, set this value to an empty list. See also RFC 5280 Section 4.2.1.12.',
      fieldGroup: 'default',
      label: 'Extended Key Usage',
    },
    extKeyUsageOids: {
      editType: 'stringArray',
      helpText: 'A comma-separated string or list of extended key usage oids.',
      fieldGroup: 'default',
      label: 'Extended Key Usage OIDs',
    },
    generateLease: {
      editType: 'boolean',
      helpText:
        'If set, certificates issued/signed against this role will have Vault leases attached to them. Defaults to "false". Certificates can be added to the CRL by "vault revoke <lease_id>" when certificates are associated with leases. It can also be done using the "pki/revoke" endpoint. However, when lease generation is disabled, invoking "pki/revoke" would be the only way to add the certificates to the CRL. When large number of certificates are generated with long lifetimes, it is recommended that lease generation be disabled, as large amount of leases adversely affect the startup time of Vault.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    issuerRef: {
      editType: 'string',
      helpText: 'Reference to the issuer used to sign requests serviced by this role.',
      fieldGroup: 'default',
      type: 'string',
    },
    keyBits: {
      editType: 'number',
      helpText:
        'The number of bits to use. Allowed values are 0 (universal default); with rsa key_type: 2048 (default), 3072, or 4096; with ec key_type: 224, 256 (default), 384, or 521; ignored with ed25519.',
      fieldGroup: 'default',
      type: 'number',
    },
    keyType: {
      editType: 'string',
      helpText:
        'The type of key to use; defaults to RSA. "rsa" "ec", "ed25519" and "any" are the only valid values.',
      possibleValues: ['rsa', 'ec', 'ed25519', 'any'],
      fieldGroup: 'default',
      type: 'string',
    },
    keyUsage: {
      editType: 'stringArray',
      helpText:
        'A comma-separated string or list of key usages (not extended key usages). Valid values can be found at https://golang.org/pkg/crypto/x509/#KeyUsage -- simply drop the "KeyUsage" part of the name. To remove all key usages from being set, set this value to an empty list. See also RFC 5280 Section 4.2.1.3.',
      fieldGroup: 'default',
      defaultValue: 'DigitalSignature,KeyAgreement,KeyEncipherment',
    },
    locality: {
      editType: 'stringArray',
      helpText: 'If set, Locality will be set to this value in certificates issued by this role.',
      fieldGroup: 'default',
      label: 'Locality/City',
    },
    maxTtl: {
      editType: 'ttl',
      helpText: 'The maximum allowed lease duration. If not set, defaults to the system maximum lease TTL.',
      fieldGroup: 'default',
      label: 'Max TTL',
    },
    noStore: {
      editType: 'boolean',
      helpText:
        'If set, certificates issued/signed against this role will not be stored in the storage backend. This can improve performance when issuing large numbers of certificates. However, certificates issued in this way cannot be enumerated or revoked, so this option is recommended only for certificates that are non-sensitive, or extremely short-lived. This option implies a value of "false" for "generate_lease".',
      fieldGroup: 'default',
      type: 'boolean',
    },
    noStoreMetadata: {
      editType: 'boolean',
      helpText:
        'If set, if a client attempts to issue or sign a certificate with attached cert_metadata to store, the issuance / signing instead fails.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    notAfter: {
      editType: 'string',
      helpText:
        'Set the not after field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ.',
      fieldGroup: 'default',
      type: 'string',
    },
    notBeforeDuration: {
      editType: 'ttl',
      helpText: 'The duration before now which the certificate needs to be backdated by.',
      fieldGroup: 'default',
      defaultValue: 30,
    },
    organization: {
      editType: 'stringArray',
      helpText: 'If set, O (Organization) will be set to this value in certificates issued by this role.',
      fieldGroup: 'default',
    },
    ou: {
      editType: 'stringArray',
      helpText:
        'If set, OU (OrganizationalUnit) will be set to this value in certificates issued by this role.',
      fieldGroup: 'default',
      label: 'Organizational Unit',
    },
    policyIdentifiers: {
      editType: 'stringArray',
      helpText:
        'A comma-separated string or list of policy OIDs, or a JSON list of qualified policy information, which must include an oid, and may include a notice and/or cps url, using the form [{"oid"="1.3.6.1.4.1.7.8","notice"="I am a user Notice"}, {"oid"="1.3.6.1.4.1.44947.1.2.4 ","cps"="https://example.com"}].',
      fieldGroup: 'default',
    },
    postalCode: {
      editType: 'stringArray',
      helpText: 'If set, Postal Code will be set to this value in certificates issued by this role.',
      fieldGroup: 'default',
    },
    province: {
      editType: 'stringArray',
      helpText: 'If set, Province will be set to this value in certificates issued by this role.',
      fieldGroup: 'default',
      label: 'Province/State',
    },
    requireCn: {
      editType: 'boolean',
      helpText: "If set to false, makes the 'common_name' field optional while generating a certificate.",
      fieldGroup: 'default',
      label: 'Require Common Name',
      type: 'boolean',
    },
    serialNumberSource: {
      defaultValue: 'json-csr',
      editType: 'string',
      fieldGroup: 'default',
      helpText:
        'Source for the certificate subject serial number. If "json-csr" (default), the value from the JSON serial_number field is used, falling back to the value in the CSR if empty. If "json", the value from the serial_number JSON field is used, ignoring the value in the CSR.',
      label: 'Serial number source',
      type: 'string',
    },
    serverFlag: {
      editType: 'boolean',
      helpText:
        'If set, certificates are flagged for server auth use. Defaults to true. See also RFC 5280 Section 4.2.1.12.',
      fieldGroup: 'default',
      defaultValue: true,
      type: 'boolean',
    },
    signatureBits: {
      editType: 'number',
      helpText:
        'The number of bits to use in the signature algorithm; accepts 256 for SHA-2-256, 384 for SHA-2-384, and 512 for SHA-2-512. Defaults to 0 to automatically detect based on key length (SHA-2-256 for RSA keys, and matching the curve size for NIST P-Curves).',
      fieldGroup: 'default',
      type: 'number',
    },
    streetAddress: {
      editType: 'stringArray',
      helpText: 'If set, Street Address will be set to this value in certificates issued by this role.',
      fieldGroup: 'default',
    },
    ttl: {
      editType: 'ttl',
      helpText:
        'The lease duration (validity period of the certificate) if no specific lease duration is requested. The lease duration controls the expiration of certificates issued by this backend. Defaults to the system default value or the value of max_ttl, whichever is shorter.',
      fieldGroup: 'default',
      label: 'TTL',
    },
    useCsrCommonName: {
      editType: 'boolean',
      helpText:
        'If set, when used with a signing profile, the common name in the CSR will be used. This does *not* include any requested Subject Alternative Names; use use_csr_sans for that. Defaults to true.',
      fieldGroup: 'default',
      defaultValue: true,
      label: 'Use CSR Common Name',
      type: 'boolean',
    },
    useCsrSans: {
      editType: 'boolean',
      helpText:
        'If set, when used with a signing profile, the SANs in the CSR will be used. This does *not* include the Common Name (cn); use use_csr_common_name for that. Defaults to true.',
      fieldGroup: 'default',
      defaultValue: true,
      label: 'Use CSR Subject Alternative Names',
      type: 'boolean',
    },
    usePss: {
      editType: 'boolean',
      helpText: 'Whether or not to use PSS signatures when using a RSA key-type issuer. Defaults to false.',
      fieldGroup: 'default',
      type: 'boolean',
    },
  },
  'pki/sign-intermediate': {
    issuerRef: {
      editType: 'string',
      helpText:
        'Reference to a existing issuer; either "default" for the configured default issuer, an identifier or the name assigned to the issuer.',
      fieldValue: 'mutableId',
      fieldGroup: 'default',
      readOnly: true,
      label: 'Issuer ref',
      type: 'string',
    },
    altNames: {
      editType: 'string',
      helpText:
        'The requested Subject Alternative Names, if any, in a comma-delimited list. May contain both DNS names and email addresses.',
      fieldGroup: 'default',
      label: 'DNS/Email Subject Alternative Names (SANs)',
      type: 'string',
    },
    commonName: {
      editType: 'string',
      helpText:
        'The requested common name; if you want more than one, specify the alternative names in the alt_names map. If not specified when signing, the common name will be taken from the CSR; other names must still be specified in alt_names or ip_sans.',
      fieldGroup: 'default',
      type: 'string',
    },
    country: {
      editType: 'stringArray',
      helpText: 'If set, Country will be set to this value.',
      fieldGroup: 'default',
    },
    csr: {
      editType: 'string',
      helpText: 'PEM-format CSR to be signed.',
      fieldGroup: 'default',
      type: 'string',
    },
    enforceLeafNotAfterBehavior: {
      editType: 'boolean',
      fieldGroup: 'default',
      helpText: "Do not truncate the NotAfter field, use the issuer's configured leaf_not_after_behavior",
      type: 'boolean',
    },
    excludeCnFromSans: {
      editType: 'boolean',
      helpText:
        'If true, the Common Name will not be included in DNS or Email Subject Alternate Names. Defaults to false (CN is included).',
      fieldGroup: 'default',
      label: 'Exclude Common Name from Subject Alternative Names (SANs)',
      type: 'boolean',
    },
    excludedDnsDomains: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'Domains for which this certificate is not allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).',
      label: 'Excluded DNS Domains',
    },
    excludedEmailAddresses: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'Email addresses for which this certificate is not allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).',
      label: 'Excluded email addresses',
    },
    excludedIpRanges: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'IP ranges for which this certificate is not allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10). Ranges must be specified in the notation of IP address and prefix length, like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.',
      label: 'Excluded IP ranges',
    },
    excludedUriDomains: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'URI domains for which this certificate is not allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).',
      label: 'Excluded URI domains',
    },
    format: {
      editType: 'string',
      helpText:
        'Format for returned data. Can be "pem", "der", or "pem_bundle". If "pem_bundle", any private key and issuing cert will be appended to the certificate pem. If "der", the value will be base64 encoded. Defaults to "pem".',
      possibleValues: ['pem', 'der', 'pem_bundle'],
      fieldGroup: 'default',
      defaultValue: 'pem',
      type: 'string',
    },
    ipSans: {
      editType: 'stringArray',
      helpText: 'The requested IP SANs, if any, in a comma-delimited list',
      fieldGroup: 'default',
      label: 'IP Subject Alternative Names (SANs)',
    },
    issuerName: {
      editType: 'string',
      helpText:
        "Provide a name to the generated or existing issuer, the name must be unique across all issuers and not be the reserved value 'default'",
      fieldGroup: 'default',
      type: 'string',
    },
    keyUsage: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'This list of key usages (not extended key usages) will be added to the existing set of key usages, CRL,CertSign, on the generated certificate. Valid values can be found at https://golang.org/pkg/crypto/x509/#KeyUsage -- simply drop the "KeyUsage" part of the name. To use the issuer for CMPv2, DigitalSignature must be set.',
    },
    locality: {
      editType: 'stringArray',
      helpText: 'If set, Locality will be set to this value.',
      fieldGroup: 'default',
      label: 'Locality/City',
    },
    maxPathLength: {
      editType: 'number',
      helpText: 'The maximum allowable path length',
      fieldGroup: 'default',
      type: 'number',
    },
    notAfter: {
      editType: 'string',
      helpText:
        'Set the not after field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ',
      fieldGroup: 'default',
      type: 'string',
    },
    notBeforeDuration: {
      editType: 'ttl',
      helpText: 'The duration before now which the certificate needs to be backdated by.',
      fieldGroup: 'default',
      defaultValue: 30,
    },
    organization: {
      editType: 'stringArray',
      helpText: 'If set, O (Organization) will be set to this value.',
      fieldGroup: 'default',
    },
    otherSans: {
      editType: 'stringArray',
      helpText:
        'Requested other SANs, in an array with the format <oid>;UTF8:<utf8 string value> for each entry.',
      fieldGroup: 'default',
      label: 'Other SANs',
    },
    ou: {
      editType: 'stringArray',
      helpText: 'If set, OU (OrganizationalUnit) will be set to this value.',
      fieldGroup: 'default',
      label: 'OU (Organizational Unit)',
    },
    permittedDnsDomains: {
      editType: 'stringArray',
      helpText:
        'Domains for which this certificate is allowed to sign or issue child certificates. If set, all DNS names (subject and alt) on child certs must be exact matches or subsets of the given domains (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).',
      fieldGroup: 'default',
      label: 'Permitted DNS Domains',
    },
    permittedEmailAddresses: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'Email addresses for which this certificate is allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).',
      label: 'Permitted email addresses',
    },
    permittedIpRanges: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'IP ranges for which this certificate is allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10). Ranges must be specified in the notation of IP address and prefix length, like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.',
      label: 'Permitted IP ranges',
    },
    permittedUriDomains: {
      editType: 'stringArray',
      fieldGroup: 'default',
      helpText:
        'URI domains for which this certificate is allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).',
      label: 'Permitted URI domains',
    },
    postalCode: {
      editType: 'stringArray',
      helpText: 'If set, Postal Code will be set to this value.',
      fieldGroup: 'default',
      label: 'Postal Code',
    },
    privateKeyFormat: {
      editType: 'string',
      helpText:
        'Format for the returned private key. Generally the default will be controlled by the "format" parameter as either base64-encoded DER or PEM-encoded DER. However, this can be set to "pkcs8" to have the returned private key contain base64-encoded pkcs8 or PEM-encoded pkcs8 instead. Defaults to "der".',
      possibleValues: ['', 'der', 'pem', 'pkcs8'],
      fieldGroup: 'default',
      defaultValue: 'der',
      type: 'string',
    },
    province: {
      editType: 'stringArray',
      helpText: 'If set, Province will be set to this value.',
      fieldGroup: 'default',
      label: 'Province/State',
    },
    serialNumber: {
      editType: 'string',
      helpText:
        "The Subject's requested serial number, if any. See RFC 4519 Section 2.31 'serialNumber' for a description of this field. If you want more than one, specify alternative names in the alt_names map using OID 2.5.4.5. This has no impact on the final certificate's Serial Number field.",
      fieldGroup: 'default',
      type: 'string',
    },
    signatureBits: {
      editType: 'number',
      helpText:
        'The number of bits to use in the signature algorithm; accepts 256 for SHA-2-256, 384 for SHA-2-384, and 512 for SHA-2-512. Defaults to 0 to automatically detect based on key length (SHA-2-256 for RSA keys, and matching the curve size for NIST P-Curves).',
      fieldGroup: 'default',
      type: 'number',
    },
    skid: {
      editType: 'string',
      helpText:
        "Value for the Subject Key Identifier field (RFC 5280 Section 4.2.1.2). This value should ONLY be used when cross-signing to mimic the existing certificate's SKID value; this is necessary to allow certain TLS implementations (such as OpenSSL) which use SKID/AKID matches in chain building to restrict possible valid chains. Specified as a string in hex format. Default is empty, allowing Vault to automatically calculate the SKID according to method one in the above RFC section.",
      fieldGroup: 'default',
      type: 'string',
    },
    streetAddress: {
      editType: 'stringArray',
      helpText: 'If set, Street Address will be set to this value.',
      fieldGroup: 'default',
      label: 'Street Address',
    },
    ttl: {
      editType: 'ttl',
      helpText:
        'The requested Time To Live for the certificate; sets the expiration date. If not specified the role default, backend default, or system default TTL is used, in that order. Cannot be larger than the mount max TTL. Note: this only has an effect when generating a CA cert or signing a CA cert, not when generating a CSR for an intermediate CA.',
      fieldGroup: 'default',
      label: 'TTL',
    },
    uriSans: {
      editType: 'stringArray',
      helpText: 'The requested URI SANs, if any, in a comma-delimited list.',
      fieldGroup: 'default',
      label: 'URI Subject Alternative Names (SANs)',
    },
    useCsrValues: {
      editType: 'boolean',
      helpText:
        'If true, then: 1) Subject information, including names and alternate names, will be preserved from the CSR rather than using values provided in the other parameters to this path; 2) Any key usages requested in the CSR will be added to the basic set of key usages used for CA certs signed by this path; for instance, the non-repudiation flag; 3) Extensions requested in the CSR will be copied into the issued certificate.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    usePss: {
      editType: 'boolean',
      helpText: 'Whether or not to use PSS signatures when using a RSA key-type issuer. Defaults to false.',
      fieldGroup: 'default',
      type: 'boolean',
    },
  },
  'pki/tidy': {
    acmeAccountSafetyBuffer: {
      editType: 'ttl',
      helpText:
        'The amount of time that must pass after creation that an account with no orders is marked revoked, and the amount of time after being marked revoked or deactivated.',
      fieldGroup: 'default',
    },
    enabled: {
      editType: 'boolean',
      helpText: 'Set to true to enable automatic tidy operations.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    intervalDuration: {
      editType: 'ttl',
      helpText:
        'Interval at which to run an auto-tidy operation. This is the time between tidy invocations (after one finishes to the start of the next). Running a manual tidy will reset this duration.',
      fieldGroup: 'default',
    },
    minStartupBackoffDuration: {
      editType: 'ttl',
      helpText: 'The minimum amount of time in seconds auto-tidy will be delayed after startup.',
      fieldGroup: 'default',
    },
    maxStartupBackoffDuration: {
      editType: 'ttl',
      helpText: 'The maximum amount of time in seconds auto-tidy will be delayed after startup.',
      fieldGroup: 'default',
    },
    issuerSafetyBuffer: {
      editType: 'ttl',
      helpText:
        "The amount of extra time that must have passed beyond issuer's expiration before it is removed from the backend storage. Defaults to 8760 hours (1 year).",
      fieldGroup: 'default',
    },
    maintainStoredCertificateCounts: {
      editType: 'boolean',
      helpText:
        'This configures whether stored certificates are counted upon initialization of the backend, and whether during normal operation, a running count of certificates stored is maintained.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    pauseDuration: {
      editType: 'string',
      helpText:
        'The amount of time to wait between processing certificates. This allows operators to change the execution profile of tidy to take consume less resources by slowing down how long it takes to run. Note that the entire list of certificates will be stored in memory during the entire tidy operation, but resources to read/process/update existing entries will be spread out over a greater period of time. By default this is zero seconds.',
      fieldGroup: 'default',
      type: 'string',
    },
    publishStoredCertificateCountMetrics: {
      editType: 'boolean',
      helpText:
        'This configures whether the stored certificate count is published to the metrics consumer. It does not affect if the stored certificate count is maintained, and if maintained, it will be available on the tidy-status endpoint.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    revocationQueueSafetyBuffer: {
      editType: 'ttl',
      helpText:
        'The amount of time that must pass from the cross-cluster revocation request being initiated to when it will be slated for removal. Setting this too low may remove valid revocation requests before the owning cluster has a chance to process them, especially if the cluster is offline.',
      fieldGroup: 'default',
    },
    safetyBuffer: {
      editType: 'ttl',
      helpText:
        'The amount of extra time that must have passed beyond certificate expiration before it is removed from the backend storage and/or revocation list. Defaults to 72 hours.',
      fieldGroup: 'default',
    },
    tidyAcme: {
      editType: 'boolean',
      helpText:
        'Set to true to enable tidying ACME accounts, orders and authorizations. ACME orders are tidied (deleted) safety_buffer after the certificate associated with them expires, or after the order and relevant authorizations have expired if no certificate was produced. Authorizations are tidied with the corresponding order. When a valid ACME Account is at least acme_account_safety_buffer old, and has no remaining orders associated with it, the account is marked as revoked. After another acme_account_safety_buffer has passed from the revocation or deactivation date, a revoked or deactivated ACME account is deleted.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    tidyCertMetadata: {
      editType: 'boolean',
      fieldGroup: 'default',
      helpText: 'Set to true to enable tidying up certificate metadata',
      type: 'boolean',
    },
    tidyCertStore: {
      editType: 'boolean',
      helpText: 'Set to true to enable tidying up the certificate store',
      fieldGroup: 'default',
      type: 'boolean',
    },
    tidyCmpv2NonceStore: {
      editType: 'boolean',
      fieldGroup: 'default',
      helpText: 'Set to true to enable tidying up the CMPv2 nonce store',
      type: 'boolean',
    },
    tidyCrossClusterRevokedCerts: {
      editType: 'boolean',
      helpText:
        'Set to true to enable tidying up the cross-cluster revoked certificate store. Only runs on the active primary node.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    tidyExpiredIssuers: {
      editType: 'boolean',
      helpText:
        'Set to true to automatically remove expired issuers past the issuer_safety_buffer. No keys will be removed as part of this operation.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    tidyMoveLegacyCaBundle: {
      editType: 'boolean',
      helpText:
        'Set to true to move the legacy ca_bundle from /config/ca_bundle to /config/ca_bundle.bak. This prevents downgrades to pre-Vault 1.11 versions (as older PKI engines do not know about the new multi-issuer storage layout), but improves the performance on seal wrapped PKI mounts. This will only occur if at least issuer_safety_buffer time has occurred after the initial storage migration. This backup is saved in case of an issue in future migrations. Operators may consider removing it via sys/raw if they desire. The backup will be removed via a DELETE /root call, but note that this removes ALL issuers within the mount (and is thus not desirable in most operational scenarios).',
      fieldGroup: 'default',
      type: 'boolean',
    },
    tidyRevocationList: {
      editType: 'boolean',
      helpText: "Deprecated; synonym for 'tidy_revoked_certs",
      fieldGroup: 'default',
      type: 'boolean',
    },
    tidyRevocationQueue: {
      editType: 'boolean',
      helpText:
        "Set to true to remove stale revocation queue entries that haven't been confirmed by any active cluster. Only runs on the active primary node",
      fieldGroup: 'default',
      type: 'boolean',
    },
    tidyRevokedCertIssuerAssociations: {
      editType: 'boolean',
      helpText:
        'Set to true to validate issuer associations on revocation entries. This helps increase the performance of CRL building and OCSP responses.',
      fieldGroup: 'default',
      type: 'boolean',
    },
    tidyRevokedCerts: {
      editType: 'boolean',
      helpText:
        'Set to true to expire all revoked and expired certificates, removing them both from the CRL and from storage. The CRL will be rotated if this causes any values to be removed.',
      fieldGroup: 'default',
      type: 'boolean',
    },
  },
  'pki/config/urls': {
    crlDistributionPoints: {
      editType: 'stringArray',
      helpText:
        'Comma-separated list of URLs to be used for the CRL distribution points attribute. See also RFC 5280 Section 4.2.1.13.',
      fieldGroup: 'default',
    },
    enableTemplating: {
      editType: 'boolean',
      helpText:
        "Whether or not to enabling templating of the above AIA fields. When templating is enabled the special values '{{issuer_id}}', '{{cluster_path}}', and '{{cluster_aia_path}}' are available, but the addresses are not checked for URI validity until issuance time. Using '{{cluster_path}}' requires /config/cluster's 'path' member to be set on all PR Secondary clusters and using '{{cluster_aia_path}}' requires /config/cluster's 'aia_path' member to be set on all PR secondary clusters.",
      fieldGroup: 'default',
      type: 'boolean',
    },
    issuingCertificates: {
      editType: 'stringArray',
      helpText:
        'Comma-separated list of URLs to be used for the issuing certificate attribute. See also RFC 5280 Section 4.2.2.1.',
      fieldGroup: 'default',
    },
    ocspServers: {
      editType: 'stringArray',
      helpText:
        'Comma-separated list of URLs to be used for the OCSP servers attribute. See also RFC 5280 Section 4.2.2.1.',
      fieldGroup: 'default',
    },
  },
};

// export object by backend name. keys of each object are model names
export default {
  kmip,
  ssh,
  pki,
};
