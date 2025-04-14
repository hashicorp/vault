/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

/**
 * These are all the secret and auth methods, including enterprise.
 * This is to be used for displaying the method name in the UI.
 * TODO: maybe add param on object that mentions secret or auth or enterprise
 */

const ALL_ENGINES = [
  {
    displayName: 'Cubbyhole',
    type: 'cubbyhole',
  },
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
    displayName: 'Key Management',
    type: 'keymgmt',
    glyph: 'key',
    category: 'cloud',
    requiredFeature: 'Key Management Secrets Engine',
    routeQueryParams: { tab: 'provider' },
  },
  {
    displayName: 'KMIP',
    type: 'kmip',
    glyph: 'lock',
    engineRoute: 'kmip.scopes.index',
    category: 'generic',
    requiredFeature: 'KMIP',
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
    displayName: 'Transform',
    type: 'transform',
    category: 'generic',
    requiredFeature: 'Transform Secrets Engine',
    glyph: 'transform-data',
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
  {
    displayName: 'AliCloud',
    type: 'alicloud',
  },
  {
    displayName: 'AppRole',
    type: 'approle',
  },
  {
    displayName: 'AWS',
    type: 'aws',
  },
  {
    displayName: 'Azure',
    type: 'azure',
  },
  {
    displayName: 'Google Cloud',
    type: 'gcp',
  },
  {
    displayName: 'GitHub',
    type: 'github',
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

export function methods() {
  return ALL_ENGINES.slice();
}

export default buildHelper(methods);
