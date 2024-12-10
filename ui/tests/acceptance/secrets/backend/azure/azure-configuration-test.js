/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, visit, currentURL } from '@ember/test-helpers';
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
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

module('Acceptance | Azure | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const flash = this.owner.lookup('service:flash-messages');
    this.store = this.owner.lookup('service:store');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashInfoSpy = spy(flash, 'info');
    this.version = this.owner.lookup('service:version');
    this.uid = uuidv4();
    return authPage.login();
  });
  module('isEnterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
    });

    test('it should show empty state and navigate to configuration view after mounting the azure engine', async function (assert) {
      const path = `azure-${this.uid}`;
      await visit('/vault/settings/mount-secret-backend');
      await mountBackend('azure', path);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${path}/configuration`,
        'navigated to configuration view'
      );
      assert.dom(GENERAL.emptyStateTitle).hasText('Azure not configured');
      assert.dom(GENERAL.emptyStateActions).doesNotContainText('Configure Azure');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it should not show "Configure" from toolbar', async function (assert) {
      const path = `azure-${this.uid}`;
      await enablePage.enable('azure', path);
      assert.dom(SES.configure).doesNotExist('Configure button does not exist.');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it should show configuration with WIF options configured', async function (assert) {
      const path = `azure-${this.uid}`;
      const type = 'azure';
      const wifAttrs = {
        subscription_id: 'subscription-id',
        tenant_id: 'tenant-id',
        client_id: 'client-id',
        identity_token_audience: 'audience',
        identity_token_ttl: 720000,
        environment: 'AZUREPUBLICCLOUD',
      };
      this.server.get(`${path}/config`, () => {
        assert.ok(true, 'request made to config when navigating to the configuration page.');
        return { data: { id: path, type, ...wifAttrs } };
      });
      await enablePage.enable(type, path);
      for (const key of expectedConfigKeys('azure-wif')) {
        const responseKeyAndValue = expectedValueOfConfigKeys(type, key);
        assert
          .dom(GENERAL.infoRowValue(key))
          .hasText(responseKeyAndValue, `value for ${key} on the ${type} config details exists.`);
      }
      // check mount configuration details are present and accurate.
      await click(SES.configurationToggle);
      assert
        .dom(GENERAL.infoRowValue('Path'))
        .hasText(`${path}/`, 'mount path is displayed in the configuration details');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it should show configuration with Azure account options configured', async function (assert) {
      const path = `azure-${this.uid}`;
      const type = 'azure';
      const azureAccountAttrs = {
        client_secret: 'client-secret',
        subscription_id: 'subscription-id',
        tenant_id: 'tenant-id',
        client_id: 'client-id',
        root_password_ttl: '20 days 20 hours',
        environment: 'AZUREPUBLICCLOUD',
      };
      this.server.get(`${path}/config`, () => {
        assert.ok(true, 'request made to config when navigating to the configuration page.');
        return { data: { id: path, type, ...azureAccountAttrs } };
      });
      await enablePage.enable(type, path);
      for (const key of expectedConfigKeys('azure')) {
        assert.dom(GENERAL.infoRowLabel(key)).exists(`${key} on the ${type} config details exists.`);
        const responseKeyAndValue = expectedValueOfConfigKeys(type, key);
        assert
          .dom(GENERAL.infoRowValue(key))
          .hasText(responseKeyAndValue, `value for ${key} on the ${type} config details exists.`);
      }
      // check mount configuration details are present and accurate.
      await click(SES.configurationToggle);
      assert
        .dom(GENERAL.infoRowValue('Path'))
        .hasText(`${path}/`, 'mount path is displayed in the configuration details');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it should show API error when configuration read fails', async function (assert) {
      assert.expect(1);
      const path = `azure-${this.uid}`;
      const type = 'azure';
      // interrupt get and return API error
      this.server.get(configUrl(type, path), () => {
        return overrideResponse(400, { errors: ['bad request'] });
      });
      await enablePage.enable(type, path);
      assert.dom(SES.error.title).hasText('Error', 'shows the secrets backend error route');
    });
  });
});
