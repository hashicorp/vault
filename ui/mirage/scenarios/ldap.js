/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.create('ldap-config', { path: 'kubernetes', backend: 'ldap-test' });
  server.create('ldap-role', 'static', { name: 'static-role' });
  server.create('ldap-role', 'dynamic', { name: 'dynamic-role' });
  // hierarchical roles
  server.create('ldap-role', 'static', { name: 'admin/child-static-role' });
  server.create('ldap-role', 'dynamic', { name: 'admin/child-dynamic-role' });
  // use same name for both role types to test edge cases
  server.create('ldap-role', 'static', { name: 'my-role' });
  server.create('ldap-role', 'dynamic', { name: 'my-role' });
  server.create('ldap-library', { name: 'test-library' });
  // mirage handler is hardcoded to accommodate hierarchical paths starting with 'admin/'
  server.create('ldap-library', { name: 'admin/test-library' });
  server.create('ldap-account-status', {
    id: 'bob.johnson',
    account: 'bob.johnson',
    available: false,
    borrower_client_token: '8b80c305eb3a7dbd161ef98f10ea60a116ce0910',
  });
  server.create('ldap-account-status', {
    id: 'mary.smith',
    account: 'mary.smith',
    available: true,
  });
}
