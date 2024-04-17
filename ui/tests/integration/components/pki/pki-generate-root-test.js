/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import Sinon from 'sinon';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_GENERATE_ROOT } from 'vault/tests/helpers/pki/pki-selectors';

module('Integration | Component | pki-generate-root', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.urls = this.store.createRecord('pki/config/urls', { id: 'pki-test' });
    this.model = this.store.createRecord('pki/action');
    this.onSave = Sinon.spy();
    this.onCancel = Sinon.spy();
  });

  test('it renders with correct sections', async function (assert) {
    await render(
      hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );

    assert.dom('h2').exists({ count: 1 }, 'One H2 title without @urls');
    assert.dom(PKI_GENERATE_ROOT.mainSectionTitle).hasText('Root parameters');
    assert.dom('[data-test-toggle-group]').exists({ count: 3 }, '3 toggle groups shown');
  });

  test('it shows the appropriate fields under the toggles', async function (assert) {
    await render(
      hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );

    await click(PKI_GENERATE_ROOT.additionalGroupToggle);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText('These fields provide more information about the client to which the certificate belongs.');
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Additional subject fields'))
      .exists({ count: 7 }, '7 form fields under Additional Fields toggle');

    await click(PKI_GENERATE_ROOT.sanGroupToggle);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'SAN fields are an extension that allow you specify additional host names (sites, IP addresses, common names, etc.) to be protected by a single certificate.'
      );
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Subject Alternative Name (SAN) Options'))
      .exists({ count: 6 }, '7 form fields under SANs toggle');

    await click(PKI_GENERATE_ROOT.keyParamsGroupToggle);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'Please choose a type to see key parameter options.',
        'Shows empty state description before type is selected'
      );
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 0 }, '0 form fields under keyParams toggle');
  });

  test('it renders the correct form fields in key params', async function (assert) {
    this.set('type', '');
    await render(
      hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );
    await click(PKI_GENERATE_ROOT.keyParamsGroupToggle);
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 0 }, '0 form fields under keyParams toggle');

    this.set('type', 'exported');
    await fillIn(GENERAL.inputByAttr('type'), this.type);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'This certificate type is exported. This means the private key will be returned in the response. Below, you will name the key and define its type and key bits.',
        `has correct description for type=${this.type}`
      );
    assert.strictEqual(this.model.type, this.type);
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 4 }, '4 form fields under keyParams toggle');
    assert.dom(GENERAL.fieldByAttr('keyName')).exists(`Key name field shown when type=${this.type}`);
    assert.dom(GENERAL.fieldByAttr('keyType')).exists(`Key type field shown when type=${this.type}`);
    assert.dom(GENERAL.fieldByAttr('keyBits')).exists(`Key bits field shown when type=${this.type}`);

    this.set('type', 'internal');
    await fillIn(GENERAL.inputByAttr('type'), this.type);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'This certificate type is internal. This means that the private key will not be returned and cannot be retrieved later. Below, you will name the key and define its type and key bits.',
        `has correct description for type=${this.type}`
      );
    assert.strictEqual(this.model.type, this.type);
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 3 }, '3 form fields under keyParams toggle');
    assert.dom(GENERAL.fieldByAttr('keyName')).exists(`Key name field shown when type=${this.type}`);
    assert.dom(GENERAL.fieldByAttr('keyType')).exists(`Key type field shown when type=${this.type}`);
    assert.dom(GENERAL.fieldByAttr('keyBits')).exists(`Key bits field shown when type=${this.type}`);

    this.set('type', 'existing');
    await fillIn(GENERAL.inputByAttr('type'), this.type);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'You chose to use an existing key. This means that we’ll use the key reference to create the CSR or root. Please provide the reference to the key.',
        `has correct description for type=${this.type}`
      );
    assert.strictEqual(this.model.type, this.type);
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 1 }, '1 form field under keyParams toggle');
    assert.dom(GENERAL.fieldByAttr('keyRef')).exists(`Key reference field shown when type=${this.type}`);

    this.set('type', 'kms');
    await fillIn(GENERAL.inputByAttr('type'), this.type);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'This certificate type is kms, meaning managed keys will be used. Below, you will name the key and tell Vault where to find it in your KMS or HSM. Learn more about managed keys.',
        `has correct description for type=${this.type}`
      );
    assert.strictEqual(this.model.type, this.type);
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 3 }, '3 form fields under keyParams toggle');
    assert.dom(GENERAL.fieldByAttr('keyName')).exists(`Key name field shown when type=${this.type}`);
    assert
      .dom(GENERAL.fieldByAttr('managedKeyName'))
      .exists(`Managed key name field shown when type=${this.type}`);
    assert
      .dom(GENERAL.fieldByAttr('managedKeyId'))
      .exists(`Managed key id field shown when type=${this.type}`);
  });

  test('it shows errors before submit if form is invalid', async function (assert) {
    const saveSpy = Sinon.spy();
    this.set('onSave', saveSpy);
    await render(
      hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );

    await click(GENERAL.saveButton);
    assert.dom(PKI_GENERATE_ROOT.formInvalidError).exists('Shows overall error form');
    assert.ok(saveSpy.notCalled);
  });

  module('URLs section', function () {
    test('it does not render when no urls passed', async function (assert) {
      await render(
        hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
        {
          owner: this.engine,
        }
      );

      assert.dom(PKI_GENERATE_ROOT.urlsSection).doesNotExist();
    });

    test('it renders when urls model passed', async function (assert) {
      await render(
        hbs`<PkiGenerateRoot @model={{this.model}} @urls={{this.urls}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
        {
          owner: this.engine,
        }
      );
      assert.dom(PKI_GENERATE_ROOT.urlsSection).exists();
      assert.dom('h2').exists({ count: 2 }, 'two H2 titles are visible on page load');
      assert.dom(PKI_GENERATE_ROOT.urlSectionTitle).hasText('Issuer URLs');
    });
  });
});
