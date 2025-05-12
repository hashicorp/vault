/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.create('login-rule', {
    name: 'Root namespace default',
    namespace: '',
    default_auth_type: 'userpass',
    backup_auth_types: ['okta'],
    disable_inheritance: true,
  });
  server.create('login-rule', {
    namespace: 'admin',
    default_auth_type: 'oidc',
    backup_auth_types: ['token'],
  });
  server.create('login-rule', { default_auth_type: 'jwt', backup_auth_types: [] });
  server.create('login-rule', { default_auth_type: '', backup_auth_types: ['oidc', 'jwt'] });
  server.create('login-rule', { default_auth_type: '', backup_auth_types: ['token'] });
}
