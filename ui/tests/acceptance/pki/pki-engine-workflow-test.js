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
import { adminPolicy, readerPolicy, updatePolicy } from 'vault/tests/helpers/pki/policy-generator';
import { runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { create } from 'ember-cli-page-object';
import flashMessage from 'vault/tests/pages/components/flash-message';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CERTIFICATES, clearRecords } from 'vault/tests/helpers/pki/pki-helpers';
import {
  PKI_CONFIGURE_CREATE,
  PKI_CONFIG_EDIT,
  PKI_ISSUER_DETAILS,
  PKI_ISSUER_LIST,
  PKI_KEYS,
  PKI_ROLE_DETAILS,
} from 'vault/tests/helpers/pki/pki-selectors';

const flash = create(flashMessage);
const { unsupportedPem } = CERTIFICATES;
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
    await runCmd([`delete sys/mounts/${this.mountPath}`]);
  });

  module('not configured', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      const pki_admin_policy = adminPolicy(this.mountPath, 'roles');
      this.pkiAdminToken = await runCmd(tokenWithPolicyCmd(`pki-admin-${this.mountPath}`, pki_admin_policy));
      await logout.visit();
      clearRecords(this.store);
    });

    test('empty state messages are correct when PKI not configured', async function (assert) {
      assert.expect(21);
      const assertEmptyState = (assert, resource) => {
        assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/${resource}`);
        assert
          .dom(GENERAL.emptyStateTitle)
          .hasText(
            'PKI not configured',
            `${resource} index renders correct empty state title when PKI not configured`
          );
        assert.dom(GENERAL.emptyStateActions).hasText('Configure PKI');
        assert
          .dom(GENERAL.emptyStateMessage)
          .hasText(
            `This PKI mount hasn't yet been configured with a certificate issuer.`,
            `${resource} index empty state message correct when PKI not configured`
          );
      };
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);

      await click(GENERAL.secretTab('Roles'));
      assertEmptyState(assert, 'roles');

      await click(GENERAL.secretTab('Issuers'));
      assertEmptyState(assert, 'issuers');

      await click(GENERAL.secretTab('Certificates'));
      assertEmptyState(assert, 'certificates');
      await click(GENERAL.secretTab('Keys'));
      assertEmptyState(assert, 'keys');
      await click(GENERAL.secretTab('Tidy'));
      assertEmptyState(assert, 'tidy');
    });
  });

  module('roles', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // Setup role-specific items
      await runCmd([
        `write ${this.mountPath}/roles/some-role \
      issuer_ref="default" \
      allowed_domains="example.com" \
      allow_subdomains=true \
      max_ttl="720h"`,
      ]);

      await runCmd([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
      const pki_admin_policy = adminPolicy(this.mountPath, 'roles');
      const pki_reader_policy = readerPolicy(this.mountPath, 'roles');
      const pki_editor_policy = updatePolicy(this.mountPath, 'roles');
      this.pkiRoleReader = await runCmd(
        tokenWithPolicyCmd(`pki-reader-${this.mountPath}`, pki_reader_policy)
      );
      this.pkiRoleEditor = await runCmd(
        tokenWithPolicyCmd(`pki-editor-${this.mountPath}`, pki_editor_policy)
      );
      this.pkiAdminToken = await runCmd(tokenWithPolicyCmd(`pki-admin-${this.mountPath}`, pki_admin_policy));
      await logout.visit();
      clearRecords(this.store);
    });

    test('shows correct items if user has all permissions', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(GENERAL.secretTab('Roles')).exists('Roles tab is present');
      await click(GENERAL.secretTab('Roles'));
      assert.dom(PKI_ROLE_DETAILS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);

      assert.dom(PKI_ROLE_DETAILS.generateCertLink).exists('Generate cert link is shown');
      await click(PKI_ROLE_DETAILS.generateCertLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/generate`);

      // Go back to details and test all the links
      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(PKI_ROLE_DETAILS.signCertLink).exists('Sign cert link is shown');
      await click(PKI_ROLE_DETAILS.signCertLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/sign`);

      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(PKI_ROLE_DETAILS.editRoleLink).exists('Edit link is shown');
      await click(PKI_ROLE_DETAILS.editRoleLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/edit`);

      await visit(`/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(PKI_ROLE_DETAILS.deleteRoleButton).exists('Delete role button is shown');
      await click(PKI_ROLE_DETAILS.deleteRoleButton);
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
      assert.dom(GENERAL.secretTab('Roles')).exists('Roles tab is present');
      await click(GENERAL.secretTab('Roles'));
      assert.dom(PKI_ROLE_DETAILS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(PKI_ROLE_DETAILS.deleteRoleButton).doesNotExist('Delete role button is not shown');
      assert.dom(PKI_ROLE_DETAILS.generateCertLink).doesNotExist('Generate cert link is not shown');
      assert.dom(PKI_ROLE_DETAILS.signCertLink).doesNotExist('Sign cert link is not shown');
      assert.dom(PKI_ROLE_DETAILS.editRoleLink).doesNotExist('Edit link is not shown');
    });

    test('it shows correct toolbar items for the user policy', async function (assert) {
      await authPage.login(this.pkiRoleEditor);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(GENERAL.secretTab('Roles')).exists('Roles tab is present');
      await click(GENERAL.secretTab('Roles'));
      assert.dom(PKI_ROLE_DETAILS.createRoleLink).exists({ count: 1 }, 'Create role link is rendered');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One role is in list');
      await click('.linked-block');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`);
      assert.dom(PKI_ROLE_DETAILS.deleteRoleButton).doesNotExist('Delete role button is not shown');
      assert.dom(PKI_ROLE_DETAILS.generateCertLink).exists('Generate cert link is shown');
      assert.dom(PKI_ROLE_DETAILS.signCertLink).exists('Sign cert link is shown');
      assert.dom(PKI_ROLE_DETAILS.editRoleLink).exists('Edit link is shown');
      await click(PKI_ROLE_DETAILS.editRoleLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/edit`,
        'Links to edit view'
      );
      await click(GENERAL.cancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/details`,
        'Cancel from edit goes to details'
      );
      await click(PKI_ROLE_DETAILS.generateCertLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/roles/some-role/generate`,
        'Generate cert button goes to generate page'
      );
      await click(GENERAL.cancelButton);
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
      assert.dom(GENERAL.emptyStateTitle).doesNotExist();
      await click(GENERAL.secretTab('Roles'));
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles`);
      await click(PKI_ROLE_DETAILS.createRoleLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/create`);
      assert.dom(GENERAL.breadcrumbs).exists({ count: 1 }, 'breadcrumbs are rendered');
      assert.dom(GENERAL.breadcrumb).exists({ count: 4 }, 'Shows 4 breadcrumbs');
      assert.dom(GENERAL.title).hasText('Create a PKI Role');

      await fillIn(GENERAL.inputByAttr('name'), roleName);
      await click(GENERAL.saveButton);
      assert.strictEqual(
        flash.latestMessage,
        `Successfully created the role ${roleName}.`,
        'renders success flash upon creation'
      );
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/${roleName}/details`);
      assert.dom(GENERAL.breadcrumb).exists({ count: 4 }, 'Shows 4 breadcrumbs');
      assert.dom(GENERAL.title).hasText(`PKI Role ${roleName}`);
    });
  });

  module('keys', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // base config pki so empty state doesn't show
      await runCmd([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
      const pki_admin_policy = adminPolicy(this.mountPath);
      const pki_reader_policy = readerPolicy(this.mountPath, 'keys', true);
      const pki_editor_policy = updatePolicy(this.mountPath, 'keys');
      this.pkiKeyReader = await runCmd(tokenWithPolicyCmd(`pki-reader-${this.mountPath}`, pki_reader_policy));
      this.pkiKeyEditor = await runCmd(tokenWithPolicyCmd(`pki-editor-${this.mountPath}`, pki_editor_policy));
      this.pkiAdminToken = await runCmd(tokenWithPolicyCmd(`pki-admin-${this.mountPath}`, pki_admin_policy));
      await logout.visit();
      clearRecords(this.store);
    });

    test('shows correct items if user has all permissions', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Keys'));
      // index page
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
      assert
        .dom(PKI_KEYS.importKey)
        .hasAttribute(
          'href',
          `/ui/vault/secrets/${this.mountPath}/pki/keys/import`,
          'import link renders with correct url'
        );
      let keyId = find(PKI_KEYS.keyId).innerText;
      assert.dom('.linked-block').exists({ count: 1 }, 'One key is in list');
      await click('.linked-block');
      // details page
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);
      assert.dom(PKI_KEYS.downloadButton).doesNotExist('does not download button for private key');

      // edit page
      await click(PKI_KEYS.keyEditLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/edit`);
      await click(GENERAL.cancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`,
        'navigates back to details on cancel'
      );
      await visit(`/vault/secrets/${this.mountPath}/pki/keys/${keyId}/edit`);
      await fillIn(GENERAL.inputByAttr('keyName'), 'test-key');
      await click(GENERAL.saveButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`,
        'navigates to details after save'
      );
      assert.dom(GENERAL.infoRowValue('Key name')).hasText('test-key', 'updates key name');

      // key generate and delete navigation
      await visit(`/vault/secrets/${this.mountPath}/pki/keys`);
      await click(PKI_KEYS.generateKey);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/create`);
      await fillIn(GENERAL.inputByAttr('type'), 'exported');
      await fillIn(GENERAL.inputByAttr('keyType'), 'rsa');
      await click(GENERAL.saveButton);
      keyId = find(GENERAL.infoRowValue('Key ID')).textContent?.trim();
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);

      assert
        .dom(PKI_KEYS.nextStepsAlert)
        .hasText(
          'Next steps This private key material will only be available once. Copy or download it now.',
          'renders banner to save private key'
        );
      assert.dom(PKI_KEYS.downloadButton).exists('renders download button');
      await click(PKI_KEYS.keyDeleteButton);
      await click(GENERAL.confirmButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys`,
        'navigates back to key list view on delete'
      );
    });

    test('it hides correct actions for user with read policy', async function (assert) {
      await authPage.login(this.pkiKeyReader);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Keys'));
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
      await isSettled();
      assert.dom(PKI_KEYS.importKey).doesNotExist();
      assert.dom(PKI_KEYS.generateKey).doesNotExist();
      assert.dom('.linked-block').exists({ count: 1 }, 'One key is in list');
      const keyId = find(PKI_KEYS.keyId).innerText;
      await click(GENERAL.menuTrigger);
      assert.dom(PKI_KEYS.popupMenuEdit).doesNotExist('popup menu edit link is not shown');
      await click(PKI_KEYS.popupMenuDetails);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);
      assert.dom(PKI_KEYS.keyDeleteButton).doesNotExist('Delete key button is not shown');
      assert.dom(PKI_KEYS.keyEditLink).doesNotExist('Edit key button does not render');
    });

    test('it shows correct toolbar items for the user with update policy', async function (assert) {
      await authPage.login(this.pkiKeyEditor);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Keys'));
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`);
      await isSettled();
      assert.dom(PKI_KEYS.importKey).exists('import action exists');
      assert.dom(PKI_KEYS.generateKey).exists('generate action exists');
      assert.dom('.linked-block').exists({ count: 1 }, 'One key is in list');
      const keyId = find(PKI_KEYS.keyId).innerText;
      await click(GENERAL.menuTrigger);
      assert.dom(PKI_KEYS.popupMenuEdit).doesNotHaveClass('disabled', 'popup menu edit link is not disabled');
      await click('.linked-block');
      assert.dom(PKI_KEYS.keyDeleteButton).doesNotExist('Delete key button is not shown');
      await click(PKI_KEYS.keyEditLink);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/edit`);
      assert.dom(GENERAL.title).hasText('Edit Key');
      await click(GENERAL.cancelButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys/${keyId}/details`);
    });
  });

  module('issuers', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      const pki_admin_policy = adminPolicy(this.mountPath);
      this.pkiAdminToken = await runCmd(tokenWithPolicyCmd(`pki-admin-${this.mountPath}`, pki_admin_policy));
      // Configure engine with a default issuer
      await runCmd([
        `write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test" name="Hashicorp Test"`,
      ]);
      await logout.visit();
      clearRecords(this.store);
    });
    test('lists the correct issuer metadata info', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(GENERAL.secretTab('Issuers')).exists();
      await click(GENERAL.secretTab('Issuers'));
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
      assert.dom('.linked-block').exists({ count: 1 });
      assert.dom('[data-test-is-root-tag="0"]').hasText('root');
      assert.dom('[data-test-serial-number="0"]').exists({ count: 1 });
      assert.dom('[data-test-common-name="0"]').exists({ count: 1 });
    });
    test('lists the correct issuer metadata info when user has only read permission', async function (assert) {
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Issuers'));
      await click(GENERAL.menuTrigger);
      await click(PKI_ISSUER_LIST.issuerPopupDetails);
      const issuerId = find(PKI_ISSUER_DETAILS.valueByName('Issuer ID')).innerText;
      const pki_issuer_denied_policy = `
      path "${this.mountPath}/*" {
        capabilities = ["create", "read", "update", "delete", "list"]
      },
      path "${this.mountPath}/issuer/${issuerId}" {
        capabilities = ["deny"]
      }
      `;
      this.token = await runCmd(
        tokenWithPolicyCmd(`pki-issuer-denied-policy-${this.mountPath}`, pki_issuer_denied_policy)
      );
      await logout.visit();
      await authPage.login(this.token);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Issuers'));
      assert.dom('[data-test-serial-number="0"]').exists({ count: 1 }, 'displays serial number tag');
      assert.dom('[data-test-common-name="0"]').doesNotExist('does not display cert common name tag');
    });

    test('details view renders correct number of info items', async function (assert) {
      assert.expect(13);
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      assert.dom(GENERAL.secretTab('Issuers')).exists('Issuers tab is present');
      await click(GENERAL.secretTab('Issuers'));
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
      assert.dom('.linked-block').exists({ count: 1 }, 'One issuer is in list');
      await click('.linked-block');
      assert.ok(
        currentURL().match(`/vault/secrets/${this.mountPath}/pki/issuers/.+/details`),
        `/vault/secrets/${this.mountPath}/pki/issuers/my-issuer/details`
      );
      assert.dom(GENERAL.title).hasText('View Issuer Certificate');
      ['Certificate', 'CA Chain', 'Common name', 'Issuer name', 'Issuer ID', 'Default key ID'].forEach(
        (label) => {
          assert
            .dom(`${PKI_ISSUER_DETAILS.defaultGroup} ${PKI_ISSUER_DETAILS.valueByName(label)}`)
            .exists({ count: 1 }, `${label} value rendered`);
        }
      );
      assert
        .dom(`${PKI_ISSUER_DETAILS.urlsGroup} ${PKI_ISSUER_DETAILS.row}`)
        .exists({ count: 3 }, 'Renders 3 info table items under URLs group');
      assert.dom(PKI_ISSUER_DETAILS.groupTitle).exists({ count: 1 }, 'only 1 group title rendered');
    });

    test('toolbar links navigate to expected routes', async function (assert) {
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Issuers'));
      await click(GENERAL.menuTrigger);
      await click(PKI_ISSUER_LIST.issuerPopupDetails);

      const issuerId = find(PKI_ISSUER_DETAILS.valueByName('Issuer ID')).innerText;
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/details`,
        'it navigates to details route'
      );
      assert
        .dom(PKI_ISSUER_DETAILS.crossSign)
        .hasAttribute('href', `/ui/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/cross-sign`);
      assert
        .dom(PKI_ISSUER_DETAILS.signIntermediate)
        .hasAttribute('href', `/ui/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/sign`);
      assert
        .dom(PKI_ISSUER_DETAILS.configure)
        .hasAttribute('href', `/ui/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/edit`);
      await click(PKI_ISSUER_DETAILS.rotateRoot);
      assert.dom(PKI_ISSUER_DETAILS.rotateModal).exists('rotate root modal opens');
      await click(PKI_ISSUER_DETAILS.rotateModalGenerate);
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
      await runCmd([`write ${this.mountPath}/root/generate/internal issuer_name="existing-issuer"`]);
      await logout.visit();
    });
    test('it renders a warning banner when parent issuer has unsupported OIDs', async function (assert) {
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/create`);
      await click(PKI_CONFIGURE_CREATE.optionByKey('import'));
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]', unsupportedPem);
      await click(PKI_CONFIGURE_CREATE.importSubmit);
      const issuerId = find(PKI_CONFIGURE_CREATE.importedIssuer).innerText;
      await click(`${PKI_CONFIGURE_CREATE.importedIssuer} a`);

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
      await click(GENERAL.saveButton);
      assert
        .dom('[data-test-rotate-error]')
        .hasText('Error issuer name already in use', 'it renders error banner');
    });
  });

  module('config', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      await runCmd([`write ${this.mountPath}/root/generate/internal issuer_name="existing-issuer"`]);
      const mixed_config_policy = `
      ${adminPolicy(this.mountPath)}
      ${readerPolicy(this.mountPath, 'config/cluster')}
      `;
      this.mixedConfigCapabilities = await runCmd(
        tokenWithPolicyCmd(`pki-reader-${this.mountPath}`, mixed_config_policy)
      );
      await logout.visit();
    });

    test('it updates config when user only has permission to some endpoints', async function (assert) {
      await authPage.login(this.mixedConfigCapabilities);
      await visit(`/vault/secrets/${this.mountPath}/pki/configuration/edit`);
      assert
        .dom(`${PKI_CONFIG_EDIT.configEditSection} [data-test-component="empty-state"]`)
        .hasText(
          `You do not have permission to set this mount's the cluster config Ask your administrator if you think you should have access to: POST /${this.mountPath}/config/cluster`
        );
      assert.dom(PKI_CONFIG_EDIT.acmeEditSection).exists();
      assert.dom(PKI_CONFIG_EDIT.urlsEditSection).exists();
      assert.dom(PKI_CONFIG_EDIT.crlEditSection).exists();
      assert.dom(`${PKI_CONFIG_EDIT.acmeEditSection} [data-test-component="empty-state"]`).doesNotExist();
      assert.dom(`${PKI_CONFIG_EDIT.urlsEditSection} [data-test-component="empty-state"]`).doesNotExist();
      assert.dom(`${PKI_CONFIG_EDIT.crlEditSection} [data-test-component="empty-state"]`).doesNotExist();
      await click(PKI_CONFIG_EDIT.crlToggleInput('expiry'));
      await click(PKI_CONFIG_EDIT.saveButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/configuration`);
      assert
        .dom('[data-test-value-div="CRL building"]')
        .hasText('Disabled', 'Successfully saves config with partial permission');
    });
  });
});
