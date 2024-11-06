/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.create('ldap-config', { path: 'kubernetes' });
  server.create('ldap-role', 'static', { name: 'static-role' });
  server.create('ldap-role', 'dynamic', { name: 'dynamic-role' });
  // hierarchical roles
  server.create('ldap-role', 'static', { name: 'admin/child-static-role' });
  server.create('ldap-role', 'dynamic', { name: 'admin/child-dynamic-role' });
  // use same name for both role types to test edge cases
  server.create('ldap-role', 'static', { name: 'my-role' });
  server.create('ldap-role', 'dynamic', { name: 'my-role' });
  server.create('ldap-library', { name: 'test-library' });
}
