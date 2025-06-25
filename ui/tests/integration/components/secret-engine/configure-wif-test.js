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
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import {
  expectedConfigKeys,
  createConfig,
  configUrl,
  fillInAzureConfig,
  fillInAwsConfig,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { capabilitiesStub } from 'vault/tests/helpers/stubs';
import engineDisplayData from 'vault/helpers/engines-display-data';
import waitForError from 'vault/tests/helpers/wait-for-error';
import AwsConfigForm from 'vault/forms/secrets/aws-config';
import AzureConfigForm from 'vault/forms/secrets/azure-config';
import GcpConfigForm from 'vault/forms/secrets/gcp-config';
import SshConfigForm from 'vault/forms/secrets/ssh-config';

const WIF_ENGINES = ALL_ENGINES.filter((e) => e.isWIF).map((e) => e.type);

module('Integration | Component | SecretEngine::ConfigureWif', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.flashDangerSpy = sinon.spy(this.flashMessages, 'danger');
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.uid = uuidv4();
    // stub capabilities so that by default user can read and update issuer
    this.server.post('/sys/capabilities-self', () => capabilitiesStub('identity/oidc/config', ['sudo']));
    this.getForm = (type, data, options) => {
      const form = {
        aws: AwsConfigForm,
        azure: AzureConfigForm,
        gcp: GcpConfigForm,
        ssh: SshConfigForm,
      }[type];
      return new form(data, options);
    };
    this.renderComponent = () =>
      render(hbs`
        <SecretEngine::ConfigureWif
          @backendPath={{this.id}}
          @displayName={{this.displayName}}
          @type={{this.type}}
          @configForm={{this.form}}
        />
      `);
  });

  module('Create view', function () {
    module('Enterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });

      for (const type of WIF_ENGINES) {
        test(`${type}: it renders default fields`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = engineDisplayData(type).displayName;
          this.form = this.getForm(type, {}, { isNew: true });
          this.type = type;

          await this.renderComponent();

          assert.dom(SES.configureForm).exists(`it lands on the ${type} configuration form`);
          assert.dom(SES.wif.accessType(type)).isChecked(`defaults to showing ${type} access type checked`);
          assert.dom(SES.wif.accessType('wif')).isNotChecked('wif access type is not checked');

          let toggleGroup = GENERAL.button('More options');
          if (type === 'aws') {
            toggleGroup = GENERAL.button('Root config options');
          }
          await click(toggleGroup);

          for (const key of expectedConfigKeys(type, true)) {
            if (key === 'configTtl' || key === 'maxTtl') {
              // because toggle.hbs passes in the name rather than the camelized attr, we have a difference of data-test=attrName vs data-test="Item name" being passed into the data-test selectors. Long-term solution we should match toggle.hbs selectors to formField.hbs selectors syntax
              assert
                .dom(GENERAL.ttl.toggle(key === 'configTtl' ? 'Config TTL' : 'Max TTL'))
                .exists(
                  `${key} shows for ${type} configuration create section when wif is not the access type.`
                );
            } else {
              assert
                .dom(GENERAL.inputByAttr(key))
                .exists(
                  `${key} shows for ${type} configuration create section when wif is not the access type`
                );
            }
          }
          assert
            .dom(GENERAL.inputByAttr('issuer'))
            .doesNotExist(`for ${type}, the issuer does not show when wif is not the access type`);
        });
      }

      for (const type of WIF_ENGINES) {
        test(`${type}: it renders wif fields when user selects wif access type`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = engineDisplayData(type).displayName;
          this.form = this.getForm(type, {}, { isNew: true });
          this.type = type;

          await this.renderComponent();
          await click(SES.wif.accessType('wif'));

          let toggleGroup = GENERAL.button('More options');
          if (type === 'aws') {
            toggleGroup = GENERAL.button('Root config options');
          }
          await click(toggleGroup);

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
      /* This module covers code that is the same for all engines. We run them once against one of the engines.*/
      module('Engine agnostic', function () {
        test('it transitions without sending a config or issuer payload on cancel', async function (assert) {
          assert.expect(3);
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.type = 'azure';
          this.form = this.getForm('azure', {}, { isNew: true });

          this.server.post(configUrl('azure', this.id), () => {
            throw new Error(
              `Request was made to post the config when it should not have been because the user canceled out of the flow.`
            );
          });
          this.server.post('/identity/oidc/config', () => {
            throw new Error(
              `Request was made to save the issuer when it should not have been because the user canceled out of the flow.`
            );
          });

          await this.renderComponent();

          await fillInAzureConfig(true);
          await click(GENERAL.cancelButton);

          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(this.flashSuccessSpy.notCalled, 'No success flash messages called.');

          assert.true(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('it throws an error if the getter isWifPluginConfigured is not defined on the model', async function (assert) {
          const promise = waitForError();
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.type = 'azure';
          // creating a config that exists but will not have the attribute isWifPluginConfigured on it
          this.form = this.getForm('ssh', {});

          await this.renderComponent();
          const err = await promise;
          assert.true(
            err.message.includes(
              `'isWifPluginConfigured' is required to be defined on the config model. Must return a boolean.`
            ),
            'asserts without isWifPluginConfigured'
          );
        });

        test('it allows user to submit the config even if API error occurs on issuer config', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.type = 'aws';
          this.form = this.getForm('aws', {}, { isNew: true });

          this.server.post(configUrl('aws', this.id), () => {
            assert.true(true, 'post request was made to config/root when issuer failed. test should pass.');
          });
          this.server.post('/identity/oidc/config', () => {
            return overrideResponse(400, { errors: ['bad request'] });
          });

          await this.renderComponent();

          await fillInAwsConfig('withWif');
          await click(GENERAL.submitButton);
          await click(SES.wif.issuerWarningSave);

          assert.true(
            this.flashDangerSpy.calledWith('Issuer was not saved: bad request'),
            'Flash message shows that issuer was not saved'
          );
          assert.true(
            this.flashSuccessSpy.calledWith(`Successfully saved ${this.id}'s configuration.`),
            'Flash message shows that root was saved even if issuer was not'
          );
          assert.true(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('it surfaces the API error if config save fails, and prevents the user from transitioning', async function (assert) {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.type = 'aws';
          this.form = this.getForm('aws', {}, { isNew: true });

          this.server.post(configUrl('aws', this.id), () => {
            return overrideResponse(400, { errors: ['bad request'] });
          });
          this.server.post(configUrl('aws-lease', this.id), () => {
            assert.true(
              true,
              'post request was made to config/lease when config/root failed. test should pass.'
            );
          });

          await this.renderComponent();
          // fill in both lease and root endpoints to ensure that both payloads are attempted to be sent
          await fillInAwsConfig('withAccess');
          await fillInAwsConfig('withLease');
          await click(GENERAL.submitButton);
          assert.dom(GENERAL.messageError).exists('API error surfaced to user');
          assert.dom(GENERAL.inlineError).exists('User shown inline error message');
        });
      });

      module('Azure specific', function (hooks) {
        hooks.beforeEach(function () {
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.type = 'azure';
          this.form = this.getForm('azure', {}, { isNew: true });
        });
        test('it clears access type inputs after toggling accessType', async function (assert) {
          await this.renderComponent();
          await fillInAzureConfig();
          await click(SES.wif.accessType('wif'));
          await fillInAzureConfig(true);
          await click(SES.wif.accessType('azure'));

          assert
            .dom(GENERAL.toggleInput('Root password TTL'))
            .isChecked('rootPasswordTtl is not cleared after toggling accessType');
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

        test('it shows the correct access type subtext', async function (assert) {
          await this.renderComponent();

          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText(
              'Choose the way to configure access to Azure. Access can be configured either using Azure account credentials or with the Plugin Workload Identity Federation (WIF).'
            );
        });

        test('it does not show aws specific note', async function (assert) {
          await this.renderComponent();

          assert
            .dom(SES.configureNote('azure'))
            .doesNotExist('Note specific to AWS does not show for Azure secret engine when configuring.');
        });
      });

      module('AWS specific', function (hooks) {
        hooks.beforeEach(function () {
          this.id = `aws-${this.uid}`;
          this.displayName = 'AWS';
          this.type = 'aws';
          this.form = this.getForm('aws', {}, { isNew: true });
        });

        test('it clears access type inputs after toggling accessType', async function (assert) {
          await this.renderComponent();
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

        test('it shows the correct access type subtext', async function (assert) {
          await this.renderComponent();

          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText(
              'Choose the way to configure access to AWS. Access can be configured either using IAM access keys or with the Plugin Workload Identity Federation (WIF).'
            );
        });

        test('it shows validation error if default lease is entered but max lease is not', async function (assert) {
          this.server.post(configUrl('aws-lease', this.id), () => {
            throw new Error(
              `Request was made to post the config/lease when it should not have been because no data was changed.`
            );
          });
          this.server.post(configUrl('aws', this.id), () => {
            throw new Error(
              `Request was made to post the config/root when it should not have been because no data was changed.`
            );
          });

          await this.renderComponent();

          await click(GENERAL.ttl.toggle('Default Lease TTL'));
          await fillIn(GENERAL.ttl.input('Default Lease TTL'), '33');
          await click(GENERAL.submitButton);
          assert
            .dom(GENERAL.inlineError)
            .hasText('Lease TTL and Max Lease TTL are both required if one of them is set.');
          assert.dom(SES.configureForm).exists('remains on the configuration form');
        });

        test('it allows user to submit root config even if API error occurs on config/lease config', async function (assert) {
          this.server.post(configUrl('aws', this.id), () => {
            assert.true(
              true,
              'post request was made to config/root when config/lease failed. test should pass.'
            );
          });
          this.server.post(configUrl('aws-lease', this.id), () => {
            return overrideResponse(400, { errors: ['bad request!!'] });
          });

          await this.renderComponent();
          // fill in both lease and root endpoints to ensure that both payloads are attempted to be sent
          await fillInAwsConfig('withAccess');
          await fillInAwsConfig('withLease');
          await click(GENERAL.submitButton);
          assert.true(
            this.flashDangerSpy.calledWith('Error saving lease configuration: bad request!!'),
            'Flash message shows that lease was not saved.'
          );
          assert.true(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('it transitions without sending a lease, root, or issuer payload on cancel', async function (assert) {
          this.server.post(configUrl('aws', this.id), () => {
            throw new Error(
              `Request was made to post the config/root when it should not have been because the user canceled out of the flow.`
            );
          });
          this.server.post(configUrl('aws-lease', this.id), () => {
            throw new Error(
              `Request was made to post the config/lease when it should not have been because the user canceled out of the flow.`
            );
          });
          this.server.post('/identity/oidc/config', () => {
            throw new Error(
              `Request was made to post the identity/oidc/config when it should not have been because the user canceled out of the flow.`
            );
          });

          await this.renderComponent();
          // fill in both lease and root endpoints to ensure that both payloads are attempted to be sent
          await fillInAwsConfig('withWif');
          await fillInAwsConfig('withLease');
          await click(GENERAL.cancelButton);

          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(this.flashSuccessSpy.notCalled, 'No success flash messages called.');
          assert.true(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('it does show aws specific note', async function (assert) {
          await this.renderComponent();

          assert.dom(SES.configureNote('aws')).exists('Note specific to AWS does show when configuring.');
        });
      });

      module('Issuer field tests', function (hooks) {
        hooks.beforeEach(function () {
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.type = 'azure';
          this.form = this.getForm('azure', {}, { isNew: true });
        });
        test('if issuer API error and user changes issuer value, shows specific warning message', async function (assert) {
          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          await fillIn(GENERAL.inputByAttr('issuer'), 'http://change.me.no.read');
          await click(GENERAL.submitButton);
          assert
            .dom(SES.wif.issuerWarningMessage)
            .hasText(
              `You are updating the global issuer config. This will overwrite Vault's current issuer if it exists and may affect other configurations using this value. Continue?`,
              'modal shows message about overwriting value if it exists'
            );
        });

        test('it shows placeholder issuer, and does not call APIs on canceling out of issuer modal', async function (assert) {
          this.server.post('/identity/oidc/config', () => {
            throw new Error(
              'Request was made to post the identity/oidc/config when it should not have been because user canceled out of the modal.'
            );
          });
          this.server.post(configUrl('azure', this.id), () => {
            throw new Error(
              `Request was made to post the config when it should not have been because the user canceled out of the flow.`
            );
          });

          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          assert
            .dom(GENERAL.inputByAttr('issuer'))
            .hasAttribute('placeholder', 'https://vault-test.com', 'shows issuer placeholder');
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('', 'shows issuer is empty when not passed');
          await fillIn(GENERAL.inputByAttr('issuer'), 'http://bar.foo');
          await click(GENERAL.submitButton);
          assert.dom(SES.wif.issuerWarningMessage).exists('issuer modal exists');
          assert
            .dom(SES.wif.issuerWarningMessage)
            .hasText(
              `You are updating the global issuer config. This will overwrite Vault's current issuer if it exists and may affect other configurations using this value. Continue?`,
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
            assert.true(true, 'post request was made to azure config');
            return {};
          });

          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('', 'issuer defaults to empty string');
          await fillIn(GENERAL.inputByAttr('issuer'), newIssuer);
          await click(GENERAL.submitButton);
          assert.dom(SES.wif.issuerWarningMessage).exists('issuer warning modal exists');

          await click(SES.wif.issuerWarningSave);
          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('shows modal when modifying the issuer, has correct payload, and shows flash message on fail', async function (assert) {
          assert.expect(7);
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

          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('');

          await fillIn(GENERAL.inputByAttr('issuer'), 'http://foo.bar');
          await fillIn(GENERAL.inputByAttr('identityTokenAudience'), 'some-value');
          await click(GENERAL.submitButton);
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
          assert.true(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('it does not clear global issuer when toggling accessType', async function (assert) {
          this.form = this.getForm('azure', { issuer: 'issuer' }, { isNew: true });

          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          assert
            .dom(GENERAL.inputByAttr('issuer'))
            .hasValue(this.form.issuer, 'issuer is what is sent in by the model on first load');
          await fillIn(GENERAL.inputByAttr('issuer'), 'http://ive-changed');
          await click(SES.wif.accessType('azure'));
          await click(SES.wif.accessType('wif'));
          assert
            .dom(GENERAL.inputByAttr('issuer'))
            .hasValue(
              this.form.issuer,
              'issuer value is still the same global value after toggling accessType'
            );
        });
      });
    });

    module('Community', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });

      for (const type of WIF_ENGINES) {
        test(`${type}: it renders fields`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = engineDisplayData(type).displayName;
          this.form = this.getForm(type, {}, { isNew: true });
          this.type = type;

          await this.renderComponent();
          assert.dom(SES.configureForm).exists(`lands on the ${type} configuration form`);
          assert
            .dom(SES.wif.accessTypeSection)
            .doesNotExist('Access type section does not render for a community user');

          let toggleGroup = GENERAL.button('More options');
          if (type === 'aws') {
            toggleGroup = GENERAL.button('Root config options');
          }
          await click(toggleGroup);

          // check all the form fields are present
          for (const key of expectedConfigKeys(type, true)) {
            if (key === 'configTtl' || key === 'maxTtl') {
              assert
                .dom(GENERAL.ttl.toggle(key === 'configTtl' ? 'Config TTL' : 'Max TTL'))
                .exists(`${key} shows for ${type} account access section.`);
            } else {
              assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for ${type} account access section.`);
            }
          }
          assert.dom(GENERAL.inputByAttr('issuer')).doesNotExist();
        });
      }
    });
  });

  module('Edit view', function () {
    module('Enterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });
      for (const type of WIF_ENGINES) {
        test(`${type}: it defaults to WIF accessType if WIF fields are already set`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = engineDisplayData(type).displayName;
          const config = createConfig(`${type}-wif`);
          this.form = this.getForm(type, config);
          this.type = type;

          await this.renderComponent();
          assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
          assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
          assert.dom(SES.wif.accessType(type)).isNotChecked(`${type} accessType is not checked`);
          assert.dom(SES.wif.accessType(type)).isDisabled(`${type} accessType is disabled`);
          assert.dom(GENERAL.inputByAttr('identityTokenAudience')).hasValue(this.form.identityTokenAudience);
          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText('You cannot edit Access Type if you have already saved access credentials.');
          assert.dom(GENERAL.ttl.input('Identity token TTL')).hasValue('2'); // 7200 on payload is 2hrs in ttl picker
        });
      }

      for (const type of WIF_ENGINES) {
        test(`${type}: it renders issuer if global issuer is already set`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          this.displayName = engineDisplayData(type).displayName;
          this.issuer = 'https://foo-bar-blah.com';
          const config = createConfig(`${type}-wif`);
          this.form = this.getForm(type, { ...config, issuer: this.issuer });
          this.type = type;

          await this.renderComponent();

          assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
          assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
          assert
            .dom(GENERAL.inputByAttr('issuer'))
            .hasValue(this.issuer, `it has the global issuer value of ${this.issuer}`);
        });
      }

      module('Azure specific', function (hooks) {
        // "clientSecret" is the only mutually exclusive Azure account attr and it's never returned from the API. Thus, we can only check for the presence of configured wif fields to determine if the accessType should be preselected to wif and disabled.
        hooks.beforeEach(function () {
          this.id = `azure-${this.uid}`;
          this.displayName = 'Azure';
          this.type = 'azure';
          const config = createConfig('azure');
          this.form = this.getForm('azure', config);
        });

        test('it allows you to change access type if no wif fields are set', async function (assert) {
          await this.renderComponent();

          assert.dom(SES.wif.accessType('azure')).isChecked('Azure accessType is checked');
          assert
            .dom(SES.wif.accessType('azure'))
            .isNotDisabled(
              'Azure accessType is not disabled because we cannot determine if client secret was set as it is not returned by the api.'
            );
          assert.dom(SES.wif.accessType('wif')).isNotChecked('WIF accessType is not checked');
          assert.dom(SES.wif.accessType('wif')).isNotDisabled('WIF accessType is disabled');
          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText(
              'Choose the way to configure access to Azure. Access can be configured either using Azure account credentials or with the Plugin Workload Identity Federation (WIF).'
            );
        });

        test('it sets access type to wif if wif fields are set', async function (assert) {
          const config = createConfig('azure-wif');
          this.form = this.getForm('azure', config);

          await this.renderComponent();

          assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
          assert
            .dom(SES.wif.accessType('azure'))
            .isDisabled('Azure accessType IS disabled because wif attributes are set.');

          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText('You cannot edit Access Type if you have already saved access credentials.');
        });

        test('it shows previously saved config information', async function (assert) {
          this.id = `azure-${this.uid}`;
          const config = createConfig('azure-generic');
          this.form = this.getForm('azure', config);

          await this.renderComponent();
          assert.dom(GENERAL.inputByAttr('subscriptionId')).hasValue(this.form.subscriptionId);
          assert.dom(GENERAL.inputByAttr('clientId')).hasValue(this.form.clientId);
          assert.dom(GENERAL.inputByAttr('tenantId')).hasValue(this.form.tenantId);
          assert
            .dom(GENERAL.inputByAttr('clientSecret'))
            .hasValue('**********', 'clientSecret is masked on edit the value');
        });

        test('it requires a double click to change the client secret', async function (assert) {
          this.id = `azure-${this.uid}`;
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
          await click(GENERAL.button('toggle-masked'));
          await fillIn(GENERAL.inputByAttr('clientSecret'), 'new-secret');
          await click(GENERAL.submitButton);
        });
      });

      module('AWS specific', function (hooks) {
        hooks.beforeEach(function () {
          this.id = `aws-${this.uid}`;
          this.type = 'aws';
          this.displayName = 'AWS';
          const config = createConfig('aws');
          this.form = this.getForm('aws', config);
        });

        test('it defaults to IAM accessType if IAM fields are already set', async function (assert) {
          await this.renderComponent();
          assert.dom(SES.wif.accessType('aws')).isChecked('IAM accessType is checked');
          assert.dom(SES.wif.accessType('aws')).isDisabled('IAM accessType is disabled');
          assert.dom(SES.wif.accessType('wif')).isNotChecked('WIF accessType is not checked');
          assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText('You cannot edit Access Type if you have already saved access credentials.');
        });

        test('it allows you to change access type if record does not have wif or iam values already set', async function (assert) {
          const config = createConfig('aws-no-access');
          this.form = this.getForm('aws', config);

          await this.renderComponent();
          assert.dom(SES.wif.accessType('wif')).isNotDisabled('WIF accessType is NOT disabled');
          assert.dom(SES.wif.accessType('aws')).isNotDisabled('IAM accessType is NOT disabled');
        });

        test('it shows previously saved root and lease information', async function (assert) {
          const config = { ...createConfig('aws'), ...createConfig('aws-lease') };
          this.form = this.getForm('aws', config);

          await this.renderComponent();

          assert.dom(GENERAL.inputByAttr('accessKey')).hasValue(this.form.accessKey);
          assert
            .dom(GENERAL.inputByAttr('secretKey'))
            .hasValue('**********', 'secretKey is masked on edit the value');

          await click(GENERAL.button('Root config options'));
          assert.dom(GENERAL.inputByAttr('region')).hasValue(this.form.region);
          assert.dom(GENERAL.inputByAttr('iamEndpoint')).hasValue(this.form.iamEndpoint);
          assert.dom(GENERAL.inputByAttr('stsEndpoint')).hasValue(this.form.stsEndpoint);
          assert.dom(GENERAL.inputByAttr('maxRetries')).hasValue('1');
          // Check lease config values
          assert.dom(GENERAL.ttl.input('Default Lease TTL')).hasValue('50');
          assert.dom(GENERAL.ttl.input('Max Lease TTL')).hasValue('55');
        });

        test('it requires a double click to change the secret key', async function (assert) {
          this.server.post(configUrl('aws', this.id), (schema, req) => {
            const payload = JSON.parse(req.requestBody);
            assert.strictEqual(
              payload.secret_key,
              'new-secret',
              'post request was made to config/root with the updated secret_key.'
            );
          });

          await this.renderComponent();
          await click(GENERAL.enableField('secretKey'));
          await click(GENERAL.button('toggle-masked'));
          await fillIn(GENERAL.inputByAttr('secretKey'), 'new-secret');
          await click(GENERAL.submitButton);
        });
      });

      module('GCP specific', function (hooks) {
        // "credentials" is the only mutually exclusive GCP account attr and it's never returned from the API. Thus, we can only check for the presence of configured wif fields to determine if the accessType should be preselected to wif and disabled.
        // If the user has configured the credentials field, the ui will not know until the user tries to save WIF fields. This is a limitation of the API and surfaced to the user in a descriptive API error.
        // We cover some of this workflow here and error testing in the gcp-configuration acceptance test.
        hooks.beforeEach(function () {
          this.id = `gcp-${this.uid}`;
          const config = createConfig('gcp');
          this.form = this.getForm('gcp', config);
          this.type = 'gcp';
          this.displayName = 'Google Cloud';
        });
        test('it allows you to change access type if no wif fields are set', async function (assert) {
          await this.renderComponent();

          assert.dom(SES.wif.accessType('gcp')).isChecked('GCP accessType is checked');
          assert
            .dom(SES.wif.accessType('gcp'))
            .isNotDisabled(
              'GCP accessType is not disabled because we cannot determine if credentials was set as it is not returned by the api.'
            );
          assert.dom(SES.wif.accessType('wif')).isNotChecked('WIF accessType is not checked');
          assert.dom(SES.wif.accessType('wif')).isNotDisabled('WIF accessType is not disabled');
          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText(
              'Choose the way to configure access to Google Cloud. Access can be configured either using Google Cloud account credentials or with the Plugin Workload Identity Federation (WIF).'
            );
        });

        test('it sets access type to wif if wif fields are set', async function (assert) {
          const config = createConfig('gcp-wif');
          this.form = this.getForm('gcp', config);

          await this.renderComponent();

          assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
          assert
            .dom(SES.wif.accessType('gcp'))
            .isDisabled('GCP accessType IS disabled because wif attributes are set.');

          assert
            .dom(SES.wif.accessTypeSubtext)
            .hasText('You cannot edit Access Type if you have already saved access credentials.');
        });

        test('it shows previously saved config information', async function (assert) {
          const config = createConfig('gcp-generic');
          this.form = this.getForm('gcp', config);
          await this.renderComponent();
          await click(GENERAL.button('More options'));
          assert.dom(GENERAL.ttl.input('Config TTL')).hasValue('100');
          assert.dom(GENERAL.ttl.input('Max TTL')).hasValue('101');
        });
      });
    });

    module('Community', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });
      for (const type of WIF_ENGINES) {
        test(`${type}:it does not show access type but defaults to type "account" fields`, async function (assert) {
          this.id = `${type}-${this.uid}`;
          const config = createConfig(`${type}-generic`);
          this.form = this.getForm(type, config);
          this.displayName = engineDisplayData(type).displayName;
          this.type = type;

          await this.renderComponent();
          assert.dom(SES.wif.accessTypeSection).doesNotExist('Access type section does not render');

          let toggleGroup = GENERAL.button('More options');
          if (type === 'aws') {
            toggleGroup = GENERAL.button('Root config options');
          }
          await click(toggleGroup);

          for (const key of expectedConfigKeys(type, true)) {
            if (key === 'secretKey' || key === 'clientSecret' || key === 'credentials') return; // these keys are not returned by the API
            // same issues noted in wif enterprise tests with how toggle.hbs passes in name vs how formField input passes in attr to data test selector
            if (key === 'configTtl') {
              assert
                .dom(GENERAL.ttl.input('Config TTL'))
                .hasValue('100', `${key} for ${type}: has the expected value set on the config`);
            } else if (key === 'maxTtl') {
              assert
                .dom(GENERAL.ttl.input('Max TTL'))
                .hasValue('101', `${key} for ${type}: has the expected value set on the config`);
            } else if (key === 'rootPasswordTtl') {
              assert
                .dom(GENERAL.ttl.input('Root password TTL'))
                .hasValue('500', `${key} for ${type}: has the expected value set on the config`);
            } else {
              assert
                .dom(GENERAL.inputByAttr(key))
                .hasValue(this.form[key], `${key} for ${type}: has the expected value set on the config`);
            }
          }
        });
      }
    });
  });
});
