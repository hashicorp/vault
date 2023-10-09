// The constants within this file represent the expected model attributes as parsed from OpenAPI
// if changes are made to the OpenAPI spec, that may result in changes that must be reflected
// here as well as ensured to not cause breaking changes within the UI.

/* Secret Engines */
const sshRole = {
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
};

/* Auth Engines */
const userpassUser = {
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
};

export default {
  sshRole,
  userpassUser,
};
