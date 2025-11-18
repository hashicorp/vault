/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import PkiIssuersSignIntermediateForm from 'vault/forms/secrets/pki/issuers/sign-intermediate';

const selectors = {
  form: '[data-test-sign-intermediate-form]',
  saveButton: '[data-test-pki-sign-intermediate-save]',
  cancelButton: '[data-test-pki-sign-intermediate-cancel]',
  formError: '[data-test-form-error]',
  resultsContainer: '[data-test-sign-intermediate-result]',
};

module('Integration | Component | pki-sign-intermediate-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.mountPath = 'pki-test';
    this.owner.lookup('service:secret-mount-path').update(this.mountPath);

    this.form = new PkiIssuersSignIntermediateForm({}, { isNew: true });
    this.issuerRef = 'some-issuer';
    this.onCancel = sinon.spy();

    this.renderComponent = () =>
      render(
        hbs`
      <PkiSignIntermediateForm
        @form={{this.form}}
        @issuerRef={{this.issuerRef}}
        @onCancel={{this.onCancel}}
      />
    `,
        { owner: this.engine }
      );

    this.signStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'pkiIssuerSignIntermediate')
      .resolves({
        serial_number: '31:52:b9:09:40',
        ca_chain: ['-----BEGIN CERTIFICATE-----'],
        issuing_ca: '-----BEGIN CERTIFICATE-----',
        certificate: '-----BEGIN CERTIFICATE-----',
      });
  });

  test('renders correctly on load', async function (assert) {
    assert.expect(10);

    await this.renderComponent();

    assert.dom(selectors.form).exists('Form is rendered');
    assert.dom(selectors.resultsContainer).doesNotExist('Results display not rendered');
    assert.dom('[data-test-field]').exists({ count: 9 }, '9 default fields shown');
    [
      'Name constraints',
      'Signing options',
      'Subject Alternative Name (SAN) Options',
      'Additional subject fields',
    ].forEach((group) => {
      assert.dom(GENERAL.button(group)).exists(`${group} renders`);
    });

    await click(GENERAL.button('Signing options'));
    ['use_pss', 'skid', 'signature_bits'].forEach((name) => {
      assert.dom(GENERAL.fieldByAttr(name)).exists();
    });
  });

  test('it shows the returned values on successful save', async function (assert) {
    assert.expect(12);

    await this.renderComponent();

    await click(selectors.saveButton);
    assert.dom(selectors.formError).hasText('There is an error with this form.', 'Shows validation errors');
    assert.dom(GENERAL.validationErrorByAttr('csr')).hasText('CSR is required.');

    await fillIn(GENERAL.inputByAttr('csr'), 'example-data');
    await click(selectors.saveButton);

    assert.true(
      this.signStub.calledWith(this.issuerRef, this.mountPath, {
        csr: 'example-data',
        format: 'pem',
        not_before_duration: 30,
        private_key_format: 'der',
      }),
      'Request made to correct endpoint on save'
    );

    [
      { label: 'Serial number' },
      { label: 'CA Chain', isCertificate: true },
      { label: 'Certificate', isCertificate: true },
      { label: 'Issuing CA', isCertificate: true },
    ].forEach(({ label, isCertificate }) => {
      assert.dom(GENERAL.infoRowLabel(label)).exists();
      if (isCertificate) {
        assert.dom(GENERAL.infoRowValue(label)).includesText('PEM Format', `${label} is isCertificate`);
      } else {
        assert.dom(GENERAL.infoRowValue(label)).hasText('31:52:b9:09:40', `Renders ${label}`);
        assert.dom(`${GENERAL.infoRowValue(label)} a`).exists(`${label} is a link`);
      }
    });
  });
});
