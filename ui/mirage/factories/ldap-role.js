/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  name: (i) => `role-${i}`,

  // static props
  static: trait({
    dn: 'cn=hashicorp,ou=Users,dc=hashicorp,dc=com',
    rotation_period: 10,
    username: 'hashicorp',
    type: 'static',
  }),

  // dynamic props
  dynamic: trait({
    creation_ldif: `dn: cn={{.Username}},ou=users,dc=learn,dc=example
    objectClass: person
    objectClass: top
    cn: learn
    sn: {{.Password | utf16le | base64}}
    memberOf: cn=dev,ou=groups,dc=learn,dc=example
    userPassword: {{.Password}}
    `,
    deletion_ldif: `dn: cn={{.Username}},ou=users,dc=learn,dc=example
    changetype: delete
    `,
    rollback_ldif: `dn: cn={{.Username}},ou=users,dc=learn,dc=example
    changetype: delete
    `,
    username_template: '{{.DisplayName}}_{{.RoleName}}',
    default_ttl: 3600,
    max_ttl: 86400,
    type: 'dynamic',
  }),
});
