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

module('Integration | Component | SecretEngine/configure-aws', function (hooks) {
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
    this.rootConfig = this.store.createRecord('aws/root-config');
    this.leaseConfig = this.store.createRecord('aws/lease-config');
    // Add backend to the configs because it's not on the testing snapshot (would come from url)
    this.rootConfig.backend = this.leaseConfig.backend = this.id;
    this.version = this.owner.lookup('service:version');

    this.renderComponent = () => {
      return render(hbs`
        <SecretEngine::ConfigureAws @rootConfig={{this.rootConfig}} @leaseConfig={{this.leaseConfig}} @backendPath={{this.id}} />
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
        assert.dom(SES.aws.accessTitle).exists('Access section is rendered');
        assert.dom(SES.aws.leaseTitle).exists('Lease section is rendered');
        assert.dom(SES.aws.accessTypeSection).exists('Access type section is rendered');
        assert.dom(SES.aws.accessType('iam')).isChecked('defaults to showing IAM access type checked');
        assert.dom(SES.aws.accessType('wif')).isNotChecked('wif access type is not checked');
        // check all the form fields are present
        await click(GENERAL.toggleGroup('Root config options'));
        for (const key of expectedConfigKeys('aws-root-create')) {
          if (key === 'secretKey') {
            assert.dom(GENERAL.maskedInput(key)).exists(`${key} shows for root section.`);
          } else {
            assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for root section.`);
          }
        }
        for (const key of expectedConfigKeys('aws-lease')) {
          assert.dom(`[data-test-ttl-form-label="${key}"]`).exists(`${key} shows for Lease section.`);
        }
      });

      test('it renders wif fields when selected', async function (assert) {
        await this.renderComponent();
        await click(SES.aws.accessType('wif'));
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
          if (key === 'secretKey') {
            assert.dom(GENERAL.maskedInput(key)).doesNotExist(`${key} does not show when wif is selected.`);
          } else {
            assert.dom(GENERAL.inputByAttr(key)).doesNotExist(`${key} does not show when wif is selected.`);
          }
        }
      });

      test('it clears wif/iam inputs after toggling accessType', async function (assert) {
        await this.renderComponent();
        await fillInAwsConfig(true, false, true); // fill in IAM fields
        await click(SES.aws.accessType('wif')); // toggle to wif
        await fillInAwsConfig(false, false, false, true); // fill in wif fields
        await click(SES.aws.accessType('iam')); // toggle to wif
        assert
          .dom(GENERAL.inputByAttr('accessKey'))
          .hasValue('', 'accessKey is cleared after toggling accessType');
        assert
          .dom(GENERAL.maskedInput('secretKey'))
          .hasValue('', 'secretKey is cleared after toggling accessType');

        await click(SES.aws.accessType('wif'));
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
        await click(SES.aws.save);
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
        await fillInAwsConfig(true, false, true);
        await click(SES.aws.save);
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
        await fillInAwsConfig(true, false, true);
        await click(SES.aws.save);

        assert.true(
          this.flashDangerSpy.calledWith('Lease configuration was not saved: bad request'),
          'Flash message shows that lease was not saved.'
        );
        assert.ok(
          this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
          'Transitioned to the configuration index route.'
        );
      });

      test('it transitions without sending a lease or root payload on cancel', async function (assert) {
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
        // fill in both lease and root endpoints to ensure that both payloads are attempted to be sent
        await fillInAwsConfig(true, false, true);
        await click(SES.aws.cancel);

        assert.true(this.flashDangerSpy.notCalled, 'No danger flash messages called.');
        assert.true(this.flashSuccessSpy.notCalled, 'No success flash messages called.');
        assert.ok(
          this.transitionStub.calledWith('vault.cluster.secrets.backend.configuration', this.id),
          'Transitioned to the configuration index route.'
        );
      });
    });
    module('isCommunity', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });
      test('it renders fields', async function (assert) {
        assert.expect(12);
        await this.renderComponent();
        assert.dom(SES.aws.rootForm).exists('it lands on the aws root configuration form.');
        assert.dom(SES.aws.accessTitle).exists('Access section is rendered');
        assert.dom(SES.aws.leaseTitle).exists('Lease section is rendered');
        assert
          .dom(SES.aws.accessTypeSection)
          .doesNotExist('Access type section does not render for a community user');
        // check all the form fields are present
        await click(GENERAL.toggleGroup('Root config options'));
        for (const key of expectedConfigKeys('aws-root-create')) {
          if (key === 'secretKey') {
            assert.dom(GENERAL.maskedInput(key)).exists(`${key} shows for root section.`);
          } else {
            assert.dom(GENERAL.inputByAttr(key)).exists(`${key} shows for root section.`);
          }
        }
        for (const key of expectedConfigKeys('aws-lease')) {
          assert.dom(`[data-test-ttl-form-label="${key}"]`).exists(`${key} shows for Lease section.`);
        }
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
        assert.dom(SES.aws.accessType('iam')).isChecked('IAM accessType is checked');
        assert.dom(SES.aws.accessType('iam')).isDisabled('IAM accessType is disabled');
        assert.dom(SES.aws.accessType('wif')).isNotChecked('WIF accessType is not checked');
        assert.dom(SES.aws.accessType('wif')).isDisabled('WIF accessType is disabled');
        assert
          .dom(SES.aws.accessTypeSubtext)
          .hasText('You cannot edit Access Type if you have already saved access credentials.');
      });

      test('it defaults to WIF accessType if WIF fields are already set', async function (assert) {
        this.rootConfig = createConfig(this.store, this.id, 'aws-wif');
        await this.renderComponent();
        assert.dom(SES.aws.accessType('wif')).isChecked('WIF accessType is checked');
        assert.dom(SES.aws.accessType('wif')).isDisabled('WIF accessType is disabled');
        assert.dom(SES.aws.accessType('iam')).isNotChecked('IAM accessType is not checked');
        assert.dom(SES.aws.accessType('iam')).isDisabled('IAM accessType is disabled');
        assert.dom(GENERAL.inputByAttr('roleArn')).hasValue(this.rootConfig.roleArn);
        assert
          .dom(SES.aws.accessTypeSubtext)
          .hasText('You cannot edit Access Type if you have already saved access credentials.');
        assert
          .dom(GENERAL.inputByAttr('identityTokenAudience'))
          .hasValue(this.rootConfig.identityTokenAudience);
        assert.dom(GENERAL.ttl.input('Identity token TTL')).hasValue('2'); // 7200 on payload is 2hrs in ttl picker
      });

      test('it allows you to change access type if record does not have wif or iam values already set', async function (assert) {
        // the model does not have to be new for a user to see the option to change the access type.
        // the access type is only disabled if the model has values already set for access type fields.
        this.rootConfig = createConfig(this.store, this.id, 'aws-no-access');
        await this.renderComponent();
        assert.dom(SES.aws.accessType('wif')).isNotDisabled('WIF accessType is NOT disabled');
        assert.dom(SES.aws.accessType('iam')).isNotDisabled('IAM accessType is NOT disabled');
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
        await fillIn(GENERAL.maskedInput('secretKey'), 'new-secret');
        await click(SES.aws.save);
      });
    });
    module('isCommunity', function (hooks) {
      hooks.beforeEach(function () {
        this.version.type = 'community';
      });

      test('it does not show access types but defaults to iam fields', async function (assert) {
        await this.renderComponent();
        assert.dom(SES.aws.accessTypeSection).doesNotExist('Access type section does not render');
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
