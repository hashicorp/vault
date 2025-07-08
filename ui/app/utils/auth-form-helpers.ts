/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * The web UI only supports logging in with these auth methods.
 * This is a subset of the methods found in the `all-engines-metadata` util,
 * which includes all the methods that can be enabled and mounted.
 */

const BASE_LOGIN_METHODS = ['github', 'jwt', 'ldap', 'oidc', 'okta', 'radius', 'token', 'userpass'];

export const ENTERPRISE_LOGIN_METHODS = ['saml'];

export const supportedTypes = (isEnterprise: boolean) => {
  return isEnterprise ? [...BASE_LOGIN_METHODS, ...ENTERPRISE_LOGIN_METHODS] : [...BASE_LOGIN_METHODS];
};

// this ensures no unexpected params are injected and submitted in the login form
// 'namespace' and 'path' are intentionally omitted because they are handled explicitly
export const POSSIBLE_FIELDS = ['role', 'jwt', 'password', 'token', 'username'];

// maps OIDC provider domain to display name for oidc-jwt auth form
export const DOMAIN_PROVIDER_MAP = {
  'github.com': 'GitHub',
  'gitlab.com': 'GitLab',
  'google.com': 'Google',
  'ping.com': 'Ping Identity',
  'okta.com': 'Okta',
  'auth0.com': 'Auth0',
  'login.microsoftonline.com': 'Azure',
};

export const ERROR_POPUP_FAILED =
  'Your web browser may have blocked or closed a pop-up window. Please check your settings and click "Sign in" to try again.';

export const ERROR_WINDOW_CLOSED = `The provider window was closed before authentication was complete. ${ERROR_POPUP_FAILED}`;

export const ERROR_JWT_LOGIN = 'OIDC login is not configured for this mount';

export const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters.  Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';

export const ERROR_TIMEOUT = 'The authentication request has timed out. Please click "Sign in" to try again.';
