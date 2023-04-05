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
    await click('[data-test-pki-config-option="generate-root"]');
    await fillIn('[data-test-input="type"]', 'exported');
    await fillIn('[data-test-input="commonName"]', 'issuer-common-0');
    await fillIn('[data-test-input="issuerName"]', 'issuer-0');
    await click('[data-test-pki-generate-root-save]');
    await click('[data-test-done]');
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
    await click(SELECTORS.configTab);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.configuration.issuerLink);
    assert.dom(SELECTORS.configuration.deleteAllIssuerModal).exists();
    await fillIn(SELECTORS.configuration.deleteAllIssuerInput, 'delete-all');
    await click(SELECTORS.configuration.deleteAllIssuerButton);
    assert.dom('[data-test-component="empty-state"]').exists();
  });

  test('it shows the correct empty state message if certificates exists after delete all issuers', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.emptyStateLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
    await click('[data-test-pki-config-option="generate-root"]');
    await fillIn('[data-test-input="type"]', 'exported');
    await fillIn('[data-test-input="commonName"]', 'issuer-common-0');
    await fillIn('[data-test-input="issuerName"]', 'issuer-0');
    await click('[data-test-pki-generate-root-save]');
    await click('[data-test-done]');
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
    await click(SELECTORS.configTab);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.configuration.issuerLink);
    assert.dom(SELECTORS.configuration.deleteAllIssuerModal).exists();
    await fillIn(SELECTORS.configuration.deleteAllIssuerInput, 'delete-all');
    await click(SELECTORS.configuration.deleteAllIssuerButton);
    assert.dom('[data-test-component="empty-state"]').exists();

    await click(SELECTORS.overviewTab);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured."
      );

    await click(SELECTORS.rolesTab);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

    await click(SELECTORS.issuersTab);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

    await click(SELECTORS.keysTab);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

    await click(SELECTORS.certsTab);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured."
      );
  });
});
