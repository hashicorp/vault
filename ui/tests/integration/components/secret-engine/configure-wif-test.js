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
import { v4 as uuidv4 } from 'uuid';
import { hbs } from 'ember-cli-htmlbars';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import {
  expectedConfigKeys,
  createConfig,
  configUrl,
  fillInAzureConfig,
  fillInAwsConfig,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { capabilitiesStub } from 'vault/tests/helpers/stubs';
import { WIF_ENGINES, allEngines } from 'vault/helpers/mountable-secret-engines';
import waitForError from 'vault/tests/helpers/wait-for-error';

const allEnginesArray = allEngines(); // saving as const so we don't invoke the method multiple times via the for loop

module('Integration | Component | SecretEngine/ConfigureWif', function (hooks) {
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
    // stub capabilities so that by default user can read and update issuer
    this.server.post('/sys/capabilities-self', () => capabilitiesStub('identity/oidc/config', ['sudo']));
  });

  module('Create view', function () {
    module('isEnterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });
      for (const type of WIF_ENGINES) {
        test(`${type}: it renders default fields`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = allEnginesArray.find((engine) => engine.type === type)?.displayName;
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config =
            type === 'aws'
              ? this.store.createRecord('aws/root-config')
              : this.store.createRecord(`${type}/config`);
          this.secondConfig = type === 'aws' ? this.store.createRecord('aws/lease-config') : null;
          this.config.backend = this.id;
          this.secondConfig ? (this.secondConfig.backend = this.id) : null; // Add backend to the configs because it's not on the testing snapshot (would come from url)
          this.type = type; // required to set on this to pass into the component

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type={{this.type}} @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          assert.dom(SES.configureForm).exists(`it lands on the ${type} configuration form`);
          assert.dom(SES.wif.accessType(type)).isChecked(`defaults to showing ${type} access type checked`);
          assert.dom(SES.wif.accessType('wif')).isNotChecked('wif access type is not checked');
          // toggle grouped fields if it exists
          const toggleGroup = document.querySelector('[data-test-toggle-group]');
          toggleGroup ? await click(toggleGroup) : null;

          for (const key of expectedConfigKeys(type, true)) {
            assert
              .dom(GENERAL.inputByAttr(key))
              .exists(
                `${key} shows for ${type} configuration create section when wif is not the access type`
              );
          }
          assert
            .dom(GENERAL.inputByAttr('issuer'))
            .doesNotExist(`for ${type}, the issuer does not show when wif is not the access type`);
        });
      }
      for (const type of WIF_ENGINES) {
        test(`${type}: it renders wif fields when user selects wif access type`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = allEnginesArray.find((engine) => engine.type === type)?.displayName;
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config =
            type === 'aws'
              ? this.store.createRecord('aws/root-config')
              : this.store.createRecord(`${type}/config`);
          this.secondConfig = type === 'aws' ? this.store.createRecord('aws/lease-config') : null;
          this.config.backend = this.id;
          this.secondConfig ? (this.secondConfig.backend = this.id) : null; // Add backend to the configs because it's not on the testing snapshot (would come from url)
          this.type = type; // required to set on this to pass into the component

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type={{this.type}} @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          await click(SES.wif.accessType('wif'));
          // check for the wif fields only
          for (const key of expectedConfigKeys(`${type}-wif`, true)) {
            if (key === 'Identity token TTL') {
              assert.dom(GENERAL.ttl.toggle(key)).exists(`${key} shows for ${type} wif section.`);
            } else {
              assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for ${type} wif section.`);
            }
          }
          assert.dom(GENERAL.inputByAttr('issuer')).exists(`issuer shows for ${type} wif section.`);
        });
      }
      /* For tests relevant to all engines, they are run once within one of the engine specific modules below */
      module('Azure specific', function () {
        test('it clears access type inputs after toggling accessType', async function (assert) {
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('azure/config');
          this.config.backend = this.id;

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
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

        test('it transitions without sending a config or issuer payload on cancel', async function (assert) {
          assert.expect(3);
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('azure/config');
          this.config.backend = this.id;

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
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

        test('it throws an error if the getter isWifPluginConfigured is not defined on the model', async function (assert) {
          const promise = waitForError();
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          // creating a config that exists but will not have the attribute isWifPluginConfigured on it
          this.config = this.store.createRecord('ssh/ca-config', { backend: this.id });
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
          const err = await promise;
          assert.ok(
            err.message.includes(
              `'isWifPluginConfigured' is required to be defined on the config model. Must return a boolean.`
            ),
            'asserts without isWifPluginConfigured'
          );
        });
      });

      module('AWS specific', function () {
        test('it clears access type inputs after toggling accessType', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('aws/root-config');
          this.config.backend = this.id;

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
          await fillInAwsConfig('aws');
          await click(SES.wif.accessType('wif'));
          await fillInAwsConfig('with-wif');
          await click(SES.wif.accessType('aws'));

          assert
            .dom(GENERAL.inputByAttr('accessKey'))
            .hasValue('', 'accessKey is cleared after toggling accessType');

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

        test('it shows validation error if default lease is entered but max lease is not', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('aws/root-config');
          this.secondConfig = this.store.createRecord('aws/lease-config');
          this.config.backend = this.secondConfig.backend = this.id; // Add backend to the configs because it's not on the testing snapshot (would come from url)

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          this.server.post(configUrl('aws-lease', this.id), () => {
            assert.false(
              true,
              'post request was made to config/lease when no data was changed. test should fail.'
            );
          });
          this.server.post(configUrl('aws', this.id), () => {
            assert.false(
              true,
              'post request was made to config/root when no data was changed. test should fail.'
            );
          });
          await click(GENERAL.ttl.toggle('Default Lease TTL'));
          await fillIn(GENERAL.ttl.input('Default Lease TTL'), '33');
          await click(GENERAL.saveButton);
          assert
            .dom(GENERAL.inlineError)
            .hasText('Lease TTL and Max Lease TTL are both required if one of them is set.');
          assert.dom(SES.configureForm).exists('remains on the configuration form');
        });

        test('it surfaces the API error if one occurs on root/config, preventing user from transitioning', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('aws/root-config');
          this.secondConfig = this.store.createRecord('aws/lease-config');
          this.config.backend = this.secondConfig.backend = this.id; // Add backend to the configs because it's not on the testing snapshot (would come from url)

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          this.server.post(configUrl('aws', this.id), () => {
            return overrideResponse(400, { errors: ['bad request'] });
          });
          this.server.post(configUrl('aws-lease', this.id), () => {
            assert.true(
              true,
              'post request was made to config/lease when config/root failed. test should pass.'
            );
          });
          // fill in both lease and root endpoints to ensure that both payloads are attempted to be sent
          await fillInAwsConfig('withAccess');
          await fillInAwsConfig('withLease');
          await click(GENERAL.saveButton);
          assert.dom(GENERAL.messageError).exists('API error surfaced to user');
          assert.dom(GENERAL.inlineError).exists('User shown inline error message');
        });

        test('it allows user to submit root config even if API error occurs on config/lease config', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('aws/root-config');
          this.secondConfig = this.store.createRecord('aws/lease-config');
          this.config.backend = this.secondConfig.backend = this.id; // Add backend to the configs because it's not on the testing snapshot (would come from url)

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          this.server.post(configUrl('aws', this.id), () => {
            assert.true(
              true,
              'post request was made to config/root when config/lease failed. test should pass.'
            );
          });
          this.server.post(configUrl('aws-lease', this.id), () => {
            return overrideResponse(400, { errors: ['bad request!!'] });
          });
          // fill in both lease and root endpoints to ensure that both payloads are attempted to be sent
          await fillInAwsConfig('withAccess');
          await fillInAwsConfig('withLease');
          await click(GENERAL.saveButton);
          assert.true(
            this.flashDangerSpy.calledWith('Lease configuration was not saved: bad request!!'),
            'Flash message shows that lease was not saved.'
          );
          assert.ok(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('it allows user to submit root config even if API error occurs on issuer config', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('aws/root-config');
          this.secondConfig = this.store.createRecord('aws/lease-config');
          this.config.backend = this.secondConfig.backend = this.id; // Add backend to the configs because it's not on the testing snapshot (would come from url)

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          this.server.post(configUrl('aws', this.id), () => {
            assert.true(true, 'post request was made to config/root when issuer failed. test should pass.');
          });
          this.server.post('/identity/oidc/config', () => {
            return overrideResponse(400, { errors: ['bad request'] });
          });
          await fillInAwsConfig('withWif');
          await click(GENERAL.saveButton);
          await click(SES.wif.issuerWarningSave);

          assert.true(
            this.flashDangerSpy.calledWith('Issuer was not saved: bad request'),
            'Flash message shows that issuer was not saved'
          );
          assert.true(
            this.flashSuccessSpy.calledWith(`Successfully saved ${this.id}'s configuration.`),
            'Flash message shows that root was saved even if issuer was not'
          );
          assert.ok(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('it transitions without sending a lease, root, or issuer payload on cancel', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('aws/root-config');
          this.secondConfig = this.store.createRecord('aws/lease-config');
          this.config.backend = this.secondConfig.backend = this.id; // Add backend to the configs because it's not on the testing snapshot (would come from url)

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          this.server.post(configUrl('aws', this.id), () => {
            assert.true(
              false,
              'post request was made to config/root when user canceled out of flow. test should fail.'
            );
          });
          this.server.post(configUrl('aws-lease', this.id), () => {
            assert.true(
              false,
              'post request was made to config/lease when user canceled out of flow. test should fail.'
            );
          });
          this.server.post('/identity/oidc/config', () => {
            assert.true(
              false,
              'post request was made to save issuer when user canceled out of flow. test should fail.'
            );
          });
          // fill in both lease and root endpoints to ensure that both payloads are attempted to be sent
          await fillInAwsConfig('withWif');
          await fillInAwsConfig('withLease');
          await click(GENERAL.cancelButton);

          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(this.flashSuccessSpy.notCalled, 'No success flash messages called.');
          assert.ok(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });
      });

      module('Issuer field tests', function () {
        // only need to test one engine because code is the same for all engines.
        test('if issuer API error and user changes issuer value, shows specific warning message', async function (assert) {
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.issuerConfig.queryIssuerError = true;
          this.config = this.store.createRecord('azure/config');
          this.config.backend = this.id;

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
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

        test('it shows placeholder issuer, and does not call APIs on canceling out of issuer modal', async function (assert) {
          this.server.post('/identity/oidc/config', () => {
            assert.notOk(true, 'request should not be made to issuer config endpoint');
          });
          this.server.post(configUrl('azure', this.id), () => {
            assert.notOk(
              true,
              'post request was made to config/ when user canceled out of flow. test should fail.'
            );
          });
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('azure/config');
          this.config.backend = this.id;

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
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
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('azure/config');
          this.config.backend = this.id;

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
          await click(SES.wif.accessType('wif'));
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('', 'issuer defaults to empty string');
          await fillIn(GENERAL.inputByAttr('issuer'), newIssuer);
          await click(GENERAL.saveButton);
          assert.dom(SES.wif.issuerWarningMessage).exists('issuer warning modal exists');

          await click(SES.wif.issuerWarningSave);
          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(
            this.flashSuccessSpy.calledWith('Issuer saved successfully'),
            'Success flash message called for issuer'
          );
          assert.true(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('shows modal when modifying the issuer, has correct payload, and shows flash message on fail', async function (assert) {
          assert.expect(7);
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord('azure/config');
          this.config.backend = this.id;

          this.server.post(configUrl('azure', this.id), () => {
            assert.true(
              true,
              'post request was made to azure config when unsetting the issuer. test should pass.'
            );
          });
          this.server.post('/identity/oidc/config', (_, req) => {
            const payload = JSON.parse(req.requestBody);
            assert.deepEqual(payload, { issuer: 'http://foo.bar' }, 'correctly sets the issuer');
            return overrideResponse(403);
          });

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
          await click(SES.wif.accessType('wif'));
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('');

          await fillIn(GENERAL.inputByAttr('issuer'), 'http://foo.bar');
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

        test('it does not clear global issuer when toggling accessType', async function (assert) {
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config = this.store.createRecord(`azure/config`);
          this.config.backend = this.id; // Add backend to the configs because it's not on the testing snapshot (would come from url)
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
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
      });
    });

    module('isCommunity', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });

      for (const type of WIF_ENGINES) {
        test(`${type}: it renders fields`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = allEnginesArray.find((engine) => engine.type === type)?.displayName;
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.config =
            type === 'aws'
              ? this.store.createRecord('aws/root-config')
              : type === 'ssh'
              ? this.store.createRecord('ssh/ca-config')
              : this.store.createRecord(`${type}/config`);
          this.secondConfig = type === 'aws' ? this.store.createRecord('aws/lease-config') : null;
          this.config.backend = this.id;
          this.secondConfig ? (this.secondConfig.backend = this.id) : null;
          this.type = type; // required to set on this to pass into the component

          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type={{this.type}} @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          assert.dom(SES.configureForm).exists(`lands on the ${type} configuration form`);
          assert
            .dom(SES.wif.accessTypeSection)
            .doesNotExist('Access type section does not render for a community user');
          // toggle grouped fields if it exists
          const toggleGroup = document.querySelector('[data-test-toggle-group]');
          toggleGroup ? await click(toggleGroup) : null;
          // check all the form fields are present
          for (const key of expectedConfigKeys(type, true)) {
            assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for ${type} account access section.`);
          }
          assert.dom(GENERAL.inputByAttr('issuer')).doesNotExist();
        });
      }
    });
  });

  module('Edit view', function () {
    module('isEnterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });
      for (const type of WIF_ENGINES) {
        test(`${type}: it defaults to WIF accessType if WIF fields are already set`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = allEnginesArray.find((engine) => engine.type === type)?.displayName;
          this.config = createConfig(this.store, this.id, `${type}-wif`);
          this.type = type;
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type={{this.type}} @model={{this.config}} @issuerConfig={{this.issuerConfig}} />
              `);
          assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
          assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
          assert.dom(SES.wif.accessType(type)).isNotChecked(`${type} accessType is not checked`);
          assert.dom(SES.wif.accessType(type)).isDisabled(`${type} accessType is disabled`);
          assert
            .dom(GENERAL.inputByAttr('identityTokenAudience'))
            .hasValue(this.config.identityTokenAudience);
          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText('You cannot edit Access Type if you have already saved access credentials.');
          assert.dom(GENERAL.ttl.input('Identity token TTL')).hasValue('2'); // 7200 on payload is 2hrs in ttl picker
        });
      }

      for (const type of WIF_ENGINES) {
        test(`${type}: it renders issuer if global issuer is already set`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = allEnginesArray.find((engine) => engine.type === type)?.displayName;
          this.config = createConfig(this.store, this.id, `${type}-wif`);
          this.issuerConfig = createConfig(this.store, this.id, 'issuer');
          this.issuerConfig.issuer = 'https://foo-bar-blah.com';
          this.type = type;
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type={{this.type}} @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);

          assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
          assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
          assert
            .dom(GENERAL.inputByAttr('issuer'))
            .hasValue(
              this.issuerConfig.issuer,
              `it has the global issuer value of ${this.issuerConfig.issuer}`
            );
        });
      }

      module('Azure specific', function () {
        test('it defaults to Azure accessType if Azure account fields are already set', async function (assert) {
          this.id = `azure-${this.uid}`;
          this.config = createConfig(this.store, this.id, 'azure');
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName='Azure' @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);

          assert.dom(SES.wif.accessType('azure')).isChecked('Azure accessType is checked');
          assert.dom(SES.wif.accessType('azure')).isDisabled('Azure accessType is disabled');
          assert.dom(SES.wif.accessType('wif')).isNotChecked('WIF accessType is not checked');
          assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText('You cannot edit Access Type if you have already saved access credentials.');
        });

        test('it allows you to change accessType if record does not have wif or azure values already set', async function (assert) {
          this.id = `azure-${this.uid}`;
          this.config = createConfig(this.store, this.id, 'azure-generic');
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName='Azure' @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);

          assert.dom(SES.wif.accessType('wif')).isNotDisabled('WIF accessType is NOT disabled');
          assert.dom(SES.wif.accessType('azure')).isNotDisabled('Azure accessType is NOT disabled');
        });

        test('it shows previously saved config information', async function (assert) {
          this.id = `azure-${this.uid}`;
          this.config = createConfig(this.store, this.id, 'azure-generic');
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName='Azure' @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
          assert.dom(GENERAL.inputByAttr('subscriptionId')).hasValue(this.config.subscriptionId);
          assert.dom(GENERAL.inputByAttr('clientId')).hasValue(this.config.clientId);
          assert.dom(GENERAL.inputByAttr('tenantId')).hasValue(this.config.tenantId);
          assert
            .dom(GENERAL.inputByAttr('clientSecret'))
            .hasValue('**********', 'clientSecret is masked on edit the value');
        });

        test('it requires a double click to change the client secret', async function (assert) {
          this.id = `azure-${this.uid}`;
          this.config = createConfig(this.store, this.id, 'azure');
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName='Azure' @type='azure' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);

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

      module('AWS specific', function () {
        test('it defaults to IAM accessType if IAM fields are already set', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.config = createConfig(this.store, this.id, 'aws');
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName='AWS' @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
          assert.dom(SES.wif.accessType('aws')).isChecked('IAM accessType is checked');
          assert.dom(SES.wif.accessType('aws')).isDisabled('IAM accessType is disabled');
          assert.dom(SES.wif.accessType('wif')).isNotChecked('WIF accessType is not checked');
          assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText('You cannot edit Access Type if you have already saved access credentials.');
        });

        test('it allows you to change access type if record does not have wif or iam values already set', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.config = createConfig(this.store, this.id, 'aws-no-access');
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName='AWS' @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}}/>
              `);
          assert.dom(SES.wif.accessType('wif')).isNotDisabled('WIF accessType is NOT disabled');
          assert.dom(SES.wif.accessType('aws')).isNotDisabled('IAM accessType is NOT disabled');
        });

        test('it shows previously saved root and lease information', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.config = createConfig(this.store, this.id, 'aws');
          this.secondConfig = createConfig(this.store, this.id, 'aws-lease');
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName='AWS' @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);

          assert.dom(GENERAL.inputByAttr('accessKey')).hasValue(this.config.accessKey);
          assert
            .dom(GENERAL.inputByAttr('secretKey'))
            .hasValue('**********', 'secretKey is masked on edit the value');

          await click(GENERAL.toggleGroup('Root config options'));
          assert.dom(GENERAL.inputByAttr('region')).hasValue(this.config.region);
          assert.dom(GENERAL.inputByAttr('iamEndpoint')).hasValue(this.config.iamEndpoint);
          assert.dom(GENERAL.inputByAttr('stsEndpoint')).hasValue(this.config.stsEndpoint);
          assert.dom(GENERAL.inputByAttr('maxRetries')).hasValue('1');
          // Check lease config values
          assert.dom(GENERAL.ttl.input('Default Lease TTL')).hasValue('50');
          assert.dom(GENERAL.ttl.input('Max Lease TTL')).hasValue('55');
        });

        test('it requires a double click to change the secret key', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.config = createConfig(this.store, this.id, 'aws');
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName='AWS' @type='aws' @model={{this.config}} @issuerConfig={{this.issuerConfig}} />
              `);

          this.server.post(configUrl('aws', this.id), (schema, req) => {
            const payload = JSON.parse(req.requestBody);
            assert.strictEqual(
              payload.secret_key,
              'new-secret',
              'post request was made to config/root with the updated secret_key.'
            );
          });

          await click(GENERAL.enableField('secretKey'));
          await click('[data-test-button="toggle-masked"]');
          await fillIn(GENERAL.inputByAttr('secretKey'), 'new-secret');
          await click(GENERAL.saveButton);
        });
      });
    });

    module('isCommunity', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });
      for (const type of WIF_ENGINES) {
        test(`${type}:it does not show access types but defaults to type account fields`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.config = createConfig(this.store, this.id, `${type}-generic`);
          this.displayName = allEnginesArray.find((engine) => engine.type === type)?.displayName;
          this.type = type;
          await render(hbs`
                <SecretEngine::ConfigureWif @backendPath={{this.id}} @displayName={{this.displayName}} @type={{this.type}} @model={{this.config}} @issuerConfig={{this.issuerConfig}} @secondModel={{this.secondConfig}}/>
              `);
          assert.dom(SES.wif.accessTypeSection).doesNotExist('Access type section does not render');
          // toggle grouped fields if it exists
          const toggleGroup = document.querySelector('[data-test-toggle-group]');
          toggleGroup ? await click(toggleGroup) : null;

          for (const key of expectedConfigKeys(type, true)) {
            if (key === 'secretKey' || key === 'clientSecret') return; // these keys are not returned by the API
            assert
              .dom(GENERAL.inputByAttr(key))
              .hasValue(this.config[key], `${key} for ${type}: has the expected value set on the config`);
          }
        });
      }
    });
  });
});
