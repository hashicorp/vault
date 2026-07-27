/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | pki | external-pki | ExternalPki::Page::Certificate', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');

    this.model = {
      engine: { id: 'pki-external-ca' },
      serial_number: 'ab:cd:ef:12:34:56',
      certLookup: {
        order_id: 'order-abc-123',
        order_status: 'pending',
        role_name: 'myrole',
        identifiers: ['example.com'],
        not_before: undefined,
        not_after: undefined,
      },
      certificate: undefined,
      responseTimestamp: new Date('2026-07-14T21:00:00Z'),
    };

    this.renderComponent = () =>
      render(
        hbs`<ExternalPki::Page::Certificate @model={{this.model}} @breadcrumbs={{array (hash label="View order")}} />`,
        {
          owner: this.engine,
        }
      );
  });

  test('it renders the last refreshed timestamp', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.textBody('Last refreshed')).hasTextContaining('Last refreshed: July 14, 2026');
  });

  test('it calls router.refresh with the certificate route when Refresh button is clicked', async function (assert) {
    const refreshStub = sinon.stub(this.router, 'refresh');
    await this.renderComponent();
    await click(GENERAL.button('Refresh'));
    assert.true(refreshStub.calledOnce, 'refresh was called once');
    const [route] = refreshStub.lastCall.args;
    assert.strictEqual(
      route,
      'vault.cluster.secrets.backend.pki.external.certificates.certificate',
      'refresh was called with the certificate route'
    );
  });

  test('orderParams: it renders order details from the certLookup request', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.infoRowValue('Order status')).hasText('Pending');
    assert.dom(GENERAL.infoRowValue('Role name')).hasText('myrole');
  });

  test('certParams: undefined certificate does not render certificate card or error banner', async function (assert) {
    this.model.certificate = undefined;
    await this.renderComponent();
    assert.dom(GENERAL.cardContainer('Certificate details')).doesNotExist();
    assert.dom(GENERAL.messageError).doesNotExist();
  });

  test('certParams: certificate with details renders certificate card', async function (assert) {
    this.model.certificate = {
      details: {
        serial_number: 'ab:cd:ef:12:34:56',
        certificate: '-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----',
      },
      error: undefined,
    };
    await this.renderComponent();
    assert
      .dom(GENERAL.cardContainer('Certificate details'))
      .exists('cert card rendered when details are present');
    assert.dom(GENERAL.infoRowValue('Serial number')).hasText('ab:cd:ef:12:34:56');
    assert.dom(GENERAL.messageError).doesNotExist('no error banner when certificate fetch succeeds');
  });

  test('certParams: certificate error with no details falls back to certLookup validity dates', async function (assert) {
    // fetch-cert failed but the serial number lookup returned validity dates.
    // certParams injects not_before/not_after from certLookup as the details so the
    // cert card still renders what little data is available alongside the error banner.
    this.model.certLookup.not_before = '2026-01-01T00:00:00Z';
    this.model.certLookup.not_after = '2027-01-01T00:00:00Z';
    this.model.certificate = {
      details: undefined,
      error: {
        status: 403,
        path: '/v1/pki/role/myrole/order/order-abc-123/fetch-cert',
        message: 'permission denied',
      },
    };
    await this.renderComponent();
    assert.dom(GENERAL.messageError).exists('error banner is shown');
    assert
      .dom(GENERAL.cardContainer('Certificate details'))
      .exists('cert card rendered with fallback not_before/not_after from certLookup');
  });

  test('certParams: certificate error with existing details is passed through without overwriting', async function (assert) {
    // If certificate.details is already populated, certParams returns it unchanged —
    // the not_before/not_after fallback is only applied when details is falsy.
    this.model.certLookup.not_before = '2026-01-01T00:00:00Z';
    this.model.certLookup.not_after = '2027-01-01T00:00:00Z';
    this.model.certificate = {
      details: { serial_number: 'ab:cd:ef:12:34:56' },
      error: { status: 400, message: 'some error' },
    };
    await this.renderComponent();
    assert
      .dom(GENERAL.cardContainer('Certificate details'))
      .exists('cert card rendered from existing details');
    assert.dom(GENERAL.infoRowValue('Serial number')).hasText('ab:cd:ef:12:34:56');
    assert.dom(GENERAL.messageError).exists('error banner rendered alongside cert details');
  });
});
