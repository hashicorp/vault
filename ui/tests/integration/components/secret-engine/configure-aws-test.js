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
  fillInAwsConfig,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { capabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | SecretEngine/ConfigureAws', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.flashDangerSpy = sinon.spy(this.flashMessages, 'danger');
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.uid = uuidv4();
    this.id = `aws-${this.uid}`;
    // using createRecord on root and lease configs to simulate a fresh mount
    this.rootConfig = this.store.createRecord('aws/root-config');
    this.leaseConfig = this.store.createRecord('aws/lease-config');
    // issuer config is never a createdRecord but the response from the API.
    this.issuerConfig = createConfig(this.store, this.id, 'issuer');
    // Add backend to the configs because it's not on the testing snapshot (would come from url)
    this.rootConfig.backend = this.leaseConfig.backend = this.id;
    this.version = this.owner.lookup('service:version');
    // stub capabilities so that by default user can read and update issuer
    this.server.post('/sys/capabilities-self', () => capabilitiesStub('identity/oidc/config', ['sudo']));

    this.renderComponent = () => {
      return render(hbs`
        <SecretEngine::ConfigureAws @rootConfig={{this.rootConfig}} @leaseConfig={{this.leaseConfig}} @issuerConfig={{this.issuerConfig}} @backendPath={{this.id}} />
        `);
    };
  });
  module('Create view', function () {
    module('isEnterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });

      test('it renders fields ', async function (assert) {
        await this.renderComponent();
        assert.dom(SES.aws.rootForm).exists('it lands on the aws root configuration form.');
        assert.dom(SES.wif.accessTitle).exists('Access section is rendered');
        assert.dom(SES.aws.leaseTitle).exists('Lease section is rendered');
        assert.dom(SES.wif.accessTypeSection).exists('Access type section is rendered');
        assert.dom(SES.wif.accessType('iam')).isChecked('defaults to showing IAM access type checked');
        assert.dom(SES.wif.accessType('wif')).isNotChecked('wif access type is not checked');
        // check all the form fields are present
        await click(GENERAL.toggleGroup('Root config options'));
        for (const key of expectedConfigKeys('aws-root-create')) {
          assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for root section.`);
        }
        for (const key of expectedConfigKeys('aws-lease')) {
          assert.dom(`[data-test-ttl-form-label="${key}"]`).exists(`${key} shows for Lease section.`);
        }
        assert.dom(GENERAL.inputByAttr('issuer')).doesNotExist();
      });

      test('it renders wif fields when selected', async function (assert) {
        await this.renderComponent();
        await click(SES.wif.accessType('wif'));
        // check for the wif fields only
        for (const key of expectedConfigKeys('aws-root-create-wif')) {
          if (key === 'Identity token TTL') {
            assert.dom(GENERAL.ttl.toggle(key)).exists(`${key} shows for wif section.`);
          } else {
            assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for wif section.`);
          }
        }
        // check iam fields do not show
        for (const key of expectedConfigKeys('aws-root-create-iam')) {
          assert.dom(GENERAL.inputByAttr(key)).doesNotExist(`${key} does not show when wif is selected.`);
        }
      });

      test('it clears wif/iam inputs after toggling accessType', async function (assert) {
        await this.renderComponent();
        await fillInAwsConfig('withAccess');
        await fillInAwsConfig('withLease');
        await click(SES.wif.accessType('wif')); // toggle to wif
        await fillInAwsConfig('withWif');
        await click(SES.wif.accessType('iam')); // toggle to wif
        assert
          .dom(GENERAL.inputByAttr('accessKey'))
          .hasValue('', 'accessKey is cleared after toggling accessType');
        assert
          .dom(GENERAL.inputByAttr('secretKey'))
          .hasValue('', 'secretKey is cleared after toggling accessType');

        await click(SES.wif.accessType('wif'));
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue('', 'issue shows no value after toggling accessType');
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasAttribute(
            'placeholder',
            'https://vault-test.com',
            'issue shows no value after toggling accessType'
          );
        assert
          .dom(GENERAL.inputByAttr('roleArn'))
          .hasValue('', 'roleArn is cleared after toggling accessType');
        assert
          .dom(GENERAL.inputByAttr('identityTokenAudience'))
          .hasValue('', 'identityTokenAudience is cleared after toggling accessType');
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
          .hasValue(this.issuerConfig.issuer, 'issuer is what is sent in my the model on first load');
        await fillIn(GENERAL.inputByAttr('issuer'), 'http://ive-changed');
        await click(SES.wif.accessType('iam'));
        await click(SES.wif.accessType('wif'));
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue(
            this.issuerConfig.issuer,
            'issuer value is still the same global value after toggling accessType'
          );
      });

      test('it shows validation error if default lease is entered but max lease is not', async function (assert) {
        assert.expect(2);
        await this.renderComponent();
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
        assert.dom(SES.aws.rootForm).exists('remains on the configuration form');
      });

      test('it surfaces the API error if one occurs on root/config, preventing user from transitioning', async function (assert) {
        assert.expect(3);
        await this.renderComponent();
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
        assert.expect(3);
        await this.renderComponent();
        this.server.post(configUrl('aws', this.id), () => {
          assert.true(
            true,
            'post request was made to config/root when config/lease failed. test should pass.'
          );
        });
        this.server.post(configUrl('aws-lease', this.id), () => {
          return overrideResponse(400, { errors: ['bad request'] });
        });
        // fill in both lease and root endpoints to ensure that both payloads are attempted to be sent
        await fillInAwsConfig('withAccess');
        await fillInAwsConfig('withLease');
        await click(GENERAL.saveButton);

        assert.true(
          this.flashDangerSpy.calledWith('Lease configuration was not saved: bad request'),
          'Flash message shows that lease was not saved.'
        );
        assert.ok(
          this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
          'Transitioned to the configuration index route.'
        );
      });

      test('it allows user to submit root config even if API error occurs on issuer config', async function (assert) {
        assert.expect(4);
        await this.renderComponent();
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
          this.flashSuccessSpy.calledWith(`Successfully saved ${this.id}'s root configuration.`),
          'Flash message shows that root was saved even if issuer was not'
        );
        assert.ok(
          this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
          'Transitioned to the configuration index route.'
        );
      });

      test('it transitions without sending a lease, root, or issuer payload on cancel', async function (assert) {
        assert.expect(3);
        await this.renderComponent();
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

      module('issuer field tests', function () {
        // the other tests where issuer is not passed do not show modals, so we only need to test when the modal should shows up
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

        test('is shows placeholder issuer, shows modal when saving changes, and does not call APIs on cancel', async function (assert) {
          this.server.post('/identity/oidc/config', () => {
            assert.notOk(true, 'request should not be made to issuer config endpoint');
          });
          this.server.post(configUrl('aws', this.id), () => {
            assert.notOk(
              true,
              'post request was made to config/root when user canceled out of flow. test should fail.'
            );
          });
          this.server.post(configUrl('aws-lease', this.id), () => {
            assert.notOk(
              true,
              'post request was made to config/lease when user canceled out of flow. test should fail.'
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
          assert.dom(SES.wif.issuerWarningModal).exists('issuer modal exists');
          assert
            .dom(SES.wif.issuerWarningMessage)
            .hasText(
              `You are updating the global issuer config. This will overwrite Vault's current issuer and may affect other configurations using this value. Continue?`,
              'modal shows message about overwriting value without the noRead: "if it exists" adage'
            );
          await click(SES.wif.issuerWarningCancel);
          assert.dom(SES.wif.issuerWarningModal).doesNotExist('issuer modal is removed on cancel');
          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(this.flashSuccessSpy.notCalled, 'No success flash messages called.');
          assert.true(this.transitionStub.notCalled, 'Does not redirect');
        });

        test('it shows modal when updating issuer and calls correct APIs on save', async function (assert) {
          const newIssuer = 'http://bar.foo';
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
          this.server.post(configUrl('aws', this.id), () => {
            assert.notOk(true, 'skips request to config/root due to no changes');
          });
          this.server.post(configUrl('aws-lease', this.id), () => {
            assert.notOk(true, 'skips request to config/lease due to no changes');
          });
          await this.renderComponent();
          await click(SES.wif.accessType('wif'));
          assert.dom(GENERAL.inputByAttr('issuer')).hasValue('', 'issuer defaults to empty string');
          await fillIn(GENERAL.inputByAttr('issuer'), newIssuer);
          await click(GENERAL.saveButton);
          assert.dom(SES.wif.issuerWarningModal).exists('issue warning modal exists');
          await click(SES.wif.issuerWarningSave);
          assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
          assert.true(
            this.flashSuccessSpy.calledWith('Issuer saved successfully'),
            'Success flash message called for issuer'
          );
          assert.ok(
            this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
            'Transitioned to the configuration index route.'
          );
        });

        test('shows modal when modifying the issuer, has correct payload, and shows flash message on fail', async function (assert) {
          assert.expect(7);
          this.issuer = 'http://foo.bar';
          this.server.post(configUrl('aws', this.id), () => {
            assert.true(
              true,
              'post request was made to config/root when unsetting the issuer. test should pass.'
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
          await fillIn(GENERAL.inputByAttr('roleArn'), 'some-other-value');
          await click(GENERAL.saveButton);
          assert.dom(SES.wif.issuerWarningModal).exists('issuer warning modal exists');

          await click(SES.wif.issuerWarningSave);
          assert.true(
            this.flashDangerSpy.calledWith('Issuer was not saved: permission denied'),
            'shows danger flash for issuer save'
          );
          assert.true(
            this.flashSuccessSpy.calledWith(`Successfully saved ${this.id}'s root configuration.`),
            "calls the root flash message not the issuer's"
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
        assert.expect(13);
        await this.renderComponent();
        assert.dom(SES.aws.rootForm).exists('it lands on the aws root configuration form.');
        assert.dom(SES.wif.accessTitle).exists('Access section is rendered');
        assert.dom(SES.aws.leaseTitle).exists('Lease section is rendered');
        assert
          .dom(SES.wif.accessTypeSection)
          .doesNotExist('Access type section does not render for a community user');
        // check all the form fields are present
        await click(GENERAL.toggleGroup('Root config options'));
        for (const key of expectedConfigKeys('aws-root-create')) {
          assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for root section.`);
        }
        for (const key of expectedConfigKeys('aws-lease')) {
          assert.dom(`[data-test-ttl-form-label="${key}"]`).exists(`${key} shows for Lease section.`);
        }
        assert.dom(GENERAL.inputByAttr('issuer')).doesNotExist();
      });
      test('it does not send issuer on save', async function (assert) {
        assert.expect(4);
        await this.renderComponent();
        this.server.post(configUrl('aws', this.id), () => {
          assert.true(true, 'post request was made to config/root. test should pass.');
        });
        this.server.post('/identity/oidc/config', () => {
          throw new Error('post request was incorrectly made to update issuer');
        });
        await fillInAwsConfig('withAccess');
        await fillInAwsConfig('withLease');
        await click(GENERAL.saveButton);
        assert.dom(SES.wif.issuerWarningModal).doesNotExist('modal should not render');
        assert.true(
          this.flashSuccessSpy.calledWith(`Successfully saved ${this.id}'s root configuration.`),
          'Flash message shows that root was saved even if issuer was not'
        );
        assert.ok(
          this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
          'Transitioned to the configuration index route.'
        );
      });
    });
  });
  module('Edit view', function (hooks) {
    hooks.beforeEach(function () {
      this.rootConfig = createConfig(this.store, this.id, 'aws');
      this.leaseConfig = createConfig(this.store, this.id, 'aws-lease');
    });
    module('isEnterprise', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'enterprise';
      });

      test('it defaults to IAM accessType if IAM fields are already set', async function (assert) {
        await this.renderComponent();
        assert.dom(SES.wif.accessType('iam')).isChecked('IAM accessType is checked');
        assert.dom(SES.wif.accessType('iam')).isDisabled('IAM accessType is disabled');
        assert.dom(SES.wif.accessType('wif')).isNotChecked('WIF accessType is not checked');
        assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
        assert
          .dom(SES.wif.accessTypeSubtext)
          .hasText('You cannot edit Access Type if you have already saved access credentials.');
      });

      test('it defaults to WIF accessType if WIF fields are already set', async function (assert) {
        this.rootConfig = createConfig(this.store, this.id, 'aws-wif');
        await this.renderComponent();
        assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
        assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
        assert.dom(SES.wif.accessType('iam')).isNotChecked('IAM accessType is not checked');
        assert.dom(SES.wif.accessType('iam')).isDisabled('IAM accessType is disabled');
        assert.dom(GENERAL.inputByAttr('roleArn')).hasValue(this.rootConfig.roleArn);
        assert
          .dom(SES.wif.accessTypeSubtext)
          .hasText('You cannot edit Access Type if you have already saved access credentials.');
        assert
          .dom(GENERAL.inputByAttr('identityTokenAudience'))
          .hasValue(this.rootConfig.identityTokenAudience);
        assert.dom(GENERAL.ttl.input('Identity token TTL')).hasValue('2'); // 7200 on payload is 2hrs in ttl picker
      });

      test('it renders issuer if global issuer is already set', async function (assert) {
        this.rootConfig = createConfig(this.store, this.id, 'aws-wif');
        this.issuerConfig = createConfig(this.store, this.id, 'issuer');
        await this.renderComponent();
        assert.dom(SES.wif.accessType('wif')).isChecked('WIF accessType is checked');
        assert.dom(SES.wif.accessType('wif')).isDisabled('WIF accessType is disabled');
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue(this.issuerConfig.issuer, 'it has the models issuer value');
      });

      test('it allows you to change access type if record does not have wif or iam values already set', async function (assert) {
        // the model does not have to be new for a user to see the option to change the access type.
        // the access type is only disabled if the model has values already set for access type fields.
        this.rootConfig = createConfig(this.store, this.id, 'aws-no-access');
        await this.renderComponent();
        assert.dom(SES.wif.accessType('wif')).isNotDisabled('WIF accessType is NOT disabled');
        assert.dom(SES.wif.accessType('iam')).isNotDisabled('IAM accessType is NOT disabled');
      });

      test('it shows previously saved root and lease information', async function (assert) {
        await this.renderComponent();
        assert.dom(GENERAL.inputByAttr('accessKey')).hasValue(this.rootConfig.accessKey);
        assert
          .dom(GENERAL.inputByAttr('secretKey'))
          .hasValue('**********', 'secretKey is masked on edit the value');

        await click(GENERAL.toggleGroup('Root config options'));
        assert.dom(GENERAL.inputByAttr('region')).hasValue(this.rootConfig.region);
        assert.dom(GENERAL.inputByAttr('iamEndpoint')).hasValue(this.rootConfig.iamEndpoint);
        assert.dom(GENERAL.inputByAttr('stsEndpoint')).hasValue(this.rootConfig.stsEndpoint);
        assert.dom(GENERAL.inputByAttr('maxRetries')).hasValue('1');
        // Check lease config values
        assert.dom(GENERAL.ttl.input('Default Lease TTL')).hasValue('50');
        assert.dom(GENERAL.ttl.input('Max Lease TTL')).hasValue('55');
      });

      test('it requires a double click to change the secret key', async function (assert) {
        await this.renderComponent();

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
    module('isCommunity', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });

      test('it does not show access types but defaults to iam fields', async function (assert) {
        await this.renderComponent();
        assert.dom(SES.wif.accessTypeSection).doesNotExist('Access type section does not render');
        assert.dom(GENERAL.inputByAttr('accessKey')).hasValue(this.rootConfig.accessKey);
        assert
          .dom(GENERAL.inputByAttr('secretKey'))
          .hasValue('**********', 'secretKey is masked on edit the value');

        await click(GENERAL.toggleGroup('Root config options'));
        assert.dom(GENERAL.inputByAttr('region')).hasValue(this.rootConfig.region);
        assert.dom(GENERAL.inputByAttr('iamEndpoint')).hasValue(this.rootConfig.iamEndpoint);
        assert.dom(GENERAL.inputByAttr('stsEndpoint')).hasValue(this.rootConfig.stsEndpoint);
        assert.dom(GENERAL.inputByAttr('maxRetries')).hasValue('1');
        // Check lease config values
        assert.dom(GENERAL.ttl.input('Default Lease TTL')).hasValue('50');
        assert.dom(GENERAL.ttl.input('Max Lease TTL')).hasValue('55');
      });
    });
  });
});
