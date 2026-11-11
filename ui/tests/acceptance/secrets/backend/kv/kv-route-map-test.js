/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentRouteName, currentURL, fillIn, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import {
  addSecretMetadataCmd,
  writeSecret,
  writeVersionedSecret,
} from 'vault/tests/helpers/kv/kv-run-commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';

const ROUTE_PREFIX = 'vault.cluster.secrets.backend.kv';

/**
 * Characterizes the public contract of the kv routes: every URL maps to the same full route name,
 * dynamic segments and query params are preserved, and route names resolve back to the same URLs.
 * None of these assertions depend on where the route files live.
 */
module('Acceptance | kv | route map', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.backend = `kv-route-map-${uuidv4()}`;
    this.secretPath = 'app/first-secret';
    this.encodedSecretPath = encodeURIComponent(this.secretPath);
    this.base = `/vault/secrets-engines/${this.backend}/kv`;
    await login();
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await writeVersionedSecret(this.backend, this.secretPath, 'hello', 'there', 3);
    await runCmd(addSecretMetadataCmd(this.backend, this.secretPath, { max_versions: 5 }));
  });

  hooks.afterEach(async function () {
    await visit('/vault/secrets-engines');
    await login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  test('the engine root redirects to the list route', async function (assert) {
    await visit(this.base);
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.list`, 'route name is list');
    assert.strictEqual(currentURL(), `${this.base}/list`, 'url is the list url');
  });

  test('list routes keep their names, urls and query params', async function (assert) {
    await visit(`${this.base}/list`);
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.list`);
    assert.dom(PAGE.list.item('app/')).exists('top level directory is listed');

    await visit(`${this.base}/list?pageFilter=app`);
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.list`);
    assert.strictEqual(currentURL(), `${this.base}/list?pageFilter=app`, 'pageFilter query param is kept');

    await visit(`${this.base}/list/app`);
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.list-directory`, 'directory list route name');
    assert.strictEqual(currentURL(), `${this.base}/list/app`, 'directory url is not rewritten');
    assert.dom(PAGE.list.item('first-secret')).exists('secret in the directory is listed');

    await visit(`${this.base}/list/app/`);
    assert.strictEqual(
      currentRouteName(),
      `${ROUTE_PREFIX}.list-directory`,
      'directory with trailing slash resolves to the same route'
    );
    assert.strictEqual(currentURL(), `${this.base}/list/app/`, 'trailing slash url is not rewritten');
  });

  test('create route keeps its name, url and initialKey query param', async function (assert) {
    await visit(`${this.base}/create`);
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.create`);
    assert.strictEqual(currentURL(), `${this.base}/create`);

    await visit(`${this.base}/create?initialKey=app/new-secret`);
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.create`);
    assert.strictEqual(
      decodeURIComponent(currentURL()),
      `${this.base}/create?initialKey=app/new-secret`,
      'initialKey query param is kept'
    );
  });

  const secretRoutes = [
    ['secret overview', '', 'secret.index', GENERAL.overviewCard.container('Current version')],
    ['secret details', '/details', 'secret.details.index', PAGE.detail.versionTimestamp],
    [
      'secret details with version query param',
      '/details?version=2',
      'secret.details.index',
      PAGE.detail.versionTimestamp,
    ],
    [
      'new version form from the latest version',
      '/details/edit?version=3',
      'secret.details.edit',
      FORM.saveBtn,
    ],
    ['new version form from a version', '/details/edit?version=2', 'secret.details.edit', FORM.saveBtn],
    ['paths', '/paths', 'secret.paths', `${PAGE.secretTab('Paths')}.active`],
    ['metadata details', '/metadata', 'secret.metadata.index', PAGE.metadata.editBtn],
    ['metadata edit', '/metadata/edit', 'secret.metadata.edit', FORM.saveBtn],
    ['version history', '/metadata/versions', 'secret.metadata.versions', PAGE.versions.linkedBlock(3)],
    ['version diff', '/metadata/diff', 'secret.metadata.diff', PAGE.diff.visualDiff],
  ];

  for (const [description, suffix, name, selector] of secretRoutes) {
    test(`secret route keeps its name and url: ${description}`, async function (assert) {
      const url = `${this.base}/${this.encodedSecretPath}${suffix}`;
      await visit(url);
      assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.${name}`, 'route name');
      assert.strictEqual(currentURL(), url, 'url is not rewritten');
      assert.dom(selector).exists('page content renders');
    });
  }

  test('the secret name segment keeps nested paths encoded as one segment', async function (assert) {
    const nested = encodeURIComponent('app/nested/deeper');
    await writeSecret(this.backend, 'app/nested/deeper', 'foo', 'bar');
    await visit(`${this.base}/${nested}/details`);
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.secret.details.index`);
    assert.strictEqual(currentURL(), `${this.base}/${nested}/details`, 'encoded slashes are preserved');
    assert.dom(PAGE.infoRowValue('foo')).exists('secret data for the nested path is rendered');
  });

  test('configuration route is reached from the Manage menu and keeps its name and url', async function (assert) {
    await visit(`${this.base}/list`);
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.configuration`);
    assert.strictEqual(currentURL(), `${this.base}/configuration`);
    assert.dom(PAGE.title).hasText(`${this.backend} configuration`);
  });

  test('configuration can be edited and exited', async function (assert) {
    await visit(`${this.base}/configuration`);
    await click(PAGE.configuration.editBtn);
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.configure`, 'edit goes to configure route');
    assert.strictEqual(currentURL(), `${this.base}/configure`);

    await fillIn('#max_versions', '7');
    await click(FORM.saveBtn);
    assert.strictEqual(
      currentRouteName(),
      `${ROUTE_PREFIX}.configuration`,
      'saving returns to the configuration route'
    );
    assert.dom(PAGE.infoRowValue('Maximum number of versions')).hasText('7', 'saved value is shown');

    await click(GENERAL.button('Exit configuration'));
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.list`, 'exit returns to the list route');
  });
});

module('Acceptance | kv | route generation', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.backend = `kv-route-gen-${uuidv4()}`;
    await login();
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
  });

  hooks.afterEach(async function () {
    await visit('/vault/secrets-engines');
    await login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  test('route names generate the same urls', async function (assert) {
    await visit(`/vault/secrets-engines/${this.backend}/kv/list`);
    const router = this.owner.lookup('service:router');
    const base = `/ui/vault/secrets-engines/${this.backend}/kv`;
    const secret = 'app%2Ffirst-secret';
    const secretRoute = `${ROUTE_PREFIX}.secret`;

    const cases = [
      [`${ROUTE_PREFIX}.list`, [], `${base}/list`],
      [`${ROUTE_PREFIX}.list-directory`, ['app/'], `${base}/list/app/`],
      [`${ROUTE_PREFIX}.create`, [], `${base}/create`],
      [`${ROUTE_PREFIX}.configuration`, [], `${base}/configuration`],
      [`${ROUTE_PREFIX}.configure`, [], `${base}/configure`],
      [`${secretRoute}.index`, [secret], `${base}/${secret}`],
      [`${secretRoute}.details`, [secret], `${base}/${secret}/details`],
      [`${secretRoute}.details.edit`, [secret], `${base}/${secret}/details/edit`],
      [`${secretRoute}.paths`, [secret], `${base}/${secret}/paths`],
      [`${secretRoute}.metadata`, [secret], `${base}/${secret}/metadata`],
      [`${secretRoute}.metadata.edit`, [secret], `${base}/${secret}/metadata/edit`],
      [`${secretRoute}.metadata.versions`, [secret], `${base}/${secret}/metadata/versions`],
      [`${secretRoute}.metadata.diff`, [secret], `${base}/${secret}/metadata/diff`],
    ];

    for (const [name, params, url] of cases) {
      // urlFor encodes params itself, so pass the decoded value and compare the real encoded path
      const decoded = params.map((p) => decodeURIComponent(p));
      assert.strictEqual(router.urlFor(name, this.backend, ...decoded), url, `${name} generates ${url}`);
    }
  });
});
