/**
 * Copyright IBM Corp. 2016, 2025
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
 * If an enterprise license is present, return all secret engines;
 * otherwise, return only the secret engines supported in OSS.
 * return filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: this.version.isEnterprise });
 */

export interface EngineDisplayData {
  category?: string; // category for NEW grouping of engines - to replace pluginCategory
  pluginCategory?: string; // The plugin category is used to group engines in the UI. e.g., 'cloud', 'infra', 'generic'
  displayName: string;
  engineRoute?: string; // engines that have their own Ember engine will have this route defined.
  glyph?: string;
  isWIF?: boolean; // flag for 'Workload Identity Federation' engines. - https://developer.hashicorp.com/hcp/docs/hcp/iam/service-principal/workload-identity-federation
  mountCategory: string[];
  requiredFeature?: string; // flag for engines that require the ADP (Advanced Data Protection) feature. - https://www.hashicorp.com/en/blog/advanced-data-protection-adp-now-available-in-hcp-vault
  requiresEnterprise?: boolean;
  isConfigurable?: boolean; // for secret engines that have additional configuration pages and actions.
  isOnlyMountable?: boolean; // The UI only supports configuration views for these secrets engines. The CLI must be used to manage other engine resources (i.e. roles, credentials).
  type: string;
  value?: string;
  configRoute?: string; // override for custom route if not "configuration.plugin-settings" (used for Ember engines)
  capabilities?: string[];
  description?: string;
  secretTypes?: string[]; // filter labels for the catalog secret type filter, e.g. 'cloud credentials', 'encryption keys'
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
  if (type === 'kv' && version === 1) {
    return false;
  }
  const engineRoute = ALL_ENGINES.find((engine) => engine.type === type)?.engineRoute;
  return !!engineRoute;
}

//  The "sys/mounts" and "sys/internal/ui/mounts" endpoints return a "secret/" key containing
//  all mounts enabled in Vault. Some types are internal Vault APIs, not user-mountable secrets engines,
//  and should be filtered in some scenarios, such as listing secrets engines.
export const INTERNAL_ENGINE_TYPES = ['system', 'identity', 'agent_registry'];

// These engines rely on a specific version being set (eg. "kv version 1" or "kv version 2") via options parameter when being mounted in Vault.
// If the engine to be mounted isn't specified here, we ignore the 'options' field.
export const VERSIONED_ENGINE_TYPES = ['vault-plugin-secrets-kv', 'kv', 'generic'];

export const ALL_ENGINES: EngineDisplayData[] = [
  {
    category: 'cloud and infrastructure',
    pluginCategory: 'cloud',
    displayName: 'AliCloud',
    glyph: 'alibaba-color',
    mountCategory: ['auth', 'secret'],
    type: 'alicloud',
    capabilities: ['dynamic', 'rotating'],
    description: 'Manage dynamic secrets in AliCloud',
    secretTypes: ['cloudCredentials'],
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
    category: 'common engines',
    pluginCategory: 'cloud',
    displayName: 'AWS',
    glyph: 'aws-color',
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'aws',
    capabilities: ['dynamic', 'rotating'],
    description: 'Generate dynamic AWS credentials with configurable IAM permissions.',
    secretTypes: ['cloudCredentials', 'apiKeysTokens'],
  },
  {
    category: 'common engines',
    pluginCategory: 'cloud',
    displayName: 'Azure',
    glyph: 'azure-color',
    isOnlyMountable: true,
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'azure',
    capabilities: ['dynamic', 'rotating'],
    description: 'Manage dynamic secrets in Azure',
    secretTypes: ['cloudCredentials', 'apiKeysTokens'],
  },
  {
    category: 'cloud and infrastructure',
    pluginCategory: 'infra',
    displayName: 'Consul',
    glyph: 'consul-color',
    mountCategory: ['secret'],
    type: 'consul',
    capabilities: ['dynamic'],
    description: 'Store and retrieve secrets from Consul',
    secretTypes: ['apiKeysTokens'],
  },
  {
    displayName: 'Cubbyhole',
    type: 'cubbyhole',
    mountCategory: ['secret'],
  },
  {
    category: 'common engines',
    pluginCategory: 'infra',
    displayName: 'Databases',
    glyph: 'database',
    mountCategory: ['secret'],
    type: 'database',
    capabilities: ['dynamic', 'rotating'],
    description: 'Generate dynamic database credentials for MySQL, PostgreSQL, MongoDB, and more.',
    secretTypes: ['databaseCredentials'],
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
    category: 'common engines',
    pluginCategory: 'cloud',
    displayName: 'Google Cloud',
    glyph: 'gcp-color',
    isOnlyMountable: true,
    isConfigurable: true,
    isWIF: true,
    mountCategory: ['auth', 'secret'],
    type: 'gcp',
    capabilities: ['dynamic', 'rotating'],
    description: 'Manage GCP service account keys and secrets',
    secretTypes: ['cloudCredentials', 'apiKeysTokens'],
  },
  {
    category: 'cryptography and data protection',
    pluginCategory: 'cloud',
    displayName: 'Google Cloud KMS',
    glyph: 'gcp-color',
    mountCategory: ['secret'],
    type: 'gcpkms',
    capabilities: ['encryption', 'signing'],
    description: 'Integrate with Google Cloud KMS for encryption',
    secretTypes: ['encryptionKeys'],
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
    category: 'common engines',
    pluginCategory: 'generic',
    displayName: 'KV',
    engineRoute: 'kv.list',
    configRoute: 'kv.configuration', // only utilized to display config data for kvv2, not in conjunction with isConfigurable as templates determine whether engine is kv v1 or v2
    glyph: 'key-values',
    mountCategory: ['secret'],
    type: 'kv',
    capabilities: ['static'],
    description: 'Store static keys, including code snippets. No auto-rotating support.',
    secretTypes: ['staticStorage'],
  },
  {
    category: 'cryptography and data protection',
    pluginCategory: 'generic',
    displayName: 'KMIP',
    engineRoute: 'kmip.scopes.index',
    configRoute: 'kmip.configuration',
    isConfigurable: true,
    glyph: 'lock',
    mountCategory: ['secret'],
    requiredFeature: 'KMIP',
    requiresEnterprise: true,
    type: 'kmip',
    capabilities: ['encryption', 'certificate authority'],
    description: 'Store and manage encryption keys using KMIP',
    secretTypes: ['encryptionKeys', 'certificatesPki'],
  },
  {
    category: 'cryptography and data protection',
    pluginCategory: 'generic',
    displayName: 'Transform',
    glyph: 'transform-data',
    mountCategory: ['secret'],
    requiredFeature: 'Transform Secrets Engine',
    requiresEnterprise: true,
    type: 'transform',
    capabilities: ['encryption', 'tokenization'],
    description: 'Perform data masking, encryption, and tokenization',
    secretTypes: ['encryptionKeys'],
  },
  {
    category: 'cryptography and data protection',
    pluginCategory: 'cloud',
    displayName: 'Key Management',
    glyph: 'key',
    mountCategory: ['secret'],
    requiredFeature: 'Key Management Secrets Engine',
    requiresEnterprise: true,
    type: 'keymgmt',
    capabilities: ['encryption', 'signing'],
    description: 'Centralize AWS, GCP, and Azure keys',
    secretTypes: ['encryptionKeys', 'cloudCredentials'],
  },
  {
    category: 'cloud and infrastructure',
    pluginCategory: 'generic',
    displayName: 'Kubernetes',
    engineRoute: 'kubernetes.overview',
    configRoute: 'kubernetes.configuration',
    glyph: 'kubernetes-color',
    isConfigurable: true,
    mountCategory: ['auth', 'secret'],
    type: 'kubernetes',
    capabilities: ['dynamic'],
    description: 'Manage Kubernetes secrets and configurations',
    secretTypes: ['apiKeysTokens'],
  },
  {
    category: 'identity and access',
    pluginCategory: 'generic',
    displayName: 'LDAP',
    isConfigurable: true,
    engineRoute: 'ldap.overview',
    configRoute: 'ldap.configuration',
    glyph: 'folder-users',
    mountCategory: ['auth', 'secret'],
    type: 'ldap',
    capabilities: ['static', 'dynamic', 'rotating'],
    description: 'Create auto-rotating static roles and dynamic roles for your AD, RACF, or OpenLDAP',
    secretTypes: ['staticStorage', 'apiKeysTokens'],
  },
  {
    category: 'cloud and infrastructure',
    pluginCategory: 'infra',
    displayName: 'Nomad',
    glyph: 'nomad-color',
    mountCategory: ['secret'],
    type: 'nomad',
    capabilities: ['dynamic'],
    description: 'Manage secrets for Nomad jobs and services',
    secretTypes: ['apiKeysTokens'],
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
    category: 'cryptography and data protection',
    pluginCategory: 'generic',
    displayName: 'Private PKI',
    isConfigurable: true,
    engineRoute: 'pki.overview',
    configRoute: 'pki.configuration',
    glyph: 'certificate',
    mountCategory: ['secret'],
    type: 'pki',
    capabilities: ['dynamic', 'certificate authority'],
    description: 'Generate and manage internal X.509 certificates',
    secretTypes: ['certificatesPki'],
  },
  {
    category: 'cryptography and data protection',
    pluginCategory: 'generic',
    displayName: 'Public PKI',
    engineRoute: 'pki.external.overview',
    glyph: 'certificate',
    mountCategory: ['secret'],
    requiresEnterprise: true,
    type: 'pki-external-ca',
    capabilities: ['dynamic', 'certificate authority'],
    description: 'Generate and manage X.509 certificates',
    secretTypes: ['certificatesPki'],
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
    category: 'cloud and infrastructure',
    pluginCategory: 'infra',
    displayName: 'RabbitMQ',
    glyph: 'rabbitmq-color',
    mountCategory: ['secret'],
    type: 'rabbitmq',
    capabilities: ['dynamic', 'rotating'],
    description: 'Generate dynamic credentials for RabbitMQ',
    secretTypes: ['apiKeysTokens'],
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
    category: 'identity and access',
    pluginCategory: 'generic',
    displayName: 'SSH',
    glyph: 'terminal-screen',
    isConfigurable: true,
    mountCategory: ['secret'],
    type: 'ssh',
    capabilities: ['dynamic', 'signing'],
    description: 'Enable secure access to servers using SSH',
    secretTypes: ['sshKeys'],
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
    category: 'identity and access',
    pluginCategory: 'generic',
    displayName: 'TOTP',
    glyph: 'history',
    mountCategory: ['secret'],
    type: 'totp',
    capabilities: ['dynamic'],
    description: 'Generate time-based one-time passwords for MFA',
    secretTypes: ['apiKeysTokens'],
  },
  {
    category: 'cryptography and data protection',
    pluginCategory: 'generic',
    displayName: 'Transit',
    glyph: 'swap-horizontal',
    mountCategory: ['secret'],
    type: 'transit',
    capabilities: ['encryption', 'signing'],
    description: 'Secure data with cryptography as a service, including post-quantum',
    secretTypes: ['encryptionKeys'],
  },
  {
    displayName: 'Token',
    type: 'token',
    glyph: 'users',
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

  // TODO: enable builtin plugins after confirming with Product
  //
  // {
  //   pluginCategory: 'generic',
  //   displayName: 'Ad',
  //   glyph: 'folder',
  //   isOnlyMountable: true,
  //   mountCategory: ['secret'],
  //   type: 'ad',
  // },
  // {
  //   pluginCategory: 'cloud',
  //   displayName: 'MongoDB Atlas',
  //   glyph: 'mongodb-color',
  //   isOnlyMountable: true,
  //   mountCategory: ['secret'],
  //   type: 'mongodbatlas',
  // },
  // {
  //   pluginCategory: 'infra',
  //   displayName: 'OpenLDAP',
  //   glyph: 'folder-users',
  //   isOnlyMountable: true,
  //   mountCategory: ['secret'],
  //   type: 'openldap',
  // },
  // {
  //   pluginCategory: 'infra',
  //   displayName: 'Terraform',
  //   glyph: 'terraform-color',
  //   isOnlyMountable: true,
  //   mountCategory: ['secret'],
  //   type: 'terraform',
  // },
];
