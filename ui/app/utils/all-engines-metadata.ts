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
  engineRoute?: string;
  glyph?: string;
  isWIF?: boolean; // flag for 'Workload Identity Federation' engines.
  mountCategory: string[];
  requiresEnterprise?: boolean;
  isConfigurable?: boolean; // for secret engines that have their own configuration page and actions. - These engines do not exist in their own Ember engine.
  isOnlyMountable?: boolean; // The UI only supports configuration views for these secrets engines. The CLI must be used to manage other engine resources (i.e. roles, credentials).
  type: string;
  value?: string;
}

/**
 * @param mountCategory - Given mount category to filter by, e.g., 'auth' or 'secret'.
 * @param isEnterprise - Optional boolean to indicate if enterprise engines should be included in the results.
 * @returns Filtered array of engines that match the given mount category
 */
export function filterEnginesByMountCategory({
  mountCategory,
  isEnterprise = false,
}: {
  mountCategory: 'auth' | 'secret';
  isEnterprise: boolean;
}) {
  return isEnterprise
    ? ALL_ENGINES.filter((engine) => engine.mountCategory.includes(mountCategory))
    : ALL_ENGINES.filter(
        (engine) => engine.mountCategory.includes(mountCategory) && !engine.requiresEnterprise
      );
}

export function isAddonEngine(type: string, version: number) {
  if (type === 'kv' && version === 1) return false;
  const engineRoute = ALL_ENGINES.find((engine) => engine.type === type)?.engineRoute;
  return !!engineRoute;
}

export const ALL_ENGINES: EngineDisplayData[] = [
  {
    category: 'cloud',
    displayName: 'AliCloud',
    glyph: 'alibaba-color',
    mountCategory: ['auth', 'secret'],
    type: 'alicloud',
  },
  {
    category: 'generic',
    displayName: 'AppRole',
    glyph: 'cpu',
    mountCategory: ['auth'],
    type: 'approle',
    value: 'approle',
  },
  {
    category: 'cloud',
    displayName: 'AWS',
    glyph: 'aws-color',
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'aws',
  },
  {
    category: 'cloud',
    displayName: 'Azure',
    glyph: 'azure-color',
    isOnlyMountable: true,
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'azure',
  },
  {
    category: 'infra',
    displayName: 'Consul',
    glyph: 'consul-color',
    mountCategory: ['secret'],
    type: 'consul',
  },
  {
    displayName: 'Cubbyhole',
    type: 'cubbyhole',
    mountCategory: ['secret'],
  },
  {
    category: 'infra',
    displayName: 'Databases',
    glyph: 'database',
    mountCategory: ['secret'],
    type: 'database',
  },
  {
    category: 'cloud',
    displayName: 'GitHub',
    glyph: 'github-color',
    mountCategory: ['auth'],
    type: 'github',
    value: 'github',
  },
  {
    category: 'cloud',
    displayName: 'Google Cloud',
    glyph: 'gcp-color',
    isOnlyMountable: true,
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'gcp',
  },
  {
    category: 'cloud',
    displayName: 'Google Cloud KMS',
    glyph: 'gcp-color',
    mountCategory: ['secret'],
    type: 'gcpkms',
  },
  {
    category: 'generic',
    displayName: 'JWT',
    glyph: 'jwt',
    mountCategory: ['auth'],
    type: 'jwt',
    value: 'jwt',
  },
  {
    category: 'generic',
    displayName: 'KV',
    engineRoute: 'kv.list',
    glyph: 'key-values',
    mountCategory: ['secret'],
    type: 'kv',
  },
  {
    category: 'generic',
    displayName: 'KMIP',
    engineRoute: 'kmip.scopes.index',
    glyph: 'lock',
    mountCategory: ['secret'],
    requiresEnterprise: true,
    type: 'kmip',
  },
  {
    category: 'generic',
    displayName: 'Transform',
    glyph: 'transform-data',
    mountCategory: ['secret'],
    requiresEnterprise: true,
    type: 'transform',
  },
  {
    category: 'cloud',
    displayName: 'Key Management',
    glyph: 'key',
    mountCategory: ['secret'],
    requiresEnterprise: true,
    type: 'keymgmt',
  },
  {
    category: 'generic',
    displayName: 'Kubernetes',
    engineRoute: 'kubernetes.overview',
    glyph: 'kubernetes-color',
    mountCategory: ['auth', 'secret'],
    type: 'kubernetes',
  },
  {
    category: 'generic',
    displayName: 'LDAP',
    engineRoute: 'ldap.overview',
    glyph: 'folder-users',
    mountCategory: ['auth', 'secret'],
    type: 'ldap',
  },
  {
    category: 'infra',
    displayName: 'Nomad',
    glyph: 'nomad-color',
    mountCategory: ['secret'],
    type: 'nomad',
  },
  {
    category: 'generic',
    displayName: 'OIDC',
    glyph: 'openid-color',
    mountCategory: ['auth'],
    type: 'oidc',
    value: 'oidc',
  },
  {
    category: 'infra',
    displayName: 'Okta',
    glyph: 'okta-color',
    mountCategory: ['auth'],
    type: 'okta',
    value: 'okta',
  },
  {
    category: 'generic',
    displayName: 'PKI Certificates',
    engineRoute: 'pki.overview',
    glyph: 'certificate',
    mountCategory: ['secret'],
    type: 'pki',
  },
  {
    category: 'infra',
    displayName: 'RADIUS',
    glyph: 'mainframe',
    mountCategory: ['auth'],
    type: 'radius',
    value: 'radius',
  },
  {
    category: 'infra',
    displayName: 'RabbitMQ',
    glyph: 'rabbitmq-color',
    mountCategory: ['secret'],
    type: 'rabbitmq',
  },
  {
    category: 'generic',
    displayName: 'SAML',
    glyph: 'saml-color',
    mountCategory: ['auth'],
    requiresEnterprise: true,
    type: 'saml',
    value: 'saml',
  },
  {
    category: 'generic',
    displayName: 'SSH',
    glyph: 'terminal-screen',
    isConfigurable: true,
    mountCategory: ['secret'],
    type: 'ssh',
  },
  {
    category: 'generic',
    displayName: 'TLS Certificates',
    glyph: 'certificate',
    mountCategory: ['auth'],
    type: 'cert',
    value: 'cert',
  },
  {
    category: 'generic',
    displayName: 'TOTP',
    glyph: 'history',
    mountCategory: ['secret'],
    type: 'totp',
  },
  {
    category: 'generic',
    displayName: 'Transit',
    glyph: 'swap-horizontal',
    mountCategory: ['secret'],
    type: 'transit',
  },
  {
    displayName: 'Token',
    type: 'token',
    mountCategory: ['auth'],
  },
  {
    category: 'generic',
    displayName: 'Userpass',
    glyph: 'users',
    mountCategory: ['auth'],
    type: 'userpass',
    value: 'userpass',
  },
];
