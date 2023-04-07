/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, fillIn, visit, isSettled } from '@ember/test-helpers';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { SELECTORS } from 'vault/tests/helpers/pki/workflow';
import { issuerPemBundle } from 'vault/tests/helpers/pki/values';

module('Acceptance | pki tidy form test', function (hooks) {
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

  test('it navigates to the tidy page from configuration toolbar', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.emptyStateLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
    await isSettled();
    await click(SELECTORS.configuration.generateRootOption);
    await fillIn(SELECTORS.configuration.typeField, 'exported');
    await fillIn(SELECTORS.configuration.generateRootCommonNameField, 'issuer-common-0');
    await fillIn(SELECTORS.configuration.generateRootIssuerNameField, 'issuer-0');
    await click(SELECTORS.configuration.generateRootSave);
    await click(SELECTORS.configuration.doneButton);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
    await isSettled();
    await click(SELECTORS.configTab);
    await isSettled();
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.configuration.tidyToolbar);
    await isSettled();
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/tidy`);
  });

  test('it returns to the configuration page after submit', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.emptyStateLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
    await isSettled();
    await click(SELECTORS.configuration.generateRootOption);
    await fillIn(SELECTORS.configuration.typeField, 'exported');
    await fillIn(SELECTORS.configuration.generateRootCommonNameField, 'issuer-common-0');
    await fillIn(SELECTORS.configuration.generateRootIssuerNameField, 'issuer-0');
    await click(SELECTORS.configuration.generateRootSave);
    await click(SELECTORS.configuration.doneButton);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
    await isSettled();
    await click(SELECTORS.configTab);
    await isSettled();
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.configuration.tidyToolbar);
    await isSettled();
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/tidy`);
    await click(SELECTORS.configuration.tidyCertStoreCheckbox);
    await click(SELECTORS.configuration.tidyRevocationCheckbox);
    await fillIn(SELECTORS.configuration.safetyBufferInput, '100');
    await fillIn(SELECTORS.configuration.safetyBufferInputDropdown, 'd');
    await click(SELECTORS.configuration.tidySave);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
  });
});
