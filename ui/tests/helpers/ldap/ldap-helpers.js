/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { visit, currentURL } from '@ember/test-helpers';

export const createSecretsEngine = (store) => {
  store.pushPayload('secret-engine', {
    modelName: 'secret-engine',
    data: {
      accessor: 'ldap_7e838627',
      path: 'ldap-test/',
      type: 'ldap',
    },
  });
  return store.peekRecord('secret-engine', 'ldap-test');
};

export const generateBreadcrumbs = (backend, childRoute) => {
  const breadcrumbs = [{ label: 'Secrets', route: 'secrets', linkExternal: true }];
  const root = { label: backend };
  if (childRoute) {
    root.route = 'overview';
    breadcrumbs.push({ label: childRoute });
  }
  breadcrumbs.splice(1, 0, root);
  return breadcrumbs;
};

const baseURL = (backend) => `/vault/secrets/${backend}/ldap/`;
const stripLeadingSlash = (uri) => (uri.charAt(0) === '/' ? uri.slice(1) : uri);

export const isURL = (uri, backend = 'ldap-test') => {
  return currentURL() === `${baseURL(backend)}${stripLeadingSlash(uri)}`;
};

export const assertURL = (assert, backend, path) => {
  assert.strictEqual(currentURL(), baseURL(backend) + path, `url is ${path}`);
};

export const visitURL = (uri, backend = 'ldap-test') => {
  return visit(`${baseURL(backend)}${stripLeadingSlash(uri)}`);
};
