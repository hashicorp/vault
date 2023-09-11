/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

const ENTERPRISE_SECRET_ENGINES = [
  {
    displayName: 'KMIP',
    type: 'kmip',
    engineRoute: 'kmip.scopes.index',
    category: 'generic',
    requiredFeature: 'KMIP',
  },
  {
    displayName: 'Transform',
    type: 'transform',
    category: 'generic',
    requiredFeature: 'Transform Secrets Engine',
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
    category: 'infra',
  },
  {
    displayName: 'Databases',
    type: 'database',
    category: 'infra',
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
    engineRoute: 'kv.list',
    category: 'generic',
  },
  {
    displayName: 'Nomad',
    type: 'nomad',
    category: 'infra',
  },
  {
    displayName: 'PKI Certificates',
    type: 'pki',
    engineRoute: 'pki.overview',
    category: 'generic',
  },
  {
    displayName: 'RabbitMQ',
    type: 'rabbitmq',
    category: 'infra',
  },
  {
    displayName: 'SSH',
    type: 'ssh',
    category: 'generic',
  },
  {
    displayName: 'Transit',
    type: 'transit',
    category: 'generic',
  },
  {
    displayName: 'TOTP',
    type: 'totp',
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

export function mountableEngines() {
  return MOUNTABLE_SECRET_ENGINES.slice();
}

export function allEngines() {
  return [...MOUNTABLE_SECRET_ENGINES, ...ENTERPRISE_SECRET_ENGINES];
}

export function isAddonEngine(type, version) {
  if (type === 'kv' && version === 1) return false;
  const engineRoute = allEngines().findBy('type', type)?.engineRoute;
  return !!engineRoute;
}

export default buildHelper(mountableEngines);
