/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import { run } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { normalizeModelName, keyForCache } from 'vault/services/store';
import clamp from 'vault/utils/clamp';
import config from 'vault/config/environment';

const { DEFAULT_PAGE_SIZE } = config.APP;

module('Unit | Service | store', function (hooks) {
  setupTest(hooks);

  test('normalizeModelName', function (assert) {
    assert.strictEqual(normalizeModelName('oneThing'), 'one-thing', 'dasherizes modelName');
  });

  test('keyForCache', function (assert) {
    const query = { id: 1 };
    const queryWithSize = { id: 1, size: 1 };
    assert.deepEqual(keyForCache(query), JSON.stringify(query), 'generated the correct cache key');
    assert.deepEqual(keyForCache(queryWithSize), JSON.stringify(query), 'excludes size from query cache');
  });

  test('clamp', function (assert) {
    assert.strictEqual(clamp('foo', 0, 100), 0, 'returns the min if passed a non-number');
    assert.strictEqual(clamp(0, 1, 100), 1, 'returns the min when passed number is less than the min');
    assert.strictEqual(clamp(200, 1, 100), 100, 'returns the max passed number is greater than the max');
    assert.strictEqual(clamp(50, 1, 100), 50, 'returns the passed number when it is in range');
  });

  test('store.storeDataset', function (assert) {
    const arr = ['one', 'two'];
    const store = this.owner.lookup('service:store');
    const query = { id: 1 };
    store.storeDataset('data', query, {}, arr);

    assert.deepEqual(store.getDataset('data', query).dataset, arr, 'it stores the array as .dataset');
    assert.deepEqual(store.getDataset('data', query).response, {}, 'it stores the response as .response');
    assert.ok(store.get('lazyCaches').has('data'), 'it stores model map');
    assert.ok(store.get('lazyCaches').get('data').has(keyForCache(query)), 'it stores data on the model map');
  });

  test('store.clearDataset with a prefix', function (assert) {
    const store = this.owner.lookup('service:store');
    const arr = ['one', 'two'];
    const arr2 = ['one', 'two', 'three', 'four'];
    store.storeDataset('data', { id: 1 }, {}, arr);
    store.storeDataset('transit-key', { id: 2 }, {}, arr2);
    assert.strictEqual(store.get('lazyCaches').size, 2, 'it stores both keys');

    store.clearDataset('transit-key');
    assert.strictEqual(store.get('lazyCaches').size, 1, 'deletes one key');
    assert.notOk(store.get('lazyCaches').has('transit-key'), 'cache is no longer stored');
  });

  test('store.clearAllDatasets', function (assert) {
    const store = this.owner.lookup('service:store');
    const arr = ['one', 'two'];
    const arr2 = ['one', 'two', 'three', 'four'];
    store.storeDataset('data', { id: 1 }, {}, arr);
    store.storeDataset('transit-key', { id: 2 }, {}, arr2);
    assert.strictEqual(store.get('lazyCaches').size, 2, 'it stores both keys');

    store.clearAllDatasets();
    assert.strictEqual(store.get('lazyCaches').size, 0, 'deletes all of the keys');
    assert.notOk(store.get('lazyCaches').has('transit-key'), 'first cache key is no longer stored');
    assert.notOk(store.get('lazyCaches').has('data'), 'second cache key is no longer stored');
  });

  test('store.getDataset', function (assert) {
    const arr = ['one', 'two'];
    const store = this.owner.lookup('service:store');
    store.storeDataset('data', { id: 1 }, {}, arr);

    assert.deepEqual(store.getDataset('data', { id: 1 }), { response: {}, dataset: arr });
  });

  test('store.constructResponse', function (assert) {
    const arr = ['one', 'two', 'three', 'fifteen', 'twelve'];
    const store = this.owner.lookup('service:store');
    store.storeDataset('data', { id: 1 }, {}, arr);

    assert.deepEqual(
      store.constructResponse('data', { id: 1, pageFilter: 't', page: 1, size: 3, responsePath: 'data' }),
      {
        data: ['two', 'three', 'fifteen'],
        meta: {
          currentPage: 1,
          lastPage: 2,
          nextPage: 2,
          prevPage: 1,
          total: 5,
          filteredTotal: 4,
          pageSize: 3,
        },
      },
      'it returns filtered results'
    );
  });

  test('store.fetchPage', async function (assert) {
    const keys = ['zero', 'one', 'two', 'three', 'four', 'five', 'six'];
    const data = {
      data: {
        keys,
      },
    };
    const store = this.owner.lookup('service:store');
    const pageSize = 2;
    const query = {
      size: pageSize,
      page: 1,
      responsePath: 'data.keys',
    };
    store.storeDataset('transit-key', query, data, keys);

    let result;
    result = await store.fetchPage('transit-key', query);
    assert.strictEqual(result.get('length'), pageSize, 'returns the correct number of items');
    assert.deepEqual(result.mapBy('id'), keys.slice(0, pageSize), 'returns the first page of items');
    assert.deepEqual(
      result.get('meta'),
      {
        nextPage: 2,
        prevPage: 1,
        currentPage: 1,
        lastPage: 4,
        total: 7,
        filteredTotal: 7,
        pageSize: 2,
      },
      'returns correct meta values'
    );

    result = await store.fetchPage('transit-key', {
      size: pageSize,
      page: 3,
      responsePath: 'data.keys',
    });
    const pageThreeEnd = 3 * pageSize;
    const pageThreeStart = pageThreeEnd - pageSize;
    assert.deepEqual(
      result.mapBy('id'),
      keys.slice(pageThreeStart, pageThreeEnd),
      'returns the third page of items'
    );

    result = await store.fetchPage('transit-key', {
      size: pageSize,
      page: 99,
      responsePath: 'data.keys',
    });

    assert.deepEqual(
      result.mapBy('id'),
      keys.slice(keys.length - 1),
      'returns the last page when the page value is beyond the of bounds'
    );

    result = await store.fetchPage('transit-key', {
      size: pageSize,
      page: 0,
      responsePath: 'data.keys',
    });
    assert.deepEqual(
      result.mapBy('id'),
      keys.slice(0, pageSize),
      'returns the first page when page value is under the bounds'
    );
  });

  test('store.lazyPaginatedQuery', function (assert) {
    const response = {
      data: ['foo'],
    };
    let queryArgs;
    const store = this.owner.factoryFor('service:store').create({
      adapterFor() {
        return {
          query(store, modelName, query) {
            queryArgs = query;
            return resolve(response);
          },
        };
      },
      fetchPage() {},
    });

    const query = { page: 1, size: 1, responsePath: 'data' };
    run(function () {
      store.lazyPaginatedQuery('transit-key', query);
    });
    assert.deepEqual(
      store.getDataset('transit-key', query),
      { response: { data: null }, dataset: ['foo'] },
      'stores returned dataset'
    );

    run(function () {
      store.lazyPaginatedQuery('secret', { page: 1, responsePath: 'data' });
    });
    assert.strictEqual(queryArgs.size, DEFAULT_PAGE_SIZE, 'calls query with DEFAULT_PAGE_SIZE');

    assert.throws(
      () => {
        store.lazyPaginatedQuery('transit-key', {});
      },
      /responsePath is required/,
      'requires responsePath'
    );
    assert.throws(
      () => {
        store.lazyPaginatedQuery('transit-key', { responsePath: 'foo' });
      },
      /page is required/,
      'requires page'
    );
  });
});
