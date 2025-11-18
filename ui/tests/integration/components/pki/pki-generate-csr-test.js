/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import PkiConfigGenerateForm from 'vault/forms/secrets/pki/config/generate';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import sinon from 'sinon';

module('Integration | Component | pki-generate-csr', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.owner.lookup('service:secretMountPath').update('pki-test');

    this.capabilitiesForStub = sinon
      .stub(this.owner.lookup('service:capabilities'), 'for')
      .resolves({ canCreate: false });

    const api = this.owner.lookup('service:api');
    this.issuersGenerateStub = sinon.stub(api.secrets, 'pkiIssuersGenerateIntermediate');
    this.generateStub = sinon.stub(api.secrets, 'pkiGenerateIntermediate');

    setRunOptions({
      rules: {
        // something strange happening here
        'link-name': { enabled: false },
      },
    });
    this.clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();

    this.form = new PkiConfigGenerateForm('PkiGenerateIntermediateRequest', {}, { isNew: true });
    this.onCancel = sinon.stub();
    this.onComplete = sinon.stub();
    this.onSave = sinon.stub();

    this.renderComponent = () =>
      render(
        hbs`
        <PkiGenerateCsr
          @form={{this.form}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
          @onComplete={{this.onComplete}}
        />
      `,
        {
          owner: this.engine,
        }
      );

    this.fillInAndSubmit = async (type) => {
      await fillIn(GENERAL.inputByAttr('type'), type);
      await fillIn(GENERAL.inputByAttr('common_name'), 'foo');
      await click(GENERAL.submitButton);
    };
  });

  hooks.afterEach(function () {
    sinon.restore(); // resets all stubs, including clipboard
  });

  test('it should render fields and save', async function (assert) {
    assert.expect(12);

    await this.renderComponent();

    const fields = [
      'type',
      'common_name',
      'exclude_cn_from_sans',
      'format',
      'subject_serial_number',
      'add_basic_constraints',
    ];
    fields.forEach((key) => {
      assert.dom(`[data-test-input="${key}"]`).exists(`${key} form field renders`);
    });

    assert.dom(GENERAL.button('Key parameters')).exists('Key parameters toggle renders');
    assert.dom(GENERAL.button('Subject Alternative Name (SAN) Options')).exists('SAN options toggle renders');
    assert
      .dom(GENERAL.button('Additional subject fields'))
      .exists('Additional subject fields toggle renders');

    await this.fillInAndSubmit('exported');
    assert.true(this.onSave.calledOnce, 'onSave action fires');
    assert.true(
      this.generateStub.calledWith('exported', 'pki-test'),
      'generateStub called with correct params'
    );
    assert.strictEqual(
      this.generateStub.lastCall.args[2].common_name,
      'foo',
      'common_name is sent in payload'
    );
  });

  test('it should display validation errors', async function (assert) {
    assert.expect(4);

    await this.renderComponent();
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('type'))
      .hasText('Type is required.', 'Type validation error renders');
    assert
      .dom(GENERAL.validationErrorByAttr('common_name'))
      .hasText('Common name is required.', 'Common name validation error renders');
    assert.dom('[data-test-alert]').hasText('There are 2 errors with this form.', 'Alert renders');

    await click('[data-test-cancel]');
    assert.ok(this.onCancel.calledOnce, 'onCancel action fires');
  });

  test('it should use correct endpoint based on issuer permissions', async function (assert) {
    assert.expect(4);

    this.capabilitiesForStub.resolves({ canCreate: true });
    this.issuersGenerateStub.rejects(getErrorResponse());

    await this.renderComponent();
    await this.fillInAndSubmit('exported');

    assert.true(
      this.capabilitiesForStub.calledWith('pkiIssuersGenerateIntermediate', {
        backend: 'pki-test',
        type: 'exported',
      })
    );
    assert.true(this.issuersGenerateStub.calledOnce, 'issuers generate endpoint used when permitted');
    assert.dom(GENERAL.messageError).exists('error message shown');

    this.capabilitiesForStub.resolves({ canCreate: false });
    this.generateStub.resolves({});
    await click(GENERAL.submitButton);

    assert.true(this.generateStub.calledOnce, 'generate endpoint used when there is no issuer permission');
  });

  test('it should show generated CSR for type=exported', async function (assert) {
    assert.expect(6);

    const data = {
      csr: '-----BEGIN CERTIFICATE REQUEST-----...-----END CERTIFICATE REQUEST-----',
      key_id: '9179de78-1275-a1cf-ebb0-a4eb2e376636',
      private_key: '-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----',
      private_key_type: 'rsa',
    };
    this.generateStub.resolves(data);

    await this.renderComponent();
    await this.fillInAndSubmit('exported');

    assert
      .dom('[data-test-next-steps-csr]')
      .hasText(
        'Next steps Copy the CSR below for a parent issuer to sign and then import the signed certificate back into this mount. The private_key is only available once. Make sure you copy and save it now.',
        'renders Next steps alert banner'
      );

    await click(`${GENERAL.infoRowValue('CSR')} ${GENERAL.copyButton}`);
    assert.strictEqual(this.clipboardSpy.firstCall.args[0], data.csr, 'copy value is csr');

    await click(`${GENERAL.infoRowValue('Key ID')} ${GENERAL.copyButton}`);
    assert.strictEqual(this.clipboardSpy.secondCall.args[0], data.key_id, 'copy value is key_id');

    await click(`${GENERAL.infoRowValue('Private key')} ${GENERAL.copyButton}`);
    assert.strictEqual(this.clipboardSpy.thirdCall.args[0], data.private_key, 'copy value is private_key');

    assert
      .dom(GENERAL.infoRowValue('Private key type'))
      .hasText(data.private_key_type, 'renders private_key_type');

    await click('[data-test-done]');
    assert.ok(this.onComplete.calledOnce, 'onComplete action fires');
  });

  test('it should show generated CSR for type=internal', async function (assert) {
    assert.expect(5);

    const data = {
      csr: '-----BEGIN CERTIFICATE REQUEST-----...-----END CERTIFICATE REQUEST-----',
      key_id: '9179de78-1275-a1cf-ebb0-a4eb2e376636',
    };
    this.generateStub.resolves(data);

    await this.renderComponent();
    await this.fillInAndSubmit('internal');

    assert
      .dom('[data-test-next-steps-csr]')
      .hasText(
        'Next steps Copy the CSR below for a parent issuer to sign and then import the signed certificate back into this mount.',
        'renders Next steps alert banner'
      );
    await click(`${GENERAL.infoRowValue('CSR')} ${GENERAL.copyButton}`);
    assert.strictEqual(this.clipboardSpy.firstCall.args[0], data.csr, 'copy value is csr');

    await click(`${GENERAL.infoRowValue('Key ID')} ${GENERAL.copyButton}`);
    assert.strictEqual(this.clipboardSpy.secondCall.args[0], data.key_id, 'copy value is key_id');

    assert.dom(GENERAL.infoRowValue('Private key')).hasText('internal', 'does not render private key');
    assert
      .dom(GENERAL.infoRowValue('Private key type'))
      .hasText('internal', 'does not render private key type');
  });
});
