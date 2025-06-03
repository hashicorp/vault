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
  mountType: 'secret' | 'auth' | 'both';
  requiresEnterprise?: boolean;
  type: string;
  value?: string;
}

export const ALL_ENGINES: EngineDisplayData[] = [
  {
    category: 'cloud',
    displayName: 'AliCloud',
    glyph: 'alibaba-color',
    mountType: 'both',
    type: 'alicloud',
  },
  {
    category: 'generic',
    displayName: 'AppRole',
    glyph: 'cpu',
    mountType: 'auth',
    type: 'approle',
    value: 'approle',
  },
  {
    category: 'cloud',
    displayName: 'AWS',
    glyph: 'aws-color',
    mountType: 'both',
    type: 'aws',
  },
  {
    category: 'cloud',
    displayName: 'Azure',
    glyph: 'azure-color',
    mountType: 'both',
    type: 'azure',
  },
  {
    category: 'infra',
    displayName: 'Consul',
    glyph: 'consul-color',
    mountType: 'secret',
    type: 'consul',
  },
  {
    displayName: 'Cubbyhole',
    type: 'cubbyhole',
    mountType: 'secret',
  },
  {
    category: 'infra',
    displayName: 'Databases',
    glyph: 'database',
    mountType: 'secret',
    type: 'database',
  },
  {
    category: 'cloud',
    displayName: 'GitHub',
    glyph: 'github-color',
    mountType: 'auth',
    type: 'github',
    value: 'github',
  },
  {
    category: 'cloud',
    displayName: 'Google Cloud',
    glyph: 'gcp-color',
    mountType: 'both',
    type: 'gcp',
  },
  {
    category: 'cloud',
    displayName: 'Google Cloud KMS',
    glyph: 'gcp-color',
    mountType: 'secret',
    type: 'gcpkms',
  },
  {
    category: 'generic',
    displayName: 'JWT',
    glyph: 'jwt',
    mountType: 'auth',
    type: 'jwt',
    value: 'jwt',
  },
  {
    category: 'generic',
    displayName: 'KV',
    glyph: 'key-values',
    mountType: 'secret',
    type: 'kv',
  },
  {
    category: 'generic',
    displayName: 'KMIP',
    glyph: 'lock',
    mountType: 'secret',
    requiresEnterprise: true,
    type: 'kmip',
  },
  {
    category: 'generic',
    displayName: 'Transform',
    glyph: 'transform-data',
    mountType: 'secret',
    requiresEnterprise: true,
    type: 'transform',
  },
  {
    category: 'cloud',
    displayName: 'Key Management',
    glyph: 'key',
    mountType: 'secret',
    requiresEnterprise: true,
    type: 'keymgmt',
  },
  {
    category: 'generic',
    displayName: 'Kubernetes',
    glyph: 'kubernetes-color',
    mountType: 'both',
    type: 'kubernetes',
  },
  {
    category: 'generic',
    displayName: 'LDAP',
    glyph: 'folder-users',
    mountType: 'both',
    type: 'ldap',
  },
  {
    category: 'infra',
    displayName: 'Nomad',
    glyph: 'nomad-color',
    mountType: 'secret',
    type: 'nomad',
  },
  {
    category: 'generic',
    displayName: 'OIDC',
    glyph: 'openid-color',
    mountType: 'auth',
    type: 'oidc',
    value: 'oidc',
  },
  {
    category: 'infra',
    displayName: 'Okta',
    glyph: 'okta-color',
    mountType: 'auth',
    type: 'okta',
    value: 'okta',
  },
  {
    category: 'generic',
    displayName: 'PKI Certificates',
    glyph: 'certificate',
    mountType: 'secret',
    type: 'pki',
  },
  {
    category: 'infra',
    displayName: 'RADIUS',
    glyph: 'mainframe',
    mountType: 'auth',
    type: 'radius',
    value: 'radius',
  },
  {
    category: 'infra',
    displayName: 'RabbitMQ',
    glyph: 'rabbitmq-color',
    mountType: 'secret',
    type: 'rabbitmq',
  },
  {
    category: 'generic',
    displayName: 'SAML',
    glyph: 'saml-color',
    mountType: 'auth',
    requiresEnterprise: true,
    type: 'saml',
    value: 'saml',
  },
  {
    category: 'generic',
    displayName: 'SSH',
    glyph: 'terminal-screen',
    mountType: 'secret',
    type: 'ssh',
  },
  {
    category: 'generic',
    displayName: 'TLS Certificates',
    glyph: 'certificate',
    mountType: 'auth',
    type: 'cert',
    value: 'cert',
  },
  {
    category: 'generic',
    displayName: 'TOTP',
    glyph: 'history',
    mountType: 'secret',
    type: 'totp',
  },
  {
    category: 'generic',
    displayName: 'Transit',
    glyph: 'swap-horizontal',
    mountType: 'secret',
    type: 'transit',
  },
  {
    displayName: 'Token',
    type: 'token',
    mountType: 'auth',
  },
  {
    category: 'generic',
    displayName: 'Userpass',
    glyph: 'users',
    mountType: 'auth',
    type: 'userpass',
    value: 'userpass',
  },
];
