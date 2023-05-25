/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { duration } from 'core/helpers/format-duration';

module('Integration | Helper | format-duration', function (hooks) {
  setupRenderingTest(hooks);

  test('it formats seconds', async function (assert) {
    assert.strictEqual(duration([3606]), '1 hour 6 seconds');
  });

  test('it format seconds and days', async function (assert) {
    assert.strictEqual(duration([93606000]), '2 years 11 months 18 days 9 hours 40 minutes');
  });

  test('it returns the integer 0', async function (assert) {
    assert.strictEqual(duration([0]), 0);
  });

  test('it returns string inputs', async function (assert) {
    this.set('number', 'arg');
    assert.strictEqual(duration(['0']), '0');
    assert.strictEqual(duration(['arg']), 'arg');
    assert.strictEqual(duration(['1245']), '1245');
  });
});
