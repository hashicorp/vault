/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { keyForCache } from 'vault/services/pagination';
import { dasherize } from '@ember/string';
import clamp from 'vault/utils/clamp';
import config from 'vault/config/environment';
import Sinon from 'sinon';

const { DEFAULT_PAGE_SIZE } = config.APP;

module('Unit | Service | pagination', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.pagination = this.owner.lookup('service:pagination');
    this.store = this.owner.lookup('service:store');
  });

  test('pagination.setLazyCacheForModel', function (assert) {
    const modelName = 'someModel';
    const key = {
      id: '',
      backend: 'database',
      responsePath: 'data.keys',
      page: 1,
      pageFilter: null,
      size: 15,
    };
    const value = {
      response: {
        request_id: '1eb6473c-8df0-924e-1c8d-e016a6420aee',
        lease_id: '',
        renewable: false,
        lease_duration: 0,
        data: {
          keys: null,
        },
        wrap_info: null,
        warnings: null,
        auth: null,
        mount_type: 'database',
        backend: 'database',
      },
      dataset: ['connection', 'connection2'],
    };
    this.pagination.setLazyCacheForModel(modelName, key, value);
    const cacheEntry = this.pagination.lazyCaches.get(dasherize(modelName));
    const actual = Object.fromEntries(cacheEntry); // convert from Map to Object for assertion
    const expected = { '{"backend":"database","id":""}': value };
    assert.propEqual(actual, expected, 'model name is dasherized and can be retrieved from lazyCache');
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

  test('pagination.storeDataset', function (assert) {
    const arr = ['one', 'two'];
    const query = { id: 1 };
    this.pagination.storeDataset('data', query, {}, arr);

    assert.deepEqual(
      this.pagination.getDataset('data', query).dataset,
      arr,
      'it stores the array as .dataset'
    );
    assert.deepEqual(
      this.pagination.getDataset('data', query).response,
      {},
      'it stores the response as .response'
    );
    assert.ok(this.pagination.get('lazyCaches').has('data'), 'it stores model map');
    assert.ok(
      this.pagination.get('lazyCaches').get('data').has(keyForCache(query)),
      'it stores data on the model map'
    );
  });

  test('pagination.clearDataset with a prefix', function (assert) {
    const arr = ['one', 'two'];
    const arr2 = ['one', 'two', 'three', 'four'];
    this.pagination.storeDataset('data', { id: 1 }, {}, arr);
    this.pagination.storeDataset('transit-key', { id: 2 }, {}, arr2);
    assert.strictEqual(this.pagination.get('lazyCaches').size, 2, 'it stores both keys');

    this.pagination.clearDataset('transit-key');
    assert.strictEqual(this.pagination.get('lazyCaches').size, 1, 'deletes one key');
    assert.notOk(this.pagination.get('lazyCaches').has('transit-key'), 'cache is no longer stored');
  });

  test('pagination.clearDataset with no args clears entire cache', function (assert) {
    const arr = ['one', 'two'];
    const arr2 = ['one', 'two', 'three', 'four'];
    this.pagination.storeDataset('data', { id: 1 }, {}, arr);
    this.pagination.storeDataset('transit-key', { id: 2 }, {}, arr2);
    assert.strictEqual(this.pagination.get('lazyCaches').size, 2, 'it stores both keys');

    this.pagination.clearDataset();
    assert.strictEqual(this.pagination.get('lazyCaches').size, 0, 'deletes all of the keys');
    assert.notOk(this.pagination.get('lazyCaches').has('transit-key'), 'first cache key is no longer stored');
    assert.notOk(this.pagination.get('lazyCaches').has('data'), 'second cache key is no longer stored');
  });

  test('pagination.getDataset', function (assert) {
    const arr = ['one', 'two'];
    this.pagination.storeDataset('data', { id: 1 }, {}, arr);

    assert.deepEqual(this.pagination.getDataset('data', { id: 1 }), { response: {}, dataset: arr });
  });

  test('pagination.constructResponse', function (assert) {
    const arr = ['one', 'two', 'three', 'fifteen', 'twelve'];
    this.pagination.storeDataset('data', { id: 1 }, {}, arr);

    assert.deepEqual(
      this.pagination.constructResponse('data', {
        id: 1,
        pageFilter: 't',
        page: 1,
        size: 3,
        responsePath: 'data',
      }),
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

  test('pagination.fetchPage', async function (assert) {
    const keys = ['zero', 'one', 'two', 'three', 'four', 'five', 'six'];
    const data = {
      data: {
        keys,
      },
    };
    const pageSize = 2;
    const query = {
      size: pageSize,
      page: 1,
      responsePath: 'data.keys',
    };
    this.pagination.storeDataset('transit-key', query, data, keys);

    let result;
    result = await this.pagination.fetchPage('transit-key', query);
    assert.strictEqual(result.get('length'), pageSize, 'returns the correct number of items');
    assert.deepEqual(
      result.map((r) => r.id),
      keys.slice(0, pageSize),
      'returns the first page of items'
    );
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

    result = await this.pagination.fetchPage('transit-key', {
      size: pageSize,
      page: 3,
      responsePath: 'data.keys',
    });
    const pageThreeEnd = 3 * pageSize;
    const pageThreeStart = pageThreeEnd - pageSize;
    assert.deepEqual(
      result.map((r) => r.id),
      keys.slice(pageThreeStart, pageThreeEnd),
      'returns the third page of items'
    );

    result = await this.pagination.fetchPage('transit-key', {
      size: pageSize,
      page: 99,
      responsePath: 'data.keys',
    });

    assert.deepEqual(
      result.map((r) => r.id),
      keys.slice(keys.length - 1),
      'returns the last page when the page value is beyond the of bounds'
    );

    result = await this.pagination.fetchPage('transit-key', {
      size: pageSize,
      page: 0,
      responsePath: 'data.keys',
    });
    assert.deepEqual(
      result.map((r) => r.id),
      keys.slice(0, pageSize),
      'returns the first page when page value is under the bounds'
    );
  });

  test('pagination.lazyPaginatedQuery', async function (assert) {
    const response = {
      data: ['foo'],
    };
    let queryArgs;
    const adapterForStub = () => {
      return {
        query(store, modelName, query) {
          queryArgs = query;
          return Promise.resolve(response);
        },
      };
    };
    Sinon.stub(this.store, 'adapterFor').callsFake(adapterForStub);
    // stub fetchPage because we test it separately
    Sinon.stub(this.pagination, 'fetchPage').callsFake(() => {});
    const query = { page: 1, size: 1, responsePath: 'data' };

    await this.pagination.lazyPaginatedQuery('transit-key', query);
    assert.deepEqual(
      this.pagination.getDataset('transit-key', query),
      { response: { data: null }, dataset: ['foo'] },
      'stores returned dataset'
    );

    await this.pagination.lazyPaginatedQuery('secret', { page: 1, responsePath: 'data' });
    assert.strictEqual(queryArgs.size, DEFAULT_PAGE_SIZE, 'calls query with DEFAULT_PAGE_SIZE');

    assert.throws(
      () => {
        this.pagination.lazyPaginatedQuery('transit-key', {});
      },
      /responsePath is required/,
      'requires responsePath'
    );
    assert.throws(
      () => {
        this.pagination.lazyPaginatedQuery('transit-key', { responsePath: 'foo' });
      },
      /page is required/,
      'requires page'
    );
  });

  test('pagination.filterData', async function (assert) {
    const dataset = [
      { id: 'foo', name: 'Foo', type: 'test' },
      { id: 'bar', name: 'Bar', type: 'test' },
      { id: 'bar-2', name: 'Bar', type: null },
    ];

    const defaultFiltering = this.pagination.filterData('foo', dataset);
    assert.deepEqual(defaultFiltering, [{ id: 'foo', name: 'Foo', type: 'test' }]);

    const filter = (data) => {
      return data.filter((d) => d.name === 'Bar' && d.type === 'test');
    };
    const customFiltering = this.pagination.filterData(filter, dataset);
    assert.deepEqual(customFiltering, [{ id: 'bar', name: 'Bar', type: 'test' }]);
  });
});
