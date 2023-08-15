/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import { currentURL, visit } from '@ember/test-helpers';
import { adminPolicy, dataPolicy, metadataPolicy } from 'vault/tests/helpers/policy-generator/kv';
import { deleteEngineCmd, mountEngineCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { writeSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';

/*
This test module tests KV permissions views, each module is is a separate tab (i.e. secret, metadata)
each sub-module is a different state, for example:
- it renders secret details
- it renders secret details after a version is deleted

And each test authenticates using varying permissions testing that view state renders as expected.
*/
// TODO: replace with workflow-* tests
module('Acceptance | kv permissions', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    this.uid = uuidv4();
    // Setup KV engine
    this.mountPath = `kv-engine-${this.uid}`;
    await runCmd(mountEngineCmd('kv-v2', this.mountPath));
    return authPage.logout();
  });

  hooks.afterEach(async function () {
    await authPage.login();
    // Cleanup engine
    await runCmd(deleteEngineCmd(this.mountPath));
  });

  module('secret tab', function (hooks) {
    hooks.beforeEach(async function () {
      // Create secret
      await authPage.login();
      this.secretPath = `my-secret-${this.uid}`;
      await writeSecret(this.mountPath, this.secretPath, 'foo', 'bar');
      // Create different policy test cases
      const kv_admin_policy = adminPolicy(this.mountPath);
      this.kvAdminToken = await runCmd(tokenWithPolicyCmd('kv-admin', kv_admin_policy));

      const no_metadata_read =
        dataPolicy({ backend: this.mountPath, secretPath: this.secretPath }) +
        metadataPolicy({ backend: this.mountPath, capabilities: ['list'] });
      this.cannotReadMetadata = await runCmd(tokenWithPolicyCmd('kv-no-metadata-read', no_metadata_read));

      const no_data_read = dataPolicy({
        backend: this.mountPath,
        secretPath: this.secretPath,
        capabilities: ['list'],
      });
      this.cannotReadData = await runCmd(tokenWithPolicyCmd('kv-no-metadata-read', no_data_read));
      await authPage.logout();
    });

    module('it renders secret details page', function () {
      test('it shows all tabs for admin policy', async function (assert) {
        assert.expect(4);
        await authPage.login(this.kvAdminToken);
        await visit(`/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.dom(PAGE.secretTab('Secret')).exists();
        assert.dom(PAGE.secretTab('Metadata')).exists();
        assert.dom(PAGE.secretTab('Version History')).exists();
      });

      test('it hides tabs when no metadata read', async function (assert) {
        assert.expect(5);
        await authPage.login(this.cannotReadMetadata);
        await visit(`/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.dom(PAGE.secretTab('Secret')).exists();
        assert.dom(PAGE.secretTab('Metadata')).exists();
        assert.dom(PAGE.secretTab('Version History')).doesNotExist();
      });

      test('it shows empty state when cannot read secret data', async function (assert) {
        assert.expect(7);
        await authPage.login(this.cannotReadData);
        await visit(`/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/kv/${this.secretPath}/details`);
        assert.dom(PAGE.secretTab('Secret')).exists();
        assert.dom(PAGE.secretTab('Metadata')).exists();
        assert.dom(PAGE.secretTab('Version History')).doesNotExist();
        assert.dom(PAGE.emptyStateTitle).hasText('You do not have permission to read this secret');
        assert
          .dom(PAGE.emptyStateMessage)
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
