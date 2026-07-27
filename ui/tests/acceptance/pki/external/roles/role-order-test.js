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

module('Acceptance | enterprise | pki | external | roles | role | order route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const api = this.owner.lookup('service:api');
    this.orderStatusStub = sinon.stub(api.secrets, 'pkiExternalCaReadRoleOrderStatus');
    this.fetchCertStub = sinon.stub(api.secrets, 'pkiExternalCaReadRoleOrderFetchCert');
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    this.roleName = 'test-role';
    this.orderId = 'test-order-123';

    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));

    this.orderURL = `/vault/secrets-engines/${this.mountPath}/pki/external/roles/${this.roleName}/${this.orderId}`;
  });

  hooks.afterEach(async function () {
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it renders correct breadcrumbs with role name as a link and "View order" as the leaf', async function (assert) {
    this.orderStatusStub.resolves({ order_status: 'completed' });

    await visit(this.orderURL);

    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.external.roles.role.order');
    // Title should be "View order", not the role name
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('View order');
    assert.dom(GENERAL.breadcrumb).exists({ count: 6 });
    assert
      .dom(GENERAL.breadcrumbs)
      .hasText(`Vault Secrets engines ${this.mountPath} Roles ${this.roleName} View order`);
  });

  test('role name breadcrumb is a link back to the role', async function (assert) {
    this.orderStatusStub.resolves({ order_status: 'completed' });
    await visit(this.orderURL);
    assert.dom(GENERAL.breadcrumbLink(this.roleName)).exists();
  });

  test('tabs are hidden on the order route', async function (assert) {
    this.orderStatusStub.resolves({ order_status: 'completed' });
    await visit(this.orderURL);
    assert.dom(GENERAL.linkTo('Details')).doesNotExist('Details tab is hidden');
    assert.dom(GENERAL.linkTo('Active orders')).doesNotExist('Active orders tab is hidden');
  });

  test('it requests order details', async function (assert) {
    this.orderStatusStub.resolves({ order_status: 'completed' });
    await visit(this.orderURL);
    assert.strictEqual(currentURL(), this.orderURL, 'it has expected URL');
    assert.true(this.orderStatusStub.calledOnce, 'order status called once');
    assert.true(this.fetchCertStub.calledOnce, 'cert fetch is called when order status is complete');
    assert.dom(GENERAL.infoRowValue('Order status')).hasText('Completed');
  });

  test('it throws order status 404', async function (assert) {
    this.orderStatusStub.rejects(getErrorResponse());
    await visit(this.orderURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.error',
      'redirects to role error route on 404'
    );
    assert.true(this.orderStatusStub.calledOnce, 'order status called once');
    assert.true(this.fetchCertStub.notCalled, 'cert fetch is NOT called when order 404s');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('View order', 'parent header is still rendered');
    assert.dom(GENERAL.breadcrumb).exists({ count: 6 }, 'parent breadcrumbs are still rendered');
  });

  test('it catches order status 403 error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.orderStatusStub.rejects(getErrorResponse(error, 403));
    await visit(this.orderURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.order',
      'renders 403 on order route'
    );
    assert.true(this.orderStatusStub.calledOnce, 'order status called once');
    assert.true(this.fetchCertStub.calledOnce, 'cert fetch is called when order 403s');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('View order', 'parent header is still rendered');
    assert.dom(GENERAL.breadcrumb).exists({ count: 6 }, 'parent breadcrumbs are still rendered');
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText(
        'Order status is unavailable You do not have "read" permissions for the path: /v1/test/error/parsing'
      );
  });

  test('it catches cert fetch 404 error', async function (assert) {
    this.orderStatusStub.resolves({ order_status: 'completed' });
    this.fetchCertStub.rejects(getErrorResponse());
    await visit(this.orderURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.order',
      'renders 404 on order route'
    );
    assert.true(this.orderStatusStub.calledOnce, 'order status called once');
    assert.true(this.fetchCertStub.calledOnce, 'cert fetch is called');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('View order', 'parent header is still rendered');
    assert.dom(GENERAL.breadcrumb).exists({ count: 6 }, 'parent breadcrumbs are still rendered');
    assert.dom(GENERAL.messageError).exists().hasText('Certificate data is unavailable');
  });

  test('it catches cert fetch 403 error', async function (assert) {
    this.orderStatusStub.resolves({ order_status: 'completed' });
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.fetchCertStub.rejects(getErrorResponse(error, 403));
    await visit(this.orderURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.order',
      'renders 403 on order route'
    );
    assert.true(this.orderStatusStub.calledOnce, 'order status called once');
    assert.true(this.fetchCertStub.calledOnce, 'cert fetch is called');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('View order', 'parent header is still rendered');
    assert.dom(GENERAL.breadcrumb).exists({ count: 6 }, 'parent breadcrumbs are still rendered');
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText(
        'Certificate data is unavailable You do not have "read" permissions for the path: /v1/test/error/parsing'
      );
  });

  test('it throws if both endpoints error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.orderStatusStub.rejects(getErrorResponse(error, 403));
    this.fetchCertStub.rejects(getErrorResponse(error, 403));
    await visit(this.orderURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.role.error',
      'redirects to role error route when both requests 403'
    );
    assert.true(this.orderStatusStub.calledOnce, 'order status called once');
    assert.true(this.fetchCertStub.calledOnce, 'cert fetch called once');
    assert.true(this.orderStatusStub.calledOnce, 'order status called once');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('View order', 'parent header is still rendered');
    assert.dom(GENERAL.breadcrumb).exists({ count: 6 }, 'parent breadcrumbs are still rendered');
  });
});
