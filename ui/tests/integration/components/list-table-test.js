/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, waitFor, find } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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
    this.data = undefined;
    this.columns = [
      { key: 'island', label: 'Islands', isSortable: true },
      { key: 'visit_length', label: 'Visit length', customTableItem: true },
      { key: 'trip_date', label: 'Date trip starts' },
      { key: 'popupMenu', label: 'Action' },
    ];
    this;

    this.renderComponent = async () => {
      return render(hbs`
          <ListTable
            @columns={{this.columns}}
            @data={{this.data}}
            @selectionKey="island"
          >
            <:customTableItem as |itemData|>
            Custom Table Item rendered!
            </:customTableItem>

            <:popupMenu as |rowData|>
            <Hds::Dropdown as |D|>
            <D.ToggleButton @text="Menu" data-test-popup-menu-trigger />
            <D.Title @text="Title Text" />
            <D.Description @text="Sample text" />
            <D.Interactive @route="components" @icon="trash" @color="critical">Delete</D.Interactive>
            </Hds::Dropdown>
            </:popupMenu>

          </ListTable>`);
    };
  });

  test('it renders and paginates data', async function (assert) {
    this.data = MOCK_DATA;
    await this.renderComponent();
    assert.dom(GENERAL.paginationInfo).hasText(`1–6 of ${this.data.length}`);

    await fillIn(GENERAL.paginationSizeSelector, '5'); // Default is 10, so change to something else
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.tableRow('Seychelles', 'island')).exists('it paginates the data');
  });

  test('it sorts table data by a sortable column', async function (assert) {
    this.data = MOCK_DATA;
    const assertSortOrder = (expectedValues, { column, page }) => {
      expectedValues.forEach((value, idx) => {
        assert
          .dom(GENERAL.tableData(value, column))
          .hasText(value, `page ${page}, row ${idx} has ${column}: ${value}`);
      });
    };

    await this.renderComponent();
    const column = find(GENERAL.icon('swap-vertical'));
    await click(column);
    assertSortOrder(['Bora Bora', 'Fiji', 'Maldives', 'Maui', 'Santorini', 'Seychelles'], {
      column: 'island',
      page: 1,
    });
  });

  test('action column renders provided yield block with popup menu', async function (assert) {
    this.data = MOCK_DATA;
    await this.renderComponent();

    assert.dom(GENERAL.tableData('Maldives', 'popupMenu')).exists('action column renders');
    assert.dom(GENERAL.menuTrigger).exists('button trigger exists for popup menu');
  });

  test('selectable column renders when isSelectable is true', async function (assert) {
    this.data = MOCK_DATA;
    await this.renderComponent();

    assert
      .dom(`${GENERAL.tableRow('Maldives')} > th`)
      .hasClass('hds-table__th--is-selectable', 'selectable column renders for row');
  });

  // check that a custom item block will render
  test('custom item renders provided yield block with customTableItem for a column has customTableItem set to true', async function (assert) {
    this.data = MOCK_DATA;
    await this.renderComponent();

    assert
      .dom(GENERAL.tableData('Maldives', 'visit_length'))
      .hasText('Custom Table Item rendered!', 'custom item renders');
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
    ``;
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
