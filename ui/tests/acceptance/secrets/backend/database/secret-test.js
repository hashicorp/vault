import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL, settled, click, visit, fillIn, findAll } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import connectionPage from 'vault/tests/pages/secrets/backend/database/connection';
import rolePage from 'vault/tests/pages/secrets/backend/database/role';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import searchSelect from 'vault/tests/pages/components/search-select';

const searchSelectComponent = create(searchSelect);

const consoleComponent = create(consoleClass);

const MODEL = {
  engineType: 'database',
  id: 'database-name',
};

const mount = async () => {
  let path = `database-${Date.now()}`;
  await mountSecrets.enable('database', path);
  await settled();
  return path;
};

const newConnection = async backend => {
  const name = `connection-${Date.now()}`;
  await connectionPage.visitCreate({ backend });
  await connectionPage.dbPlugin('mongodb-database-plugin');
  await connectionPage.name(name);
  await connectionPage.url(`mongodb://127.0.0.1:4321/${name}`);
  await connectionPage.toggleVerify();
  await connectionPage.save();
  await connectionPage.enable();
  return name;
};

module('Acceptance | secrets/database/*', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    this.server = apiStub({ usePassthrough: true });
    return authPage.login();
  });
  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('can enable the database secrets engine', async function(assert) {
    let backend = `database-${Date.now()}`;
    await mountSecrets.enable('database', backend);
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'Mounts and redirects to connection list page'
    );
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert
      .dom('.is-active[data-test-secret-list-tab="Connections"]')
      .exists('Has Connections tab which is active');
    await click('[data-test-tab="overview"]');
    assert.equal(currentURL(), `/vault/secrets/${backend}/overview`, 'Tab links to overview page');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state also exists on overview page');
    assert.dom('[data-test-secret-list-tab="Roles"]').exists('Has Roles tab');
  });

  test('Connection create and edit form happy path works as expected', async function(assert) {
    const backend = await mount();
    const connectionDetails = {
      plugin: 'mongodb-database-plugin',
      id: 'horses-db',
      fields: [
        { label: 'Connection Name', name: 'name', value: 'horses-db' },
        { label: 'Connection url', name: 'connection_url', value: 'mongodb://127.0.0.1:235/horses' },
        { label: 'Username', name: 'username', value: 'user', hideOnShow: true },
        { label: 'Password', name: 'password', password: 'so-secure', hideOnShow: true },
        { label: 'Write concern', name: 'write_concern' },
      ],
    };
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'Mounts and redirects to connection list page'
    );
    await connectionPage.createLink();
    assert.equal(currentURL(), `/vault/secrets/${backend}/create`, 'Create link goes to create page');
    assert
      .dom('[data-test-empty-state-title')
      .hasText('No plugin selected', 'No plugin is selected by default and empty state shows');
    await connectionPage.dbPlugin(connectionDetails.plugin);
    assert.dom('[data-test-empty-state]').doesNotExist('Empty state goes away after plugin selected');
    connectionDetails.fields.forEach(async ({ name, value }) => {
      assert
        .dom(`[data-test-input="${name}"]`)
        .exists(`Field ${name} exists for ${connectionDetails.plugin}`);
      if (value) {
        await fillIn(`[data-test-input="${name}"]`, value);
      }
    });
    // uncheck verify for the save step to work
    await connectionPage.toggleVerify();
    await connectionPage.save();
    assert
      .dom('[data-test-modal-title]')
      .hasText('Rotate your root credentials?', 'Modal appears asking to ');
    await connectionPage.enable();
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/show/${connectionDetails.id}`,
      'Saves connection and takes you to show page'
    );
    connectionDetails.fields.forEach(({ label, name, value, hideOnShow }) => {
      if (hideOnShow) {
        assert
          .dom(`[data-test-row-value="${label}"]`)
          .doesNotExist(`Does not show ${name} value on show page for ${connectionDetails.plugin}`);
      } else if (!value) {
        assert.dom(`[data-test-row-value="${label}"]`).hasText('Default');
      } else {
        assert.dom(`[data-test-row-value="${label}"]`).hasText(value);
      }
    });
    // go back and edit write_concern
    await connectionPage.edit();
    assert.dom(`[data-test-input="name"]`).hasAttribute('readonly');
    assert.dom(`[data-test-input="plugin_name"]`).hasAttribute('readonly');
    // assert password is hidden
    findAll('.CodeMirror')[0].CodeMirror.setValue(JSON.stringify({ wtimeout: 5000 }));
    // uncheck verify for the save step to work
    await connectionPage.toggleVerify();
    await connectionPage.save();
    assert
      .dom(`[data-test-row-value="Write concern"]`)
      .hasText('{ "wtimeout": 5000 }', 'Write concern is now showing on the table');
  });

  test('buttons show up for managing connection', async function(assert) {
    const backend = await mount();
    const connection = await newConnection(backend);
    await connectionPage.visitShow({ backend, id: connection });
    assert
      .dom('[data-test-database-connection-delete]')
      .hasText('Delete connection', 'Delete connection button exists with correct text');
    assert
      .dom('[data-test-database-connection-reset]')
      .hasText('Reset connection', 'Reset button exists with correct text');
    assert.dom('[data-test-secret-create]').hasText('Add role', 'Add role button exists with correct text');
    assert.dom('[data-test-edit-link]').hasText('Edit configuration', 'Edit button exists with correct text');
    const CONNECTION_VIEW_ONLY = `
path "${backend}/*" {
  capabilities = ["deny"]
}
path "${backend}/config" {
  capabilities = ["list"]
}
path "${backend}/config/*" {
  capabilities = ["read"]
}
    `;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=database`,
      `write sys/policies/acl/test-policy policy=${btoa(CONNECTION_VIEW_ONLY)}`,
      'write -field=client_token auth/token/create policies=test-policy ttl=1h',
    ]);
    let token = consoleComponent.lastTextOutput;
    await logout.visit();
    await authPage.login(token);
    await connectionPage.visitShow({ backend, id: connection });
    assert.equal(currentURL(), `/vault/secrets/${backend}/show/${connection}`, 'Allows reading connection');
    assert
      .dom('[data-test-database-connection-delete]')
      .doesNotExist('Delete button does not show due to permissions');
    assert
      .dom('[data-test-database-connection-reset]')
      .doesNotExist('Reset button does not show due to permissions');
    assert.dom('[data-test-secret-create]').doesNotExist('Add role button does not show due to permissions');
    assert.dom('[data-test-edit-link]').doesNotExist('Edit button does not show due to permissions');
    await visit(`/vault/secrets/${backend}/overview`);
    assert.dom('[data-test-selectable-card="Connections"]').exists('Connections card exists on overview');
    assert
      .dom('[data-test-selectable-card="Roles"]')
      .doesNotExist('Roles card does not exist on overview w/ policy');
    assert.dom('.title-number').hasText('1', 'Lists the correct number of connections');
  });

  test('Role create form', async function(assert) {
    const backend = await mount();
    // Connection needed for role fields
    await newConnection(backend);
    await rolePage.visitCreate({ backend });
    await rolePage.name('bar');
    assert
      .dom('[data-test-component="empty-state"]')
      .exists({ count: 2 }, 'Two empty states exist before selections made');
    await clickTrigger('#database');
    assert.equal(searchSelectComponent.options.length, 1, 'list shows existing connections so far');
    await selectChoose('#database', '.ember-power-select-option', 0);
    assert
      .dom('[data-test-component="empty-state"]')
      .exists({ count: 2 }, 'Two empty states exist before selections made');
    await rolePage.roleType('static');
    assert.dom('[data-test-component="empty-state"]').doesNotExist('Empty states go away');
    assert.dom('[data-test-input="username"]').exists('Username field appears for static role');
    assert
      .dom('[data-test-toggle-input="Rotation period"]')
      .exists('Rotation period field appears for static role');
    await rolePage.roleType('dynamic');
    assert
      .dom('[data-test-toggle-input="Generated credentials’s Time-to-Live (TTL)"]')
      .exists('TTL field exists for dynamic');
    assert
      .dom('[data-test-toggle-input="Generated credentials’s maximum Time-to-Live (Max TTL)"]')
      .exists('Max TTL field exists for dynamic');
  });

  test('root and limited access', async function(assert) {
    this.set('model', MODEL);
    let backend = 'database';
    const NO_ROLES_POLICY = `
      path "database/roles/*" {
        capabilities = ["delete"]
      }
      path "database/static-roles/*" {
        capabilities = ["delete"]
      }
      path "database/config/*" {
        capabilities = ["list", "create", "read", "update"]
      }
      path "database/creds/*" {
        capabilities = ["list", "create", "read", "update"]
      }
    `;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=database`,
      `write sys/policies/acl/test-policy policy=${btoa(NO_ROLES_POLICY)}`,
      'write -field=client_token auth/token/create policies=test-policy ttl=1h',
    ]);
    let token = consoleComponent.lastTextOutput;

    // test root user flow
    await settled();

    // await click('[data-test-secret-backend-row="database"]');
    // skipping the click because occasionally is shows up on the second page and cannot be found
    await visit(`/vault/secrets/database/overview`);
    await settled();
    assert.dom('[data-test-component="empty-state"]').exists('renders empty state');
    assert.dom('[data-test-secret-list-tab="Connections"]').exists('renders connections tab');
    assert.dom('[data-test-secret-list-tab="Roles"]').exists('renders connections tab');

    await click('[data-test-secret-create="connections"]');
    assert.equal(currentURL(), '/vault/secrets/database/create');

    // Login with restricted policy
    await logout.visit();
    await authPage.login(token);
    await settled();
    // skipping the click because occasionally is shows up on the second page and cannot be found
    await visit(`/vault/secrets/database/overview`);
    assert.dom('[data-test-tab="overview"]').exists('renders overview tab');
    assert.dom('[data-test-secret-list-tab="Connections"]').exists('renders connections tab');
    assert
      .dom('[data-test-secret-list-tab="Roles]')
      .doesNotExist(`does not show the roles tab because it does not have permissions`);
    assert
      .dom('[data-test-selectable-card="Connections"]')
      .exists({ count: 1 }, 'renders only the connection card');

    await click('[data-test-action-text="Configure new"]');
    assert.equal(currentURL(), '/vault/secrets/database/create?itemType=connection');
  });
});
