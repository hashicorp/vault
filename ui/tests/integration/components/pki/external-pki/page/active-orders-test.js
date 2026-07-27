/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module(
  'Integration | Component | pki | external-pki | ExternalPki::Page::RolesRoleActiveOrders',
  function (hooks) {
    setupRenderingTest(hooks);
    setupEngine(hooks, 'pki');

    hooks.beforeEach(function () {
      this.model = {
        engine: { id: 'pki-external-ca' },
        activeOrders: [],
        role_name: 'myrole',
      };
      this.router = this.owner.lookup('service:router');
      this.renderComponent = () =>
        render(hbs`<ExternalPki::Page::RolesRoleActiveOrders @model={{this.model}} />`, {
          owner: this.engine,
        });
    });

    test('it renders empty state when no active orders', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.emptyStateTitle).exists().hasText('No active orders');
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText(
          'In progress orders will appear here once created. Lookup a specific order by its ID or navigate to Recent orders to view recently created and completed orders. Lookup order'
        );
      assert.dom(GENERAL.linkTo('API docs: Create a new order')).exists();
      assert.dom(GENERAL.inputSearch('Filter by order ID')).doesNotExist();
      assert.dom(GENERAL.button('Refresh')).doesNotExist();
    });

    test('it renders list of active orders', async function (assert) {
      this.model.activeOrders = ['order-123', 'order-456', 'order-789'];

      await this.renderComponent();
      assert.dom(GENERAL.inputSearch('Filter by order ID')).exists('search input is rendered');
      assert.dom(GENERAL.button('Refresh')).exists();
      assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'displays all orders');
      assert.dom(GENERAL.linkTo('order-123')).exists().hasText('order-123');
      assert.dom(GENERAL.linkTo('order-456')).exists().hasText('order-456');
      assert.dom(GENERAL.linkTo('order-789')).exists().hasText('order-789');
    });

    test('it filters orders by search input', async function (assert) {
      this.model.activeOrders = ['order-123', 'order-456', 'order-789', 'test-abc-order'];
      await this.renderComponent();
      assert.dom(GENERAL.listItem()).exists({ count: 4 }, 'initially displays all orders');
      // Filter by "456"
      await fillIn(GENERAL.inputSearch('Filter by order ID'), '456');
      assert.dom(GENERAL.listItem()).exists({ count: 1 }, 'displays only matching order');
      assert.dom(GENERAL.pagination).hasTextContaining('1–1 of 1');
      // Filter by "order-"
      await fillIn(GENERAL.inputSearch('Filter by order ID'), 'order-');
      assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'display 3 matching orders');
      assert.dom(GENERAL.linkTo('order-123')).exists().hasText('order-123');
      assert.dom(GENERAL.linkTo('order-456')).exists().hasText('order-456');
      assert.dom(GENERAL.linkTo('order-789')).exists().hasText('order-789');
      assert.dom(GENERAL.pagination).hasTextContaining('1–3 of 3');
      // Filter by ORDER-1
      await fillIn(GENERAL.inputSearch('Filter by order ID'), 'ORDER-1');
      assert.dom(GENERAL.listItem()).exists({ count: 1 }, 'display 1 matching orders');
      assert.dom(GENERAL.linkTo('order-123')).hasText('order-123', 'filter is case INsensitive');
      // Clear search input
      await fillIn(GENERAL.inputSearch('Filter by order ID'), '');
      assert.dom(GENERAL.listItem()).exists({ count: 4 }, 'shows all orders again');
      assert.dom(GENERAL.pagination).hasTextContaining('1–4 of 4');
    });

    test('it shows empty state when search has no results', async function (assert) {
      this.model.activeOrders = ['order-123', 'order-456'];

      await this.renderComponent();
      await fillIn(GENERAL.inputSearch('Filter by order ID'), 'nope');
      assert.dom(GENERAL.inputSearch('Filter by order ID')).exists();
      assert.dom(GENERAL.button('Refresh')).exists();
      assert.dom(GENERAL.listItem()).doesNotExist();
      assert.dom(GENERAL.emptyStateTitle).exists().hasText('No orders matching ID: nope');
    });

    test('it navigates to order when lookup by ID button is clicked', async function (assert) {
      this.model.activeOrders = [];
      const transitionStub = sinon.stub(this.router, 'transitionTo');
      await this.renderComponent();
      await fillIn(GENERAL.inputSearch('orderId'), 'test-order-123');
      await click(GENERAL.button('Lookup order'));
      assert.true(transitionStub.calledOnce, 'transitionTo was called once');
      const [route, mount, role, orderId] = transitionStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.pki.external.roles.role.order',
        'transitionTo called with route'
      );
      assert.strictEqual(mount, 'pki-external-ca', 'transitionTo called with mount');
      assert.strictEqual(role, 'myrole', 'transitionTo called with role');
      assert.strictEqual(orderId, 'test-order-123', 'transitionTo called with orderId');
    });

    test('it calls refresh when refresh button is clicked', async function (assert) {
      this.model.activeOrders = ['order-123', 'order-456'];
      const refreshStub = sinon.stub(this.router, 'refresh');
      await this.renderComponent();
      await click(GENERAL.button('Refresh'));
      assert.true(refreshStub.calledOnce, 'refresh was called once');
      const [route] = refreshStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.pki.external.roles.role.active-orders',
        'refresh was called with correct route'
      );
    });
  }
);
