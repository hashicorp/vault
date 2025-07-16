/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { duration } from 'core/helpers/format-duration';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Helper | format-duration', function (hooks) {
  setupRenderingTest(hooks);

  test('it formats-duration in template view', async function (assert) {
    await render(hbs`<p data-test-format-duration>Date: {{format-duration 3606 }}</p>`);

    assert
      .dom('[data-test-format-duration]')
      .includesText('1 hour 6 seconds', 'it renders the duration in hours and seconds');
  });

  test('it formats seconds', async function (assert) {
    assert.strictEqual(duration([3606]), '1 hour 6 seconds');
  });

  test('it format seconds and days', async function (assert) {
    assert.strictEqual(duration([93606000]), '2 years 11 months 18 days 9 hours 40 minutes');
  });

  test('it returns the integer 0', async function (assert) {
    assert.strictEqual(duration([0]), 0);
  });

  test('it returns plain or non-parsable string inputs', async function (assert) {
    assert.strictEqual(duration(['0']), '0 seconds'); // assume seconds for '0' string values only
    assert.strictEqual(duration(['arg']), 'arg');
    assert.strictEqual(duration(['1245']), '1245');
    assert.strictEqual(duration(['11y']), '11y');
  });

  test('it formats duration string inputs', async function (assert) {
    assert.strictEqual(duration(['0s']), '0 seconds');
    assert.strictEqual(duration(['5s']), '5 seconds');
    assert.strictEqual(duration(['545h']), '22 days 17 hours');
    assert.strictEqual(duration(['8h']), '8 hours');
    assert.strictEqual(duration(['3m']), '3 minutes');
    assert.strictEqual(duration(['10d']), '10 days');
  });
});
