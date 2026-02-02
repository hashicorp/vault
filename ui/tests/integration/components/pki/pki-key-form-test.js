/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_KEY_FORM, PKI_KEYS } from 'vault/tests/helpers/pki/pki-selectors';
import PkiKeyForm from 'vault/forms/secrets/pki/key';

module('Integration | Component | pki key form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki'); // https://github.com/ember-engines/ember-engines/pull/653

  hooks.beforeEach(function () {
    this.backend = 'pki-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const { secrets } = this.owner.lookup('service:api');
    this.response = { key_id: 'test-key-id' };
    this.writeStub = sinon.stub(secrets, 'pkiWriteKey').resolves(this.response);
    this.genInternalStub = sinon.stub(secrets, 'pkiGenerateInternalKey').resolves(this.response);
    this.genExportedStub = sinon
      .stub(secrets, 'pkiGenerateExportedKey')
      .resolves({ ...this.response, private_key: 'private-key' });

    this.capabilitiesStub = sinon
      .stub(this.owner.lookup('service:capabilities'), 'for')
      .resolves({ canUpdate: true, canDelete: true });

    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.form = new PkiKeyForm({}, { isNew: true });

    this.renderComponent = () => render(hbs`<PkiKeyForm @form={{this.form}} />`, { owner: this.engine });
  });

  test('it should render fields and show validation messages', async function (assert) {
    assert.expect(7);

    await this.renderComponent();

    assert.dom(GENERAL.inputByAttr('key_name')).exists('renders name input');
    assert.dom(GENERAL.inputByAttr('type')).exists('renders type input');
    assert.dom(GENERAL.inputByAttr('key_type')).exists('renders key type input');
    assert.dom(GENERAL.inputByAttr('key_bits')).exists('renders key bits input');

    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('type'))
      .hasTextContaining('Type is required.', 'renders presence validation for type of key');
    assert
      .dom(GENERAL.validationErrorByAttr('key_type'))
      .hasTextContaining('Please select a key type.', 'renders selection prompt for key type');
    assert
      .dom(PKI_KEY_FORM.validationError)
      .hasTextContaining('There are 2 errors with this form.', 'renders correct form error count');
  });

  test('it generates a key type=exported', async function (assert) {
    assert.expect(9);

    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('key_name'), 'test-key');
    await fillIn(GENERAL.inputByAttr('type'), 'exported');
    assert.dom(GENERAL.inputByAttr('key_bits')).isDisabled('key bits disabled when no key type selected');
    await fillIn(GENERAL.inputByAttr('key_type'), 'rsa');
    await click(GENERAL.submitButton);

    assert.true(
      this.genExportedStub.calledWith(this.backend, {
        key_name: 'test-key',
        key_type: 'rsa',
        key_bits: 2048,
      }),
      'generates exported key with correct params'
    );
    assert.true(
      this.transitionStub.notCalled,
      'does not transition to key details when private_key is returned'
    );
    assert.true(
      this.capabilitiesStub.calledWith('pkiKey', { backend: this.backend, keyId: this.response.key_id }),
      'checks capabilities for new key'
    );
    assert.dom(PKI_KEYS.keyDeleteButton).exists('renders delete button for new key after generation');
    assert.dom(GENERAL.button('Download')).exists('renders download button for private key after generation');
    assert.dom(PKI_KEYS.keyEditLink).exists('renders edit link for new key after generation');
    assert.dom(GENERAL.infoRowValue('Key ID')).hasText(this.response.key_id, 'key id renders');
    assert.dom('[data-test-certificate-card]').exists('Certificate card renders for the private key');
  });

  test('it generates a key type=internal', async function (assert) {
    assert.expect(3);

    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('key_name'), 'test-key');
    await fillIn(GENERAL.inputByAttr('type'), 'internal');
    assert.dom(GENERAL.inputByAttr('key_bits')).isDisabled('key bits disabled when no key type selected');
    await fillIn(GENERAL.inputByAttr('key_type'), 'rsa');
    await click(GENERAL.submitButton);

    assert.true(
      this.genInternalStub.calledWith(this.backend, {
        key_name: 'test-key',
        key_type: 'rsa',
        key_bits: 2048,
      }),
      'generates internal key with correct params'
    );
    assert.true(
      this.transitionStub.calledWith(
        'vault.cluster.secrets.backend.pki.keys.key.details',
        this.response.key_id
      ),
      'transitions to key details page on save'
    );
  });

  test('it should edit key', async function (assert) {
    this.form = new PkiKeyForm({ key_id: 'test-edit-id', key_name: 'FooBar', key_type: 'rsa' });

    await this.renderComponent();

    assert.dom(GENERAL.inputByAttr('key_name')).hasValue('FooBar', 'name field has correct initial value');
    assert.dom(GENERAL.inputByAttr('key_type')).hasValue('rsa', 'key type field has correct initial value');
    assert.dom(GENERAL.inputByAttr('key_type')).isDisabled('key type is not editable');

    await fillIn(GENERAL.inputByAttr('key_name'), 'BarBaz');
    await click(GENERAL.submitButton);

    assert.true(
      this.writeStub.calledWith('test-edit-id', this.backend, {
        key_name: 'BarBaz',
        key_type: 'rsa',
      }),
      'updates key with correct params'
    );
    assert.true(
      this.transitionStub.calledWith(
        'vault.cluster.secrets.backend.pki.keys.key.details',
        this.response.key_id
      ),
      'transitions to key details page on save'
    );
  });
});
