/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, currentRouteName, settled, fillIn, waitUntil, find, click } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';

import authPage from 'vault/tests/pages/auth';
import scopesPage from 'vault/tests/pages/secrets/backend/kmip/scopes';
import rolesPage from 'vault/tests/pages/secrets/backend/kmip/roles';
import credentialsPage from 'vault/tests/pages/secrets/backend/kmip/credentials';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { allEngines } from 'vault/helpers/mountable-secret-engines';
import { runCmd } from 'vault/tests/helpers/commands';
import { v4 as uuidv4 } from 'uuid';

// port has a lower limit of 1024
const getRandomPort = () => Math.floor(Math.random() * 5000 + 1024);

const mount = async (backend) => {
  const res = await runCmd(`write sys/mounts/${backend} type=kmip`);
  await settled();
  if (res.includes('Error')) {
    throw new Error(`Error mounting secrets engine: ${res}`);
  }
  return backend;
};

const mountWithConfig = async (backend) => {
  const addr = `127.0.0.1:${getRandomPort()}`;
  await mount(backend);
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
    await authPage.login();
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
    await mountSecrets.selectType(engine.type);
    await mountSecrets.path(this.backend).submit();
    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.${engine.engineRoute}`,
      `Transitions to ${engine.displayName} route on mount success`
    );
    assert.ok(scopesPage.isEmpty, 'renders empty state');
  });

  test('it can configure a KMIP secrets engine', async function (assert) {
    const backend = await mount(this.backend);
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
    const backend = await mountWithConfig(this.backend);
    await settled();
    const scope = `scope-for-can-create-role`;
    await settled();
    const res = await runCmd([`write ${backend}/scope/${scope} -force`]);
    await settled();
    if (res.includes('Error')) {
      throw new Error(`Error creating scope: ${res}`);
    }
    const role = `role-new-role`;
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

    await rolesPage.roleName(role);
    await settled();
    await rolesPage.submit();
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/kmip/scopes/${scope}/roles`,
      'redirects to roles list'
    );

    assert.strictEqual(rolesPage.listItemLinks.length, 1, 'renders a single role');
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
    await rolesPage.cancelLink();
    await settled();
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
});
