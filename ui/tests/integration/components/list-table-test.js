/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const MOCK_DATA = [
  { island: 'Maldives', visit_length: 5, trip_date: '2025-06-22T00:00:00.000Z' },
  { island: 'Bora Bora', visit_length: 7, trip_date: '2025-03-15T00:00:00.000Z' },
  { island: 'Fiji', visit_length: 10, trip_date: '2025-09-08T00:00:00.000Z' },
  { island: 'Santorini', visit_length: 4, trip_date: '2026-04-10T00:00:00.000Z' },
  { island: 'Maui', visit_length: 8, trip_date: '2026-01-18T00:00:00.000Z' },
  { island: 'Seychelles', visit_length: 6, trip_date: '2025-12-03T00:00:00.000Z' },
];
module('Integration | Component | list-table', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(async function () {
    this.data = MOCK_DATA;
    this.onSelectionChange = undefined;
    this.selectionKeyField = undefined;
    this.columns = [
      { key: 'island', label: 'Islands', isSortable: true },
      { key: 'visit_length', label: 'Visit length', customTableItem: true },
      { key: 'trip_date', label: 'Date trip starts' },
      { key: 'popupMenu', label: 'Action' },
    ];

    this.renderComponent = async () => {
      return render(hbs`
          <ListTable
            @columns={{this.columns}}
            @data={{this.data}}
            @selectionKeyField={{this.selectionKeyField}}
            @onSelectionChange={{this.onSelectionChange}}
          >
            <:customTableItem as |itemData|>
              <Hds::BadgeCount @text={{itemData.visit_length}} @type="outlined" />            
            </:customTableItem>

            <:popupMenu as |rowData|>
              <Hds::Dropdown as |D|>
              <D.ToggleButton @text="Menu" data-test-popup-menu-trigger />
              <D.Title @text={{rowData.island}} />
              <D.Description @text="Sample text" />
              <D.Interactive @route="components" @icon="trash" @color="critical">Delete</D.Interactive>
              </Hds::Dropdown>
            </:popupMenu>

          </ListTable>`);
    };
  });

  test('it renders and paginates data', async function (assert) {
    await this.renderComponent();
    assert.dom('input[type="checkbox"]').doesNotExist('table is not selectable by default');
    assert.dom(GENERAL.paginationInfo).hasText(`1–6 of ${this.data.length}`);
    // Default is 10, so change to something else to test pagination
    await fillIn(GENERAL.paginationSizeSelector, '5');
    assert.dom(GENERAL.paginationInfo).hasText(`1–5 of ${this.data.length}`);
    assert.dom(GENERAL.tableRow()).exists({ count: 5 }, 'only 5 rows render');
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.paginationInfo).hasText(`6–6 of ${this.data.length}`);
    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'only 1 row renders on second page');
    assert.dom(GENERAL.tableData(0, 'island')).hasText('Seychelles', 'second page has expected row');
  });

  test('it does not render popup menu if @columns does not include a popupMenu key', async function (assert) {
    this.columns = [
      { key: 'island', label: 'Islands', isSortable: true },
      { key: 'visit_length', label: 'Visit length', customTableItem: true },
      { key: 'trip_date', label: 'Date trip starts' },
    ];
    await this.renderComponent();
    assert.dom(GENERAL.menuTrigger).doesNotExist();
  });

  test('it stringifies object and array values for non-custom columns', async function (assert) {
    this.columns = [
      { key: 'island', label: 'Islands' },
      { key: 'trip_details', label: 'Trip details' },
      { key: 'tags', label: 'Tags' },
    ];
    this.data = [
      {
        island: 'Maldives',
        trip_details: { hotel: 'Atoll Inn', nights: 5 },
        tags: ['beach', 'snorkel'],
      },
    ];

    await this.renderComponent();
    assert.dom(GENERAL.tableData(0, 'trip_details')).hasText('{ "hotel": "Atoll Inn", "nights": 5 }');
    assert.dom(GENERAL.tableData(0, 'tags')).hasText('[ "beach", "snorkel" ]');
  });

  test('it does not render popup menu if parent does not yield one', async function (assert) {
    await render(hbs`
          <ListTable
            @columns={{this.columns}}
            @data={{this.data}}
            @selectionKeyField={{this.selectionKeyField}}
            @onSelectionChange={{this.onSelectionChange}}
          />`);
    assert.dom(GENERAL.menuTrigger).doesNotExist();
  });

  test('it sorts table data by a sortable column', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.icon('swap-vertical'));
    const expectedOrder = ['Bora Bora', 'Fiji', 'Maldives', 'Maui', 'Santorini', 'Seychelles'];
    expectedOrder.forEach((island, idx) => {
      assert.dom(GENERAL.tableData(idx, 'island')).hasText(island);
    });
  });

  test('action column renders provided yield block with popup menu', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.menuTrigger).exists({ count: this.data.length }, 'popup trigger exists for each item');
    await click(`${GENERAL.tableRow(2)} ${GENERAL.menuTrigger}`);
    assert.dom('li').hasText(this.data[2].island, 'popup menu renders relevant row data');
  });

  test('selectable checkboxes render and are selectable when selectionKeyField is provided', async function (assert) {
    this.selectionKeyField = 'island';
    const count = this.data.length + 1;
    this.onSelectionChange = sinon.spy();
    await this.renderComponent();
    assert
      .dom('input[type="checkbox"]')
      .exists({ count }, 'it renders a checkbox for each row plus the header to select all');
    assert
      .dom(`${GENERAL.tableRow(0)} input[type="checkbox"]`)
      .hasAttribute(
        'aria-label',
        `Select row ${this.data[0][this.selectionKeyField]}`,
        'selection aria label suffix uses selectionKeyField in value'
      );
    await click(`${GENERAL.tableRow(0)} input[type="checkbox"]`);
    await click(`${GENERAL.tableRow(2)} input[type="checkbox"]`);
    assert.true(this.onSelectionChange.calledTwice, 'onSelectionChange is called twice');
    const [callbackArgs] = this.onSelectionChange.lastCall.args;
    const { selectionKey, selectedRowsKeys, selectableRowsStates } = callbackArgs;
    const lastItemSelected = this.data[2];
    assert.strictEqual(selectionKey, lastItemSelected.island, 'selectionKey is last selected row');
    assert.propEqual(selectedRowsKeys, ['Maldives', 'Fiji'], 'callback passes selectedRowKeys');
    const expectedRowStates = [
      { selectionKey: 'Maldives', isSelected: true },
      { selectionKey: 'Bora Bora', isSelected: false },
      { selectionKey: 'Fiji', isSelected: true },
      { selectionKey: 'Santorini', isSelected: false },
      { selectionKey: 'Maui', isSelected: false },
      { selectionKey: 'Seychelles', isSelected: false },
    ];
    assert.propEqual(selectableRowsStates, expectedRowStates, 'callback contains selectableRowsStates');
  });

  test('it is still selectable when selection callback is not provided', async function (assert) {
    this.selectionKeyField = 'island';
    await this.renderComponent();
    assert.dom('input[type="checkbox"]').exists({ count: this.data.length + 1 });
    const firstRowCheckbox = '[data-test-table-row="0"] input[type="checkbox"]';
    await click(firstRowCheckbox);
    assert.dom(firstRowCheckbox).isChecked('row checkbox can be toggled without @onSelectionChange');
  });

  test('custom item renders provided yield block with customTableItem for a column has customTableItem set to true', async function (assert) {
    await this.renderComponent();
    assert
      .dom(`${GENERAL.tableData(0, 'visit_length')} .hds-badge-count`)
      .hasText('5', 'custom table item renders yielded badge');
  });

  test('it resets pagination when data changes', async function (assert) {
    const moreData = [
      { island: 'Tahiti', visit_length: 12, trip_date: '2025-05-10T00:00:00.000Z' },
      { island: 'Barbados', visit_length: 6, trip_date: '2025-08-25T00:00:00.000Z' },
      { island: 'Cyprus', visit_length: 9, trip_date: '2026-03-12T00:00:00.000Z' },
      { island: 'Jamaica', visit_length: 7, trip_date: '2025-11-05T00:00:00.000Z' },
      { island: 'Crete', visit_length: 11, trip_date: '2026-06-18T00:00:00.000Z' },
      { island: 'Aruba', visit_length: 5, trip_date: '2025-10-14T00:00:00.000Z' },
    ];
    this.data = [...MOCK_DATA, ...moreData];
    await this.renderComponent();
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.paginationInfo).hasText(`11–12 of ${this.data.length}`, 'it navigates to next page');
    // Changing the @data arg should trigger an update and reset pagination
    this.set('data', [
      { island: 'Palawan', visit_length: 9, trip_date: '2025-11-14T00:00:00.000Z' },
      { island: 'Mykonos', visit_length: 3, trip_date: '2026-02-28T00:00:00.000Z' },
    ]);

    // There's a workaround using next() from @ember/runloop because the Hds::Pagination::Numbered component
    // doesn't re-render when @currentPage updates. When that's fixed at the source we should be able to remove waitFor
    await waitFor(GENERAL.paginationInfo);
    assert.dom(GENERAL.paginationInfo).hasText(`1–2 of ${this.data.length}`);
    assert.dom(GENERAL.paginationSizeSelector).hasValue('10', 'page selector is unchanged when data updates');
  });
});
