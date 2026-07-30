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

module('Acceptance | enterprise | pki | external | certificates | certificate route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const api = this.owner.lookup('service:api');
    this.certLookupStub = sinon.stub(api.secrets, 'pkiExternalCaReadLookupCert');
    this.fetchCertStub = sinon.stub(api.secrets, 'pkiExternalCaReadRoleOrderFetchCert');
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    this.serialNumber = 'ab:cd:ef:12:34:56';

    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));

    this.certificateURL = `/vault/secrets-engines/${this.mountPath}/pki/external/certificates/${this.serialNumber}`;
  });

  hooks.afterEach(async function () {
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it renders breadcrumbs and header without tabs for certificate lookup', async function (assert) {
    this.certLookupStub.resolves({
      order_id: 'order-abc-123',
      order_status: 'pending',
      role_name: 'myrole',
    });

    await visit(this.certificateURL);

    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText('View order');
    assert.dom(GENERAL.breadcrumb).exists({ count: 4 });
    assert.dom(GENERAL.breadcrumbs).hasText(`Vault Secrets engines ${this.mountPath} ${this.serialNumber}`);
    ['Overview', 'Roles', 'Recent orders', 'DNS providers', 'ACME accounts'].forEach((t) => {
      assert.dom(GENERAL.linkTo(t)).doesNotExist(`${t} tab does not render`);
    });
  });

  test('it requests a certificate lookup and displays response timestamp', async function (assert) {
    this.certLookupStub.resolves({
      order_id: 'order-abc-123',
      order_status: 'pending',
      role_name: 'myrole',
      identifiers: ['example.com'],
    });
    sinon.stub(timestamp, 'now').returns(new Date('2026-07-20T22:26:14.142Z'));

    await visit(this.certificateURL);

    assert.strictEqual(currentURL(), this.certificateURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.certificates.certificate',
      'navigated to certificate route'
    );
    assert.true(this.certLookupStub.calledOnce, 'cert lookup called once');
    assert.dom('h1').hasText('View order', 'page title is displayed');
    assert.dom(GENERAL.textBody('Last refreshed')).hasTextContaining('Last refreshed July 20, 2026');
  });

  test('it does not fetch certificate when order status is not "completed"', async function (assert) {
    this.certLookupStub.resolves({
      order_id: 'order-abc-123',
      order_status: 'pending',
      role_name: 'myrole',
    });

    await visit(this.certificateURL);

    assert.true(this.certLookupStub.calledOnce, 'cert lookup called once');
    assert.true(this.fetchCertStub.notCalled, 'certificate fetch not attempted for non-completed order');
  });

  test('it fetches certificate when order status is "completed"', async function (assert) {
    this.certLookupStub.resolves({
      order_id: 'order-abc-123',
      order_status: 'completed',
      role_name: 'myrole',
    });
    this.fetchCertStub.resolves({
      serial_number: this.serialNumber,
      certificate: '-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----',
    });

    await visit(this.certificateURL);

    assert.true(this.certLookupStub.calledOnce, 'cert lookup called once');
    assert.true(this.fetchCertStub.calledOnce, 'certificate fetch called once');
    const [role, orderId, mount] = this.fetchCertStub.lastCall.args;
    assert.strictEqual(role, 'myrole', 'cert request called with expected role');
    assert.strictEqual(orderId, 'order-abc-123', 'cert request called with expected order ID');
    assert.strictEqual(mount, this.mountPath, 'cert request called with expected mount path');
  });

  test('it handles certificate fetch 403 error', async function (assert) {
    this.certLookupStub.resolves({
      order_id: 'order-abc-123',
      order_status: 'completed',
      role_name: 'myrole',
    });
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.fetchCertStub.rejects(getErrorResponse(error, 403));

    await visit(this.certificateURL);

    assert.true(this.certLookupStub.calledOnce, 'cert lookup called once');
    assert.true(this.fetchCertStub.calledOnce, 'certificate fetch attempted');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.certificates.certificate',
      'stays on certificate route despite cert fetch failure'
    );
    assert.dom(GENERAL.infoRowValue('Order status')).hasText('Completed', 'order details still displayed');
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText(
        'Certificate data is unavailable You do not have "read" permissions for the path: /v1/test/error/parsing'
      );
  });

  test('it handles certificate fetch 400 error', async function (assert) {
    this.certLookupStub.resolves({
      order_id: 'order-abc-123',
      order_status: 'completed',
      role_name: 'myrole',
    });
    const error = { errors: ['order has status expired, must be completed to fetch cert'] };
    this.fetchCertStub.rejects(getErrorResponse(error, 400));

    await visit(this.certificateURL);

    assert.true(this.certLookupStub.calledOnce, 'cert lookup called once');
    assert.true(this.fetchCertStub.calledOnce, 'certificate fetch attempted');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.certificates.certificate',
      'stays on certificate route despite cert fetch failure'
    );
    assert.dom(GENERAL.infoRowValue('Order status')).hasText('Completed', 'order details still displayed');
    assert
      .dom(GENERAL.messageError)
      .exists()
      .hasText('Certificate data is unavailable order has status expired, must be completed to fetch cert');
  });

  test('it handles cert lookup 404 error', async function (assert) {
    this.certLookupStub.rejects(getErrorResponse());

    await visit(this.certificateURL);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'redirects to error route on 404'
    );
    assert.true(this.certLookupStub.calledOnce, 'cert lookup called once');
    assert.true(this.fetchCertStub.notCalled, 'certificate fetch not attempted after lookup error');
    assert.dom(GENERAL.pageError.title(404)).exists().hasText('ERROR 404 Not found');
  });

  test('it handles cert lookup 403 permission denied error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.certLookupStub.rejects(getErrorResponse(error, 403));

    await visit(this.certificateURL);

    assert.strictEqual(currentURL(), this.certificateURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'redirects to error route on 403'
    );
    assert.true(this.certLookupStub.calledOnce, 'cert lookup called once');
    assert.true(this.fetchCertStub.notCalled, 'certificate fetch not attempted after lookup error');
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');
  });
});
