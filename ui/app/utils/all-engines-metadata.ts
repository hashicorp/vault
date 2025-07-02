/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Metadata configuration for secret and auth engines, including enterprise.
 *
 * This file defines and exports engine metadata, including its
 * displayName, mountCategory, requiresEnterprise, and other relevant properties. It serves as a
 * centralized source of truth for engine-related configurations.
 *
 * Key responsibilities:
 * - Define metadata for all engines.
 * - Provide utility functions or constants for accessing engine-specific data.
 * - Facilitate dynamic engine rendering and behavior based on metadata.
 *
 * Example usage:
 * // If an enterprise license is present, return all secret engines;
 * // otherwise, return only the secret engines supported in OSS.
 * return filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: this.version.isEnterprise });
 */

export interface EngineDisplayData {
  pluginCategory?: string; // The plugin category is used to group engines in the UI. e.g., 'cloud', 'infra', 'generic'
  displayName: string;
  engineRoute?: string;
  glyph?: string;
  isWIF?: boolean; // flag for 'Workload Identity Federation' engines.
  mountCategory: string[];
  requiredFeature?: string; // flag for engines that require the ADP (Advanced Data Protection) feature. - https://www.hashicorp.com/en/blog/advanced-data-protection-adp-now-available-in-hcp-vault
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
    pluginCategory: 'cloud',
    displayName: 'AliCloud',
    glyph: 'alibaba-color',
    mountCategory: ['auth', 'secret'],
    type: 'alicloud',
  },
  {
    pluginCategory: 'generic',
    displayName: 'AppRole',
    glyph: 'cpu',
    mountCategory: ['auth'],
    type: 'approle',
    value: 'approle',
  },
  {
    pluginCategory: 'cloud',
    displayName: 'AWS',
    glyph: 'aws-color',
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'aws',
  },
  {
    pluginCategory: 'cloud',
    displayName: 'Azure',
    glyph: 'azure-color',
    isOnlyMountable: true,
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'azure',
  },
  {
    pluginCategory: 'infra',
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
    pluginCategory: 'infra',
    displayName: 'Databases',
    glyph: 'database',
    mountCategory: ['secret'],
    type: 'database',
  },
  {
    pluginCategory: 'cloud',
    displayName: 'GitHub',
    glyph: 'github-color',
    mountCategory: ['auth'],
    type: 'github',
    value: 'github',
  },
  {
    pluginCategory: 'cloud',
    displayName: 'Google Cloud',
    glyph: 'gcp-color',
    isOnlyMountable: true,
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'gcp',
  },
  {
    pluginCategory: 'cloud',
    displayName: 'Google Cloud KMS',
    glyph: 'gcp-color',
    mountCategory: ['secret'],
    type: 'gcpkms',
  },
  {
    pluginCategory: 'generic',
    displayName: 'JWT',
    glyph: 'jwt',
    mountCategory: ['auth'],
    type: 'jwt',
    value: 'jwt',
  },
  {
    pluginCategory: 'generic',
    displayName: 'KV',
    engineRoute: 'kv.list',
    glyph: 'key-values',
    mountCategory: ['secret'],
    type: 'kv',
  },
  {
    pluginCategory: 'generic',
    displayName: 'KMIP',
    engineRoute: 'kmip.scopes.index',
    glyph: 'lock',
    mountCategory: ['secret'],
    requiredFeature: 'KMIP',
    requiresEnterprise: true,
    type: 'kmip',
  },
  {
    pluginCategory: 'generic',
    displayName: 'Transform',
    glyph: 'transform-data',
    mountCategory: ['secret'],
    requiredFeature: 'Transform Secrets Engine',
    requiresEnterprise: true,
    type: 'transform',
  },
  {
    pluginCategory: 'cloud',
    displayName: 'Key Management',
    glyph: 'key',
    mountCategory: ['secret'],
    requiredFeature: 'Key Management Secrets Engine',
    requiresEnterprise: true,
    type: 'keymgmt',
  },
  {
    pluginCategory: 'generic',
    displayName: 'Kubernetes',
    engineRoute: 'kubernetes.overview',
    glyph: 'kubernetes-color',
    mountCategory: ['auth', 'secret'],
    type: 'kubernetes',
  },
  {
    pluginCategory: 'generic',
    displayName: 'LDAP',
    engineRoute: 'ldap.overview',
    glyph: 'folder-users',
    mountCategory: ['auth', 'secret'],
    type: 'ldap',
  },
  {
    pluginCategory: 'infra',
    displayName: 'Nomad',
    glyph: 'nomad-color',
    mountCategory: ['secret'],
    type: 'nomad',
  },
  {
    pluginCategory: 'generic',
    displayName: 'OIDC',
    glyph: 'openid-color',
    mountCategory: ['auth'],
    type: 'oidc',
    value: 'oidc',
  },
  {
    pluginCategory: 'infra',
    displayName: 'Okta',
    glyph: 'okta-color',
    mountCategory: ['auth'],
    type: 'okta',
    value: 'okta',
  },
  {
    pluginCategory: 'generic',
    displayName: 'PKI Certificates',
    engineRoute: 'pki.overview',
    glyph: 'certificate',
    mountCategory: ['secret'],
    type: 'pki',
  },
  {
    pluginCategory: 'infra',
    displayName: 'RADIUS',
    glyph: 'mainframe',
    mountCategory: ['auth'],
    type: 'radius',
    value: 'radius',
  },
  {
    pluginCategory: 'infra',
    displayName: 'RabbitMQ',
    glyph: 'rabbitmq-color',
    mountCategory: ['secret'],
    type: 'rabbitmq',
  },
  {
    pluginCategory: 'generic',
    displayName: 'SAML',
    glyph: 'saml-color',
    mountCategory: ['auth'],
    requiresEnterprise: true,
    type: 'saml',
    value: 'saml',
  },
  {
    pluginCategory: 'generic',
    displayName: 'SSH',
    glyph: 'terminal-screen',
    isConfigurable: true,
    mountCategory: ['secret'],
    type: 'ssh',
  },
  {
    pluginCategory: 'generic',
    displayName: 'TLS Certificates',
    glyph: 'certificate',
    mountCategory: ['auth'],
    type: 'cert',
    value: 'cert',
  },
  {
    pluginCategory: 'generic',
    displayName: 'TOTP',
    glyph: 'history',
    mountCategory: ['secret'],
    type: 'totp',
  },
  {
    pluginCategory: 'generic',
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
    pluginCategory: 'generic',
    displayName: 'Userpass',
    glyph: 'users',
    mountCategory: ['auth'],
    type: 'userpass',
    value: 'userpass',
  },
];
