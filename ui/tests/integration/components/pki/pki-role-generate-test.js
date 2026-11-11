/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, fillIn, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import Sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { TTL_PICKER } from 'vault/tests/helpers/components/ttl-picker-selectors';
import { PKI_ROLE_GENERATE } from 'vault/tests/helpers/pki/pki-selectors';
import PkiCertificateForm from 'vault/forms/secrets/pki/certificate';

module('Integration | Component | pki-role-generate', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.role = 'my-role';
    this.backend = 'pki-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.form = new PkiCertificateForm('PkiIssueWithRoleRequest', {}, { isNew: true });
    this.mode = 'generate';
    this.onSuccess = Sinon.spy();
    const { secrets } = this.owner.lookup('service:api');
    const response = { serial_number: 'abcd-efgh-ijkl', certificate: '---CERT---' };
    this.issueStub = Sinon.stub(secrets, 'pkiIssueWithRole').resolves(response);
    this.signStub = Sinon.stub(secrets, 'pkiSignWithRole').resolves(response);

    this.renderComponent = () =>
      render(
        hbs`<PkiRoleGenerate @form={{this.form}} @role={{this.role}} @mode={{this.mode}} @onSuccess={{this.onSuccess}} />`,
        { owner: this.engine }
      );
  });

  test('it should render the component with the form by default', async function (assert) {
    assert.expect(4);

    await this.renderComponent();
    assert.dom(PKI_ROLE_GENERATE.form).exists('shows the cert generate form');
    assert.dom(GENERAL.inputByAttr('common_name')).exists('shows the common name field');
    assert.dom(GENERAL.button('Subject Alternative Name (SAN) Options')).exists('toggle exists');
    await fillIn(GENERAL.inputByAttr('common_name'), 'example.com');
    assert.strictEqual(this.form.data.common_name, 'example.com', 'Filling in the form updates the model');
  });

  test('it should render correctly for each mode', async function (assert) {
    assert.expect(4);

    this.mode = 'generate';
    await this.renderComponent();
    assert.dom(GENERAL.submitButton).hasText('Generate', 'shows correct text for submit button');
    await click(GENERAL.submitButton);
    assert.true(
      this.issueStub.calledWith(this.role, this.backend),
      'makes request to issue endpoint in generate mode'
    );

    this.mode = 'sign';
    await this.renderComponent();
    assert.dom(GENERAL.submitButton).hasText('Sign', 'shows correct text for submit button');
    await click(GENERAL.submitButton);
    assert.true(
      this.signStub.calledWith(this.role, this.backend),
      'makes request to sign endpoint in sign mode'
    );
  });

  test('it should generate cert and display details', async function (assert) {
    assert.expect(5);

    await this.renderComponent();
    await click(GENERAL.submitButton);
    assert.dom(PKI_ROLE_GENERATE.form).doesNotExist('Does not show the form');
    assert.dom(PKI_ROLE_GENERATE.downloadButton).exists('shows the download button');
    assert.dom(PKI_ROLE_GENERATE.revokeButton).exists('shows the revoke button');
    assert.dom(GENERAL.infoRowValue('Certificate')).exists({ count: 1 }, 'shows certificate info row');
    assert
      .dom(GENERAL.infoRowValue('Serial number'))
      .hasText('abcd-efgh-ijkl', 'shows serial number info row');
  });

  test('it blocks a TTL longer than Vault supports', async function (assert) {
    assert.expect(6);

    await this.renderComponent();
    await fillIn(GENERAL.ttl.input('TTL'), '9999999999');
    await fillIn(TTL_PICKER.unitInputByLabel('TTL'), 'd');
    await click(GENERAL.submitButton);

    assert.dom(TTL_PICKER.error).hasText('TTL is too large. Maximum allowed value is 106,751 days.');
    assert.dom(PKI_ROLE_GENERATE.alert).hasText('There is an error with this form.');
    assert.true(this.issueStub.notCalled, 'no certificate is issued');
    assert.dom(PKI_ROLE_GENERATE.form).exists('the form stays open');

    await fillIn(GENERAL.ttl.input('TTL'), '106751');
    await click(GENERAL.submitButton);
    assert.true(this.issueStub.calledOnce, 'the certificate is issued once the TTL is valid');
    assert.strictEqual(this.issueStub.firstCall.args[2].ttl, '2562024h', 'the maximum whole days is sent');
  });

  test('it blocks a TTL too long to be stored as a number', async function (assert) {
    assert.expect(2);

    await this.renderComponent();
    // The TTL picker turns this many digits into Infinity, which it sends as the string "Infinityh".
    await fillIn(GENERAL.ttl.input('TTL'), '9'.repeat(400));
    await fillIn(TTL_PICKER.unitInputByLabel('TTL'), 'd');
    await click(GENERAL.submitButton);

    assert.dom(PKI_ROLE_GENERATE.alert).hasText('There is an error with this form.');
    assert.true(this.issueStub.notCalled, 'no certificate is issued');
  });
});
