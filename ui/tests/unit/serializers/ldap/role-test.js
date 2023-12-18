/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Serializer | ldap/role', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    const store = this.owner.lookup('service:store');
    this.model = store.createRecord('ldap/role', {
      backend: 'ldap',
      name: 'test-role',
      dn: 'cn=hashicorp,ou=Users,dc=hashicorp,dc=com',
      rotation_period: '24h',
      username: 'hashicorp',
      creation_ldif: 'foo',
      deletion_ldif: 'bar',
      rollback_ldif: 'baz',
      username_template: 'default',
      default_ttl: '1h',
      max_ttl: '24h',
    });
  });

  test('it should serialize attributes based on type', async function (assert) {
    assert.expect(11);

    const serializeAndAssert = (type) => {
      this.model.type = type;
      const payload = this.model.serialize();
      // intentionally not using fieldsForType from model to detect any drift
      const fieldsForType = {
        static: ['username', 'dn', 'rotation_period'],
        dynamic: [
          'default_ttl',
          'max_ttl',
          'username_template',
          'creation_ldif',
          'deletion_ldif',
          'rollback_ldif',
        ],
      }[type];

      assert.strictEqual(
        Object.keys(payload).length,
        fieldsForType.length,
        `Correct number of keys exist in serialized payload for ${type} role type`
      );
      Object.keys(payload).forEach((key) => {
        assert.true(
          fieldsForType.includes(key),
          `${key} property exists in serialized payload for ${type} role type`
        );
      });
    };

    serializeAndAssert('static');
    serializeAndAssert('dynamic');
  });
});
