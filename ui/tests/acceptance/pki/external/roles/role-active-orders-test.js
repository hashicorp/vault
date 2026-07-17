/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { currentRouteName, currentURL, visit, click } from '@ember/test-helpers';
import sinon from 'sinon';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Acceptance | enterprise | pki | external | roles | role | active-orders route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.activeOrdersListStub = sinon.stub(api.secrets, 'pkiExternalCaListRoleActiveOrders');
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    this.roleName = 'test-role';

    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.activeOrdersURL = `/vault/secrets-engines/${this.mountPath}/pki/external/roles/${this.roleName}/active-orders`;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it renders breadcrumbs for role active orders', async function (assert) {
    await visit(this.activeOrdersURL);
    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.roleName);
    assert.dom(GENERAL.breadcrumb).exists({ count: 5 });
    assert.dom(GENERAL.breadcrumbs).hasText(`Vault Secrets engines ${this.mountPath} Roles ${this.roleName}`);
    assert.dom(GENERAL.linkTo('Active orders')).exists().hasClass('active');
    assert.dom(GENERAL.linkTo('Details')).exists().doesNotHaveClass('active');

    // Navigate to a role details
    await click(GENERAL.linkTo('Details'));
    assert.dom(GENERAL.linkTo('Details')).exists().hasClass('active');
    assert.dom(GENERAL.linkTo('Active orders')).exists().doesNotHaveClass('active');
  });

  test('it fetches and displays active orders', async function (assert) {
    this.activeOrdersListStub.resolves({
      keys: ['order-abc123', 'order-def456', 'order-ghi789'],
    });

    await visit(this.activeOrdersURL);

    assert.strictEqual(currentURL(), this.activeOrdersURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.active-orders',
      'navigated to active orders route'
    );
    assert.true(this.activeOrdersListStub.calledOnce, 'active orders list called once');

    // Verify orders are displayed
    assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'displays all orders');
    assert.dom(GENERAL.linkTo('order-abc123')).exists();
    assert.dom(GENERAL.linkTo('order-def456')).exists();
    assert.dom(GENERAL.linkTo('order-ghi789')).exists();
  });

  test('it handles empty orders list (404)', async function (assert) {
    this.activeOrdersListStub.rejects(getErrorResponse());
    await visit(this.activeOrdersURL);
    assert.strictEqual(currentURL(), this.activeOrdersURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.active-orders',
      'stays on active orders route'
    );
    assert.true(this.activeOrdersListStub.calledOnce, 'active orders list called once');
    assert.dom('h1').hasText(this.roleName, 'role name is displayed');
    assert.dom(GENERAL.emptyStateTitle).exists().hasText('No active orders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'In progress orders will appear here once created. Lookup a specific order by its ID or navigate to Recent orders to view recently created and completed orders. Lookup order'
      );
    assert.dom(GENERAL.linkTo('API docs: Create a new order')).exists();
  });

  test('it handles 403 permission denied error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.activeOrdersListStub.rejects(getErrorResponse(error, 403));

    await visit(this.activeOrdersURL);

    assert.strictEqual(currentURL(), this.activeOrdersURL, 'it has expected URL');
    assert.true(this.activeOrdersListStub.calledOnce, 'active orders list called once');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.error',
      'redirects to role error route'
    );
    assert.dom('h1').hasText(this.roleName, 'role name is displayed');
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');
  });

  test('it handles 500 internal server error', async function (assert) {
    const error = { errors: ['Internal server error'] };
    this.activeOrdersListStub.rejects(getErrorResponse(error, 500));

    await visit(this.activeOrdersURL);

    assert.strictEqual(currentURL(), this.activeOrdersURL, 'it has expected URL');
    assert.true(this.activeOrdersListStub.calledOnce, 'active orders list called once');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.error',
      'redirects to role error route'
    );
    assert.dom('h1').hasText(this.roleName, 'role name is displayed');
    assert.dom(GENERAL.pageError.title(500)).exists().hasText('ERROR 500 Error');
  });

  test('it navigates to individual order details', async function (assert) {
    this.activeOrdersListStub.resolves({
      keys: ['order-abc123'],
    });
    await visit(this.activeOrdersURL);
    await click(GENERAL.linkTo('order-abc123'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/orders/order-abc123`,
      'navigates to order details'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.orders.order',
      'transitions to order route'
    );
  });
});
