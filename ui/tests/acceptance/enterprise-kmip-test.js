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

const mount = async () => {
  let path = `kmip-${Date.now()}`;
  await uiConsole.runCommands([`write sys/mounts/${path} type=kmip`, `write ${path}/config -force`]);
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

  test('it mounts KMIP secrets engine and scopes pages render appropriately', async function(assert) {
    let path = `kmip-${Date.now()}`;
    await mountSecrets.enable('kmip', path);

    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes`,
      'mounts and redirects to the kmip scopes page'
    );
    assert.ok(scopesPage.isEmpty, 'renders empty state');
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

  test('it can create and delete a scope', async function(assert) {
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

    // delete the scope
    await scopesPage.listItemLinks.objectAt(0).menuToggle();
    await scopesPage.delete();
    await scopesPage.confirmDelete();
    assert.equal(scopesPage.listItemLinks.length, 0, 'renders a single scope');
  });

  test('it can create and delete a role', async function(assert) {
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
    await rolesPage
      .detailDelete()
      .delete()
      .confirmDelete();
    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles`,
      'redirects to the roles list'
    );
    assert.ok(rolesPage.isEmpty, 'renders an empty roles page');
  });

  test('it can create and delete a credential', async function(assert) {
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

    // revoke the credentials
    await credentialsPage.listItemLinks.objectAt(0).menuToggle();
    await credentialsPage.delete();
    await credentialsPage.confirmDelete();
    assert.equal(credentialsPage.listItemLinks.length, 0, 'renders no credentials');
    assert.ok(credentialsPage.isEmpty, 'renders empty');
  });

  test('it can revoke from the credentials show page', async function(assert) {
    let { path, scope, role, serial } = await generateCreds();
    await credentialsPage.visitDetail({ backend: path, scope, role, serial });
    await credentialsPage
      .detailRevoke()
      .delete()
      .confirmDelete();

    assert.equal(
      currentURL(),
      `/vault/secrets/${path}/kmip/scopes/${scope}/roles/${role}/credentials`,
      'redirects to the credentials list'
    );
    assert.ok(credentialsPage.isEmpty, 'renders an empty credentials page');
  });
});
