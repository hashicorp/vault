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
import { adminPolicy, dataPolicy, metadataPolicy } from 'vault/tests/helpers/policy-generator/kv';
import { tokenWithPolicy, runCommands, writeSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { SELECTORS } from 'vault/tests/helpers/kv/kv-general-selectors';

/* 
This test module tests KV permissions views, each module is is a separate tab (i.e. secret, metadata)
each sub-module is a different state, for example: 
- it renders secret details
- it renders secret details after a version is deleted

And each test authenticates using varying permissions testing that view state renders as expected.
*/

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

  module('secret tab', function (hooks) {
    hooks.beforeEach(async function () {
      // Create secret
      await authPage.login();
      this.secretPath = `my-secret-${uuidv4()}`;
      await writeSecret(this.mountPath, this.secretPath, 'foo', 'bar');

      // Create different policy test cases
      const kv_admin_policy = adminPolicy(this.mountPath);
      this.kvAdminToken = await tokenWithPolicy('kv-admin', kv_admin_policy);

      const no_metadata_read =
        dataPolicy({ backend: this.mountPath, secretPath: this.secretPath }) +
        metadataPolicy({ backend: this.mountPath, capabilities: ['list'] });
      this.cannotReadMetadata = await tokenWithPolicy('kv-no-metadata-read', no_metadata_read);

      const no_data_read = dataPolicy({
        backend: this.mountPath,
        secretPath: this.secretPath,
        capabilities: ['list'],
      });
      this.cannotReadData = await tokenWithPolicy('kv-no-metadata-read', no_data_read);
      await logout.visit();
    });

    module('it renders secret details page', function () {
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
        await authPage.login(this.cannotReadMetadata);
        await visit(`/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.dom(SELECTORS.secretTab('Secret')).exists();
        assert.dom(SELECTORS.secretTab('Metadata')).exists();
        assert.dom(SELECTORS.secretTab('Version History')).doesNotExist();
        assert.dom(SELECTORS.secretTab('Version Diff')).doesNotExist();
      });

      test('it shows empty state when cannot read secret data', async function (assert) {
        assert.expect(7);
        await authPage.login(this.cannotReadData);
        await visit(`/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.dom(SELECTORS.secretTab('Secret')).exists();
        assert.dom(SELECTORS.secretTab('Metadata')).exists();
        assert.dom(SELECTORS.secretTab('Version History')).doesNotExist();
        assert.dom(SELECTORS.secretTab('Version Diff')).doesNotExist();
        assert.dom(SELECTORS.emptyStateTitle).hasText('You do not have permission to read this secret');
        assert
          .dom(SELECTORS.emptyStateMessage)
          .hasText(
            'Your policies may permit you to write a new version of this secret, but do not allow you to read its current contents.'
          );
      });
    });

    module('it renders secret details page after deleting a version', function () {
      // TODO delete secret and test different policy views
    });
  });
});
