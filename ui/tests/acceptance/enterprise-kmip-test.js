/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  currentURL,
  currentRouteName,
  settled,
  fillIn,
  visit,
  waitUntil,
  find,
  click,
} from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import scopesPage from 'vault/tests/pages/secrets/backend/kmip/scopes';
import rolesPage from 'vault/tests/pages/secrets/backend/kmip/roles';
import credentialsPage from 'vault/tests/pages/secrets/backend/kmip/credentials';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { allEngines } from 'vault/helpers/mountable-secret-engines';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { v4 as uuidv4 } from 'uuid';

// port has a lower limit of 1024
const getRandomPort = () => Math.floor(Math.random() * 5000 + 1024);

const mountWithConfig = async (backend) => {
  const addr = `127.0.0.1:${getRandomPort()}`;
  await runCmd(mountEngineCmd('kmip', backend), false);
  const res = await runCmd(`write ${backend}/config listen_addrs=${addr}`);
  if (res.includes('Error')) {
    throw new Error(`Error configuring KMIP: ${res}`);
  }
  return backend;
};

const createScope = async (backend) => {
  await mountWithConfig(backend);
  await settled();
  const scope = `scope-${uuidv4()}`;
  await settled();
  const res = await runCmd([`write ${backend}/scope/${scope} -force`]);
  await settled();
  if (res.includes('Error')) {
    throw new Error(`Error creating scope: ${res}`);
  }
  return { backend, scope };
};

const createRole = async (backend) => {
  const { scope } = await createScope(backend);
  await settled();
  const role = `role-${uuidv4()}`;
  const res = await runCmd([`write ${backend}/scope/${scope}/role/${role} operation_all=true`]);
  await settled();
  if (res.includes('Error')) {
    throw new Error(`Error creating role: ${res}`);
  }
  return { backend, scope, role };
};

const generateCreds = async (backend) => {
  const { scope, role } = await createRole(backend);
  await settled();
  const serial = await runCmd([
    `write ${backend}/scope/${scope}/role/${role}/credential/generate format=pem -field=serial_number`,
  ]);
  await settled();
  if (serial.includes('Error')) {
    throw new Error(`Credential generation failed with error: ${serial}`);
  }
  return { backend, scope, role, serial };
};

module('Acceptance | Enterprise | KMIP secrets', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.backend = `kmip-${uuidv4()}`;
    await login();
    return;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.backend}`], false);
  });

  test('it should enable KMIP & transitions to addon engine route after mount success', async function (assert) {
    // test supported backends that ARE ember engines (enterprise only engines are tested individually)
    const engine = allEngines().find((e) => e.type === 'kmip');

    await mountSecrets.visit();
    await mountBackend(engine.type, `${engine.type}-${uuidv4()}`);

    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.${engine.engineRoute}`,
      `Transitions to ${engine.displayName} route on mount success`
    );
    assert.ok(scopesPage.isEmpty, 'renders empty state');
  });

  test('it can configure a KMIP secrets engine', async function (assert) {
    await runCmd(mountEngineCmd('kmip', this.backend));
    const backend = this.backend;
    await scopesPage.visit({ backend });
    await settled();
    await scopesPage.configurationLink();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/configuration`,
      'configuration navigates to the config page'
    );
    assert.ok(scopesPage.isEmpty, 'config page renders empty state');

    await scopesPage.configureLink();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/configure`,
      'configuration navigates to the configure page'
    );
    const addr = `127.0.0.1:${getRandomPort()}`;
    await fillIn('[data-test-string-list-input="0"]', addr);
    await scopesPage.submit();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/configuration`,
      'redirects to configuration page after saving config'
    );
    assert.notOk(scopesPage.isEmpty, 'configuration page no longer renders empty state');
    assert.dom('[data-test-value-div="Listen addrs"]').hasText(addr, 'renders the correct listen address');
  });

  test('it can revoke from the credentials show page', async function (assert) {
    const { backend, scope, role, serial } = await generateCreds(this.backend);
    await settled();
    await credentialsPage.visitDetail({ backend, scope, role, serial });
    await settled();
    await waitUntil(() => find('[data-test-confirm-action-trigger]'));
    assert.dom('[data-test-confirm-action-trigger]').exists('delete button exists');
    await credentialsPage.delete().confirmDelete();
    await settled();

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/${scope}/roles/${role}/credentials`,
      'redirects to the credentials list'
    );
    assert.ok(credentialsPage.isEmpty, 'renders an empty credentials page');
  });

  test('it can create a scope', async function (assert) {
    const backend = await mountWithConfig(this.backend);
    await scopesPage.visit({ backend });
    await settled();
    await scopesPage.createLink();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/create`,
      'navigates to the kmip scope create page'
    );

    // create scope
    await scopesPage.scopeName('foo');
    await settled();
    await scopesPage.submit();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes`,
      'navigates to the kmip scopes page after create'
    );
    assert.strictEqual(scopesPage.listItemLinks.length, 1, 'renders a single scope');
  });

  test('it navigates to kmip scopes view using breadcrumbs', async function (assert) {
    const backend = await mountWithConfig(this.backend);
    await scopesPage.visitCreate({ backend });
    await settled();
    await click(GENERAL.breadcrumbLink(backend));

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.kmip.scopes.index',
      'Breadcrumb transitions to scopes list'
    );
  });

  test('it can delete a scope from the list', async function (assert) {
    const { backend } = await createScope(this.backend);
    await scopesPage.visit({ backend });
    await settled();
    // delete the scope
    await scopesPage.listItemLinks.objectAt(0).menuToggle();
    await settled();
    await scopesPage.delete();
    await settled();
    await scopesPage.confirmDelete();
    await settled();
    assert.strictEqual(scopesPage.listItemLinks.length, 0, 'no scopes');
    assert.ok(scopesPage.isEmpty, 'renders the empty state');
  });

  test('it can create a role', async function (assert) {
    // moving create scope here to help with flaky test
    const scope = `scope-for-can-create-role`;
    const role = `role-new-role`;
    const backend = await mountWithConfig(this.backend);
    await settled();
    await runCmd([`write ${backend}/scope/${scope} -force`], true);
    await rolesPage.visit({ backend, scope });
    await settled();
    assert.ok(rolesPage.isEmpty, 'renders the empty role page');
    await rolesPage.create();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/${scope}/roles/create`,
      'links to the role create form'
    );
    // check that the role form looks right
    assert.dom(GENERAL.inputByAttr('operationNone')).isChecked('allows role to perform roles by default');
    assert.dom(GENERAL.inputByAttr('operationAll')).isChecked('operationAll is checked by default');
    assert.dom('[data-test-kmip-section]').exists({ count: 2 });
    assert.dom('[data-test-kmip-operations]').exists({ count: 4 });

    await rolesPage.roleName(role);
    await settled();
    await click(GENERAL.saveButton);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/${scope}/roles`,
      'redirects to roles list'
    );

    assert.strictEqual(rolesPage.listItemLinks.length, 1, 'renders a single role');
    await rolesPage.visitDetail({ backend, scope, role });
    // check that the role details looks right
    assert.dom('h2').exists({ count: 3 }, 'renders correct section headings');
    assert.dom('[data-test-inline-error-message]').hasText('This role allows all KMIP operations');
    ['Managed Cryptographic Objects', 'Object Attributes', 'Server', 'Other'].forEach((title) => {
      assert.dom(`[data-test-row-label="${title}"]`).exists(`Renders allowed operations row for: ${title}`);
    });
  });

  test('it navigates to kmip roles view using breadcrumbs', async function (assert) {
    const { backend, scope, role } = await createRole(this.backend);
    await settled();
    await rolesPage.visitDetail({ backend, scope, role });
    // navigate to scope from role
    await click(GENERAL.breadcrumbLink(scope));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.kmip.scope.roles',
      'Breadcrumb transitions to scope details'
    );
    await rolesPage.visitDetail({ backend, scope, role });
    // navigate to scopes from role
    await click(GENERAL.breadcrumbLink(backend));
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.kmip.scopes.index',
      'Breadcrumb transitions to scopes list'
    );
  });

  test('it can delete a role from the list', async function (assert) {
    const { backend, scope } = await createRole(this.backend);
    await rolesPage.visit({ backend, scope });
    await settled();
    // delete the role
    await rolesPage.listItemLinks.objectAt(0).menuToggle();
    await settled();
    await rolesPage.delete();
    await settled();
    await rolesPage.confirmDelete();
    await settled();
    assert.strictEqual(rolesPage.listItemLinks.length, 0, 'renders no roles');
    assert.ok(rolesPage.isEmpty, 'renders empty');
  });

  test('it can delete a role from the detail page', async function (assert) {
    const { backend, scope, role } = await createRole(this.backend);
    await settled();
    await rolesPage.visitDetail({ backend, scope, role });
    await settled();
    await waitUntil(() => find('[data-test-kmip-link-edit-role]'));
    await rolesPage.detailEditLink();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/${scope}/roles/${role}/edit`,
      'navigates to role edit'
    );
    await click(GENERAL.cancelButton);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/${scope}/roles/${role}`,
      'cancel navigates to role show'
    );
    await rolesPage.delete().confirmDelete();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/${scope}/roles`,
      'redirects to the roles list'
    );
    assert.ok(rolesPage.isEmpty, 'renders an empty roles page');
  });

  test('it can create a credential', async function (assert) {
    const { backend, scope, role } = await createRole(this.backend);
    await credentialsPage.visit({ backend, scope, role });
    await settled();
    assert.ok(credentialsPage.isEmpty, 'renders empty creds page');
    await credentialsPage.generateCredentialsLink();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/${scope}/roles/${role}/credentials/generate`,
      'navigates to generate credentials'
    );
    await credentialsPage.submit();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.kmip.credentials.show',
      'generate redirects to the show page'
    );
    await credentialsPage.backToRoleLink();
    await settled();
    assert.strictEqual(credentialsPage.listItemLinks.length, 1, 'renders a single credential');
  });

  test('it can revoke a credential from the list', async function (assert) {
    const { backend, scope, role } = await generateCreds(this.backend);
    await credentialsPage.visit({ backend, scope, role });
    // revoke the credentials
    await settled();
    await credentialsPage.listItemLinks.objectAt(0).menuToggle();
    await settled();
    await credentialsPage.delete().confirmDelete();
    await settled();
    assert.strictEqual(credentialsPage.listItemLinks.length, 0, 'renders no credentials');
    assert.ok(credentialsPage.isEmpty, 'renders empty');
  });

  // the kmip/role model relies on openApi so testing the form via an acceptance test
  module('kmip role edit form', function (hooks) {
    hooks.beforeEach(async function () {
      this.store = this.owner.lookup('service:store');
      this.scope = 'my-scope';
      this.name = 'my-role';
      await login();
      await runCmd(mountEngineCmd('kmip', this.backend), false);
      await runCmd([`write ${this.backend}/scope/${this.scope} -force`]);
      await rolesPage.visit({ backend: this.backend, scope: this.scope });
      this.setModel = async () => {
        await click(GENERAL.saveButton);
        await visit(`/vault/secrets/${this.backend}/kmip/scopes/${this.scope}/roles/${this.name}`);
        this.model = this.store.peekRecord('kmip/role', this.name);
      };
    });

    // "operationNone" is the attr name for the 'Allow this role to perform KMIP operations' toggle
    // operationNone = false => the toggle is ON and KMIP operations are allowed
    // operationNone = true => the toggle is OFF and KMIP operations are not allowed
    test('it submits when operationNone is toggled on', async function (assert) {
      assert.expect(3);

      await click('[data-test-role-create]');
      await fillIn(GENERAL.inputByAttr('role'), this.name);
      assert.dom(GENERAL.inputByAttr('operationAll')).isChecked('operationAll is checked by default');
      await this.setModel();
      assert.true(this.model.operationAll, 'operationAll is true');
      assert.strictEqual(this.model.operationNone, undefined, 'operationNone is unset');
    });

    test('it submits when operationNone is toggled off', async function (assert) {
      assert.expect(4);

      await click('[data-test-role-create]');
      await fillIn(GENERAL.inputByAttr('role'), this.name);
      await click(GENERAL.inputByAttr('operationNone'));
      assert
        .dom(GENERAL.inputByAttr('operationNone'))
        .isNotChecked('Allow this role to perform KMIP operations is toggled off');
      assert
        .dom(GENERAL.inputByAttr('operationAll'))
        .doesNotExist('clicking the toggle hides KMIP operation checkboxes');
      await this.setModel();
      assert.strictEqual(this.model.operationAll, undefined, 'operationAll is unset');
      assert.true(this.model.operationNone, 'operationNone is true');
    });

    test('it submits when operationAll is unchecked', async function (assert) {
      assert.expect(2);

      await click('[data-test-role-create]');
      await fillIn(GENERAL.inputByAttr('role'), this.name);
      await click(GENERAL.inputByAttr('operationAll'));
      await this.setModel();

      assert.strictEqual(this.model.operationAll, undefined, 'operationAll is unset');
      assert.true(this.model.operationNone, 'operationNone is true');
    });

    test('it submits individually selected operations', async function (assert) {
      assert.expect(6);

      await click('[data-test-role-create]');
      await fillIn(GENERAL.inputByAttr('role'), this.name);
      await click(GENERAL.inputByAttr('operationAll'));
      await click(GENERAL.inputByAttr('operationGet'));
      await click(GENERAL.inputByAttr('operationGetAttributes'));
      assert.dom(GENERAL.inputByAttr('operationAll')).isNotChecked();
      assert.dom(GENERAL.inputByAttr('operationCreate')).isNotChecked(); // unchecking operationAll deselects the other checkboxes
      await this.setModel();
      assert.strictEqual(this.model.operationAll, undefined, 'operationAll is unset');
      assert.strictEqual(this.model.operationNone, undefined, 'operationNone is unset');
      assert.true(this.model.operationGet, 'operationGet is true');
      assert.true(this.model.operationGetAttributes, 'operationGetAttributes is true');
    });
  });
});
