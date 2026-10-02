/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const BREADCRUMBS = [
  { label: 'Vault', route: 'vault', icon: 'vault' },
  { label: 'Access', route: 'vault.cluster.access' },
  { label: 'Entities' },
];

const COLUMNS = [
  { key: 'name', label: 'Name', isSortable: true },
  { key: 'id', label: 'ID', customTableItem: true },
  { key: 'popupMenu', label: 'Actions', width: '8%' },
];

const MOCK_ITEMS = [
  { id: 'aaa', name: 'alice' },
  { id: 'bbb', name: 'bob' },
  { id: 'ccc', name: 'carol' },
];

module('Integration | Component | page/list-view', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.config = {
      title: 'Entities',
      description: 'Entities represent distinct users or systems.',
      breadcrumbs: BREADCRUMBS,
      columns: COLUMNS,
      noDataTitle: 'No entities yet',
      noDataDescription: 'Create an entity to get started.',
      filteredEmptyTitle: 'No entities matching',
    };
    this.model = MOCK_ITEMS;
    this.page = 1;
  });

  // ── Test 1 — page title and breadcrumbs ─────────────────────────────────────
  test('it renders the page title and breadcrumbs', async function (assert) {
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Entities', 'renders @config.title in the page header');
    assert.dom(GENERAL.breadcrumb).exists({ count: 3 }, 'renders all 3 breadcrumb items');
  });

  // ── Test 2 — description ─────────────────────────────────────────────────────
  test('it renders @config.description when provided', async function (assert) {
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    assert
      .dom(GENERAL.hdsPageHeaderDescription)
      .hasText('Entities represent distinct users or systems.', 'renders @config.description');
  });

  // ── Test 3 — headerActions block ─────────────────────────────────────────────
  test('it renders the <:headerActions> named block', async function (assert) {
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      >
        <:headerActions>
          <Hds::Button @color="primary" @icon="plus" @text="Create entity" data-test-button="create-entity" />
        </:headerActions>
      </Page::ListView>
    `);

    assert
      .dom(GENERAL.button('create-entity'))
      .exists('renders the primary create button from <:headerActions>');
  });

  // ── Test 4 — table renders when data exists ──────────────────────────────────
  test('it renders the table when data is present', async function (assert) {
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      >
        <:popupMenu as |rowData|>
          <span data-test-popup-menu-trigger={{rowData.id}}>menu</span>
        </:popupMenu>
      </Page::ListView>
    `);

    assert.dom(GENERAL.tableRow()).exists({ count: MOCK_ITEMS.length }, 'renders one row per data item');
    assert.dom(GENERAL.emptyStateTitle).doesNotExist('no empty state when data is present');
  });

  // ── Test 5 — no-data empty state ─────────────────────────────────────────────
  test('it renders the no-data empty state when data is empty and no filter', async function (assert) {
    this.model = [];

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    assert.dom(GENERAL.emptyStateTitle).hasText('No entities yet', 'renders @config.noDataTitle');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Create an entity to get started.', 'renders @config.noDataDescription');
  });

  // ── Test 6 — filtered empty state ────────────────────────────────────────────
  test('it renders the filtered empty state when a filter returns no results', async function (assert) {
    // pageFilter is internal state — drive it via fillIn on the FilterInput.
    // The model has items but none match "zzz", so the filtered-empty state appears.
    this.config = {
      ...this.config,
      filter: { type: 'text', placeholder: 'Filter entities', ariaLabel: 'Filter', filterKey: 'name' },
    };

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    await fillIn(GENERAL.filterInput, 'zzz');

    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No entities matching "zzz"', 'renders filtered empty title with filter value appended');
  });

  test('it renders filtered-empty state for a caller-filtered model', async function (assert) {
    this.model = [];

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @filterValue="zzz"
        @page={{this.page}}
      />
    `);

    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText(
        'No entities matching "zzz"',
        'uses the caller-provided filter value in the empty-state title'
      );
  });

  // ── Test 7 — toolbar renders when toolbarActions block provided ───────────────
  test('it renders the toolbar when <:toolbarActions> is provided', async function (assert) {
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      >
        <:toolbarActions>
          <input placeholder="Filter entities" data-test-filter-input />
        </:toolbarActions>
      </Page::ListView>
    `);

    assert.dom(GENERAL.filterInput).exists('renders the toolbar content from <:toolbarActions>');
  });

  // ── Test 8 — popupMenu block renders once per row ─────────────────────────────
  test('it renders the <:popupMenu> block once per data row', async function (assert) {
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      >
        <:popupMenu as |rowData|>
          <button type="button" data-test-popup-menu-trigger={{rowData.id}}>menu</button>
        </:popupMenu>
      </Page::ListView>
    `);

    assert
      .dom(GENERAL.menuTrigger)
      .exists({ count: MOCK_ITEMS.length }, 'popup menu trigger rendered once per row');
  });

  // ── Test 9 — customTableItem block renders for custom columns ─────────────────
  test('it renders the <:customTableItem> block for columns with customTableItem: true', async function (assert) {
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      >
        <:customTableItem as |rowData column|>
          {{#if (eq column.key "id")}}
            <span data-test-custom-id={{rowData.id}}>{{rowData.id}}</span>
          {{/if}}
        </:customTableItem>
        <:popupMenu as |rowData|>
          <span data-test-popup-menu-trigger={{rowData.id}}>menu</span>
        </:popupMenu>
      </Page::ListView>
    `);

    assert
      .dom('[data-test-custom-id="aaa"]')
      .exists('customTableItem block fired for the "id" column on first row');
    assert.dom('[data-test-custom-id="ccc"]').exists('customTableItem block fired for all rows');
  });

  // ── Test 10 — confirmModal block renders when provided ────────────────────────
  test('it renders the <:confirmModal> named block', async function (assert) {
    this.showModal = true;

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      >
        <:popupMenu as |rowData|>
          <button type="button" data-test-popup-menu-trigger={{rowData.id}}
            {{on "click" (fn (mut this.showModal) true)}}>menu</button>
        </:popupMenu>
        <:confirmModal>
          {{#if this.showModal}}
            <div data-test-confirm-modal>
              <p data-test-confirm-action-title>Delete this entity?</p>
              <button type="button" data-test-confirm-button
                {{on "click" (fn (mut this.showModal) false)}}>Confirm</button>
            </div>
          {{/if}}
        </:confirmModal>
      </Page::ListView>
    `);

    assert.dom(GENERAL.confirmModal).exists('confirm modal content rendered from <:confirmModal> block');
    assert
      .dom(GENERAL.confirmTitle)
      .hasText('Delete this entity?', 'confirm modal title matches expected text');

    await click(GENERAL.confirmButton);
    assert.dom(GENERAL.confirmModal).doesNotExist('confirm modal is dismissed after confirm click');
  });

  // ── Test 11 — typing a filter from page 2 shows filtered results, never calls transitionTo ──
  test('typing a filter while on page 2 shows filtered results without triggering a router transition', async function (assert) {
    // Root cause of the search-from-page-2 bug (VAULT-50816-style):
    // previously onFilterChange called router.transitionTo({ queryParams: { page: 1 } }),
    // which triggered a full route model refresh (page has refreshModel: true), destroyed
    // the Page::ListView component instance, and wiped @tracked pageFilter back to ''.
    // Fix: paginate() already resets to page 1 when page > lastPage, so no transition is needed.
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    // Build 25 items so page 2 is reachable (pageSize defaults to 10).
    const items = Array.from({ length: 25 }, (_, i) => ({ id: `id-${i}`, name: `item-${i}` }));
    // Add a clearly identifiable item that will survive the filter.
    items.push({ id: 'id-target', name: 'alice-target' });
    this.model = items;
    this.config = {
      ...this.config,
      filter: { type: 'text', placeholder: 'Filter entities', ariaLabel: 'Filter', filterKey: 'name' },
    };

    // Start on page 2.
    this.page = 2;
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    await fillIn(GENERAL.filterInput, 'alice-target');

    // The filter must remain in place — the input still shows the typed value.
    assert.dom(GENERAL.filterInput).hasValue('alice-target', 'filter input retains typed value');

    // The table must show only the matching row (paginate resets to page 1 for filtered results).
    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'only the matching row is shown after filtering');

    // No router transition must have been triggered — transitioning would destroy the component.
    assert.ok(
      transitionStub.notCalled,
      'router.transitionTo was NOT called — no route refresh on filter change'
    );
  });

  // ── Test 11b — filter on page 1 also does NOT trigger a router transition ──────
  test('typing a filter while on page 1 also does not call router.transitionTo', async function (assert) {
    // Confirms the no-transition guarantee holds regardless of starting page.
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.config = {
      ...this.config,
      filter: { type: 'text', placeholder: 'Filter entities', ariaLabel: 'Filter', filterKey: 'name' },
    };

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    await fillIn(GENERAL.filterInput, 'alice');

    assert.dom(GENERAL.filterInput).hasValue('alice', 'filter input retains typed value on page 1');
    assert.ok(transitionStub.notCalled, 'router.transitionTo was not called when starting from page 1');
  });

  // ── Test 12 — filteredData filters by config.filter.filterKey ────────────────
  test('it filters rows by config.filter.filterKey when the user types in the filter input', async function (assert) {
    // pageFilter is internal state — drive it via fillIn on the FilterInput.
    // Only "alice" should survive the filter — "bob" and "carol" do not match.
    this.config = {
      ...this.config,
      filter: { type: 'text', placeholder: 'Filter', ariaLabel: 'Filter', filterKey: 'name' },
    };

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    await fillIn(GENERAL.filterInput, 'alice');

    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'only the matching row is rendered');
  });

  // ── Test 13 — confirmDelete calls deleteAction, shows flash, refreshes ────────
  test('it calls deleteAction and shows a success flash when the confirm button is clicked', async function (assert) {
    const deleteAction = sinon.stub().resolves();
    const refreshStub = sinon.stub(this.owner.lookup('service:router'), 'refresh');
    const flash = this.owner.lookup('service:flash-messages');
    const successStub = sinon.stub(flash, 'success');

    this.config = {
      ...this.config,
      rowActions: [
        {
          label: 'Delete',
          dataTest: 'delete',
          kind: 'modal',
          color: 'critical',
          modal: {
            title: 'Delete entity',
            titleItemDisplayKey: 'name',
            body: 'This will permanently delete',
            itemDisplayKey: 'name',
            deleteAction,
          },
        },
      ],
    };

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    // Open the ··· menu and trigger the modal
    await click(GENERAL.menuTrigger);
    await click('[data-test-popup-menu="delete"]');

    assert.dom('[data-test-confirm-action-title]').hasText('Delete entity alice?');

    // Confirm
    await click(GENERAL.confirmButton);

    assert.ok(deleteAction.calledOnce, 'deleteAction was called once');
    assert.ok(
      deleteAction.calledWith(sinon.match({ id: 'aaa' })),
      'deleteAction received the first row item'
    );
    assert.ok(successStub.calledOnce, 'success flash was shown');
    assert.ok(refreshStub.calledOnce, 'router.refresh was called');
    assert.dom('[data-test-confirm-modal]').doesNotExist('modal is closed after confirm');
  });

  // ── Test 14 — row actions honour the per-row capability key ─────────────────
  test('it only renders a row action when the row has the capability it requires', async function (assert) {
    this.config = {
      ...this.config,
      rowActions: [
        { label: 'View', dataTest: 'view', kind: 'action', actionName: 'closeModal' },
        {
          label: 'Delete',
          dataTest: 'delete',
          kind: 'action',
          actionName: 'closeModal',
          capability: 'canDelete',
        },
      ],
    };
    this.model = [
      { id: 'aaa', name: 'alice', capabilities: { canDelete: true } },
      { id: 'bbb', name: 'bob', capabilities: { canDelete: false } },
    ];

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    await click('[data-test-popup-menu-trigger="aaa"]');
    assert
      .dom(`${GENERAL.tableRow(0)} ${GENERAL.menuItem('view')}`)
      .exists('ungated action renders for a permitted row');
    assert
      .dom(`${GENERAL.tableRow(0)} ${GENERAL.menuItem('delete')}`)
      .exists('gated action renders when the row has the capability');

    await click('[data-test-popup-menu-trigger="bbb"]');
    assert
      .dom(`${GENERAL.tableRow(1)} ${GENERAL.menuItem('view')}`)
      .exists('ungated action renders for a restricted row');
    assert
      .dom(`${GENERAL.tableRow(1)} ${GENERAL.menuItem('delete')}`)
      .doesNotExist('gated action is hidden when the row lacks the capability');
  });

  // ── Test 15 — sort sorts the full dataset before pagination ──────────────────
  test('sorting a column sorts the entire dataset before pagination, not just the visible page', async function (assert) {
    // Build 12 items whose names, when sorted A→Z, place "item-00" first and
    // "zzz-name" last.  With pageSize = 10 the default order puts "zzz-name" on
    // page 2 (item 13).  After passing @sortBy="name" @sortOrder="desc",
    // "zzz-name" must appear on page 1 (the global maximum must be shown first).
    //
    // Sort state is driven by args (not internal clicks) because the integration
    // test has no router — @onSortChange would have nowhere to write back to.
    // Acceptance tests cover the full click→URL→re-render round-trip.
    const items = Array.from({ length: 12 }, (_, i) => ({
      id: `id-${i}`,
      name: `item-${String(i).padStart(2, '0')}`,
    }));
    // Add one item that should sort last (asc) / first (desc).
    items.push({ id: 'id-zzz', name: 'zzz-name' });
    this.model = items;
    this.sortBy = undefined;
    this.sortOrder = undefined;
    this.config = { ...this.config, columns: [{ key: 'name', label: 'Name', isSortable: true }] };

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
        @sortBy={{this.sortBy}}
        @sortOrder={{this.sortOrder}}
      />
    `);

    // Before sorting: 13 items, page 1 shows 10. "zzz-name" is on page 2.
    assert.dom(GENERAL.tableRow()).exists({ count: 10 }, 'page 1 shows 10 rows before sort');

    // Simulate the controller setting sortBy/sortOrder (as if the user clicked the header).
    this.set('sortBy', 'name');
    this.set('sortOrder', 'desc');

    // After desc sort the first row on page 1 must be "zzz-name" (global maximum).
    assert
      .dom(`${GENERAL.tableRow(0)} [data-test-table-data="name"]`)
      .hasText('zzz-name', 'first row on page 1 is the globally largest name after descending sort');
  });

  // ── Test 16 — sort state persists across page transitions ────────────────────
  test('sort column and direction persist when the page arg changes', async function (assert) {
    // Simulates what happens when the user navigates to page 2 via the
    // pagination control and the route rerenders with a new @page value.
    // The sort indicator must remain active (blue arrow) rather than resetting.
    const items = Array.from({ length: 12 }, (_, i) => ({
      id: `id-${i}`,
      name: `item-${String(i).padStart(2, '0')}`,
    }));
    this.model = items;
    this.config = { ...this.config, columns: [{ key: 'name', label: 'Name', isSortable: true }] };

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    // Sort ascending — sort button should become active.
    await click(GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true }));

    const sortButtonSel = GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true });

    // Confirm sort indicator is active (HDS adds aria-sort="ascending" to the th).
    assert
      .dom('.hds-advanced-table__th:nth-child(1)')
      .hasAttribute('aria-sort', 'ascending', 'sort indicator is ascending before page change');

    // Simulate a page transition — Page::ListView receives a new @page arg from
    // the URL query param update.  The component instance (and its @tracked
    // sortColumn/sortDirection) must survive because it is not destroyed on a
    // simple arg change.
    this.set('page', 2);

    // Sort indicator must still be active.
    assert
      .dom('.hds-advanced-table__th:nth-child(1)')
      .hasAttribute('aria-sort', 'ascending', 'sort indicator stays active after navigating to page 2');

    // The sort button element must still be present (header did not reset).
    assert.dom(sortButtonSel).exists('sort button still rendered on page 2');
  });

  // ── Test 16 — natural / numeric alphanumeric sorting ──────────────────────────
  test('natural sorting orders numeric segments numerically rather than lexicographically', async function (assert) {
    this.model = [
      { id: '1', name: 'pebble-docker-18' },
      { id: '2', name: 'pebble-docker-19' },
      { id: '3', name: 'pebble-docker-2' },
      { id: '4', name: 'pebble-docker-20' },
      { id: '5', name: 'pebble-docker-3' },
    ];
    this.config = { ...this.config, columns: [{ key: 'name', label: 'Name', isSortable: true }] };

    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{1}}
        @sortBy="name"
        @sortOrder="asc"
      />
    `);

    // Under natural sorting: pebble-docker-2, pebble-docker-3, pebble-docker-18, pebble-docker-19, pebble-docker-20
    assert.dom(`${GENERAL.tableRow(0)} [data-test-table-data="name"]`).hasText('pebble-docker-2');
    assert.dom(`${GENERAL.tableRow(1)} [data-test-table-data="name"]`).hasText('pebble-docker-3');
    assert.dom(`${GENERAL.tableRow(2)} [data-test-table-data="name"]`).hasText('pebble-docker-18');
    assert.dom(`${GENERAL.tableRow(3)} [data-test-table-data="name"]`).hasText('pebble-docker-19');
    assert.dom(`${GENERAL.tableRow(4)} [data-test-table-data="name"]`).hasText('pebble-docker-20');
  });
});
