/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { currentRouteName, currentURL, visit } from '@ember/test-helpers';
import sinon from 'sinon';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

const TABS = ['Overview', 'Roles', 'Recent orders', 'DNS providers', 'ACME accounts'];

// This test asserts tab state navigation for each route in pki.external
module('Acceptance | enterprise | pki | external | navigation', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.acmeListStub = sinon.stub(api.secrets, 'pkiExternalCaListConfigAcmeAccount');
    this.dnsListStub = sinon.stub(api.secrets, 'pkiExternalCaListConfigDns');
    this.rolesListStub = sinon.stub(api.secrets, 'pkiExternalCaListRole');
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.engineURL = `vault/secrets-engines/${this.mountPath}/pki/external`;
    this.assertTabState = (assert, activeTab) => {
      const inactive = TABS.filter((t) => t !== activeTab);
      inactive.forEach((t) => {
        assert.dom(GENERAL.linkTo(t)).exists().doesNotHaveClass('active', `${t} is inactive`);
      });
      assert.dom(GENERAL.linkTo(activeTab)).exists().hasClass('active', `${activeTab} is active`);
    };
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('only "Overview" tab renders when no resources exist but user has permission to list everything', async function (assert) {
    // getErrorResponse() throws 404 by default
    this.acmeListStub.rejects(getErrorResponse());
    this.dnsListStub.rejects(getErrorResponse());
    this.rolesListStub.rejects(getErrorResponse());
    await visit(this.engineURL); // navigate to index route to test it redirects to overview
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/overview`,
      'it navigates to overview'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.overview',
      'it redirects to overview route'
    );
    assert.dom(GENERAL.linkTo('Overview')).exists().hasClass('active');
    const hidden = TABS.filter((t) => t !== 'Overview');
    hidden.forEach((t) => assert.dom(GENERAL.linkTo(t)).doesNotExist());
  });

  test('All tabs render when user does NOT have permission to list ACME accounts', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.acmeListStub.rejects(getErrorResponse(error, 403));
    await visit(this.engineURL);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/overview`,
      'it navigates to overview'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.overview',
      'it redirects to overview route'
    );
    this.assertTabState(assert, 'Overview');
  });

  test('All tabs render when user does NOT have permission to list DNS providers', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.dnsListStub.rejects(getErrorResponse(error, 403));
    await visit(this.engineURL);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/overview`,
      'it navigates to overview'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.overview',
      'it redirects to overview route'
    );
    this.assertTabState(assert, 'Overview');
  });

  test('All tabs render when user does NOT have permission to list roles', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.rolesListStub.rejects(getErrorResponse(error, 403));
    await visit(this.engineURL);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/overview`,
      'it navigates to overview'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.overview',
      'it redirects to overview route'
    );
    this.assertTabState(assert, 'Overview');
  });

  module('configured', function (hooks) {
    hooks.beforeEach(async function () {
      this.acmeListStub.resolves({ keys: ['my-acme-account'] });
      this.rolesListStub.rejects(getErrorResponse());
      this.dnsListStub.rejects(getErrorResponse());
      // assertion helpers
      this.baseCrumbs = `Vault Secrets engines ${this.mountPath}`;
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
      assert.dom(GENERAL.breadcrumb).exists({ count: 5 });
      assert.dom(GENERAL.breadcrumbs).hasText(`${this.baseCrumbs} Roles ${roleName}`);
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

    test('it navigates to external order', async function (assert) {
      const orderID = '123';

      await visit(`${this.engineURL}/orders/${orderID}`);
      assert.strictEqual(
        currentURL(),
        `vault/secrets-engines/${this.mountPath}/pki/external/orders/${orderID}`,
        'it navigates to order'
      );
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.external.orders.order');

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
});
