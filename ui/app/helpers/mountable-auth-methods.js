/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';

const MOUNTABLE_AUTH_METHODS = [
  {
    displayName: 'AliCloud',
    value: 'alicloud',
    type: 'alicloud',
    category: 'cloud',
  },
  {
    displayName: 'AppRole',
    value: 'approle',
    type: 'approle',
    category: 'generic',
  },
  {
    displayName: 'AWS',
    value: 'aws',
    type: 'aws',
    category: 'cloud',
    glyph: 'aws-color',
  },
  {
    displayName: 'Azure',
    value: 'azure',
    type: 'azure',
    category: 'cloud',
    glyph: 'azure-color',
  },
  {
    displayName: 'Google Cloud',
    value: 'gcp',
    type: 'gcp',
    category: 'cloud',
    glyph: 'gcp-color',
  },
  {
    displayName: 'GitHub',
    value: 'github',
    type: 'github',
    category: 'cloud',
    glyph: 'github-color',
  },
  {
    displayName: 'JWT',
    value: 'jwt',
    type: 'jwt',
    glyph: 'auth',
    category: 'generic',
  },
  {
    displayName: 'OIDC',
    value: 'oidc',
    type: 'oidc',
    glyph: 'auth',
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
    glyph: 'auth',
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
    glyph: 'auth',
    category: 'infra',
  },
  {
    displayName: 'TLS Certificates',
    value: 'cert',
    type: 'cert',
    category: 'generic',
  },
  {
    displayName: 'Username & Password',
    value: 'userpass',
    type: 'userpass',
    category: 'generic',
  },
];

export function methods() {
  return MOUNTABLE_AUTH_METHODS.slice();
}

export default buildHelper(methods);
