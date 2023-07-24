/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { currentURL, visit } from '@ember/test-helpers';
import {
  adminPolicy,
  dataSecretPathCreateReadUpdate,
  metadataListOnly,
} from 'vault/tests/helpers/policy-generator/kv';
import { tokenWithPolicy, runCommands, writeSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { SELECTORS } from 'vault/tests/helpers/kv/kv-general-selectors';

// This test module should test KV permissions views, each sub-module is a separate tab (i.e. secret, metadata)

module('Acceptance | kv permissions', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    // Setup KV engine
    const mountPath = `kv-engine-${uuidv4()}`;
    await enablePage.enable('kv', mountPath);
    this.mountPath = mountPath;
    await logout.visit();
  });

  hooks.afterEach(async function () {
    await logout.visit();
    await authPage.login();
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.mountPath}`]);
    await logout.visit();
  });

  module('secret', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      this.secretPath = `my-secret-${uuidv4()}`;
      await writeSecret(this.mountPath, this.secretPath, 'foo', 'bar');

      const kv_admin_policy = adminPolicy(this.mountPath);
      this.kvAdminToken = await tokenWithPolicy('kv-admin', kv_admin_policy);

      const no_metadata_read_policy =
        metadataListOnly(this.mountPath) + dataSecretPathCreateReadUpdate(this.mountPath, this.secretPath);
      this.kvNoMetadataRead = await tokenWithPolicy('kv-no-metadata-read', no_metadata_read_policy);

      await logout.visit();
    });

    test('it shows all tabs for admin policy', async function (assert) {
      assert.expect(5);
      await authPage.login(this.kvAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
      assert.dom(SELECTORS.secretTab('Secret')).exists();
      assert.dom(SELECTORS.secretTab('Metadata')).exists();
      assert.dom(SELECTORS.secretTab('Version History')).exists();
      assert.dom(SELECTORS.secretTab('Version Diff')).exists();
    });

    test('it hides tabs when no metadata read', async function (assert) {
      assert.expect(5);
      await authPage.login(this.kvNoMetadataRead);
      await visit(`/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
      assert.dom(SELECTORS.secretTab('Secret')).exists();
      assert.dom(SELECTORS.secretTab('Metadata')).exists();
      assert.dom(SELECTORS.secretTab('Version History')).doesNotExist();
      assert.dom(SELECTORS.secretTab('Version Diff')).doesNotExist();
    });
  });
});
