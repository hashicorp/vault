/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { deepCopyValue } from 'vault/forms/v2/utils/deep-copy';

module('Unit | forms/v2/utils | deep-copy', function () {
  test('it deep-copies nested arrays and plain objects', function (assert) {
    const original = {
      level1: {
        level2: {
          value: 'original',
          list: [{ id: 1 }, { id: 2 }],
        },
      },
    };

    const cloned = deepCopyValue(original);
    cloned.level1.level2.value = 'changed';
    cloned.level1.level2.list[0].id = 42;

    assert.strictEqual(original.level1.level2.value, 'original', 'nested object value is not mutated');
    assert.strictEqual(original.level1.level2.list[0].id, 1, 'nested array object value is not mutated');
  });

  test('it preserves function references while cloning containers', function (assert) {
    const validatorFn = () => true;
    const original = {
      validations: [
        {
          validator: validatorFn,
        },
      ],
    };

    const cloned = deepCopyValue(original);

    assert.notStrictEqual(cloned, original, 'returns a new top-level object');
    assert.notStrictEqual(cloned.validations, original.validations, 'returns a new nested array');
    assert.strictEqual(cloned.validations[0].validator, validatorFn, 'function reference is preserved');
  });

  test('it handles circular references', function (assert) {
    const original = { name: 'root' };
    original.self = original;

    const cloned = deepCopyValue(original);

    assert.notStrictEqual(cloned, original, 'returns a new object');
    assert.strictEqual(cloned.self, cloned, 'preserves circular structure in clone');
  });
});
