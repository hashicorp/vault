/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.create('ldap-config', { path: 'kubernetes' });
  server.create('ldap-role', 'static', { name: 'static-role' });
  server.create('ldap-role', 'dynamic', { name: 'dynamic-role' });
  server.create('ldap-library', { name: 'test-library' });
}
