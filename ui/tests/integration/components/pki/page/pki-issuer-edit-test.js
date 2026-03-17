/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import PkiIssuerForm from 'vault/forms/secrets/pki/issuers/issuer';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

const selectors = {
  leaf: '[data-test-field="leaf_not_after_behavior"] select',
  leafOption: '[data-test-field="leaf_not_after_behavior"] option',
  usageCert: '[data-test-usage="Issuing certificates"]',
  usageCrl: '[data-test-usage="Signing CRLs"]',
  usageOcsp: '[data-test-usage="Signing OCSPs"]',
  manualChain: '[data-test-input="manual_chain"] [data-test-string-list-input="0"]',
  certUrls: '[data-test-input="issuing_certificates"][data-test-string-list-input]',
  certUrl1: '[data-test-input="issuing_certificates"] [data-test-string-list-input="0"]',
  certUrl2: '[data-test-input="issuing_certificates"] [data-test-string-list-input="1"]',
  certUrlAdd: '[data-test-input="issuing_certificates"] [data-test-string-list-button="add"]',
  certUrlRemove: '[data-test-input="issuing_certificates"] [data-test-string-list-button="delete"]',
  crlDist: '[data-test-input="crl_distribution_points"] [data-test-string-list-input="0"]',
  ocspServers: '[data-test-input="ocsp_servers"]  [data-test-string-list-input="0"]',
  cancel: '[data-test-cancel]',
  error: '[data-test-message-error]',
  alert: '[data-test-inline-error-message]',
};

module('Integration | Component | pki | Page::PkiIssuerEditPage::PkiIssuerEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    const router = this.owner.lookup('service:router');
    const transitionSpy = sinon.stub(router, 'transitionTo');
    this.transitionCalled = () =>
      transitionSpy.calledWith('vault.cluster.secrets.backend.pki.issuers.issuer.details');

    this.writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'pkiWriteIssuer').resolves();

    // value pulled from secretMountPath service for write request
    this.backend = 'pki';
    this.owner.lookup('service:secretMountPath').update(this.backend);

    this.issuerRef = 'test-issuer';
    this.data = {
      issuer_id: 'test',
      issuer_name: 'foo-bar',
      leaf_not_after_behavior: 'err',
      usage: 'read-only,issuing-certificates,ocsp-signing',
      manual_chain: 'issuer_ref',
      issuing_certificates: ['http://localhost', 'http://localhost:8200'],
      crl_distribution_points: ['http://localhost'],
      ocsp_servers: ['http://localhost'],
    };
    this.form = new PkiIssuerForm(this.data);

    this.renderComponent = () =>
      render(hbs`<Page::PkiIssuerEdit @form={{this.form}} @issuerRef={{this.issuerRef}} />`, {
        owner: this.engine,
      });

    this.update = async () => {
      await fillIn(GENERAL.inputByAttr('issuer_name'), 'bar-baz');
      await click(selectors.usageCrl);
      await click(selectors.certUrlRemove);
    };
  });

  test('it should populate fields with values', async function (assert) {
    await this.renderComponent();

    const {
      issuer_name,
      leaf_not_after_behavior,
      manual_chain,
      issuing_certificates,
      crl_distribution_points,
      ocsp_servers,
    } = this.data;

    assert.dom(GENERAL.inputByAttr('issuer_name')).hasValue(issuer_name, 'Issuer name field populates');
    assert.dom(selectors.leaf).hasValue(leaf_not_after_behavior, 'Leaf not after behavior option selected');
    assert
      .dom(selectors.leafOption)
      .hasText(
        'Error if the computed NotAfter exceeds that of this issuer in all circumstances (leaf, CA issuance and ACME)',
        'Correct text renders for leaf option'
      );
    assert.dom(selectors.usageCert).isChecked('Usage issuing certificates is checked');
    assert.dom(selectors.usageCrl).isNotChecked('Usage signing crls is not checked');
    assert.dom(selectors.usageOcsp).isChecked('Usage signing ocsps is checked');
    assert.dom(selectors.manualChain).hasValue(manual_chain, 'Manual chain field populates');
    assert.dom(selectors.certUrl1).hasValue(issuing_certificates[0], 'Issuing certificate populates');
    assert.dom(selectors.certUrl2).hasValue(issuing_certificates[1], 'Issuing certificate populates');
    assert.dom(selectors.crlDist).hasValue(crl_distribution_points[0], 'Crl distribution points populate');
    assert.dom(selectors.ocspServers).hasValue(ocsp_servers[0], 'Ocsp servers populate');
  });

  test('it should update issuer', async function (assert) {
    assert.expect(2);

    await this.renderComponent();
    await this.update();
    await click(GENERAL.submitButton);

    const payload = {
      ...this.data,
      issuer_name: 'bar-baz',
      usage: 'read-only,issuing-certificates,ocsp-signing,crl-signing',
      issuing_certificates: ['http://localhost:8200'],
    };
    assert.true(
      this.writeStub.calledWith(this.issuerRef, this.backend, payload),
      'API called with updated issuer data'
    );
    assert.ok(this.transitionCalled(), 'Transitions to details route on save success');
  });

  test('it should show error messages', async function (assert) {
    this.writeStub.rejects(getErrorResponse({ errors: ['Some error occurred'] }, 400));

    await this.renderComponent();
    await click(GENERAL.submitButton);

    assert
      .dom(selectors.alert)
      .hasText('There was an error submitting this form.', 'Inline error alert renders');
    assert.dom(selectors.error).hasTextContaining('Some error occurred', 'Error message renders');
  });
});
