/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { v4 as uuidv4 } from 'uuid';
import { click, currentRouteName, currentURL, fillIn, visit } from '@ember/test-helpers';
import ldapMirageScenario from 'vault/mirage/scenarios/ldap';
import ldapHandlers from 'vault/mirage/handlers/ldap';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { LDAP_SELECTORS } from 'vault/tests/helpers/ldap/ldap-selectors';

const ROUTE_PREFIX = 'vault.cluster.secrets.backend.ldap';

/**
 * Pins the public URL and full route name of every LDAP route. These must not change when the
 * engine is inlined into the main app.
 */
module('Acceptance | ldap | routing', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    ldapHandlers(this.server);
    ldapMirageScenario(this.server);
    this.backend = `ldap-test-${uuidv4()}`;
    this.base = `/vault/secrets-engines/${this.backend}/ldap`;
    this.assertRoute = (assert, path, routeName, message) => {
      assert.strictEqual(currentURL(), `${this.base}/${path}`, message || `url is ${path}`);
      assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.${routeName}`, `route name is ${routeName}`);
    };
    this.configure = () =>
      runCmd([
        mountEngineCmd('ldap', this.backend),
        `write ${this.backend}/config binddn=foo bindpass=bar url=http://localhost:8208`,
      ]);
    await login();
  });

  hooks.afterEach(async function () {
    await runCmd(deleteEngineCmd(this.backend));
  });

  module('configured mount', function (hooks) {
    hooks.beforeEach(function () {
      return this.configure();
    });

    test('it resolves top level routes', async function (assert) {
      const routes = {
        overview: 'overview',
        configuration: 'configuration',
        configure: 'configure',
        roles: 'roles.index',
        'roles/create': 'roles.create',
        libraries: 'libraries.index',
        'libraries/create': 'libraries.create',
      };
      for (const [path, routeName] of Object.entries(routes)) {
        await visit(`${this.base}/${path}`);
        this.assertRoute(assert, path, routeName);
      }
    });

    test('it resolves role routes', async function (assert) {
      const routes = {
        'roles/static/static-role/details': 'roles.role.details',
        'roles/static/static-role/edit': 'roles.role.edit',
        'roles/static/static-role/credentials': 'roles.role.credentials',
        'roles/dynamic/dynamic-role/details': 'roles.role.details',
        'roles/dynamic/dynamic-role/edit': 'roles.role.edit',
        'roles/dynamic/dynamic-role/credentials': 'roles.role.credentials',
      };
      for (const [path, routeName] of Object.entries(routes)) {
        await visit(`${this.base}/${path}`);
        this.assertRoute(assert, path, routeName);
      }
    });

    test('it redirects a role url without a tab to role details', async function (assert) {
      await visit(`${this.base}/roles/static/static-role`);
      this.assertRoute(assert, 'roles/static/static-role/details', 'roles.role.details');
      await visit(`${this.base}/roles/dynamic/dynamic-role`);
      this.assertRoute(assert, 'roles/dynamic/dynamic-role/details', 'roles.role.details');
    });

    test('it resolves hierarchical role routes', async function (assert) {
      await visit(`${this.base}/roles/static/subdirectory/admin/`);
      this.assertRoute(assert, 'roles/static/subdirectory/admin/', 'roles.subdirectory');
      assert.dom(LDAP_SELECTORS.roleItem('static', 'child-static-role')).exists();

      await visit(`${this.base}/roles/dynamic/subdirectory/admin/`);
      this.assertRoute(assert, 'roles/dynamic/subdirectory/admin/', 'roles.subdirectory');
      assert.dom(LDAP_SELECTORS.roleItem('dynamic', 'child-dynamic-role')).exists();

      await click(LDAP_SELECTORS.roleItem('dynamic', 'child-dynamic-role'));
      assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.roles.role.details`);
      assert.strictEqual(
        currentURL(),
        `${this.base}/roles/dynamic/admin%2Fchild-dynamic-role/details`,
        'child role name is encoded in the url'
      );
    });

    test('it resolves library routes', async function (assert) {
      const routes = {
        'libraries/test-library/details/accounts': 'libraries.library.details.accounts',
        'libraries/test-library/details/configuration': 'libraries.library.details.configuration',
        'libraries/test-library/edit': 'libraries.library.edit',
      };
      for (const [path, routeName] of Object.entries(routes)) {
        await visit(`${this.base}/${path}`);
        this.assertRoute(assert, path, routeName);
      }
    });

    test('it redirects incomplete library urls to the accounts tab', async function (assert) {
      const accounts = 'libraries/test-library/details/accounts';
      await visit(`${this.base}/libraries/test-library`);
      this.assertRoute(assert, accounts, 'libraries.library.details.accounts');
      await visit(`${this.base}/libraries/test-library/details`);
      this.assertRoute(assert, accounts, 'libraries.library.details.accounts');
    });

    test('it redirects a library check-out deep link to the accounts tab', async function (assert) {
      await visit(`${this.base}/libraries/test-library/check-out`);
      this.assertRoute(
        assert,
        'libraries/test-library/details/accounts',
        'libraries.library.details.accounts',
        'check-out is only reachable from the accounts tab'
      );
    });

    test('it resolves hierarchical library routes', async function (assert) {
      await visit(`${this.base}/libraries/subdirectory/admin/`);
      this.assertRoute(assert, 'libraries/subdirectory/admin/', 'libraries.subdirectory');
      assert.dom(LDAP_SELECTORS.libraryItem('test-library')).exists();

      await visit(`${this.base}/libraries/admin%2Ftest-library/details/accounts`);
      this.assertRoute(
        assert,
        'libraries/admin%2Ftest-library/details/accounts',
        'libraries.library.details.accounts'
      );
    });

    test('it keeps the role list page filter in the url', async function (assert) {
      await visit(`${this.base}/roles?pageFilter=static`);
      assert.strictEqual(currentURL(), `${this.base}/roles?pageFilter=static`);
      assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.roles.index`);
      assert.dom(LDAP_SELECTORS.roleItem('static', 'static-role')).exists();
      assert.dom(LDAP_SELECTORS.roleItem('dynamic', 'dynamic-role')).doesNotExist('filter is applied');
    });

    test('it filters the library list without changing the url', async function (assert) {
      await visit(`${this.base}/libraries`);
      assert.dom(LDAP_SELECTORS.libraryItem('test-library')).exists();

      await fillIn('input[aria-label="Filter libraries"]', 'no-such-library');
      assert.strictEqual(currentURL(), `${this.base}/libraries`, 'filter is not a query param');
      assert
        .dom(LDAP_SELECTORS.libraryItem('test-library'))
        .doesNotExist('non-matching filter hides library');

      await fillIn('input[aria-label="Filter libraries"]', 'test');
      assert.dom(LDAP_SELECTORS.libraryItem('test-library')).exists('matching filter shows library');
    });

    test('it redirects the mount root url to overview', async function (assert) {
      await visit(`/vault/secrets-engines/${this.backend}`);
      this.assertRoute(assert, 'overview', 'overview');
    });

    test('it navigates between routes with the tabs and breadcrumbs', async function (assert) {
      await visit(`${this.base}/overview`);
      await click('[data-test-tab="roles"]');
      this.assertRoute(assert, 'roles', 'roles.index');
      await click('[data-test-tab="libraries"]');
      this.assertRoute(assert, 'libraries', 'libraries.index');
      await click('[data-test-tab="overview"]');
      this.assertRoute(assert, 'overview', 'overview');

      await visit(`${this.base}/configuration`);
      await click(GENERAL.breadcrumbLink(this.backend));
      this.assertRoute(assert, 'overview', 'overview');
    });
  });

  module('unconfigured mount', function (hooks) {
    hooks.beforeEach(function () {
      return runCmd(mountEngineCmd('ldap', this.backend));
    });

    test('it redirects configuration to configure', async function (assert) {
      await visit(`${this.base}/configuration`);
      this.assertRoute(assert, 'configure', 'configure');
    });

    test('it stays on overview and offers configuration', async function (assert) {
      await visit(`${this.base}/overview`);
      this.assertRoute(assert, 'overview', 'overview');
      await click(`${GENERAL.emptyStateActions} a`);
      this.assertRoute(assert, 'configure', 'configure');
    });
  });
});
