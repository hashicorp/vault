/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * These are all the secret and auth methods, including enterprise.
 * This is to be used for displaying the method name in the UI.
 * TODO: maybe add param on object that mentions secret or auth or enterprise
 */

export const ALL_ENGINES = [
  {
    displayName: 'AliCloud',
    type: 'alicloud',
    glyph: 'alibaba-color',
    category: 'cloud',
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
    displayName: 'Cubbyhole',
    type: 'cubbyhole',
  },
  {
    displayName: 'Databases',
    type: 'database',
    category: 'infra',
    glyph: 'database',
  },
  {
    displayName: 'GitHub',
    value: 'github',
    type: 'github',
    category: 'cloud',
    glyph: 'github-color',
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
    displayName: 'JWT',
    value: 'jwt',
    type: 'jwt',
    glyph: 'jwt',
    category: 'generic',
  },
  {
    displayName: 'KV',
    type: 'kv',
    glyph: 'key-values',
    category: 'generic',
  },
  {
    displayName: 'Kubernetes',
    type: 'kubernetes',
    category: 'generic',
    glyph: 'kubernetes-color',
  },
  {
    displayName: 'LDAP',
    type: 'ldap',
    category: 'generic',
    glyph: 'folder-users',
  },
  {
    displayName: 'Nomad',
    type: 'nomad',
    glyph: 'nomad-color',
    category: 'infra',
  },
  {
    displayName: 'OIDC',
    value: 'oidc',
    type: 'oidc',
    glyph: 'openid-color',
    category: 'generic',
  },
  {
    displayName: 'Okta',
    value: 'okta',
    type: 'okta',
    category: 'infra',
    glyph: 'okta-color',
  },
  {
    displayName: 'PKI Certificates',
    type: 'pki',
    glyph: 'certificate',
    category: 'generic',
  },
  {
    displayName: 'RADIUS',
    value: 'radius',
    type: 'radius',
    glyph: 'mainframe',
    category: 'infra',
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
    displayName: 'TLS Certificates',
    value: 'cert',
    type: 'cert',
    category: 'generic',
    glyph: 'certificate',
  },
  {
    displayName: 'TOTP',
    type: 'totp',
    glyph: 'history',
    category: 'generic',
  },
  {
    displayName: 'Transit',
    type: 'transit',
    glyph: 'swap-horizontal',
    category: 'generic',
  },
  {
    displayName: 'Username & Password',
    value: 'userpass',
    type: 'userpass',
    category: 'generic',
    glyph: 'users',
  },
];
