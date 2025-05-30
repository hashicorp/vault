/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.create('login-rule', {
    name: 'Root namespace default',
    namespace: '',
    default_auth_type: 'userpass',
    backup_auth_types: ['okta', 'token'],
    disable_inheritance: true,
  });
  server.create('login-rule', {
    namespace: 'admin',
    default_auth_type: 'oidc',
    backup_auth_types: ['token'],
  });
  server.create('login-rule', { default_auth_type: 'jwt', backup_auth_types: [] }); // namespace-2
  server.create('login-rule', { default_auth_type: '', backup_auth_types: ['oidc', 'jwt'] }); // namespace-3
  server.create('login-rule', { default_auth_type: '', backup_auth_types: ['token'] }); // namespace-4
}
