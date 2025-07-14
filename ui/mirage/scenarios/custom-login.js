/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.create('login-rule', {
    name: 'root-rule',
    namespace_path: '',
    default_auth_type: 'okta',
    backup_auth_types: ['token'],
    disable_inheritance: false,
  });
  server.create('login-rule', {
    namespace_path: 'admin',
    default_auth_type: 'oidc',
    backup_auth_types: ['token'],
  });
  server.create('login-rule', {
    name: 'ns-rule',
    namespace_path: 'test-ns',
    default_auth_type: 'ldap',
    backup_auth_types: [],
    disable_inheritance: true,
  });
  // generated with defaults set by ui/mirage/factories/login-rule.js
  server.create('login-rule', { default_auth_type: 'jwt', backup_auth_types: [] }); // namespace-2
  server.create('login-rule', { default_auth_type: '', backup_auth_types: ['oidc', 'jwt'] }); // namespace-3
  server.create('login-rule', { default_auth_type: '', backup_auth_types: ['token'] }); // namespace-4
}
