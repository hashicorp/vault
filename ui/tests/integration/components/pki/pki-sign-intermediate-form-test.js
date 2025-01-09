/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';

const selectors = {
  form: '[data-test-sign-intermediate-form]',
  csrInput: '[data-test-input="csr"]',
  toggleGroup: (group) => `[data-test-toggle-group="${group}"]`,
  fieldByName: (name) => `[data-test-field="${name}"]`,
  saveButton: '[data-test-pki-sign-intermediate-save]',
  cancelButton: '[data-test-pki-sign-intermediate-cancel]',
  fieldError: '[data-test-inline-alert]',
  formError: '[data-test-form-error]',
  resultsContainer: '[data-test-sign-intermediate-result]',
  rowByName: (name) => `[data-test-row-label="${name}"]`,
  valueByName: (name) => `[data-test-value-div="${name}"]`,
};
module('Integration | Component | pki-sign-intermediate-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.model = this.store.createRecord('pki/sign-intermediate', { issuerRef: 'some-issuer' });
    this.onCancel = Sinon.spy();
  });

  test('renders correctly on load', async function (assert) {
    assert.expect(10);
    await render(hbs`<PkiSignIntermediateForm @onCancel={{this.onCancel}} @model={{this.model}} />`, {
      owner: this.engine,
    });

    assert.dom(selectors.form).exists('Form is rendered');
    assert.dom(selectors.resultsContainer).doesNotExist('Results display not rendered');
    assert.dom('[data-test-field]').exists({ count: 9 }, '9 default fields shown');
    [
      'Name constraints',
      'Signing options',
      'Subject Alternative Name (SAN) Options',
      'Additional subject fields',
    ].forEach((group) => {
      assert.dom(selectors.toggleGroup(group)).exists(`${group} renders`);
    });

    await click(selectors.toggleGroup('Signing options'));
    ['usePss', 'skid', 'signatureBits'].forEach((name) => {
      assert.dom(selectors.fieldByName(name)).exists();
    });
  });

  test('it shows the returned values on successful save', async function (assert) {
    assert.expect(13);
    await render(hbs`<PkiSignIntermediateForm @onCancel={{this.onCancel}} @model={{this.model}} />`, {
      owner: this.engine,
    });

    this.server.post(`/pki-test/issuer/some-issuer/sign-intermediate`, function (schema, req) {
      const payload = JSON.parse(req.requestBody);
      assert.strictEqual(payload.csr, 'example-data', 'Request made to correct endpoint on save');
      return {
        request_id: 'some-id',
        data: {
          serial_number: '31:52:b9:09:40',
          ca_chain: ['-----BEGIN CERTIFICATE-----'],
          issuing_ca: '-----BEGIN CERTIFICATE-----',
          certificate: '-----BEGIN CERTIFICATE-----',
        },
      };
    });
    await click(selectors.saveButton);
    assert.dom(selectors.formError).hasText('There is an error with this form.', 'Shows validation errors');
    assert.dom(selectors.csrInput).hasClass('has-error-border');
    assert.dom(selectors.fieldError).hasText('CSR is required.');

    await fillIn(selectors.csrInput, 'example-data');
    await click(selectors.saveButton);
    [
      { label: 'Serial number' },
      { label: 'CA Chain', isCertificate: true },
      { label: 'Certificate', isCertificate: true },
      { label: 'Issuing CA', isCertificate: true },
    ].forEach(({ label, isCertificate }) => {
      assert.dom(selectors.rowByName(label)).exists();
      if (isCertificate) {
        assert.dom(selectors.valueByName(label)).includesText('PEM Format', `${label} is isCertificate`);
      } else {
        assert.dom(selectors.valueByName(label)).hasText('31:52:b9:09:40', `Renders ${label}`);
        assert.dom(`${selectors.valueByName(label)} a`).exists(`${label} is a link`);
      }
    });
  });
});
