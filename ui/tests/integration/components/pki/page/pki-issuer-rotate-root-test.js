/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { hbs } from 'ember-cli-htmlbars';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { PKI_CONFIGURE_CREATE, PKI_CONFIG_EDIT } from 'vault/tests/helpers/pki/pki-selectors';

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
    this.api = this.owner.lookup('service:api');

    this.onCancel = sinon.spy();
    this.onComplete = sinon.spy();
    this.rotateStub = sinon.stub(this.api.secrets, 'pkiRotateRoot').resolves();

    this.breadcrumbs = [{ label: 'rotate root' }];
    this.oldRootData = {
      certificate: loadedCert,
      issuer_id: 'old-issuer-id',
      issuer_name: 'old-issuer',
    };
    this.store.pushPayload('pki/issuer', { modelName: 'pki/issuer', data: this.oldRootData });
    this.oldRoot = this.store.peekRecord('pki/issuer', 'old-issuer-id');
    this.certData = parseCertificate(loadedCert);

    this.returnedData = {
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

    this.renderComponent = () =>
      render(
        hbs`
          <Page::PkiIssuerRotateRoot
            @oldRoot={{this.oldRoot}}
            @certData={{this.certData}}
            @parsingErrors={{this.parsingErrors}}
            @breadcrumbs={{this.breadcrumbs}}
            @onCancel={{this.onCancel}}
            @onComplete={{this.onComplete}}
          />
      `,
        { owner: this.engine }
      );

    this.customizeAndSubmit = async () => {
      await click(SELECTORS.customRadioSelect);
      await fillIn(GENERAL.inputByAttr('type'), 'exported');
      await fillIn(GENERAL.inputByAttr('common_name'), 'foo');
      await click(GENERAL.submitButton);
    };
  });

  test('it renders', async function (assert) {
    assert.expect(17);

    await this.renderComponent();

    assert.dom(GENERAL.title).hasText('Generate New Root');
    assert.dom(SELECTORS.oldRadioSelect).isChecked('defaults to use-old-settings');
    assert.dom(SELECTORS.rotateRootForm).exists('it renders old settings form');
    assert
      .dom(GENERAL.inputByAttr('common_name'))
      .hasValue(this.certData.common_name, 'common name prefilled with root cert cn');
    assert.dom(SELECTORS.toggle).hasText('Old root settings', 'toggle renders correct text');
    assert.dom(GENERAL.inputByAttr('issuer_name')).exists('renders issuer name input');
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
      .dom(PKI_CONFIG_EDIT.stringListInput('permitted_dns_domains', 0))
      .hasValue(
        this.certData.permitted_dns_domains.split(',')[0],
        'form is prefilled with values from old root'
      );
    await click(GENERAL.cancelButton);
    assert.ok(this.onCancel.calledOnce, 'custom form calls @onCancel passed from parent');
    await click(SELECTORS.oldRadioSelect);
    await click(GENERAL.cancelButton);
    assert.ok(this.onCancel.calledTwice, 'old root settings form calls @onCancel from parent');

    // validations
    await fillIn(GENERAL.inputByAttr('common_name'), '');
    await fillIn(GENERAL.inputByAttr('issuer_name'), 'default');
    await click(GENERAL.submitButton);
    assert.dom(SELECTORS.validationError).hasText('There are 2 errors with this form.');
    assert.dom(GENERAL.validationErrorByAttr('common_name')).exists();
    assert.dom(GENERAL.validationErrorByAttr('issuer_name')).exists();
  });

  test('it sends request to rotate/internal on save when using old root settings', async function (assert) {
    assert.expect(1);

    await this.renderComponent();

    await click(GENERAL.submitButton);

    assert.true(
      this.rotateStub.calledWith('internal', this.backend),
      'rotateStub called with correct params'
    );
  });

  test.each(
    `it sends request to rotate with correct type on save with custom root settings`,
    ['internal', 'exported', 'existing', 'kms'],
    async function (assert, type) {
      assert.expect(1);

      await this.renderComponent();

      await click(SELECTORS.customRadioSelect);
      await fillIn(GENERAL.inputByAttr('type'), type);
      await click(GENERAL.submitButton);

      assert.true(this.rotateStub.calledWith(type, this.backend), 'rotateStub called with correct params');
    }
  );

  test('it renders details after save for exported key type', async function (assert) {
    assert.expect(10);

    this.rotateStub.resolves({
      ...this.returnedData,
      private_key: `-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAtc9yU`,
      private_key_type: 'rsa',
    });

    await this.renderComponent();
    await this.customizeAndSubmit();

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
    assert.dom(`${GENERAL.infoRowValue('Private key type')}`).hasText('rsa');
    assert.dom(GENERAL.infoRowValue('Serial number')).hasText(this.returnedData.serial_number);
    assert.dom(GENERAL.infoRowValue('Key ID')).hasText(this.returnedData.key_id);

    await click(PKI_CONFIGURE_CREATE.doneButton);
    assert.true(this.onComplete.calledOnce, 'clicking done fires @onComplete from parent');
  });

  test('it renders details after save for internal key type', async function (assert) {
    assert.expect(13);

    this.rotateStub.resolves(this.returnedData);

    await this.renderComponent();
    await this.customizeAndSubmit();

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
    assert.dom(`${GENERAL.infoRowValue('Private key')} div`).hasText('internal');
    assert.dom(`${GENERAL.infoRowValue('Private key type')} div`).hasText('internal');
    assert.dom(GENERAL.infoRowValue('Serial number')).hasText(this.returnedData.serial_number);
    assert.dom(GENERAL.infoRowValue('Key ID')).hasText(this.returnedData.key_id);

    await click(SELECTORS.doneButton);
    assert.ok(this.onComplete.calledOnce, 'clicking done fires @onComplete from parent');
  });
});
