import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentRouteName, currentURL, typeIn, visit, waitUntil } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
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
  setupControlGroup,
  writeSecret,
  writeVersionedSecret,
} from 'vault/tests/helpers/kv/kv-run-commands';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import controlGroup from 'vault/tests/pages/components/control-group';
import { CONTROL_GROUP_PREFIX, TOKEN_SEPARATOR } from 'vault/services/control-group';

const controlGroupComponent = create(controlGroup);

const secretPath = `my-#:$=?-secret`;
// This doesn't encode in a normal way, so hardcoding it here until we sort that out
const secretPathUrlEncoded = `my-%23:$=%3F-secret`;
const navToBackend = (backend) => {
  return visit(`/vault/secrets/${backend}/kv/list`);
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
        tokenWithPolicyCmd('admin', personas.admin(this.backend) + personas.admin(this.emptyBackend))
      );
      await authPage.login(token);
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState', async function (assert) {
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
      assert.dom(PAGE.secretTab('list')).hasText('Secrets');
      assert.dom(PAGE.secretTab('list')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar renders');
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');
      // Page content correct
      assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
      assert.dom(PAGE.emptyStateActions).hasText('Create secret');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // Click empty state CTA
      await click(`${PAGE.emptyStateActions} a`);
      // TODO: initialKey should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      // TODO: pageFilter should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      // TODO: initialKey should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      // TODO: pageFilter should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret', async function (assert) {
      assert.expect(36);
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
      assert.dom(PAGE.toolbar).exists('toolbar renders');
      assert.dom(PAGE.toolbarAction).exists({ count: 2 }, 'correct number of toolbar actions render');

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
    test('versioned secret nav, tabs, breadcrumbs', async function (assert) {
      assert.expect(43);
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
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState', async function (assert) {
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
      assert.dom(PAGE.secretTab('list')).hasText('Secrets');
      assert.dom(PAGE.secretTab('list')).hasClass('active');
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
      // TODO: initialKey should not show on query params if empty
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=`,
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      // TODO: pageFilter should not show on query params if empty
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list?pageFilter=`,
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret', async function (assert) {
      assert.expect(19);
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
      assert.dom(PAGE.toolbar).exists('toolbar renders');
      assert.dom(PAGE.toolbarAction).exists({ count: 2 }, 'correct number of toolbar actions render');

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
    test('versioned secret nav, tabs, breadcrumbs', async function (assert) {
      assert.expect(29);
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
      assert.dom(PAGE.secretTab('Secret')).hasText('Secret');
      assert.dom(PAGE.secretTab('Secret')).hasClass('active');
      assert.dom(PAGE.secretTab('Metadata')).hasText('Metadata');
      assert.dom(PAGE.secretTab('Metadata')).doesNotHaveClass('active');
      // TODO: hide tab
      // assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version history tab not shown');
      // TODO: hide dropdown
      // assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown hidden');
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

      // data-reader can't navigate to older versions, but they can go to page directly
      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`);
      // TODO: hide version dropdown
      // assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown does not exist');
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
      // TODO: version alert should exist
      // assert.dom(FORM.versionAlert).exists('Shows version alert');
      await click(FORM.cancelBtn);

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
    });
    test('breadcrumbs & page titles are correct', async function (assert) {
      assert.expect(30);
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

      await click(PAGE.detail.createNewVersion);
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'app/nested/secret', 'edit']);
      assert.dom(PAGE.title).hasText('Create New Version', 'title correct on create new version');

      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, 'app', 'nested', 'secret', 'metadata']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title correct on metadata');
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
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState', async function (assert) {
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
      assert.dom(PAGE.secretTab('list')).hasText('Secrets');
      assert.dom(PAGE.secretTab('list')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar renders');
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');
      // Page content correct
      assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
      assert.dom(PAGE.emptyStateActions).hasText('Create secret');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // Click empty state CTA
      await click(`${PAGE.emptyStateActions} a`);
      // TODO: initialKey should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      // TODO: pageFilter should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      // TODO: initialKey should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      // TODO: pageFilter should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret', async function (assert) {
      assert.expect(26);
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

      // TODO: overview card with pageFilter
      // assert
      //   .dom(PAGE.list.overviewInput)
      //   .hasValue('app/', 'overview card is pre-filled with directory param');
      // await typeIn(PAGE.list.overviewInput, 'nested/secret');
      await typeIn(PAGE.list.overviewInput, 'app/nested/secret'); // this is a workaround
      await click(PAGE.list.overviewButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details?version=1`
      );
      assertCorrectBreadcrumbs(assert, ['secret', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assert.dom(PAGE.toolbar).exists('toolbar renders');
      assert.dom(PAGE.toolbarAction).exists({ count: 2 }, 'correct number of toolbar actions render');

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
    test('versioned secret nav, tabs, breadcrumbs', async function (assert) {
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
      assert.dom(PAGE.secretTab('Secret')).hasText('Secret');
      assert.dom(PAGE.secretTab('Secret')).hasClass('active');
      assert.dom(PAGE.secretTab('Metadata')).hasText('Metadata');
      assert.dom(PAGE.secretTab('Metadata')).doesNotHaveClass('active');
      // TODO: version dropdown hidden
      // assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown hidden');
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

      // TODO: version dropdown should be hidden
      // assert.dom(PAGE.detail.versionDropdown).doesNotExist('version dropdown hidden');
      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`);
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
      // TODO: should show
      // assert.dom(FORM.versionAlert).exists('Shows version alert');
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
      assert.dom(PAGE.metadata.editBtn).doesNotExist('edit button hidden');
    });
    test('breadcrumbs & page titles are correct', async function (assert) {
      assert.expect(26);
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
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState', async function (assert) {
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
      assert.dom(PAGE.secretTab('list')).hasText('Secrets');
      assert.dom(PAGE.secretTab('list')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar renders');
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');
      // Page content correct
      assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
      assert.dom(PAGE.emptyStateActions).hasText('Create secret');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // Click empty state CTA
      await click(`${PAGE.emptyStateActions} a`);
      // TODO: initialKey should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      // TODO: pageFilter should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      // TODO: initialKey should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      // TODO: pageFilter should not show on query params if empty
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret', async function (assert) {
      assert.expect(36);
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
      assert.dom(PAGE.toolbar).exists('toolbar renders');
      // TODO: verify create new shouldn't show
      assert.dom(PAGE.toolbarAction).exists({ count: 1 }, 'correct number of toolbar actions render');
      // TODO: add version to dropdown when no data
      // assert.dom(PAGE.detail.versionDropdown).hasText('Version 1');

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
    test('versioned secret nav, tabs, breadcrumbs', async function (assert) {
      assert.expect(34);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));
      // TODO: url should have query param
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details`,
        'Url includes version query param'
      );
      assert.dom(PAGE.title).hasText(secretPath, 'Goes to secret detail view');
      assert.dom(PAGE.secretTab('Secret')).hasText('Secret');
      assert.dom(PAGE.secretTab('Secret')).hasClass('active');
      assert.dom(PAGE.secretTab('Metadata')).hasText('Metadata');
      assert.dom(PAGE.secretTab('Metadata')).doesNotHaveClass('active');
      assert.dom(PAGE.secretTab('Version History')).hasText('Version History');
      assert.dom(PAGE.secretTab('Version History')).doesNotHaveClass('active');
      assert
        .dom(PAGE.detail.versionDropdown)
        .hasText('Version current', 'Version dropdown shows current version');
      assert.dom(PAGE.detail.createNewVersion).doesNotExist('Create new version button not shown');
      // TODO: should the created metadata show?
      assert.dom(PAGE.detail.versionCreated).doesNotExist('Version created text not shown');
      assert.dom(PAGE.infoRowValue('foo')).doesNotExist('does not render current data');
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'Shows empty state on secret detail');

      await click(PAGE.detail.versionDropdown);
      await click(`${PAGE.detail.version(1)} a`);
      // TODO: version param missing
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`,
        'Goes to detail view for version 1'
      );
      // TODO: version number missing
      // assert.dom(PAGE.detail.versionDropdown).hasText('Version 1', 'Version dropdown shows selected version');
      // TODO: versionCreated missing
      // assert.dom(PAGE.detail.versionCreated).containsText('Version 1 created');
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
    test('breadcrumbs & page titles are correct', async function (assert) {
      assert.expect(33);
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
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState', async function (assert) {
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
      assert.dom(PAGE.secretTab('list')).hasText('Secrets');
      assert.dom(PAGE.secretTab('list')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar renders');
      assert.dom(PAGE.list.filter).doesNotExist('List filter input is not rendered');
      // Page content correct
      assert.dom(PAGE.list.overviewCard).exists('Overview card renders');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      // TODO: qp should not be present if empty
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=`,
        `goes to /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      // TODO: qp should not be present if empty
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list?pageFilter=`,
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret', async function (assert) {
      assert.expect(19);
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
      assert.dom(PAGE.toolbar).exists('toolbar renders');
      assert.dom(PAGE.toolbarAction).exists({ count: 1 }, 'correct number of toolbar actions render');

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
    test('versioned secret nav, tabs, breadcrumbs', async function (assert) {
      assert.expect(26);
      const backend = this.backend;
      await navToBackend(backend);

      await typeIn(PAGE.list.overviewInput, secretPath);
      await click(PAGE.list.overviewButton);
      // TODO: url should include version param
      // assert.strictEqual(
      //   currentURL(),
      //   `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=3`,
      //   'Url includes version query param'
      // );
      assert.dom(PAGE.title).hasText(secretPath, 'Goes to secret detail view');
      assert.dom(PAGE.secretTab('Secret')).hasText('Secret');
      assert.dom(PAGE.secretTab('Secret')).hasClass('active');
      assert.dom(PAGE.secretTab('Metadata')).hasText('Metadata');
      assert.dom(PAGE.secretTab('Metadata')).doesNotHaveClass('active');
      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab does not render');
      assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown does not render');
      assert.dom(PAGE.detail.createNewVersion).hasText('Create new version', 'Create version button shows');
      assert.dom(PAGE.detail.versionCreated).doesNotExist('Version created info is not rendered');
      assert.dom(PAGE.infoRowValue('foo')).doesNotExist('current data not rendered');
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'empty state shows');

      await click(PAGE.detail.createNewVersion);
      // TODO: url should include version param
      // assert.strictEqual(
      //   currentURL(),
      //   `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details/edit?version=3`,
      //   'Url includes version query param'
      // );
      assert.dom(FORM.versionAlert).doesNotExist('Does not show version alert for current version');
      // TODO: show this warning
      // assert
      //   .dom(FORM.noReadAlert)
      //   .hasText(
      //     'You do not have read permissions. If a secret exists at this path creating a new secret will overwrite it.',
      //     'Shows warning about no read permissions'
      //   );
      // TODO: input should be disabled
      // assert.dom(FORM.inputByAttr('path')).isDisabled();

      await click(FORM.cancelBtn);
      // TODO: version qp should exist
      // assert.strictEqual(
      //   currentURL(),
      //   `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=3`,
      //   'Goes back to detail view'
      // );

      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=1`);
      assert.dom(PAGE.detail.versionDropdown).doesNotExist('Version dropdown does not exist');
      assert.dom(PAGE.detail.versionCreated).doesNotExist('version created data not rendered');
      assert.dom(PAGE.infoRowValue('key-1')).doesNotExist('does not render previous data');

      await click(PAGE.detail.createNewVersion);
      // TODO: qp should exist
      // assert.strictEqual(
      //   currentURL(),
      //   `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details/edit?version=1`,
      //   'Url includes version query param'
      // );
      // TODO: path should be disabled
      // assert.dom(FORM.inputByAttr('path')).isDisabled();
      assert.dom(FORM.keyInput()).hasValue('', 'form does not pre-populate');
      assert.dom(FORM.maskedValueInput()).hasValue('', 'form does not pre-populate');
      // TODO: should show
      // assert.dom(FORM.noReadAlert).exists('Shows no read alert');
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
    test('breadcrumbs & page titles are correct', async function (assert) {
      assert.expect(32);
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

      await click(PAGE.secretTab('Version History'));
      assertCorrectBreadcrumbs(assert, ['secrets', backend, secretPath, 'version history']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for version history');
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
    });
    const storageKey = (accessor, path) => {
      return `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
    };
    test('can access nested secret', async function (assert) {
      assert.expect(38);
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
      const accessor = controlGroupComponent.accessor;
      const controlGroupToken = controlGroupComponent.token;

      await authPage.loginUsername('authorizer', 'password');
      await visit(`/vault/access/control-groups/${accessor}`);
      await controlGroupComponent.authorize();

      await authPage.login(this.userToken);
      localStorage.setItem(
        storageKey(accessor, `${backend}/data/app/nested/secret`),
        JSON.stringify({
          accessor,
          token: controlGroupToken,
          creation_path: `${backend}/data/app/nested/secret`,
          uiParams: {
            url: `/vault/secrets/${backend}/kv/app%2Fnested%2F/directory`,
          },
        })
      );
      await visit(`/vault/access/control-groups/${accessor}`);
      await click(`[data-test-navigate-button]`);
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
      assert.dom(PAGE.toolbar).exists('toolbar renders');
      assert.dom(PAGE.toolbarAction).exists({ count: 2 }, 'correct number of toolbar actions render');

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
  });
});
