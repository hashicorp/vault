/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { click, currentURL, fillIn, find, isSettled, visit } from '@ember/test-helpers';
import { SELECTORS } from 'vault/tests/helpers/pki/workflow';
import { adminPolicy, readerPolicy, updatePolicy } from 'vault/tests/helpers/policy-generator/pki';
import { tokenWithPolicy, runCommands, clearRecords } from 'vault/tests/helpers/pki/pki-run-commands';
import { unsupportedPem } from 'vault/tests/helpers/pki/values';

/**
 * This test module should test the PKI workflow, including:
 * - link between pages and confirm that the url is as expected
 * - log in as user with a policy and ensure expected UI elements are shown/hidden
 */
module('Acceptance | pki workflow', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    await authPage.login();
    // Setup PKI engine
    const mountPath = `pki-workflow-${uuidv4()}`;
    await enablePage.enable('pki', mountPath);
    this.mountPath = mountPath;
    await logout.visit();
    clearRecords(this.store);
  });

  hooks.afterEach(async function () {
    await logout.visit();
    await authPage.login();
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.mountPath}`]);
  });

  module('not configured', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      const pki_admin_policy = adminPolicy(this.mountPath, 'roles');
      this.pkiAdminToken = await tokenWithPolicy(`pki-admin-${this.mountPath}`, pki_admin_policy);
      await logout.visit();
      clearRecords(this.store);
    });

    test('empty state messages are correct when PKI not configured', async function (assert) {
      assert.expect(21);
      const assertEmptyState = (assert, resource) => {
        assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/${resource}`);
        assert
          .dom(SELECTORS.emptyStateTitle)
          .hasText(
            'PKI not configured',
            `${resource} index renders correct empty state title when PKI not configured`
          );
        assert.dom(SELECTORS.emptyStateLink).hasText('Configure PKI');
        assert
          .dom(SELECTORS.emptyStateMessage)
          .hasText(
            `This PKI mount hasn't yet been configured with a certificate issuer.`,
            `${resource} index empty state message correct when PKI not configured`
          );
      };
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);

      await click(SELECTORS.rolesTab);
      assertEmptyState(assert, 'roles');

      await click(SELECTORS.issuersTab);
      assertEmptyState(assert, 'issuers');

      await click(SELECTORS.certsTab);
      assertEmptyState(assert, 'certificates');
      await click(SELECTORS.keysTab);
      assertEmptyState(assert, 'keys');
      await click(SELECTORS.tidyTab);
      assertEmptyState(assert, 'tidy');
    });
  });

  module('roles', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // Setup role-specific items
      await runCommands([
        `write ${this.mountPath}/roles/some-role \
      issuer_ref="default" \
      allowed_domains="example.com" \
      allow_subdomains=true \
      max_ttl="720h"`,
      ]);
      await runCommands([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
      const pki_admin_policy = adminPolicy(this.mountPath, 'roles');
      const pki_reader_policy = readerPolicy(this.mountPath, 'roles');
      const pki_editor_policy = updatePolicy(this.mountPath, 'roles');
      this.pkiRoleReader = await tokenWithPolicy(`pki-reader-${this.mountPath}`, pki_reader_policy);
      this.pkiRoleEditor = await tokenWithPolicy(`pki-editor-${this.mountPath}`, pki_editor_policy);
      this.pkiAdminToken = await tokenWithPolicy(`pki-admin-${this.mountPath}`, pki_admin_policy);
      await logout.visit();
      clearRecords(this.store);
    });

    test('shows correct items if user has all permissions', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.rolesTab).exists('Roles tab is present');
      await click(SELECTORS.rolesTab);
      assert.dom(SELECTORS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);

      assert.dom(SELECTORS.generateCertLink).exists('Generate cert link is shown');
      await click(SELECTORS.generateCertLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/generate`);

      // Go back to details and test all the links
      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.signCertLink).exists('Sign cert link is shown');
      await click(SELECTORS.signCertLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/sign`);

      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.editRoleLink).exists('Edit link is shown');
      await click(SELECTORS.editRoleLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/edit`);

      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.deleteRoleButton).exists('Delete role button is shown');
      await click(SELECTORS.deleteRoleButton);
      await click('[data-test-confirm-button]');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles`,
        'redirects to roles list after deletion'
      );
    });

    test('it does not show toolbar items the user does not have permission to see', async function (assert) {
      await authPage.login(this.pkiRoleReader);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.rolesTab).exists('Roles tab is present');
      await click(SELECTORS.rolesTab);
      assert.dom(SELECTORS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.deleteRoleButton).doesNotExist('Delete role button is not shown');
      assert.dom(SELECTORS.generateCertLink).doesNotExist('Generate cert link is not shown');
      assert.dom(SELECTORS.signCertLink).doesNotExist('Sign cert link is not shown');
      assert.dom(SELECTORS.editRoleLink).doesNotExist('Edit link is not shown');
    });

    test('it shows correct toolbar items for the user policy', async function (assert) {
      await authPage.login(this.pkiRoleEditor);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.rolesTab).exists('Roles tab is present');
      await click(SELECTORS.rolesTab);
      assert.dom(SELECTORS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(SELECTORS.deleteRoleButton).doesNotExist('Delete role button is not shown');
      assert.dom(SELECTORS.generateCertLink).exists('Generate cert link is shown');
      assert.dom(SELECTORS.signCertLink).exists('Sign cert link is shown');
      assert.dom(SELECTORS.editRoleLink).exists('Edit link is shown');
      await click(SELECTORS.editRoleLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/edit`,
        'Links to edit view'
      );
      await click(SELECTORS.roleForm.roleCancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`,
        'Cancel from edit goes to details'
      );
      await click(SELECTORS.generateCertLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/generate`,
        'Generate cert button goes to generate page'
      );
      await click(SELECTORS.generateCertForm.cancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`,
        'Cancel from generate goes to details'
      );
    });

    test('create role happy path', async function (assert) {
      const roleName = 'another-role';
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.rolesTab).exists('Roles tab is present');
      await click(SELECTORS.rolesTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      await click(SELECTORS.createRoleLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/create`);
      assert.dom(SELECTORS.breadcrumbContainer).exists({ count: 1 }, 'breadcrumbs are rendered');
      assert.dom(SELECTORS.breadcrumbs).exists({ count: 4 }, 'Shows 4 breadcrumbs');
      assert.dom(SELECTORS.pageTitle).hasText('Create a PKI Role');

      await fillIn(SELECTORS.roleForm.roleName, roleName);
      await click(SELECTORS.roleForm.roleCreateButton);

      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/${roleName}/details`);
      assert.dom(SELECTORS.breadcrumbs).exists({ count: 4 }, 'Shows 4 breadcrumbs');
      assert.dom(SELECTORS.pageTitle).hasText(`PKI Role ${roleName}`);
    });
  });

  module('keys', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // base config pki so empty state doesn't show
      await runCommands([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
      const pki_admin_policy = adminPolicy(this.mountPath);
      const pki_reader_policy = readerPolicy(this.mountPath, 'keys', true);
      const pki_editor_policy = updatePolicy(this.mountPath, 'keys');
      this.pkiKeyReader = await tokenWithPolicy(`pki-reader-${this.mountPath}`, pki_reader_policy);
      this.pkiKeyEditor = await tokenWithPolicy(`pki-editor-${this.mountPath}`, pki_editor_policy);
      this.pkiAdminToken = await tokenWithPolicy(`pki-admin-${this.mountPath}`, pki_admin_policy);
      await logout.visit();
      clearRecords(this.store);
    });

    test('shows correct items if user has all permissions', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.keysTab);
      // index page
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
      assert
        .dom(SELECTORS.keyPages.importKey)
        .hasAttribute(
          'href',
          `/ui/vault/secrets/${this.mountPath}/pki/keys/import`,
          'import link renders with correct url'
        );
      let keyId = find(SELECTORS.keyPages.keyId).innerText;
      assert.dom('.linked-block').exists({ count: 1 }, 'One key is in list');
      await click('.linked-block');
      // details page
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);
      assert.dom(SELECTORS.keyPages.downloadButton).doesNotExist('does not download button for private key');

      // edit page
      await click(SELECTORS.keyPages.keyEditLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/edit`);
      await click(SELECTORS.keyForm.keyCancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`,
        'navigates back to details on cancel'
      );
      await visit(`/vault/secrets/${this.mountPath}/pki/keys/${keyId}/edit`);
      await fillIn(SELECTORS.keyForm.keyNameInput, 'test-key');
      await click(SELECTORS.keyForm.keyCreateButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`,
        'navigates to details after save'
      );
      assert.dom(SELECTORS.keyPages.keyNameValue).hasText('test-key', 'updates key name');

      // key generate and delete navigation
      await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
      await click(SELECTORS.keyPages.generateKey);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/create`);
      await fillIn(SELECTORS.keyForm.typeInput, 'exported');
      await fillIn(SELECTORS.keyForm.keyTypeInput, 'rsa');
      await click(SELECTORS.keyForm.keyCreateButton);
      keyId = find(SELECTORS.keyPages.keyIdValue).innerText;
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);

      assert
        .dom(SELECTORS.keyPages.nextStepsAlert)
        .hasText(
          'Next steps This private key material will only be available once. Copy or download it now.',
          'renders banner to save private key'
        );
      assert.dom(SELECTORS.keyPages.downloadButton).exists('renders download button');
      await click(SELECTORS.keyPages.keyDeleteButton);
      await click(SELECTORS.keyPages.confirmDelete);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys`,
        'navigates back to key list view on delete'
      );
    });

    test('it hide corrects actions for user with read policy', async function (assert) {
      await authPage.login(this.pkiKeyReader);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.keysTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
      await isSettled();
      assert.dom(SELECTORS.keyPages.importKey).doesNotExist();
      assert.dom(SELECTORS.keyPages.generateKey).doesNotExist();
      assert.dom('.linked-block').exists({ count: 1 }, 'One key is in list');
      const keyId = find(SELECTORS.keyPages.keyId).innerText;
      await click(SELECTORS.keyPages.popupMenuTrigger);
      assert.dom(SELECTORS.keyPages.popupMenuEdit).hasClass('disabled', 'popup menu edit link is disabled');
      await click(SELECTORS.keyPages.popupMenuDetails);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);
      assert.dom(SELECTORS.keyPages.keyDeleteButton).doesNotExist('Delete key button is not shown');
      assert.dom(SELECTORS.keyPages.keyEditLink).doesNotExist('Edit key button does not render');
    });

    test('it shows correct toolbar items for the user with update policy', async function (assert) {
      await authPage.login(this.pkiKeyEditor);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.keysTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
      await isSettled();
      assert.dom(SELECTORS.keyPages.importKey).exists('import action exists');
      assert.dom(SELECTORS.keyPages.generateKey).exists('generate action exists');
      assert.dom('.linked-block').exists({ count: 1 }, 'One key is in list');
      const keyId = find(SELECTORS.keyPages.keyId).innerText;
      await click(SELECTORS.keyPages.popupMenuTrigger);
      assert
        .dom(SELECTORS.keyPages.popupMenuEdit)
        .doesNotHaveClass('disabled', 'popup menu edit link is not disabled');
      await click('.linked-block');
      assert.dom(SELECTORS.keyPages.keyDeleteButton).doesNotExist('Delete key button is not shown');
      await click(SELECTORS.keyPages.keyEditLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/edit`);
      assert.dom(SELECTORS.keyPages.title).hasText('Edit Key');
      await click(SELECTORS.keyForm.keyCancelButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);
    });
  });

  module('issuers', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      const pki_admin_policy = adminPolicy(this.mountPath);
      this.pkiAdminToken = await tokenWithPolicy(`pki-admin-${this.mountPath}`, pki_admin_policy);
      // Configure engine with a default issuer
      await runCommands([
        `write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test" name="Hashicorp Test"`,
      ]);
      await logout.visit();
      clearRecords(this.store);
    });
    test('lists the correct issuer metadata info', async function (assert) {
      assert.expect(6);
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.issuersTab).exists('Issuers tab is present');
      await click(SELECTORS.issuersTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One issuer is in list');
      assert.dom('[data-test-is-root-tag="0"]').hasText('root');
      assert.dom('[data-test-serial-number="0"]').exists({ count: 1 }, 'displays serial number tag');
      assert.dom('[data-test-common-name="0"]').exists({ count: 1 }, 'displays cert common name tag');
    });
    test('lists the correct issuer metadata info when user has only read permission', async function (assert) {
      assert.expect(2);
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      await click(SELECTORS.issuerPopupMenu);
      await click(SELECTORS.issuerPopupDetails);
      const issuerId = find(SELECTORS.issuerDetails.valueByName('Issuer ID')).innerText;
      const pki_issuer_denied_policy = `
      path "${this.mountPath}/*" {
        capabilities = ["create", "read", "update", "delete", "list"]
      },
      path "${this.mountPath}/issuer/${issuerId}" {
        capabilities = ["deny"]
      }
      `;
      this.token = await tokenWithPolicy(
        `pki-issuer-denied-policy-${this.mountPath}`,
        pki_issuer_denied_policy
      );
      await logout.visit();
      await authPage.login(this.token);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      assert.dom('[data-test-serial-number="0"]').exists({ count: 1 }, 'displays serial number tag');
      assert.dom('[data-test-common-name="0"]').doesNotExist('does not display cert common name tag');
    });

    test('details view renders correct number of info items', async function (assert) {
      assert.expect(13);
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(SELECTORS.issuersTab).exists('Issuers tab is present');
      await click(SELECTORS.issuersTab);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One issuer is in list');
      await click('.linked-block');
      assert.ok(
        currentURL().match(`/vault/secrets/${this.mountPath}/pki/issuers/.+/details`),
        `/vault/secrets/${this.mountPath}/pki/issuers/my-issuer/details`
      );
      assert.dom(SELECTORS.issuerDetails.title).hasText('View Issuer Certificate');
      ['Certificate', 'CA Chain', 'Common name', 'Issuer name', 'Issuer ID', 'Default key ID'].forEach(
        (label) => {
          assert
            .dom(`${SELECTORS.issuerDetails.defaultGroup} ${SELECTORS.issuerDetails.valueByName(label)}`)
            .exists({ count: 1 }, `${label} value rendered`);
        }
      );
      assert
        .dom(`${SELECTORS.issuerDetails.urlsGroup} ${SELECTORS.issuerDetails.row}`)
        .exists({ count: 3 }, 'Renders 3 info table items under URLs group');
      assert.dom(SELECTORS.issuerDetails.groupTitle).exists({ count: 1 }, 'only 1 group title rendered');
    });

    test('toolbar links navigate to expected routes', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      await click(SELECTORS.issuerPopupMenu);
      await click(SELECTORS.issuerPopupDetails);

      const issuerId = find(SELECTORS.issuerDetails.valueByName('Issuer ID')).innerText;
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/details`,
        'it navigates to details route'
      );
      assert
        .dom(SELECTORS.issuerDetails.crossSign)
        .hasAttribute('href', `/ui/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/cross-sign`);
      assert
        .dom(SELECTORS.issuerDetails.signIntermediate)
        .hasAttribute('href', `/ui/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/sign`);
      assert
        .dom(SELECTORS.issuerDetails.configure)
        .hasAttribute('href', `/ui/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/edit`);
      await click(SELECTORS.issuerDetails.rotateRoot);
      assert.dom(SELECTORS.issuerDetails.rotateModal).exists('rotate root modal opens');
      await click(SELECTORS.issuerDetails.rotateModalGenerate);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/rotate-root`,
        'it navigates to root rotate form'
      );
      assert
        .dom('[data-test-input="commonName"]')
        .hasValue('Hashicorp Test', 'form prefilled with parent issuer cn');
    });
  });

  module('rotate', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      await runCommands([`write ${this.mountPath}/root/generate/internal issuer_name="existing-issuer"`]);
      await logout.visit();
    });
    test('it renders a warning banner when parent issuer has unsupported OIDs', async function (assert) {
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(SELECTORS.configuration.optionByKey('import'));
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', unsupportedPem);
      await click(SELECTORS.configuration.importSubmit);
      const issuerId = find(SELECTORS.configuration.importedIssuer).innerText;
      await click(`${SELECTORS.configuration.importedIssuer} a`);

      // navigating directly to route because the rotate button is not visible for non-root issuers
      // but we're just testing that route model was parsed and passed as expected
      await visit(`/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/rotate-root`);
      assert
        .dom('[data-test-parsing-warning]')
        .hasTextContaining(
          'Not all of the certificate values can be parsed and transferred to a new root',
          'it renders warning banner'
        );
      assert.dom('[data-test-input="commonName"]').hasValue('fancy-cert-unsupported-subj-and-ext-oids');
      await fillIn('[data-test-input="issuerName"]', 'existing-issuer');
      await click('[data-test-pki-rotate-root-save]');
      assert
        .dom('[data-test-rotate-error]')
        .hasText('Error issuer name already in use', 'it renders error banner');
    });
  });

  module('config', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      await runCommands([`write ${this.mountPath}/root/generate/internal issuer_name="existing-issuer"`]);
      const mixed_config_policy = `
      ${adminPolicy(this.mountPath)}
      ${readerPolicy(this.mountPath, 'config/cluster')}
      `;
      this.mixedConfigCapabilities = await tokenWithPolicy(
        `pki-reader-${this.mountPath}`,
        mixed_config_policy
      );
      await logout.visit();
    });

    test('it updates config when user only has permission to some endpoints', async function (assert) {
      await authPage.login(this.mixedConfigCapabilities);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/edit`);
      assert
        .dom(`${SELECTORS.configEdit.configEditSection} [data-test-component="empty-state"]`)
        .hasText(
          `You do not have permission to set this mount's the cluster config Ask your administrator if you think you should have access to: POST /${this.mountPath}/config/cluster`
        );
      assert.dom(SELECTORS.configEdit.acmeEditSection).exists();
      assert.dom(SELECTORS.configEdit.urlsEditSection).exists();
      assert.dom(SELECTORS.configEdit.crlEditSection).exists();
      assert.dom(`${SELECTORS.acmeEditSection} [data-test-component="empty-state"]`).doesNotExist();
      assert.dom(`${SELECTORS.urlsEditSection} [data-test-component="empty-state"]`).doesNotExist();
      assert.dom(`${SELECTORS.crlEditSection} [data-test-component="empty-state"]`).doesNotExist();
      await click(SELECTORS.configEdit.crlToggleInput('expiry'));
      await click(SELECTORS.configEdit.saveButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      assert
        .dom('[data-test-value-div="CRL building"]')
        .hasText('Disabled', 'Successfully saves config with partial permission');
    });
  });
});
