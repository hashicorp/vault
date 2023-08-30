/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  // static props
  static: trait({
    last_vault_rotation: '2023-07-31T10:32:49.744033-06:00',
    password: 'fQ428N5JVeB2MbINwBCIbPh2ffhkJP0jZT3SfopZO0xRmbOaKRa6bwtAw3d2m4DR',
    username: 'foobar',
    rotation_period: 86400,
    ttl: 71365,
    type: 'static',
  }),

  // dynamic props
  dynamic: trait({
    distinguished_names: [
      'cn=v_userpass-test_dynamic-role_mrx3r26XIj_1690836430,ou=users,dc=learn,dc=example',
    ],
    username: 'v_userpass-test_dynamic-role_mrx3r26XIj_1690836430',
    password: 'YE2qe1vpiBtEvjCSr7BmI0NhSPPmrizngNYxa3lEebMFvAussxHf3PWfDVJPxXj1',
    lease_id: 'ldap/creds/dynamic-role/SZN8HcuieCbdDobD7jTb6V9X',
    lease_duration: 3600,
    renewable: true,
    type: 'dynamic',
  }),
});
