/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | cluster', function (hooks) {
  setupTest(hooks);

  test('cluster api urls', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });
    adapter.health();
    assert.strictEqual(url, '/v1/sys/health', 'health url OK');
    assert.deepEqual(
      {
        standbycode: 200,
        sealedcode: 200,
        uninitcode: 200,
        drsecondarycode: 200,
        performancestandbycode: 200,
      },
      options.data,
      'health data params OK'
    );
    assert.strictEqual(method, 'GET', 'health method OK');

    adapter.sealStatus();
    assert.strictEqual(url, '/v1/sys/seal-status', 'health url OK');
    assert.strictEqual(method, 'GET', 'seal-status method OK');

    const data = { someData: 1 };
    adapter.unseal(data);
    assert.strictEqual(url, '/v1/sys/unseal', 'unseal url OK');
    assert.strictEqual(method, 'PUT', 'unseal method OK');
    assert.deepEqual({ data, unauthenticated: true }, options, 'unseal options OK');

    adapter.initCluster(data);
    assert.strictEqual(url, '/v1/sys/init', 'init url OK');
    assert.strictEqual(method, 'PUT', 'init method OK');
    assert.deepEqual({ data, unauthenticated: true }, options, 'init options OK');
  });

  test('cluster replication api urls', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    adapter.replicationStatus();
    assert.strictEqual(url, '/v1/sys/replication/status', 'replication:status url OK');
    assert.strictEqual(method, 'GET', 'replication:status method OK');
    assert.deepEqual({ unauthenticated: true }, options, 'replication:status options OK');

    adapter.replicationAction('recover', 'dr');
    assert.strictEqual(url, '/v1/sys/replication/recover', 'replication: recover url OK');
    assert.strictEqual(method, 'POST', 'replication:recover method OK');

    adapter.replicationAction('reindex', 'dr');
    assert.strictEqual(url, '/v1/sys/replication/reindex', 'replication: reindex url OK');
    assert.strictEqual(method, 'POST', 'replication:reindex method OK');

    adapter.replicationAction('enable', 'dr', 'primary');
    assert.strictEqual(url, '/v1/sys/replication/dr/primary/enable', 'replication:dr primary:enable url OK');
    assert.strictEqual(method, 'POST', 'replication:primary:enable method OK');
    adapter.replicationAction('enable', 'performance', 'primary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/primary/enable',
      'replication:performance primary:enable url OK'
    );

    adapter.replicationAction('enable', 'dr', 'secondary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/enable',
      'replication:dr secondary:enable url OK'
    );
    assert.strictEqual(method, 'POST', 'replication:secondary:enable method OK');
    adapter.replicationAction('enable', 'performance', 'secondary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/secondary/enable',
      'replication:performance secondary:enable url OK'
    );

    adapter.replicationAction('disable', 'dr', 'primary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/primary/disable',
      'replication:dr primary:disable url OK'
    );
    assert.strictEqual(method, 'POST', 'replication:primary:disable method OK');
    adapter.replicationAction('disable', 'performance', 'primary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/primary/disable',
      'replication:performance primary:disable url OK'
    );

    adapter.replicationAction('disable', 'dr', 'secondary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/disable',
      'replication: drsecondary:disable url OK'
    );
    assert.strictEqual(method, 'POST', 'replication:secondary:disable method OK');
    adapter.replicationAction('disable', 'performance', 'secondary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/secondary/disable',
      'replication: performance:disable url OK'
    );

    adapter.replicationAction('demote', 'dr', 'primary');
    assert.strictEqual(url, '/v1/sys/replication/dr/primary/demote', 'replication: dr primary:demote url OK');
    assert.strictEqual(method, 'POST', 'replication:primary:demote method OK');
    adapter.replicationAction('demote', 'performance', 'primary');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/primary/demote',
      'replication: performance primary:demote url OK'
    );

    adapter.replicationAction('promote', 'performance', 'secondary');
    assert.strictEqual(method, 'POST', 'replication:secondary:promote method OK');
    assert.strictEqual(
      url,
      '/v1/sys/replication/performance/secondary/promote',
      'replication:performance secondary:promote url OK'
    );

    adapter.replicationDrPromote();
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/promote',
      'replication:dr secondary:promote url OK'
    );
    assert.strictEqual(method, 'PUT', 'replication:dr secondary:promote method OK');
    adapter.replicationDrPromote({}, { checkStatus: true });
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/promote',
      'replication:dr secondary:promote url OK'
    );
    assert.strictEqual(method, 'GET', 'replication:dr secondary:promote method OK');
  });

  test('cluster generateDrOperationToken', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:cluster').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    // Generate token progress
    adapter.generateDrOperationToken({ key: 'foo' }, {});
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/generate-operation-token/update',
      'progress url correct'
    );
    assert.strictEqual(method, 'POST', 'progress method OK');
    assert.deepEqual({ data: { key: 'foo' }, unauthenticated: true }, options, 'progress payload OK');

    // CheckStatus / Read generation progress
    adapter.generateDrOperationToken({ key: 'foo' }, { checkStatus: true });
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/generate-operation-token/attempt',
      'checkStatus url correct'
    );
    assert.strictEqual(method, 'GET', 'checkStatus method OK');
    assert.deepEqual({ data: { key: 'foo' }, unauthenticated: true }, options, 'checkStatus options OK');

    // Cancel generation
    adapter.generateDrOperationToken({}, { cancel: true });
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/generate-operation-token/attempt',
      'Cancel url correct'
    );
    assert.strictEqual(method, 'DELETE', 'cancel method OK');
    assert.deepEqual({ data: {}, unauthenticated: true }, options, 'cancel options OK');

    // pgp_key
    adapter.generateDrOperationToken({ pgp_key: 'yes' }, {});
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/generate-operation-token/attempt',
      'pgp_key url correct'
    );
    assert.strictEqual(method, 'POST', 'method ok when pgp_key on data');
    assert.deepEqual({ data: { pgp_key: 'yes' }, unauthenticated: true }, options, 'pgp_key options OK');

    // data.attempt
    adapter.generateDrOperationToken({ attempt: true }, {});
    assert.strictEqual(
      url,
      '/v1/sys/replication/dr/secondary/generate-operation-token/attempt',
      'data.attempt url correct'
    );
    assert.strictEqual(method, 'POST', 'data.attempt method OK');
    assert.deepEqual({ data: { attempt: true }, unauthenticated: true }, options, 'data.attempt options OK');
  });
});
