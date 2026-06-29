/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { currentRouteName, currentURL, visit } from '@ember/test-helpers';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';

const TABS = ['Overview', 'Roles', 'Recent orders', 'DNS providers', 'ACME accounts'];

module('Acceptance | enterprise pki external navigation', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    // Setup External PKI engine
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.engineURL = `vault/secrets-engines/${this.mountPath}/pki/external`;
    this.baseCrumbs = `Vault Secrets engines ${this.mountPath}`;
    this.assertTabState = (assert, activeTab) => {
      const inactive = TABS.filter((t) => t !== activeTab);
      inactive.forEach((t) => {
        assert.dom(GENERAL.linkTo(t)).exists().doesNotHaveClass('active', `${t} is inactive`);
      });
      assert.dom(GENERAL.linkTo(activeTab)).exists().hasClass('active', `${activeTab} is active`);
    };
  });

  test('it navigates to external overview', async function (assert) {
    await visit(this.engineURL);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/overview`,
      'it navigates to overview'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.overview',
      'navigating to pki.external.index redirects to overview'
    );
    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.mountPath);
    assert.dom(GENERAL.breadcrumb).exists({ count: 3 });
    assert.dom(GENERAL.breadcrumbs).hasText(this.baseCrumbs);
    this.assertTabState(assert, 'Overview');
  });

  test('it navigates to external roles', async function (assert) {
    await visit(`${this.engineURL}/roles`);
    assert.strictEqual(
      currentURL(),
      `vault/secrets-engines/${this.mountPath}/pki/external/roles`,
      'it navigates to roles'
    );
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.external.roles.index');

    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.mountPath);
    assert.dom(GENERAL.breadcrumb).exists({ count: 4 });
    assert.dom(GENERAL.breadcrumbs).hasText(`${this.baseCrumbs} Roles`);
    this.assertTabState(assert, 'Roles');
  });

  test('it navigates to external role details and active orders', async function (assert) {
    const roleName = 'myrole';
    await visit(`${this.engineURL}/roles/${roleName}/details`);
    assert.strictEqual(
      currentURL(),
      `vault/secrets-engines/${this.mountPath}/pki/external/roles/${roleName}/details`,
      'it navigates to role details'
    );
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.external.roles.role.details');
    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(roleName);
    assert.dom(GENERAL.breadcrumb).exists({ count: 5 });
    assert.dom(GENERAL.breadcrumbs).hasText(`${this.baseCrumbs} Roles ${roleName}`);
    TABS.forEach((t) => assert.dom(GENERAL.linkTo(t)).doesNotExist());
    assert.dom(GENERAL.linkTo('Details')).exists().hasClass('active');
    assert.dom(GENERAL.linkTo('Active orders')).exists().doesNotHaveClass('active');

    // Navigate to a role's active orders
    await visit(`${this.engineURL}/roles/${roleName}/active-orders`);
    assert.strictEqual(
      currentURL(),
      `vault/secrets-engines/${this.mountPath}/pki/external/roles/${roleName}/active-orders`,
      'it navigates to role active-orders'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.active-orders'
    );
    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(roleName);
    assert.dom(GENERAL.breadcrumb).exists({ count: 6 });
    assert.dom(GENERAL.breadcrumbs).hasText(`${this.baseCrumbs} Roles ${roleName} Active orders`);
    TABS.forEach((t) => assert.dom(GENERAL.linkTo(t)).doesNotExist());
    assert.dom(GENERAL.linkTo('Details')).exists().doesNotHaveClass('active');
    assert.dom(GENERAL.linkTo('Active orders')).exists().hasClass('active');
  });

  test('it navigates to external orders', async function (assert) {
    await visit(`${this.engineURL}/orders`);
    assert.strictEqual(
      currentURL(),
      `vault/secrets-engines/${this.mountPath}/pki/external/orders`,
      'it navigates to orders'
    );
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.external.orders.index');
    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.mountPath);
    assert.dom(GENERAL.breadcrumb).exists({ count: 4 });
    assert.dom(GENERAL.breadcrumbs).hasText(`${this.baseCrumbs} Recent orders`);
    this.assertTabState(assert, 'Recent orders');
  });

  test('it navigates to external order details', async function (assert) {
    const orderID = '123';

    await visit(`${this.engineURL}/orders/${orderID}/details`);
    assert.strictEqual(
      currentURL(),
      `vault/secrets-engines/${this.mountPath}/pki/external/orders/${orderID}/details`,
      'it navigates to order details'
    );
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.external.orders.order.details');

    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(orderID);
    assert.dom(GENERAL.breadcrumb).exists({ count: 5 });
    assert.dom(GENERAL.breadcrumbs).hasText(`${this.baseCrumbs} Orders ${orderID}`);
    TABS.forEach((t) => {
      assert.dom(GENERAL.linkTo(t)).doesNotExist();
    });
  });

  test('it navigates to external DNS providers', async function (assert) {
    await visit(`${this.engineURL}/dns-providers`);
    assert.strictEqual(
      currentURL(),
      `vault/secrets-engines/${this.mountPath}/pki/external/dns-providers`,
      'it navigates to dns-providers'
    );
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.external.dns-providers');

    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.mountPath);
    assert.dom(GENERAL.breadcrumb).exists({ count: 4 });
    assert.dom(GENERAL.breadcrumbs).hasText(`${this.baseCrumbs} DNS providers`);
    this.assertTabState(assert, 'DNS providers');
  });

  test('it navigates to external ACME accounts', async function (assert) {
    await visit(`${this.engineURL}/acme-accounts`);
    assert.strictEqual(
      currentURL(),
      `vault/secrets-engines/${this.mountPath}/pki/external/acme-accounts`,
      'it navigates to acme-accounts'
    );
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.external.acme-accounts');

    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.mountPath);
    assert.dom(GENERAL.breadcrumb).exists({ count: 4 });
    assert.dom(GENERAL.breadcrumbs).hasText(`${this.baseCrumbs} ACME accounts`);
    this.assertTabState(assert, 'ACME accounts');
  });
});
