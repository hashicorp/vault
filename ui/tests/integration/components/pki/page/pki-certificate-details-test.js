/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import sinon from 'sinon';

module('Integration | Component | pki | Page::PkiCertificateDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    const downloadService = this.owner.lookup('service:download');
    this.downloadSpy = sinon.stub(downloadService, 'pem');

    const routerService = this.owner.lookup('service:router');
    this.routerSpy = sinon.stub(routerService, 'transitionTo');

    this.backend = 'pki';
    this.owner.lookup('service:secretMountPath').update(this.backend);

    this.revokeResponse = {
      revocation_time: 1673972804,
      revocation_time_rfc3339: '2023-01-17T16:26:44.960933411Z',
    };
    this.revokeStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'pkiRevoke')
      .resolves(this.revokeResponse);

    const sn = '4d:b6:ed:90:d6:b0:d4:bb:8e:5d:73:6a:6f:32:dc:8c:71:7c:db:5f';
    this.cert = {
      certificate: CERTIFICATES.rootPem,
      serial_number: sn,
      revocation_time: 0,
      revocation_time_rfc3339: '',
    };
    this.generatedCert = {
      certificate: CERTIFICATES.rootPem,
      ca_chain: ['-----BEGIN CERTIFICATE-----'],
      issuing_ca: '-----BEGIN CERTIFICATE-----',
      private_key: '-----BEGIN PRIVATE KEY-----',
      private_key_type: 'rsa',
      serial_number: sn,
    };
    this.certData = this.cert;
    this.canRevoke = true;
    this.onBack = sinon.spy();
    this.onRevoke = sinon.spy();

    this.renderComponent = () =>
      render(
        hbs`
          <Page::PkiCertificateDetails
            @certData={{this.certData}}
            @canRevoke={{this.canRevoke}}
            @onBack={{this.onBack}}
            @onRevoke={{this.onRevoke}}
          />
        `,
        { owner: this.engine }
      );
  });

  test('it should render actions and fields for base cert', async function (assert) {
    assert.expect(6);

    await this.renderComponent();
    assert
      .dom('[data-test-component="info-table-row"]')
      .exists({ count: 8 }, 'Correct number of fields render when certificate has not been revoked');
    assert
      .dom(`${GENERAL.infoRowValue('Certificate')} [data-test-certificate-card]`)
      .exists('Certificate card renders for certificate');
    assert.dom(`${GENERAL.infoRowValue('Serial number')} code`).exists('Serial number renders as monospace');

    await click('[data-test-pki-cert-download-button]');
    const { serial_number, certificate } = this.certData;
    assert.ok(
      this.downloadSpy.calledWith(serial_number.replace(/(\s|:)+/g, '-'), certificate),
      'Download pem method called with correct args'
    );

    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    assert.true(
      this.revokeStub.calledWith(this.backend, { serial_number: this.certData.serial_number }),
      'Revoke request called with correct params'
    );
    assert.dom(GENERAL.infoRowValue('Revocation time')).exists('Revocation time is displayed');
  });

  test('it should render actions and fields for generated cert', async function (assert) {
    assert.expect(10);

    this.certData = this.generatedCert;
    await this.renderComponent();

    assert.dom('[data-test-cert-detail-next-steps]').exists('Private key next steps warning shows');
    assert
      .dom('[data-test-component="info-table-row"]')
      .exists({ count: 12 }, 'Correct number of fields render when certificate has not been revoked');
    assert
      .dom(`${GENERAL.infoRowValue('Certificate')} [data-test-certificate-card]`)
      .exists('Certificate card renders for certificate');
    assert.dom(`${GENERAL.infoRowValue('Serial number')} code`).exists('Serial number renders as monospace');
    assert
      .dom(`${GENERAL.infoRowValue('CA chain')} [data-test-certificate-card]`)
      .exists('Certificate card renders for CA Chain');
    assert
      .dom(`${GENERAL.infoRowValue('Issuing CA')} [data-test-certificate-card]`)
      .exists('Certificate card renders for Issuing CA');
    assert
      .dom(`${GENERAL.infoRowValue('Private key')} [data-test-certificate-card]`)
      .exists('Certificate card renders for private key');

    await click('[data-test-pki-cert-download-button]');
    const { serial_number, certificate } = this.certData;
    assert.ok(
      this.downloadSpy.calledWith(serial_number.replace(/(\s|:)+/g, '-'), certificate),
      'Download pem method called with correct args'
    );

    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    assert.true(
      this.revokeStub.calledWith(this.backend, { serial_number: this.certData.serial_number }),
      'Revoke request called with correct params'
    );
    assert.dom(GENERAL.infoRowValue('Revocation time')).exists('Revocation time is displayed');
  });

  test('it should render back button', async function (assert) {
    assert.expect(1);

    await this.renderComponent();
    await click('[data-test-pki-cert-details-back]');
    assert.true(this.onBack.calledOnce, 'onBack action is called when back button is clicked');
  });

  test('it should send action on revoke if provided', async function (assert) {
    assert.expect(1);

    await this.renderComponent();
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.true(this.onRevoke.calledOnce, 'onRevoke action is called when certificate is revoked');
  });
});
