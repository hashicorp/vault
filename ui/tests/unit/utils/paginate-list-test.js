/**
 * Copyright IBM Corp. 2016, 2025
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
      filteredTotal: 20,
      pageSize: 3,
    };
    assert.deepEqual(meta, expectedMeta, 'returns correct meta data');
  });

  test('it should return remaining results on last page', async function (assert) {
    const paginatedData = paginate(this.items, { page: 7, pageSize: 3 });
    const expected = [18, 19];
    assert.deepEqual(paginatedData, expected, 'returns correct items for last page');
  });

  test('filteredTotal reflects total matching items', function (assert) {
    // 20 items, filter matches first 6 (0-5), paginate to page 1 with size 4
    const data = Array.from({ length: 20 }, (_, i) => ({ id: i, name: i < 6 ? `match-${i}` : `skip-${i}` }));
    const { meta } = paginate(data, { page: 1, pageSize: 4, filter: 'match', filterKey: 'name' });
    assert.strictEqual(meta.filteredTotal, 6, 'filteredTotal is total matching items across all pages');
    assert.strictEqual(meta.lastPage, 2, 'lastPage is based on filteredTotal');
  });

  test('it should reset to page 1 when page exceeds lastPage', function (assert) {
    // 20 items, pageSize 10 = 2 pages; requesting page 5 should fall back to page 1
    const paginatedData = paginate(this.items, { page: 5, pageSize: 10 });
    assert.deepEqual(
      paginatedData,
      this.items.slice(0, 10),
      'returns first page of items when page is out of bounds'
    );
    assert.strictEqual(
      paginatedData.meta.currentPage,
      1,
      'currentPage in meta is 1, not the out-of-bounds page'
    );
  });

  test('meta currentPage matches actual data shown when page is out of bounds', function (assert) {
    const { meta } = paginate(this.items, { page: 99, pageSize: 5 });
    assert.strictEqual(
      meta.currentPage,
      1,
      'currentPage in meta reflects actual page shown, not requested page'
    );
    assert.strictEqual(meta.lastPage, 4, 'lastPage is computed correctly');
  });
});
