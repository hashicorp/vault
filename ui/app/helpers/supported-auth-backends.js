/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

/**
 * These are all the auth methods with which a user can log into the UI.
 * This is a subset of the methods found in the `mountable-auth-methods` helper,
 * which lists all the methods that can be mounted.
 */

const SUPPORTED_AUTH_BACKENDS = [
  {
    type: 'token',
    typeDisplay: 'Token',
    description: 'Token authentication.',
    tokenPath: 'id',
    displayNamePath: 'display_name',
    formAttributes: ['token'],
  },
  {
    type: 'userpass',
    typeDisplay: 'Username',
    description: 'A simple username and password backend.',
    tokenPath: 'client_token',
    displayNamePath: 'metadata.username',
    formAttributes: ['username', 'password'],
  },
  {
    type: 'ldap',
    typeDisplay: 'LDAP',
    description: 'LDAP authentication.',
    tokenPath: 'client_token',
    displayNamePath: 'metadata.username',
    formAttributes: ['username', 'password'],
  },
  {
    type: 'okta',
    typeDisplay: 'Okta',
    description: 'Authenticate with your Okta username and password.',
    tokenPath: 'client_token',
    displayNamePath: 'metadata.username',
    formAttributes: ['username', 'password'],
  },
  {
    type: 'jwt',
    typeDisplay: 'JWT',
    description: 'Authenticate using JWT or OIDC provider.',
    tokenPath: 'client_token',
    displayNamePath: 'display_name',
    formAttributes: ['role', 'jwt'],
  },
  {
    type: 'oidc',
    typeDisplay: 'OIDC',
    description: 'Authenticate using JWT or OIDC provider.',
    tokenPath: 'client_token',
    displayNamePath: 'display_name',
    formAttributes: ['role', 'jwt'],
  },
  {
    type: 'radius',
    typeDisplay: 'RADIUS',
    description: 'Authenticate with your RADIUS username and password.',
    tokenPath: 'client_token',
    displayNamePath: 'metadata.username',
    formAttributes: ['username', 'password'],
  },
  {
    type: 'github',
    typeDisplay: 'GitHub',
    description: 'GitHub authentication.',
    tokenPath: 'client_token',
    displayNamePath: ['metadata.org', 'metadata.username'],
    formAttributes: ['token'],
  },
];

const ENTERPRISE_AUTH_METHODS = [
  {
    type: 'saml',
    typeDisplay: 'SAML',
    description: 'Authenticate using SAML provider.',
    tokenPath: 'client_token',
    displayNamePath: 'display_name',
    formAttributes: ['role'],
  },
];

export function supportedAuthBackends() {
  return [...SUPPORTED_AUTH_BACKENDS];
}

export function allSupportedAuthBackends() {
  return [...SUPPORTED_AUTH_BACKENDS, ...ENTERPRISE_AUTH_METHODS];
}

export default buildHelper(supportedAuthBackends);
