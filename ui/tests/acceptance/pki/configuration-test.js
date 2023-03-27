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
    await click('[data-test-delete-all-issuers-link]');
    assert.dom('[data-test-modal-background="Delete All Issuers?"]').exists();
    await fillIn('[data-test-delete-all-issuers-input]', 'delete-all');
    await click('[data-test-delete-all-issuers-button]');

    assert
      .dom(SELECTORS.emptyStateTitle)
      .hasText('PKI not configured', `renders correct empty state title after issuers and keys are deleted`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer.",
        `renders correct empty state message after issuers and keys are deleted`
      );
  });

  test('it shows the correct empty state message if roles still exists but no issuers + keys exist', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.emptyStateLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
    await click(SELECTORS.configuration.optionByKey('import'));
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', this.pemBundle);
    await click('[data-test-pki-import-pem-bundle]');
    await runCommands([
      `write ${this.mountPath}/roles/some-role \
      issuer_ref="default" \
      allowed_domains="example.com" \
      allow_subdomains=true \
      max_ttl="720h"`,
    ]);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click('[data-test-delete-all-issuers-link]');
    assert.dom('[data-test-modal-background="Delete All Issuers?"]').exists();
    await fillIn('[data-test-delete-all-issuers-input]', 'delete-all');
    await click('[data-test-delete-all-issuers-button]');

    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing roles. Use the CLI to perform any operations with them until an issuer is configured.",
        'renders correct empty state message when roles still exist, but no issuers or keys exist on overview page'
      );
    await visit(`/vault/secrets/${this.mountPath}/pki/roles`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing roles. Use the CLI to perform any operations with them until an issuer is configured.",
        'renders correct empty state message when roles still exist, but no issuers or keys exist on roles page'
      );
    await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");
    await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");
    await visit(`/vault/secrets/${this.mountPath}/pki/certificates`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");
  });

  test('it shows the correct empty state message if roles and certificates still exists but no issuers + keys exist', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.emptyStateLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
    await click(SELECTORS.configuration.optionByKey('import'));
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', this.pemBundle);
    await click('[data-test-pki-import-pem-bundle]');
    await runCommands([
      `write ${this.mountPath}/roles/some-role \
      issuer_ref="default" \
      allowed_domains="example.com" \
      allow_subdomains=true \
      max_ttl="720h"`,
    ]);
    await runCommands([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click('[data-test-delete-all-issuers-link]');
    assert.dom('[data-test-modal-background="Delete All Issuers?"]').exists();
    await fillIn('[data-test-delete-all-issuers-input]', 'delete-all');
    await click('[data-test-delete-all-issuers-button]');

    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing roles and certificates. Use the CLI to perform any operations with them until an issuer is configured.",
        'renders correct empty state message when roles and certificates still exist, but no issuers or keys exist on overview page'
      );
    await visit(`/vault/secrets/${this.mountPath}/pki/roles`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing roles. Use the CLI to perform any operations with them until an issuer is configured.",
        'renders correct empty state message when roles and certificates still exist, but no issuers or keys exist on roles page'
      );
    await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");
    await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");
    await visit(`/vault/secrets/${this.mountPath}/pki/certificates`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured.",
        'renders correct empty state message when roles and certificates still exist, but no issuers or keys exist on roles page'
      );
  });
  test('it shows the correct empty state message if certificates still exists but no issuers + keys exist', async function (assert) {
    await authPage.login(this.pkiAdminToken);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click(SELECTORS.emptyStateLink);
    assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
    await click(SELECTORS.configuration.optionByKey('import'));
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', this.pemBundle);
    await click('[data-test-pki-import-pem-bundle]');
    await runCommands([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
    await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
    await click('[data-test-delete-all-issuers-link]');
    assert.dom('[data-test-modal-background="Delete All Issuers?"]').exists();
    await fillIn('[data-test-delete-all-issuers-input]', 'delete-all');
    await click('[data-test-delete-all-issuers-button]');

    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

    await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured.",
        'renders correct empty state message when roles and certificates still exist, but no issuers or keys exist on overview page'
      );
    await visit(`/vault/secrets/${this.mountPath}/pki/roles`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer.",
        'renders correct empty state message when roles and certificates still exist, but no issuers or keys exist on roles page'
      );
    await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");
    await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");
    await visit(`/vault/secrets/${this.mountPath}/pki/certificates`);
    assert
      .dom(SELECTORS.emptyStateMessage)
      .hasText(
        "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured.",
        'renders correct empty state message when roles and certificates still exist, but no issuers or keys exist on roles page'
      );
  });
});
