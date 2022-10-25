import { currentURL, currentRouteName, settled, fillIn, waitUntil, find } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';

import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import authPage from 'vault/tests/pages/auth';
import scopesPage from 'vault/tests/pages/secrets/backend/kmip/scopes';
import rolesPage from 'vault/tests/pages/secrets/backend/kmip/roles';
import credentialsPage from 'vault/tests/pages/secrets/backend/kmip/credentials';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';

const uiConsole = create(consoleClass);

const getRandomPort = () => {
  let a = Math.floor(100000 + Math.random() * 900000);
  a = String(a);
  return a.substring(0, 4);
};

const mount = async (shouldConfig = true) => {
  const now = Date.now();
  let path = `kmip-${now}`;
  let addr = `127.0.0.1:${getRandomPort()}`; // use random port
  await settled();
  let commands = shouldConfig
    ? [`write sys/mounts/${path} type=kmip`, `write ${path}/config listen_addrs=${addr}`]
    : [`write sys/mounts/${path} type=kmip`];
  await uiConsole.runCommands(commands);
  await settled();
  let res = uiConsole.lastLogOutput;
  if (res.includes('Error')) {
    throw new Error(`Error mounting secrets engine: ${res}`);
  }
  return path;
};

const createScope = async () => {
  let path = await mount();
  await settled();
  let scope = `scope-${Date.now()}`;
  await settled();
  await uiConsole.runCommands([`write ${path}/scope/${scope} -force`]);
  await settled();
  let res = uiConsole.lastLogOutput;
  if (res.includes('Error')) {
    throw new Error(`Error creating scope: ${res}`);
  }
  return { path, scope };
};

const createRole = async () => {
  let { path, scope } = await createScope();
  await settled();
  let role = `role-${Date.now()}`;
  await uiConsole.runCommands([`write ${path}/scope/${scope}/role/${role} operation_all=true`]);
  await settled();
  let res = uiConsole.lastLogOutput;
  if (res.includes('Error')) {
    throw new Error(`Error creating role: ${res}`);
  }
  return { path, scope, role };
};

const generateCreds = async () => {
  let { path, scope, role } = await createRole();
  await settled();
  await uiConsole.runCommands([
    `write ${path}/scope/${scope}/role/${role}/credential/generate format=pem -field=serial_number`,
  ]);
  let serial = uiConsole.lastLogOutput;
  if (serial.includes('Error')) {
    throw new Error(`Credential generation failed with error: ${serial}`);
  }
  return { path, scope, role, serial };
};
module('Acceptance | Enterprise | KMIP secrets', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    return;
  });

  test('it enables KMIP secrets engine', async function (assert) {
    let path = `kmip-${Date.now()}`;
    await mountSecrets.enable('kmip', path);
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes`,
      'mounts and redirects to the kmip scopes page'
    );
    assert.ok(scopesPage.isEmpty, 'renders empty state');
  });

  test('it can configure a KMIP secrets engine', async function (assert) {
    let path = await mount(false);
    await scopesPage.visit({ backend: path });
    await settled();
    await scopesPage.configurationLink();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/configuration`,
      'configuration navigates to the config page'
    );
    assert.ok(scopesPage.isEmpty, 'config page renders empty state');

    await scopesPage.configureLink();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/configure`,
      'configuration navigates to the configure page'
    );
    let addr = `127.0.0.1:${getRandomPort()}`;
    await fillIn('[data-test-string-list-input="0"]', addr);

    await scopesPage.submit();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/configuration`,
      'redirects to configuration page after saving config'
    );
    assert.notOk(scopesPage.isEmpty, 'configuration page no longer renders empty state');
  });

  test('it can revoke from the credentials show page', async function (assert) {
    let { path, scope, role, serial } = await generateCreds();
    await settled();
    await credentialsPage.visitDetail({ backend: path, scope, role, serial });
    await settled();
    await waitUntil(() => find('[data-test-confirm-action-trigger]'));
    assert.dom('[data-test-confirm-action-trigger]').exists('delete button exists');
    await credentialsPage.delete().confirmDelete();
    await settled();

    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}/credentials`,
      'redirects to the credentials list'
    );
    assert.ok(credentialsPage.isEmpty, 'renders an empty credentials page');
  });

  test('it can create a scope', async function (assert) {
    let path = await mount(this);
    await scopesPage.visit({ backend: path });
    await settled();
    await scopesPage.createLink();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/create`,
      'navigates to the kmip scope create page'
    );

    // create scope
    await scopesPage.scopeName('foo');
    await settled();
    await scopesPage.submit();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes`,
      'navigates to the kmip scopes page after create'
    );
    assert.equal(scopesPage.listItemLinks.length, 1, 'renders a single scope');
  });

  test('it can delete a scope from the list', async function (assert) {
    let { path } = await createScope();
    await scopesPage.visit({ backend: path });
    await settled();
    // delete the scope
    await scopesPage.listItemLinks.objectAt(0).menuToggle();
    await settled();
    await scopesPage.delete();
    await settled();
    await scopesPage.confirmDelete();
    await settled();
    assert.equal(scopesPage.listItemLinks.length, 0, 'no scopes');
    assert.ok(scopesPage.isEmpty, 'renders the empty state');
  });

  test('it can create a role', async function (assert) {
    // moving create scope here to help with flaky test
    let path = await mount();
    await settled();
    let scope = `scope-for-can-create-role`;
    await settled();
    await uiConsole.runCommands([`write ${path}/scope/${scope} -force`]);
    await settled();
    let res = uiConsole.lastLogOutput;
    if (res.includes('Error')) {
      throw new Error(`Error creating scope: ${res}`);
    }
    let role = `role-new-role`;
    await rolesPage.visit({ backend: path, scope });
    await settled();
    assert.ok(rolesPage.isEmpty, 'renders the empty role page');
    await rolesPage.create();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/create`,
      'links to the role create form'
    );

    await rolesPage.roleName(role);
    await settled();
    await rolesPage.submit();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles`,
      'redirects to roles list'
    );

    assert.equal(rolesPage.listItemLinks.length, 1, 'renders a single role');
  });

  test('it can delete a role from the list', async function (assert) {
    let { path, scope } = await createRole();
    await rolesPage.visit({ backend: path, scope });
    await settled();
    // delete the role
    await rolesPage.listItemLinks.objectAt(0).menuToggle();
    await settled();
    await rolesPage.delete();
    await settled();
    await rolesPage.confirmDelete();
    await settled();
    assert.equal(rolesPage.listItemLinks.length, 0, 'renders no roles');
    assert.ok(rolesPage.isEmpty, 'renders empty');
  });

  test('it can delete a role from the detail page', async function (assert) {
    let { path, scope, role } = await createRole();
    await settled();
    await rolesPage.visitDetail({ backend: path, scope, role });
    await settled();
    await waitUntil(() => find('[data-test-kmip-link-edit-role]'));
    await rolesPage.detailEditLink();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}/edit`,
      'navigates to role edit'
    );
    await rolesPage.cancelLink();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}`,
      'cancel navigates to role show'
    );
    await rolesPage.delete().confirmDelete();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles`,
      'redirects to the roles list'
    );
    assert.ok(rolesPage.isEmpty, 'renders an empty roles page');
  });

  test('it can create a credential', async function (assert) {
    // TODO come back and figure out why issue here with test
    let { path, scope, role } = await createRole();
    await credentialsPage.visit({ backend: path, scope, role });
    await settled();
    assert.ok(credentialsPage.isEmpty, 'renders empty creds page');
    await credentialsPage.generateCredentialsLink();
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}/credentials/generate`,
      'navigates to generate credentials'
    );
    await credentialsPage.submit();
    await settled();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.kmip.credentials.show',
      'generate redirects to the show page'
    );
    await credentialsPage.backToRoleLink();
    await settled();
    assert.equal(credentialsPage.listItemLinks.length, 1, 'renders a single credential');
  });

  test('it can revoke a credential from the list', async function (assert) {
    let { path, scope, role } = await generateCreds();
    await credentialsPage.visit({ backend: path, scope, role });
    // revoke the credentials
    await settled();
    await credentialsPage.listItemLinks.objectAt(0).menuToggle();
    await settled();
    await credentialsPage.delete().confirmDelete();
    await settled();
    assert.equal(credentialsPage.listItemLinks.length, 0, 'renders no credentials');
    assert.ok(credentialsPage.isEmpty, 'renders empty');
  });
});
