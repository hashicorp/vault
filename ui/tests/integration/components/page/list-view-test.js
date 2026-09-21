/**
 * Copyright IBM Corp. 2016, 2025
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

  // ── Test 11 — filter input triggers router transition only when page > 1 ──────
  test('it resets page to 1 via router.transitionTo when the filter input changes on page > 1', async function (assert) {
    // The transition is skipped when already on page 1 to avoid spurious model
    // re-fetches. It only fires when the user is on a later page so the QP is
    // actually reset to 1.
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.config = {
      ...this.config,
      filter: { type: 'text', placeholder: 'Filter entities', ariaLabel: 'Filter', filterKey: 'name' },
    };

    // Start on page 2 — transition should fire.
    this.page = 2;
    await render(hbs`
      <Page::ListView
        @config={{this.config}}
        @model={{this.model}}
        @page={{this.page}}
      />
    `);

    await fillIn(GENERAL.filterInput, 'alice');

    assert.ok(transitionStub.calledOnce, 'router.transitionTo was called once when starting from page 2');
    assert.deepEqual(
      transitionStub.firstCall.args,
      [{ queryParams: { page: 1 } }],
      'transitions with page reset to 1 (filter value is internal, not a QP)'
    );
  });

  // ── Test 11b — filter on page 1 does NOT trigger a router transition ──────────
  test('it does not call router.transitionTo when the filter changes on page 1', async function (assert) {
    // Typing while already on page 1 must not trigger a route transition — doing
    // so would cause a model re-fetch and a visible "page refresh" on every keystroke.
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

    assert.ok(transitionStub.notCalled, 'router.transitionTo was not called when already on page 1');
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
});
