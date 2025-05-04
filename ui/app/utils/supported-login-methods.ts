/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * The web UI only supports logging in with these auth methods.
 * The method data is all related to logic for authenticating via that method.
 * This is a subset of the methods found in the `mountable-auth-methods` util,
 * which lists all the methods that can be enabled and mounted.
 */

export const BASE_LOGIN_METHODS = [
  {
    type: 'token',
    displayName: 'Token',
  },
  {
    type: 'userpass',
    displayName: 'Userpass',
  },
  {
    type: 'ldap',
    displayName: 'LDAP',
  },
  {
    type: 'okta',
    displayName: 'Okta',
  },
  {
    type: 'jwt',
    displayName: 'JWT',
  },
  {
    type: 'oidc',
    displayName: 'OIDC',
  },
  {
    type: 'radius',
    displayName: 'RADIUS',
  },
  {
    type: 'github',
    displayName: 'GitHub',
  },
];

export const ENTERPRISE_LOGIN_METHODS = [
  {
    type: 'saml',
    displayName: 'SAML',
  },
];

export const ALL_LOGIN_METHODS = [...BASE_LOGIN_METHODS, ...ENTERPRISE_LOGIN_METHODS];

export const supportedTypes = (isEnterprise: boolean) =>
  isEnterprise ? ALL_LOGIN_METHODS.map((m) => m.type) : BASE_LOGIN_METHODS.map((m) => m.type);

// this ensures no unexpected params are injected and submitted in the login form
// 'namespace' and 'path' are intentionally omitted because they are handled explicitly
export const POSSIBLE_FIELDS = ['role', 'jwt', 'password', 'token', 'username'];
