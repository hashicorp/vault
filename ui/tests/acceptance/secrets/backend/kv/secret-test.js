/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, visit, settled, currentURL, currentRouteName, fillIn, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { writeSecret, writeVersionedSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SS } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { createSecret } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

const deleteEngine = async function (enginePath, assert) {
  await login();

  const response = await runCmd([`delete sys/mounts/${enginePath}`]);
  assert.strictEqual(
    response,
    `Success! Data deleted (if it existed) at: sys/mounts/${enginePath}`,
    'Engine successfully deleted'
  );
};

module('Acceptance | secrets/secret/create, read, delete', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.uid = uuidv4();
    await login();
  });

  module('mount and configure', function () {
    // no further configuration needed
    test('it can mount a KV 2 secret engine with config metadata', async function (assert) {
      assert.expect(4);
      const enginePath = `kv-secret-${this.uid}`;
      const maxVersion = '101';
      await mountSecrets.visit();
      await click(GENERAL.cardContainer('kv'));
      await fillIn(GENERAL.inputByAttr('path'), enginePath);

      await fillIn(GENERAL.inputByAttr('kv_config.max_versions'), maxVersion);
      await click(GENERAL.inputByAttr('kv_config.cas_required'));
      await click(GENERAL.ttl.toggle('Automate secret deletion'));
      await fillIn(GENERAL.selectByAttr('ttl-unit'), 's');
      await fillIn(GENERAL.ttl.input('Automate secret deletion'), '1');

      await click(GENERAL.submitButton);

      await click(GENERAL.dropdownToggle('Manage'));
      await click(GENERAL.menuItem('Configure'));
      await click(GENERAL.tabLink('plugin-settings'));

      assert
        .dom(PAGE.infoRowValue('Maximum number of versions'))
        .hasText(maxVersion, 'displays the max version set when configuring the secret-engine');
      assert
        .dom(PAGE.infoRowValue('Require check and set'))
        .hasText('Yes', 'displays the cas set when configuring the secret-engine');
      assert
        .dom(PAGE.infoRowValue('Automate secret deletion'))
        .hasText('1 second', 'displays the delete version after set when configuring the secret-engine');
      // [BANDAID] avoid error from missing param for links in SecretEdit > KeyValueHeader
      await visit('/vault/secrets-engines');
      await deleteEngine(enginePath, assert);
    });

    // https://github.com/hashicorp/vault/issues/5994
    test('v1 key named keys', async function (assert) {
      assert.expect(2);
      await runCmd(['vault write sys/mounts/test type=kv', 'refresh', 'vault write test/a keys=a keys=b']);
      await showPage.visit({ backend: 'test', id: 'a' });
      assert.ok(showPage.editIsPresent, 'renders the page properly');
      // [BANDAID] avoid error from missing param for links in SecretEdit > KeyValueHeader
      await visit('/vault/secrets-engines');
      await deleteEngine('test', assert);
    });
  });

  module('kv v2', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = `kvv2-${this.uid}`;
      await runCmd([`write sys/mounts/${this.backend} type=kv options=version=2`]);
    });
    hooks.afterEach(async function () {
      await runCmd([`delete sys/mounts/${this.backend}`]);
    });

    test('it can create a secret when check-and-set is required', async function (assert) {
      const secretPath = 'foo/bar';
      const output = await runCmd(`write ${this.backend}/config cas_required=true`);
      assert.strictEqual(
        output,
        `Success! Data written to: ${this.backend}/config`,
        'Engine successfully updated'
      );
      await visit(`/vault/secrets-engines/kv/list`);
      await writeSecret(this.backend, secretPath, 'foo', 'bar');
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.index',
        'redirects to the overview page'
      );
    });

    test('it navigates to version history and to a specific version', async function (assert) {
      assert.expect(4);
      const secretPath = `specific-version`;
      await writeVersionedSecret(this.backend, secretPath, 'foo', 'bar', 4);
      await click(PAGE.secretTab('Secret'));
      assert
        .dom(PAGE.detail.versionTimestamp)
        .includesText('Version 4 created', 'shows version created time');

      await click(PAGE.secretTab('Version History'));
      assert.dom(PAGE.versions.linkedBlock()).exists({ count: 4 }, 'Lists 4 versions in history');
      assert.dom(PAGE.versions.icon(4)).includesText('Current', 'shows current version on v4');
      await click(PAGE.versions.linkedBlock(2));

      assert.strictEqual(
        currentURL(),
        `/vault/secrets-engines/${this.backend}/kv/${secretPath}/details?version=2`,
        'redirects to the show page with queryParam version=2'
      );
    });
  });

  module('kv v1', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = `kv-v1-${this.uid}`;
      // mount version 1 engine
      await mountSecrets.visit();
      await click(GENERAL.cardContainer('kv'));
      await fillIn(GENERAL.inputByAttr('path'), this.backend);
      await click(GENERAL.button('Method Options'));
      await mountSecrets.version(1);
      await click(GENERAL.submitButton);
    });

    hooks.afterEach(async function () {
      await runCmd([`delete sys/mounts/${this.backend}`]);
    });

    test('version 1 performs the correct capabilities lookup', async function (assert) {
      // TODO: while this should pass it doesn't really do anything anymore for us as v1 and v2 are completely separate.
      const secretPath = 'foo/bar';
      await click(SS.createSecretLink);
      await createSecret(secretPath, 'foo', 'bar');
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.show',
        'redirects to the show page'
      );
      assert.ok(showPage.editIsPresent, 'shows the edit button');
    });

    test('version 1 token without read permissions can create and update a secret', async function (assert) {
      const updatePersonaToken = await runCmd(
        tokenWithPolicyCmd(
          'read-all',
          `
    path "${this.backend}/*" {
      capabilities = ["create", "update", "list"]
    }
    # used to delete the engine after test done in afterEach hook
    path "sys/mounts/${this.backend}" {
      capabilities = ["delete"]
    }
    `
        )
      );

      await login(updatePersonaToken);
      await visit(`/vault/secrets-engines/${this.backend}/list`);
      await click(SS.createSecretLink);
      await createSecret('test', 'foo', 'bar');
      await click('[data-test-secret-edit]', 'can click edit button');
      // edit only without read permissions
      assert
        .dom('[data-test-secret-no-read-permissions] .hds-alert__description')
        .hasText(
          'You do not have read permissions. If a secret exists at this path creating a new secret will overwrite it.',
          'Displays warning about no read permissions'
        );
      await fillIn('[data-test-secret-key]', 'new');
      await fillIn('[data-test-secret-value] textarea', 'new');
      await click(GENERAL.submitButton);
      assert.dom(GENERAL.latestFlashContent).includesText('Secret test updated successfully.');
    });

    // https://github.com/hashicorp/vault/issues/5960
    test('version 1: nested paths creation maintains ability to navigate the tree', async function (assert) {
      const enginePath = this.backend;
      await runCmd([
        `write ${enginePath}/1/2 foo=bar`,
        `write ${enginePath}/1/2/3/4 foo=bar`,
        `write ${enginePath}/1/2/3/4a foo=bar`,
        'refresh',
      ]);
      await settled();
      // navigate to farthest leaf
      await visit(`/vault/secrets-engines/${enginePath}/list`);
      assert.dom('[data-test-component="navigate-input"]').hasNoValue();
      assert.dom(SS.secretLink()).exists({ count: 1 });
      await click(SS.secretLink('1/'));
      assert.dom('[data-test-component="navigate-input"]').hasValue('1/');
      assert.dom(SS.secretLink()).exists({ count: 2 });
      await click(SS.secretLink('1/2/'));
      assert.dom('[data-test-component="navigate-input"]').hasValue('1/2/');
      assert.dom(SS.secretLink()).exists({ count: 1 });
      await click(SS.secretLink('1/2/3/'));
      assert.dom('[data-test-component="navigate-input"]').hasValue('1/2/3/');
      assert.dom(SS.secretLink()).exists({ count: 2 });

      // delete the items
      await click(SS.secretLinkMenu('1/2/3/4'));
      await click(`${SS.secretLink('1/2/3/4')} ${GENERAL.confirmTrigger}`);
      await click(GENERAL.confirmButton);
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.list');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets-engines/${enginePath}/list/1/2/3/`,
        'remains on the page'
      );
      assert.dom(SS.secretLink()).exists({ count: 1 });

      await listPage.secrets.objectAt(0).menuToggle();
      await click(GENERAL.confirmTrigger);
      await click(GENERAL.confirmButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets-engines/${enginePath}/list/1/2/3/`,
        'remains on the page'
      );
      assert.dom(GENERAL.emptyStateTitle).hasText('No secrets under "1/2/3/".');

      await fillIn('[data-test-component="navigate-input"]', '1/2/');
      assert.dom(GENERAL.emptyStateTitle).hasText('No secrets under "1/2/".');

      await click('[data-test-list-root-link]');
      assert.strictEqual(currentURL(), `/vault/secrets-engines/${enginePath}/list`);
      assert.dom(SS.secretLink()).exists({ count: 1 });
    });

    test('first level secrets redirect properly upon deletion', async function (assert) {
      const secretPath = 'test';
      await click(SS.createSecretLink);
      await createSecret(secretPath, 'foo', 'bar');
      await click(GENERAL.confirmTrigger);
      await click(GENERAL.confirmButton);
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.list-root',
        'redirected to the list page on delete'
      );
    });

    test('paths are properly encoded', async function (assert) {
      const backend = this.backend;
      const paths = [
        '(',
        ')',
        '"',
        // "'",
        '!',
        '#',
        '$',
        '&',
        '*',
        '+',
        '@',
        '{',
        '|',
        '}',
        '~',
        '[',
        '\\',
        ']',
        '^',
        '_',
      ].map((char) => `${char}some`);
      assert.expect(paths.length * 2 + 1);
      const secretPath = '2';
      const commands = paths.map((path) => `write '${backend}/${path}/${secretPath}' 3=4`);
      await runCmd([...commands, 'refresh']);
      for (const path of paths) {
        await listPage.visit({ backend, id: path });
        assert.dom(SS.secretLinkATag()).hasText('2', `${path}: secret is displayed properly`);
        await click(SS.secretLink());
        assert.strictEqual(
          currentRouteName(),
          'vault.cluster.secrets.backend.show',
          `${path}: show page renders correctly`
        );
      }
      // [BANDAID] avoid error from missing param for links in SecretEdit > KeyValueHeader
      await visit('/vault/secrets-engines');
      await deleteEngine(backend, assert);
    });

    test('KVv1 handles secret with % in path correctly', async function (assert) {
      const enginePath = this.backend;
      const secretPath = 'per%cent/%fu ll';
      const [firstPath, secondPath] = secretPath.split('/');
      const commands = [`write '${enginePath}/${secretPath}' 3=4`, `refresh`];
      await runCmd(commands);
      await listPage.visitRoot({ backend: enginePath });
      await settled();

      assert.dom(SS.secretLink(`${firstPath}/`)).exists('First section item exists');
      await click(SS.secretLink(`${firstPath}/`));

      assert.strictEqual(
        currentURL(),
        `/vault/secrets-engines/${enginePath}/list/${encodeURIComponent(firstPath)}/`,
        'First part of path is encoded in URL'
      );
      assert.dom(SS.secretLink(secretPath)).exists('Link to secret exists');
      await click(SS.secretLink(secretPath));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets-engines/${enginePath}/show/${encodeURIComponent(firstPath)}/${encodeURIComponent(
          secondPath
        )}`,
        'secret path is encoded in URL'
      );
      assert.dom('h1').hasText(secretPath, 'Path renders correctly on show page');
      await click(SS.crumb(firstPath));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets-engines/${enginePath}/list/${encodeURIComponent(firstPath)}/`,
        'Breadcrumb link encodes correctly'
      );
    });

    // the web cli does not handle a quote as part of a path, so we test it here via the UI
    test('creating a secret with a single or double quote works properly', async function (assert) {
      assert.expect(6);
      const backend = this.backend;
      const paths = ["'some", '"some'];
      for (const path of paths) {
        await listPage.visitRoot({ backend });
        await click(SS.createSecretLink);
        await createSecret(`${path}/2`, 'foo', 'bar');
        await listPage.visit({ backend, id: path });
        assert.dom(SS.secretLinkATag()).hasText('2', `${path}: secret is displayed properly`);
        await click(SS.secretLink());
        assert.strictEqual(
          currentRouteName(),
          'vault.cluster.secrets.backend.show',
          `${path}: show page renders correctly`
        );
        assert.dom(GENERAL.hdsPageHeaderTitle).hasText(`${path}/2`, 'shows correct page title');
      }
    });

    test('filter clears on nav', async function (assert) {
      const backend = this.backend;
      await runCmd([
        `vault write sys/mounts/${backend} type=kv`,
        `refresh`,
        `vault write ${backend}/filter/foo keys=a keys=b`,
        `vault write ${backend}/filter/foo1 keys=a keys=b`,
        `vault write ${backend}/filter/foo2 keys=a keys=b`,
      ]);
      await listPage.visit({ backend, id: 'filter' });
      assert.strictEqual(listPage.secrets.length, 3, 'renders three secrets');
      await listPage.filterInput('filter/foo1');
      assert.strictEqual(listPage.secrets.length, 1, 'renders only one secret');
      await listPage.secrets.objectAt(0).click();
      await click(SS.crumb('filter'));
      assert.strictEqual(listPage.secrets.length, 3, 'renders three secrets');
      assert.strictEqual(listPage.filterInputValue, 'filter/', 'pageFilter has been reset');
    });

    test('it can edit via the JSON input', async function (assert) {
      const content = JSON.stringify({ foo: 'fa', bar: 'boo' });
      const secretPath = `kv-json-${this.uid}`;
      await listPage.visitRoot({ backend: this.backend });
      await click(SS.createSecretLink);
      await fillIn(SS.secretPath('create'), secretPath);
      await click(GENERAL.toggleInput('json'));

      await waitFor('.cm-editor');
      const editor = codemirror();
      setCodeEditorValue(editor, content);

      await click(GENERAL.submitButton);

      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.show',
        'redirects to the show page'
      );
      assert.ok(showPage.editIsPresent, 'shows the edit button');
      assert
        .dom('.hds-code-block')
        .includesText(
          `Secret Data ${JSON.stringify({ bar: 'boo', foo: 'fa' }, null, 2).replace(/\n\s*/g, ' ').trim()}`,
          'shows the secret data'
        );
    });
  });
});
