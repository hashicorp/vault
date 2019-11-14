import { currentURL, currentRouteName } from '@ember/test-helpers';
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

const mount = async (shouldConfig = true) => {
  let path = `kmip-${Date.now()}`;
  let commands = shouldConfig
    ? [`write sys/mounts/${path} type=kmip`, `write ${path}/config -force`]
    : [`write sys/mounts/${path} type=kmip`];
  await uiConsole.runCommands(commands);
  return path;
};

const createScope = async () => {
  let path = await mount();
  let scope = `scope-${Date.now()}`;
  await uiConsole.runCommands([`write ${path}/scope/${scope} -force`]);
  return { path, scope };
};

const createRole = async () => {
  let { path, scope } = await createScope();
  let role = `role-${Date.now()}`;
  await uiConsole.runCommands([`write ${path}/scope/${scope}/role/${role} operation_all=true`]);
  return { path, scope, role };
};

const generateCreds = async () => {
  let { path, scope, role } = await createRole();
  await uiConsole.runCommands([
    `write ${path}/scope/${scope}/role/${role}/credential/generate format=pem
    -field=serial_number`,
  ]);
  let serial = uiConsole.lastLogOutput;
  return { path, scope, role, serial };
};

module('Acceptance | Enterprise | KMIP secrets', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it enables KMIP secrets engine', async function(assert) {
    let path = `kmip-${Date.now()}`;
    await mountSecrets.enable('kmip', path);

    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes`,
      'mounts and redirects to the kmip scopes page'
    );
    assert.ok(scopesPage.isEmpty, 'renders empty state');
  });

  test('it can configure a KMIP secrets engine', async function(assert) {
    let path = await mount(false);
    await scopesPage.visit({ backend: path });
    await scopesPage.configurationLink();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/configuration`,
      'configuration navigates to the config page'
    );
    assert.ok(scopesPage.isEmpty, 'config page renders empty state');

    await scopesPage.configureLink();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/configure`,
      'configuration navigates to the configure page'
    );
    await scopesPage.submit();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/configuration`,
      'redirects to configuration page after saving config'
    );
    assert.notOk(scopesPage.isEmpty, 'configuration page no longer renders empty state');
  });

  test('it can create a scope', async function(assert) {
    let path = await mount(this);
    await scopesPage.visit({ backend: path });
    await scopesPage.createLink();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/create`,
      'navigates to the kmip scope create page'
    );

    // create scope
    await scopesPage.scopeName('foo');
    await scopesPage.submit();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes`,
      'navigates to the kmip scopes page after create'
    );
    assert.equal(scopesPage.listItemLinks.length, 1, 'renders a single scope');
  });

  test('it can delete a scope from the list', async function(assert) {
    let { path } = await createScope(this);
    await scopesPage.visit({ backend: path });
    // delete the scope
    await scopesPage.listItemLinks.objectAt(0).menuToggle();
    await scopesPage.delete();
    await scopesPage.confirmDelete();
    assert.equal(scopesPage.listItemLinks.length, 0, 'no scopes');
    assert.ok(scopesPage.isEmpty, 'renders the empty state');
  });

  test('it can create a role', async function(assert) {
    let { path, scope } = await createScope(this);
    let role = `role-${Date.now()}`;
    await rolesPage.visit({ backend: path, scope });
    assert.ok(rolesPage.isEmpty, 'renders the empty role page');
    await rolesPage.create();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/create`,
      'links to the role create form'
    );

    await rolesPage.roleName(role);
    await rolesPage.submit();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles`,
      'redirects to roles list'
    );

    assert.equal(rolesPage.listItemLinks.length, 1, 'renders a single role');
  });

  test('it can delete a role from the list', async function(assert) {
    let { path, scope } = await createRole();
    await rolesPage.visit({ backend: path, scope });
    // delete the role
    await rolesPage.listItemLinks.objectAt(0).menuToggle();
    await rolesPage.delete();
    await rolesPage.confirmDelete();
    assert.equal(rolesPage.listItemLinks.length, 0, 'renders no roles');
    assert.ok(rolesPage.isEmpty, 'renders empty');
  });

  test('it can delete a role from the detail page', async function(assert) {
    let { path, scope, role } = await createRole(this);
    await rolesPage.visitDetail({ backend: path, scope, role });
    await rolesPage.detailEditLink();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}/edit`,
      'navigates to role edit'
    );
    await rolesPage.cancelLink();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}`,
      'cancel navigates to role show'
    );
    await rolesPage.delete().confirmDelete();

    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles`,
      'redirects to the roles list'
    );
    assert.ok(rolesPage.isEmpty, 'renders an empty roles page');
  });

  test('it can create a credential', async function(assert) {
    let { path, scope, role } = await createRole();
    await credentialsPage.visit({ backend: path, scope, role });
    assert.ok(credentialsPage.isEmpty, 'renders empty creds page');
    await credentialsPage.generateCredentialsLink();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}/credentials/generate`,
      'navigates to generate credentials'
    );
    await credentialsPage.submit();
    assert.equal(
      currentRouteName(),
      'vault.cluster.secrets.backend.kmip.credentials.show',
      'generate redirects to the show page'
    );
    await credentialsPage.backToRoleLink();

    assert.equal(credentialsPage.listItemLinks.length, 1, 'renders a single credential');
  });

  test('it can revoke a credential from the list', async function(assert) {
    let { path, scope, role } = await generateCreds();
    await credentialsPage.visit({ backend: path, scope, role });
    // revoke the credentials
    await credentialsPage.listItemLinks.objectAt(0).menuToggle();
    await credentialsPage.delete().confirmDelete();
    assert.equal(credentialsPage.listItemLinks.length, 0, 'renders no credentials');
    assert.ok(credentialsPage.isEmpty, 'renders empty');
  });

  test('it can revoke from the credentials show page', async function(assert) {
    let { path, scope, role, serial } = await generateCreds();
    await credentialsPage.visitDetail({ backend: path, scope, role, serial });
    await credentialsPage.delete().confirmDelete();

    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}/credentials`,
      'redirects to the credentials list'
    );
    assert.ok(credentialsPage.isEmpty, 'renders an empty credentials page');
  });
});
