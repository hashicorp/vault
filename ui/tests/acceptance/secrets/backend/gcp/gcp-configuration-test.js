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
  fillInGcpConfig,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

module('Acceptance | GCP | configuration', function (hooks) {
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
    this.type = 'gcp';
    this.path = `GCP-${this.uid}`;
    return authPage.login();
  });

  test('it should prompt configuration after mounting the GCP engine', async function (assert) {
    await visit('/vault/settings/mount-secret-backend');
    await mountBackend(this.type, this.path);

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/configuration`,
      'navigated to configuration view'
    );
    assert.dom(GENERAL.emptyStateTitle).hasText('Google Cloud not configured');
    assert.dom(GENERAL.emptyStateActions).hasText('Configure Google Cloud');
    // cleanup
    await runCmd(`delete sys/mounts/${this.path}`);
  });

  test('it should transition to configure page on click "Configure" from toolbar', async function (assert) {
    await enablePage.enable(this.type, this.path);
    await click(SES.configure);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/configuration/edit`,
      'navigated to configuration edit view'
    );
    // cleanup
    await runCmd(`delete sys/mounts/${this.path}`);
  });

  module('Community', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'community';
    });

    module('details', function () {
      test('it should show configuration details with GCP account options configured', async function (assert) {
        const gcpAccountAttrs = {
          credentials: '{"some-key":"some-value"}',
          ttl: '1 minute 40 seconds',
          max_ttl: '1 minute 41 seconds',
        };
        this.server.get(`${this.path}/config`, () => {
          assert.true(true, 'request made to config when navigating to the configuration page.');
          return { data: { id: this.path, type: this.type, ...gcpAccountAttrs } };
        });
        await enablePage.enable(this.type, this.path);
        for (const key of expectedConfigKeys(this.type)) {
          if (key === 'Credentials') continue; // not returned by the API
          const responseKeyAndValue = expectedValueOfConfigKeys(this.type, key);
          assert
            .dom(GENERAL.infoRowValue(key))
            .hasText(responseKeyAndValue, `value for ${key} on the ${this.type} config details exists.`);
        }
        // check mount configuration details are present and accurate.
        await click(SES.configurationToggle);
        assert
          .dom(GENERAL.infoRowValue('Path'))
          .hasText(`${this.path}/`, 'mount path is displayed in the configuration details');
        // cleanup
        await runCmd(`delete sys/mounts/${this.path}`);
      });
    });

    module('create', function () {
      test('it should save gcp account accessType options', async function (assert) {
        await enablePage.enable(this.type, this.path);
        await click(SES.configTab);
        await click(SES.configure);
        await fillInGcpConfig();
        await click(GENERAL.saveButton);
        assert.true(
          this.flashSuccessSpy.calledWith(`Successfully saved ${this.path}'s configuration.`),
          'Success flash message is rendered showing the GCP model configuration was saved.'
        );

        assert
          .dom(GENERAL.infoRowValue('Config TTL'))
          .hasText('2 hours', 'Config TTL, a generic account specific field, has been set.');
        assert
          .dom(GENERAL.infoRowValue('Max TTL'))
          .hasText('2 hours 16 minutes 40 seconds', 'Max TTL, a generic field, has been set.');
        assert
          .dom(GENERAL.infoRowValue('Credentials'))
          .doesNotExist('credentials are not shown in the configuration details');
        // cleanup
        await runCmd(`delete sys/mounts/${this.path}`);
      });
    });

    module('edit', function (hooks) {
      hooks.beforeEach(async function () {
        const genericAttrs = {
          configTtl: '2h',
          maxTtl: '4h',
        };
        this.server.get(`${this.path}/config`, () => {
          return { data: { id: this.path, type: this.type, ...genericAttrs } };
        });
        await enablePage.enable(this.type, this.path);
      });

      test('it should save credentials', async function (assert) {
        assert.expect(3);
        const credentials = '{"some-key":"some-value"}';
        await click(SES.configure);

        this.server.post(configUrl('gcp', this.path), (schema, req) => {
          const payload = JSON.parse(req.requestBody);
          assert.strictEqual(credentials, payload.credentials, 'credentials are sent in post request');
          assert.strictEqual(
            undefined,
            payload.configTtl,
            'config_ttl is not included in payload if value has not been updated'
          );
          assert.strictEqual(
            undefined,
            payload.maxTtl,
            'max_ttl is not included in payload if value has not been updated'
          );
        });

        await click(GENERAL.textToggle);
        await fillIn(GENERAL.textToggleTextarea, credentials);
        await click(GENERAL.saveButton);
        // cleanup
        await runCmd(`delete sys/mounts/${this.path}`);
      });

      test('it should not save credentials if it has NOT been changed', async function (assert) {
        assert.expect(3);
        await click(SES.configure);

        this.server.post(configUrl('gcp', this.path), (schema, req) => {
          const payload = JSON.parse(req.requestBody);
          assert.strictEqual(payload.credentials, undefined, 'credentials are not sent in post request');

          assert.strictEqual(
            payload.ttl,
            '10800s',
            'config_ttl is included in payload because the value was updated'
          );
          assert.strictEqual(
            payload.maxTtl,
            undefined,
            'max_ttl is not included in payload if value has not been updated'
          );
        });

        await click(GENERAL.toggleGroup('More options'));
        await click(GENERAL.ttl.toggle('Config TTL'));
        await fillIn(GENERAL.ttl.input('Config TTL'), '10800');
        await click(GENERAL.saveButton);
        // cleanup
        await runCmd(`delete sys/mounts/${this.path}`);
      });
    });

    module('Error handling', function () {
      test('it prevents transition and shows api error if config errored on save', async function (assert) {
        await enablePage.enable(this.type, this.path);

        this.server.post(configUrl(this.type, this.path), () => {
          return overrideResponse(400, { errors: ['my goodness, that did not work!'] });
        });

        await click(SES.configTab);
        await click(SES.configure);
        await fillInGcpConfig();
        await click(GENERAL.saveButton);

        assert
          .dom(GENERAL.messageError)
          .hasText('Error my goodness, that did not work!', 'API error shows on form');
        assert.strictEqual(
          currentURL(),
          `/vault/secrets/${this.path}/configuration/edit`,
          'the form did not transition because the save failed.'
        );
        // cleanup
        await runCmd(`delete sys/mounts/${this.path}`);
      });

      test('it should show API error when configuration read fails', async function (assert) {
        this.server.get(configUrl(this.type, this.path), () => {
          return overrideResponse(400, { errors: ['bad request'] });
        });
        await enablePage.enable(this.type, this.path);
        assert.dom(SES.error.title).hasText('Error', 'shows the secrets backend error route');
      });
    });
  });

  module('Enterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
      this.wifAttrs = {
        service_account_email: 'service-email',
        identity_token_audience: 'audience',
        identity_token_ttl: 720000,
        max_ttl: 14400,
        ttl: 3600,
      };
    });

    test('it should show configuration details with WIF options configured', async function (assert) {
      this.server.get(`${this.path}/config`, () => {
        assert.true(true, 'request made to config when navigating to the configuration page.');
        return { data: { id: this.path, type: this.type, ...this.wifAttrs } };
      });
      await enablePage.enable(this.type, this.path);
      for (const key of expectedConfigKeys('gcp-wif')) {
        const responseKeyAndValue = expectedValueOfConfigKeys(this.type, key);
        assert
          .dom(GENERAL.infoRowValue(key))
          .hasText(responseKeyAndValue, `value for ${key} on the ${this.type} config details exists.`);
      }
      // check mount configuration details are present and accurate.
      await click(SES.configurationToggle);
      assert
        .dom(GENERAL.infoRowValue('Path'))
        .hasText(`${this.path}/`, 'mount path is displayed in the configuration details');
      // cleanup
      await runCmd(`delete sys/mounts/${this.path}`);
    });

    module('Error handling', function () {
      test('it shows API error if user previously set credentials but tries to edit the configuration with wif fields', async function (assert) {
        await enablePage.enable(this.type, this.path);
        await click(SES.configTab);
        await click(SES.configure);
        await fillInGcpConfig();
        await click(GENERAL.saveButton); // save GCP credentials

        await click(SES.configure); // navigate so you can edit that configuration
        await fillInGcpConfig(true);
        await click(GENERAL.saveButton); // try and save wif fields
        assert
          .dom(GENERAL.messageError)
          .hasText(
            `Error only one of 'credentials' or 'identity_token_audience' can be set`,
            'api error about conflicting fields is shown'
          );
        assert.dom(GENERAL.inlineError).hasText('There was an error submitting this form.');
        // cleanup
        await runCmd(`delete sys/mounts/${this.path}`);
      });
    });
  });
});
