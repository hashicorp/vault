/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

const ENTERPRISE_SECRET_ENGINES = [
  {
    displayName: 'KMIP',
    type: 'kmip',
    glyph: 'lock',
    engineRoute: 'kmip.scopes.index',
    category: 'generic',
    requiredFeature: 'KMIP',
  },
  {
    displayName: 'Transform',
    type: 'transform',
    category: 'generic',
    requiredFeature: 'Transform Secrets Engine',
    glyph: 'transform-data',
  },
  {
    displayName: 'Key Management',
    type: 'keymgmt',
    glyph: 'key',
    category: 'cloud',
    requiredFeature: 'Key Management Secrets Engine',
    routeQueryParams: { tab: 'provider' },
  },
];

const MOUNTABLE_SECRET_ENGINES = [
  {
    displayName: 'AliCloud',
    type: 'alicloud',
    glyph: 'alibaba-color',
    category: 'cloud',
  },
  {
    displayName: 'AWS',
    type: 'aws',
    category: 'cloud',
    glyph: 'aws-color',
  },
  {
    displayName: 'Azure',
    type: 'azure',
    category: 'cloud',
    glyph: 'azure-color',
  },
  {
    displayName: 'Consul',
    type: 'consul',
    glyph: 'consul-color',
    category: 'infra',
  },
  {
    displayName: 'Databases',
    type: 'database',
    category: 'infra',
    glyph: 'database',
  },
  {
    displayName: 'Google Cloud',
    type: 'gcp',
    category: 'cloud',
    glyph: 'gcp-color',
  },
  {
    displayName: 'Google Cloud KMS',
    type: 'gcpkms',
    category: 'cloud',
    glyph: 'gcp-color',
  },
  {
    displayName: 'KV',
    type: 'kv',
    glyph: 'key-values',
    engineRoute: 'kv.list',
    category: 'generic',
  },
  {
    displayName: 'Nomad',
    type: 'nomad',
    glyph: 'nomad-color',
    category: 'infra',
  },
  {
    displayName: 'PKI Certificates',
    type: 'pki',
    glyph: 'certificate',
    engineRoute: 'pki.overview',
    category: 'generic',
  },
  {
    displayName: 'RabbitMQ',
    type: 'rabbitmq',
    glyph: 'rabbitmq-color',
    category: 'infra',
  },
  {
    displayName: 'SSH',
    type: 'ssh',
    glyph: 'terminal-screen',
    category: 'generic',
  },
  {
    displayName: 'Transit',
    type: 'transit',
    glyph: 'swap-horizontal',
    category: 'generic',
  },
  {
    displayName: 'TOTP',
    type: 'totp',
    glyph: 'history',
    category: 'generic',
  },
  {
    displayName: 'LDAP',
    type: 'ldap',
    engineRoute: 'ldap.overview',
    category: 'generic',
    glyph: 'folder-users',
  },
  {
    displayName: 'Kubernetes',
    type: 'kubernetes',
    engineRoute: 'kubernetes.overview',
    category: 'generic',
    glyph: 'kubernetes-color',
  },
];

// A list of Workload Identity Federation engines.
export const WIF_ENGINES = ['aws', 'azure'];

export function wifEngines() {
  return WIF_ENGINES.slice();
}

// The UI only supports configuration views for these secrets engines. The CLI must be used to manage other engine resources (i.e. roles, credentials).
export const CONFIGURATION_ONLY = ['azure', 'gcp'];

export function configurationOnly() {
  return CONFIGURATION_ONLY.slice();
}

// Secret engines that have their own configuration page and actions
// These engines do not exist in their own Ember engine.
export const CONFIGURABLE_SECRET_ENGINES = ['aws', 'azure', 'gcp', 'ssh'];

export function configurableSecretEngines() {
  return CONFIGURABLE_SECRET_ENGINES.slice();
}

export function mountableEngines() {
  return MOUNTABLE_SECRET_ENGINES.slice();
}
// secret engines that have not other views than the mount view and mount details view
export const UNSUPPORTED_ENGINES = ['alicloud', 'consul', 'gcpkms', 'nomad', 'rabbitmq', 'totp'];

export function unsupportedEngines() {
  return UNSUPPORTED_ENGINES.slice();
}

export function allEngines() {
  return [...MOUNTABLE_SECRET_ENGINES, ...ENTERPRISE_SECRET_ENGINES];
}

export function isAddonEngine(type, version) {
  if (type === 'kv' && version === 1) return false;
  const engineRoute = allEngines().find((engine) => engine.type === type)?.engineRoute;
  return !!engineRoute;
}

export default buildHelper(mountableEngines);
