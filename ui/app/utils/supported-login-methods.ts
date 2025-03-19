/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * The web UI only supports logging in with these auth methods.
 * The method data is all related to logic for authenticating via that method.
 * This is a subset of the methods found in the `mountable-auth-methods` helper,
 * which lists all the methods that can be enabled and mounted.
 */

interface MethodData {
  type: string;
  typeDisplay: string;
  description: string;
  tokenPath: string;
  displayNamePath: string | string[];
  formAttributes: string[];
}

export const BASE_LOGIN_METHODS = [
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

const ENTERPRISE_LOGIN_METHODS = [
  {
    type: 'saml',
    typeDisplay: 'SAML',
    description: 'Authenticate using SAML provider.',
    tokenPath: 'client_token',
    displayNamePath: 'display_name',
    formAttributes: ['role'],
  },
];

export const ALL_LOGIN_METHODS = [...BASE_LOGIN_METHODS, ...ENTERPRISE_LOGIN_METHODS];

export const findLoginMethod = (authType: string) =>
  ALL_LOGIN_METHODS.find((m: MethodData) => m.type === authType);
