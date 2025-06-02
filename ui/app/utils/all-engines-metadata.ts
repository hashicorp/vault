/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * These are all the secret and auth methods, including enterprise.
 * If filtering for enterprise-only engines or a specific category is required, that logic should be done within the consumer of this utility.
 */

export interface EngineDisplayData {
  category?: string;
  displayName: string;
  glyph?: string;
  requiresEnterprise?: boolean;
  type: string;
  value?: string;
}

export const ALL_ENGINES: EngineDisplayData[] = [
  {
    category: 'cloud',
    displayName: 'AliCloud',
    glyph: 'alibaba-color',
    type: 'alicloud',
  },
  {
    category: 'generic',
    displayName: 'AppRole',
    glyph: 'cpu',
    type: 'approle',
    value: 'approle',
  },
  {
    category: 'cloud',
    displayName: 'AWS',
    glyph: 'aws-color',
    type: 'aws',
  },
  {
    category: 'cloud',
    displayName: 'Azure',
    glyph: 'azure-color',
    type: 'azure',
  },
  {
    category: 'infra',
    displayName: 'Consul',
    glyph: 'consul-color',
    type: 'consul',
  },
  {
    displayName: 'Cubbyhole',
    type: 'cubbyhole',
  },
  {
    category: 'infra',
    displayName: 'Databases',
    glyph: 'database',
    type: 'database',
  },
  {
    category: 'cloud',
    displayName: 'GitHub',
    glyph: 'github-color',
    type: 'github',
    value: 'github',
  },
  {
    category: 'cloud',
    displayName: 'Google Cloud',
    glyph: 'gcp-color',
    type: 'gcp',
  },
  {
    category: 'cloud',
    displayName: 'Google Cloud KMS',
    glyph: 'gcp-color',
    type: 'gcpkms',
  },
  {
    category: 'generic',
    displayName: 'JWT',
    glyph: 'jwt',
    type: 'jwt',
    value: 'jwt',
  },
  {
    category: 'generic',
    displayName: 'KV',
    glyph: 'key-values',
    type: 'kv',
  },
  {
    category: 'generic',
    displayName: 'KMIP',
    glyph: 'lock',
    requiresEnterprise: true,
    type: 'kmip',
  },
  {
    category: 'generic',
    displayName: 'Transform',
    glyph: 'transform-data',
    requiresEnterprise: true,
    type: 'transform',
  },
  {
    category: 'cloud',
    displayName: 'Key Management',
    glyph: 'key',
    requiresEnterprise: true,
    type: 'keymgmt',
  },
  {
    category: 'generic',
    displayName: 'Kubernetes',
    glyph: 'kubernetes-color',
    type: 'kubernetes',
  },
  {
    category: 'generic',
    displayName: 'LDAP',
    glyph: 'folder-users',
    type: 'ldap',
  },
  {
    category: 'infra',
    displayName: 'Nomad',
    glyph: 'nomad-color',
    type: 'nomad',
  },
  {
    category: 'generic',
    displayName: 'OIDC',
    glyph: 'openid-color',
    type: 'oidc',
    value: 'oidc',
  },
  {
    category: 'infra',
    displayName: 'Okta',
    glyph: 'okta-color',
    type: 'okta',
    value: 'okta',
  },
  {
    category: 'generic',
    displayName: 'PKI Certificates',
    glyph: 'certificate',
    type: 'pki',
  },
  {
    category: 'infra',
    displayName: 'RADIUS',
    glyph: 'mainframe',
    type: 'radius',
    value: 'radius',
  },
  {
    category: 'infra',
    displayName: 'RabbitMQ',
    glyph: 'rabbitmq-color',
    type: 'rabbitmq',
  },
  {
    category: 'generic',
    displayName: 'SAML',
    glyph: 'saml-color',
    requiresEnterprise: true,
    type: 'saml',
    value: 'saml',
  },
  {
    category: 'generic',
    displayName: 'SSH',
    glyph: 'terminal-screen',
    type: 'ssh',
  },
  {
    category: 'generic',
    displayName: 'TLS Certificates',
    glyph: 'certificate',
    type: 'cert',
    value: 'cert',
  },
  {
    category: 'generic',
    displayName: 'TOTP',
    glyph: 'history',
    type: 'totp',
  },
  {
    category: 'generic',
    displayName: 'Transit',
    glyph: 'swap-horizontal',
    type: 'transit',
  },
  {
    displayName: 'Token',
    type: 'token',
  },
  {
    category: 'generic',
    displayName: 'Userpass',
    glyph: 'users',
    type: 'userpass',
    value: 'userpass',
  },
];
