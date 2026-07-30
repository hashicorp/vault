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
import { assertTabState } from 'vault/tests/helpers/pki/assertion-helpers';

const ROLE_TABS = ['Overview', 'Active orders'];
module('Acceptance | enterprise | pki | external | roles | role | overview route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.roleReadStub = sinon.stub(api.secrets, 'pkiExternalCaReadRole');
    this.activeOrdersListStub = sinon.stub(api.secrets, 'pkiExternalCaListRoleActiveOrders');
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    this.roleName = 'test-role';

    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.roleOverviewURL = `/vault/secrets-engines/${this.mountPath}/pki/external/roles/${this.roleName}/overview`;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it navigates to role overview route', async function (assert) {
    this.roleReadStub.resolves({ name: this.roleName });
    this.activeOrdersListStub.withArgs(this.roleName, this.mountPath).rejects(getErrorResponse());
    await visit(this.roleOverviewURL);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/roles/${this.roleName}/overview`,
      'it navigates to url'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.overview',
      'it navigates to route'
    );
    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.roleName);
    assert.dom(GENERAL.breadcrumb).exists({ count: 5 });
    assert.dom(GENERAL.breadcrumbs).hasText(`Vault Secrets engines ${this.mountPath} Roles ${this.roleName}`);
    assertTabState(assert, 'Overview', ROLE_TABS);
    // Navigate to a active orders
    await click(GENERAL.linkTo('Active orders'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.active-orders'
    );
    assertTabState(assert, 'Active orders', ROLE_TABS);
  });

  test('it fetches and displays role details', async function (assert) {
    this.roleReadStub.resolves({
      name: this.roleName,
      acme_account_name: 'production-account',
      dns_provider_name: 'aws-route53-prod',
      allowed_domains: ['example.com', '*.example.com'],
      allow_subdomains: true,
    });
    await visit(this.roleOverviewURL);
    assert.strictEqual(currentURL(), this.roleOverviewURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.overview',
      'navigated to role overview route'
    );
    assert.true(this.roleReadStub.calledOnce, 'role read called once');
    assert.true(this.activeOrdersListStub.calledOnce, 'active orders read called once');
    assert.dom('h1').hasText(this.roleName, 'role name is header');
    assert.dom(GENERAL.infoRowValue('ACME account name')).hasText('production-account');
    assert.dom(GENERAL.infoRowLabel('Name')).doesNotExist();
    assert
      .dom('[data-test-component="info-table-row"]')
      .exists({ count: 4 }, 'it renders every config param EXCEPT "name"');
  });

  test('it handles role read 404 error', async function (assert) {
    this.roleReadStub.withArgs(this.roleName, this.mountPath).rejects(getErrorResponse());
    await visit(this.roleOverviewURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'redirects to role error route'
    );
    assert.true(this.roleReadStub.calledOnce, 'role read called once');
    assert.true(this.activeOrdersListStub.notCalled, 'active-orders is not called');
    assert.dom('h1').hasText(this.mountPath, 'mount path is header');
    assert.dom(GENERAL.pageError.title(404)).exists().hasText('ERROR 404 Not found');
  });

  test('it handles role read 403 permission denied error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.roleReadStub.withArgs(this.roleName, this.mountPath).rejects(getErrorResponse(error, 403));
    await visit(this.roleOverviewURL);
    assert.strictEqual(currentURL(), this.roleOverviewURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'redirects to role.error route'
    );
    assert.true(this.roleReadStub.calledOnce, 'role read called once');
    assert.true(this.activeOrdersListStub.notCalled, 'active-orders is not called');
    assert.dom('h1').hasText(this.mountPath, 'mount path is header');
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');

    await click(GENERAL.linkTo('Back'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.index',
      'Back link navigates to roles index route'
    );
  });

  test('it catches active-order 404 error', async function (assert) {
    this.roleReadStub.resolves({ name: this.roleName });
    this.activeOrdersListStub.withArgs(this.roleName, this.mountPath).rejects(getErrorResponse());
    await visit(this.roleOverviewURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.overview',
      'navigated to role overview route'
    );
    assert.true(this.roleReadStub.calledOnce, 'role read called once');
    assert.true(this.activeOrdersListStub.calledOnce, 'active-orders read called once');
    assert.dom(GENERAL.overviewCard.container('Active orders')).exists();
  });

  test('it catches active-order 403 error', async function (assert) {
    this.roleReadStub.resolves({ name: this.roleName });
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.activeOrdersListStub.withArgs(this.roleName, this.mountPath).rejects(getErrorResponse(error, 403));
    await visit(this.roleOverviewURL);
    assert.strictEqual(currentURL(), this.roleOverviewURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.overview',
      'navigated to role overview route'
    );
    assert.true(this.roleReadStub.calledOnce, 'role read called once');
    assert.true(this.activeOrdersListStub.calledOnce, 'active-orders read called once');
    assert.dom(GENERAL.overviewCard.container('Active orders')).doesNotExist();
  });
});
