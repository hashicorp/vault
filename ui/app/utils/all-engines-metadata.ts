/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * These are all the secret and auth methods, including enterprise.
 */

export interface EngineDisplayData {
  category?: string;
  displayName: string;
  glyph?: string;
  mountGroup: string[];
  requiresEnterprise?: boolean;
  type: string;
  value?: string;
}

/**
 * @param mountGroup - Given mount group to filter by, e.g., 'auth' or 'secret'.
 * @param isEnterprise - Optional boolean to indicate if enterprise engines should be included in the results.
 * @returns Filtered array of engines that match the given mount type
 */
export function filterEnginesByMountType({
  mountGroup,
  isEnterprise = false,
}: {
  mountGroup: 'auth' | 'secret';
  isEnterprise: boolean;
}) {
  return isEnterprise
    ? ALL_ENGINES.filter((engine) => engine.mountGroup.includes(mountGroup))
    : ALL_ENGINES.filter((engine) => engine.mountGroup.includes(mountGroup) && !engine.requiresEnterprise);
}

export const ALL_ENGINES: EngineDisplayData[] = [
  {
    category: 'cloud',
    displayName: 'AliCloud',
    glyph: 'alibaba-color',
    mountGroup: ['auth', 'secret'],
    type: 'alicloud',
  },
  {
    category: 'generic',
    displayName: 'AppRole',
    glyph: 'cpu',
    mountGroup: ['auth'],
    type: 'approle',
    value: 'approle',
  },
  {
    category: 'cloud',
    displayName: 'AWS',
    glyph: 'aws-color',
    mountGroup: ['auth', 'secret'],
    type: 'aws',
  },
  {
    category: 'cloud',
    displayName: 'Azure',
    glyph: 'azure-color',
    mountGroup: ['auth', 'secret'],
    type: 'azure',
  },
  {
    category: 'infra',
    displayName: 'Consul',
    glyph: 'consul-color',
    mountGroup: ['secret'],
    type: 'consul',
  },
  {
    displayName: 'Cubbyhole',
    type: 'cubbyhole',
    mountGroup: ['secret'],
  },
  {
    category: 'infra',
    displayName: 'Databases',
    glyph: 'database',
    mountGroup: ['secret'],
    type: 'database',
  },
  {
    category: 'cloud',
    displayName: 'GitHub',
    glyph: 'github-color',
    mountGroup: ['auth'],
    type: 'github',
    value: 'github',
  },
  {
    category: 'cloud',
    displayName: 'Google Cloud',
    glyph: 'gcp-color',
    mountGroup: ['auth', 'secret'],
    type: 'gcp',
  },
  {
    category: 'cloud',
    displayName: 'Google Cloud KMS',
    glyph: 'gcp-color',
    mountGroup: ['secret'],
    type: 'gcpkms',
  },
  {
    category: 'generic',
    displayName: 'JWT',
    glyph: 'jwt',
    mountGroup: ['auth'],
    type: 'jwt',
    value: 'jwt',
  },
  {
    category: 'generic',
    displayName: 'KV',
    glyph: 'key-values',
    mountGroup: ['secret'],
    type: 'kv',
  },
  {
    category: 'generic',
    displayName: 'KMIP',
    glyph: 'lock',
    mountGroup: ['secret'],
    requiresEnterprise: true,
    type: 'kmip',
  },
  {
    category: 'generic',
    displayName: 'Transform',
    glyph: 'transform-data',
    mountGroup: ['secret'],
    requiresEnterprise: true,
    type: 'transform',
  },
  {
    category: 'cloud',
    displayName: 'Key Management',
    glyph: 'key',
    mountGroup: ['secret'],
    requiresEnterprise: true,
    type: 'keymgmt',
  },
  {
    category: 'generic',
    displayName: 'Kubernetes',
    glyph: 'kubernetes-color',
    mountGroup: ['auth', 'secret'],
    type: 'kubernetes',
  },
  {
    category: 'generic',
    displayName: 'LDAP',
    glyph: 'folder-users',
    mountGroup: ['auth', 'secret'],
    type: 'ldap',
  },
  {
    category: 'infra',
    displayName: 'Nomad',
    glyph: 'nomad-color',
    mountGroup: ['secret'],
    type: 'nomad',
  },
  {
    category: 'generic',
    displayName: 'OIDC',
    glyph: 'openid-color',
    mountGroup: ['auth'],
    type: 'oidc',
    value: 'oidc',
  },
  {
    category: 'infra',
    displayName: 'Okta',
    glyph: 'okta-color',
    mountGroup: ['auth'],
    type: 'okta',
    value: 'okta',
  },
  {
    category: 'generic',
    displayName: 'PKI Certificates',
    glyph: 'certificate',
    mountGroup: ['secret'],
    type: 'pki',
  },
  {
    category: 'infra',
    displayName: 'RADIUS',
    glyph: 'mainframe',
    mountGroup: ['auth'],
    type: 'radius',
    value: 'radius',
  },
  {
    category: 'infra',
    displayName: 'RabbitMQ',
    glyph: 'rabbitmq-color',
    mountGroup: ['secret'],
    type: 'rabbitmq',
  },
  {
    category: 'generic',
    displayName: 'SAML',
    glyph: 'saml-color',
    mountGroup: ['auth'],
    requiresEnterprise: true,
    type: 'saml',
    value: 'saml',
  },
  {
    category: 'generic',
    displayName: 'SSH',
    glyph: 'terminal-screen',
    mountGroup: ['secret'],
    type: 'ssh',
  },
  {
    category: 'generic',
    displayName: 'TLS Certificates',
    glyph: 'certificate',
    mountGroup: ['auth'],
    type: 'cert',
    value: 'cert',
  },
  {
    category: 'generic',
    displayName: 'TOTP',
    glyph: 'history',
    mountGroup: ['secret'],
    type: 'totp',
  },
  {
    category: 'generic',
    displayName: 'Transit',
    glyph: 'swap-horizontal',
    mountGroup: ['secret'],
    type: 'transit',
  },
  {
    displayName: 'Token',
    type: 'token',
    mountGroup: ['auth'],
  },
  {
    category: 'generic',
    displayName: 'Userpass',
    glyph: 'users',
    mountGroup: ['auth'],
    type: 'userpass',
    value: 'userpass',
  },
];
