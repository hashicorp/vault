/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import {
  click,
  currentRouteName,
  currentURL,
  find,
  findAll,
  fillIn,
  typeIn,
  visit,
  waitUntil,
} from '@ember/test-helpers';
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
import { personas } from 'vault/tests/helpers/kv/policy-generator';
import {
  addSecretMetadataCmd,
  clearRecords,
  writeSecret,
  writeVersionedSecret,
} from 'vault/tests/helpers/kv/kv-run-commands';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupControlGroup, grantAccess } from 'vault/tests/helpers/control-groups';

const secretPath = `my-#:$=?-secret`;
// This doesn't encode in a normal way, so hardcoding it here until we sort that out
const secretPathUrlEncoded = `my-%23:$=%3F-secret`;
// these are rendered individually by each page component, assigning a const here for consistency
const ALL_TABS = ['Overview', 'Secret', 'Metadata', 'Paths', 'Version History'];
const navToBackend = async (backend) => {
  await visit(`/vault/secrets`);
  return click(PAGE.backends.link(backend));
};
const assertCorrectBreadcrumbs = (assert, expected) => {
  assert.dom(PAGE.breadcrumbs).hasText(expected.join(' '));
  const breadcrumbs = findAll(PAGE.breadcrumb);
  expected.forEach((text, idx) => {
    assert.dom(breadcrumbs[idx]).hasText(text, `position ${idx} breadcrumb includes text ${text}`);
  });
};
const assertDetailTabs = (assert, current, hidden = []) => {
  ALL_TABS.forEach((tab) => {
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
// patchLatest is only available for enterprise so it's not included here
const DETAIL_TOOLBARS = ['delete', 'destroy', 'copy', 'versionDropdown', 'createNewVersion'];
const assertDetailsToolbar = (assert, expected = DETAIL_TOOLBARS) => {
  assert
    .dom(PAGE.toolbarAction)
    .exists({ count: expected.length }, 'correct number of toolbar actions render');
  expected.forEach((toolbar) => {
    assert.dom(PAGE.detail[toolbar]).exists(`${toolbar} action exists`);
  });
  const unexpected = DETAIL_TOOLBARS.filter((t) => !expected.includes(t));
  unexpected.forEach((toolbar) => {
    assert.dom(PAGE.detail[toolbar]).doesNotExist(`${toolbar} action doesNotExist`);
  });
};

const patchRedirectTest = (test, testCase) => {
  // only run this test on enterprise so we are testing permissions specifically and not enterprise vs CE (which also redirects)
  test(`enterprise: patch route redirects for users without permissions (${testCase})`, async function (assert) {
    await visit(`/vault/secrets/${this.backend}/kv/app%2Fnested%2Fsecret/patch`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.backend}/kv/app%2Fnested%2Fsecret`,
      'redirects to index'
    );
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.kv.secret.index');
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
    this.version = this.owner.lookup('service:version');
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

  test('KVv2 handles secret with % and space in path correctly', async function (assert) {
    // To check this bug no longer happens: https://github.com/hashicorp/vault/issues/11616
    await navToBackend(this.backend);
    await click(PAGE.list.createSecret);
    const pathWithSpace = 'per%centfu ll';
    await typeIn(GENERAL.inputByAttr('path'), pathWithSpace);
    await fillIn(FORM.keyInput(), 'someKey');
    await fillIn(FORM.maskedValueInput(), 'someValue');
    await click(FORM.saveBtn);
    assert.dom(PAGE.title).hasText(pathWithSpace, 'title is full path without any encoding/decoding.');
    assert
      .dom(PAGE.breadcrumbAtIdx(1))
      .hasText(this.backend, 'breadcrumb before secret path is backend path');
    assert
      .dom(PAGE.breadcrumbCurrentAtIdx(2))
      .hasText('per%centfu ll', 'the current breadcrumb is value of the secret path');

    await click(PAGE.breadcrumbAtIdx(1));
    assert
      .dom(`${PAGE.list.item(pathWithSpace)} [data-test-path]`)
      .hasText(pathWithSpace, 'the list item is shown correctly');

    await typeIn(PAGE.list.filter, 'per%');
    await click('[data-test-kv-list-filter-submit]');
    assert
      .dom(`${PAGE.list.item(pathWithSpace)} [data-test-path]`)
      .hasText(pathWithSpace, 'the list item is shown correctly after filtering');

    await click(PAGE.list.item(pathWithSpace));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.backend}/kv/${encodeURIComponent(pathWithSpace)}`,
      'Path is encoded in the URL'
    );
  });

  test('KVv2 handles nested secret with % and space in path correctly', async function (assert) {
    await navToBackend(this.backend);
    await click(PAGE.list.createSecret);
    const nestedPathWithSpace = 'per%/centfu ll';
    await typeIn(GENERAL.inputByAttr('path'), nestedPathWithSpace);
    await fillIn(FORM.keyInput(), 'someKey');
    await fillIn(FORM.maskedValueInput(), 'someValue');
    await click(FORM.saveBtn);
    assert
      .dom(PAGE.title)
      .hasText(
        nestedPathWithSpace,
        'title is of the full nested path (directory included) without any encoding/decoding.'
      );
    assert.dom(PAGE.breadcrumbAtIdx(2)).hasText('per%');
    assert
      .dom(PAGE.breadcrumbCurrentAtIdx(3))
      .hasText('centfu ll', 'the current breadcrumb is value centfu ll');

    await click(PAGE.breadcrumbAtIdx(1));
    assert
      .dom(`${PAGE.list.item('per%/')} [data-test-path]`)
      .hasText('per%/', 'the directory item is shown correctly');

    await typeIn(PAGE.list.filter, 'per%/');
    await click('[data-test-kv-list-filter-submit]');
    assert
      .dom(`${PAGE.list.item('centfu ll')} [data-test-path]`)
      .hasText('centfu ll', 'the list item is shown correctly after filtering');

    await click(PAGE.list.item('centfu ll'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.backend}/kv/${encodeURIComponent(nestedPathWithSpace)}`,
      'Path is encoded in the URL'
    );
  });

  test('KVv2 handles nested secret with a percent-encoded data octet in path correctly', async function (assert) {
    // To check this bug no longer happens: https://github.com/hashicorp/vault/issues/25905
    await navToBackend(this.backend);
    await click(PAGE.list.createSecret);
    const pathDataOctet = 'hello/foo%2fbar/world';
    await typeIn(GENERAL.inputByAttr('path'), pathDataOctet);
    await fillIn(FORM.keyInput(), 'someKey');
    await fillIn(FORM.maskedValueInput(), 'someValue');
    await click(FORM.saveBtn);
    assert
      .dom(PAGE.title)
      .hasText(
        pathDataOctet,
        'title is of the full nested path (directory included) without any encoding/decoding.'
      );
    assert
      .dom(PAGE.breadcrumbAtIdx(2))
      .hasText('hello', 'hello is the first directory and shows up as a separate breadcrumb');
    assert
      .dom(PAGE.breadcrumbAtIdx(3))
      .hasText('foo%2fbar', 'foo%2fbar is the second directory and shows up as a separate breadcrumb');
    assert.dom(PAGE.breadcrumbCurrentAtIdx(4)).hasText('world', 'the current breadcrumb is value world');

    await click(PAGE.breadcrumbAtIdx(2));
    assert
      .dom(`${PAGE.list.item('foo%2fbar/')} [data-test-path]`)
      .hasText('foo%2fbar/', 'the directory item is shown correctly');

    await click(PAGE.list.item('foo%2fbar/'));
    await click(PAGE.list.item('world'));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.backend}/kv/${encodeURIComponent(pathDataOctet)}`,
      'Path is encoded in the URL'
    );
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
      assert.expect(23);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // CONFIGURATION TAB
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'Configuration']);
      assert.dom(PAGE.secretTab('Configuration')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Secrets')).hasText('Secrets');
      assert.dom(PAGE.secretTab('Secrets')).doesNotHaveClass('active');
      // SECRETS TAB
      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      assert.dom(PAGE.secretTab('Secrets')).hasText('Secrets');
      assert.dom(PAGE.secretTab('Secrets')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbar).exists({ count: 1 }, 'toolbar renders');
      assert.dom(PAGE.list.filter).doesNotExist('List filter does not show because no secrets exists.');
      // Page content correct
      assert.dom(PAGE.emptyStateTitle).hasText('No secrets yet');
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret (a)', async function (assert) {
      // enterprise has "Patch latest version" in the toolbar which adds an assertion
      const count = this.version.isEnterprise ? 47 : 46;
      assert.expect(count);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Navigate through list items
      await click(PAGE.list.item('app/'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/app/`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('nested/')).exists('Shows nested secret');

      await click(PAGE.list.item('nested/'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/app/nested/`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/nested/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('secret')).exists('Shows deeply nested secret');

      await click(PAGE.list.item('secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');

      await click(PAGE.secretTab('Secret'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret']);
      const expectedToolbar = this.version.isEnterprise
        ? [...DETAIL_TOOLBARS, 'patchLatest']
        : DETAIL_TOOLBARS;
      assertDetailsToolbar(assert, expectedToolbar);

      await click(PAGE.breadcrumbAtIdx(3));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/nested/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(2));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(1));
      assert.true(currentURL().startsWith(`/vault/secrets/${backend}/kv/list`), 'links back to list root');
    });
    test('it redirects from LIST, SHOW and EDIT views using old non-engine url to ember engine url (a)', async function (assert) {
      assert.expect(4);
      const backend = this.backend;
      // create with initialKey
      await visit(`/vault/secrets/${backend}/create/test`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=test`,
        `navigated to ${currentURL()}`
      );
      // Reported bug, backported fix https://github.com/hashicorp/vault/pull/24281
      // list for directory
      await visit(`/vault/secrets/${backend}/list/app/`);
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list/app/`, `navigates to list`);
      // show for secret
      await visit(`/vault/secrets/${backend}/show/app/nested/secret`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret`,
        `navigates to overview`
      );
      // edit for secret
      await visit(`/vault/secrets/${backend}/edit/app/nested/secret`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details/edit?version=1`,
        `navigates to edit`
      );
    });
    test('versioned secret nav, tabs (a)', async function (assert) {
      assert.expect(27);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`,
        'navigates to overview'
      );
      await click(PAGE.secretTab('Secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details?version=3`,
        'Url includes version query param'
      );
      assert.dom(PAGE.title).hasText(secretPath, 'title is correct on detail view');
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
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`,
        'Goes back to overview'
      );
      await click(PAGE.secretTab('Secret'));
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
      await click(FORM.cancelBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/metadata`,
        `cancel btn goes back to metadata page`
      );
    });
    test('breadcrumbs, tabs & page titles are correct (a)', async function (assert) {
      assert.expect(123);
      // only need to assert hrefs one test, no need for this function to be global
      const assertTabHrefs = (assert, page) => {
        ALL_TABS.forEach((tab) => {
          const baseUrl = `/ui/vault/secrets/${backend}/kv`;
          const hrefs = {
            Overview: `${baseUrl}/${secretPathUrlEncoded}`,
            Secret:
              page === 'Secret'
                ? `${baseUrl}/${secretPathUrlEncoded}/details?version=3`
                : `${baseUrl}/${secretPathUrlEncoded}/details`,
            Metadata: `${baseUrl}/${secretPathUrlEncoded}/metadata`,
            Paths: `${baseUrl}/${secretPathUrlEncoded}/paths`,
            'Version History': `${baseUrl}/${secretPathUrlEncoded}/metadata/versions`,
          };
          assert
            .dom(PAGE.secretTab(tab))
            .hasAttribute('href', hrefs[tab], `${tab} tab for page: ${page} has expected href`);
        });
      };
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));

      // PAGE COMPONENTS RENDER THEIR OWN TABS, ASSERT EACH HREF ON EACH PAGE
      // overview tab
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.index',
        'navs to overview'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath]);
      assertDetailTabs(assert, 'Overview');
      assertTabHrefs(assert, 'Overview');
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret overview');

      // secret tab
      await click(PAGE.secretTab('Secret'));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.details.index',
        'navs to details'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath]);
      assertDetailTabs(assert, 'Secret');
      assertTabHrefs(assert, 'Secret');
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      await click(PAGE.detail.createNewVersion);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.details.edit',
        'navs to create'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Edit']);
      assert.dom(PAGE.title).hasText('Create New Version', 'correct page title for secret edit');

      // metadata tab
      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Metadata'));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.metadata.index',
        'navs to metadata'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
      assertDetailTabs(assert, 'Metadata');
      assertTabHrefs(assert, 'Metadata');
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      await click(PAGE.metadata.editBtn);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.metadata.edit',
        'navs to metadata.edit'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata', 'Edit']);
      assert.dom(PAGE.title).hasText('Edit Secret Metadata', 'correct page title for metadata edit');

      // paths tab
      await click(PAGE.breadcrumbAtIdx(3));
      await click(PAGE.secretTab('Paths'));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.paths',
        'navs to paths'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Paths']);
      assertDetailTabs(assert, 'Paths');
      assertTabHrefs(assert, 'Paths');
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      // version history tab
      await click(PAGE.secretTab('Version History'));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.metadata.versions',
        'navs to version history'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Version History']);
      assertDetailTabs(assert, 'Version History');
      assertTabHrefs(assert, 'Version History');
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for version history');
    });
    // only run this test on enterprise so we are testing permissions specifically and not enterprise vs CE (which also redirects)
    test('enterprise: patch route does not redirect for users with permissions (a)', async function (assert) {
      await visit(`/vault/secrets/${this.backend}/kv/app%2Fnested%2Fsecret/patch`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/kv/app%2Fnested%2Fsecret/patch`,
        'redirects to index'
      );
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.kv.secret.patch');
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          `data-reader-${this.backend}`,
          personas.dataReader(this.backend) + personas.dataReader(this.emptyBackend)
        ),
        createTokenCmd(`data-reader-${this.backend}`),
      ]);
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState (dr)', async function (assert) {
      assert.expect(16);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // Breadcrumbs correct
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
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

      await typeIn(PAGE.list.overviewInput, 'directory/');
      await click(PAGE.list.overviewButton);
      assert
        .dom('[data-test-inline-error-message]')
        .hasText('You do not have the required permissions or the directory does not exist.');

      // click toolbar CTA
      await visit(`/vault/secrets/${backend}/kv/list`);
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
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert
        .dom(PAGE.list.filter)
        .doesNotExist('List filter input does not render because no list capabilities');

      await typeIn(PAGE.list.overviewInput, 'app/nested/secret');
      await click(PAGE.list.overviewButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret`,
        `navigated to secret overview ${currentURL()}`
      );
      await click(PAGE.secretTab('Secret'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['copy']);

      await click(PAGE.breadcrumbAtIdx(3));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/nested/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(2));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(1));
      assert.true(currentURL().startsWith(`/vault/secrets/${backend}/kv/list`), 'links back to list root');
    });
    test('versioned secret nav, tabs, breadcrumbs (dr)', async function (assert) {
      assert.expect(32);
      const backend = this.backend;
      await navToBackend(backend);

      // Navigate to secret
      await typeIn(PAGE.list.overviewInput, secretPath);
      await click(PAGE.list.overviewButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`,
        'navigates to secret overview'
      );
      await click(PAGE.secretTab('Secret'));
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
      assert.dom(PAGE.title).hasText(secretPath);
      assert.dom(PAGE.toolbarAction).doesNotExist('no toolbar actions available on metadata');
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('Request custom metadata?');
      await click(PAGE.metadata.requestData);
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'Configuration']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'title correct on config page');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'title correct on secrets list');

      await typeIn(PAGE.list.overviewInput, 'app/nested/secret');
      await click(PAGE.list.overviewButton);
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title correct on secret detail');

      assert.dom(PAGE.detail.createNewVersion).doesNotExist('cannot create new version');

      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret', 'Metadata']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title correct on metadata');

      assert.dom(PAGE.metadata.editBtn).doesNotExist('cannot edit metadata');

      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret', 'Paths']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'correct page title for paths');

      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab not shown');
    });
    patchRedirectTest(test, 'dr');
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          `data-reader-list-${this.backend}`,
          personas.dataListReader(this.backend) + personas.dataListReader(this.emptyBackend)
        ),
        createTokenCmd(`data-reader-list-${this.backend}`),
      ]);

      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState (dlr)', async function (assert) {
      assert.expect(15);
      const backend = this.emptyBackend;
      await navToBackend(backend);

      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // Breadcrumbs correct
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
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
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret (dlr)', async function (assert) {
      assert.expect(32);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Navigate through list items
      await click(PAGE.list.item('app/'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/app/`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      assert.dom(PAGE.list.filter).doesNotExist('List filter hidden since no nested list access');

      assert
        .dom(PAGE.list.overviewInput)
        .hasValue('app/', 'overview card is pre-filled with directory param');
      await typeIn(PAGE.list.overviewInput, 'nested/secret');
      await click(PAGE.list.overviewButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret`,
        `navigated to overview`
      );
      await click(PAGE.secretTab('Secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details?version=1`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['delete', 'copy']);

      await click(PAGE.breadcrumbAtIdx(3));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/nested/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(2));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(1));
      assert.true(currentURL().startsWith(`/vault/secrets/${backend}/kv/list`), 'links back to list root');
    });
    test('versioned secret nav, tabs, breadcrumbs (dlr)', async function (assert) {
      assert.expect(32);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`,
        'navigates to overview'
      );
      await click(PAGE.secretTab('Secret'));
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
      assert.dom(PAGE.title).hasText(secretPath);
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('Request custom metadata?');
      await click(PAGE.metadata.requestData);
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'Configuration']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'correct page title for secret list');

      await click(PAGE.list.item(secretPath));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      assert.dom(PAGE.detail.createNewVersion).doesNotExist('cannot create new version');

      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      assert.dom(PAGE.metadata.editBtn).doesNotExist('cannot edit metadata');

      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab not shown');
    });
    patchRedirectTest(test, 'dlr');
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          `metadata-maintainer-${this.backend}`,
          personas.metadataMaintainer(this.backend) + personas.metadataMaintainer(this.emptyBackend)
        ),
        createTokenCmd(`metadata-maintainer-${this.backend}`),
      ]);
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('empty backend - breadcrumbs, title, tabs, emptyState (mm)', async function (assert) {
      assert.expect(15);
      const backend = this.emptyBackend;
      await navToBackend(backend);
      // URL correct
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');
      // Breadcrumbs correct
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
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
      assert.dom(PAGE.list.createSecret).hasText('Create secret');

      // click toolbar CTA
      await click(PAGE.list.createSecret);
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/create`),
        `url includes /vault/secrets/${backend}/kv/create`
      );

      // Click cancel btn
      await click(FORM.cancelBtn);
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list`),
        `url includes /vault/secrets/${backend}/kv/list`
      );
    });
    test('can access nested secret (mm)', async function (assert) {
      assert.expect(42);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Navigate through list items
      await click(PAGE.list.item('app/'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/app/`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('nested/')).exists('Shows nested secret');

      await click(PAGE.list.item('nested/'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/app/nested/`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/nested/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('secret')).exists('Shows deeply nested secret');

      await click(PAGE.list.item('secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret`,
        `goes to overview`
      );
      await click(PAGE.secretTab('Secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details`,
        `Goes to URL without version`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['delete', 'destroy', 'versionDropdown']);
      assert.dom(PAGE.detail.versionDropdown).hasText('Version 1', 'Shows version timestamp');

      await click(PAGE.breadcrumbAtIdx(3));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/nested/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(2));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(1));
      assert.true(currentURL().startsWith(`/vault/secrets/${backend}/kv/list`), 'links back to list root');
    });
    test('versioned secret nav, tabs, breadcrumbs (mm)', async function (assert) {
      assert.expect(40);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.list.item(secretPath));

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`,
        'navs to overview'
      );
      await click(PAGE.secretTab('Secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details`,
        'Url does not include version query param'
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata', 'Edit']);
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'Configuration']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'correct page title for secret list');

      await click(PAGE.list.item(secretPath));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      await click(PAGE.metadata.editBtn);
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata', 'Edit']);
      assert.dom(PAGE.title).hasText('Edit Secret Metadata', 'correct page title for metadata edit');

      await click(PAGE.breadcrumbAtIdx(3));
      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      await click(PAGE.secretTab('Version History'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Version History']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for version history');
    });
    patchRedirectTest(test, 'mm');
  });

  module('secret-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          `secret-creator-${this.backend}`,
          personas.secretCreator(this.backend) + personas.secretCreator(this.emptyBackend)
        ),
        createTokenCmd(`secret-creator-${this.backend}`),
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
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
      assert.expect(24);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.list.filter).doesNotExist('List filter input is not rendered');

      // Navigate to secret
      await typeIn(PAGE.list.overviewInput, 'app/nested/secret');
      await click(PAGE.list.overviewButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret`,
        'goes to overview'
      );
      await click(PAGE.secretTab('Secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details`,
        'goes to secret detail page'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['createNewVersion']);

      await click(PAGE.breadcrumbAtIdx(3));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/nested/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(2));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(1));
      assert.true(currentURL().startsWith(`/vault/secrets/${backend}/kv/list`), 'links back to list root');
    });
    test('versioned secret nav, tabs, breadcrumbs (sc)', async function (assert) {
      assert.expect(39);
      const backend = this.backend;
      await navToBackend(backend);

      await typeIn(PAGE.list.overviewInput, secretPath);
      await click(PAGE.list.overviewButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`,
        'Goes to overview'
      );
      await click(PAGE.secretTab('Secret'));
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
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`,
        'Goes back to overview'
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
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
      assert.dom(PAGE.title).hasText(secretPath);
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('You do not have access to read custom metadata', 'shows correct empty state');
      assert.dom(PAGE.metadata.editBtn).doesNotExist('edit metadata button does not render');
    });
    test('breadcrumbs & page titles are correct (sc)', async function (assert) {
      assert.expect(39);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'Configuration']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'correct page title for secret list');

      await typeIn(PAGE.list.overviewInput, secretPath);
      await click(PAGE.list.overviewButton);
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      await click(PAGE.secretTab('Secret'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret detail');

      await click(PAGE.detail.createNewVersion);
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Edit']);
      assert.dom(PAGE.title).hasText('Create New Version', 'correct page title for secret edit');

      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');

      assert.dom(PAGE.metadata.editBtn).doesNotExist('cannot edit metadata');

      await click(PAGE.breadcrumbAtIdx(2));
      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab not shown');
    });
    patchRedirectTest(test, 'sc');
  });

  module('enterprise controlled access persona', function (hooks) {
    hooks.beforeEach(async function () {
      // Set up control group scenario
      const userPolicy = `
path "${this.backend}/data/*" {
  capabilities = ["create", "read", "update", "delete", "list", "patch"]
  control_group = {
    max_ttl = "24h"
    factor "ops_manager" {
      controlled_capabilities = ["read", "patch"]
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

path "${this.backend}/subkeys/*" {
  capabilities = ["read"]
}
`;
      const { userToken } = await setupControlGroup({ userPolicy, backend: this.backend });
      this.userToken = userToken;
      await authPage.login(userToken);
      clearRecords(this.store);
      return;
    });
    test('can access nested secret (cg)', async function (assert) {
      assert.expect(44);
      const backend = this.backend;
      await navToBackend(backend);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'title text correct');
      assert.dom(PAGE.emptyStateTitle).doesNotExist('No empty state');
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.list.filter).hasNoValue('List filter input is empty');

      // Navigate through list items
      await click(PAGE.list.item('app/'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/app/`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('nested/')).exists('Shows nested secret');

      await click(PAGE.list.item('nested/'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/app/nested/`,
        `navigated to ${currentURL()}`
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      assert.dom(PAGE.list.filter).hasValue('app/nested/', 'List filter input is prefilled');
      assert.dom(PAGE.list.item('secret')).exists('Shows deeply nested secret');

      // For some reason when we click on the item in tests it throws a global control group error
      // But not when we visit the page directly
      await visit(`/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details`);
      assert.true(
        await waitUntil(() => currentRouteName() === 'vault.cluster.access.control-group-accessor'),
        'redirects to access control group route'
      );
      await grantAccess({
        apiPath: `${backend}/data/app/nested/secret`,
        originUrl: `/vault/secrets/${backend}/kv/list/app/nested/`,
        userToken: this.userToken,
        backend: this.backend,
      });
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/app/nested/`,
        'navigates to list url where secret is'
      );
      await click(PAGE.list.item('secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret`,
        'goes to overview'
      );

      await click(PAGE.secretTab('Secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Fnested%2Fsecret/details?version=1`,
        'goes to secret details'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'app', 'nested', 'secret']);
      assert.dom(PAGE.title).hasText('app/nested/secret', 'title is full secret path');
      assertDetailsToolbar(assert, ['delete', 'copy', 'createNewVersion', 'patchLatest']);

      await click(PAGE.breadcrumbAtIdx(3));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/nested/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(2));
      assert.true(
        currentURL().startsWith(`/vault/secrets/${backend}/kv/list/app/`),
        'links back to list directory'
      );

      await click(PAGE.breadcrumbAtIdx(1));
      assert.true(currentURL().startsWith(`/vault/secrets/${backend}/kv/list`), 'links back to list root');
    });
    test('breadcrumbs & page titles are correct (cg)', async function (assert) {
      assert.expect(42);
      const backend = this.backend;
      await navToBackend(backend);
      await click(PAGE.secretTab('Configuration'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, 'Configuration']);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'correct page title for configuration');

      await click(PAGE.secretTab('Secrets'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend]);
      assert.dom(PAGE.title).hasText(`${backend} version 2`, 'correct page title for secret list');

      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/details`);

      assert.true(
        await waitUntil(() => currentRouteName() === 'vault.cluster.access.control-group-accessor'),
        'redirects to access control group route'
      );
      await grantAccess({
        apiPath: `${backend}/data/${encodeURIComponent(secretPath)}`,
        originUrl: `/vault/secrets/${backend}/kv/list`,
        userToken: this.userToken,
        backend: this.backend,
      });

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list`,
        'navigates back to list url after authorized'
      );
      await click(PAGE.list.item(secretPath));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`,
        'Goes to overview'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret overview');

      await click(PAGE.secretTab('Metadata'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Metadata']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for metadata');
      assert.dom(PAGE.metadata.editBtn).doesNotExist('cannot edit metadata');

      await click(PAGE.secretTab('Paths'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Paths']);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for paths');

      assert.dom(PAGE.secretTab('Version History')).doesNotExist('Version History tab not shown');

      await click(PAGE.secretTab('Secret'));
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath]);
      assert.dom(PAGE.title).hasText(secretPath, 'correct page title for secret details');
      await click(PAGE.detail.createNewVersion);
      assertCorrectBreadcrumbs(assert, ['Secrets', backend, secretPath, 'Edit']);
      assert.dom(PAGE.title).hasText('Create New Version', 'correct page title for secret edit');
    });
    test('can request custom_metadata from data endpoint (cg)', async function (assert) {
      // custom metadata is empty
      assert.expect(3);
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`);
      await click(PAGE.secretTab('Metadata'));
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('Request custom metadata?');
      await click(PAGE.metadata.requestData);
      assert
        .dom(GENERAL.messageError)
        .hasTextContaining(
          `Control Group Error A Control Group was encountered at ${backend}/data/${secretPath}.`
        );
      const url = find('[data-test-control-error="href"]').innerText;
      await visit(url);
      await grantAccess({
        apiPath: `${backend}/data/${encodeURIComponent(secretPath)}`,
        originUrl: `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/metadata`,
        userToken: this.userToken,
        backend: this.backend,
      });
      await click(PAGE.metadata.requestData);
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('No custom metadata', 'empty state updates when access is granted');
    });
    test('can patch a secret (cg)', async function (assert) {
      assert.expect(3);
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`);
      await click(GENERAL.overviewCard.actionText('Patch secret'));
      await fillIn(FORM.keyInput('new'), 'newkey');
      await fillIn(FORM.valueInput('new'), 'newvalue');
      await click(FORM.saveBtn);
      assert
        .dom(GENERAL.messageError)
        .hasTextContaining(
          `Control Group Error A Control Group was encountered at ${backend}/data/${secretPath}.`
        );
      assert
        .dom(GENERAL.messageError)
        .hasTextContaining(
          'You can re-submit the form once access is granted. Ask your authorizer when to attempt saving again.'
        );
      const url = find('[data-test-control-error="href"]').innerText;
      await visit(url);
      await grantAccess({
        apiPath: `${backend}/data/${encodeURIComponent(secretPath)}`,
        originUrl: `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/patch`,
        userToken: this.userToken,
        backend: this.backend,
      });
      // we have to refill the data because granting access reloads the form
      // however in the real world it's likely access is authorized in a separate browser
      // once granted, the user can click "submit" the form will save successfully.
      await fillIn(FORM.keyInput('new'), 'newkey');
      await fillIn(FORM.valueInput('new'), 'newvalue');
      await click(FORM.saveBtn);
      assert.dom(GENERAL.overviewCard.container('Subkeys')).hasTextContaining('Keys foo newkey');
    });
    test('can read custom_metadata from data endpoint (cg)', async function (assert) {
      assert.expect(3);
      // login is root user and make custom metadata since console can't be used to pass an object
      await authPage.login();
      await visit(`/vault/secrets/${this.backend}/kv/${secretPathUrlEncoded}/metadata/edit`);
      await fillIn(FORM.keyInput(), 'special');
      await fillIn(FORM.valueInput(), 'secret');
      await click(FORM.saveBtn);
      await authPage.login(this.userToken);

      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/${secretPathUrlEncoded}`);

      await click(PAGE.secretTab('Metadata'));
      assert
        .dom(`${PAGE.metadata.customMetadataSection} ${PAGE.emptyStateTitle}`)
        .hasText('Request custom metadata?');
      await click(PAGE.metadata.requestData);
      assert
        .dom(GENERAL.messageError)
        .hasTextContaining(
          `Control Group Error A Control Group was encountered at ${backend}/data/${secretPath}.`
        );
      const url = find('[data-test-control-error="href"]').innerText;
      await visit(url);
      await grantAccess({
        apiPath: `${backend}/data/${encodeURIComponent(secretPath)}`,
        originUrl: `/vault/secrets/${backend}/kv/${secretPathUrlEncoded}/metadata`,
        userToken: this.userToken,
        backend: this.backend,
      });
      await click(PAGE.metadata.requestData);
      assert.dom(PAGE.infoRowValue('special')).hasText('secret', 'it renders custom metadata');
    });
  });

  // patch is technically enterprise only but stubbing the version so these tests can run on both CE and enterprise
  module('patch-persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd([
        createPolicyCmd(
          `secret-patcher-${this.backend}`,
          personas.secretPatcher(this.backend) + personas.secretPatcher(this.emptyBackend)
        ),
        createTokenCmd(`secret-patcher-${this.backend}`),
      ]);
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });

    test('it navigates to patch a secret from overview', async function (assert) {
      this.version.type = 'enterprise';
      await navToBackend(this.backend);
      await click(PAGE.list.item(secretPath));
      await click(GENERAL.overviewCard.actionText('Patch secret'));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.patch',
        'navs to patch'
      );
      assertCorrectBreadcrumbs(assert, ['Secrets', this.backend, secretPath, 'Patch']);
      assert.dom(PAGE.title).hasText('Patch Secret to New Version');
      await click(FORM.cancelBtn);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.index',
        'navs back to overview'
      );
    });

    test('overview subkeys card is hidden for community edition', async function (assert) {
      this.version.type = 'community';
      await navToBackend(this.backend);
      await click(PAGE.list.item(secretPath));
      assert.dom(GENERAL.overviewCard.container('Subkeys')).doesNotExist();
    });

    test('it does not redirect for ent', async function (assert) {
      this.version.type = 'enterprise';
      await visit(`/vault/secrets/${this.backend}/kv/app%2Fnested%2Fsecret/patch`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/kv/app%2Fnested%2Fsecret/patch`,
        'redirects to index'
      );
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.kv.secret.patch');
    });

    test('it redirects for community edition', async function (assert) {
      this.version.type = 'community';
      await visit(`/vault/secrets/${this.backend}/kv/app%2Fnested%2Fsecret/patch`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/kv/app%2Fnested%2Fsecret`,
        'redirects to index'
      );
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.kv.secret.index');
    });
  });
});
