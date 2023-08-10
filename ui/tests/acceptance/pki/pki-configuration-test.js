/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, currentURL, fillIn, visit, isSettled, waitUntil, find } from '@ember/test-helpers';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { SELECTORS } from 'vault/tests/helpers/pki/workflow';
import { issuerPemBundle } from 'vault/tests/helpers/pki/values';

module('Acceptance | pki configuration test', function (hooks) {
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
  });

  module('delete all issuers modal and empty states', function (hooks) {
    setupMirage(hooks);

    test('it shows the delete all issuers modal', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(SELECTORS.configuration.configureButton);
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
      await click(SELECTORS.configuration.issuerLink);
      await isSettled();
      assert.dom(SELECTORS.configuration.deleteAllIssuerModal).exists();
      await fillIn(SELECTORS.configuration.deleteAllIssuerInput, 'delete-all');
      await click(SELECTORS.configuration.deleteAllIssuerButton);
      assert.dom(SELECTORS.configuration.deleteAllIssuerModal).doesNotExist();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
    });

    test('it shows the correct empty state message if certificates exists after delete all issuers', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(SELECTORS.configuration.configureButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(SELECTORS.configuration.generateRootOption);
      await fillIn(SELECTORS.configuration.typeField, 'exported');
      await fillIn(SELECTORS.configuration.generateRootCommonNameField, 'issuer-common-0');
      await fillIn(SELECTORS.configuration.generateRootIssuerNameField, 'issuer-0');
      await click(SELECTORS.configuration.generateRootSave);
      await click(SELECTORS.configuration.doneButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.configTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(SELECTORS.configuration.issuerLink);
      assert.dom(SELECTORS.configuration.deleteAllIssuerModal).exists();
      await fillIn(SELECTORS.configuration.deleteAllIssuerInput, 'delete-all');
      await click(SELECTORS.configuration.deleteAllIssuerButton);
      await isSettled();
      assert.dom(SELECTORS.configuration.deleteAllIssuerModal).doesNotExist();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      await isSettled();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await waitUntil(() => currentURL() === `/vault/secrets/${this.mountPath}/pki/overview`);
      await isSettled();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured."
        );

      await visit(`/vault/secrets/${this.mountPath}/pki/roles`);
      await isSettled();
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
      await isSettled();
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
      await isSettled();
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/certificates`);
      await isSettled();
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured."
        );
    });

    test('it shows the correct empty state message if roles and certificates exists after delete all issuers', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(SELECTORS.configuration.configureButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(SELECTORS.configuration.generateRootOption);
      await fillIn(SELECTORS.configuration.typeField, 'exported');
      await fillIn(SELECTORS.configuration.generateRootCommonNameField, 'issuer-common-0');
      await fillIn(SELECTORS.configuration.generateRootIssuerNameField, 'issuer-0');
      await click(SELECTORS.configuration.generateRootSave);
      await click(SELECTORS.configuration.doneButton);
      await runCommands([
        `write ${this.mountPath}/roles/some-role \
        issuer_ref="default" \
        allowed_domains="example.com" \
        allow_subdomains=true \
        max_ttl="720h"`,
      ]);
      await runCommands([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.configTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(SELECTORS.configuration.issuerLink);
      assert.dom(SELECTORS.configuration.deleteAllIssuerModal).exists();
      await fillIn(SELECTORS.configuration.deleteAllIssuerInput, 'delete-all');
      await click(SELECTORS.configuration.deleteAllIssuerButton);
      await isSettled();
      assert.dom(SELECTORS.configuration.deleteAllIssuerModal).doesNotExist();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      await isSettled();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await waitUntil(() => currentURL() === `/vault/secrets/${this.mountPath}/pki/overview`);
      await isSettled();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing roles and certificates. Use the CLI to perform any operations with them until an issuer is configured."
        );

      await visit(`/vault/secrets/${this.mountPath}/pki/roles`);
      await isSettled();
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing roles. Use the CLI to perform any operations with them until an issuer is configured."
        );

      await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
      await isSettled();
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
      await isSettled();
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/certificates`);
      await isSettled();
      assert
        .dom(SELECTORS.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured."
        );
    });

    // test coverage for ed25519 certs not displaying because the verify() function errors
    test('it generates and displays a root issuer of key type = ed25519', async function (assert) {
      assert.expect(4);
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      await click(SELECTORS.generateIssuerDropdown);
      await click(SELECTORS.generateIssuerRoot);
      await fillIn(SELECTORS.configuration.inputByName('type'), 'internal');
      await fillIn(SELECTORS.configuration.inputByName('commonName'), 'my-certificate');
      await click(SELECTORS.configuration.keyParamsGroupToggle);
      await fillIn(SELECTORS.configuration.inputByName('keyType'), 'ed25519');
      await click(SELECTORS.configuration.generateRootSave);

      const issuerId = find(SELECTORS.configuration.saved.issuerLink).innerHTML;
      await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
      assert.dom(SELECTORS.issuerListItem(issuerId)).exists();
      assert
        .dom('[data-test-common-name="0"]')
        .hasText('my-certificate', 'parses certificate metadata in the list view');
      await click(SELECTORS.issuerListItem(issuerId));
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/details`);
      assert.dom(SELECTORS.configuration.saved.commonName).exists('renders issuer details');
    });
  });
});
