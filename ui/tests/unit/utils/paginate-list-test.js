/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { paginate } from 'core/utils/paginate-list';
import { module, test } from 'qunit';

module('Unit | Utility | paginate-list', function (hooks) {
  hooks.beforeEach(function () {
    this.items = Array.from({ length: 20 }, (val, i) => i);
  });

  test('it should return data for non or empty arrays', function (assert) {
    let items = 'not an array';
    assert.strictEqual(paginate(items), items, 'returns the same data when input is not an array');

    items = [];
    assert.deepEqual(paginate(items), items, 'returns the same data when input is an empty array');
  });

  test('it should use default page and size', function (assert) {
    const paginatedData = paginate(this.items, { page: 1 });
    assert.strictEqual(paginatedData.length, 15, 'returns 15 items as default when no page size is set');
    assert.deepEqual(this.items.slice(0, 15), paginatedData, 'returns first page of items by default');
  });

  test('it should return items for given page and size', function (assert) {
    const paginatedData = paginate(this.items, { page: 3, pageSize: 5 });
    assert.strictEqual(paginatedData.length, 5, 'returns correct number of items based on size');
    assert.deepEqual(this.items.slice(10, 15), paginatedData, 'returns correct items for given page');
  });

  test('it should return remaining items on last page', function (assert) {
    const paginatedData = paginate(this.items, { page: 3, pageSize: 8 });
    assert.strictEqual(paginatedData.length, 4, 'returns correct number of items on last page');
    assert.deepEqual(this.items.slice(16), paginatedData, 'returns correct items for last page');
  });

  test('it should filter items and then paginate', function (assert) {
    let data = ['Test', 'foo', 'test', 'bar', 'test'];
    let expected = ['Test', 'test'];
    const options = { page: 1, pageSize: 2, filter: 'test' };

    let paginatedData = paginate(data, options);
    assert.deepEqual(paginatedData, expected, 'returns correct number of filtered items');

    data = data.map((id) => ({ id }));
    expected = [{ id: 'Test' }, { id: 'test' }];
    paginatedData = paginate(data, { ...options, filterKey: 'id' });
    assert.deepEqual(paginatedData, expected, 'returns correct number of filtered objects');
  });

  test('it should add meta data to returned object', function (assert) {
    const { meta } = paginate(this.items, { page: 2, pageSize: 3 });
    const expectedMeta = {
      currentPage: 2,
      lastPage: 7,
      nextPage: 3,
      prevPage: 1,
      total: 20,
      filteredTotal: 3,
      pageSize: 3,
    };
    assert.deepEqual(meta, expectedMeta, 'returns correct meta data');
  });

  test('it should return remaining results on last page', async function (assert) {
    const paginatedData = paginate(this.items, { page: 7, pageSize: 3 });
    const expected = [18, 19];
    assert.deepEqual(paginatedData, expected, 'returns correct items for last page');
  });
});
