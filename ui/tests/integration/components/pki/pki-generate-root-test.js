/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_GENERATE_ROOT } from 'vault/tests/helpers/pki/pki-selectors';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import sinon from 'sinon';

module('Integration | Component | pki-generate-root', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.owner.lookup('service:secretMountPath').update('pki-test');

    this.capabilitiesForStub = sinon
      .stub(this.owner.lookup('service:capabilities'), 'for')
      .resolves({ canCreate: false });

    const api = this.owner.lookup('service:api');
    this.issuersGenerateStub = sinon.stub(api.secrets, 'pkiIssuersGenerateRoot');
    this.generateStub = sinon.stub(api.secrets, 'pkiGenerateRoot');
    this.rotateStub = sinon.stub(api.secrets, 'pkiRotateRoot');

    this.withUrls = true;
    this.canSetUrls = true;
    this.onCancel = sinon.stub();
    this.onComplete = sinon.stub();
    this.onSave = sinon.stub();

    this.renderComponent = () =>
      render(
        hbs`
          <PkiGenerateRoot
            @withUrls={{this.withUrls}}
            @canSetUrls={{this.canSetUrls}}
            @rotateCertData={{this.rotateCertData}}
            @capabilities={{this.capabilities}}
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

  test('it renders with correct sections', async function (assert) {
    this.withUrls = false;
    await this.renderComponent();

    assert.dom('h2').exists({ count: 1 }, 'One H2 title without @urls');
    assert.dom(PKI_GENERATE_ROOT.mainSectionTitle).hasText('Root parameters');
    assert.dom(GENERAL.button('Key parameters')).exists('Key parameters toggle renders');
    assert.dom(GENERAL.button('Subject Alternative Name (SAN) Options')).exists('SAN options toggle renders');
    assert
      .dom(GENERAL.button('Additional subject fields'))
      .exists('Additional subject fields toggle renders');
  });

  test('it shows the appropriate fields under the toggles', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.button('Additional subject fields'));
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText('These fields provide more information about the client to which the certificate belongs.');
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Additional subject fields'))
      .exists({ count: 7 }, '7 form fields under Additional Fields toggle');

    await click(GENERAL.button('Subject Alternative Name (SAN) Options'));
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'SAN fields are an extension that allow you specify additional host names (sites, IP addresses, common names, etc.) to be protected by a single certificate.'
      );
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Subject Alternative Name (SAN) Options'))
      .exists({ count: 6 }, '7 form fields under SANs toggle');

    await click(GENERAL.button('Key parameters'));
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
    await this.renderComponent();

    await click(GENERAL.button('Key parameters'));
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 0 }, '0 form fields under keyParams toggle');

    this.type = 'exported';
    await fillIn(GENERAL.inputByAttr('type'), this.type);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'This certificate type is exported. This means the private key will be returned in the response. Below, you will name the key and define its type and key bits.',
        `has correct description for type=${this.type}`
      );
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 4 }, '4 form fields under keyParams toggle');
    assert.dom(GENERAL.fieldByAttr('key_name')).exists(`Key name field shown when type=${this.type}`);
    assert.dom(GENERAL.fieldByAttr('key_type')).exists(`Key type field shown when type=${this.type}`);
    assert.dom(GENERAL.fieldByAttr('key_bits')).exists(`Key bits field shown when type=${this.type}`);

    this.type = 'internal';
    await fillIn(GENERAL.inputByAttr('type'), this.type);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'This certificate type is internal. This means that the private key will not be returned and cannot be retrieved later. Below, you will name the key and define its type and key bits.',
        `has correct description for type=${this.type}`
      );
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 3 }, '3 form fields under keyParams toggle');
    assert.dom(GENERAL.fieldByAttr('key_name')).exists(`Key name field shown when type=${this.type}`);
    assert.dom(GENERAL.fieldByAttr('key_type')).exists(`Key type field shown when type=${this.type}`);
    assert.dom(GENERAL.fieldByAttr('key_bits')).exists(`Key bits field shown when type=${this.type}`);

    this.set('type', 'existing');
    await fillIn(GENERAL.inputByAttr('type'), this.type);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'You chose to use an existing key. This means that we’ll use the key reference to create the CSR or root. Please provide the reference to the key.',
        `has correct description for type=${this.type}`
      );
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 1 }, '1 form field under keyParams toggle');
    assert.dom(GENERAL.fieldByAttr('key_ref')).exists(`Key reference field shown when type=${this.type}`);

    this.set('type', 'kms');
    await fillIn(GENERAL.inputByAttr('type'), this.type);
    assert
      .dom(PKI_GENERATE_ROOT.toggleGroupDescription)
      .hasText(
        'This certificate type is kms, meaning managed keys will be used. Below, you will name the key and tell Vault where to find it in your KMS or HSM. Learn more about managed keys.',
        `has correct description for type=${this.type}`
      );
    assert
      .dom(PKI_GENERATE_ROOT.groupFields('Key parameters'))
      .exists({ count: 3 }, '3 form fields under keyParams toggle');
    assert.dom(GENERAL.fieldByAttr('key_name')).exists(`Key name field shown when type=${this.type}`);
    assert
      .dom(GENERAL.fieldByAttr('managed_key_name'))
      .exists(`Managed key name field shown when type=${this.type}`);
    assert
      .dom(GENERAL.fieldByAttr('managed_key_id'))
      .exists(`Managed key id field shown when type=${this.type}`);
  });

  test('it shows errors before submit if form is invalid', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.submitButton);
    assert.dom(PKI_GENERATE_ROOT.formInvalidError).exists('Shows overall error form');
    assert.true(this.onSave.notCalled);
  });

  test('it should use correct endpoint based on issuer permissions', async function (assert) {
    assert.expect(6);

    this.rotateCertData = parseCertificate(CERTIFICATES.loadedCert);
    this.rotateStub.rejects(getErrorResponse());

    await this.renderComponent();
    await this.fillInAndSubmit('exported');

    assert.true(this.rotateStub.calledOnce, 'rotate endpoint used when rotating');
    assert.dom(GENERAL.messageError).exists('error message shown');

    this.rotateCertData = undefined;
    this.capabilitiesForStub.resolves({ canCreate: true });
    this.issuersGenerateStub.rejects(getErrorResponse());
    await this.renderComponent();
    await this.fillInAndSubmit('exported');

    assert.true(
      this.capabilitiesForStub.calledWith('pkiIssuersGenerateRoot', {
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

  module('URLs section', function () {
    test('it does not render when no urls passed', async function (assert) {
      this.withUrls = false;
      await this.renderComponent();
      assert.dom(PKI_GENERATE_ROOT.urlsSection).doesNotExist();
    });

    test('it renders when urls model passed', async function (assert) {
      await this.renderComponent();
      assert.dom(PKI_GENERATE_ROOT.urlsSection).exists();
      assert.dom('h2').exists({ count: 2 }, 'two H2 titles are visible on page load');
      assert.dom(PKI_GENERATE_ROOT.urlSectionTitle).hasText('Issuer URLs');
    });
  });
});
