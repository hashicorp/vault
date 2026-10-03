/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { PKI_CONFIGURE_CREATE } from 'vault/tests/helpers/pki/pki-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const { issuerPemBundle } = CERTIFICATES;
module('Integration | Component | PkiImportPemBundle', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki'); // https://github.com/ember-engines/ember-engines/pull/653

  hooks.beforeEach(function () {
    this.backend = 'pki-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.pemBundle = issuerPemBundle;
    this.onComplete = sinon.stub();
    this.onCancel = sinon.stub();
    this.onSave = sinon.stub();

    const response = {
      request_id: 'test',
      data: {
        mapping: { 'issuer-id': 'key-id' },
      },
    };
    const api = this.owner.lookup('service:api');
    this.issuersImportStub = sinon.stub(api.secrets, 'pkiIssuersImportBundle').resolves(response);
    this.importStub = sinon.stub(api.secrets, 'pkiConfigureCa').resolves(response);

    this.renderComponent = () =>
      render(
        hbs`
        <PkiImportPemBundle
          @useIssuer={{this.useIssuer}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
          @onComplete={{this.onComplete}}
      />`,
        { owner: this.engine }
      );
  });

  test('it renders import form', async function (assert) {
    assert.expect(2);

    await this.renderComponent();

    assert.dom('[data-test-pki-import-pem-bundle-form]').exists('renders form');
    assert.dom('[data-test-component="text-file"]').exists('renders text file input');
  });

  test('it sends correct payload to import endpoint', async function (assert) {
    assert.expect(2);

    await this.renderComponent();

    await click(GENERAL.textToggle);
    await fillIn(GENERAL.maskedInput, this.pemBundle);
    await click(PKI_CONFIGURE_CREATE.importSubmit);
    assert.true(this.importStub.calledWith(this.backend, { pem_bundle: this.pemBundle }));
    assert.true(this.onSave.calledOnce, 'onSave callback fires on save success');
  });

  test('it hits correct endpoint when userIssuer=true', async function (assert) {
    assert.expect(2);

    this.useIssuer = true;
    await this.renderComponent();

    await click(GENERAL.textToggle);
    await fillIn(GENERAL.maskedInput, this.pemBundle);
    await click(PKI_CONFIGURE_CREATE.importSubmit);
    assert.true(this.issuersImportStub.calledWith(this.backend, { pem_bundle: this.pemBundle }));
    assert.true(this.onSave.calledOnce, 'onSave callback fires on save success');
  });

  test('it shows the bundle mapping on success', async function (assert) {
    assert.expect(10);

    this.importStub.resolves({
      imported_issuers: ['issuer-id', 'another-issuer'],
      imported_keys: ['key-id', 'another-key'],
      mapping: { 'issuer-id': 'key-id', 'another-issuer': null },
    });
    await this.renderComponent();

    await click(GENERAL.textToggle);
    await fillIn(GENERAL.maskedInput, this.pemBundle);
    await click(PKI_CONFIGURE_CREATE.importSubmit);

    assert.true(this.onSave.calledOnce, 'onSave callback fires on save success');
    assert.true(this.importStub.calledOnce, 'import endpoint called once');

    assert
      .dom('[data-test-import-pair]')
      .exists({ count: 3 }, 'Shows correct number of rows for imported items');
    // Check that each row has expected values
    assert.dom('[data-test-import-pair="issuer-id_key-id"] [data-test-imported-issuer]').hasText('issuer-id');
    assert.dom('[data-test-import-pair="issuer-id_key-id"] [data-test-imported-key]').hasText('key-id');
    assert
      .dom('[data-test-import-pair="another-issuer_"] [data-test-imported-issuer]')
      .hasText('another-issuer');
    assert.dom('[data-test-import-pair="another-issuer_"] [data-test-imported-key]').hasText('None');
    assert.dom('[data-test-import-pair="_another-key"] [data-test-imported-issuer]').hasText('None');
    assert.dom('[data-test-import-pair="_another-key"] [data-test-imported-key]').hasText('another-key');
    await click('[data-test-done]');

    assert.true(this.onComplete.calledOnce, 'onComplete callback fires on done button click');
  });

  test('it should fire callback on cancel', async function (assert) {
    assert.expect(1);
    await this.renderComponent();
    await click('[data-test-pki-ca-cert-cancel]');
    assert.true(this.onCancel.calledOnce, 'onCancel callback fires on cancel click');
  });
});
