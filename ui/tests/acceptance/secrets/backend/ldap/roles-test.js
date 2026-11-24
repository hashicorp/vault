/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, visit, waitFor } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupApplicationTest } from 'ember-qunit';
import { module, test } from 'qunit';
import sinon from 'sinon';
import { v4 as uuidv4 } from 'uuid';
import ldapHandlers from 'vault/mirage/handlers/ldap';
import ldapMirageScenario from 'vault/mirage/scenarios/ldap';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { assertURL, isURL, visitURL } from 'vault/tests/helpers/ldap/ldap-helpers';
import { LDAP_SELECTORS } from 'vault/tests/helpers/ldap/ldap-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Acceptance | ldap | roles', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    ldapHandlers(this.server);
    ldapMirageScenario(this.server);
    this.backend = `ldap-test-${uuidv4()}`;
    await login();
    // mount & configure
    await runCmd([
      mountEngineCmd('ldap', this.backend),
      `write ${this.backend}/config binddn=foo bindpass=bar url=http://localhost:8208`,
    ]);
    return visitURL('roles', this.backend);
  });

  hooks.afterEach(async function () {
    await runCmd(deleteEngineCmd(this.backend));
  });

  test('it should transition to create role route on toolbar link click', async function (assert) {
    await click('[data-test-toolbar-action="role"]');
    assert.true(
      isURL('roles/create', this.backend),
      'Transitions to role create route on toolbar link click'
    );
  });

  test('it should transition to role details route on list item click', async function (assert) {
    let path;

    await click(LDAP_SELECTORS.roleItem('dynamic', 'dynamic-role'));
    path = 'roles/dynamic/dynamic-role/details';
    assertURL(assert, this.backend, path);
    await click(GENERAL.breadcrumbLink('Roles'));

    await click(LDAP_SELECTORS.roleItem('static', 'static-role'));
    path = 'roles/static/static-role/details';
    assertURL(assert, this.backend, path);
    await click(GENERAL.breadcrumbLink('Roles'));

    // edge case, roles of different type with same name
    await click(LDAP_SELECTORS.roleItem('dynamic', 'my-role'));
    path = 'roles/dynamic/my-role/details';
    assertURL(assert, this.backend, path);
    await click(GENERAL.breadcrumbLink('Roles'));

    await click(LDAP_SELECTORS.roleItem('static', 'my-role'));
    path = 'roles/static/my-role/details';
    assertURL(assert, this.backend, path);
  });

  test('it should transition to routes from list item action menu', async function (assert) {
    assert.expect(3);

    for (const action of ['edit', 'get-creds', 'details']) {
      await click(LDAP_SELECTORS.roleMenu('dynamic', 'dynamic-role'));
      await click(LDAP_SELECTORS.action(action));
      const uri = action === 'get-creds' ? 'credentials' : action;
      assert.true(
        isURL(`roles/dynamic/dynamic-role/${uri}`, this.backend),
        `Transitions to ${uri} route on list item action menu click`
      );
      await click(GENERAL.breadcrumbLink('Roles'));
    }
  });

  test('it should transition to routes from role details page header actions', async function (assert) {
    await click(LDAP_SELECTORS.roleItem('dynamic', 'dynamic-role'));
    await click(GENERAL.button('Get credentials'));
    assert.true(
      isURL('roles/dynamic/dynamic-role/credentials', this.backend),
      'Transitions to credentials route from toolbar link'
    );

    await click(GENERAL.breadcrumbLink('dynamic-role'));
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Edit role'));
    assert.true(
      isURL('roles/dynamic/dynamic-role/edit', this.backend),
      'Transitions to edit route from toolbar link'
    );
  });

  // Vault handles rotate-role as an update operation, so a create-only token is denied.
  test('it should offer rotate credentials only with update capability on rotate-role', async function (assert) {
    const stubRotateCapabilities = (rotateCapabilities) => {
      this.server.post('/sys/capabilities-self', (_, { requestBody }) => {
        const { paths } = JSON.parse(requestBody);
        const data = paths.reduce((obj, path) => {
          obj[path] = path.includes('/rotate-role/') ? rotateCapabilities : ['root'];
          return obj;
        }, {});
        return { ...data, data };
      });
    };

    // Leave the roles page first so revisiting it re-runs the model hook with the new stub.
    stubRotateCapabilities(['update']);
    await visitURL('overview', this.backend);
    await visitURL('roles', this.backend);
    await click(LDAP_SELECTORS.roleMenu('static', 'static-role'));
    assert.dom(LDAP_SELECTORS.action('rotate-creds')).exists('list menu offers rotate with update');
    await click(LDAP_SELECTORS.roleItem('static', 'static-role'));
    await click(GENERAL.dropdownToggle('Manage'));
    assert.dom(GENERAL.menuItem('Rotate credentials')).exists('details menu offers rotate with update');

    stubRotateCapabilities(['create']);
    await visitURL('overview', this.backend);
    await visitURL('roles', this.backend);
    await click(LDAP_SELECTORS.roleMenu('static', 'static-role'));
    assert.dom(LDAP_SELECTORS.action('rotate-creds')).doesNotExist('list menu hides rotate with create only');
    await click(LDAP_SELECTORS.roleItem('static', 'static-role'));
    await click(GENERAL.dropdownToggle('Manage'));
    assert
      .dom(GENERAL.menuItem('Rotate credentials'))
      .doesNotExist('details menu hides rotate with create only');
  });

  test('it should clear roles page filter value on route exit', async function (assert) {
    await fillIn('[data-test-filter-input]', 'foo');
    assert
      .dom('[data-test-filter-input]')
      .hasValue('foo', 'Roles page filter value set after model refresh and rerender');
    await waitFor(GENERAL.emptyStateTitle);
    await click('[data-test-tab="libraries"]');
    await click('[data-test-tab="roles"]');
    assert.dom('[data-test-filter-input]').hasNoValue('Roles page filter value cleared on route exit');
  });

  // A token that cannot read the mount config cannot tell whether the mount is self-managed.
  test('it should flag the mount type as unknown when the config cannot be read', async function (assert) {
    this.server.get('/:backend/config', () => overrideResponse(403));
    // The engine reads the config once on entry, so leave it before visiting again.
    await visit('/vault/secrets');
    await visitURL('roles/create', this.backend);

    assert
      .dom(`${GENERAL.inlineAlertByAttr('mount-type-unknown')} .hds-alert__title`)
      .hasText('Unable to determine LDAP configuration mode', 'the create page flags the unknown mount type');
    assert.dom(GENERAL.radioCardByAttr('dynamic')).exists('both role types are still offered');
  });

  // Only a 403 means the token lacks permission; other failures must not claim it does.
  test('it should not flag the mount type when the config read fails for another reason', async function (assert) {
    this.server.get('/:backend/config', () => overrideResponse(500));
    await visit('/vault/secrets');
    await visitURL('roles/create', this.backend);

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create Role', 'the create page renders');
    assert
      .dom(GENERAL.inlineAlertByAttr('mount-type-unknown'))
      .doesNotExist('a server error is not reported as missing permission');
  });

  module('self-managed mounts', function (hooks) {
    hooks.beforeEach(async function () {
      // Flip the mount config written in the outer hook to self-managed. Setting the record
      // directly rather than writing via the CLI keeps `self_managed` a real boolean: Vault coerces
      // it per its field schema and returns JSON, but the mirage handler stores the request body
      // verbatim, so a CLI write would leave the string "true" here.
      this.server.db.ldapConfigs.update({ backend: this.backend }, { self_managed: true });
      // The engine's application route resolves the config once on entry, so leave the engine to
      // force the updated config to be re-read when each test navigates back in.
      await visit('/vault/secrets');
    });

    test('it hides the role type choice and titles the page for static roles', async function (assert) {
      await visitURL('roles/create', this.backend);

      assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create static role', 'page is titled for static roles');
      assert
        .dom(GENERAL.radioCardByAttr('static'))
        .doesNotExist('role type choice is not offered on a self-managed mount');
      assert.dom(GENERAL.radioCardByAttr('dynamic')).doesNotExist('dynamic roles are not offered');
    });

    test('it requires a distinguished name and password before submitting', async function (assert) {
      await visitURL('roles/create', this.backend);

      await fillIn(GENERAL.inputByAttr('name'), 'self-managed-role');
      await fillIn(GENERAL.inputByAttr('username'), 'some-user');
      await click(GENERAL.submitButton);

      assert
        .dom(GENERAL.validationErrorByAttr('dn'))
        .hasText('Enter the distinguished name (DN) of the LDAP account', 'dn is required');
      assert
        .dom(GENERAL.validationErrorByAttr('password'))
        .hasText('Enter the current password for this LDAP account', 'password is required');
      assert.true(isURL('roles/create', this.backend), 'submission is blocked, still on the create page');
    });

    test('it creates a static role', async function (assert) {
      await visitURL('roles/create', this.backend);

      await fillIn(GENERAL.inputByAttr('name'), 'self-managed-role');
      await fillIn(GENERAL.inputByAttr('dn'), 'cn=self-managed,dc=example,dc=com');
      await fillIn(GENERAL.inputByAttr('username'), 'some-user');
      await fillIn(GENERAL.inputByAttr('password'), 'super-secret');
      await click(GENERAL.submitButton);

      assertURL(assert, this.backend, 'roles/static/self-managed-role/details');
    });
  });

  module('root-managed mounts', function () {
    test('it offers both role types and does not require a password', async function (assert) {
      await visitURL('roles/create', this.backend);

      assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create Role', 'page title is unchanged');
      assert.dom(GENERAL.radioCardByAttr('static')).exists('static role type is offered');
      assert.dom(GENERAL.radioCardByAttr('dynamic')).exists('dynamic role type is offered');

      await fillIn(GENERAL.inputByAttr('name'), 'root-managed-role');
      await fillIn(GENERAL.inputByAttr('username'), 'some-user');
      await click(GENERAL.submitButton);

      assertURL(assert, this.backend, 'roles/static/root-managed-role/details');
    });
  });

  // The details page holds the revealed credential on the component, not on the shared route
  // model, so viewing it must not carry the plaintext into the edit form.
  test('it does not carry a revealed password into the edit form', async function (assert) {
    await visitURL('roles/static/static-role/details', this.backend);
    await click(GENERAL.button('toggle-masked'));
    assert.dom(GENERAL.maskedInput).exists('the password is revealed on the details page');

    await visitURL('roles/static/static-role/edit', this.backend);
    await click(GENERAL.enableField('password'));

    assert
      .dom(GENERAL.inputByAttr('password'))
      .hasNoValue('the edit form does not pre-fill with the revealed password');
  });

  module('subdirectory', function () {
    test('it navigates to hierarchical roles', async function (assert) {
      let path;
      // hierarchical paths
      await click(LDAP_SELECTORS.roleItem('dynamic', 'admin/'));
      path = 'roles/dynamic/subdirectory/admin/';
      assertURL(assert, this.backend, path);

      await click(LDAP_SELECTORS.roleItem('dynamic', 'child-dynamic-role'));
      path = 'roles/dynamic/admin%2Fchild-dynamic-role/details';
      assertURL(assert, this.backend, path);

      // navigate out via breadcrumbs to test
      await click(GENERAL.breadcrumbLink('admin'));
      path = 'roles/dynamic/subdirectory/admin/';
      assertURL(assert, this.backend, path);

      await click(GENERAL.breadcrumbLink('Roles'));

      await click(LDAP_SELECTORS.roleItem('static', 'admin/'));
      path = 'roles/static/subdirectory/admin/';
      assertURL(assert, this.backend, path);
      await click(LDAP_SELECTORS.roleItem('static', 'child-static-role'));
      path = 'roles/static/admin%2Fchild-static-role/details';
      assertURL(assert, this.backend, path);

      // navigate out via breadcrumbs to test
      await click(GENERAL.breadcrumbLink('admin'));
      path = 'roles/static/subdirectory/admin/';
      assertURL(assert, this.backend, path);
    });

    test('it should transition to subdirectory from hierarchical role popup menu', async function (assert) {
      assert.expect(4);

      await click(LDAP_SELECTORS.roleMenu('dynamic', 'admin/'));
      for (const action of ['edit', 'get-creds', 'details']) {
        assert.dom(LDAP_SELECTORS.action(action)).doesNotExist(`${action} does not render in popup menu`);
      }
      await click(LDAP_SELECTORS.action('subdirectory'));
      assertURL(assert, this.backend, 'roles/dynamic/subdirectory/admin/');
    });

    test('it should clear roles page filter value on route exit', async function (assert) {
      await visitURL('roles/static/subdirectory/admin/', this.backend);
      await fillIn('[data-test-filter-input]', 'foo');
      assert
        .dom('[data-test-filter-input]')
        .hasValue('foo', 'Roles page filter value set after model refresh and rerender');
      await waitFor(GENERAL.emptyStateTitle);
      await click('[data-test-tab="libraries"]');
      await click('[data-test-tab="roles"]');
      assert.dom('[data-test-filter-input]').hasNoValue('Roles page filter value cleared on route exit');
    });

    test('it calls API with correct parameter order for hierarchical paths', async function (assert) {
      // Get the API service and stub the SDK methods
      const owner = this.owner;
      const apiService = owner.lookup('service:api');
      const ldapListStaticRolePathStub = sinon.stub(apiService.secrets, 'ldapListStaticRolePath');
      const ldapListRolePathStub = sinon.stub(apiService.secrets, 'ldapListRolePath');

      // Configure stubs to return empty lists
      ldapListStaticRolePathStub.resolves({ keys: [] });
      ldapListRolePathStub.resolves({ keys: [] });

      // Navigate to static role subdirectory
      await visitURL('roles/static/subdirectory/admin/', this.backend);

      // Verify static role API was called with correct parameter order
      assert.true(ldapListStaticRolePathStub.calledOnce, 'ldapListStaticRolePath was called');
      const [staticPath, staticMount] = ldapListStaticRolePathStub.firstCall.args;
      assert.strictEqual(staticPath, 'admin/', 'First parameter is the hierarchical path');
      assert.strictEqual(staticMount, this.backend, 'Second parameter is the mount path');

      // Navigate to dynamic role subdirectory
      await visitURL('roles/dynamic/subdirectory/admin/', this.backend);

      // Verify dynamic role API was called with correct parameter order
      const [dynamicPath] = ldapListRolePathStub.firstCall.args;
      assert.strictEqual(dynamicPath, 'admin/', 'First parameter is the hierarchical path for dynamic roles');

      // Restore stubs
      ldapListStaticRolePathStub.restore();
      ldapListRolePathStub.restore();
    });
  });
});
