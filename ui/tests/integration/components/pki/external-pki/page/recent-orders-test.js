/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const SELECTORS = {
  filterTag: '[data-test-filter-tag]',
};
module('Integration | Component | pki | external-pki | ExternalPki::Page::RecentOrders', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.model = { engine: { id: 'pki-external-ca' }, query: { within: '24h' } };

    this.renderComponent = () =>
      render(hbs`<ExternalPki::Page::RecentOrders @model={{this.model}} />`, { owner: this.engine });
  });

  test('it renders empty state when no recent orders', async function (assert) {
    this.model.recentOrders = [];
    await this.renderComponent();
    // Filter toolbar state
    assert
      .dom(GENERAL.inputSearch('Filter by order ID'))
      .exists()
      .hasAttribute('placeholder', 'No orders to filter')
      .hasAttribute('aria-label', 'No orders to filter');
    assert.dom(GENERAL.dropdownToggle('Role')).exists().isDisabled();
    assert.dom(GENERAL.dropdownToggle('Status')).exists().isDisabled();
    assert
      .dom(GENERAL.dropdownToggle('Created in last'))
      .exists()
      .isNotDisabled()
      .hasText('Created in last: day');
    await click(GENERAL.dropdownToggle('Created in last'));
    assert.dom(GENERAL.menuItem('1 day')).hasAttribute('aria-selected', 'true');
    assert.dom(GENERAL.button('Refresh')).exists();
    assert.dom(GENERAL.emptyStateTitle).hasText('No recent orders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .includesText(
        'No orders have been created in the last day (24h). Select a different time period or lookup an archived order by its ID.'
      );
    assert.dom(GENERAL.inputSearch('orderId')).exists().hasAttribute('placeholder', '01936d8e-7c3...');
    assert.dom(GENERAL.button('Lookup order')).exists();
  });

  test('it shows empty state when filters have no results', async function (assert) {
    this.model.recentOrders = [
      {
        order_id: 'order-123',
        role_name: 'dev-server',
        order_status: 'submitted',
        identifiers: 'example.com',
        creation_date: '2026-07-14T20:00:00Z',
        last_update: '2026-07-14T20:05:00Z',
      },
      {
        order_id: 'order-456',
        role_name: 'web-server',
        order_status: 'completed',
        identifiers: 'example.com',
        creation_date: '2026-07-14T20:00:00Z',
        last_update: '2026-07-14T20:05:00Z',
      },
      {
        order_id: 'order-456',
        role_name: 'prod-server',
        order_status: 'expired',
        identifiers: 'example.com',
        creation_date: '2026-07-14T20:00:00Z',
        last_update: '2026-07-14T20:05:00Z',
      },
    ];
    await this.renderComponent();
    await fillIn(GENERAL.inputSearch('Filter by order ID'), '789');
    assert.dom(GENERAL.emptyStateTitle).hasText('No recent orders matching: 789');
    // Clear order ID, select dropdown filters to test their empty state
    await fillIn(GENERAL.inputSearch('Filter by order ID'), '');
    await click(GENERAL.dropdownToggle('Role'));
    await click(GENERAL.menuItem('web-server'));
    await click(GENERAL.dropdownToggle('Status'));
    await click(GENERAL.menuItem('Pending'));
    assert.dom(GENERAL.emptyStateTitle).hasText('No matching orders');
    assert.dom(GENERAL.emptyStateMessage).hasText('Clear or update filters to view recent orders.');
  });

  test('it navigates to order when lookup by ID button is clicked', async function (assert) {
    this.model.recentOrders = [];
    const transitionStub = sinon.stub(this.router, 'transitionTo');
    await this.renderComponent();
    await fillIn(GENERAL.inputSearch('orderId'), 'test-order-123');
    await click(GENERAL.button('Lookup order'));
    assert.true(transitionStub.calledOnce, 'transitionTo was called once');
    assert.true(
      transitionStub.calledWith(
        'vault.cluster.secrets.backend.pki.external.orders.order',
        'pki-external-ca',
        'test-order-123'
      ),
      'transitionTo was called with correct arguments'
    );
  });

  test('it renders time query dropdown', async function (assert) {
    this.model.recentOrders = [];
    this.model.query = { within: '72h' };
    await this.renderComponent();
    assert.dom(GENERAL.dropdownToggle('Created in last')).hasText('Created in last: 3 days');
    await click(GENERAL.dropdownToggle('Created in last'));
    assert.dom(GENERAL.dropdownToggle('Created in last')).hasAttribute('aria-expanded', 'true');
    assert.dom(GENERAL.menuItem('1 hour')).exists();
    assert.dom(GENERAL.menuItem('1 day')).exists();
    assert.dom(GENERAL.menuItem('3 days')).exists().hasAttribute('aria-selected', 'true');
    assert.dom(GENERAL.menuItem('5 days')).exists();
    assert.dom(GENERAL.menuItem('1 week')).exists();
  });

  test('it does not remove "1" from formatted query for multiple units', async function (assert) {
    this.model.recentOrders = [];
    this.model.query = { within: '26h' };
    await this.renderComponent();
    assert
      .dom(GENERAL.dropdownToggle('Created in last'))
      .hasText(
        'Created in last: 1 day 2 hours',
        '"1" is not stripped from formatted string when followed by another unit'
      );
  });

  module('with orders', function (hooks) {
    hooks.beforeEach(function () {
      // Helper function to generate orders from status list
      let count = 100;
      const generateOrders = (statuses) => {
        const roles = ['web-server', 'api-server'];
        const baseTime = new Date('2026-07-14T20:00:00Z');

        return statuses.map((status, index) => {
          const minutesOffset = index * 10;
          const creationDate = new Date(baseTime.getTime() - minutesOffset * 60000);
          const lastUpdate = new Date(creationDate.getTime() + 5 * 60000);

          return {
            order_id: `order-${count++}`,
            role_name: roles[index % 2],
            order_status: status,
            identifiers: `${status}.example.com`,
            creation_date: creationDate.toISOString(),
            last_update: lastUpdate.toISOString(),
          };
        });
      };

      // All possible order statuses
      this.orderStatuses = [
        'new',
        'submitted',
        'awaiting-challenge-fulfillment',
        'vault-challenge-fulfillment',
        'vault-challenge-propagating',
        'notify-acme-server-challenges-completed',
        'processing-challenge',
        'fetching-certificate',
        'completed',
        'revoked',
        'expired',
        'error',
        'undefined',
      ];

      this.model.recentOrders = generateOrders(this.orderStatuses);
    });

    test('it renders list of recent orders', async function (assert) {
      await this.renderComponent();
      assert
        .dom(GENERAL.inputSearch('Filter by order ID'))
        .exists()
        .hasAttribute('placeholder', 'Filter by order ID')
        .hasAttribute('aria-label', 'Filter by order ID');
      assert.dom(GENERAL.dropdownToggle('Role')).exists().isNotDisabled();
      assert.dom(GENERAL.dropdownToggle('Status')).exists().isNotDisabled();
      assert.dom(GENERAL.dropdownToggle('Created in last')).exists();
      assert.dom(GENERAL.button('Refresh')).exists();
      assert.dom(GENERAL.listItem()).exists({ count: 10 }, 'displays first 10 orders');
      assert.dom(GENERAL.pagination).hasTextContaining('1–10 of 13');
      // Assert table row renders values as expected
      assert
        .dom(`${GENERAL.tableData(0, 'order_id')} [data-test-link-to]`)
        .exists('it renders link for order_id');
      assert
        .dom(`${GENERAL.tableData(0, 'role_name')} [data-test-link-to]`)
        .exists('it renders link for role_name');
      assert
        .dom(`${GENERAL.tableData(0, 'order_status')} ${GENERAL.badge('order_status')}`)
        .exists()
        .hasText('Pending');
      assert.dom(GENERAL.tableData(0, 'creation_date')).exists().hasTextContaining('07/14/2026');
      assert.dom(GENERAL.tableData(0, 'last_update')).exists().hasTextContaining('07/14/2026');
    });

    test('it shows filter tags container', async function (assert) {
      await this.renderComponent();
      assert.dom('[data-test-filter-tag-container]').exists();
      assert.dom('[data-test-filter-tag-container]').includesText('Filters applied:');
      assert.dom('[data-test-filter-tag-container]').includesText('None');
    });

    test('it filters by status and maps API statuses', async function (assert) {
      await this.renderComponent();
      const expectedMappedStatuses = { Pending: 8, Failed: 1, Expired: 1, Revoked: 1, Issued: 1, Unknown: 1 };
      await click(GENERAL.dropdownToggle('Status'));
      assert.dom('[data-test-popup-menu]').exists({ count: 6 });
      // Close menu again because loop below re-opens
      await click(GENERAL.dropdownToggle('Status'));
      for (const status in expectedMappedStatuses) {
        const count = expectedMappedStatuses[status];
        await click(GENERAL.dropdownToggle('Status'));
        assert.dom(GENERAL.dropdownToggle('Status')).hasAttribute('aria-expanded', 'true');
        assert.dom(GENERAL.menuItem(status)).exists(`status dropdown includes ${status}`);
        await click(GENERAL.menuItem(status));
        assert.dom(SELECTORS.filterTag).hasText(status);
        assert.dom(GENERAL.listItem()).exists({ count }, `table renders expected row number: ${count}`);
        assert
          .dom(GENERAL.dropdownToggle('Status'))
          .hasAttribute('aria-expanded', 'false', 'selecting a status closes the dropdown');
      }
    });

    test('it clears status filter when tag is dismissed', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.dropdownToggle('Status'));
      await click(GENERAL.menuItem('Pending'));
      assert.dom(GENERAL.listItem()).exists({ count: 8 });
      assert.dom(SELECTORS.filterTag).hasText('Pending');
      // Dismiss tag
      await click(`${SELECTORS.filterTag} button`);
      assert.dom(GENERAL.listItem()).exists({ count: 10 }, 'shows all orders again');
      assert.dom(SELECTORS.filterTag).doesNotExist();
    });

    test('it filters orders by order ID search', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.listItem()).exists({ count: 10 }, 'displays first 10 orders');
      assert.dom(GENERAL.pagination).hasTextContaining('1–10 of 13');
      await fillIn(GENERAL.inputSearch('Filter by order ID'), '101');
      await waitFor(GENERAL.listItem());
      assert.dom(SELECTORS.filterTag).doesNotExist();
      assert.dom(GENERAL.listItem()).exists({ count: 1 }, 'displays only matching order');
      assert.dom(GENERAL.linkTo('order-101')).exists();
      await fillIn(GENERAL.inputSearch('Filter by order ID'), 'order-11');
      await waitFor(GENERAL.listItem());
      assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'displays 3 matching orders');
      assert.dom(GENERAL.linkTo('order-110')).exists();
      assert.dom(GENERAL.linkTo('order-111')).exists();
      assert.dom(GENERAL.linkTo('order-112')).exists();

      await fillIn(GENERAL.inputSearch('Filter by order ID'), '');
      await waitFor(GENERAL.listItem());
      assert.dom(GENERAL.listItem()).exists({ count: 10 }, 'displays first 10 orders');
      assert.dom(GENERAL.pagination).hasTextContaining('1–10 of 13');
    });

    // Order IDs ARE case sensitive, but for filtering purposes we don't need to be that strict
    test('filter is case insensitive', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.inputSearch('Filter by order ID'), 'ORDER-101');
      assert.dom(GENERAL.listItem()).exists({ count: 1 }, 'displays 1 matching order');
      assert.dom(GENERAL.linkTo('order-101')).exists();
    });

    test('it filters orders by role', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.dropdownToggle('Role'));
      assert.dom(GENERAL.dropdownToggle('Role')).hasAttribute('aria-expanded', 'true');
      assert.dom(GENERAL.menuItem('web-server')).exists();
      assert.dom(GENERAL.menuItem('api-server')).exists();
      await click(GENERAL.menuItem('web-server'));
      assert
        .dom(GENERAL.dropdownToggle('Role'))
        .hasAttribute('aria-expanded', 'false', 'dropdown closes after selecting role');
      assert.dom(GENERAL.listItem()).exists({ count: 7 }, 'displays 7 web-server orders');
      assert.dom(SELECTORS.filterTag).hasText('web-server');
      // Check a few rows for accuracy
      assert.dom(GENERAL.linkTo('order-100')).exists();
      assert.dom(GENERAL.linkTo('order-102')).exists();
      assert.dom(GENERAL.linkTo('order-103')).doesNotExist();
      // Select other role to make sure list updates
      await click(GENERAL.dropdownToggle('Role'));
      await click(GENERAL.menuItem('api-server'));
      assert.dom(GENERAL.listItem()).exists({ count: 6 }, 'displays 6 api-server orders');
      // Check a few rows for accuracy
      assert.dom(SELECTORS.filterTag).hasText('api-server');
      assert.dom(GENERAL.linkTo('order-101')).exists();
      assert.dom(GENERAL.linkTo('order-103')).exists();
      assert.dom(GENERAL.linkTo('order-104')).doesNotExist();
    });

    test('it clears role filter when tag is dismissed', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.dropdownToggle('Role'));
      await click(GENERAL.menuItem('web-server'));
      assert.dom(SELECTORS.filterTag).hasText('web-server');
      await click(`${SELECTORS.filterTag} button`);
      assert.dom(GENERAL.listItem()).exists({ count: 10 }, 'shows all orders again');
      assert.dom(SELECTORS.filterTag).doesNotExist();
    });

    test('it searches role names in dropdown', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.dropdownToggle('Role'));
      assert.dom(GENERAL.menuItem('web-server')).exists();
      assert.dom(GENERAL.menuItem('api-server')).exists();

      await fillIn('#roleNameSearch', 'web');
      assert.dom(GENERAL.menuItem('web-server')).exists();
      assert.dom(GENERAL.menuItem('api-server')).doesNotExist();

      await fillIn('#roleNameSearch', 'api');
      assert.dom(GENERAL.menuItem('web-server')).doesNotExist();
      assert.dom(GENERAL.menuItem('api-server')).exists();
    });

    test('it updates dropdown when query param changes', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.dropdownToggle('Created in last')).hasText('Created in last: day');
      await click(GENERAL.dropdownToggle('Created in last'));
      assert.dom(GENERAL.dropdownToggle('Created in last')).hasAttribute('aria-expanded', 'true');
      assert.dom(GENERAL.menuItem('1 hour')).exists();
      assert.dom(GENERAL.menuItem('1 day')).exists().hasAttribute('aria-selected', 'true');
      assert.dom(GENERAL.menuItem('3 days')).exists();
      assert.dom(GENERAL.menuItem('5 days')).exists();
      assert.dom(GENERAL.menuItem('1 week')).exists();
      await click(GENERAL.dropdownToggle('Created in last'));
      // Change query param
      this.set('model', { ...this.model, query: { within: '78h' } });
      assert.dom(GENERAL.dropdownToggle('Created in last')).hasText('Created in last: 3 days 6 hours');
    });

    test('it calls refresh when refresh button is clicked', async function (assert) {
      this.model.recentOrders = [
        {
          order_id: 'order-123',
          role_name: 'web-server',
          order_status: 'pending',
          identifiers: 'example.com',
          creation_date: '2026-07-14T20:00:00Z',
          last_update: '2026-07-14T20:05:00Z',
        },
      ];

      const refreshStub = sinon.stub(this.router, 'refresh');
      await this.renderComponent();
      await click(GENERAL.button('Refresh'));
      assert.true(refreshStub.calledOnce, 'refresh was called once');
      assert.true(
        refreshStub.calledWith('vault.cluster.secrets.backend.pki.external.orders'),
        'refresh was called with correct route'
      );
    });
  });

  module('combined filters', function (hooks) {
    hooks.beforeEach(function () {
      this.model.recentOrders = [
        {
          order_id: 'order-123',
          role_name: 'web-server',
          order_status: 'submitted',
          identifiers: 'example.com',
          creation_date: '2026-07-14T20:00:00Z',
          last_update: '2026-07-14T20:05:00Z',
        },
        {
          order_id: 'order-456',
          role_name: 'api-server',
          order_status: 'completed',
          identifiers: 'api.example.com',
          creation_date: '2026-07-14T19:00:00Z',
          last_update: '2026-07-14T19:30:00Z',
        },
        {
          order_id: 'order-789',
          role_name: 'web-server',
          order_status: 'completed',
          identifiers: 'test.example.com',
          creation_date: '2026-07-14T18:00:00Z',
          last_update: '2026-07-14T18:15:00Z',
        },
      ];
    });

    test('it applies multiple filters together', async function (assert) {
      await this.renderComponent();

      // Apply role filter
      await click(GENERAL.dropdownToggle('Role'));
      await click(GENERAL.menuItem('web-server'));
      assert.dom(GENERAL.listItem()).exists({ count: 2 });
      assert.dom(SELECTORS.filterTag).exists({ count: 1 });

      // Apply status filter
      await click(GENERAL.dropdownToggle('Status'));
      await click(GENERAL.menuItem('Issued'));
      assert.dom(GENERAL.listItem()).exists({ count: 1 }, 'only one order matches both filters');
      assert.dom(GENERAL.linkTo('order-789')).exists();
      assert.dom(SELECTORS.filterTag).exists({ count: 2 });

      // Apply order ID search
      await fillIn(GENERAL.inputSearch('Filter by order ID'), '789');
      await waitFor(GENERAL.listItem());
      assert.dom(GENERAL.listItem()).exists({ count: 1 });
      assert.dom(GENERAL.linkTo('order-789')).exists();
      assert.dom(SELECTORS.filterTag).exists({ count: 2 });
    });

    test('it shows clear all filters button when multiple filters are applied', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.button('Clear filters')).doesNotExist();
      await click(GENERAL.dropdownToggle('Role'));
      await click(GENERAL.menuItem('web-server'));
      assert.dom(GENERAL.button('Clear filters')).doesNotExist('not shown with only one filter');
      await click(GENERAL.dropdownToggle('Status'));
      await click(GENERAL.menuItem('Issued'));
      assert.dom(GENERAL.button('Clear filters')).exists('shown with multiple filters');
    });

    test('it clears all filters when clear all button is clicked', async function (assert) {
      await this.renderComponent();
      await click(GENERAL.dropdownToggle('Role'));
      await click(GENERAL.menuItem('web-server'));
      await click(GENERAL.dropdownToggle('Status'));
      await click(GENERAL.menuItem('Issued'));
      assert.dom(SELECTORS.filterTag).exists({ count: 2 });
      assert.dom(GENERAL.listItem()).exists({ count: 1 });
      await click(GENERAL.button('Clear filters'));
      assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'shows all orders again');
      assert.dom(SELECTORS.filterTag).doesNotExist();
    });
  });
});
