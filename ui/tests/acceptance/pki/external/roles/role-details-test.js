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

module('Acceptance | enterprise | pki | external | roles | role | details route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.roleReadStub = sinon.stub(api.secrets, 'pkiExternalCaReadRole');
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    this.roleName = 'test-role';

    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.roleDetailsURL = `/vault/secrets-engines/${this.mountPath}/pki/external/roles/${this.roleName}/details`;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it renders breadcrumbs for role details', async function (assert) {
    await visit(this.roleDetailsURL);
    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.roleName);
    assert.dom(GENERAL.breadcrumb).exists({ count: 5 });
    assert.dom(GENERAL.breadcrumbs).hasText(`Vault Secrets engines ${this.mountPath} Roles ${this.roleName}`);
    assert.dom(GENERAL.linkTo('Details')).exists().hasClass('active');
    assert.dom(GENERAL.linkTo('Active orders')).exists().doesNotHaveClass('active');
    // Navigate to a active orders
    await click(GENERAL.linkTo('Active orders'));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.active-orders'
    );
    assert.dom(GENERAL.linkTo('Details')).exists().doesNotHaveClass('active');
    assert.dom(GENERAL.linkTo('Active orders')).exists().hasClass('active');
  });

  test('it fetches and displays role details', async function (assert) {
    this.roleReadStub.resolves({
      name: this.roleName,
      acme_account_name: 'production-account',
      dns_provider_name: 'aws-route53-prod',
      allowed_domains: ['example.com', '*.example.com'],
      allow_subdomains: true,
    });

    await visit(this.roleDetailsURL);
    assert.strictEqual(currentURL(), this.roleDetailsURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.details',
      'navigated to role details route'
    );
    assert.true(this.roleReadStub.calledOnce, 'role read called once');
    assert.dom('h1').hasText(this.roleName, 'role name is displayed');
    assert.dom(GENERAL.infoRowValue('ACME account name')).hasText('production-account');
    assert.dom(GENERAL.infoRowLabel('Name')).doesNotExist();
    assert
      .dom('[data-test-component="info-table-row"]')
      .exists({ count: 4 }, 'it renders every config param EXCEPT "name"');
  });

  test('it handles 404 error', async function (assert) {
    this.roleReadStub.withArgs(this.roleName, this.mountPath).rejects(getErrorResponse());
    await visit(this.roleDetailsURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.error',
      'redirects to role error route'
    );
    assert.true(this.roleReadStub.calledOnce, 'role read called once');
    assert.dom('h1').hasText(this.roleName, 'role name is displayed');
    assert.dom(GENERAL.pageError.title(404)).exists().hasText('ERROR 404 Not found');
  });

  test('it handles 403 permission denied error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.roleReadStub.withArgs(this.roleName, this.mountPath).rejects(getErrorResponse(error, 403));
    await visit(this.roleDetailsURL);
    assert.strictEqual(currentURL(), this.roleDetailsURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.error',
      'redirects to role.error route'
    );
    assert.true(this.roleReadStub.calledOnce, 'role read called once');
    assert.dom('h1').hasText(this.roleName, 'role name is displayed');
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');
  });
});
