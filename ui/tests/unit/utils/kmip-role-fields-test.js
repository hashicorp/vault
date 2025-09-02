/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import {
  nonOperationFields,
  operationFields,
  operationFieldsWithoutSpecial,
} from 'vault/utils/model-helpers/kmip-role-fields';

module('Unit | Util | kmip role fields', function (hooks) {
  setupTest(hooks);

  [
    {
      name: 'when fields is empty',
      fields: [],
      opFields: [],
      nonOpFields: [],
      opWithoutSpecial: [],
    },
    {
      name: 'when no op fields',
      fields: ['foo', 'bar'],
      opFields: [],
      nonOpFields: ['foo', 'bar'],
      opWithoutSpecial: [],
    },
    {
      name: 'when op fields',
      fields: ['foo', 'bar', 'operationFoo', 'operationBar', 'operationAll'],
      opFields: ['operationFoo', 'operationBar', 'operationAll'],
      nonOpFields: ['foo', 'bar'],
      opWithoutSpecial: ['operationFoo', 'operationBar'],
    },
  ].forEach(({ name, fields, opFields, nonOpFields, opWithoutSpecial }) => {
    test(`${name}`, function (assert) {
      const originalFields = JSON.parse(JSON.stringify(fields));
      assert.deepEqual(operationFields(fields), opFields, 'operation fields correct');
      assert.deepEqual(nonOperationFields(fields), nonOpFields, 'non operation fields');
      assert.deepEqual(
        operationFieldsWithoutSpecial(fields),
        opWithoutSpecial,
        'operation fields without special'
      );
      assert.deepEqual(fields, originalFields, 'does not mutate the original');
    });
  });
});
