/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
/* eslint-disable no-useless-escape */
import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import {
  click,
  currentURL,
  fillIn,
  findAll,
  setupOnerror,
  typeIn,
  visit,
  triggerKeyEvent,
} from '@ember/test-helpers';
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
import { clearRecords, writeSecret, writeVersionedSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import codemirror from 'vault/tests/helpers/codemirror';
import { personas } from 'vault/tests/helpers/kv/policy-generator';

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
      assert.expect(23);
      const backend = this.backend;
      const [root, subdirectory, secret] = this.fullSecretPath.split('/');

      await visit(`/vault/secrets/${backend}/kv/list`);
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');

      await typeIn(PAGE.list.overviewInput, `${root}/no-access/`);
      assert
        .dom(PAGE.list.overviewButton)
        .hasText('View list', 'shows list and not secret because search is a directory');
      await click(PAGE.list.overviewButton);
      assert.dom(PAGE.emptyStateTitle).hasText(`There are no secrets matching "${root}/no-access/".`);

      await visit(`/vault/secrets/${backend}/kv/list`);
      await typeIn(PAGE.list.overviewInput, `${root}/`); // add slash because this is a directory
      await click(PAGE.list.overviewButton);

      // URL correct
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/${root}/`,
        'visits list-directory of root'
      );

      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} version 2`);
      // Tabs correct
      assert.dom(PAGE.secretTab('Secrets')).hasText('Secrets');
      assert.dom(PAGE.secretTab('Secrets')).hasClass('active');
      assert.dom(PAGE.secretTab('Configuration')).hasText('Configuration');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('active');
      // Toolbar correct
      assert.dom(PAGE.toolbarAction).exists({ count: 1 }, 'toolbar only renders create secret action');
      assert.dom(PAGE.list.filter).hasValue(`${root}/`);
      // List content correct
      assert.dom(PAGE.list.item(`${subdirectory}/`)).exists('renders linked block for subdirectory');
      await click(PAGE.list.item(`${subdirectory}/`));
      assert.dom(PAGE.list.item(secret)).exists('renders linked block for child secret');
      await click(PAGE.list.item(secret));
      assert
        .dom(GENERAL.overviewCard.container('Current version'))
        .hasText(`Current version The current version of this secret. 1`);
      // Secret details visible
      await click(PAGE.secretTab('Secret'));
      assert.dom(PAGE.title).hasText(this.fullSecretPath);
      assert.dom(PAGE.secretTab('Secret')).hasText('Secret');
      assert.dom(PAGE.secretTab('Secret')).hasClass('active');
      assert.dom(PAGE.secretTab('Metadata')).hasText('Metadata');
      assert.dom(PAGE.secretTab('Metadata')).doesNotHaveClass('active');
      assert.dom(PAGE.secretTab('Version History')).hasText('Version History');
      assert.dom(PAGE.secretTab('Version History')).doesNotHaveClass('active');
      assert.dom(PAGE.detail.copy).exists();
      assert.dom(PAGE.detail.versionDropdown).exists();
    });

    test('it navigates back to engine index route via breadcrumbs from secret details', async function (assert) {
      assert.expect(6);
      const backend = this.backend;
      const [root, subdirectory, secret] = this.fullSecretPath.split('/');

      await visit(`vault/secrets/${backend}/kv/${encodeURIComponent(this.fullSecretPath)}/details?version=1`);
      // navigate back through crumbs
      let previousCrumb = findAll('[data-test-breadcrumbs] li').length - 2;
      await click(PAGE.breadcrumbAtIdx(previousCrumb));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/${root}/${subdirectory}/`,
        'goes back to subdirectory list'
      );
      assert.dom(PAGE.list.filter).hasValue(`${root}/${subdirectory}/`);
      assert.dom(PAGE.list.item(secret)).exists('renders linked block for child secret');

      // back again
      previousCrumb = findAll('[data-test-breadcrumbs] li').length - 2;
      await click(PAGE.breadcrumbAtIdx(previousCrumb));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/list/${root}/`,
        'goes back to root directory'
      );
      assert.dom(PAGE.list.item(`${subdirectory}/`)).exists('renders linked block for subdirectory');

      // and back to the engine list view
      previousCrumb = findAll('[data-test-breadcrumbs] li').length - 2;
      await click(PAGE.breadcrumbAtIdx(previousCrumb));
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
      await typeIn(PAGE.list.overviewInput, `${root}/${subdirectory}`); // intentionally leave out trailing slash
      await click(PAGE.list.overviewButton);
      assert.dom(PAGE.error.title).hasText('404 Not Found');
      assert
        .dom(PAGE.error.message)
        .hasText(
          `Sorry, we were unable to find any content at /v1/${backend}/metadata/${root}/${subdirectory}.`
        );

      assert.dom(PAGE.breadcrumbAtIdx(0)).hasText('Secrets');
      assert.dom(PAGE.breadcrumbAtIdx(1)).hasText(backend);
      assert.dom(PAGE.secretTab('Secrets')).doesNotHaveClass('is-active');
      assert.dom(PAGE.secretTab('Configuration')).doesNotHaveClass('is-active');
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

      assert.dom(PAGE.detail.delete).exists('renders delete button');
      await click(PAGE.detail.delete);
      assert
        .dom(PAGE.detail.deleteModal)
        .hasTextContaining('Delete this version This deletes a specific version of the secret');
      assert.dom(PAGE.detail.deleteOption).isDisabled('disables version specific option');
      assert.dom(PAGE.detail.deleteOptionLatest).isEnabled('enables version specific option');
    });

    test('it renders the delete action and disables delete latest version option', async function (assert) {
      assert.expect(4);
      const testSecret = 'delete-version-only';
      await visit(`/vault/secrets/${this.backend}/kv/${testSecret}/details`);

      assert.dom(PAGE.detail.delete).exists('renders delete button');
      await click(PAGE.detail.delete);
      assert
        .dom(PAGE.detail.deleteModal)
        .hasTextContaining('Delete this version This deletes a specific version of the secret');

      assert.dom(PAGE.detail.deleteOption).isEnabled('enables version specific option');
      assert.dom(PAGE.detail.deleteOptionLatest).isDisabled('disables version specific option');
    });

    test('it hides destroy option without version number', async function (assert) {
      assert.expect(1);
      const testSecret = 'destroy-version-only';
      await visit(`/vault/secrets/${this.backend}/kv/${testSecret}/details`);

      assert.dom(PAGE.detail.destroy).doesNotExist();
    });

    test('it renders the destroy metadata action and expected modal copy', async function (assert) {
      assert.expect(2);

      const testSecret = 'destroy-metadata-only';
      await visit(`/vault/secrets/${this.backend}/kv/${testSecret}/metadata`);
      assert.dom(PAGE.metadata.deleteMetadata).exists('renders delete metadata button');
      await click(PAGE.metadata.deleteMetadata);
      assert
        .dom(PAGE.detail.deleteModal)
        .hasText(
          'Delete metadata and secret data? This will permanently delete the metadata and versions of the secret. All version history will be removed. This cannot be undone. Confirm Cancel'
        );
    });
  });

  test('no ghost item after editing metadata', async function (assert) {
    await visit(`/vault/secrets/${this.backend}/kv/list/edge/`);
    assert.dom(PAGE.list.item()).exists({ count: 2 }, 'two secrets are listed');
    await click(PAGE.list.item('two'));
    await click(PAGE.secretTab('Metadata'));
    await click(PAGE.metadata.editBtn);
    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.valueInput(), 'bar');
    await click(FORM.saveBtn);
    await click(PAGE.breadcrumbAtIdx(2));
    assert.dom(PAGE.list.item()).exists({ count: 2 }, 'two secrets are listed');
  });

  test('advanced secret values default to JSON display', async function (assert) {
    await visit(`/vault/secrets/${this.backend}/kv/create`);
    await fillIn(FORM.inputByAttr('path'), 'complex');

    await click(FORM.toggleJson);

    assert.strictEqual(
      codemirror().getValue(),
      `{
  \"\": \"\"
}`,
      'JSON editor displays correct empty object'
    );
    codemirror().setValue('{ "foo3": { "name": "bar3" } }');
    await click(FORM.saveBtn);

    // Details view
    await click(PAGE.secretTab('Secret'));
    assert.dom(FORM.toggleJson).isNotDisabled('JSON toggle is not disabled');
    assert.dom(FORM.toggleJson).isChecked("JSON toggle is checked 'on'");

    assert
      .dom(GENERAL.codeBlock('secret-data'))
      .hasText('Version data { "foo3": { "name": "bar3" } }', 'Values are displayed in the details view');

    // New version view
    await click(PAGE.detail.createNewVersion);
    assert.dom(FORM.toggleJson).isNotDisabled();
    assert.dom(FORM.toggleJson).isChecked();
    assert.deepEqual(
      codemirror().getValue(),
      `{
  "foo3": {
    "name": "bar3"
  }
}`,
      'Values are displayed in the new version view'
    );
  });

  test('on enter the JSON editor cursor goes to the next line', async function (assert) {
    // see issue here: https://github.com/hashicorp/vault/issues/27524
    const predictedCursorPosition = JSON.stringify({ line: 3, ch: 0, sticky: null });
    await visit(`/vault/secrets/${this.backend}/kv/create`);
    await fillIn(FORM.inputByAttr('path'), 'json jump');

    await click(FORM.toggleJson);
    codemirror().setCursor({ line: 2, ch: 1 });
    await triggerKeyEvent(GENERAL.codemirrorTextarea, 'keydown', 'Enter');
    const actualCursorPosition = JSON.stringify(codemirror().getCursor());
    assert.strictEqual(actualCursorPosition, predictedCursorPosition, 'the cursor stayed on the next line');
  });

  test('viewing advanced secret data versions displays the correct version data', async function (assert) {
    assert.expect(2);
    const expectedDataV1 = `{
  "foo1": {
    "name": "bar1"
  }
}`;
    const expectedDataV2 = `{
  "foo2": {
    "name": "bar2"
  }
}`;

    await visit(`/vault/secrets/${this.backend}/kv/create`);
    await fillIn(FORM.inputByAttr('path'), 'complex_version_test');

    await click(FORM.toggleJson);
    codemirror().setValue('{ "foo1": { "name": "bar1" } }');
    await click(FORM.saveBtn);

    // Create another version
    await click(GENERAL.overviewCard.actionText('Create new'));
    codemirror().setValue('{ "foo2": { "name": "bar2" } }');
    await click(FORM.saveBtn);

    // View the first version and make sure the secret data is correct
    await click(PAGE.secretTab('Secret'));
    await click(PAGE.detail.versionDropdown);
    await click(`${PAGE.detail.version(1)} a`);
    assert
      .dom(GENERAL.codeBlock('secret-data'))
      .hasText(`Version data ${expectedDataV1}`, 'Version one data is displayed');

    // Navigate back the second version and make sure the secret data is correct
    await click(PAGE.detail.versionDropdown);
    await click(`${PAGE.detail.version(2)} a`);
    assert
      .dom(GENERAL.codeBlock('secret-data'))
      .hasText(`Version data ${expectedDataV2}`, 'Version two data is displayed');
  });

  test('does not register as advanced when value includes {', async function (assert) {
    await visit(`/vault/secrets/${this.backend}/kv/create`);
    await fillIn(FORM.inputByAttr('path'), 'not-advanced');

    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.maskedValueInput(), '{bar}');
    await click(FORM.saveBtn);
    await click(GENERAL.overviewCard.actionText('Create new'));
    assert.dom(FORM.toggleJson).isNotDisabled();
    assert.dom(FORM.toggleJson).isNotChecked();
  });

  // patch is technically enterprise only but stubbing the version so these tests run on both CE and enterprise
  module('patch-persona', function (hooks) {
    hooks.beforeEach(async function () {
      this.patchSecret = 'patch-secret';
      this.owner.lookup('service:version').type = 'enterprise';
      this.store = this.owner.lookup('service:store');
      await writeSecret(this.backend, this.patchSecret, 'foo', 'bar');
      await writeSecret(this.backend, 'my-destroyed-secret', 'foo', 'bar');
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

    test('it patches a secret from the overview page', async function (assert) {
      await visit(`/vault/secrets/${this.backend}/kv/${this.patchSecret}`);
      assert.dom(GENERAL.overviewCard.content('Subkeys')).hasText('Keys foo');

      await click(GENERAL.overviewCard.actionText('Patch secret'));
      await click(FORM.patchEdit(0));
      await fillIn(FORM.valueInput(0), 'newvalue');
      await fillIn(FORM.keyInput('new'), 'newkey');
      await fillIn(FORM.valueInput('new'), 'newvalue');
      await click(FORM.saveBtn);
      assert.dom(GENERAL.overviewCard.content('Subkeys')).hasText('Keys foo newkey');
    });

    test('it patches a secret from the secret details', async function (assert) {
      await visit(`/vault/secrets/${this.backend}/kv/${this.patchSecret}`);
      assert.dom(GENERAL.overviewCard.content('Subkeys')).hasText('Keys foo');
      await click(PAGE.secretTab('Secret'));
      await click(PAGE.detail.patchLatest);
      await click(FORM.patchEdit(0));
      await fillIn(FORM.valueInput(0), 'newvalue');
      await fillIn(FORM.keyInput('new'), 'newkey');
      await fillIn(FORM.valueInput('new'), 'newvalue');
      await click(FORM.saveBtn);
      assert.dom(GENERAL.overviewCard.content('Subkeys')).hasText('Keys foo newkey');
    });

    // testing both adding and deleting a key here because the writeSecret helper only creates a single key/value pair
    test('it adds and deletes a key', async function (assert) {
      await visit(`/vault/secrets/${this.backend}/kv/${this.patchSecret}`);
      // add a new key
      assert.dom(GENERAL.overviewCard.content('Subkeys')).hasText('Keys foo');
      await click(GENERAL.overviewCard.actionText('Patch secret'));
      await fillIn(FORM.keyInput('new'), 'newkey');
      await fillIn(FORM.valueInput('new'), 'newvalue');
      await click(FORM.saveBtn);
      assert.dom(GENERAL.overviewCard.content('Subkeys')).hasText('Keys foo newkey');

      // deletes a key
      await click(GENERAL.overviewCard.actionText('Patch secret'));
      await click(FORM.patchDelete());
      await click(FORM.saveBtn);
      assert.dom(GENERAL.overviewCard.content('Subkeys')).hasText('Keys newkey');
    });

    test('patching a destroyed secret is not allowed', async function (assert) {
      assert.expect(5);
      const secret = 'my-destroyed-secret';
      await visit(`/vault/secrets/${this.backend}/kv/${secret}`);
      assert.dom(GENERAL.overviewCard.actionText('Patch secret')).exists();
      await click(PAGE.secretTab('Secret'));
      assert.dom(PAGE.detail.patchLatest).exists();
      await click(PAGE.detail.destroy);
      await click(PAGE.detail.deleteConfirm);
      // check overview
      assert
        .dom(GENERAL.overviewCard.actionText('Patch secret'))
        .doesNotExist('overview patch action is hidden for destroyed versions');
      await click(PAGE.secretTab('Secret'));
      // check secret tab
      assert
        .dom(PAGE.detail.patchLatest)
        .doesNotExist('toolbar patch action is hidden for destroyed versions');
      // check navigating directly
      await visit(`/vault/secrets/${this.backend}/kv/${secret}/patch`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/kv/${secret}`,
        'destroyed secrets redirect'
      );
    });
  });
});

// NAMESPACE TESTS
module('Acceptance | Enterprise | kv-v2 workflow | edge cases', function (hooks) {
  setupApplicationTest(hooks);

  const navToEngine = async (backend) => {
    await click('[data-test-sidebar-nav-link="Secrets Engines"]');
    return await click(PAGE.backends.link(backend));
  };

  const assertDeleteActions = (assert, expected = ['delete', 'destroy']) => {
    ['delete', 'destroy', 'undelete'].forEach((toolbar) => {
      if (expected.includes(toolbar)) {
        assert.dom(PAGE.detail[toolbar]).exists(`${toolbar} toolbar action exists`);
      } else {
        assert.dom(PAGE.detail[toolbar]).doesNotExist(`${toolbar} toolbar action not rendered`);
      }
    });
  };

  const assertVersionDropdown = async (assert, deleted = [], versions = [2, 1]) => {
    assert.dom(PAGE.detail.versionDropdown).hasText(`Version ${versions[0]}`);
    await click(PAGE.detail.versionDropdown);
    versions.forEach((num) => {
      assert.dom(PAGE.detail.version(num)).exists(`renders version ${num} link in dropdown`);
    });
    // also asserts destroyed icon
    deleted.forEach((num) => {
      assert.dom(`${PAGE.detail.version(num)} [data-test-icon="x-square"]`);
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
      assert.expect(16);
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
      await click(PAGE.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), secret);
      assert.dom(FORM.toggleMetadata).exists('Shows metadata toggle when creating new secret');
      await fillIn(FORM.keyInput(), 'foo');
      await fillIn(FORM.maskedValueInput(), 'woahsecret');
      await click(FORM.saveBtn);
      assert
        .dom(GENERAL.overviewCard.container('Current version'))
        .hasText(`Current version Create new The current version of this secret. 1`);

      // Create a new version
      await click(GENERAL.overviewCard.actionText('Create new'));
      assert.dom(FORM.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(FORM.inputByAttr('path')).hasValue(secret);
      assert.dom(FORM.toggleMetadata).doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(FORM.keyInput()).hasValue('foo');
      assert.dom(FORM.maskedValueInput()).hasValue('woahsecret');
      await fillIn(FORM.keyInput(1), 'foo-two');
      await fillIn(FORM.maskedValueInput(1), 'supersecret');
      await click(FORM.saveBtn);
      assert
        .dom(GENERAL.overviewCard.container('Current version'))
        .hasText(`Current version Create new The current version of this secret. 2`);

      // Check details
      await click(PAGE.secretTab('Secret'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secret}/details?namespace=${ns}&version=2`,
        'navigates to details'
      );
      await assertVersionDropdown(assert);
      assert
        .dom(`${PAGE.detail.version(2)} [data-test-icon="check-circle"]`)
        .exists('renders current version icon');
      assert.dom(PAGE.infoRowValue('foo-two')).hasText('***********');
      await click(PAGE.infoRowToggleMasked('foo-two'));
      assert.dom(PAGE.infoRowValue('foo-two')).hasText('supersecret', 'secret value shows after toggle');
    });

    test('namespace: it manages state throughout delete, destroy and undelete operations', async function (assert) {
      assert.expect(36);
      const backend = this.backend;
      const ns = this.namespace;
      const secret = 'my-delete-secret';
      await writeVersionedSecret(backend, secret, 'foo', 'bar', 2, ns);
      await navToEngine(backend);

      await click(PAGE.list.item(secret));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secret}?namespace=${ns}`,
        'navigates to overview'
      );

      // correct toolbar options & details show
      await click(PAGE.secretTab('Secret'));
      assertDeleteActions(assert);
      await assertVersionDropdown(assert);
      // delete flow
      await click(PAGE.detail.delete);
      await click(PAGE.detail.deleteOption);
      await click(PAGE.detail.deleteConfirm);
      assert
        .dom(GENERAL.overviewCard.container('Current version'))
        .hasTextContaining(
          'Current version Deleted Create new The current version of this secret was deleted'
        );

      await click(PAGE.secretTab('Secret'));
      // check empty state and toolbar
      assertDeleteActions(assert, ['undelete', 'destroy']);
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('Version 2 of this secret has been deleted', 'Shows deleted message');
      assert.dom(PAGE.detail.versionTimestamp).includesText('Version 2 deleted');
      await assertVersionDropdown(assert, [2]); // important to test dropdown versions are accurate

      // navigate to sibling route to make sure empty state remains for details tab
      await click(PAGE.secretTab('Version History'));
      assert.dom(PAGE.versions.linkedBlock()).exists({ count: 2 });

      // back to secret tab to confirm deleted state
      await click(PAGE.secretTab('Secret'));
      // if this assertion fails, the view is rendering a stale model
      assert.dom(PAGE.emptyStateTitle).exists('still renders empty state!!');
      await assertVersionDropdown(assert, [2]);

      // undelete flow
      await click(PAGE.detail.undelete);
      assert
        .dom(GENERAL.overviewCard.container('Current version'))
        .hasTextContaining('Current version Create new The current version of this secret.');
      // details update accordingly
      await click(PAGE.secretTab('Secret'));
      assertDeleteActions(assert, ['delete', 'destroy']);
      assert.dom(PAGE.infoRow).exists('shows secret data');
      assert.dom(PAGE.detail.versionTimestamp).includesText('Version 2 created');

      // destroy flow
      await click(PAGE.detail.destroy);
      await click(PAGE.detail.deleteConfirm);
      await click(PAGE.secretTab('Secret'));
      assertDeleteActions(assert, []);
      assert
        .dom(PAGE.emptyStateTitle)
        .hasText('Version 2 of this secret has been permanently destroyed', 'Shows destroyed message');

      // navigate to sibling route to make sure empty state remains for details tab
      await click(PAGE.secretTab('Version History'));
      assert.dom(PAGE.versions.linkedBlock()).exists({ count: 2 });

      // back to secret tab to confirm destroyed state
      await click(PAGE.secretTab('Secret'));
      // if this assertion fails, the view is rendering a stale model
      assert.dom(PAGE.emptyStateTitle).exists('still renders empty state!!');
      await assertVersionDropdown(assert, [2]);
    });
  });
});
