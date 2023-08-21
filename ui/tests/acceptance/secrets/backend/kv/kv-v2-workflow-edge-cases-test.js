import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentURL, findAll, setupOnerror, typeIn, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import {
  createPolicyCmd,
  deleteEngineCmd,
  mountEngineCmd,
  runCmd,
  createTokenCmd,
} from 'vault/tests/helpers/commands';
import { dataPolicy, metadataPolicy } from 'vault/tests/helpers/policy-generator/kv';
import { writeSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';

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
      assert.dom(PAGE.secretTab('list-directory')).hasText('Secrets');
      assert.dom(PAGE.secretTab('list-directory')).hasClass('active');
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
});
