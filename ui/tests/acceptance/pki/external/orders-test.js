/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentRouteName, currentURL, visit } from '@ember/test-helpers';
import sinon from 'sinon';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Acceptance | enterprise | pki | external | orders route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    this.api = this.owner.lookup('service:api');
    this.recentOrdersListStub = sinon.stub(this.api.secrets, 'pkiExternalCaListLookupOrdersRecent');
    this.mountPath = `pki-external-ca-${uuidv4()}`;

    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.ordersURL = `/vault/secrets-engines/${this.mountPath}/pki/external/orders`;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it renders breadcrumbs for recent orders', async function (assert) {
    this.recentOrdersListStub.resolves({
      keys: [],
      key_info: {},
    });

    await visit(this.ordersURL);

    assert.dom(GENERAL.breadcrumb).exists({ count: 4 });
    assert.dom(GENERAL.breadcrumbs).hasText(`Vault Secrets engines ${this.mountPath} Recent orders`);
    assert.dom(GENERAL.breadcrumbLink('Vault')).exists();
    assert.dom(GENERAL.breadcrumbLink('Secrets engines')).exists();
    assert.dom(GENERAL.breadcrumbLink(this.mountPath)).exists();
    assert.dom(GENERAL.currentBreadcrumb('Recent orders')).exists();
  });

  test('it sets default query param when not provided', async function (assert) {
    const apiSpy = sinon.spy(this.api, 'addQueryParams');
    // Restore stub so we can assert addQueryParams is called with expected query
    this.recentOrdersListStub.restore();
    await visit(this.ordersURL);
    assert.strictEqual(currentURL(), this.ordersURL, 'it navigates to url');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.orders.index',
      'navigated to orders index route'
    );
    assert.dom(GENERAL.dropdownToggle('Created in last')).exists().hasText('Created in last: hour');
    const [, query] = apiSpy.lastCall.args;
    assert.propEqual(query, { within: '1h' }, 'request is made with default query');
  });

  test('it respects provided query param', async function (assert) {
    this.recentOrdersListStub.resolves({
      keys: [],
      key_info: {},
    });
    await visit(`${this.ordersURL}?within=24h`);
    assert.strictEqual(currentURL(), `${this.ordersURL}?within=24h`, 'URL includes within=24h query param');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.orders.index',
      'navigated to orders index route'
    );
    assert.true(this.recentOrdersListStub.calledOnce, 'recent orders list request is made');
    assert.dom(GENERAL.dropdownToggle('Created in last')).exists().hasText('Created in last: day');
  });

  test('it updates query param when time period selected from dropdown', async function (assert) {
    this.recentOrdersListStub.resolves({
      keys: [],
      key_info: {},
    });
    await visit(`${this.ordersURL}?within=26h`);
    assert.strictEqual(currentURL(), `${this.ordersURL}?within=26h`, 'initial URL has within=26h');
    assert.true(this.recentOrdersListStub.calledOnce, 'initial API call made');
    assert.dom(GENERAL.dropdownToggle('Created in last')).exists().hasText('Created in last: 1 day 2 hours');
    // Select a different query from dropdown
    await click(GENERAL.dropdownToggle('Created in last'));
    await click(GENERAL.menuItem('1 week'));
    assert
      .dom(GENERAL.dropdownToggle('Created in last'))
      .hasAttribute('aria-expanded', 'false', 'dropdown closes after selecting a query');
    assert.strictEqual(currentURL(), `${this.ordersURL}?within=168h`, 'URL updated to within=168h');
    assert.true(this.recentOrdersListStub.calledTwice, 'API called again with new param');
    assert.dom(GENERAL.dropdownToggle('Created in last')).exists().hasText('Created in last: 7 days');
  });

  test('it handles empty orders list (404)', async function (assert) {
    this.recentOrdersListStub.rejects(getErrorResponse());
    await visit(this.ordersURL);
    assert.strictEqual(currentURL(), this.ordersURL, 'it navigates to url');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.orders.index',
      'stays on orders index route'
    );
    assert.true(this.recentOrdersListStub.calledOnce, 'recent orders list called once');
    assert.dom(GENERAL.emptyStateTitle).exists().hasText('No recent orders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'No orders have been created in the last hour (1h). Select a different time period or lookup an archived order by its ID.'
      );
  });

  test('it handles 403 permission denied error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.recentOrdersListStub.rejects(getErrorResponse(error, 403));
    await visit(this.ordersURL);
    assert.strictEqual(currentURL(), this.ordersURL, 'navigates to URL');
    assert.true(this.recentOrdersListStub.calledOnce, 'recent orders list called once');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'redirects to external error route'
    );
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');
  });

  test('it handles 500 internal server error', async function (assert) {
    const error = { errors: ['Internal server error'] };
    this.recentOrdersListStub.rejects(getErrorResponse(error, 500));
    await visit(this.ordersURL);
    assert.strictEqual(currentURL(), this.ordersURL, 'navigates to URL');
    assert.true(this.recentOrdersListStub.calledOnce, 'recent orders list called once');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'redirects to external error route'
    );
    assert.dom(GENERAL.pageError.title(500)).exists().hasText('ERROR 500 Error');
  });
});
