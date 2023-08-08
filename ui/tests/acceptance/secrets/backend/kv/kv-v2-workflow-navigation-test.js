import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import {
  createPolicyCmd,
  deleteEngineCmd,
  mountEngineCmd,
  runCmd,
  createTokenCmd,
  mountAuthCmd,
  tokenWithPolicyCmd,
} from 'vault/tests/helpers/commands';
import {
  adminPolicy,
  dataPolicy,
  metadataListPolicy,
  metadataPolicy,
} from 'vault/tests/helpers/policy-generator/kv';
import {
  addSecretMetadataCmd,
  writeSecret,
  writeVersionedSecret,
} from 'vault/tests/helpers/kv/kv-run-commands';
import { click, currentURL, visit } from '@ember/test-helpers';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';

const secretPath = `my-#:$=?-secret`;
// This doesn't encode in a normal way, so hardcoding it here until we sort that out
const secretPathUrlEncoded = `my-%23:$=%3F-secret`;
const navToBackend = async (backend) => {
  await visit(`/vault/secrets/${backend}/kv/list`);
  // await click(`[data-test-auth-backend-link="${backend}"]`);
  return;
};
const assertCorrectBreadcrumbs = (assert, expected) => {
  assert.dom(PAGE.breadcrumb).exists({ count: expected.length }, 'correct number of breadcrumbs');
  const breadcrumbs = document.querySelectorAll(PAGE.breadcrumb);
  expected.forEach((text, idx) => {
    assert.dom(breadcrumbs[idx]).includesText(text, `position ${idx} breadcrumb includes text ${text}`);
  });
};

/**
 * This test set is for testing the navigation, breadcrumbs, and tabs
 */
module('Acceptance | kv-v2 workflow | navigation', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const uid = uuidv4();
    this.emptyBackend = `kv-empty-${uid}`;
    this.backend = `kv-nav-${uid}`;
    await authPage.login();
    await runCmd(mountEngineCmd('kv-v2', this.emptyBackend), false);
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await writeSecret(this.backend, 'app/nested/secret', 'foo', 'bar');
    await writeVersionedSecret(this.backend, secretPath, 'foo', 'bar', 3);
    await runCmd(addSecretMetadataCmd(this.backend, secretPath, { max_versions: 5, cas_required: true }));
    return;
  });

  hooks.afterEach(async function () {
    await authPage.login();
    await runCmd(deleteEngineCmd(this.backend));
    await runCmd(deleteEngineCmd(this.emptyBackend));
    return;
  });

  module('admin persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        tokenWithPolicyCmd('admin', adminPolicy(this.backend) + adminPolicy(this.emptyBackend))
      );
      await authPage.login(token);
    });
    test('backend nav, tabs, & empty states', async function (assert) {
      assert.expect(24);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // Secrets list page has correct breadcrumbs, toolbar, and contents
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      assert.dom(PAGE.secretTab('list')).hasText('Secrets');
      assert.dom(PAGE.secretTab('list')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Configuration page has correct breadcrumbs, toolbar, and contents
      await click(PAGE.secretTab('Configuration'));
      // when on the list page, Secrets tab selector is `list`, but from config it's `Secrets`
      assert.dom(PAGE.secretTab('Secrets')).doesNotHaveClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasClass('active');
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/configuration`,
        `URL is /vault/secrets/${backend}/kv/configuration`
      );
      // TODO: shows config info & actions

      await click(PAGE.secretTab('Secrets'));
      assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
      assert.dom(PAGE.emptyStateActions).hasText('Create secret');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      await click(`${PAGE.emptyStateActions} a`);
      // TODO: initialKey should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      await click(FORM.cancelBtn);
      // TODO: pageFilter should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );

      await click(PAGE.list.createSecret);
      // TODO: initialKey should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );
    });
    test('navigates to nested secret', async function (assert) {
      assert.expect(29);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['secret', backend]);

      await click(PAGE.list.item('app/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app']);
      assert.dom(PAGE.list.item('nested/')).exists('Shows nested secret');

      await click(PAGE.list.item('nested/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2Fnested%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested']);
      assert.dom(PAGE.list.item('secret')).exists('Shows deeply nested secret');

      await click(PAGE.list.item('secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details?version=1`
      );
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');

      await click(PAGE.breadcrumbAtIdx(3));
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/app%2Fnested%2F/directory`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(2));
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/app%2F/directory`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(1));
      assert.ok(currentURL().startsWith(`/vault/secrets/${backend}/kv/list`), 'links back to list root');
    });
    test('secret nav, tabs, & empty states', async function (assert) {
      assert.expect(55);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=3`,
        'Url includes version query param'
      );
      assert.dom(PAGE.title).hasText(secretPath, 'Goes to secret detail view');
      assert.dom(PAGE.secretTab('Secret')).hasText('Secret');
      assert.dom(PAGE.secretTab('Secret')).hasClass('active');
      assert.dom(PAGE.secretTab('Metadata')).hasText('Metadata');
      assert.dom(PAGE.secretTab('Metadata')).doesNotHaveClass('active');
      assert.dom(PAGE.secretTab('Version History')).hasText('Version History');
      assert.dom(PAGE.secretTab('Version History')).doesNotHaveClass('active');
      assert.dom(PAGE.secretTab('Version Diff')).hasText('Version Diff');
      assert.dom(PAGE.secretTab('Version Diff')).doesNotHaveClass('active');
      assert.dom(PAGE.detail.versionDropdown).hasText('Version 3', 'Version dropdown shows current version');
      assert.dom(PAGE.detail.createNewVersion).hasText('Create new version', 'Create version button shows');
      assert.dom(PAGE.detail.versionCreated).containsText('Version 3 created');
      assert.dom(PAGE.infoRowValue('foo')).exists('renders current data');

      await click(PAGE.detail.createNewVersion);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details/edit?version=3`,
        'Url includes version query param'
      );
      assert.dom(FORM.versionAlert).doesNotExist('Does not show version alert for current version');
      assert.dom(FORM.inputByAttr('path')).isDisabled();

      await click(FORM.cancelBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=3`,
        'Goes back to detail view'
      );

      await click(PAGE.detail.versionDropdown);
      await click(`${PAGE.detail.version(1)} a`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`,
        'Goes to detail view for version 1'
      );
      assert.dom(PAGE.detail.versionDropdown).hasText('Version 1', 'Version dropdown shows selected version');
      assert.dom(PAGE.detail.versionCreated).containsText('Version 1 created');
      assert.dom(PAGE.infoRowValue('key-1')).exists('renders previous data');

      await click(PAGE.detail.createNewVersion);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details/edit?version=1`,
        'Url includes version query param'
      );
      assert.dom(FORM.inputByAttr('path')).isDisabled();
      assert.dom(FORM.keyInput()).hasValue('key-1', 'pre-populates form with selected version data');
      assert.dom(FORM.maskedValueInput()).hasValue('val-1', 'pre-populates form with selected version data');
      assert.dom(FORM.versionAlert).exists('Shows version alert');
      await click(FORM.cancelBtn);

      await click(PAGE.secretTab('Metadata'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/metadata`,
        `goes to metadata page`
      );
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata']);
      assert.dom(PAGE.title).hasText(secretPath);
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('No custom metadata');
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateActions}`)
        .hasText('Add metadata', 'empty state has metadata CTA');
      assert.dom(PAGE.metadata.editBtn).hasText('Edit metadata');

      await click(PAGE.metadata.editBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/metadata/edit`,
        `goes to metadata edit page`
      );
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata', 'edit']);
      await click(FORM.cancelBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/metadata`,
        `cancel btn goes back to metadata page`
      );
    });
    test('breadcrumbs & page titles are correct', async function (assert) {
      assert.expect(38);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'configuration']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      // TODO: config edit button

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);

      await click(PAGE.list.item(secretPath));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath);

      await click(PAGE.detail.createNewVersion);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'edit']);
      assert.dom(PAGE.title).hasText('Create New Version');

      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata']);
      assert.dom(PAGE.title).hasText(secretPath);

      await click(PAGE.metadata.editBtn);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata', 'edit']);
      assert.dom(PAGE.title).hasText('Edit Secret Metadata');

      await click(PAGE.breadcrumbAtIdx(3));
      await click(PAGE.secretTab('Version History'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'version history']);
      // TODO add title to version history page
      // assert.dom(PAGE.title).hasText(secretPath);

      // TODO: fix this selector
      // await click(PAGE.secretTab('Version Diff'));
      // assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'version diff']);
      // assert.dom(PAGE.title).hasText(secretPath);
    });
    test.skip('toolbar actions are correct', async function (assert) {
      assert.expect(0);
    });
    test.skip('can access nested secret', async function (assert) {
      assert.expect(0);
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd('data-reader', dataPolicy({ backend: this.backend, capabilities: ['read'] })),
        createTokenCmd('data-reader'),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'data-reader-list',
          dataPolicy({ backend: this.backend, capabilities: ['read'] }) + metadataListPolicy(this.backend)
        ),
        createTokenCmd('data-reader-list'),
      ]);

      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'metadata-maintainer',
          metadataPolicy({ backend: this.backend }) + metadataListPolicy(this.backend)
        ),
        createTokenCmd('metadata-maintainer'),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('secret-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'secret-creator',
          dataPolicy({ backend: this.backend, capabilities: ['create', 'update'] })
        ),
        createTokenCmd('secret-creator'),
      ]);
      await authPage.login(token);
    });
    // Copy test outline from admin persona
  });

  module('enterprise controlled access persona', function (hooks) {
    hooks.beforeEach(async function () {
      // Set up control group scenario
      const adminUser = 'admin';
      const adminPassword = 'password';
      const userpassMount = 'userpass';
      const POLICY = `
path "${this.backend}/data/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
  control_group = {
    max_ttl = "24h"
    factor "approver" {
      controlled_capabilities = ["write"]
      identity {
          group_names = ["managers"]
          approvals = 1
      }
    }
  }
}

path "${this.backend}/*" {
  capabilities = ["list"]
}

// Can we allow this so user can self-authorize?
path "sys/control-group/authorize" {
  capabilities = ["update"]
}

path "sys/control-group/request" {
  capabilities = ["update"]
}
`;
      const AUTHORIZER_POLICY = `
path "sys/control-group/authorize" {
  capabilities = ["update"]
}

path "sys/control-group/request" {
  capabilities = ["update"]
}
`;
      const userpassAccessor = await runCmd([
        //enable userpass, create user and associated entity
        mountAuthCmd('userpass', userpassMount),
        `write auth/${userpassMount}/users/${adminUser} password=${adminPassword} policies=default`,
        `write identity/entity name=${adminUser} policies=default`,
        // write policies for control group + authorization
        createPolicyCmd('kv-control-group', POLICY),
        createPolicyCmd('authorizer', AUTHORIZER_POLICY),
        // read out mount to get the accessor
        `read -field=accessor sys/internal/ui/mounts/auth/${userpassMount}`,
      ]);
      const authorizerEntityId = await runCmd(`write -field=id identity/lookup/entity name=${adminUser}`);
      const userToken = await runCmd([
        // create alias for authorizor and add them to the managers group
        `write identity/alias mount_accessor=${userpassAccessor} entity_id=${authorizerEntityId} name=${adminUser}`,
        `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=authorizer`,
        // create a token to request access to kv/foo
        'write -field=client_token auth/token/create policies=kv-control-group',
      ]);
      this.userToken = userToken;
      await authPage.login(userToken);
    });
    // Copy test outline from admin persona
  });
});
