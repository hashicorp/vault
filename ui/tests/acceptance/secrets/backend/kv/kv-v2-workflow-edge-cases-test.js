/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentURL, fillIn, findAll, setupOnerror, typeIn, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import {
  createPolicyCmd,
  deleteEngineCmd,
  mountEngineCmd,
  runCmd,
  createTokenCmd,
} from 'vault/tests/helpers/commands';
import {
  dataPolicy,
  deleteVersionsPolicy,
  destroyVersionsPolicy,
  metadataListPolicy,
  metadataPolicy,
} from 'vault/tests/helpers/kv/policy-generator';
import {
  KV_FORM,
  KV_WORKFLOW,
  clearRecords,
  writeSecret,
  writeVersionedSecret,
} from 'vault/tests/helpers/kv/kv-selectors';
import codemirror from 'vault/tests/helpers/codemirror';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { KV_METADATA_DETAILS } from 'vault/tests/helpers/components/kv/page/secret/metadata/details-selectors';
import { KV_SECRET } from 'vault/tests/helpers/components/kv/page/secret/details-selectors';
import { KV_LIST } from 'vault/tests/helpers/components/kv/page/list-selectors';

/**
 * This test set is for testing edge cases, such as specific bug fixes or reported user workflows
 */
module('Acceptance | kv-v2 workflow | edge cases', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const uid = uuidv4();
    this.backend = `kv-edge-${uid}`;
    this.rootSecret = 'root-directory';
    this.fullSecretPath = `${this.rootSecret}/nested/child-secret`;
    await authPage.login();
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await writeSecret(this.backend, this.fullSecretPath, 'foo', 'bar');
    await writeSecret(this.backend, 'edge/one', 'foo', 'bar');
    await writeSecret(this.backend, 'edge/two', 'foo', 'bar');
    return;
  });

  hooks.afterEach(async function () {
    await authPage.login();
    await runCmd(deleteEngineCmd(this.backend));
    return;
  });

  module('persona with read and list access on the secret level', function (hooks) {
    // see github issue for more details https://github.com/hashicorp/vault/issues/5362
    hooks.beforeEach(async function () {
      const secretPath = `${this.rootSecret}/*`; // user has LIST and READ access within this root secret directory
      const capabilities = ['list', 'read'];
      const backend = this.backend;
      const token = await runCmd([
        createPolicyCmd(
          `nested-secret-list-reader-${this.backend}`,
          metadataPolicy({ backend, secretPath, capabilities }) +
            dataPolicy({ backend, secretPath, capabilities })
        ),
        createTokenCmd(`nested-secret-list-reader-${this.backend}`),
      ]);
      await authPage.login(token);
    });

    test('it can navigate to secrets within a secret directory', async function (assert) {
      assert.expect(21);
      const backend = this.backend;
      const [root, subdirectory, secret] = this.fullSecretPath.split('/');

      await visit(`/vault/secrets/${backend}/kv/list`);
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');

      await typeIn(KV_LIST.overviewInput, `${root}/no-access/`);
      assert
        .dom(KV_LIST.overviewButton)
        .hasText('View list', 'shows list and not secret because search is a directory');
      await click(KV_LIST.overviewButton);
      assert.dom(GENERAL.emptyStateTitle).hasText(`There are no secrets matching "${root}/no-access/".`);

      await visit(`/vault/secrets/${backend}/kv/list`);
      await typeIn(KV_LIST.overviewInput, `${root}/`); // add slash because this is a directory
      await click(KV_LIST.overviewButton);

      // URL correct
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/${root}/`,
        'visits list-directory of root'
      );

      // Title correct
      assert.dom(GENERAL.title).hasText(`${backend} version 2`);
      // Tabs correct
      assert.dom(GENERAL.tab('Secrets')).hasText('Secrets');
      assert.dom(GENERAL.tab('Secrets')).hasClass('active');
      assert.dom(GENERAL.tab('Configuration')).hasText('Configuration');
      assert.dom(GENERAL.tab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(KV_WORKFLOW.toolbarAction).exists({ count: 1 }, 'toolbar only renders create secret action');
      assert.dom(KV_LIST.filter).hasValue(`${root}/`);
      // List content correct
      assert.dom(KV_LIST.item(`${subdirectory}/`)).exists('renders linked block for subdirectory');
      await click(KV_LIST.item(`${subdirectory}/`));
      assert.dom(KV_LIST.item(secret)).exists('renders linked block for child secret');
      await click(KV_LIST.item(secret));
      // Secret details visible
      assert.dom(GENERAL.title).hasText(this.fullSecretPath);
      assert.dom(GENERAL.tab('Secret')).hasText('Secret');
      assert.dom(GENERAL.tab('Secret')).hasClass('active');
      assert.dom(GENERAL.tab('Metadata')).hasText('Metadata');
      assert.dom(GENERAL.tab('Metadata')).doesNotHaveClass('active');
      assert.dom(GENERAL.tab('Version History')).hasText('Version History');
      assert.dom(GENERAL.tab('Version History')).doesNotHaveClass('active');
      assert.dom(KV_WORKFLOW.toolbarAction).exists({ count: 4 }, 'toolbar renders all actions');
    });

    test('it navigates back to engine index route via breadcrumbs from secret details', async function (assert) {
      assert.expect(6);
      const backend = this.backend;
      const [root, subdirectory, secret] = this.fullSecretPath.split('/');

      await visit(`vault/secrets/${backend}/kv/${encodeURIComponent(this.fullSecretPath)}/details?version=1`);
      // navigate back through crumbs
      let previousCrumb = findAll('[data-test-breadcrumbs] li').length - 2;
      await click(GENERAL.breadcrumbAtIdx(previousCrumb));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/${root}/${subdirectory}/`,
        'goes back to subdirectory list'
      );
      assert.dom(KV_LIST.filter).hasValue(`${root}/${subdirectory}/`);
      assert.dom(KV_LIST.item(secret)).exists('renders linked block for child secret');

      // back again
      previousCrumb = findAll('[data-test-breadcrumbs] li').length - 2;
      await click(GENERAL.breadcrumbAtIdx(previousCrumb));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/${root}/`,
        'goes back to root directory'
      );
      assert.dom(KV_LIST.item(`${subdirectory}/`)).exists('renders linked block for subdirectory');

      // and back to the engine list view
      previousCrumb = findAll('[data-test-breadcrumbs] li').length - 2;
      await click(GENERAL.breadcrumbAtIdx(previousCrumb));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list`,
        'navigates back to engine list from crumbs'
      );
    });

    test('it handles errors when attempting to view details of a secret that is a directory', async function (assert) {
      assert.expect(7);
      const backend = this.backend;
      const [root, subdirectory] = this.fullSecretPath.split('/');
      setupOnerror((error) => assert.strictEqual(error.httpStatus, 404), '404 error is thrown'); // catches error so qunit test doesn't fail

      await visit(`/vault/secrets/${backend}/kv/list`);
      await typeIn(KV_LIST.overviewInput, `${root}/${subdirectory}`); // intentionally leave out trailing slash
      await click(KV_LIST.overviewButton);
      assert.dom(KV_WORKFLOW.error.title).hasText('404 Not Found');
      assert
        .dom(KV_WORKFLOW.error.message)
        .hasText(`Sorry, we were unable to find any content at /v1/${backend}/data/${root}/${subdirectory}.`);

      assert.dom(GENERAL.breadcrumbAtIdx(0)).hasText('secrets');
      assert.dom(GENERAL.breadcrumbAtIdx(1)).hasText(backend);
      assert.dom(GENERAL.tab('Secrets')).doesNotHaveClass('is-active');
      assert.dom(GENERAL.tab('Configuration')).doesNotHaveClass('is-active');
    });
  });

  module('destruction without read', function (hooks) {
    hooks.beforeEach(async function () {
      const backend = this.backend;
      const testSecrets = [
        'data-delete-only',
        'delete-version-only',
        'destroy-version-only',
        'destroy-metadata-only',
      ];

      // user has different permissions for each secret path
      const token = await runCmd([
        createPolicyCmd(
          `destruction-no-read-${this.backend}`,
          dataPolicy({ backend, secretPath: 'data-delete-only', capabilities: ['delete'] }) +
            deleteVersionsPolicy({ backend, secretPath: 'delete-version-only' }) +
            destroyVersionsPolicy({ backend, secretPath: 'destroy-version-only' }) +
            metadataPolicy({ backend, secretPath: 'destroy-metadata-only', capabilities: ['delete'] }) +
            metadataListPolicy(backend)
        ),
        createTokenCmd(`destruction-no-read-${this.backend}`),
      ]);
      for (const secret of testSecrets) {
        await writeVersionedSecret(backend, secret, 'foo', 'bar', 2);
      }
      await authPage.login(token);
    });

    test('it renders the delete action and disables delete this version option', async function (assert) {
      assert.expect(4);
      const testSecret = 'data-delete-only';
      await visit(`/vault/secrets/${this.backend}/kv/${testSecret}/details`);

      assert.dom(KV_SECRET.delete).exists('renders delete button');
      await click(KV_SECRET.delete);
      assert
        .dom(KV_SECRET.deleteModal)
        .hasTextContaining('Delete this version This deletes a specific version of the secret');
      assert.dom(KV_SECRET.deleteOption).isDisabled('disables version specific option');
      assert.dom(KV_SECRET.deleteOptionLatest).isEnabled('enables version specific option');
    });

    test('it renders the delete action and disables delete latest version option', async function (assert) {
      assert.expect(4);
      const testSecret = 'delete-version-only';
      await visit(`/vault/secrets/${this.backend}/kv/${testSecret}/details`);

      assert.dom(KV_SECRET.delete).exists('renders delete button');
      await click(KV_SECRET.delete);
      assert
        .dom(KV_SECRET.deleteModal)
        .hasTextContaining('Delete this version This deletes a specific version of the secret');

      assert.dom(KV_SECRET.deleteOption).isEnabled('enables version specific option');
      assert.dom(KV_SECRET.deleteOptionLatest).isDisabled('disables version specific option');
    });

    test('it hides destroy option without version number', async function (assert) {
      assert.expect(1);
      const testSecret = 'destroy-version-only';
      await visit(`/vault/secrets/${this.backend}/kv/${testSecret}/details`);

      assert.dom(KV_SECRET.destroy).doesNotExist();
    });

    test('it renders the destroy metadata action and expected modal copy', async function (assert) {
      assert.expect(2);

      const testSecret = 'destroy-metadata-only';
      await visit(`/vault/secrets/${this.backend}/kv/${testSecret}/metadata`);
      assert.dom(KV_METADATA_DETAILS.deleteMetadata).exists('renders delete metadata button');
      await click(KV_METADATA_DETAILS.deleteMetadata);
      assert
        .dom(KV_SECRET.deleteModal)
        .hasText(
          'Delete metadata and secret data? This will permanently delete the metadata and versions of the secret. All version history will be removed. This cannot be undone. Confirm Cancel'
        );
    });
  });

  test('no ghost item after editing metadata', async function (assert) {
    await visit(`/vault/secrets/${this.backend}/kv/list/edge/`);
    assert.dom(KV_LIST.item()).exists({ count: 2 }, 'two secrets are listed');
    await click(KV_LIST.item('two'));
    await click(GENERAL.tab('Metadata'));
    await click(KV_METADATA_DETAILS.editBtn);
    await fillIn(KV_FORM.keyInput(), 'foo');
    await fillIn(KV_FORM.valueInput(), 'bar');
    await click(GENERAL.saveButton);
    await click(GENERAL.breadcrumbAtIdx(2));
    assert.dom(KV_LIST.item()).exists({ count: 2 }, 'two secrets are listed');
  });

  test('advanced secret values default to JSON display', async function (assert) {
    const obscuredData = `{
  "foo3": {
    "name": "********"
  }
}`;
    await visit(`/vault/secrets/${this.backend}/kv/create`);
    await fillIn(GENERAL.inputByAttr('path'), 'complex');

    await click(KV_FORM.toggleJson);
    assert.strictEqual(codemirror().getValue(), '{ "": "" }');
    codemirror().setValue('{ "foo3": { "name": "bar3" } }');
    await click(GENERAL.saveButton);

    // Details view
    assert.dom(KV_FORM.toggleJson).isNotDisabled();
    assert.dom(KV_FORM.toggleJson).isChecked();
    assert.strictEqual(
      codemirror().getValue(),
      obscuredData,
      'Value is obscured by default on details view when advanced'
    );
    await click('[data-test-toggle-input="revealValues"]');
    assert.false(codemirror().getValue().includes('*'), 'Value unobscured after toggle');

    // New version view
    await click(KV_SECRET.createNewVersion);
    assert.dom(KV_FORM.toggleJson).isNotDisabled();
    assert.dom(KV_FORM.toggleJson).isChecked();
    assert.false(codemirror().getValue().includes('*'), 'Values are not obscured on edit view');
  });

  test('viewing advanced secret data versions displays the correct version data', async function (assert) {
    assert.expect(2);
    const obscuredDataV1 = `{
  "foo1": {
    "name": "********"
  }
}`;
    const obscuredDataV2 = `{
  "foo2": {
    "name": "********"
  }
}`;

    await visit(`/vault/secrets/${this.backend}/kv/create`);
    await fillIn(GENERAL.inputByAttr('path'), 'complex_version_test');

    await click(KV_FORM.toggleJson);
    codemirror().setValue('{ "foo1": { "name": "bar1" } }');
    await click(GENERAL.saveButton);

    // Create another version
    await click(KV_SECRET.createNewVersion);
    codemirror().setValue('{ "foo2": { "name": "bar2" } }');
    await click(GENERAL.saveButton);

    // View the first version and make sure the secret data is correct
    await click(KV_SECRET.versionDropdown);
    await click(`${KV_SECRET.version(1)} a`);
    assert.strictEqual(codemirror().getValue(), obscuredDataV1, 'Version one data is displayed');

    // Navigate back the second version and make sure the secret data is correct
    await click(KV_SECRET.versionDropdown);
    await click(`${KV_SECRET.version(2)} a`);
    assert.strictEqual(codemirror().getValue(), obscuredDataV2, 'Version two data is displayed');
  });

  test('does not register as advanced when value includes {', async function (assert) {
    await visit(`/vault/secrets/${this.backend}/kv/create`);
    await fillIn(GENERAL.inputByAttr('path'), 'not-advanced');

    await fillIn(KV_FORM.keyInput(), 'foo');
    await fillIn(KV_FORM.maskedValueInput(), '{bar}');
    await click(GENERAL.saveButton);
    await click(KV_SECRET.createNewVersion);
    assert.dom(KV_FORM.toggleJson).isNotDisabled();
    assert.dom(KV_FORM.toggleJson).isNotChecked();
  });
});

// NAMESPACE TESTS
module('Acceptance | Enterprise | kv-v2 workflow | edge cases', function (hooks) {
  setupApplicationTest(hooks);

  const navToEngine = async (backend) => {
    await click('[data-test-sidebar-nav-link="Secrets Engines"]');
    return await click(KV_WORKFLOW.backends.link(backend));
  };

  const assertDeleteActions = (assert, expected = ['delete', 'destroy']) => {
    ['delete', 'destroy', 'undelete'].forEach((toolbar) => {
      if (expected.includes(toolbar)) {
        assert.dom(KV_SECRET[toolbar]).exists(`${toolbar} toolbar action exists`);
      } else {
        assert.dom(KV_SECRET[toolbar]).doesNotExist(`${toolbar} toolbar action not rendered`);
      }
    });
  };

  const assertVersionDropdown = async (assert, deleted = [], versions = [2, 1]) => {
    assert.dom(KV_SECRET.versionDropdown).hasText(`Version ${versions[0]}`);
    await click(KV_SECRET.versionDropdown);
    versions.forEach((num) => {
      assert.dom(KV_SECRET.version(num)).exists(`renders version ${num} link in dropdown`);
    });
    // also asserts destroyed icon
    deleted.forEach((num) => {
      assert.dom(`${KV_SECRET.version(num)} [data-test-icon="x-square"]`);
    });
  };

  // each test uses a different secret path
  hooks.beforeEach(async function () {
    const uid = uuidv4();
    this.store = this.owner.lookup('service:store');
    this.backend = `kv-enterprise-edge-${uid}`;
    this.namespace = `ns-${uid}`;
    await authPage.login();
    await runCmd([`write sys/namespaces/${this.namespace} -force`]);
    return;
  });

  hooks.afterEach(async function () {
    await authPage.login();
    await runCmd([`delete /sys/auth/${this.namespace}`]);
    await runCmd(deleteEngineCmd(this.backend));
    return;
  });

  module('admin persona', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.loginNs(this.namespace);
      // mount engine within namespace
      await runCmd(mountEngineCmd('kv-v2', this.backend), false);
      clearRecords(this.store);
      return;
    });
    hooks.afterEach(async function () {
      // visit logout with namespace query param because we're transitioning from within an engine
      // and navigating directly to /vault/auth caused test context routing problems :(
      await visit(`/vault/logout?namespace=${this.namespace}`);
      await authPage.namespaceInput(''); // clear login form namespace input
    });

    test('namespace: it can create a secret and new secret version', async function (assert) {
      assert.expect(15);
      const backend = this.backend;
      const ns = this.namespace;
      const secret = 'my-create-secret';
      await navToEngine(backend);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list?namespace=${ns}`,
        'navigates to list'
      );
      // Create first version of secret
      await click(KV_LIST.createSecret);
      await fillIn(GENERAL.inputByAttr('path'), secret);
      assert.dom(KV_FORM.toggleMetadata).exists('Shows metadata toggle when creating new secret');
      await fillIn(KV_FORM.keyInput(), 'foo');
      await fillIn(KV_FORM.maskedValueInput(), 'woahsecret');
      await click(GENERAL.saveButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secret}/details?namespace=${ns}&version=1`,
        'navigates to details'
      );

      // Create a new version
      await click(KV_SECRET.createNewVersion);
      assert.dom(GENERAL.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(GENERAL.inputByAttr('path')).hasValue(secret);
      assert
        .dom(KV_FORM.toggleMetadata)
        .doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(KV_FORM.keyInput()).hasValue('foo');
      assert.dom(KV_FORM.maskedValueInput()).hasValue('woahsecret');
      await fillIn(KV_FORM.keyInput(1), 'foo-two');
      await fillIn(KV_FORM.maskedValueInput(1), 'supersecret');
      await click(GENERAL.saveButton);

      // Check details
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secret}/details?namespace=${ns}&version=2`,
        'navigates to details'
      );
      await assertVersionDropdown(assert);
      assert
        .dom(`${KV_SECRET.version(2)} [data-test-icon="check-circle"]`)
        .exists('renders current version icon');
      assert.dom(GENERAL.infoRowValue('foo-two')).hasText('***********');
      await click(KV_WORKFLOW.infoRowToggleMasked('foo-two'));
      assert.dom(GENERAL.infoRowValue('foo-two')).hasText('supersecret', 'secret value shows after toggle');
    });

    test('namespace: it manages state throughout delete, destroy and undelete operations', async function (assert) {
      assert.expect(34);
      const backend = this.backend;
      const ns = this.namespace;
      const secret = 'my-delete-secret';
      await writeVersionedSecret(backend, secret, 'foo', 'bar', 2, ns);
      await navToEngine(backend);

      await click(KV_LIST.item(secret));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secret}/details?namespace=${ns}&version=2`,
        'navigates to details'
      );

      // correct toolbar options & details show
      assertDeleteActions(assert);
      await assertVersionDropdown(assert);
      // delete flow
      await click(KV_SECRET.delete);
      await click(KV_SECRET.deleteOption);
      await click(KV_SECRET.deleteConfirm);
      // check empty state and toolbar
      assertDeleteActions(assert, ['undelete', 'destroy']);
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('Version 2 of this secret has been deleted', 'Shows deleted message');
      assert.dom(KV_SECRET.versionTimestamp).includesText('Version 2 deleted');
      await assertVersionDropdown(assert, [2]); // important to test dropdown versions are accurate

      // navigate to sibling route to make sure empty state remains for details tab
      await click(GENERAL.tab('Version History'));
      assert.dom(KV_WORKFLOW.versions.linkedBlock()).exists({ count: 2 });

      // back to secret tab to confirm deleted state
      await click(GENERAL.tab('Secret'));
      // if this assertion fails, the view is rendering a stale model
      assert.dom(GENERAL.emptyStateTitle).exists('still renders empty state!!');
      await assertVersionDropdown(assert, [2]);

      // undelete flow
      await click(KV_SECRET.undelete);
      // details update accordingly
      assertDeleteActions(assert, ['delete', 'destroy']);
      assert.dom(KV_WORKFLOW.infoRow).exists('shows secret data');
      assert.dom(KV_SECRET.versionTimestamp).includesText('Version 2 created');

      // destroy flow
      await click(KV_SECRET.destroy);
      await click(KV_SECRET.deleteConfirm);
      assertDeleteActions(assert, []);
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('Version 2 of this secret has been permanently destroyed', 'Shows destroyed message');

      // navigate to sibling route to make sure empty state remains for details tab
      await click(GENERAL.tab('Version History'));
      assert.dom(KV_WORKFLOW.versions.linkedBlock()).exists({ count: 2 });

      // back to secret tab to confirm destroyed state
      await click(GENERAL.tab('Secret'));
      // if this assertion fails, the view is rendering a stale model
      assert.dom(GENERAL.emptyStateTitle).exists('still renders empty state!!');
      await assertVersionDropdown(assert, [2]);
    });
  });
});
