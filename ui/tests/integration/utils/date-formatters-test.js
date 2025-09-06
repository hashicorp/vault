/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { parseAPITimestamp } from 'core/utils/date-formatters';

module('Integration | Util | date formatters utils', function (hooks) {
  setupTest(hooks);

  test('parseAPITimestamp: it returns a date object when a format is not passed', async function (assert) {
    const timestamp = '2025-05-01T00:00:00Z';
    const parsed = parseAPITimestamp(timestamp);
    assert.true(parsed instanceof Date, 'parsed timestamp is a date object');
    assert.strictEqual(parsed.getFullYear(), 2025, 'parsed timestamp is correct year');
    assert.strictEqual(parsed.getMonth(), 4, 'parsed timestamp is correct month (months are 0-indexed)');
    assert.strictEqual(parsed.getDate(), 1, 'parsed timestamp is first of the month');
  });

  test('parseAPITimestamp: it formats midnight timestamps in UTC', async function (assert) {
    const timestamp = '2025-05-01T00:00:00Z';
    const formatted = parseAPITimestamp(timestamp, 'MM dd yyyy');
    // if parseISO was not used, this would return 04 30 2025 if new Date() is invoked in the pacific timezone
    assert.strictEqual(formatted, '05 01 2025', 'it returns the expected year, month and day');
  });

  test('parseAPITimestamp: it returns the original value if timestamp is not a string', async function (assert) {
    const unix = new Date().getTime();
    assert.strictEqual(parseAPITimestamp(unix), unix);
  });

  test('parseAPITimestamp: it formats end of the day timestamps in UTC', async function (assert) {
    const timestamp = '2025-09-30T23:59:59Z';
    const formatted = parseAPITimestamp(timestamp, 'MM dd yyyy');
    assert.strictEqual(formatted, '09 30 2025', 'it formats the date in UTC');
  });
});
