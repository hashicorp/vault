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

  hooks.beforeEach(function () {
    this.timestampStub = sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2023-01-13T09:30:15')));
    this.store = this.owner.lookup('service:store');
    this.modelName = 'clients/activity';
    const mockNow = timestamp.now();
    this.startDate = subMonths(mockNow, 6);
    this.endDate = mockNow;
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

  test('it sends without query if no dates provided', async function (assert) {
    assert.expect(1);

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      assert.propEqual(req.queryParams, {});
    });

    this.store.queryRecord(this.modelName, { foo: 'bar' });
  });

  test('it sends without query if no valid dates provided', async function (assert) {
    assert.expect(1);

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      assert.propEqual(req.queryParams, {});
    });

    this.store.queryRecord(this.modelName, { start_time: 'bar' });
  });

  test('it handles empty query gracefully', async function (assert) {
    assert.expect(1);

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      assert.propEqual(req.queryParams, {});
    });

    this.store.queryRecord(this.modelName, {});
  });

  test('it adds the passed namespace to the request header', async function (assert) {
    assert.expect(2);
    const queryParams = {
      start_time: { timestamp: this.startDate.toISOString() },
      end_time: { timestamp: this.endDate.toISOString() },
      // the adapter does not do any more transformations, so it must be called
      // with the combined current + selected namespace
      namespace: 'foobar/baz',
    };
    this.server.get('sys/internal/counters/activity', (schema, req) => {
      assert.propEqual(req.queryParams, {
        start_time: this.startDate.toISOString(),
        end_time: this.endDate.toISOString(),
      });
      assert.strictEqual(req.requestHeaders['X-Vault-Namespace'], 'foobar/baz');
    });

    this.store.queryRecord(this.modelName, queryParams);
  });

  module('exportData', function (hooks) {
    hooks.beforeEach(function () {
      this.adapter = this.store.adapterFor('clients/activity');
    });
    test('it requests with correct params when no query', async function (assert) {
      assert.expect(1);

      this.server.get('sys/internal/counters/activity/export', (schema, req) => {
        assert.propEqual(req.queryParams, { format: 'csv' });
      });

      await this.adapter.exportData();
    });

    test('it requests with correct params when start only', async function (assert) {
      assert.expect(1);

      this.server.get('sys/internal/counters/activity/export', (schema, req) => {
        assert.propEqual(req.queryParams, { format: 'csv', start_time: '2024-04-01T00:00:00.000Z' });
      });

      await this.adapter.exportData({ start_time: '2024-04-01T00:00:00.000Z' });
    });

    test('it requests with correct params when all params', async function (assert) {
      assert.expect(2);

      this.server.get('sys/internal/counters/activity/export', (schema, req) => {
        assert.strictEqual(req.requestHeaders['X-Vault-Namespace'], 'foo/bar');
        assert.propEqual(req.queryParams, {
          format: 'json',
          start_time: '2024-04-01T00:00:00.000Z',
          end_time: '2024-05-31T00:00:00.000Z',
        });
      });

      await this.adapter.exportData({
        start_time: '2024-04-01T00:00:00.000Z',
        end_time: '2024-05-31T00:00:00.000Z',
        format: 'json',
        namespace: 'foo/bar',
      });
    });
  });
});
