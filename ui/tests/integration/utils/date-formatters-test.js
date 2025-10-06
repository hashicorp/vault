/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { buildISOTimestamp, isSameMonthUTC, parseAPITimestamp } from 'core/utils/date-formatters';

module('Integration | Util | date formatters utils', function (hooks) {
  setupTest(hooks);

  test('parseAPITimestamp: it returns a date object when a format is not passed', async function (assert) {
    const timestamp = '2025-05-01T00:00:00Z';
    const parsed = parseAPITimestamp(timestamp);
    assert.true(parsed instanceof Date, 'parsed timestamp is a date object');
    assert.strictEqual(parsed.getUTCFullYear(), 2025, 'parsed timestamp is correct year');
    assert.strictEqual(parsed.getUTCMonth(), 4, 'parsed timestamp is correct month (months are 0-indexed)');
    assert.strictEqual(parsed.getUTCDate(), 1, 'parsed timestamp is first of the month');
    assert.strictEqual(
      parsed.toISOString().replace('.000', ''),
      timestamp,
      'parsed ISO is the same date (in UTC)'
    );
  });

  test('parseAPITimestamp: it formats midnight timestamps in UTC', async function (assert) {
    const timestamp = '2025-05-01T00:00:00Z';
    const formatted = parseAPITimestamp(timestamp, 'MM dd yyyy');
    // if parseISO was not used, this would return 04 30 2025 if new Date() is invoked in the pacific timezone
    assert.strictEqual(formatted, '05 01 2025', 'it returns the expected year, month and day');
  });

  test('parseAPITimestamp: it formats end of the day timestamps in UTC', async function (assert) {
    const timestamp = '2025-09-30T23:59:59Z';
    const formatted = parseAPITimestamp(timestamp, 'MM dd yyyy');
    assert.strictEqual(formatted, '09 30 2025', 'it formats the date in UTC');
  });

  test('parseAPITimestamp: it returns null for invalid timestamps', function (assert) {
    const unix = new Date().getTime();
    assert.strictEqual(parseAPITimestamp(unix), null, 'it returns null for unix arg');
    assert.strictEqual(parseAPITimestamp(null), null, 'it returns null for null arg');
    assert.strictEqual(parseAPITimestamp(undefined), null, 'it returns null for undefined arg');
    assert.strictEqual(parseAPITimestamp(''), null, 'it returns null for an empty string arg');
    assert.strictEqual(parseAPITimestamp('invalid'), null, 'it returns null for an invalid string');
  });

  test('parseAPITimestamp: it handles future dates to prep for the next y2k', function (assert) {
    const futureDate = '9999-12-31T23:59:59Z';
    const parsed = parseAPITimestamp(futureDate);
    assert.true(parsed instanceof Date, 'parsed future date is a date object');
    assert.strictEqual(parsed.getUTCFullYear(), 9999, 'parsed future date has correct year');
    assert.strictEqual(parsed.getUTCMonth(), 11, 'parsed future date has correct month');
    assert.strictEqual(parsed.getUTCDate(), 31, 'parsed future date has correct day');
  });

  test('buildISOTimestamp: it formats an ISO timestamp for the start of the month', async function (assert) {
    const timestamp = buildISOTimestamp({ monthIdx: 0, year: 2025, isEndDate: false });
    assert.strictEqual(
      timestamp,
      '2025-01-01T00:00:00Z',
      'it returns an ISO string for the first of the month at midnight'
    );
  });

  test('buildISOTimestamp: it formats an ISO timestamp for the end of the month', async function (assert) {
    const timestamp = buildISOTimestamp({ monthIdx: 0, year: 2025, isEndDate: true });
    assert.strictEqual(timestamp, '2025-01-31T23:59:59Z', 'ISO string is for the last day and hour');
  });

  test('buildISOTimestamp: it formats an ISO timestamp for leap years', async function (assert) {
    const timestamp = buildISOTimestamp({ monthIdx: 1, year: 2024, isEndDate: true });
    assert.strictEqual(timestamp, '2024-02-29T23:59:59Z');
  });

  test('isSameMonthUTC: it returns true for timestamps in the same month', function (assert) {
    const timestampA = '2025-07-01T00:00:00Z';
    const timestampB = '2025-07-15T23:48:09Z';
    assert.true(isSameMonthUTC(timestampA, timestampB));
  });

  test('isSameMonthUTC: it returns false for timestamps in different months', function (assert) {
    const timestampA = '2025-09-01T00:00:00Z';
    const timestampB = '2025-07-15T23:48:09Z';
    assert.false(isSameMonthUTC(timestampA, timestampB));
  });

  test('isSameMonthUTC: it returns false for timestamps in different years', function (assert) {
    const timestampA = '2025-12-31T23:59:59Z';
    const timestampB = '2026-12-01T00:00:00Z';
    assert.false(isSameMonthUTC(timestampA, timestampB));
  });

  test('isSameMonthUTC: it returns true when both timestamps are the same', function (assert) {
    const timestamp = '2025-07-15T23:48:09Z';
    assert.true(isSameMonthUTC(timestamp, timestamp));
  });

  test('isSameMonthUTC: it returns true for timestamps on the first and last day of the month', function (assert) {
    const start = '2025-07-01T00:00:00Z';
    const end = '2025-07-31T23:59:59Z';
    assert.true(isSameMonthUTC(start, end));
  });

  test('isSameMonthUTC: it returns false for timestamps a second apart on different days', function (assert) {
    const endJuly = '2025-07-31T23:59:59Z';
    const startAugust = '2025-08-01T00:00:00Z';
    assert.false(isSameMonthUTC(endJuly, startAugust));
  });

  test('isSameMonthUTC: it returns true if passed a timestamp with a timezone offset', function (assert) {
    const utc = '2025-07-01T00:00:00Z';
    const localTimezone = '2025-07-01T07:00:00+07:00'; // same time in UTC+7
    assert.true(isSameMonthUTC(utc, localTimezone));
  });

  test('isSameMonthUTC: it returns false for invalid inputs', function (assert) {
    assert.false(isSameMonthUTC(null, '2025-07-01T00:00:00Z'), 'null input returns false');
    assert.false(isSameMonthUTC('2025-07-01T00:00:00Z', undefined), 'undefined input returns false');
    assert.false(isSameMonthUTC(12345, '2025-07-01T00:00:00Z'), 'non-string input returns false');
  });
});
