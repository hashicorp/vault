/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
  toggleSigningOptions: '[data-test-toggle-group="Signing options"]',
  toggleSANOptions: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  toggleAdditionalFields: '[data-test-toggle-group="Additional subject fields"]',
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
    assert.expect(9);
    await render(hbs`<PkiSignIntermediateForm @onCancel={{this.onCancel}} @model={{this.model}} />`, {
      owner: this.engine,
    });

    assert.dom(selectors.form).exists('Form is rendered');
    assert.dom(selectors.resultsContainer).doesNotExist('Results display not rendered');
    assert.dom('[data-test-field]').exists({ count: 8 }, '8 default fields shown');
    assert.dom(selectors.toggleSigningOptions).exists();
    assert.dom(selectors.toggleSANOptions).exists();
    assert.dom(selectors.toggleAdditionalFields).exists();

    await click(selectors.toggleSigningOptions);
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
          ca_chain: ['-----root pem------'],
          issuing_ca: '-----issuing ca------',
          certificate: '-----certificate------',
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
      { label: 'CA Chain', masked: true },
      { label: 'Certificate', masked: true },
      { label: 'Issuing CA', masked: true },
    ].forEach(({ label, masked }) => {
      assert.dom(selectors.rowByName(label)).exists();
      if (masked) {
        assert.dom(selectors.valueByName(label)).hasText('***********', `${label} is masked`);
      } else {
        assert.dom(selectors.valueByName(label)).hasText('31:52:b9:09:40', `Renders ${label}`);
        assert.dom(`${selectors.valueByName(label)} a`).exists(`${label} is a link`);
      }
    });
  });
});
