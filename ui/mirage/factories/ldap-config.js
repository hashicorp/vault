/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  backend: 'ldap-test',
  binddn: 'cn=vault,ou=Users,dc=hashicorp,dc=com',
  bindpass: 'pa$$w0rd',
  url: 'ldaps://127.0.0.11',
  password_policy: 'default',
  schema: 'openldap',
  starttls: false,
  insecure_tls: false,
  certificate:
    '-----BEGIN CERTIFICATE-----\nMIIDNTCCAh2gApGgAwIBAgIULNEk+01LpkDeJujfsAgIULNEkAgIULNEckApGgAwIBAg+01LpkDeJuj\n-----END CERTIFICATE-----',
  client_tls_cert:
    '-----BEGIN CERTIFICATE-----\nMIIDNTCCAh2gApGgAwIBAgIULNEk+01LpkDeJujfsAgIULNEkAgIULNEckApGgAwIBAg+01LpkDeJuj\n-----END CERTIFICATE-----',
  client_tls_key: '47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=',
  userdn: 'ou=Users,dc=hashicorp,dc=com',
  userattr: 'cn',
  upndomain: 'vault@hashicorp.com',
  connection_timeout: 90,
  request_timeout: 30,
});
