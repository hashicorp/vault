/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { subMonths, fromUnixTime, addMonths } from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import timestamp from 'core/utils/timestamp';

module('Unit | Adapter | clients activity', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    this.timestampStub = sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2023-01-13T09:30:15')));
  });
  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'clients/activity';
    this.startDate = subMonths(this.timestampStub(), 6);
    this.endDate = this.timestampStub();
    this.readableUnix = (unix) => parseAPITimestamp(fromUnixTime(unix).toISOString(), 'MMMM dd yyyy');
  });

  test('it does not format if both params are timestamp strings', async function (assert) {
    assert.expect(1);
    const queryParams = {
      start_time: { timestamp: this.startDate.toISOString() },
      end_time: { timestamp: this.endDate.toISOString() },
    };
    this.server.get('sys/internal/counters/activity', (schema, req) => {
      assert.propEqual(req.queryParams, {
        start_time: this.startDate.toISOString(),
        end_time: this.endDate.toISOString(),
      });
    });

    this.store.queryRecord(this.modelName, queryParams);
  });

  test('it formats start_time if only end_time is a timestamp string', async function (assert) {
    assert.expect(2);
    const twoMonthsAhead = addMonths(this.startDate, 2);
    const month = twoMonthsAhead.getMonth();
    const year = twoMonthsAhead.getFullYear();
    const queryParams = {
      start_time: {
        monthIdx: month,
        year,
      },
      end_time: {
        timestamp: this.endDate.toISOString(),
      },
    };

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      const { start_time, end_time } = req.queryParams;
      const readableStart = this.readableUnix(start_time);
      assert.strictEqual(
        readableStart,
        `September 01 2022`,
        `formatted unix start time is the first of the month: ${readableStart}`
      );
      assert.strictEqual(end_time, this.endDate.toISOString(), 'end time is a timestamp string');
    });
    this.store.queryRecord(this.modelName, queryParams);
  });

  test('it formats end_time only if only start_time is a timestamp string', async function (assert) {
    assert.expect(2);
    const twoMothsAgo = subMonths(this.endDate, 2);
    const endMonth = twoMothsAgo.getMonth();
    const year = twoMothsAgo.getFullYear();
    const queryParams = {
      start_time: {
        timestamp: this.startDate.toISOString(),
      },
      end_time: {
        monthIdx: endMonth,
        year,
      },
    };

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      const { start_time, end_time } = req.queryParams;
      const readableEnd = this.readableUnix(end_time);
      assert.strictEqual(start_time, this.startDate.toISOString(), 'start time is a timestamp string');
      assert.strictEqual(
        readableEnd,
        `November 30 2022`,
        `formatted unix end time is the last day of the month: ${readableEnd}`
      );
    });

    this.store.queryRecord(this.modelName, queryParams);
  });

  test('it formats both params if neither are a timestamp', async function (assert) {
    assert.expect(2);
    const startDate = subMonths(this.startDate, 2);
    const endDate = addMonths(this.endDate, 2);
    const startMonth = startDate.getMonth();
    const startYear = startDate.getFullYear();
    const endMonth = endDate.getMonth();
    const endYear = endDate.getFullYear();
    const queryParams = {
      start_time: {
        monthIdx: startMonth,
        year: startYear,
      },
      end_time: {
        monthIdx: endMonth,
        year: endYear,
      },
    };

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      const { start_time, end_time } = req.queryParams;
      const readableEnd = this.readableUnix(end_time);
      const readableStart = this.readableUnix(start_time);
      assert.strictEqual(
        readableStart,
        `May 01 2022`,
        `formatted unix start time is the first of the month: ${readableStart}`
      );
      assert.strictEqual(
        readableEnd,
        `March 31 2023`,
        `formatted unix end time is the last day of the month: ${readableEnd}`
      );
    });

    this.store.queryRecord(this.modelName, queryParams);
  });

  test('it sends current billing period boolean if provided', async function (assert) {
    assert.expect(1);

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      assert.propEqual(
        req.queryParams,
        { current_billing_period: 'true' },
        'it passes current_billing_period to query record'
      );
    });

    this.store.queryRecord(this.modelName, { current_billing_period: true });
  });
});
