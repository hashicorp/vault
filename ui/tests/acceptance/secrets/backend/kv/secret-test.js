/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, visit, settled, currentURL, currentRouteName, fillIn } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { writeSecret, writeVersionedSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';

const consoleComponent = create(consoleClass);

const deleteEngine = async function (enginePath, assert) {
  await logout.visit();
  await authPage.login();
  await consoleComponent.runCommands([`delete sys/mounts/${enginePath}`]);
  const response = consoleComponent.lastLogOutput;
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
    await authPage.login();
  });

  module('mount and configure', function () {
    // no further configuration needed
    test('it can mount a KV 2 secret engine with config metadata', async function (assert) {
      assert.expect(4);
      const enginePath = `kv-secret-${this.uid}`;
      const maxVersion = '101';
      await mountSecrets.visit();
      await click('[data-test-mount-type="kv"]');

      await click('[data-test-mount-next]');

      await fillIn('[data-test-input="path"]', enginePath);
      await fillIn('[data-test-input="maxVersions"]', maxVersion);
      await click('[data-test-input="casRequired"]');
      await click('[data-test-toggle-label="Automate secret deletion"]');
      await fillIn('[data-test-select="ttl-unit"]', 's');
      await fillIn('[data-test-ttl-value="Automate secret deletion"]', '1');
      await click('[data-test-mount-submit="true"]');

      await click(PAGE.secretTab('Configuration'));

      assert
        .dom(PAGE.infoRowValue('Maximum number of versions'))
        .hasText(maxVersion, 'displays the max version set when configuring the secret-engine');
      assert
        .dom(PAGE.infoRowValue('Require check and set'))
        .hasText('Yes', 'displays the cas set when configuring the secret-engine');
      assert
        .dom(PAGE.infoRowValue('Automate secret deletion'))
        .hasText('1 second', 'displays the delete version after set when configuring the secret-engine');
      await deleteEngine(enginePath, assert);
    });

    // https://github.com/hashicorp/vault/issues/5994
    test('v1 key named keys', async function (assert) {
      assert.expect(2);
      await consoleComponent.runCommands([
        'vault write sys/mounts/test type=kv',
        'refresh',
        'vault write test/a keys=a keys=b',
      ]);
      await showPage.visit({ backend: 'test', id: 'a' });
      assert.ok(showPage.editIsPresent, 'renders the page properly');
      await deleteEngine('test', assert);
    });
  });

  module('kv v2', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = `kvv2-${this.uid}`;
      await consoleComponent.runCommands([`write sys/mounts/${this.backend} type=kv options=version=2`]);
    });
    hooks.afterEach(async function () {
      await consoleComponent.runCommands([`delete sys/mounts/${this.backend}`]);
    });
    test('it can create a secret when check-and-set is required', async function (assert) {
      const secretPath = 'foo/bar';
      await consoleComponent.runCommands(`write ${this.backend}/config cas_required=true`);
      assert.strictEqual(
        consoleComponent.lastLogOutput,
        `Success! Data written to: ${this.backend}/config`,
        'Engine successfully updated'
      );
      await visit(`/vault/secrets/kv/list`);
      await writeSecret(this.backend, secretPath, 'foo', 'bar');
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.kv.secret.details.index',
        'redirects to the show page'
      );
      assert.dom(PAGE.detail.createNewVersion).exists('shows the edit button');
    });
    test('it navigates to version history and to a specific version', async function (assert) {
      assert.expect(4);
      const secretPath = `specific-version`;
      await writeVersionedSecret(this.backend, secretPath, 'foo', 'bar', 4);
      assert
        .dom(PAGE.detail.versionTimestamp)
        .includesText('Version 4 created', 'shows version created time');

      await click(PAGE.secretTab('Version History'));
      assert.dom(PAGE.versions.linkedBlock()).exists({ count: 4 }, 'Lists 4 versions in history');
      assert.dom(PAGE.versions.icon(4)).includesText('Current', 'shows current version on v4');
      await click(PAGE.versions.linkedBlock(2));

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/kv/${secretPath}/details?version=2`,
        'redirects to the show page with queryParam version=2'
      );
    });
  });

  module('kv v1', function (hooks) {
    hooks.beforeEach(async function () {
      this.backend = `kv-v1-${this.uid}`;
      // mount version 1 engine
      await mountSecrets.visit();
      await mountSecrets.selectType('kv');
      await mountSecrets.next().path(this.backend).toggleOptions().version(1).submit();
    });
    hooks.afterEach(async function () {
      await consoleComponent.runCommands([`delete sys/mounts/${this.backend}`]);
    });
    test('version 1 performs the correct capabilities lookup and does not show metadata tab', async function (assert) {
      const secretPath = 'foo/bar';
      await listPage.create();
      await editPage.createSecret(secretPath, 'foo', 'bar');
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.show',
        'redirects to the show page'
      );
      assert.ok(showPage.editIsPresent, 'shows the edit button');
      // check for metadata tab should not exist on KV version 1
      assert.dom('[data-test-secret-metadata-tab]').doesNotExist('does not show metadata tab');
    });
    // https://github.com/hashicorp/vault/issues/5960
    test('version 1: nested paths creation maintains ability to navigate the tree', async function (assert) {
      const enginePath = this.backend;
      const secretPath = '1/2/3/4';
      await listPage.create();
      await editPage.createSecret(secretPath, 'foo', 'bar');

      // setup an ancestor for when we delete
      await listPage.visitRoot({ backend: enginePath });
      await listPage.secrets.filterBy('text', '1/')[0].click();
      await listPage.create();
      await editPage.createSecret('1/2', 'foo', 'bar');

      // lol we have to do this because ember-cli-page-object doesn't like *'s in visitable
      await listPage.visitRoot({ backend: enginePath });
      await listPage.secrets.filterBy('text', '1/')[0].click();
      await listPage.secrets.filterBy('text', '2/')[0].click();
      await listPage.secrets.filterBy('text', '3/')[0].click();
      await listPage.create();

      await editPage.createSecret(secretPath + 'a', 'foo', 'bar');
      await listPage.visitRoot({ backend: enginePath });
      await listPage.secrets.filterBy('text', '1/')[0].click();
      await listPage.secrets.filterBy('text', '2/')[0].click();
      const secretLink = listPage.secrets.filterBy('text', '3/')[0];
      assert.ok(secretLink, 'link to the 3/ branch displays properly');

      await listPage.secrets.filterBy('text', '3/')[0].click();
      await listPage.secrets.objectAt(0).menuToggle();
      await settled();
      await listPage.delete();
      await listPage.confirmDelete();
      await settled();
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.list');
      assert.strictEqual(currentURL(), `/vault/secrets/${enginePath}/list/1/2/3/`, 'remains on the page');

      await listPage.secrets.objectAt(0).menuToggle();
      await listPage.delete();
      await listPage.confirmDelete();
      await settled();
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.list');
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${enginePath}/list/1/`,
        'navigates to the ancestor created earlier'
      );
    });
    test('first level secrets redirect properly upon deletion', async function (assert) {
      const secretPath = 'test';
      await listPage.create();
      await editPage.createSecret(secretPath, 'foo', 'bar');
      await showPage.deleteSecretV1();
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
        //"'",
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
      await consoleComponent.runCommands([...commands, 'refresh']);
      for (const path of paths) {
        await listPage.visit({ backend, id: path });
        assert.ok(listPage.secrets.filterBy('text', '2')[0], `${path}: secret is displayed properly`);
        await listPage.secrets.filterBy('text', '2')[0].click();
        assert.strictEqual(
          currentRouteName(),
          'vault.cluster.secrets.backend.show',
          `${path}: show page renders correctly`
        );
      }
      await deleteEngine(backend, assert);
    });

    test('UI handles secret with % in path correctly', async function (assert) {
      const enginePath = this.backend;
      const secretPath = 'per%cent/%fu ll';
      const [firstPath, secondPath] = secretPath.split('/');
      const commands = [`write '${enginePath}/${secretPath}' 3=4`, `refresh`];
      await consoleComponent.runCommands(commands);
      await listPage.visitRoot({ backend: enginePath });

      assert.dom(`[data-test-secret-link="${firstPath}/"]`).exists('First section item exists');
      await click(`[data-test-secret-link="${firstPath}/"]`);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${enginePath}/list/${encodeURIComponent(firstPath)}/`,
        'First part of path is encoded in URL'
      );
      assert.dom(`[data-test-secret-link="${secretPath}"]`).exists('Link to secret exists');
      await click(`[data-test-secret-link="${secretPath}"]`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${enginePath}/show/${encodeURIComponent(firstPath)}/${encodeURIComponent(
          secondPath
        )}`,
        'secret path is encoded in URL'
      );
      assert.dom('h1').hasText(secretPath, 'Path renders correctly on show page');
      await click(`[data-test-secret-breadcrumb="${firstPath}"]`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${enginePath}/list/${encodeURIComponent(firstPath)}/`,
        'Breadcrumb link encodes correctly'
      );
    });

    // the web cli does not handle a quote as part of a path, so we test it here via the UI
    test('creating a secret with a single or double quote works properly', async function (assert) {
      assert.expect(6);
      const backend = this.backend;
      // await consoleComponent.runCommands(`write sys/mounts/${backend} type=kv`);
      const paths = ["'some", '"some'];
      for (const path of paths) {
        await listPage.visitRoot({ backend });
        await listPage.create();
        await editPage.createSecret(`${path}/2`, 'foo', 'bar');
        await listPage.visit({ backend, id: path });
        assert.ok(listPage.secrets.filterBy('text', '2')[0], `${path}: secret is displayed properly`);
        await listPage.secrets.filterBy('text', '2')[0].click();
        assert.strictEqual(
          currentRouteName(),
          'vault.cluster.secrets.backend.show',
          `${path}: show page renders correctly`
        );
        assert.dom('h1.title').hasText(`${path}/2`, 'shows correct page title');
      }
    });

    test('filter clears on nav', async function (assert) {
      const backend = this.backend;
      await consoleComponent.runCommands([
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
      await showPage.breadcrumbs.filterBy('text', 'filter')[0].click();
      assert.strictEqual(listPage.secrets.length, 3, 'renders three secrets');
      assert.strictEqual(listPage.filterInputValue, 'filter/', 'pageFilter has been reset');
    });

    test('it can edit via the JSON input', async function (assert) {
      const content = JSON.stringify({ foo: 'fa', bar: 'boo' });
      const secretPath = `kv-json-${this.uid}`;
      await listPage.visitRoot({ backend: this.backend });
      await listPage.create();
      await editPage.path(secretPath).toggleJSON();
      const instance = document.querySelector('.CodeMirror').CodeMirror;
      instance.setValue(content);
      await editPage.save();

      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.show',
        'redirects to the show page'
      );
      assert.ok(showPage.editIsPresent, 'shows the edit button');
      const savedInstance = document.querySelector('.CodeMirror').CodeMirror;
      assert.strictEqual(
        savedInstance.options.value,
        JSON.stringify({ bar: 'boo', foo: 'fa' }, null, 2),
        'saves the content'
      );
    });
  });
});
