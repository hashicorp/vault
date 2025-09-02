/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, findAll, render, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

const MOCK_DATA = [
  { island: 'Maldives', visit_length: 5, is_booked: false, trip_date: '2025-06-22T00:00:00.000Z' },
  { island: 'Bora Bora', visit_length: 7, is_booked: true, trip_date: '2025-03-15T00:00:00.000Z' },
  { island: 'Fiji', visit_length: 10, is_booked: true, trip_date: '2025-09-08T00:00:00.000Z' },
  { island: 'Santorini', visit_length: 4, is_booked: false, trip_date: '2026-04-10T00:00:00.000Z' },
  { island: 'Maui', visit_length: 8, is_booked: true, trip_date: '2026-01-18T00:00:00.000Z' },
  { island: 'Seychelles', visit_length: 6, is_booked: false, trip_date: '2025-12-03T00:00:00.000Z' },
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
    assert.dom(GENERAL.paginationInfo).hasText(`1–5 of ${this.data.length}`);
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
    assertSortOrder(['Bora Bora', 'Fiji', 'Maldives', 'Maui', 'Santorini'], { column: 'island', page: 1 });
    await click(GENERAL.nextPage);
    assertSortOrder(['Seychelles'], { column: 'island', page: 2 });

    await click(GENERAL.prevPage);
    await click(secondColumn);
    assertSortOrder(['4', '5', '6', '7', '8'], { column: 'visit_length', page: 1 });
    await click(GENERAL.nextPage);
    assertSortOrder(['10'], { column: 'visit_length', page: 2 });

    await click(GENERAL.prevPage);
    await click(thirdColumn);
    assertSortOrder(['false', 'false', 'false', 'true', 'true'], { column: 'is_booked', page: 1 });
    await click(GENERAL.nextPage);
    assertSortOrder(['true'], { column: 'is_booked', page: 2 });

    await click(GENERAL.prevPage);
    await click(fourthColumn);
    assertSortOrder(
      [
        '2025-03-15T00:00:00.000Z',
        '2025-06-22T00:00:00.000Z',
        '2025-09-08T00:00:00.000Z',
        '2025-12-03T00:00:00.000Z',
        '2026-01-18T00:00:00.000Z',
      ],
      {
        column: 'trip_date',
        page: 1,
      }
    );
    await click(GENERAL.nextPage);
    assertSortOrder(['2026-04-10T00:00:00.000Z'], {
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
    const firstPage = ['Fiji', 'Maui', 'Bora Bora', 'Seychelles', 'Maldives'];
    firstPage.forEach((value, idx) => {
      assert.dom(GENERAL.tableData(idx, 'island')).hasText(value, `page 1, row ${idx} has ${value}`);
    });
    await click(GENERAL.nextPage);
    const secondPage = ['Santorini'];
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

  test('it resets pagination when data changes', async function (assert) {
    // We need more than 5 rows, so here's more mock data!
    const moreData = [
      { island: 'Tahiti', visit_length: 12, is_booked: true, trip_date: '2025-05-10T00:00:00.000Z' },
      { island: 'Barbados', visit_length: 6, is_booked: false, trip_date: '2025-08-25T00:00:00.000Z' },
      { island: 'Cyprus', visit_length: 9, is_booked: true, trip_date: '2026-03-12T00:00:00.000Z' },
      { island: 'Jamaica', visit_length: 7, is_booked: false, trip_date: '2025-11-05T00:00:00.000Z' },
      { island: 'Crete', visit_length: 11, is_booked: true, trip_date: '2026-06-18T00:00:00.000Z' },
      { island: 'Aruba', visit_length: 5, is_booked: false, trip_date: '2025-10-14T00:00:00.000Z' },
    ];
    this.data = [...MOCK_DATA, ...moreData];
    this.showPaginationSizeSelector = true;
    await this.renderComponent();
    await fillIn(GENERAL.paginationSizeSelector, '10'); // Default is 5, so change to something else
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.paginationInfo).hasText(`11–12 of ${this.data.length}`, 'it navigates to next page');
    assert.dom(GENERAL.tableRow()).exists({ count: 2 }, '2 row renders');
    // Changing the @data arg should trigger an update and reset pagination
    // We have to use `this.set` to trigger did-update
    this.set('data', [
      { island: 'Palawan', visit_length: 9, is_booked: true, trip_date: '2025-11-14T00:00:00.000Z' },
      { island: 'Mykonos', visit_length: 3, is_booked: false, trip_date: '2026-02-28T00:00:00.000Z' },
    ]);

    // There's a workaround using next() from @ember/runloop because the Hds::Pagination::Numbered component
    // doesn't re-render when @currentPage updates. When that's fixed at the source we should be able to remove waitFor
    await waitFor(GENERAL.paginationInfo);
    assert.dom(GENERAL.paginationInfo).hasText(`1–2 of ${this.data.length}`);
    assert.dom(GENERAL.tableRow()).exists({ count: 2 }, '2 rows render');
    assert.dom(GENERAL.paginationSizeSelector).hasValue('10', 'page selector is unchanged when data updates');
  });
});
