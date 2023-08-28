import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentRouteName, currentURL, typeIn, visit, waitUntil } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import {
  createPolicyCmd,
  deleteEngineCmd,
  mountEngineCmd,
  runCmd,
  createTokenCmd,
  tokenWithPolicyCmd,
} from 'vault/tests/helpers/commands';
import { personas } from 'vault/tests/helpers/policy-generator/kv';
import {
  addSecretMetadataCmd,
  clearRecords,
  writeSecret,
  writeVersionedSecret,
} from 'vault/tests/helpers/kv/kv-run-commands';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { setupControlGroup, grantAccess } from 'vault/tests/helpers/control-groups';

const secretPath = `my-#:$=?-secret`;
// This doesn't encode in a normal way, so hardcoding it here until we sort that out
const secretPathUrlEncoded = `my-%23:$=%3F-secret`;
const navToBackend = async (backend) => {
  await visit(`/vault/secrets`);
  return click(PAGE.backends.link(backend));
};
const assertCorrectBreadcrumbs = (assert, expected) => {
  assert.dom(PAGE.breadcrumb).exists({ count: expected.length }, 'correct number of breadcrumbs');
  const breadcrumbs = document.querySelectorAll(PAGE.breadcrumb);
  expected.forEach((text, idx) => {
    assert.dom(breadcrumbs[idx]).includesText(text, `position ${idx} breadcrumb includes text ${text}`);
  });
};
const assertDetailTabs = (assert, current, hidden = []) => {
  const allTabs = ['Secret', 'Metadata', 'Paths', 'Version History'];
  allTabs.forEach((tab) => {
    if (hidden.includes(tab)) {
      assert.dom(PAGE.secretTab(tab)).doesNotExist(`${tab} tab does not render`);
      return;
    }
    assert.dom(PAGE.secretTab(tab)).hasText(tab);
    if (current === tab) {
      assert.dom(PAGE.secretTab(tab)).hasClass('active');
    } else {
      assert.dom(PAGE.secretTab(tab)).doesNotHaveClass('active');
    }
  });
};
const DETAIL_TOOLBARS = ['delete', 'destroy', 'copy', 'versionDropdown', 'createNewVersion'];
const assertDetailsToolbar = (assert, expected = DETAIL_TOOLBARS) => {
  assert
    .dom(PAGE.toolbarAction)
    .exists({ count: expected.length }, 'correct number of toolbar actions render');
  DETAIL_TOOLBARS.forEach((toolbar) => {
    if (expected.includes(toolbar)) {
      assert.dom(PAGE.detail[toolbar]).exists(`${toolbar} toolbar action exists`);
    } else {
      assert.dom(PAGE.detail[toolbar]).doesNotExist(`${toolbar} toolbar action not rendered`);
    }
  });
};

/**
 * This test set is for testing the navigation, breadcrumbs, and tabs.
 * Letter(s) in parenthesis at the end are shorthand for the persona,
 * for ease of tracking down specific tests failures from CI
 */
module('Acceptance | kv-v2 workflow | navigation', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const uid = uuidv4();
    this.store = this.owner.lookup('service:store');
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
        tokenWithPolicyCmd('admin', personas.admin(this.backend) + personas.admin(this.emptyBackend))
      );
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState (a)', async function (assert) {
      assert.expect(18);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // Breadcrumbs correct
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      // Tabs correct
      assert.dom(PAGE.secretTab('Secrets')).hasText('Secrets');
      assert.dom(PAGE.secretTab('Secrets')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar renders');
      assert.dom(PAGE.list.filter).doesNotExist('List filter does not show because no secrets exists.');
      // Page content correct
      assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
      assert.dom(PAGE.emptyStateActions).hasText('Create secret');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // Click empty state CTA
      await click(`${PAGE.emptyStateActions} a`);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret (a)', async function (assert) {
      assert.expect(40);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['secret', backend]);
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Navigate through list items
      await click(PAGE.list.item('app/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('nested/')).exists('Shows nested secret');

      await click(PAGE.list.item('nested/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2Fnested%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/nested/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('secret')).exists('Shows deeply nested secret');

      await click(PAGE.list.item('secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details?version=1`
      );
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert);

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
    test('versioned secret nav, tabs, breadcrumbs (a)', async function (assert) {
      assert.expect(45);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=3`,
        'Url includes version query param'
      );
      assert.dom(PAGE.title).hasText(secretPath, 'title is correct on detail view');
      assertDetailTabs(assert, 'Secret');
      assert.dom(PAGE.detail.versionDropdown).hasText('Version 3', 'Version dropdown shows current version');
      assert.dom(PAGE.detail.createNewVersion).hasText('Create new version', 'Create version button shows');
      assert.dom(PAGE.detail.versionTimestamp).containsText('Version 3 created');
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
      assert.dom(PAGE.detail.versionTimestamp).containsText('Version 1 created');
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
    test('breadcrumbs & page titles are correct (a)', async function (assert) {
      assert.expect(45);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'configuration']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for secret list');

      await click(PAGE.list.item(secretPath));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      await click(PAGE.detail.createNewVersion);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'edit']);
      assert.dom(PAGE.title).hasText('Create New Version', 'correct page title for secret edit');

      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      await click(PAGE.metadata.editBtn);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata', 'edit']);
      assert.dom(PAGE.title).hasText('Edit Secret Metadata', 'correct page title for metadata edit');

      await click(PAGE.breadcrumbAtIdx(3));
      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      await click(PAGE.secretTab('Version History'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'version history']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for version history');
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'data-reader',
          personas.dataReader(this.backend) + personas.dataReader(this.emptyBackend)
        ),
        createTokenCmd('data-reader'),
      ]);
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState (dr)', async function (assert) {
      assert.expect(15);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // Breadcrumbs correct
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      // Tabs correct
      assert.dom(PAGE.secretTab('Secrets')).hasText('Secrets');
      assert.dom(PAGE.secretTab('Secrets')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar renders');
      assert
        .dom(PAGE.list.filter)
        .doesNotExist('list filter input does not render because no list capabilities');
      // Page content correct
      assert
        .dom(PAGE.emptyStateTitle)
        .doesNotExist('empty state does not render because no metadata access to list');
      assert.dom(PAGE.list.overviewCard).exists('renders overview card');

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create`,
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list`,
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret (dr)', async function (assert) {
      assert.expect(23);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['secret', backend]);
      assert
        .dom(PAGE.list.filter)
        .doesNotExist('List filter input does not render because no list capabilities');

      await typeIn(PAGE.list.overviewInput, 'app/nested/secret');
      await click(PAGE.list.overviewButton);

      // Goes to correct detail view
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details?version=1`
      );
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['copy']);

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
    test('versioned secret nav, tabs, breadcrumbs (dr)', async function (assert) {
      assert.expect(28);
      const backend = this.backend;
      await navToBackend(backend);

      // Navigate to secret
      await typeIn(PAGE.list.overviewInput, secretPath);
      await click(PAGE.list.overviewButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=3`,
        'Url includes version query param'
      );
      assert.dom(PAGE.title).hasText(secretPath, 'Goes to secret detail view');
      assertDetailTabs(assert, 'Secret', ['Version History']);
      assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown hidden');
      assert.dom(PAGE.detail.createNewVersion).doesNotExist('unable to create a new version');
      assert.dom(PAGE.detail.versionTimestamp).containsText('Version 3 created');
      assert.dom(PAGE.infoRowValue('foo')).exists('renders current data');

      // data-reader can't navigate to older versions, but they can go to page directly
      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`);
      assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown does not exist');
      assert.dom(PAGE.detail.versionTimestamp).containsText('Version 1 created');
      assert.dom(PAGE.infoRowValue('key-1')).exists('renders previous data');

      assert.dom(PAGE.detail.createNewVersion).doesNotExist('cannot create new version');

      await click(PAGE.secretTab('Metadata'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/metadata`,
        `goes to metadata page`
      );
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata']);
      assert.dom(PAGE.title).hasText(secretPath);
      assert.dom(PAGE.toolbarAction).doesNotExist('no toolbar actions available on metadata');
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('No custom metadata');
      assert
        .dom(`${PAGE.metadata.secretMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('You do not have access to secret metadata');
      assert.dom(PAGE.metadata.editBtn).doesNotExist('edit button hidden');
    });
    test('breadcrumbs & page titles are correct (dr)', async function (assert) {
      assert.expect(35);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'configuration']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title correct on config page');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title correct on secrets list');

      await typeIn(PAGE.list.overviewInput, 'app/nested/secret');
      await click(PAGE.list.overviewButton);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title correct on secret detail');

      assert.dom(PAGE.detail.createNewVersion).doesNotExist('cannot create new version');

      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'app', 'nested', 'secret', 'metadata']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title correct on metadata');

      assert.dom(PAGE.metadata.editBtn).doesNotExist('cannot edit metadata');

      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'app', 'nested', 'secret', 'paths']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'correct page title for paths');

      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab not shown');
    });
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'data-reader-list',
          personas.dataListReader(this.backend) + personas.dataListReader(this.emptyBackend)
        ),
        createTokenCmd('data-reader-list'),
      ]);

      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState (dlr)', async function (assert) {
      assert.expect(18);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // Breadcrumbs correct
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      // Tabs correct
      assert.dom(PAGE.secretTab('Secrets')).hasText('Secrets');
      assert.dom(PAGE.secretTab('Secrets')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar renders');
      assert.dom(PAGE.list.filter).doesNotExist('List filter does not show because no secrets exists.');
      // Page content correct
      assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
      assert.dom(PAGE.emptyStateActions).hasText('Create secret');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // Click empty state CTA
      await click(`${PAGE.emptyStateActions} a`);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret (dlr)', async function (assert) {
      assert.expect(31);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['secret', backend]);
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Navigate through list items
      await click(PAGE.list.item('app/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      assert.dom(PAGE.list.filter).doesNotExist('List filter hidden since no nested list access');

      assert
        .dom(PAGE.list.overviewInput)
        .hasValue('app/', 'overview card is pre-filled with directory param');
      await typeIn(PAGE.list.overviewInput, 'nested/secret');
      await click(PAGE.list.overviewButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details?version=1`
      );
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['delete', 'copy']);

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
    test('versioned secret nav, tabs, breadcrumbs (dlr)', async function (assert) {
      assert.expect(28);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=3`,
        'Url includes version query param'
      );
      assert.dom(PAGE.title).hasText(secretPath, 'Goes to secret detail view');
      assertDetailTabs(assert, 'Secret', ['Version History']);
      assert.dom(PAGE.detail.versionDropdown).doesNotExist('does not show version dropdown');
      assert.dom(PAGE.detail.createNewVersion).doesNotExist('unable to create a new version');
      assert.dom(PAGE.detail.versionTimestamp).containsText('Version 3 created');
      assert.dom(PAGE.infoRowValue('foo')).exists('renders current data');

      assert.dom(PAGE.detail.createNewVersion).doesNotExist('cannot create new version');

      // data-list-reader can't navigate to older versions, but they can go to page directly
      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`);
      assert.dom(PAGE.detail.versionDropdown).doesNotExist('no version dropdown');
      assert.dom(PAGE.detail.versionTimestamp).containsText('Version 1 created');
      assert.dom(PAGE.infoRowValue('key-1')).exists('renders previous data');

      assert.dom(PAGE.detail.createNewVersion).doesNotExist('cannot create new version from old version');

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
        .dom(`${PAGE.metadata.secretMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('You do not have access to secret metadata');
      assert.dom(PAGE.metadata.editBtn).doesNotExist('edit button hidden');
    });
    test('breadcrumbs & page titles are correct (dlr)', async function (assert) {
      assert.expect(29);
      const backend = this.backend;
      await navToBackend(backend);

      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'configuration']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for secret list');

      await click(PAGE.list.item(secretPath));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      assert.dom(PAGE.detail.createNewVersion).doesNotExist('cannot create new version');

      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      assert.dom(PAGE.metadata.editBtn).doesNotExist('cannot edit metadata');

      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab not shown');
    });
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'metadata-maintainer',
          personas.metadataMaintainer(this.backend) + personas.metadataMaintainer(this.emptyBackend)
        ),
        createTokenCmd('metadata-maintainer'),
      ]);
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState (mm)', async function (assert) {
      assert.expect(18);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // Breadcrumbs correct
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      // Tabs correct
      assert.dom(PAGE.secretTab('Secrets')).hasText('Secrets');
      assert.dom(PAGE.secretTab('Secrets')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar only renders create secret action');
      assert.dom(PAGE.list.filter).doesNotExist('List filter does not show because no secrets exists.');
      // Page content correct
      assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
      assert.dom(PAGE.emptyStateActions).hasText('Create secret');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // Click empty state CTA
      await click(`${PAGE.emptyStateActions} a`);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret (mm)', async function (assert) {
      assert.expect(41);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['secret', backend]);
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Navigate through list items
      await click(PAGE.list.item('app/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('nested/')).exists('Shows nested secret');

      await click(PAGE.list.item('nested/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2Fnested%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/nested/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('secret')).exists('Shows deeply nested secret');

      await click(PAGE.list.item('secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details`,
        `Goes to URL with version`
      );
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['delete', 'destroy', 'versionDropdown']);
      assert.dom(PAGE.detail.versionDropdown).hasText('Version 1', 'Shows version timestamp');

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
    test('versioned secret nav, tabs, breadcrumbs (mm)', async function (assert) {
      assert.expect(37);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details`,
        'Url includes version query param'
      );
      assert.dom(PAGE.title).hasText(secretPath, 'Goes to secret detail view');
      assertDetailTabs(assert, 'Secret');
      assert.dom(PAGE.detail.versionDropdown).hasText('Version 3', 'Version dropdown shows current version');
      assert.dom(PAGE.detail.createNewVersion).doesNotExist('Create new version button not shown');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('Version created text not shown');
      assert.dom(PAGE.infoRowValue('foo')).doesNotExist('does not render current data');
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'Shows empty state on secret detail');

      await click(PAGE.detail.versionDropdown);
      await click(`${PAGE.detail.version(1)} a`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`,
        'Goes to detail view for version 1'
      );
      assert.dom(PAGE.detail.versionDropdown).hasText('Version 1', 'Version dropdown shows selected version');

      assert.dom(PAGE.infoRowValue('key-1')).doesNotExist('does not render previous data');
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText(
          'You do not have permission to read this secret',
          'Shows empty state on secret detail for older version'
        );

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
    test('breadcrumbs & page titles are correct (mm)', async function (assert) {
      assert.expect(39);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'configuration']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for secret list');

      await click(PAGE.list.item(secretPath));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      await click(PAGE.metadata.editBtn);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata', 'edit']);
      assert.dom(PAGE.title).hasText('Edit Secret Metadata', 'correct page title for metadata edit');

      await click(PAGE.breadcrumbAtIdx(3));
      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      await click(PAGE.secretTab('Version History'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'version history']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for version history');
    });
  });

  module('secret-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          'secret-creator',
          personas.secretCreator(this.backend) + personas.secretCreator(this.emptyBackend)
        ),
        createTokenCmd('secret-creator'),
      ]);
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState (sc)', async function (assert) {
      assert.expect(15);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // Breadcrumbs correct
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      // Tabs correct
      assert.dom(PAGE.secretTab('Secrets')).hasText('Secrets');
      assert.dom(PAGE.secretTab('Secrets')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar only renders create secret action');
      assert.dom(PAGE.list.filter).doesNotExist('List filter input is not rendered');
      // Page content correct
      assert.dom(PAGE.list.overviewCard).exists('Overview card renders');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create`,
        `goes to /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list`,
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret (sc)', async function (assert) {
      assert.expect(23);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['secret', backend]);
      assert.dom(PAGE.list.filter).doesNotExist('List filter input is not rendered');

      // Navigate to secret
      await typeIn(PAGE.list.overviewInput, 'app/nested/secret');
      await click(PAGE.list.overviewButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details`,
        'goes to secret detail page'
      );
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['createNewVersion']);

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
    test('versioned secret nav, tabs, breadcrumbs (sc)', async function (assert) {
      assert.expect(36);
      const backend = this.backend;
      await navToBackend(backend);

      await typeIn(PAGE.list.overviewInput, secretPath);
      await click(PAGE.list.overviewButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details`,
        'Goes to detail view'
      );
      assert.dom(PAGE.title).hasText(secretPath, 'Goes to secret detail view');
      assertDetailTabs(assert, 'Secret', ['Version History']);
      assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown does not render');
      assert.dom(PAGE.detail.createNewVersion).hasText('Create new version', 'Create version button shows');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('Version created info is not rendered');
      assert.dom(PAGE.infoRowValue('foo')).doesNotExist('current data not rendered');
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'empty state shows');

      await click(PAGE.detail.createNewVersion);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details/edit`,
        'Goes to edit page'
      );
      assert.dom(FORM.versionAlert).doesNotExist('Does not show version alert for current version');
      assert
        .dom(FORM.noReadAlert)
        .hasText(
          'Warning You do not have read permissions for this secret data. Saving will overwrite the existing secret.',
          'Shows warning about no read permissions'
        );
      assert.dom(FORM.inputByAttr('path')).isDisabled();

      await click(FORM.cancelBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details`,
        'Goes back to detail view'
      );

      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`);
      assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown does not exist');
      assert.dom(PAGE.detail.versionTimestamp).doesNotExist('version created data not rendered');
      assert.dom(PAGE.infoRowValue('key-1')).doesNotExist('does not render previous data');

      await click(PAGE.detail.createNewVersion);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details/edit?version=1`,
        'Url includes version query param'
      );
      assert.dom(FORM.inputByAttr('path')).isDisabled();
      assert.dom(FORM.keyInput()).hasValue('', 'form does not pre-populate');
      assert.dom(FORM.maskedValueInput()).hasValue('', 'form does not pre-populate');
      assert.dom(FORM.noReadAlert).exists('Shows no read alert');
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
        .hasText('You do not have access to read custom metadata', 'shows correct empty state');
      assert.dom(PAGE.metadata.editBtn).doesNotExist('edit metadata button does not render');
    });
    test('breadcrumbs & page titles are correct (sc)', async function (assert) {
      assert.expect(34);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'configuration']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for secret list');

      await typeIn(PAGE.list.overviewInput, secretPath);
      await click(PAGE.list.overviewButton);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      await click(PAGE.detail.createNewVersion);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'edit']);
      assert.dom(PAGE.title).hasText('Create New Version', 'correct page title for secret edit');

      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      assert.dom(PAGE.metadata.editBtn).doesNotExist('cannot edit metadata');

      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab not shown');
    });
  });

  module('enterprise controlled access persona', function (hooks) {
    hooks.beforeEach(async function () {
      // Set up control group scenario
      const userPolicy = `
path "${this.backend}/data/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
  control_group = {
    max_ttl = "24h"
    factor "ops_manager" {
      controlled_capabilities = ["read"]
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
`;
      const { userToken } = await setupControlGroup({ userPolicy });
      this.userToken = userToken;
      await authPage.login(userToken);
      clearRecords(this.store);
      return;
    });
    test('can access nested secret (cg)', async function (assert) {
      assert.expect(42);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['secret', backend]);
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Navigate through list items
      await click(PAGE.list.item('app/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('nested/')).exists('Shows nested secret');

      await click(PAGE.list.item('nested/'));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/app%2Fnested%2F/directory`);
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/nested/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('secret')).exists('Shows deeply nested secret');

      // For some reason when we click on the item in tests it throws a global control group error
      // But not when we visit the page directly
      await visit(`/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details`);
      assert.ok(
        await waitUntil(() => currentRouteName() === 'vault.cluster.access.control-group-accessor'),
        'redirects to access control group route'
      );
      await grantAccess({
        apiPath: `${backend}/data/app/nested/secret`,
        originUrl: `/vault/secrets/${backend}/kv/app%2Fnested%2F/directory`,
        userToken: this.userToken,
      });
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2F/directory`,
        'navigates to list url where secret is'
      );
      await click(PAGE.list.item('secret'));

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details?version=1`,
        'goes to secret details'
      );
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['delete', 'copy', 'createNewVersion']);

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
    test('breadcrumbs & page titles are correct (cg)', async function (assert) {
      assert.expect(36);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'configuration']);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} Version 2`, 'correct page title for secret list');

      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details`);

      assert.ok(
        await waitUntil(() => currentRouteName() === 'vault.cluster.access.control-group-accessor'),
        'redirects to access control group route'
      );

      await grantAccess({
        apiPath: `${backend}/data/${encodeURIComponent(secretPath)}`,
        originUrl: `/vault/secrets/${backend}/kv/list`,
        userToken: this.userToken,
      });

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list`,
        'navigates back to list url after authorized'
      );
      await click(PAGE.list.item(secretPath));

      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      assert.dom(PAGE.metadata.editBtn).doesNotExist('cannot edit metadata');

      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab not shown');

      await click(PAGE.secretTab('Secret'));
      await click(PAGE.detail.createNewVersion);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'edit']);
      assert.dom(PAGE.title).hasText('Create New Version', 'correct page title for secret edit');
    });
  });
});
