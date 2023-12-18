/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Serializer | ldap/library', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it should normalize and serialize disable_check_in_enforcement value', async function (assert) {
    assert.expect(4);

    const model = this.store.createRecord('ldap/library', {
      backend: 'ldap-test',
      name: 'test-library',
    });
    const cases = [
      { value: false, transformed: 'Enabled' },
      { value: true, transformed: 'Disabled' },
    ];

    cases.forEach(({ value, transformed }) => {
      const normalized = this.store.normalize('ldap/library', { disable_check_in_enforcement: value });
      assert.strictEqual(
        normalized.data.attributes.disable_check_in_enforcement,
        transformed,
        `Normalizes ${value} value to ${transformed}`
      );
      model.disable_check_in_enforcement = transformed;
      const { disable_check_in_enforcement } = model.serialize();
      assert.strictEqual(disable_check_in_enforcement, value, `Serializes ${transformed} value to ${value}`);
    });
  });
});
