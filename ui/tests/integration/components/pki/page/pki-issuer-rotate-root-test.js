/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { hbs } from 'ember-cli-htmlbars';
import camelizeKeys from 'vault/utils/camelize-object-keys';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { PKI_CONFIGURE_CREATE } from 'vault/tests/helpers/pki/pki-selectors';

const SELECTORS = {
  nextSteps: '[data-test-rotate-next-steps]',
  toolbarCrossSign: '[data-test-pki-issuer-cross-sign]',
  toolbarSignInt: '[data-test-pki-issuer-sign-int]',
  toolbarDownload: '[data-test-issuer-download]',
  oldRadioSelect: 'input#use-old-root-settings',
  customRadioSelect: 'input#customize-new-root-certificate',
  toggle: '[data-test-details-toggle]',
  validationError: '[data-test-pki-rotate-root-validation-error]',
  rotateRootForm: '[data-test-pki-rotate-old-settings-form]',
  doneButton: '[data-test-done]',
  // root form
  generateRootForm: '[data-test-pki-config-generate-root-form]',
};
const { loadedCert } = CERTIFICATES;

module('Integration | Component | page/pki-issuer-rotate-root', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'test-pki';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.onCancel = sinon.spy();
    this.onComplete = sinon.spy();
    this.breadcrumbs = [{ label: 'rotate root' }];
    this.oldRootData = {
      certificate: loadedCert,
      issuer_id: 'old-issuer-id',
      issuer_name: 'old-issuer',
    };
    this.parsedRootCert = camelizeKeys(parseCertificate(loadedCert));
    this.store.pushPayload('pki/issuer', { modelName: 'pki/issuer', data: this.oldRootData });
    this.oldRoot = this.store.peekRecord('pki/issuer', 'old-issuer-id');
    this.newRootModel = this.store.createRecord('pki/action', {
      actionType: 'rotate-root',
      type: 'internal',
      ...this.parsedRootCert, // copy old root settings over to new one
    });

    this.returnedData = {
      id: 'response-id',
      certificate: loadedCert,
      expiration: 1682735724,
      issuer_id: 'some-issuer-id',
      issuer_name: 'my issuer',
      issuing_ca: loadedCert,
      key_id: 'my-key-id',
      key_name: 'my-key',
      serial_number: '3a:3c:17:..',
    };
    setRunOptions({
      rules: {
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it renders', async function (assert) {
    assert.expect(17);
    await render(
      hbs`
      <Page::PkiIssuerRotateRoot
        @oldRoot={{this.oldRoot}}
        @newRootModel={{this.newRootModel}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onComplete={{this.onComplete}}
      />
    `,
      { owner: this.engine }
    );
    assert.dom(GENERAL.title).hasText('Generate New Root');
    assert.dom(SELECTORS.oldRadioSelect).isChecked('defaults to use-old-settings');
    assert.dom(SELECTORS.rotateRootForm).exists('it renders old settings form');
    assert
      .dom(GENERAL.inputByAttr('commonName'))
      .hasValue(this.parsedRootCert.commonName, 'common name prefilled with root cert cn');
    assert.dom(SELECTORS.toggle).hasText('Old root settings', 'toggle renders correct text');
    assert.dom(GENERAL.inputByAttr('issuerName')).exists('renders issuer name input');
    assert.strictEqual(findAll('[data-test-row-label]').length, 0, 'it hides the old root info table rows');
    await click(SELECTORS.toggle);
    assert.strictEqual(findAll('[data-test-row-label]').length, 19, 'it shows the old root info table rows');
    assert
      .dom(GENERAL.infoRowValue('Issuer name'))
      .hasText(this.oldRoot.issuerName, 'renders correct issuer data');
    await click(SELECTORS.toggle);
    assert.strictEqual(findAll('[data-test-row-label]').length, 0, 'it hides again');

    // customize form
    await click(SELECTORS.customRadioSelect);
    assert.dom(SELECTORS.generateRootForm).exists('it renders generate root form');
    assert
      .dom(GENERAL.inputByAttr('permittedDnsDomains'))
      .hasValue(this.parsedRootCert.permittedDnsDomains, 'form is prefilled with values from old root');
    await click(GENERAL.cancelButton);
    assert.ok(this.onCancel.calledOnce, 'custom form calls @onCancel passed from parent');
    await click(SELECTORS.oldRadioSelect);
    await click(GENERAL.cancelButton);
    assert.ok(this.onCancel.calledTwice, 'old root settings form calls @onCancel from parent');

    // validations
    await fillIn(GENERAL.inputByAttr('commonName'), '');
    await fillIn(GENERAL.inputByAttr('issuerName'), 'default');
    await click(GENERAL.saveButton);
    assert.dom(SELECTORS.validationError).hasText('There are 2 errors with this form.');
    assert
      .dom(GENERAL.inputByAttr('commonName'))
      .hasClass('has-error-border', 'common name has error border');
    assert
      .dom(GENERAL.inputByAttr('issuerName'))
      .hasClass('has-error-border', 'issuer name has error border');
  });

  test('it sends request to rotate/internal on save when using old root settings', async function (assert) {
    assert.expect(1);
    this.server.post(`/${this.backend}/root/rotate/internal`, () => {
      assert.ok('request made to correct default endpoint type=internal');
    });
    await render(
      hbs`
      <Page::PkiIssuerRotateRoot
        @oldRoot={{this.oldRoot}}
        @newRootModel={{this.newRootModel}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onComplete={{this.onComplete}}
      />
    `,
      { owner: this.engine }
    );
    await click(GENERAL.saveButton);
  });

  function testEndpoint(test, type) {
    test(`it sends request to rotate/${type} endpoint on save with custom root settings`, async function (assert) {
      assert.expect(1);
      this.server.post(`/${this.backend}/root/rotate/${type}`, () => {
        assert.ok('request is made to correct endpoint');
      });
      await render(
        hbs`
      <Page::PkiIssuerRotateRoot
        @oldRoot={{this.oldRoot}}
        @newRootModel={{this.newRootModel}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onComplete={{this.onComplete}}
      />
    `,
        { owner: this.engine }
      );
      await click(SELECTORS.customRadioSelect);
      await fillIn(GENERAL.inputByAttr('type'), type);
      await click(GENERAL.saveButton);
    });
  }
  testEndpoint(test, 'internal');
  testEndpoint(test, 'exported');
  testEndpoint(test, 'existing');
  testEndpoint(test, 'kms');

  test('it renders details after save for exported key type', async function (assert) {
    assert.expect(10);
    const keyData = {
      private_key: `-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAtc9yU`,
      private_key_type: 'rsa',
    };
    this.store.pushPayload('pki/action', {
      modelName: 'pki/action',
      data: { ...this.returnedData, ...keyData },
    });
    this.newRootModel = this.store.peekRecord('pki/action', 'response-id');
    await render(
      hbs`
        <Page::PkiIssuerRotateRoot
          @oldRoot={{this.oldRoot}}
          @newRootModel={{this.newRootModel}}
          @breadcrumbs={{this.breadcrumbs}}
          @onCancel={{this.onCancel}}
          @onComplete={{this.onComplete}}
        />
      `,
      { owner: this.engine }
    );
    assert.dom(GENERAL.title).hasText('View Issuer Certificate');
    assert
      .dom(SELECTORS.nextSteps)
      .hasText(
        'Next steps Your new root has been generated. Make sure to copy and save the private_key as it is only available once. If you’re ready, you can begin cross-signing issuers now. If not, the option to cross-sign is available when you use this certificate. Cross-sign issuers'
      );
    assert.dom(GENERAL.infoRowValue('Certificate')).exists();
    assert.dom(GENERAL.infoRowValue('Issuer name')).exists();
    assert.dom(GENERAL.infoRowValue('Issuing CA')).exists();
    assert.dom(GENERAL.infoRowValue('Private key')).exists();
    assert.dom(`${GENERAL.infoRowValue('Private key type')} span`).hasText('rsa');
    assert.dom(GENERAL.infoRowValue('Serial number')).hasText(this.returnedData.serial_number);
    assert.dom(GENERAL.infoRowValue('Key ID')).hasText(this.returnedData.key_id);

    await click(PKI_CONFIGURE_CREATE.doneButton);
    assert.ok(this.onComplete.calledOnce, 'clicking done fires @onComplete from parent');
  });

  test('it renders details after save for internal key type', async function (assert) {
    assert.expect(13);
    this.store.pushPayload('pki/action', {
      modelName: 'pki/action',
      data: this.returnedData,
    });
    this.newRootModel = this.store.peekRecord('pki/action', 'response-id');
    await render(
      hbs`
        <Page::PkiIssuerRotateRoot
          @oldRoot={{this.oldRoot}}
          @newRootModel={{this.newRootModel}}
          @breadcrumbs={{this.breadcrumbs}}
          @onCancel={{this.onCancel}}
          @onComplete={{this.onComplete}}
        />
      `,
      { owner: this.engine }
    );
    assert.dom(GENERAL.title).hasText('View Issuer Certificate');
    assert.dom(SELECTORS.toolbarCrossSign).exists();
    assert.dom(SELECTORS.toolbarSignInt).exists();
    assert.dom(SELECTORS.toolbarDownload).exists();
    assert
      .dom(SELECTORS.nextSteps)
      .hasText(
        'Next steps Your new root has been generated. If you’re ready, you can begin cross-signing issuers now. If not, the option to cross-sign is available when you use this certificate. Cross-sign issuers'
      );
    assert.dom(GENERAL.infoRowValue('Certificate')).exists();
    assert.dom(GENERAL.infoRowValue('Issuer name')).exists();
    assert.dom(GENERAL.infoRowValue('Issuing CA')).exists();
    assert.dom(`${GENERAL.infoRowValue('Private key')} span`).hasText('internal');
    assert.dom(`${GENERAL.infoRowValue('Private key type')} span`).hasText('internal');
    assert.dom(GENERAL.infoRowValue('Serial number')).hasText(this.returnedData.serial_number);
    assert.dom(GENERAL.infoRowValue('Key ID')).hasText(this.returnedData.key_id);

    await click(SELECTORS.doneButton);
    assert.ok(this.onComplete.calledOnce, 'clicking done fires @onComplete from parent');
  });
});
