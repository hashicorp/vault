/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import {
  click,
  visit,
  settled,
  currentURL,
  currentRouteName,
  fillIn,
  triggerKeyEvent,
  typeIn,
} from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { module, skip, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';

const consoleComponent = create(consoleClass);

const writeSecret = async function (backend, path, key, val) {
  await listPage.visitRoot({ backend });
  await listPage.create();
  return editPage.createSecret(path, key, val);
};

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

const mountEngineGeneratePolicyToken = async (enginePath, secretPath, policy, version = 2) => {
  await consoleComponent.runCommands([
    // delete any kv previously written here so that tests can be re-run
    `delete ${enginePath}/metadata/${secretPath}`,
    // delete any previous mount with same name
    `delete sys/mounts/${enginePath}`,
    // mount engine and generate policy
    `write sys/mounts/${enginePath} type=kv options=version=${version}`,
    `write sys/policies/acl/kv-v2-test-policy policy=${btoa(policy)}`,
    'write -field=client_token auth/token/create policies=kv-v2-test-policy',
  ]);
  return consoleComponent.lastLogOutput;
};

module('Acceptance | secrets/secret/create, read, delete', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.uid = uuidv4();
    this.server = apiStub({ usePassthrough: true });
    await authPage.login();
  });

  hooks.afterEach(async function () {
    this.server.shutdown();
  });

  test('it creates a secret and redirects', async function (assert) {
    assert.expect(5);
    const secretPath = `kv-path-${this.uid}`;
    const path = `kv-engine-${this.uid}`;
    await enablePage.enable('kv', path);
    await listPage.visitRoot({ backend: path });
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'navigates to the list page'
    );
    await listPage.create();
    await settled();
    await editPage.toggleMetadata();
    await settled();
    assert.ok(editPage.hasMetadataFields, 'shows the metadata form');
    await editPage.createSecret(secretPath, 'foo', 'bar');
    await settled();

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    await deleteEngine(path, assert);
  });

  test('it can create a secret when check-and-set is required', async function (assert) {
    assert.expect(3);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = 'foo/bar';
    await mountSecrets.visit();
    await mountSecrets.enable('kv', enginePath);
    await consoleComponent.runCommands(`write ${enginePath}/config cas_required=true`);
    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    await deleteEngine(enginePath, assert);
  });

  test('it can create a secret with a non default max version and add metadata', async function (assert) {
    assert.expect(4);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = 'maxVersions';
    const maxVersions = 101;
    await mountSecrets.visit();
    await mountSecrets.enable('kv', enginePath);
    await settled();
    await editPage.startCreateSecret();
    await editPage.path(secretPath);
    await editPage.toggleMetadata();
    await settled();
    await editPage.maxVersion(maxVersions);
    await settled();
    await editPage.save();
    await settled();
    await editPage.metadataTab();
    await settled();
    const savedMaxVersions = Number(
      document.querySelector('[data-test-value-div="Maximum versions"]').innerText
    );
    assert.strictEqual(
      maxVersions,
      savedMaxVersions,
      'max_version displays the saved number set when creating the secret'
    );
    // add metadata
    await click('[data-test-add-custom-metadata]');
    await fillIn('[data-test-kv-key]', 'key');
    await fillIn('[data-test-kv-value]', 'value');
    await click('[data-test-save-metadata]');
    const key = document.querySelector('[data-test-row-label="key"]').innerText;
    const value = document.querySelector('[data-test-row-value="key"]').innerText;
    assert.strictEqual(key, 'key', 'metadata key displays after adding it.');
    assert.strictEqual(value, 'value', 'metadata value displays after adding it.');
    await deleteEngine(enginePath, assert);
  });

  skip('it can handle validation on custom metadata', async function (assert) {
    assert.expect(3);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = 'customMetadataValidations';

    await mountSecrets.visit();
    await mountSecrets.enable('kv', enginePath);
    await settled();
    await editPage.startCreateSecret();
    await editPage.path(secretPath);
    await editPage.toggleMetadata();
    await settled();
    await typeIn('[data-test-kv-value]', 'invalid\\/');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('Custom values cannot contain a backward slash.', 'will not allow backward slash in value.');
    await fillIn('[data-test-kv-value]', ''); // clear previous contents
    await typeIn('[data-test-kv-value]', 'removed!');
    assert.dom('[data-test-inline-error-message]').doesNotExist('inline error goes away');
    await click('[data-test-secret-save]');
    assert
      .dom('[data-test-message-error]')
      .includesText(
        'custom_metadata validation failed: length of key',
        'shows API error that is not captured by validation'
      );
    await deleteEngine(enginePath, assert);
  });

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

    await click('[data-test-configuration-tab]');

    const cas = document.querySelector('[data-test-value-div="Require Check and Set"]').innerText;
    const deleteVersionAfter = document.querySelector(
      '[data-test-value-div="Automate secret deletion"]'
    ).innerText;
    const savedMaxVersion = document.querySelector(
      '[data-test-value-div="Maximum number of versions"]'
    ).innerText;

    assert.strictEqual(
      maxVersion,
      savedMaxVersion,
      'displays the max version set when configuring the secret-engine'
    );
    assert.strictEqual(cas.trim(), 'Yes', 'displays the cas set when configuring the secret-engine');
    assert.strictEqual(
      deleteVersionAfter.trim(),
      '1 second',
      'displays the delete version after set when configuring the secret-engine'
    );
    await deleteEngine(enginePath, assert);
  });

  test('it can create a secret and metadata can be created and edited', async function (assert) {
    assert.expect(2);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = 'metadata';
    const maxVersions = 101;
    await mountSecrets.visit();
    await mountSecrets.enable('kv', enginePath);
    await settled();
    await editPage.startCreateSecret();
    await editPage.path(secretPath);
    await editPage.toggleMetadata();
    await settled();
    await fillIn('[data-test-input="maxVersions"]', maxVersions);

    await editPage.save();
    await settled();
    await editPage.metadataTab();
    await settled();
    const savedMaxVersions = Number(document.querySelectorAll('[data-test-value-div]')[0].innerText);
    assert.strictEqual(
      maxVersions,
      savedMaxVersions,
      'max_version displays the saved number set when creating the secret'
    );
    await deleteEngine(enginePath, assert);
  });

  test('it disables save when validation errors occur', async function (assert) {
    assert.expect(5);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = 'not-duplicate';
    await mountSecrets.visit();
    await mountSecrets.enable('kv', enginePath);
    await settled();
    await editPage.startCreateSecret();
    await typeIn('[data-test-secret-path="true"]', 'beep');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'A secret with this path already exists.',
        'when duplicate path it shows correct error message'
      );

    await editPage.toggleMetadata();
    await settled();
    await typeIn('[data-test-input="maxVersions"]', 'abc');
    assert
      .dom('[data-test-input="maxVersions"]')
      .hasClass('has-error-border', 'shows border error on input with error');
    assert.dom('[data-test-secret-save]').isDisabled('Save button is disabled');
    await fillIn('[data-test-input="maxVersions"]', 20); // fillIn replaces the text, whereas typeIn only adds to it.
    await triggerKeyEvent('[data-test-input="maxVersions"]', 'keyup', 65);
    await editPage.path(secretPath);
    await triggerKeyEvent('[data-test-secret-path="true"]', 'keyup', 65);
    await click('[data-test-secret-save]');
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath}/show/${secretPath}`,
      'navigates to show secret'
    );
    await deleteEngine(enginePath, assert);
  });

  test('it navigates to version history and to a specific version', async function (assert) {
    assert.expect(6);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = `specific-version`;
    await mountSecrets.visit();
    await mountSecrets.enable('kv', enginePath);
    await settled();
    await listPage.visitRoot({ backend: enginePath });
    await settled();
    await listPage.create();
    await settled();
    await editPage.createSecret(secretPath, 'foo', 'bar');
    await click('[data-test-popup-menu-trigger="version"]');

    assert.dom('[data-test-created-time]').includesText('Version created ', 'shows version created time');

    await click('[data-test-version-history]');

    assert
      .dom('[data-test-list-item-content]')
      .includesText('Version 1 Current', 'shows version one data on the version history as current');
    assert.dom('[data-test-list-item-content]').exists({ count: 1 }, 'renders a single version');

    await click('.linked-block');
    await click('button.button.masked-input-toggle');
    assert.dom('[data-test-masked-input]').hasText('bar', 'renders secret on the secret version show page');
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath}/show/${secretPath}?version=1`,
      'redirects to the show page with queryParam version=1'
    );
    await deleteEngine(enginePath, assert);
  });

  test('version 1 performs the correct capabilities lookup and does not show metadata tab', async function (assert) {
    assert.expect(4);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = 'foo/bar';
    // mount version 1 engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets.next().path(enginePath).toggleOptions().version(1).submit();
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
    await deleteEngine(enginePath, assert);
  });

  // https://github.com/hashicorp/vault/issues/5960
  test('version 1: nested paths creation maintains ability to navigate the tree', async function (assert) {
    assert.expect(6);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = '1/2/3/4';
    // mount version 1 engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets.next().path(enginePath).toggleOptions().version(1).submit();
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
    await deleteEngine(enginePath, assert);
  });

  test('first level secrets redirect properly upon deletion', async function (assert) {
    assert.expect(2);
    const enginePath = `kv-secret-${this.uid}`;
    const secretPath = 'test';
    // mount version 1 engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets.next().path(enginePath).toggleOptions().version(1).submit();
    await listPage.create();
    await editPage.createSecret(secretPath, 'foo', 'bar');
    await showPage.deleteSecretV1();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'redirected to the list page on delete'
    );
    await deleteEngine(enginePath, assert);
  });

  // https://github.com/hashicorp/vault/issues/5994
  test('version 1: key named keys', async function (assert) {
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

  test('it redirects to the path ending in / for list pages', async function (assert) {
    assert.expect(3);
    const secretPath = `foo/bar/kv-list-${this.uid}`;
    await consoleComponent.runCommands(['vault write sys/mounts/secret type=kv']);
    await listPage.visitRoot({ backend: 'secret' });
    await listPage.create();
    await editPage.createSecret(secretPath, 'foo', 'bar');
    await settled();
    // use visit helper here because ids with / in them get encoded
    await visit('/vault/secrets/secret/list/foo/bar');
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.list');
    assert.ok(currentURL().endsWith('/'), 'redirects to the path ending in a slash');
    await deleteEngine('secret', assert);
  });

  test('it can edit via the JSON input', async function (assert) {
    assert.expect(4);
    const content = JSON.stringify({ foo: 'fa', bar: 'boo' });
    const secretPath = `kv-json-${this.uid}`;
    await consoleComponent.runCommands(['vault write sys/mounts/secret type=kv']);
    await listPage.visitRoot({ backend: 'secret' });
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
    await deleteEngine('secret', assert);
  });

  test('paths are properly encoded', async function (assert) {
    const backend = `kv-encoding-${this.uid}`;
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
    await consoleComponent.runCommands([`write sys/mounts/${backend} type=kv`, ...commands]);
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

  test('create secret with space shows version data and shows space warning', async function (assert) {
    assert.expect(4);
    const enginePath = `kv-engine-${this.uid}`;
    const secretPath = 'space space';
    // mount version 2
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets.next().path(enginePath).submit();
    await settled();
    await listPage.create();
    await editPage.createSecretDontSave(secretPath, 'foo', 'bar');
    // to trigger warning need to hit keyup on the secret path
    await triggerKeyEvent('[data-test-secret-path="true"]', 'keyup', 65);

    assert.dom('[data-test-whitespace-warning]').exists('renders warning about their being a space');
    await settled();
    await click('[data-test-secret-save]');

    await click('[data-test-popup-menu-trigger="version"]');

    await click('[data-test-version-history]');

    assert.dom('[data-test-list-item-content]').exists('renders the version and not an error state');
    // click on version
    await click('[data-test-popup-menu-trigger="true"]');
    await click('[data-test-version]');

    // perform encode function that should be done by the encodePath
    const encodedSecretPath = secretPath.replace(/ /g, '%20');
    assert.strictEqual(currentURL(), `/vault/secrets/${enginePath}/show/${encodedSecretPath}?version=1`);
    await deleteEngine(enginePath, assert);
  });

  test('UI handles secret with % in path correctly', async function (assert) {
    assert.expect(7);
    const enginePath = `kv-engine-${this.uid}`;
    const secretPath = 'per%cent/%fu ll';
    const [firstPath, secondPath] = secretPath.split('/');
    const commands = [`write sys/mounts/${enginePath} type=kv`, `write '${enginePath}/${secretPath}' 3=4`];
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
      `/vault/secrets/${enginePath}/show/${encodeURIComponent(firstPath)}/${encodeURIComponent(secondPath)}`,
      'secret path is encoded in URL'
    );
    assert.dom('h1').hasText(secretPath, 'Path renders correctly on show page');
    await click(`[data-test-secret-breadcrumb="${firstPath}"]`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath}/list/${encodeURIComponent(firstPath)}/`,
      'Breadcrumb link encodes correctly'
    );
    await deleteEngine(enginePath, assert);
  });

  // the web cli does not handle a quote as part of a path, so we test it here via the UI
  test('creating a secret with a single or double quote works properly', async function (assert) {
    assert.expect(5);
    const backend = `kv-quotes-${this.uid}`;
    await consoleComponent.runCommands(`write sys/mounts/${backend} type=kv`);
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
    }
    await deleteEngine(backend, assert);
  });

  test('filter clears on nav', async function (assert) {
    assert.expect(5);
    const backend = 'test';
    await consoleComponent.runCommands([
      'vault write sys/mounts/test type=kv',
      'refresh',
      'vault write test/filter/foo keys=a keys=b',
      'vault write test/filter/foo1 keys=a keys=b',
      'vault write test/filter/foo2 keys=a keys=b',
    ]);
    await listPage.visit({ backend, id: 'filter' });
    assert.strictEqual(listPage.secrets.length, 3, 'renders three secrets');
    await listPage.filterInput('filter/foo1');
    assert.strictEqual(listPage.secrets.length, 1, 'renders only one secret');
    await listPage.secrets.objectAt(0).click();
    await showPage.breadcrumbs.filterBy('text', 'filter')[0].click();
    assert.strictEqual(listPage.secrets.length, 3, 'renders three secrets');
    assert.strictEqual(listPage.filterInputValue, 'filter/', 'pageFilter has been reset');
    await deleteEngine(backend, assert);
  });

  // All policy tests below this line
  test('version 2 with restricted policy still allows creation and does not show metadata tab', async function (assert) {
    assert.expect(4);
    const enginePath = 'dont-show-metadata-tab';
    const secretPath = 'dont-show-metadata-tab-secret-path';
    const V2_POLICY = `
      path "${enginePath}/metadata/*" {
        capabilities = ["list"]
      }
      path "${enginePath}/data/${secretPath}" {
        capabilities = ["create", "read", "update"]
      }
    `;
    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, V2_POLICY);
    await logout.visit();
    await authPage.login(userToken);

    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.ok(showPage.editIsPresent, 'shows the edit button');
    //check for metadata tab which should not show because you don't have read capabilities
    assert.dom('[data-test-secret-metadata-tab]').doesNotExist('does not show metadata tab');
    await deleteEngine(enginePath, assert);
  });

  test('version 2 with no access to data but access to metadata shows metadata tab', async function (assert) {
    assert.expect(5);
    const enginePath = 'kv-metadata-access-only';
    const secretPath = 'nested/kv-metadata-access-only-secret-name';
    const V2_POLICY = `
      path "${enginePath}/metadata/nested/*" {
        capabilities = ["read", "update"]
      }
    `;

    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, V2_POLICY);
    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    await logout.visit();
    await authPage.login(userToken);
    await settled();
    await visit(`/vault/secrets/${enginePath}/show/${secretPath}`);
    assert.dom('[data-test-empty-state-title]').hasText('You do not have permission to read this secret.');
    assert.dom('[data-test-secret-metadata-tab]').exists('Metadata tab exists');
    await editPage.metadataTab();
    await settled();
    assert.dom('[data-test-empty-state-title]').hasText('No custom metadata');
    assert.dom('[data-test-add-custom-metadata]').exists('it shows link to edit metadata');

    await deleteEngine(enginePath, assert);
  });

  // TODO VAULT-16258: revisit when KV-V2 is engine
  test.skip('version 2: with metadata no read or list but with delete access and full access to the data endpoint', async function (assert) {
    assert.expect(12);
    const enginePath = 'no-metadata-read';
    const secretPath = 'no-metadata-read-secret-name';
    const V2_POLICY_NO_LIST = `
      path "${enginePath}/metadata/*" {
        capabilities = ["update","delete"]
      }
      path "${enginePath}/data/*" {
        capabilities = ["create", "read", "update", "delete"]
      }
    `;
    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, V2_POLICY_NO_LIST);
    await listPage.visitRoot({ backend: enginePath });
    // confirm they see an empty state and not the get-credentials card
    assert.dom('[data-test-empty-state-title]').hasText('No secrets in this backend');
    await settled();
    await listPage.create();
    await settled();
    await editPage.createSecretWithMetadata(secretPath, 'secret-key', 'secret-value', 101);
    await settled();
    await logout.visit();
    await settled();
    await authPage.login(userToken);
    await settled();
    // test if metadata tab there with no read access message and no ability to edit.
    await click(`[data-test-auth-backend-link=${enginePath}]`);
    assert
      .dom('[data-test-get-credentials]')
      .exists(
        'They do not have list access so when logged in under the restricted policy they see the get-credentials-card'
      );

    await visit(`/vault/secrets/${enginePath}/show/${secretPath}`);

    assert
      .dom('[data-test-value-div="secret-key"]')
      .exists('secret view page and info table row with secret-key value');

    // Create new version
    assert.dom('[data-test-secret-edit]').doesNotHaveClass('disabled', 'Create new version is not disabled');
    await click('[data-test-secret-edit]');

    // create new version should not include version in the URL
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath}/edit/${secretPath}`,
      'edit route does not include version query param'
    );
    // Update key
    await editPage.secretKey('newKey');
    await editPage.secretValue('some-value');
    await editPage.save();
    assert.dom('[data-test-value-div="newKey"]').exists('Info row table exists at newKey');

    // check metadata tab
    await click('[data-test-secret-metadata-tab]');

    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'In order to edit secret metadata access, the UI requires read permissions; otherwise, data may be deleted. Edits can still be made via the API and CLI.'
      );
    // destroy the version
    await click('[data-test-secret-tab]');

    await click('[data-test-delete-open-modal]');

    assert.dom('.modal.is-active').exists('Modal appears');
    assert.dom('[data-test-delete-modal="destroy-all-versions"]').exists(); // we have a if Ember.testing catch in the delete action because it breaks things in testing
    // we can however destroy the versions
    await click('#destroy-all-versions');

    await click('[data-test-modal-delete]');

    assert.strictEqual(currentURL(), `/vault/secrets/${enginePath}/list`, 'brings you back to the list page');
    await visit(`/vault/secrets/${enginePath}/show/${secretPath}`);

    assert.dom('[data-test-secret-not-found]').exists('secret no longer found');
    await deleteEngine(enginePath, assert);
  });

  // KV delete operations testing
  test('version 2 with policy with destroy capabilities shows modal', async function (assert) {
    assert.expect(5);
    const enginePath = 'kv-v2-destroy-capabilities';
    const secretPath = 'kv-v2-destroy-capabilities-secret-path';
    const V2_POLICY = `
      path "${enginePath}/destroy/*" {
        capabilities = ["update"]
      }
      path "${enginePath}/metadata/*" {
        capabilities = ["list", "update", "delete"]
      }
      path "${enginePath}/data/${secretPath}" {
        capabilities = ["create", "read", "update"]
      }
    `;
    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, V2_POLICY);
    await logout.visit();
    await authPage.login(userToken);

    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    await click('[data-test-delete-open-modal]');

    assert.dom('[data-test-delete-modal="destroy-version"]').exists('destroy this version option shows');
    assert.dom('[data-test-delete-modal="destroy-all-versions"]').exists('destroy all versions option shows');
    assert.dom('[data-test-delete-modal="delete-version"]').doesNotExist('delete version does not show');

    // because destroy requires a page refresh (making the test suite run in a loop) this action is caught in ember testing and does not refresh.
    // therefore to show new state change after modal closes we jump to the metadata tab and then back.
    await click('#destroy-version');
    await settled(); // eslint-disable-line
    await click('[data-test-modal-delete]');
    await settled(); // eslint-disable-line
    await click('[data-test-secret-metadata-tab]');
    await settled(); // eslint-disable-line
    await click('[data-test-secret-tab]');
    await settled(); // eslint-disable-line
    assert
      .dom('[data-test-empty-state-title]')
      .includesText('Version 1 of this secret has been permanently destroyed');
    await deleteEngine(enginePath, assert);
  });

  test('version 2 with policy with only delete option does not show modal and undelete is an option', async function (assert) {
    assert.expect(5);
    const enginePath = 'kv-v2-only-delete';
    const secretPath = 'kv-v2-only-delete-secret-path';
    const V2_POLICY = `
      path "${enginePath}/delete/*" {
        capabilities = ["update"]
      }
      path "${enginePath}/undelete/*" {
        capabilities = ["update"]
      }
      path "${enginePath}/metadata/*" {
        capabilities = ["list","read","create","update"]
      }
      path "${enginePath}/data/${secretPath}" {
        capabilities = ["create", "read"]
      }
    `;
    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, V2_POLICY);
    await logout.visit();
    await authPage.login(userToken);
    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    assert.dom('[data-test-delete-open-modal]').doesNotExist('delete version does not show');
    assert.dom('[data-test-secret-v2-delete="true"]').exists('drop down delete shows');
    await showPage.deleteSecretV2();
    // unable to reload page in test scenario so going to list and back to secret to confirm deletion
    const url = `/vault/secrets/${enginePath}/list`;
    await visit(url);

    await click(`[data-test-secret-link="${secretPath}"]`);
    await settled(); // eslint-disable-line
    assert.dom('[data-test-component="empty-state"]').exists('secret has been deleted');
    assert.dom('[data-test-secret-undelete]').exists('undelete button shows');
    await deleteEngine(enginePath, assert);
  });

  test('version 2: policy includes "delete" capability for secret path but does not have "update" to /delete endpoint', async function (assert) {
    assert.expect(4);
    const enginePath = 'kv-v2-soft-delete-only';
    const secretPath = 'kv-v2-delete-capability-not-path';
    const policy = `
      path "${enginePath}/data/${secretPath}" { capabilities = ["create","read","update","delete","list"] }
      path "${enginePath}/metadata/*" { capabilities = ["create","update","delete","list","read"] }
      path "${enginePath}/undelete/*" { capabilities = ["update"] }
    `;
    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, policy);
    await logout.visit();
    await authPage.login(userToken);
    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    // create multiple versions
    await click('[data-test-secret-edit]');
    await editPage.editSecret('foo2', 'bar2');
    await click('[data-test-secret-edit]');
    await editPage.editSecret('foo3', 'bar3');
    // delete oldest version
    await click('[data-test-popup-menu-trigger="version"]');
    await click('[data-test-version-dropdown-link="1"]');
    await click('[data-test-delete-open-modal]');
    assert
      .dom('[data-test-type-select="delete-version"]')
      .hasText('Delete latest version', 'modal reads that it will delete latest version');
    await click('input#delete-version');
    await click('[data-test-modal-delete]');
    await visit(`/vault/secrets/${enginePath}/show/${secretPath}?version=3`);
    assert
      .dom('[data-test-empty-state-title]')
      .hasText(
        'Version 3 of this secret has been deleted',
        'empty state renders latest version has been deleted'
      );
    await visit(`/vault/secrets/${enginePath}/show/${secretPath}?version=1`);
    assert.dom('[data-test-delete-open-modal]').hasText('Delete', 'version 1 has not been deleted');
    await deleteEngine(enginePath, assert);
  });

  test('version 2: policy has "update" to /delete endpoint but not "delete" capability for secret path', async function (assert) {
    assert.expect(5);
    const enginePath = 'kv-v2-can-delete-version';
    const secretPath = 'kv-v2-delete-path-not-capability';
    const policy = `
      path "${enginePath}/data/${secretPath}" { capabilities = ["create","read","update","list"] }
      path "${enginePath}/metadata/*" { capabilities = ["create","update","delete","list","read"] }
      path "${enginePath}/undelete/*" { capabilities = ["update"] }
      path "${enginePath}/delete/*" { capabilities = ["update"] }
    `;
    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, policy);
    await logout.visit();
    await authPage.login(userToken);
    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    // create multiple versions
    await click('[data-test-secret-edit]');
    await editPage.editSecret('foo2', 'bar2');
    await click('[data-test-secret-edit]');
    await editPage.editSecret('foo3', 'bar3');
    // delete oldest version
    await click('[data-test-popup-menu-trigger="version"]');
    await click('[data-test-version-dropdown-link="1"]');
    await click('[data-test-delete-open-modal]');
    assert
      .dom('[data-test-type-select="delete-version"]')
      .hasText('Delete this version', 'delete option refers to "this" version');
    assert
      .dom('[data-test-delete-modal="delete-version"]')
      .hasTextContaining('Version 1', 'modal reads that it will delete version 1');
    await click('input#delete-version');
    await click('[data-test-modal-delete]');
    await visit(`/vault/secrets/${enginePath}/show/${secretPath}?version=3`);
    assert.dom('[data-test-delete-open-modal]').hasText('Delete', 'latest version (3) has not been deleted');
    await visit(`/vault/secrets/${enginePath}/show/${secretPath}?version=1`);
    assert
      .dom('[data-test-empty-state-title]')
      .hasText(
        'Version 1 of this secret has been deleted',
        'empty state renders oldest version (1) has been deleted'
      );
    await deleteEngine(enginePath, assert);
  });

  test('version 2 with path forward slash will show delete button', async function (assert) {
    assert.expect(2);
    const enginePath = 'kv-v2-forward-slash';
    const secretPath = 'forward/slash';
    const V2_POLICY = `
      path "${enginePath}/delete/${secretPath}" {
        capabilities = ["update"]
      }
      path "${enginePath}/metadata/*" {
        capabilities = ["list","read","create","update"]
      }
      path "${enginePath}/data/${secretPath}" {
        capabilities = ["create", "read"]
      }
    `;
    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, V2_POLICY);
    await logout.visit();
    await authPage.login(userToken);
    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    assert.dom('[data-test-secret-v2-delete="true"]').exists('drop down delete shows');
    await deleteEngine(enginePath, assert);
  });

  test('version 2 with engine with forward slash will show delete button', async function (assert) {
    assert.expect(2);
    const enginePath = 'forward/slash';
    const secretPath = 'secret-name';
    const V2_POLICY = `
      path "${enginePath}/delete/${secretPath}" {
        capabilities = ["update"]
      }
      path "${enginePath}/metadata/*" {
        capabilities = ["list","read","create","update"]
      }
      path "${enginePath}/data/*" {
        capabilities = ["create", "read"]
      }
    `;
    const userToken = await mountEngineGeneratePolicyToken(enginePath, secretPath, V2_POLICY);
    await logout.visit();
    await authPage.login(userToken);
    await writeSecret(enginePath, secretPath, 'foo', 'bar');
    assert.dom('[data-test-secret-v2-delete="true"]').exists('drop down delete shows');
    await deleteEngine(enginePath, assert);
  });

  const setupNoRead = async function (backend, canReadMeta = false) {
    const V2_WRITE_ONLY_POLICY = `
      path "${backend}/+/+" {
        capabilities = ["create", "update", "list"]
      }
      path "${backend}/+" {
        capabilities = ["list"]
      }
    `;

    const V2_WRITE_WITH_META_READ_POLICY = `
      path "${backend}/+/+" {
        capabilities = ["create", "update", "list"]
      }
      path "${backend}/metadata/+" {
        capabilities = ["read"]
      }
      path "${backend}/+" {
        capabilities = ["list"]
      }
    `;
    const V1_WRITE_ONLY_POLICY = `
     path "${backend}/+" {
        capabilities = ["create", "update", "list"]
      }
    `;

    const version = backend === 'kv-v2' ? 2 : 1;
    let policy;
    if (backend === 'kv-v2' && canReadMeta) {
      policy = V2_WRITE_WITH_META_READ_POLICY;
    } else if (backend === 'kv-v2') {
      policy = V2_WRITE_ONLY_POLICY;
    } else if (backend === 'kv-v1') {
      policy = V1_WRITE_ONLY_POLICY;
    }

    return await mountEngineGeneratePolicyToken(backend, 'nonexistent-secret', policy, version);
  };
  test('write without read: version 2', async function (assert) {
    assert.expect(5);
    const backend = 'kv-v2';
    const userToken = await setupNoRead(backend);
    await writeSecret(backend, 'secret', 'foo', 'bar');
    await logout.visit();
    await authPage.login(userToken);

    await showPage.visit({ backend, id: 'secret' });
    assert.ok(showPage.noReadIsPresent, 'shows no read empty state');
    assert.ok(showPage.editIsPresent, 'shows the edit button');

    await editPage.visitEdit({ backend, id: 'secret' });
    assert.notOk(editPage.hasMetadataFields, 'hides the metadata form');

    await editPage.editSecret('bar', 'baz');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    await deleteEngine(backend, assert);
  });

  test('write without read: version 2 with metadata read', async function (assert) {
    assert.expect(5);
    const backend = 'kv-v2';
    const userToken = await setupNoRead(backend, true);
    await writeSecret(backend, 'secret', 'foo', 'bar');
    await logout.visit();
    await authPage.login(userToken);

    await showPage.visit({ backend, id: 'secret' });
    assert.ok(showPage.noReadIsPresent, 'shows no read empty state');
    assert.ok(showPage.editIsPresent, 'shows the edit button');

    await editPage.visitEdit({ backend, id: 'secret' });
    assert
      .dom('[data-test-warning-no-read-permissions]')
      .hasText(
        'You do not have read permissions. If a secret exists here creating a new secret will overwrite it.'
      );

    await editPage.editSecret('bar', 'baz');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    await deleteEngine(backend, assert);
  });

  test('write without read: version 1', async function (assert) {
    assert.expect(4);
    const backend = 'kv-v1';
    const userToken = await setupNoRead(backend);
    await writeSecret(backend, 'secret', 'foo', 'bar');
    await logout.visit();
    await authPage.login(userToken);

    await showPage.visit({ backend, id: 'secret' });
    assert.ok(showPage.noReadIsPresent, 'shows no read empty state');
    assert.ok(showPage.editIsPresent, 'shows the edit button');

    await editPage.visitEdit({ backend, id: 'secret' });
    await editPage.editSecret('bar', 'baz');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    await deleteEngine(backend, assert);
  });
});
