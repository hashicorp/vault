/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-settled-after-test-helper */
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, currentURL, fillIn, visit, settled, find, waitFor, waitUntil } from '@ember/test-helpers';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import {
  PKI_CONFIGURE_CREATE,
  PKI_DELETE_ALL_ISSUERS,
  PKI_GENERATE_ROOT,
  PKI_ISSUER_LIST,
} from 'vault/tests/helpers/pki/pki-selectors';

const { issuerPemBundle } = CERTIFICATES;
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
    await runCmd([`delete sys/mounts/${this.mountPath}`]);
  });

  module('delete all issuers modal and empty states', function (hooks) {
    setupMirage(hooks);

    test('it shows the delete all issuers modal', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(PKI_CONFIGURE_CREATE.configureButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await settled();
      await click(PKI_CONFIGURE_CREATE.generateRootOption);
      await fillIn(GENERAL.inputByAttr('type'), 'exported');
      await fillIn(GENERAL.inputByAttr('commonName'), 'issuer-common-0');
      await fillIn(GENERAL.inputByAttr('issuerName'), 'issuer-0');
      await click(GENERAL.saveButton);
      await click(PKI_CONFIGURE_CREATE.doneButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await settled();
      await click(GENERAL.secretTab('Configuration'));
      await settled();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(PKI_DELETE_ALL_ISSUERS.issuerLink);
      await settled();
      await waitFor(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal, { timeout: 5000 });
      assert.dom(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal).exists();
      await fillIn(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerInput, 'delete-all');
      await click(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerButton);
      await settled();
      await waitUntil(() => !find(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal));

      assert.dom(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal).doesNotExist();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
    });

    test('it shows the correct empty state message if certificates exists after delete all issuers', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(PKI_CONFIGURE_CREATE.configureButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration/create`,
        'goes to pki configure page'
      );
      await click(PKI_CONFIGURE_CREATE.generateRootOption);
      await fillIn(GENERAL.inputByAttr('type'), 'exported');
      await fillIn(GENERAL.inputByAttr('commonName'), 'issuer-common-0');
      await fillIn(GENERAL.inputByAttr('issuerName'), 'issuer-0');
      await click(GENERAL.saveButton);
      await click(PKI_CONFIGURE_CREATE.doneButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/overview`,
        'goes to overview page'
      );
      await click(GENERAL.secretTab('Configuration'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration`,
        'goes to configuration page'
      );
      await click(PKI_DELETE_ALL_ISSUERS.issuerLink);
      await waitFor(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal);
      assert.dom(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal).exists();
      await fillIn(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerInput, 'delete-all');
      await click(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerButton);
      await waitUntil(() => !find(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal));
      assert.dom(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal).doesNotExist('delete all issuers modal closes');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/configuration`,
        'is still on configuration page'
      );
      await settled();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await settled();
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/overview`,
        'goes to overview page'
      );
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured."
        );

      await visit(`/vault/secrets/${this.mountPath}/pki/roles`);
      await settled();
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
      await settled();
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
      await settled();
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/certificates`);
      await settled();
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured."
        );
    });

    test('it shows the correct empty state message if roles and certificates exists after delete all issuers', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      // Configure PKI
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(PKI_CONFIGURE_CREATE.configureButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(PKI_CONFIGURE_CREATE.generateRootOption);
      await fillIn(GENERAL.inputByAttr('type'), 'exported');
      await fillIn(GENERAL.inputByAttr('commonName'), 'issuer-common-0');
      await fillIn(GENERAL.inputByAttr('issuerName'), 'issuer-0');
      await click(GENERAL.saveButton);
      await click(PKI_CONFIGURE_CREATE.doneButton);
      // Create role and root CA"
      await runCmd([
        `write ${this.mountPath}/roles/some-role \
        issuer_ref="default" \
        allowed_domains="example.com" \
        allow_subdomains=true \
        max_ttl="720h"`,
      ]);
      await runCmd([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Configuration'));
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      await click(PKI_DELETE_ALL_ISSUERS.issuerLink);
      await waitFor(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal);
      assert.dom(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal).exists();
      await fillIn(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerInput, 'delete-all');
      await click(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerButton);
      await settled();
      await waitUntil(() => !find(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal));
      assert.dom(PKI_DELETE_ALL_ISSUERS.deleteAllIssuerModal).doesNotExist();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      await settled();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await settled();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing roles and certificates. Use the CLI to perform any operations with them until an issuer is configured."
        );

      await visit(`/vault/secrets/${this.mountPath}/pki/roles`);
      await settled();
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing roles. Use the CLI to perform any operations with them until an issuer is configured."
        );

      await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
      await settled();
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
      await settled();
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText("This PKI mount hasn't yet been configured with a certificate issuer.");

      await visit(`/vault/secrets/${this.mountPath}/pki/certificates`);
      await settled();
      assert
        .dom(GENERAL.emptyStateMessage)
        .hasText(
          "This PKI mount hasn't yet been configured with a certificate issuer. There are existing certificates. Use the CLI to perform any operations with them until an issuer is configured."
        );
    });

    // test coverage for ed25519 certs not displaying because the verify() function errors
    test('it generates and displays a root issuer of key type = ed25519', async function (assert) {
      assert.expect(4);
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Issuers'));
      await click(PKI_ISSUER_LIST.generateIssuerDropdown);
      await click(PKI_ISSUER_LIST.generateIssuerRoot);
      await fillIn(GENERAL.inputByAttr('type'), 'internal');
      await fillIn(GENERAL.inputByAttr('commonName'), 'my-certificate');
      await click(PKI_GENERATE_ROOT.keyParamsGroupToggle);
      await fillIn(GENERAL.inputByAttr('keyType'), 'ed25519');
      await click(GENERAL.saveButton);

      const issuerId = find(PKI_GENERATE_ROOT.saved.issuerLink).innerHTML;
      await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);
      assert.dom(PKI_ISSUER_LIST.issuerListItem(issuerId)).exists();
      assert
        .dom('[data-test-common-name="0"]')
        .hasText('my-certificate', 'parses certificate metadata in the list view');
      await click(PKI_ISSUER_LIST.issuerListItem(issuerId));
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/details`);
      assert.dom(PKI_GENERATE_ROOT.saved.commonName).exists('renders issuer details');
    });
  });
});
