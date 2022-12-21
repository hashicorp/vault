import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import { SELECTORS } from 'vault/tests/helpers/pki/workflow';

const consoleComponent = create(consoleClass);

const tokenWithPolicy = async function (name, policy) {
  await consoleComponent.runCommands([
    `write sys/policies/acl/${name} policy=${btoa(policy)}`,
    `write -field=client_token auth/token/create policies=${name}`,
  ]);
  return consoleComponent.lastLogOutput;
};

const runCommands = async function (commands) {
  try {
    await consoleComponent.runCommands(commands);
    const res = consoleComponent.lastLogOutput;
    if (res.includes('Error')) {
      throw new Error(res);
    }
    return res;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error(
      `The following occurred when trying to run the command(s):\n ${commands.join('\n')} \n\n ${
        consoleComponent.lastLogOutput
      }`
    );
    throw error;
  }
};

/**
 * This test module should test the PKI workflow, including:
 * - link between pages and confirm that the url is as expected
 * - log in as user with a policy and ensure expected UI elements are shown/hidden
 */
module('Acceptance | pki workflow', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    // Setup PKI engine
    const mountPath = `pki-workflow-${new Date().getTime()}`;
    await enablePage.enable('pki', mountPath);
    await runCommands([
      `write ${mountPath}/roles/some-role \
    issuer_ref="default" \
    allowed_domains="example.com" \
    allow_subdomains=true \
    max_ttl="720h"`,
    ]);
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
      const pki_admin_policy = `
          path "${this.mountPath}/*" {
            capabilities = ["create", "read", "update", "delete", "list"]
          },
        `;
      const pki_reader_policy = `
        path "${this.mountPath}/roles" {
          capabilities = ["read", "list"]
        },
        path "${this.mountPath}/roles/*" {
          capabilities = ["read", "list"]
        },
      `;
      const pki_editor_policy = `
        path "${this.mountPath}/roles" {
          capabilities = ["read", "list"]
        },
        path "${this.mountPath}/roles/*" {
          capabilities = ["read", "update"]
        },
        path "${this.mountPath}/issue/*" {
          capabilities = ["update"]
        },
        path "${this.mountPath}/sign/*" {
          capabilities = ["update"]
        },
      `;
      this.pkiRoleReader = await tokenWithPolicy('pki-reader', pki_reader_policy);
      this.pkiRoleEditor = await tokenWithPolicy('pki-editor', pki_editor_policy);
      this.pkiAdminToken = await tokenWithPolicy('pki-admin', pki_admin_policy);
      await logout.visit();
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
      await click(`${SELECTORS.deleteRoleButton} [data-test-confirm-action-trigger]`);
      await click(`[data-test-confirm-button]`);
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
      assert.dom(SELECTORS.pageTitle).hasText('Create a PKI role');

      await fillIn(SELECTORS.roleForm.roleName, roleName);
      await click(SELECTORS.roleForm.roleCreateButton);

      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/${roleName}/details`);
      assert.dom(SELECTORS.breadcrumbs).exists({ count: 4 }, 'Shows 4 breadcrumbs');
      assert.dom(SELECTORS.pageTitle).hasText(`PKI Role ${roleName}`);
    });
  });

  module('issuers', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // Configure engine with a default issuer
      await runCommands([`write ${this.mountPath}/root/generate/internal common_name="Hashicorp Test"`]);
      await logout.visit();
    });
    test('details view renders correct number of info items', async function (assert) {
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
      assert.dom(SELECTORS.issuerDetails.title).hasText('View issuer certificate');
      assert
        .dom(`${SELECTORS.issuerDetails.defaultGroup} ${SELECTORS.issuerDetails.row}`)
        .exists({ count: 9 }, 'Renders 9 info table items under default group');
      assert
        .dom(`${SELECTORS.issuerDetails.urlsGroup} ${SELECTORS.issuerDetails.row}`)
        .exists({ count: 4 }, 'Renders 4 info table items under URLs group');
      assert.dom(SELECTORS.issuerDetails.groupTitle).exists({ count: 1 }, 'only 1 group title rendered');
    });
  });
});
