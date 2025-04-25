/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import FormField from 'vault/utils/forms/field';
import { module, test } from 'qunit';

module('Unit | Utility | forms | field', function (hooks) {
  hooks.beforeEach(function () {
    this.options = {
      label: 'Test field',
      subText: 'This is a form field',
      fieldValue: 'custom.path',
      editType: 'textarea',
    };
  });

  test('it should set name, type and field options', async function (assert) {
    const field = new FormField('test', 'string', this.options);
    assert.strictEqual(field.name, 'test', 'Name is set correctly');
    assert.strictEqual(field.type, 'string', 'Type is set correctly');
    assert.deepEqual(field.options, this.options, 'Options are set correctly');
  });

  test('it should default field value', async function (assert) {
    this.options.fieldValue = undefined;
    const field = new FormField('test', 'string', this.options);
    assert.strictEqual(field.options.fieldValue, 'test', 'Default field value is set correctly');
  });
});
