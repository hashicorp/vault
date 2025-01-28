/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, visit, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

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

module('Acceptance | GCP | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.uid = uuidv4();
    this.type = 'gcp';
    return authPage.login();
  });
  module('isEnterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
    });

    test('it should show empty state and navigate to configuration view after mounting the GCP engine', async function (assert) {
      const path = `GCP-${this.uid}`;
      await visit('/vault/settings/mount-secret-backend');
      await mountBackend(this.type, path);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${path}/configuration`,
        'navigated to configuration view'
      );
      assert.dom(GENERAL.emptyStateTitle).hasText('Google Cloud not configured');
      assert.dom(GENERAL.emptyStateActions).doesNotContainText('Configure GCP');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it should not show "Configure" from toolbar', async function (assert) {
      const path = `GCP-${this.uid}`;
      await enablePage.enable(this.type, path);
      assert.dom(SES.configure).doesNotExist('Configure button does not exist.');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it should show configuration details with WIF options configured', async function (assert) {
      const path = `GCP-${this.uid}`;
      const wifAttrs = {
        service_account_email: 'service-email',
        identity_token_audience: 'audience',
        identity_token_ttl: 720000,
        max_ttl: 14400,
        ttl: 3600,
      };
      this.server.get(`${path}/config`, () => {
        assert.true(true, 'request made to config when navigating to the configuration page.');
        return { data: { id: path, type: this.type, ...wifAttrs } };
      });
      await enablePage.enable(this.type, path);
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
        .hasText(`${path}/`, 'mount path is displayed in the configuration details');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it should show configuration details with GCP account options configured', async function (assert) {
      const path = `GCP-${this.uid}`;
      const GCPAccountAttrs = {
        credentials: '{"some-key":"some-value"}',
        ttl: '1 hour',
        max_ttl: '4 hours',
      };
      this.server.get(`${path}/config`, () => {
        assert.true(true, 'request made to config when navigating to the configuration page.');
        return { data: { id: path, type: this.type, ...GCPAccountAttrs } };
      });
      await enablePage.enable(this.type, path);
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
        .hasText(`${path}/`, 'mount path is displayed in the configuration details');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it should show API error when configuration read fails', async function (assert) {
      assert.expect(1);
      const path = `GCP-${this.uid}`;
      // interrupt get and return API error
      this.server.get(configUrl(this.type, path), () => {
        return overrideResponse(400, { errors: ['bad request'] });
      });
      await enablePage.enable(this.type, path);
      assert.dom(SES.error.title).hasText('Error', 'shows the secrets backend error route');
    });
  });
});
