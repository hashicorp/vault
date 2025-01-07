/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { render, click, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { v4 as uuidv4 } from 'uuid';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import {
  expectedConfigKeys,
  createConfig,
  configUrl,
  fillInAzureConfig,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { capabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | SecretEngine/ConfigureAzure', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.flashDangerSpy = sinon.spy(this.flashMessages, 'danger');
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.uid = uuidv4();
    this.id = `azure-${this.uid}`;
    this.config = this.store.createRecord('azure/config');
    this.issuerConfig = createConfig(this.store, this.id, 'issuer');
    this.config.backend = this.id; // Add backend to the configs because it's not on the testing snapshot (would come from url)
    // stub capabilities so that by default user can read and update issuer
    this.server.post('/sys/capabilities-self', () => capabilitiesStub('identity/oidc/config', ['sudo']));

    this.renderComponent = () => {
      return render(hbs`
        <SecretEngine::ConfigureAzure @model={{this.config}} @issuerConfig={{this.issuerConfig}} @backendPath={{this.id}} />
        `);
    };
  });
  module('Create view', function () {
    module('isEnterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });

      test('it renders default fields, showing access type options for enterprise users', async function (assert) {
        await this.renderComponent();
        assert.dom(SES.configureForm).exists('it lands on the Azure configuration form.');
        assert.dom(SES.wif.accessType('azure')).isChecked('defaults to showing Azure access type checked');
        assert.dom(SES.wif.accessType('wif')).isNotChecked('wif access type is not checked');
        // check all the form fields are present
        for (const key of expectedConfigKeys('azure-camelCase')) {
          assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for root section`);
        }
        assert.dom(GENERAL.inputByAttr('issuer')).doesNotExist();
      });

      test('it renders wif fields when user selects wif access type', async function (assert) {
        await this.renderComponent();
        await click(SES.wif.accessType('wif'));
        // check for the wif fields only
        for (const key of expectedConfigKeys('azure-wif-camelCase')) {
          if (key === 'Identity token TTL') {
            assert.dom(GENERAL.ttl.toggle(key)).exists(`${key} shows for wif section.`);
          } else {
            assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for wif section.`);
          }
        }
        assert.dom(GENERAL.inputByAttr('issuer')).exists('issuer shows for wif section.');
      });

      test('it clears wif/azure-account inputs after toggling accessType', async function (assert) {
        await this.renderComponent();
        await fillInAzureConfig('azure');
        await click(SES.wif.accessType('wif'));
        await fillInAzureConfig('withWif');
        await click(SES.wif.accessType('azure'));

        assert
          .dom(GENERAL.toggleInput('Root password TTL'))
          .isNotChecked('rootPasswordTtl is cleared after toggling accessType');
        assert
          .dom(GENERAL.inputByAttr('clientSecret'))
          .hasValue('', 'clientSecret is cleared after toggling accessType');
        await click(SES.wif.accessType('wif'));
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue('', 'issuer shows no value after toggling accessType');
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasAttribute(
            'placeholder',
            'https://vault-test.com',
            'issuer shows no value after toggling accessType'
          );
        assert
          .dom(GENERAL.inputByAttr('identityTokenAudience'))
          .hasValue('', 'idTokenAudience is cleared after toggling accessType');
        assert
          .dom(GENERAL.toggleInput('Identity token TTL'))
          .isNotChecked('identityTokenTtl is cleared after toggling accessType');
      });

      test('it does not clear global issuer when toggling accessType', async function (assert) {
        this.issuerConfig = createConfig(this.store, this.id, 'issuer');
        await this.renderComponent();
        await click(SES.wif.accessType('wif'));
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue(this.issuerConfig.issuer, 'issuer is what is sent in by the model on first load');
        await fillIn(GENERAL.inputByAttr('issuer'), 'http://ive-changed');
        await click(SES.wif.accessType('azure'));
        await click(SES.wif.accessType('wif'));
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue(
            this.issuerConfig.issuer,
            'issuer value is still the same global value after toggling accessType'
          );
      });

      test('it transitions without sending a config or issuer payload on cancel', async function (assert) {
        assert.expect(3);
        await this.renderComponent();
        this.server.post(configUrl('azure', this.id), () => {
          assert.notOk(
            true,
            'post request was made to config when user canceled out of flow. test should fail.'
          );
        });
        this.server.post('/identity/oidc/config', () => {
          assert.notOk(
            true,
            'post request was made to save issuer when user canceled out of flow. test should fail.'
          );
        });
        await fillInAzureConfig('withWif');
        await click(GENERAL.cancelButton);

        assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
        assert.true(this.flashSuccessSpy.notCalled, 'No success flash messages called.');

        assert.ok(
          this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
          'Transitioned to the configuration index route.'
        );
      });

      module('issuer field tests', function () {
        test('if issuer API error and user changes issuer value, shows specific warning message', async function (assert) {
          this.issuerConfig.queryIssuerError = true;
          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          await fillIn(GENERAL.inputByAttr('issuer'), 'http://change.me.no.read');
          await click(GENERAL.saveButton);
          assert
            .dom(SES.wif.issuerWarningMessage)
            .hasText(
              `You are updating the global issuer config. This will overwrite Vault's current issuer if it exists and may affect other configurations using this value. Continue?`,
              'modal shows message about overwriting value if it exists'
            );
        });

        test('is shows placeholder issuer, and does not call APIs on canceling out of issuer modal', async function (assert) {
          this.server.post('/identity/oidc/config', () => {
            assert.notOk(true, 'request should not be made to issuer config endpoint');
          });
          this.server.post(configUrl('azure', this.id), () => {
            assert.notOk(
              true,
              'post request was made to config/ when user canceled out of flow. test should fail.'
            );
          });
          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          assert
            .dom(GENERAL.inputByAttr('issuer'))
            .hasAttribute('placeholder', 'https://vault-test.com', 'shows issuer placeholder');
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('', 'shows issuer is empty when not passed');
          await fillIn(GENERAL.inputByAttr('issuer'), 'http://bar.foo');
          await click(GENERAL.saveButton);
          assert.dom(SES.wif.issuerWarningMessage).exists('issuer modal exists');
          assert
            .dom(SES.wif.issuerWarningMessage)
            .hasText(
              `You are updating the global issuer config. This will overwrite Vault's current issuer and may affect other configurations using this value. Continue?`,
              'modal shows message about overwriting value without the noRead: "if it exists" adage'
            );
          await click(SES.wif.issuerWarningCancel);
          assert.dom(SES.wif.issuerWarningMessage).doesNotExist('issuer modal is removed on cancel');
          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(this.flashSuccessSpy.notCalled, 'No success flash messages called.');
          assert.true(this.transitionStub.notCalled, 'Does not redirect');
        });

        test('it shows modal when updating issuer and calls correct APIs on save', async function (assert) {
          const newIssuer = `http://bar.${uuidv4()}`;
          this.server.post('/identity/oidc/config', (schema, req) => {
            const payload = JSON.parse(req.requestBody);
            assert.deepEqual(payload, { issuer: newIssuer }, 'payload for issuer is correct');
            return {
              id: 'identity-oidc-config', // id needs to match the id on secret-engine-helpers createIssuerConfig
              data: null,
              warnings: [
                'If "issuer" is set explicitly, all tokens must be validated against that address, including those issued by secondary clusters. Setting issuer to "" will restore the default behavior of using the cluster\'s api_addr as the issuer.',
              ],
            };
          });
          this.server.post(configUrl('azure', this.id), () => {
            assert.notOk(true, 'skips request to config because the model was not changed');
          });
          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('', 'issuer defaults to empty string');
          await fillIn(GENERAL.inputByAttr('issuer'), newIssuer);
          await click(GENERAL.saveButton);

          assert.dom(SES.wif.issuerWarningMessage).exists('issue warning modal exists');

          await click(SES.wif.issuerWarningSave);
          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(
            this.flashSuccessSpy.calledWith('Issuer saved successfully'),
            'Success flash message called for Azure issuer'
          );
          assert.true(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('shows modal when modifying the issuer, has correct payload, and shows flash message on fail', async function (assert) {
          assert.expect(7);
          this.issuer = 'http://foo.bar';
          this.server.post(configUrl('azure', this.id), () => {
            assert.true(
              true,
              'post request was made to azure config when unsetting the issuer. test should pass.'
            );
          });
          this.server.post('/identity/oidc/config', (_, req) => {
            const payload = JSON.parse(req.requestBody);
            assert.deepEqual(payload, { issuer: this.issuer }, 'correctly sets the issuer');
            return overrideResponse(403);
          });

          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('');
          await fillIn(GENERAL.inputByAttr('issuer'), this.issuer);
          await fillIn(GENERAL.inputByAttr('identityTokenAudience'), 'some-value');
          await click(GENERAL.saveButton);
          assert.dom(SES.wif.issuerWarningMessage).exists('issuer warning modal exists');
          await click(SES.wif.issuerWarningSave);

          assert.true(
            this.flashDangerSpy.calledWith('Issuer was not saved: permission denied'),
            'shows danger flash for issuer save'
          );
          assert.true(
            this.flashSuccessSpy.calledWith(`Successfully saved ${this.id}'s configuration.`),
            "calls the config flash message not the issuer's"
          );
          assert.ok(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });
      });
    });
    module('isCommunity', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });

      test('it renders fields', async function (assert) {
        assert.expect(8);
        await this.renderComponent();
        assert.dom(SES.configureForm).exists('t lands on the Azure configuration form');
        assert
          .dom(SES.wif.accessTypeSection)
          .doesNotExist('Access type section does not render for a community user');
        // check all the form fields are present
        for (const key of expectedConfigKeys('azure-camelCase')) {
          assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for azure account creds section.`);
        }
        assert.dom(GENERAL.inputByAttr('issuer')).doesNotExist();
      });

      test('it does not send issuer on save', async function (assert) {
        assert.expect(4);
        await this.renderComponent();
        this.server.post(configUrl('azure', this.id), () => {
          assert.true(true, 'post request was made to config. test should pass.');
        });
        this.server.post('/identity/oidc/config', () => {
          throw new Error('post request was incorrectly made to update issuer');
        });
        await fillInAzureConfig('azure');
        await click(GENERAL.saveButton);
        assert.dom(SES.wif.issuerWarningMessage).doesNotExist('modal should not render');
        assert.true(
          this.flashSuccessSpy.calledWith(`Successfully saved ${this.id}'s configuration.`),
          'Flash message shows that config was saved even if issuer was not.'
        );
        assert.ok(
          this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
          'Transitioned to the configuration index route.'
        );
      });
    });
  });

  module('Edit view', function () {
    module('isEnterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });

      test('it defaults to Azure accessType if Azure account fields are already set', async function (assert) {
        this.config = createConfig(this.store, this.id, 'azure');
        await this.renderComponent();
        assert.dom(SES.wif.accessType('azure')).isChecked('Azure accessType is checked');
        assert.dom(SES.wif.accessType('azure')).isDisabled('Azure accessType is disabled');
        assert.dom(SES.wif.accessType('wif')).isNotChecked('WIF accessType is not checked');
        assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
        assert
          .dom(SES.wif.accessTypeSubtext)
          .hasText('You cannot edit Access Type if you have already saved access credentials.');
      });

      test('it defaults to WIF accessType if WIF fields are already set', async function (assert) {
        this.config = createConfig(this.store, this.id, 'azure-wif');
        await this.renderComponent();
        assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
        assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
        assert.dom(SES.wif.accessType('azure')).isNotChecked('azure accessType is not checked');
        assert.dom(SES.wif.accessType('azure')).isDisabled('azure accessType is disabled');
        assert.dom(GENERAL.inputByAttr('identityTokenAudience')).hasValue(this.config.identityTokenAudience);
        assert
          .dom(SES.wif.accessTypeSubtext)
          .hasText('You cannot edit Access Type if you have already saved access credentials.');
        assert.dom(GENERAL.ttl.input('Identity token TTL')).hasValue('2'); // 7200 on payload is 2hrs in ttl picker
      });

      test('it renders issuer if global issuer is already set', async function (assert) {
        this.config = createConfig(this.store, this.id, 'azure-wif');
        this.issuerConfig = createConfig(this.store, this.id, 'issuer');
        this.issuerConfig.issuer = 'https://foo-bar-blah.com';
        await this.renderComponent();
        assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
        assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue(
            this.issuerConfig.issuer,
            `it has the global issuer value of ${this.issuerConfig.issuer}`
          );
      });

      test('it allows you to change accessType if record does not have wif or azure values already set', async function (assert) {
        // the model does not have to be new for a user to see the option to change the access type.
        // the access type is only disabled if the model has values already set for access type fields.
        this.config = createConfig(this.store, this.id, 'azure-generic');
        await this.renderComponent();
        assert.dom(SES.wif.accessType('wif')).isNotDisabled('WIF accessType is NOT disabled');
        assert.dom(SES.wif.accessType('azure')).isNotDisabled('Azure accessType is NOT disabled');
      });

      test('it shows previously saved config information', async function (assert) {
        this.config = createConfig(this.store, this.id, 'azure-generic');
        await this.renderComponent();
        assert.dom(GENERAL.inputByAttr('subscriptionId')).hasValue(this.config.subscriptionId);
        assert.dom(GENERAL.inputByAttr('clientId')).hasValue(this.config.clientId);
        assert.dom(GENERAL.inputByAttr('tenantId')).hasValue(this.config.tenantId);
        assert
          .dom(GENERAL.inputByAttr('clientSecret'))
          .hasValue('**********', 'clientSecret is masked on edit the value');
      });

      test('it requires a double click to change the client secret', async function (assert) {
        this.config = createConfig(this.store, this.id, 'azure');
        await this.renderComponent();

        this.server.post(configUrl('azure', this.id), (schema, req) => {
          const payload = JSON.parse(req.requestBody);
          assert.strictEqual(
            payload.client_secret,
            'new-secret',
            'post request was made to azure/config with the updated client_secret.'
          );
        });

        await click(GENERAL.enableField('clientSecret'));
        await click('[data-test-button="toggle-masked"]');
        await fillIn(GENERAL.inputByAttr('clientSecret'), 'new-secret');
        await click(GENERAL.saveButton);
      });
    });
    module('isCommunity', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });

      test('it does not show access types but defaults to Azure account fields', async function (assert) {
        this.config = createConfig(this.store, this.id, 'azure-generic');
        await this.renderComponent();
        assert.dom(SES.wif.accessTypeSection).doesNotExist('Access type section does not render');
        assert.dom(GENERAL.inputByAttr('clientId')).hasValue(this.config.clientId);
        assert.dom(GENERAL.inputByAttr('subscriptionId')).hasValue(this.config.subscriptionId);
        assert.dom(GENERAL.inputByAttr('tenantId')).hasValue(this.config.tenantId);
      });
    });
  });
});
