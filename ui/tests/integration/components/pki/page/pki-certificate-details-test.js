/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';

module('Integration | Component | pki | Page::PkiCertificateDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const downloadService = this.owner.lookup('service:download');
    this.downloadSpy = sinon.stub(downloadService, 'pem');

    const routerService = this.owner.lookup('service:router');
    this.routerSpy = sinon.stub(routerService, 'transitionTo');

    this.owner.lookup('service:secretMountPath').update('pki');

    const store = this.owner.lookup('service:store');
    const id = '4d:b6:ed:90:d6:b0:d4:bb:8e:5d:73:6a:6f:32:dc:8c:71:7c:db:5f';
    store.pushPayload('pki/certificate/base', {
      modelName: 'pki/certificate/base',
      data: {
        certificate: '-----BEGIN CERTIFICATE-----',
        common_name: 'example.com Intermediate Authority',
        issue_date: 1673540867000,
        serial_number: id,
        parsed_certificate: {
          not_valid_after: 1831220897000,
          not_valid_before: 1673540867000,
        },
      },
    });
    store.pushPayload('pki/certificate/generate', {
      modelName: 'pki/certificate/generate',
      data: {
        certificate: '-----BEGIN CERTIFICATE-----',
        ca_chain: '-----BEGIN CERTIFICATE-----',
        issuer_ca: '-----BEGIN CERTIFICATE-----',
        private_key: '-----BEGIN PRIVATE KEY-----',
        private_key_type: 'rsa',
        common_name: 'example.com Intermediate Authority',
        issue_date: 1673540867000,
        serial_number: id,
        parsed_certificate: {
          not_valid_after: 1831220897000,
          not_valid_before: 1673540867000,
        },
      },
    });
    this.model = store.peekRecord('pki/certificate/base', id);
    this.generatedModel = store.peekRecord('pki/certificate/generate', id);

    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
        'pki/revoke': ['root'],
      },
    }));
  });

  test('it should render actions and fields for base cert', async function (assert) {
    assert.expect(6);

    this.server.post('/pki/revoke', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      assert.strictEqual(
        data.serial_number,
        this.model.serialNumber,
        'Revoke request made with serial number'
      );
      return {
        data: {
          revocation_time: 1673972804,
          revocation_time_rfc3339: '2023-01-17T16:26:44.960933411Z',
        },
      };
    });

    await render(hbs`<Page::PkiCertificateDetails @model={{this.model}} />`, { owner: this.engine });
    assert
      .dom('[data-test-component="info-table-row"]')
      .exists({ count: 5 }, 'Correct number of fields render when certificate has not been revoked');
    assert
      .dom('[data-test-value-div="Certificate"] [data-test-certificate-card]')
      .exists('Certificate card renders for certificate');
    assert.dom('[data-test-value-div="Serial number"] code').exists('Serial number renders as monospace');

    await click('[data-test-pki-cert-download-button]');
    const { serialNumber, certificate } = this.model;
    assert.ok(
      this.downloadSpy.calledWith(serialNumber.replace(/(\s|:)+/g, '-'), certificate),
      'Download pem method called with correct args'
    );

    await click('[data-test-confirm-action-trigger]');
    await click('[data-test-confirm-button]');

    assert.dom('[data-test-value-div="Revocation time"]').exists('Revocation time is displayed');
  });

  test('it should render actions and fields for generated cert', async function (assert) {
    assert.expect(10);

    this.server.post('/pki/revoke', (schema, req) => {
      const data = JSON.parse(req.requestBody);
      assert.strictEqual(
        data.serial_number,
        this.model.serialNumber,
        'Revoke request made with serial number'
      );
      return {
        data: {
          revocation_time: 1673972804,
          revocation_time_rfc3339: '2023-01-17T16:26:44.960933411Z',
        },
      };
    });

    await render(hbs`<Page::PkiCertificateDetails @model={{this.generatedModel}} />`, { owner: this.engine });
    assert.dom('[data-test-cert-detail-next-steps]').exists('Private key next steps warning shows');
    assert
      .dom('[data-test-component="info-table-row"]')
      .exists({ count: 9 }, 'Correct number of fields render when certificate has not been revoked');
    assert
      .dom('[data-test-value-div="Certificate"] [data-test-certificate-card]')
      .exists('Certificate card renders for certificate');
    assert.dom('[data-test-value-div="Serial number"] code').exists('Serial number renders as monospace');
    assert
      .dom('[data-test-value-div="CA Chain"] [data-test-certificate-card]')
      .exists('Certificate card renders for CA Chain');
    assert
      .dom('[data-test-value-div="Issuing CA"] [data-test-certificate-card]')
      .exists('Certificate card renders for Issuing CA');
    assert
      .dom('[data-test-value-div="Private key"] [data-test-certificate-card]')
      .exists('Certificate card renders for private key');

    await click('[data-test-pki-cert-download-button]');
    const { serialNumber, certificate } = this.model;
    assert.ok(
      this.downloadSpy.calledWith(serialNumber.replace(/(\s|:)+/g, '-'), certificate),
      'Download pem method called with correct args'
    );

    await click('[data-test-confirm-action-trigger]');
    await click('[data-test-confirm-button]');

    assert.dom('[data-test-value-div="Revocation time"]').exists('Revocation time is displayed');
  });

  test('it should render back button', async function (assert) {
    assert.expect(1);

    this.cancel = () => assert.ok('onBack action is triggered');

    await render(hbs`<Page::PkiCertificateDetails @model={{this.model}} @onBack={{this.cancel}} />`, {
      owner: this.engine,
    });
    await click('[data-test-pki-cert-details-back]');
  });

  test('it should send action on revoke if provided', async function (assert) {
    assert.expect(1);

    this.server.post('/pki/revoke', () => ({
      data: {
        revocation_time: 1673972804,
        revocation_time_rfc3339: '2023-01-17T16:26:44.960933411Z',
      },
    }));

    this.revoked = () => assert.ok('onRevoke action is triggered');

    await render(hbs`<Page::PkiCertificateDetails @model={{this.model}} @onRevoke={{this.revoked}} />`, {
      owner: this.engine,
    });
    await click('[data-test-confirm-action-trigger]');
    await click('[data-test-confirm-button]');
  });
});
