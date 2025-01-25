/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, visit, currentURL, fillIn } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { spy } from 'sinon';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import {
  expectedConfigKeys,
  expectedValueOfConfigKeys,
  configUrl,
  createConfig,
  fillInAzureConfig,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

module('Acceptance | Azure | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const flash = this.owner.lookup('service:flash-messages');
    this.store = this.owner.lookup('service:store');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashDangerSpy = spy(flash, 'danger');
    this.flashInfoSpy = spy(flash, 'info');
    this.version = this.owner.lookup('service:version');
    this.uid = uuidv4();
    this.type = 'azure';
    return authPage.login();
  });

  test('it should prompt configuration after mounting the azure engine', async function (assert) {
    const path = `azure-${this.uid}`;
    await visit('/vault/settings/mount-secret-backend');
    await mountBackend(this.type, path);

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/configuration`,
      'navigated to configuration view'
    );
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('Azure not configured', "empty state title is 'Azure not configured'");
    assert.dom(GENERAL.emptyStateActions).hasText('Configure Azure');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should transition to configure page on click "Configure" from toolbar', async function (assert) {
    const path = `azure-${this.uid}`;
    await enablePage.enable(this.type, path);
    await click(SES.configure);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/configuration/edit`,
      'navigated to configuration edit view'
    );
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  module('isCommunity', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'community';
    });

    module('details', function () {
      test('it should show configuration with Azure account options configured', async function (assert) {
        const path = `azure-${this.uid}`;
        const azureAccountAttrs = {
          client_secret: 'client-secret',
          subscription_id: 'subscription-id',
          tenant_id: 'tenant-id',
          client_id: 'client-id',
          root_password_ttl: '20 days 20 hours',
          environment: 'AZUREPUBLICCLOUD',
        };
        this.server.get(`${path}/config`, () => {
          assert.true(true, 'request made to config when navigating to the configuration page.');
          return { data: { id: path, type: this.type, ...azureAccountAttrs } };
        });
        await enablePage.enable(this.type, path);
        for (const key of expectedConfigKeys('azure')) {
          if (key === 'Client secret') continue; // client-secret is not returned by the API
          assert.dom(GENERAL.infoRowLabel(key)).exists(`${key} on the ${this.type} config details exists.`);
          const responseKeyAndValue = expectedValueOfConfigKeys(this.type, key);
          assert
            .dom(GENERAL.infoRowValue(key))
            .hasText(responseKeyAndValue, `value for ${key} on the ${this.type} config details exists.`);
        }
        // check mount configuration details are present and accurate.
        await click(SES.configurationToggle);
        assert
          .dom(GENERAL.infoRowValue('Path'))
          .hasText(`${path}/`, 'mount path is displayed in the configuration details');
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });
    });

    module('create', function () {
      test('it should save azure account accessType options', async function (assert) {
        assert.expect(3);
        const path = `azure-${this.uid}`;
        await enablePage.enable(this.type, path);

        this.server.post('/identity/oidc/config', () => {
          throw new Error(
            `Request was made to return the issuer when it should not have been because user is on CE.`
          );
        });

        await click(SES.configTab);
        await click(SES.configure);
        await fillInAzureConfig(this.type);
        await click(GENERAL.saveButton);
        assert.true(
          this.flashSuccessSpy.calledWith(`Successfully saved ${path}'s configuration.`),
          'Success flash message is rendered showing the azure model configuration was saved.'
        );
        assert
          .dom(GENERAL.infoRowValue('Root password TTL'))
          .hasText(
            '1 hour 26 minutes 40 seconds',
            'Root password TTL, an azure account specific field, has been set.'
          );
        assert
          .dom(GENERAL.infoRowValue('Subscription ID'))
          .hasText('subscription-id', 'Subscription ID, a generic field, has been set.');
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });
    });

    module('edit', function (hooks) {
      hooks.beforeEach(async function () {
        const path = `azure-${this.uid}`;
        const type = 'azure';
        const genericAttrs = {
          // attributes that can be used for either wif or azure account access type
          subscription_id: 'subscription-id',
          tenant_id: 'tenant-id',
          client_id: 'client-id',
          environment: 'AZUREPUBLICCLOUD',
        };
        this.server.get(`${path}/config`, () => {
          return { data: { id: path, type, ...genericAttrs } };
        });
        await enablePage.enable(type, path);
      });

      test('it should not save client secret if it has NOT been changed', async function (assert) {
        assert.expect(2);
        await click(SES.configure);
        const url = currentURL();
        const path = url.split('/')[3]; // get path from url because we can't pass the path from beforeEach hook to individual test.
        this.server.post(configUrl('azure', path), (schema, req) => {
          const payload = JSON.parse(req.requestBody);
          assert.strictEqual(
            undefined,
            payload.client_secret,
            'client_secret is not sent if it has not been changed'
          );
          assert.strictEqual(
            payload.subscription_id,
            'subscription-id-updated',
            'subscription_id is included with updated value in the payload'
          );
        });
        await fillIn(GENERAL.inputByAttr('subscriptionId'), 'subscription-id-updated');
        await click(GENERAL.enableField('clientSecret'));
        await click(GENERAL.saveButton);
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });

      test('it should save client secret if it HAS been changed', async function (assert) {
        assert.expect(2);
        await click(SES.configure);
        const url = currentURL();
        const path = url.split('/')[3]; // get path from url because we can't pass the path from beforeEach hook to individual test.
        this.server.post(configUrl('azure', path), (schema, req) => {
          const payload = JSON.parse(req.requestBody);
          assert.strictEqual(
            payload.client_secret,
            'client-secret-updated',
            'client_secret is sent on payload because user updated its value'
          );
          assert.strictEqual(
            payload.subscription_id,
            'subscription-id-updated-again',
            'subscription_id is included with updated value in the payload'
          );
        });
        await fillIn(GENERAL.inputByAttr('subscriptionId'), 'subscription-id-updated-again');
        await click(GENERAL.enableField('clientSecret'));
        await click('[data-test-button="toggle-masked"]');
        await fillIn(GENERAL.inputByAttr('clientSecret'), 'client-secret-updated');
        await click(GENERAL.saveButton);
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });
    });

    module('Error handling', function () {
      test('it prevents transition and shows api error if config errored on save', async function (assert) {
        const path = `azure-${this.uid}`;
        await enablePage.enable('azure', path);

        this.server.post(configUrl('azure', path), () => {
          return overrideResponse(400, { errors: ['welp, that did not work!'] });
        });

        await click(SES.configTab);
        await click(SES.configure);
        await fillInAzureConfig('azure');
        await click(GENERAL.saveButton);

        assert.dom(GENERAL.messageError).hasText('Error welp, that did not work!', 'API error shows on form');
        assert.strictEqual(
          currentURL(),
          `/vault/secrets/${path}/configuration/edit`,
          'the form did not transition because the save failed.'
        );
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });

      test('it should show API error when configuration read fails', async function (assert) {
        assert.expect(1);
        const path = `azure-${this.uid}`;
        // interrupt get and return API error
        this.server.get(configUrl(this.type, path), () => {
          return overrideResponse(400, { errors: ['bad request'] });
        });
        await enablePage.enable(this.type, path);
        assert.dom(SES.error.title).hasText('Error', 'shows the secrets backend error route');
      });
    });
  });

  module('isEnterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
    });

    module('details', function () {
      test('it should save WIF configuration options', async function (assert) {
        const path = `azure-${this.uid}`;
        const wifAttrs = {
          subscription_id: 'subscription-id',
          tenant_id: 'tenant-id',
          client_id: 'client-id',
          identity_token_audience: 'audience',
          identity_token_ttl: 720000,
          environment: 'AZUREPUBLICCLOUD',
        };
        this.server.get(`${path}/config`, () => {
          assert.true(true, 'request made to config when navigating to the configuration page.');
          return { data: { id: path, type: this.type, ...wifAttrs } };
        });
        await enablePage.enable(this.type, path);
        for (const key of expectedConfigKeys('azure-wif')) {
          const responseKeyAndValue = expectedValueOfConfigKeys(this.type, key);
          assert
            .dom(GENERAL.infoRowValue(key))
            .hasText(responseKeyAndValue, `value for ${key} on the ${this.type} config details exists.`);
        }
        // check mount configuration details are present and accurate.
        await click(SES.configurationToggle);
        assert
          .dom(GENERAL.infoRowValue('Path'))
          .hasText(`${path}/`, 'mount path is displayed in the configuration details');
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });

      test('it should not show issuer if no WIF configuration data is returned', async function (assert) {
        const path = `azure-${this.uid}`;
        this.server.get(`${path}/config`, (schema, req) => {
          const payload = JSON.parse(req.requestBody);
          assert.true(true, 'request made to config/root when navigating to the configuration page.');
          return { data: { id: path, type: this.type, attributes: payload } };
        });
        this.server.get(`identity/oidc/config`, () => {
          throw new Error(`Request was made to return the issuer when it should not have been.`);
        });
        await createConfig(this.store, path, this.type); // create the azure account config in the store
        await enablePage.enable(this.type, path);
        await click(SES.configTab);

        assert.dom(GENERAL.infoRowLabel('Issuer')).doesNotExist(`Issuer does not exists on config details.`);
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });
    });

    module('create', function () {
      test('it should transition and save issuer if model was not changed but issuer was', async function (assert) {
        assert.expect(3);
        const path = `azure-${this.uid}`;
        await enablePage.enable(this.type, path);
        const newIssuer = `http://new.issuer.${this.uid}`;

        this.server.post('/identity/oidc/config', (schema, req) => {
          const payload = JSON.parse(req.requestBody);
          assert.deepEqual(payload, { issuer: newIssuer }, 'payload for issuer is correct');
          return {
            id: 'identity-oidc-config',
            data: null,
            warnings: [
              'If "issuer" is set explicitly, all tokens must be validated against that address, including those issued by secondary clusters. Setting issuer to "" will restore the default behavior of using the cluster\'s api_addr as the issuer.',
            ],
          };
        });
        this.server.post(configUrl(this.type, path), () => {
          throw new Error('post request was incorrectly made to update the azure config model');
        });

        await click(SES.configTab);
        await click(SES.configure);
        await click(SES.wif.accessType('wif'));
        await fillIn(GENERAL.inputByAttr('issuer'), newIssuer);
        await click(GENERAL.saveButton);
        await click(SES.wif.issuerWarningSave);
        assert.true(
          this.flashSuccessSpy.calledWith(`Issuer saved successfully`),
          'Shows issuer saved message'
        );

        assert
          .dom(GENERAL.emptyStateTitle)
          .hasText(
            'Azure not configured',
            'Empty state message is displayed because the model was not saved only the issuer'
          );
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });

      test('it should transition and save model if the model was changed but issuer was not', async function (assert) {
        assert.expect(4);
        const path = `azure-${this.uid}`;
        await enablePage.enable(this.type, path);

        this.server.post('/identity/oidc/config', () => {
          throw new Error('post request was incorrectly made to update the issuer');
        });

        await click(SES.configTab);
        await click(SES.configure);
        await fillInAzureConfig('withWif');
        await click(GENERAL.saveButton);
        assert.dom(SES.wif.issuerWarningModal).doesNotExist('issuer warning modal does not show');
        assert.true(
          this.flashSuccessSpy.calledWith(`Successfully saved ${path}'s configuration.`),
          'Success flash message is rendered showing the azure model configuration was saved.'
        );
        assert
          .dom(GENERAL.infoRowValue('Identity token audience'))
          .hasText('azure-audience', 'Identity token audience has been set.');
        assert
          .dom(GENERAL.infoRowValue('Identity token TTL'))
          .hasText('2 hours', 'Identity token TTL has been set.');
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });

      test('it should transition and show issuer error if model saved but issuer encountered an error', async function (assert) {
        const path = `azure-${this.uid}`;
        const oldIssuer = 'http://old.issuer';
        await enablePage.enable(this.type, path);
        this.server.get('/identity/oidc/config', () => {
          return { issuer: 'http://old.issuer' };
        });
        this.server.post('/identity/oidc/config', () => {
          return overrideResponse(400, { errors: ['bad request'] });
        });

        await click(SES.configTab);
        await click(SES.configure);
        await fillInAzureConfig('withWif');
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue(oldIssuer, 'issuer defaults to previously saved value');
        await fillIn(GENERAL.inputByAttr('issuer'), 'http://new.issuererrors');
        await click(GENERAL.saveButton);
        await click(SES.wif.issuerWarningSave);
        assert.true(
          this.flashSuccessSpy.calledWith(`Successfully saved ${path}'s configuration.`),
          'Success flash message is rendered showing the azure model configuration was saved.'
        );
        assert.true(
          this.flashDangerSpy.calledWith(`Issuer was not saved: bad request`),
          'Danger flash message is rendered showing the issuer was not saved.'
        );
        assert
          .dom(GENERAL.infoRowValue('Identity token audience'))
          .hasText('azure-audience', 'Identity token audience has been set.');
        assert
          .dom(GENERAL.infoRowValue('Identity token TTL'))
          .hasText('2 hours', 'Identity token TTL has been set.');
        assert
          .dom(GENERAL.infoRowValue('Issuer'))
          .hasText(oldIssuer, 'Issuer is shows the old saved value not the new value that errors on save.');
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });

      test('it should NOT transition and show error if model errored but issuer was saved', async function (assert) {
        const path = `azure-${this.uid}`;
        const newIssuer = `http://new.issuer.${this.uid}`;
        const oldIssuer = 'http://old.issuer';
        await enablePage.enable(this.type, path);
        this.server.get('/identity/oidc/config', () => {
          return { issuer: oldIssuer };
        });
        this.server.post(configUrl(this.type, path), () => {
          return overrideResponse(400, { errors: ['bad request'] });
        });

        await click(SES.configTab);
        await click(SES.configure);
        await fillInAzureConfig('withWif');
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue(oldIssuer, 'issuer defaults to previously saved value');
        await fillIn(GENERAL.inputByAttr('issuer'), newIssuer);
        await click(GENERAL.saveButton);
        await click(SES.wif.issuerWarningSave);
        assert.true(
          this.flashSuccessSpy.calledWith(`Issuer saved successfully`),
          'Success flash message is rendered showing the issuer configuration was saved.'
        );
        assert.dom(GENERAL.messageError).hasText('Error bad request', 'Error message is displayed.');
        assert.strictEqual(
          currentURL(),
          `/vault/secrets/${path}/configuration/edit`,
          'stays on the edit page'
        );
        assert
          .dom(GENERAL.inputByAttr('issuer'))
          .hasValue(newIssuer, 'issuer is updated to newly saved value');
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });
    });

    module('edit', function () {
      test('it should update WIF attributes', async function (assert) {
        const path = `azure-${this.uid}`;
        await enablePage.enable(this.type, path);
        await click(SES.configTab);
        await click(SES.configure);
        await fillInAzureConfig('withWif');
        await click(GENERAL.saveButton); // finished creating attributes, go back and edit them.
        assert
          .dom(GENERAL.infoRowValue('Identity token audience'))
          .hasText('azure-audience', `value for identity token audience shows on the config details view.`);
        await click(SES.configure);
        await fillIn(GENERAL.inputByAttr('identityTokenAudience'), 'new-audience');
        await click(GENERAL.saveButton);
        assert
          .dom(GENERAL.infoRowValue('Identity token audience'))
          .hasText('new-audience', `value for identity token audience shows on the config details view.`);
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      });
    });
  });
});
