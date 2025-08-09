/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

const MOCK_DATA = [
  { island: 'Maldives', visit_length: 5, is_booked: false, trip_date: '2025-06-22T00:00:00.000Z' },
  { island: 'Bora Bora', visit_length: 7, is_booked: true, trip_date: '2025-03-15T00:00:00.000Z' },
  { island: 'Fiji', visit_length: 10, is_booked: true, trip_date: '2025-09-08T00:00:00.000Z' },
  { island: 'Seychelles', visit_length: 6, is_booked: false, trip_date: '2025-12-03T00:00:00.000Z' },
  { island: 'Maui', visit_length: 8, is_booked: true, trip_date: '2026-01-18T00:00:00.000Z' },
];
module('Integration | Component | clients/table', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(async function () {
    this.data = undefined;
    this.columns = [
      { key: 'island', label: 'Islands', isSortable: true },
      { key: 'visit_length', label: 'Visit length', isSortable: true },
      { key: 'is_booked', label: 'Vacation booking status', isSortable: true },
      { key: 'trip_date', label: 'Date trip starts', isSortable: true },
    ];
    this;
    this.initiallySortBy = undefined;
    this.setPageSize = undefined;
    this.showPaginationSizeSelector = undefined;
    // helper function to setup table with X number of simple objects for pagination tests
    this.mockMoreData = (recordNumber) => {
      const record = (i) => ({ id: i, name: `record-${i}` });
      this.columns = [
        { key: 'id', label: 'ID', isSortable: true },
        { key: 'name', label: 'Name', isSortable: true },
      ];
      this.data = Array.from({ length: recordNumber }, (_, i) => record(i));
    };
    this.renderComponent = async () => {
      return render(hbs`
          <Clients::Table
            @columns={{this.columns}}
            @data={{this.data}}
            @initiallySortBy={{this.initiallySortBy}}
            @setPageSize={{this.setPageSize}}
            @showPaginationSizeSelector={{this.showPaginationSizeSelector}}
          />`);
    };
  });

  test('it renders default empty state when no data exists', async function (assert) {
    await this.renderComponent();
    assert.dom(CLIENT_COUNT.card('table empty state')).hasText('No data to display');
  });

  test('it renders yielded empty state block when no data exists', async function (assert) {
    await render(hbs`
      <Clients::Table @columns={{this.columns}} @data={{this.data}}>
      <:emptyState>Oh no, there's no data!</:emptyState>
      </Clients::Table>`);
    assert.dom(CLIENT_COUNT.card('table empty state')).hasText("Oh no, there's no data!");
  });

  test('it renders and paginates data', async function (assert) {
    this.data = MOCK_DATA;
    await this.renderComponent();
    assert.dom(CLIENT_COUNT.card('table empty state')).doesNotExist();
    assert.dom(GENERAL.paginationInfo).hasText(`1–3 of ${this.data.length}`);
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.tableData(0, 'island')).hasText('Seychelles', 'it paginates the data');
  });

  test('it sorts table data', async function (assert) {
    this.data = MOCK_DATA;
    const assertSortOrder = (expectedValues, { column, page }) => {
      expectedValues.forEach((value, idx) => {
        assert
          .dom(GENERAL.tableData(idx, column))
          .hasText(value, `page ${page}, row ${idx} has ${column}: ${value}`);
      });
    };

    await this.renderComponent();
    const [firstColumn, secondColumn, thirdColumn, fourthColumn] = findAll(GENERAL.icon('swap-vertical'));
    await click(firstColumn);
    assertSortOrder(['Bora Bora', 'Fiji', 'Maldives'], { column: 'island', page: 1 });
    await click(GENERAL.nextPage);
    assertSortOrder(['Maui', 'Seychelles'], { column: 'island', page: 2 });

    await click(GENERAL.prevPage);
    await click(secondColumn);
    assertSortOrder(['5', '6', '7'], { column: 'visit_length', page: 1 });
    await click(GENERAL.nextPage);
    assertSortOrder(['8', '10'], { column: 'visit_length', page: 2 });

    await click(GENERAL.prevPage);
    await click(thirdColumn);
    assertSortOrder(['false', 'false', 'true'], { column: 'is_booked', page: 1 });
    await click(GENERAL.nextPage);
    assertSortOrder(['true', 'true'], { column: 'is_booked', page: 2 });

    await click(GENERAL.prevPage);
    await click(fourthColumn);
    assertSortOrder(['2025-03-15T00:00:00.000Z', '2025-06-22T00:00:00.000Z', '2025-09-08T00:00:00.000Z'], {
      column: 'trip_date',
      page: 1,
    });
    await click(GENERAL.nextPage);
    assertSortOrder(['2025-12-03T00:00:00.000Z', '2026-01-18T00:00:00.000Z'], {
      column: 'trip_date',
      page: 2,
    });
  });

  test('it pre-sorts table data if @initiallySortBy is set', async function (assert) {
    this.data = MOCK_DATA;
    this.initiallySortBy = { column: 'visit_length', direction: 'desc' };
    await this.renderComponent();
    assert.dom(GENERAL.tableColumnHeader(2, { isAdvanced: true })).hasAttribute('aria-sort', 'descending');
    assert
      .dom(`${GENERAL.tableColumnHeader(2, { isAdvanced: true })} ${GENERAL.icon('arrow-down')}`)
      .exists();
    const firstPage = ['Fiji', 'Maui', 'Bora Bora'];
    firstPage.forEach((value, idx) => {
      assert.dom(GENERAL.tableData(idx, 'island')).hasText(value, `page 1, row ${idx} has ${value}`);
    });
    await click(GENERAL.nextPage);
    const secondPage = ['Seychelles', 'Maldives'];
    secondPage.forEach((value, idx) => {
      assert.dom(GENERAL.tableData(idx, 'island')).hasText(value, `page 2, row ${idx} has ${value}`);
    });
  });

  test('it sets page size if @setPageSize has a value', async function (assert) {
    this.mockMoreData(15);
    this.setPageSize = 8; // component default for testing is 3, so set to anything but 3
    await this.renderComponent();
    assert.dom(GENERAL.paginationInfo).hasText(`1–${this.setPageSize} of ${this.data.length}`);
    let idx = 7; // 8th item, table items are 0-indexed
    assert
      .dom(GENERAL.tableData(idx, 'name'))
      .hasText(this.data[idx].name, 'last row is 8th item in dataset');
    await click(GENERAL.nextPage);
    idx = 8; // 9th item, table items are 0-indexed
    assert
      .dom(GENERAL.tableData(0, 'name'))
      .hasText(this.data[idx].name, 'first row on page 2 is 9th item in dataset');
  });

  test('it renders size selector if @showPaginationSizeSelector is true', async function (assert) {
    this.mockMoreData(10);
    this.setPageSize = 5;
    this.showPaginationSizeSelector = true;
    await this.renderComponent();

    assert.dom(GENERAL.tableRow()).exists({ count: 5 }, '5 rows render');
    assert.dom(GENERAL.paginationInfo).hasText(`1–${this.setPageSize} of ${this.data.length}`);
    assert.dom(GENERAL.paginationSizeSelector).hasValue('5');
    let idx = 4; // rows and ids are 0-indexed
    assert.dom(GENERAL.tableData(idx, 'id')).hasText(`${idx}`, `last row is ${idx + 1}th item in dataset`);
    await fillIn(GENERAL.paginationSizeSelector, '10');
    idx = 9;
    assert.dom(GENERAL.tableRow()).exists({ count: 10 }, '10 rows render');
    assert.dom(GENERAL.paginationSizeSelector).hasValue('10', 'it updates the size selector to 10');
    assert.dom(GENERAL.tableData(idx, 'id')).hasText(`${idx}`, `last row is ${idx + 1}th item in dataset`);
  });

  test('it renders "Deleted" badge for "mount_type" keys if the value is "deleted mount"', async function (assert) {
    this.columns = [
      { key: 'mount_type', label: 'Mount type' },
      { key: 'mount_path', label: 'Mount path' },
    ];
    this.data = [
      { mount_type: 'deleted mount', mount_path: 'auth/userpass/' },
      { mount_type: 'ns_token', mount_path: 'auth/token/' },
    ];
    await this.renderComponent();
    assert.dom('.hds-badge').exists({ count: 1 }, 'only one badge renders');
    assert
      .dom(`${GENERAL.tableData(0, 'mount_type')} .hds-badge`)
      .exists('it renders a badge for the deleted mount')
      .hasText('Deleted');
  });
});
