/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { subMonths, fromUnixTime } from 'date-fns';
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
      start_time: this.startDate.toISOString(),
      end_time: this.endDate.toISOString(),
    };
    this.server.get('sys/internal/counters/activity', (schema, req) => {
      assert.propEqual(req.queryParams, {
        start_time: this.startDate.toISOString(),
        end_time: this.endDate.toISOString(),
      });
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

    this.store.queryRecord(this.modelName, { start_time: 'bar', end_time: 'baz' });
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
      start_time: this.startDate.toISOString(),
      end_time: this.endDate.toISOString(),
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
