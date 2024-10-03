/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

/**
 * These are all the auth methods that can be mounted.
 * Some methods may not be available for login via the UI,
 * which are in the `supported-auth-backends` helper.
 */

const ENTERPRISE_AUTH_METHODS = [
  {
    displayName: 'SAML',
    value: 'saml',
    type: 'saml',
    category: 'generic',
    glyph: 'saml-color',
  },
];

const MOUNTABLE_AUTH_METHODS = [
  {
    displayName: 'AliCloud',
    value: 'alicloud',
    type: 'alicloud',
    category: 'cloud',
    glyph: 'alibaba-color',
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
  return MOUNTABLE_AUTH_METHODS.slice();
}

export function allMethods() {
  return [...MOUNTABLE_AUTH_METHODS, ...ENTERPRISE_AUTH_METHODS];
}

export default buildHelper(methods);
