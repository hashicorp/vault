/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { formatRFC3339, isSameDay, isSameMonth, isSameYear } from 'date-fns';
import { parseAPITimestamp, formatChartDate } from 'core/utils/date-formatters';

module('Integration | Util | date formatters utils', function (hooks) {
  setupTest(hooks);

  test('parseAPITimestamp: parses API timestamp string irrespective of timezone', async function (assert) {
    assert.expect(6);
    const DATE = new Date('2012-06-10T15:30:45');
    const API_TIMESTAMP = formatRFC3339(DATE).split('T')[0].concat('T00:00:00Z');
    const UNIX_TIME = DATE.getTime();
    assert.strictEqual(
      parseAPITimestamp(UNIX_TIME),
      undefined,
      'it returns undefined if timestamp is not a string'
    );

    const parsedTimestamp = parseAPITimestamp(API_TIMESTAMP);

    assert.true(parsedTimestamp instanceof Date, 'parsed timestamp is a date object');
    assert.true(isSameYear(parsedTimestamp, DATE), 'parsed timestamp is correct year');
    assert.true(isSameMonth(parsedTimestamp, DATE), 'parsed timestamp is correct month');
    assert.true(isSameDay(parsedTimestamp, DATE), 'parsed timestamp is correct day');

    const formattedTimestamp = parseAPITimestamp(API_TIMESTAMP, 'MM yyyy');
    assert.strictEqual(formattedTimestamp, '06 2012', 'it formats the date');
  });

  test('formatChartDate: expand chart date to full month and year', async function (assert) {
    assert.expect(1);
    const chartDate = '03/21';
    assert.strictEqual(formatChartDate(chartDate), 'March 2021', 'it re-formats the date');
  });
});
