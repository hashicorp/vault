/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { currentRouteName, currentURL, visit } from '@ember/test-helpers';
import sinon from 'sinon';
import timestamp from 'core/utils/timestamp';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Acceptance | enterprise | pki | external | orders | order route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.orderReadStub = sinon.stub(api.secrets, 'pkiExternalCaReadLookupOrder');
    this.fetchCertStub = sinon.stub(api.secrets, 'pkiExternalCaReadRoleOrderFetchCert');
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    this.orderId = 'test-order-123';

    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));

    // assertion helpers
    this.orderURL = `/vault/secrets-engines/${this.mountPath}/pki/external/orders/${this.orderId}`;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it renders breadcrumbs for order from orders list', async function (assert) {
    this.orderReadStub.resolves({
      order_id: this.orderId,
      status: 'valid',
      challenges: {
        'example.com': [
          {
            challenge_status: 'valid',
            challenge_type: 'dns-01',
            expires: '2026-07-24T21:34:36Z',
            requires_manual_fulfillment: 'false',
          },
        ],
      },
    });

    await visit(this.orderURL);

    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText('View order');
    assert.dom(GENERAL.breadcrumb).exists({ count: 5 });
    assert
      .dom(GENERAL.breadcrumbs)
      .hasText(`Vault Secrets engines ${this.mountPath} Recent orders ${this.orderId}`);
  });

  test('it requests an order and displays response timestamp', async function (assert) {
    this.orderReadStub.resolves({
      order_id: this.orderId,
      order_status: 'pending',
      identifiers: ['example.com', 'test.example.com'],
      challenges: {
        'example.com': [
          {
            challenge_status: 'valid',
            challenge_type: 'dns-01',
            expires: '2026-07-24T21:34:36Z',
            requires_manual_fulfillment: 'false',
          },
        ],
        'test.example.com': [
          {
            challenge_status: 'pending',
            challenge_type: 'http-01',
            expires: '2026-07-25T21:34:36Z',
            requires_manual_fulfillment: 'true',
          },
        ],
      },
    });
    sinon.stub(timestamp, 'now').returns(new Date('2026-07-20T22:26:14.142Z'));
    await visit(this.orderURL);

    assert.strictEqual(currentURL(), this.orderURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.orders.order',
      'navigated to order route'
    );
    assert.true(this.orderReadStub.calledOnce, 'order read called once');
    assert.dom('h1').hasText('View order', 'page title is displayed');
    assert.dom(GENERAL.textBody('Last refreshed')).hasTextContaining('Last refreshed: July 20, 2026');
  });

  test('it does not fetch certificate when order status is not "completed"', async function (assert) {
    this.orderReadStub.resolves({
      order_id: this.orderId,
      order_status: 'pending',
      role_name: 'myrole',
    });

    await visit(this.orderURL);

    assert.true(this.orderReadStub.calledOnce, 'order read called once');
    assert.true(this.fetchCertStub.notCalled, 'certificate read not attempted for non-completed order');
  });

  test('it fetches certificate when order status is "completed"', async function (assert) {
    const serialNumber = 'ab:cd:ef:12:34:56';
    this.orderReadStub.resolves({
      order_id: this.orderId,
      order_status: 'completed',
      role_name: 'myrole',
    });
    this.fetchCertStub.resolves({
      serial_number: serialNumber,
      certificate: '-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----',
    });

    await visit(this.orderURL);

    assert.true(this.orderReadStub.calledOnce, 'order read called once');
    assert.true(this.fetchCertStub.calledOnce, 'certificate read called once');
    const [role, orderId, mount] = this.fetchCertStub.lastCall.args;
    assert.strictEqual(role, 'myrole', 'cert request called with expected role');
    assert.strictEqual(orderId, this.orderId, 'cert request called with expected order ID');
    assert.strictEqual(mount, this.mountPath, 'cert request called with expected mount path');
  });

  test('it handles certificate 403 error', async function (assert) {
    this.orderReadStub.resolves({
      order_id: this.orderId,
      order_status: 'completed',
      role_name: 'myrole',
    });
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.fetchCertStub.rejects(getErrorResponse(error, 403));

    await visit(this.orderURL);
    assert.true(this.orderReadStub.calledOnce, 'order read called once');
    assert.true(this.fetchCertStub.calledOnce, 'certificate read attempted');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.orders.order',
      'stays on order route despite cert fetch failure'
    );
    assert.dom(GENERAL.infoRowValue('Order status')).hasText('Completed', 'order details still displayed');
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText(
        'Certificate data is unavailable You do not have "read" permissions for the path: /v1/test/error/parsing'
      );
  });

  test('it handles certificate 400 error', async function (assert) {
    this.orderReadStub.resolves({
      order_id: this.orderId,
      order_status: 'completed',
      role_name: 'myrole',
    });
    const error = { errors: ['order has status expired, must be completed to fetch cert'] };
    this.fetchCertStub.rejects(getErrorResponse(error, 400));

    await visit(this.orderURL);
    assert.true(this.orderReadStub.calledOnce, 'order read called once');
    assert.true(this.fetchCertStub.calledOnce, 'certificate read attempted');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.orders.order',
      'stays on order route despite cert fetch failure'
    );
    assert.dom(GENERAL.infoRowValue('Order status')).hasText('Completed', 'order details still displayed');
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText('Certificate data is unavailable order has status expired, must be completed to fetch cert');
  });

  test('it handles order 404 error', async function (assert) {
    this.orderReadStub.rejects(getErrorResponse());

    await visit(this.orderURL);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'redirects to error route on 404'
    );
    assert.true(this.orderReadStub.calledOnce, 'order read called once');
    assert.true(this.fetchCertStub.notCalled, 'certificate read not attempted after order error');
    assert.dom(GENERAL.pageError.title(404)).exists().hasText('ERROR 404 Not found');
  });

  test('it handles order 403 permission denied error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.orderReadStub.rejects(getErrorResponse(error, 403));

    await visit(this.orderURL);

    assert.strictEqual(currentURL(), this.orderURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'redirects to error route on 403'
    );
    assert.true(this.orderReadStub.calledOnce, 'order read called once');
    assert.true(this.fetchCertStub.notCalled, 'certificate read not attempted after order error');
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');
  });
});
