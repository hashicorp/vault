/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Transform | comma string', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.transform = this.owner.lookup('transform:comma-string');
  });

  test('it serializes correctly for API', function (assert) {
    const serialized = this.transform.serialize('one,two,three');
    assert.propEqual(serialized, ['one', 'two', 'three'], 'it serializes from string to array');
    assert.propEqual(
      this.transform.serialize(['not a string']),
      ['not a string'],
      'it returns original value if not a string'
    );
    assert.propEqual(
      this.transform.serialize('no commas'),
      ['no commas'],
      'it splits a string without commas'
    );
  });

  test('it deserializes correctly from API', function (assert) {
    const deserialized = this.transform.deserialize(['one', 'two', 'three']);
    assert.strictEqual(deserialized, 'one,two,three', 'it deserializes from array to string');
    assert.strictEqual(
      this.transform.deserialize('not an array'),
      'not an array',
      'it returns original value if not an array'
    );
  });
});
