/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/index';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { deleteAuthCmd, mountAuthCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | settings/auth/configure', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  test('it redirects to section options when there are no other sections', async function (assert) {
    const path = `approle-config-${this.uid}`;
    const type = 'approle';
    await enablePage.enable(type, path);
    await page.visit({ path });
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.auth.configure.section');
    assert.strictEqual(
      currentURL(),
      `/vault/settings/auth/configure/${path}/options`,
      'loads the options route'
    );
  });

  test('it redirects to the first section', async function (assert) {
    const path = `aws-redirect-${this.uid}`;
    const type = 'aws';
    await enablePage.enable(type, path);
    await page.visit({ path });
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.auth.configure.section');
    assert.strictEqual(
      currentURL(),
      `/vault/settings/auth/configure/${path}/client`,
      'loads the first section for the type of auth method'
    );
  });

  module('configure route does not require read on per-mount sys/auth/* path', function (hooks) {
    hooks.beforeEach(async function () {
      this.path = `approle-regression-${this.uid}`;
      const policyName = `configure-auth-regression-${this.uid}`;
      const policy = `
path "sys/auth/*" {
  capabilities = ["create", "update", "delete", "sudo", "patch"]
}
path "sys/auth*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo", "patch"]
}
      `.trim();

      await runCmd(mountAuthCmd('approle', this.path));
      const token = await runCmd(tokenWithPolicyCmd(policyName, policy));
      await login(token);
    });

    hooks.afterEach(async function () {
      await login();
      await runCmd(deleteAuthCmd(this.path));
    });

    test('it loads the configure route without a 403', async function (assert) {
      await page.visit({ path: this.path });
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.settings.auth.configure.section',
        'lands on the configure section route'
      );
      assert.dom(GENERAL.submitButton).exists('options form rendered with submit button');
      assert.dom('[data-test-field]').exists('tune fields are present in the options form');
    });

    test('it routes to the error page for a non-existent auth mount path', async function (assert) {
      await page.visit({ path: 'nonexistent-mount' });
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.error',
        'routes to the error page when the mount does not exist'
      );
    });
  });
});
