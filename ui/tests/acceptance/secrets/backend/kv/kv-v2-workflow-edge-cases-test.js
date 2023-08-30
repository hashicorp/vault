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
} from 'vault/tests/helpers/policy-generator/kv';
import { clearRecords, writeSecret, writeVersionedSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';

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
          'nested-secret-list-reader',
          metadataPolicy({ backend, secretPath, capabilities }) +
            dataPolicy({ backend, secretPath, capabilities })
        ),
        createTokenCmd('nested-secret-list-reader'),
      ]);
      await authPage.login(token);
    });

    test('it can navigate to secrets within a secret directory', async function (assert) {
      assert.expect(19);
      const backend = this.backend;
      const [root, subdirectory, secret] = this.fullSecretPath.split('/');

      await visit(`/vault/secrets/${backend}/kv/list`);
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list`, 'lands on secrets list page');

      await typeIn(PAGE.list.overviewInput, `${root}/`); // add slash because this is a directory
      await click(PAGE.list.overviewButton);

      // URL correct
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${root}%2F/directory`,
        'visits list-directory of root'
      );

      // Title correct
      assert.dom(PAGE.title).hasText(`${backend} Version 2`);
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
      // Secret details visible
      assert.dom(PAGE.title).hasText(this.fullSecretPath);
      assert.dom(PAGE.secretTab('Secret')).hasText('Secret');
      assert.dom(PAGE.secretTab('Secret')).hasClass('active');
      assert.dom(PAGE.secretTab('Metadata')).hasText('Metadata');
      assert.dom(PAGE.secretTab('Metadata')).doesNotHaveClass('active');
      assert.dom(PAGE.secretTab('Version History')).hasText('Version History');
      assert.dom(PAGE.secretTab('Version History')).doesNotHaveClass('active');
      assert.dom(PAGE.toolbarAction).exists({ count: 5 }, 'toolbar renders all actions');
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
        `/vault/secrets/${backend}/kv/${root}%2F${subdirectory}%2F/directory`,
        'goes back to subdirectory list'
      );
      assert.dom(PAGE.list.filter).hasValue(`${root}/${subdirectory}/`);
      assert.dom(PAGE.list.item(secret)).exists('renders linked block for child secret');

      // back again
      previousCrumb = findAll('[data-test-breadcrumbs] li').length - 2;
      await click(PAGE.breadcrumbAtIdx(previousCrumb));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${root}%2F/directory`,
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
        .hasText(`Sorry, we were unable to find any content at /v1/${backend}/data/${root}/${subdirectory}.`);

      assert.dom(PAGE.breadcrumbAtIdx(0)).hasText('secrets');
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
          'destruction-no-read',
          dataPolicy({ backend, secretPath: 'data-delete-only', capabilities: ['delete'] }) +
            deleteVersionsPolicy({ backend, secretPath: 'delete-version-only' }) +
            destroyVersionsPolicy({ backend, secretPath: 'destroy-version-only' }) +
            metadataPolicy({ backend, secretPath: 'destroy-metadata-only', capabilities: ['delete'] }) +
            metadataListPolicy(backend)
        ),
        createTokenCmd('destruction-no-read'),
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
          'Delete metadata? This will permanently delete the metadata and versions of the secret. All version history will be removed. This cannot be undone. Confirm Cancel'
        );
    });
  });
});

module('Acceptance | Enterprise | kv-v2 workflow | edge cases', function (hooks) {
  setupApplicationTest(hooks);

  const navToEngine = async (backend) => {
    await click('[data-test-sidebar-nav-link="Secrets engines"]');
    return await click(PAGE.backends.link(backend));
  };
  hooks.beforeEach(async function () {
    const uid = uuidv4();
    this.store = this.owner.lookup('service:store');
    this.backend = `kv-enterprise-edge-${uid}`;
    this.namespace = `ns-${uid}`;
    this.secretPath = 'my-secret';
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

    test('it can create a new secret version in a namespace', async function (assert) {
      assert.expect(15);
      const backend = this.backend;
      const ns = this.namespace;
      const secret = this.secretPath;
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
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secret}/details?namespace=${ns}&version=1`,
        'navigates to details'
      );

      // Create a new version
      await click(PAGE.detail.createNewVersion);
      assert.dom(FORM.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(FORM.inputByAttr('path')).hasValue(secret);
      assert.dom(FORM.toggleMetadata).doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(FORM.keyInput()).hasValue('foo');
      assert.dom(FORM.maskedValueInput()).hasValue('woahsecret');
      await fillIn(FORM.keyInput(1), 'foo-two');
      await fillIn(FORM.maskedValueInput(1), 'supersecret');
      await click(FORM.saveBtn);

      // Check details
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secret}/details?namespace=${ns}&version=2`,
        'navigates to details'
      );
      assert.dom(PAGE.detail.versionDropdown).hasText('Version 2');
      await click(PAGE.detail.versionDropdown);
      assert.dom(PAGE.detail.version(1)).exists('renders version 1 link in dropdown');
      assert.dom(PAGE.detail.version(2)).exists('renders version 2 link in dropdown');
      assert
        .dom(`${PAGE.detail.version(2)} [data-test-icon="check-circle"]`)
        .exists('renders current version icon');
      assert.dom(PAGE.infoRowValue('foo-two')).hasText('***********');
      await click(PAGE.infoRowToggleMasked('foo-two'));
      assert.dom(PAGE.infoRowValue('foo-two')).hasText('supersecret', 'secret value shows after toggle');
    });
  });
});
