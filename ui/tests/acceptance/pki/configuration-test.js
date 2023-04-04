/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { SELECTORS } from 'vault/tests/helpers/pki/workflow';
import { issuerPemBundle } from 'vault/tests/helpers/pki/values';

module('Acceptance | pki configuration', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.pemBundle = issuerPemBundle;
    await authPage.login();
    // Setup PKI engine
    const mountPath = `pki-workflow-${uuidv4()}`;
    await enablePage.enable('pki', mountPath);
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

  test('it shows the delete all issuers modal', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.emptyStateLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
    await click(SELECTORS.configuration.optionByKey('import'));
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', this.pemBundle);
    await click('[data-test-pki-import-pem-bundle]');
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.configuration.issuerLink);
    assert.dom(SELECTORS.configuration.deleteAllIssuerModal).exists();
    await fillIn(SELECTORS.configuration.deleteAllIssuerInput, 'delete-all');
    await click(SELECTORS.configuration.deleteAllIssuerButton);
    assert.dom('[data-test-component="empty-state"]').exists();
  });
});
