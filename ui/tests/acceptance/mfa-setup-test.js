/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { create } from 'ember-cli-page-object';
import { v4 as uuidv4 } from 'uuid';
import { setupApplicationTest } from 'ember-qunit';
import { click, fillIn, visit, waitFor } from '@ember/test-helpers';
import { login, loginMethod } from 'vault/tests/helpers/auth/auth-helpers';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const consoleComponent = create(consoleClass);
const USER = 'end-user';
const PASSWORD = 'mypassword';
const POLICY_NAME = 'identity_policy';

const writePolicy = async function (path) {
  await visit('/vault/settings/auth/enable');
  await mountBackend('userpass', path);
  const identityPolicy = `path "identity/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
  }`;
  await consoleComponent.runCommands([
    `write sys/policies/acl/${POLICY_NAME} policy=${btoa(identityPolicy)}`,
  ]);
};

const writeUserWithPolicy = async function (path) {
  await consoleComponent.runCommands([
    `write auth/${path}/users/${USER} password=${PASSWORD} policies=${POLICY_NAME}`,
  ]);
};

const setupUser = async function (path) {
  await writePolicy(path);
  await writeUserWithPolicy(path);
  await click(GENERAL.submitButton);
};

module('Acceptance | mfa-setup', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const path = `userpass-${uuidv4()}`;
    await login();
    await setupUser(path);
    await loginMethod(
      { username: USER, password: PASSWORD, path },
      { authType: 'userpass', toggleOptions: true }
    );
    await click(GENERAL.button('user-menu-trigger'));
    await click('[data-test-user-menu-item="mfa"]');
  });

  test('it closes the dropdown after navigating', async function (assert) {
    assert
      .dom(GENERAL.button('user-menu-trigger'))
      .hasAttribute('aria-expanded', 'false', 'dropdown closes after navigating to MFA');
  });

  test('it should login through MFA and post to generate and be able to restart the setup', async function (assert) {
    assert.expect(5);
    // the network requests required in this test
    this.server.post('/identity/mfa/method/totp/generate', (scheme, req) => {
      const json = JSON.parse(req.requestBody);
      assert.strictEqual(json.method_id, '123', 'sends the UUID value');
      return {
        data: {
          barcode:
            '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',
          url: 'otpauth://totp/Vault:26606dbe-d8ea-82ca-41b0-1250a4484079?algorithm=SHA1&digits=6&issuer=Vault&period=30&secret=FID3WRPRRADQDN3CGPVVOLKCXTZZPSML',
          lease_duration: 0,
        },
      };
    });
    this.server.post('/identity/mfa/method/totp/admin-destroy', (scheme, req) => {
      const json = JSON.parse(req.requestBody);
      assert.strictEqual(json.method_id, '123', 'sends the UUID value');
      // returns nothing
      return {};
    });
    await fillIn('[data-test-input="uuid"]', 123);
    await click('[data-test-verify]');
    await waitFor('[data-test-qrcode]', { timeout: 5000 });
    assert.dom('[data-test-qrcode]').exists('the qrCode is shown.');
    assert.dom('[data-test-mfa-enabled-warning]').doesNotExist('warning does not show.');
    await click('[data-test-restart]');
    await waitFor('[data-test-step-one]', { timeout: 5000 });
    assert.dom('[data-test-step-one]').exists('back to step one.');
  });

  test('it should show a warning if you enter in the same UUID without restarting the setup', async function (assert) {
    assert.expect(2);
    // the network requests required in this test
    this.server.post('/identity/mfa/method/totp/generate', () => {
      return {
        data: null,
        warnings: ['Entity already has a secret for MFA method “”'],
      };
    });

    await fillIn('[data-test-input="uuid"]', 123);
    await click('[data-test-verify]');
    await waitFor('[data-test-mfa-enabled-warning]', { timeout: 5000 });
    assert.dom('[data-test-qrcode]').doesNotExist('the qrCode is not shown.');
    assert.dom('[data-test-mfa-enabled-warning]').exists('the mfa-enabled warning shows.');
  });
});
