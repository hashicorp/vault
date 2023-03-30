/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { hbs } from 'ember-cli-htmlbars';
import { loadedCert } from 'vault/tests/helpers/pki/values';
import camelizeKeys from 'vault/utils/camelize-object-keys';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import { SELECTORS as S } from 'vault/tests/helpers/pki/pki-generate-root';

const SELECTORS = {
  toggle: '[data-test-details-toggle]',
  oldRadioSelect: 'input#use-old-root-settings',
  customRadioSelect: 'input#customize-new-root-certificate',
  input: (attr) => `[data-test-input="${attr}"]`,
  infoRowValue: (attr) => `[data-test-value-div="${attr}"]`,
  validationError: '[data-test-pki-rotate-root-validation-error]',
  rotateRootForm: '[data-test-pki-rotate-old-settings-form]',
  rotateRootSave: '[data-test-pki-rotate-root-save]',
  rotateRootCancel: '[data-test-pki-rotate-root-cancel]',
  // root form
  generateRootForm: '[data-test-pki-config-generate-root-form]',
  ...S,
};

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
    this.store.pushPayload('pki/issuer', {
      modelName: 'pki/issuer',
      data: this.oldRootData,
    });

    this.oldRoot = this.store.peekRecord('pki/issuer', 'old-issuer-id');
    this.newRootModel = this.store.createRecord('pki/action', {
      actionType: 'rotate-root',
      type: 'internal',
      ...this.parsedRootCert, // copy old root settings over to new one
    });
  });

  test('it renders', async function (assert) {
    assert.expect(16);
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
    assert.dom(SELECTORS.oldRadioSelect).isChecked('defaults to use-old-settings');
    assert.dom(SELECTORS.rotateRootForm).exists('it renders old settings form');
    assert
      .dom(SELECTORS.input('commonName'))
      .hasValue(this.parsedRootCert.commonName, 'common name prefilled with root cert cn');
    assert.dom(SELECTORS.toggle).hasText('Old root settings', 'toggle renders correct text');
    assert.dom(SELECTORS.input('issuerName')).exists('renders issuer name input');
    assert.strictEqual(findAll('[data-test-row-label]').length, 0, 'it hides the old root info table rows');
    await click(SELECTORS.toggle);
    assert.strictEqual(findAll('[data-test-row-label]').length, 11, 'it shows the old root info table rows');
    assert
      .dom(SELECTORS.infoRowValue('Issuer name'))
      .hasText(this.oldRoot.issuerName, 'renders correct issuer data');
    await click(SELECTORS.toggle);
    assert.strictEqual(findAll('[data-test-row-label]').length, 0, 'it hides again');
    await click(SELECTORS.customRadioSelect);
    assert.dom(SELECTORS.generateRootForm).exists('it renders generate root form');
    assert
      .dom(SELECTORS.input('permittedDnsDomains'))
      .hasValue(this.parsedRootCert.permittedDnsDomains, 'form is prefilled with values from old root');
    await click(SELECTORS.generateRootCancel);
    assert.ok(this.onCancel.calledOnce, 'custom form calls @onCancel passed from parent');
    await click(SELECTORS.oldRadioSelect);
    await click(SELECTORS.rotateRootCancel);
    assert.ok(this.onCancel.calledTwice, 'old root settings form calls @onCancel from parent');

    // validations
    await fillIn(SELECTORS.input('commonName'), '');
    await fillIn(SELECTORS.input('issuerName'), 'default');
    await click(SELECTORS.rotateRootSave);
    assert.dom(SELECTORS.validationError).hasText('There are 2 errors with this form.');
    assert.dom(SELECTORS.input('commonName')).hasClass('has-error-border', 'common name has error border');
    assert.dom(SELECTORS.input('issuerName')).hasClass('has-error-border', 'issuer name has error border');
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
    await click(SELECTORS.rotateRootSave);
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
      await fillIn(SELECTORS.typeField, type);
      await click(SELECTORS.generateRootSave);
    });
  }
  testEndpoint(test, 'internal');
  testEndpoint(test, 'exported');
  testEndpoint(test, 'existing');
  testEndpoint(test, 'kms');
});
