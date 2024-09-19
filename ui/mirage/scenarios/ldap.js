/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.create('ldap-config', { path: 'kubernetes', backend: 'ldap-test' });
  server.create('ldap-role', 'static', { name: 'static-role' });
  server.create('ldap-role', 'dynamic', { name: 'dynamic-role' });
  server.create('ldap-library', { name: 'test-library' });
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
