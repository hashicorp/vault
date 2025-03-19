/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * These are all the auth methods that can be mounted (enabled) in the UI.
 * Token is not included because it cannot be enabled or disabled.
 * The method data is all related to displaying each auth type.
 * The `supported-login-methods` util handles method data related to authenticating
 */

const ENTERPRISE_AUTH_BACKENDS = [
  {
    displayName: 'SAML',
    value: 'saml',
    type: 'saml',
    category: 'generic',
    glyph: 'saml-color',
  },
];

const BASE_AUTH_BACKENDS = [
  {
    displayName: 'AliCloud',
    value: 'alicloud',
    type: 'alicloud',
    category: 'cloud',
    glyph: 'alibaba-color',
  },
  {
    displayName: 'AppRole',
    value: 'approle',
    type: 'approle',
    category: 'generic',
    glyph: 'cpu',
  },
  {
    displayName: 'AWS',
    value: 'aws',
    type: 'aws',
    category: 'cloud',
    glyph: 'aws-color',
  },
  {
    displayName: 'Azure',
    value: 'azure',
    type: 'azure',
    category: 'cloud',
    glyph: 'azure-color',
  },
  {
    displayName: 'Google Cloud',
    value: 'gcp',
    type: 'gcp',
    category: 'cloud',
    glyph: 'gcp-color',
  },
  {
    displayName: 'GitHub',
    value: 'github',
    type: 'github',
    category: 'cloud',
    glyph: 'github-color',
  },
  {
    displayName: 'JWT',
    value: 'jwt',
    type: 'jwt',
    glyph: 'jwt',
    category: 'generic',
  },
  {
    displayName: 'OIDC',
    value: 'oidc',
    type: 'oidc',
    glyph: 'openid-color',
    category: 'generic',
  },
  {
    displayName: 'Kubernetes',
    value: 'kubernetes',
    type: 'kubernetes',
    category: 'infra',
    glyph: 'kubernetes-color',
  },
  {
    displayName: 'LDAP',
    value: 'ldap',
    type: 'ldap',
    glyph: 'folder-users',
    category: 'infra',
  },
  {
    displayName: 'Okta',
    value: 'okta',
    type: 'okta',
    category: 'infra',
    glyph: 'okta-color',
  },
  {
    displayName: 'RADIUS',
    value: 'radius',
    type: 'radius',
    glyph: 'mainframe',
    category: 'infra',
  },
  {
    displayName: 'TLS Certificates',
    value: 'cert',
    type: 'cert',
    category: 'generic',
    glyph: 'certificate',
  },
  {
    displayName: 'Username & Password',
    value: 'userpass',
    type: 'userpass',
    category: 'generic',
    glyph: 'users',
  },
];

const ALL_AUTH_BACKENDS = [...BASE_AUTH_BACKENDS, ...ENTERPRISE_AUTH_BACKENDS];

// The UI supports management of these auth methods (i.e. configuring roles or users)
// otherwise only configuration (enabling and mounting) of the method is supported.
export const MANAGED_AUTH_BACKENDS = ['cert', 'kubernetes', 'ldap', 'okta', 'radius', 'userpass'];

export const findAuthMethod = (authType: string) => ALL_AUTH_BACKENDS.find((m) => m.type === authType);

export const mountableAuthMethods = (isEnterprise: boolean) => {
  return isEnterprise ? ALL_AUTH_BACKENDS : BASE_AUTH_BACKENDS;
};
