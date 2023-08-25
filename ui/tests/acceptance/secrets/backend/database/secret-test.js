/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL, settled, click, visit, fillIn, typeIn } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import connectionPage from 'vault/tests/pages/secrets/backend/database/connection';
import rolePage from 'vault/tests/pages/secrets/backend/database/role';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import searchSelect from 'vault/tests/pages/components/search-select';
import { deleteEngineCmd, mountEngineCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';

const searchSelectComponent = create(searchSelect);

const newConnection = async (backend, plugin = 'mongodb-database-plugin') => {
  const name = `connection-${Date.now()}`;
  await connectionPage.visitCreate({ backend });
  await connectionPage.dbPlugin(plugin);
  await connectionPage.name(name);
  await connectionPage.connectionUrl(`mongodb://127.0.0.1:4321/${name}`);
  await connectionPage.toggleVerify();
  await connectionPage.save();
  await connectionPage.enable();
  return name;
};

const navToConnection = async (backend, connection) => {
  await visit('/vault/secrets');
  await click(`[data-test-secrets-backend-link="${backend}"]`);
  await click('[data-test-secret-list-tab="Connections"]');
  await click(`[data-test-secret-link="${connection}"]`);
  return;
};

const connectionTests = [
  {
    name: 'elasticsearch-connection',
    plugin: 'elasticsearch-database-plugin',
    elasticUser: 'username',
    elasticPassword: 'password',
    url: 'http://127.0.0.1:9200',
    assertCount: 9,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert.dom('[data-test-input="ca_cert"]').exists(`CA certificate field exists for ${name}`);
      assert.dom('[data-test-input="ca_path"]').exists(`CA path field exists for ${name}`);
      assert.dom('[data-test-input="client_cert"]').exists(`Client certificate field exists for ${name}`);
      assert.dom('[data-test-input="client_key"]').exists(`Client key field exists for ${name}`);
      assert.dom('[data-test-input="tls_server_name"]').exists(`TLS server name field exists for ${name}`);
      assert.dom('[data-test-input="insecure"]').exists(`Insecure checkbox exists for ${name}`);
      assert
        .dom('[data-test-toggle-input="show-username_template"]')
        .exists(`Username template toggle exists for ${name}`);
    },
  },
  {
    name: 'mongodb-connection',
    plugin: 'mongodb-database-plugin',
    url: `mongodb://127.0.0.1:4321/test`,
    assertCount: 5,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert.dom('[data-test-input="write_concern"]').exists(`Write concern field exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mssql-connection',
    plugin: 'mssql-database-plugin',
    url: `mssql://127.0.0.1:4321/test`,
    assertCount: 6,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mysql-connection',
    plugin: 'mysql-database-plugin',
    url: `{{username}}:{{password}}@tcp(127.0.0.1:3306)/test`,
    assertCount: 7,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mysql-aurora-connection',
    plugin: 'mysql-aurora-database-plugin',
    url: `{{username}}:{{password}}@tcp(127.0.0.1:3306)/test`,
    assertCount: 7,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mysql-rds-connection',
    plugin: 'mysql-rds-database-plugin',
    url: `{{username}}:{{password}}@tcp(127.0.0.1:3306)/test`,
    assertCount: 7,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mysql-legacy-connection',
    plugin: 'mysql-legacy-database-plugin',
    url: `{{username}}:{{password}}@tcp(127.0.0.1:3306)/test`,
    assertCount: 7,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'postgresql-connection',
    plugin: 'postgresql-database-plugin',
    url: `postgresql://{{username}}:{{password}}@localhost:5432/postgres?sslmode=disable`,
    assertCount: 7,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
      assert
        .dom('[data-test-toggle-input="show-username_template"]')
        .exists(`Username template toggle exists for ${name}`);
    },
  },
];

module('Acceptance | secrets/database/*', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.backend = `database-testing`;
    await authPage.login();
    return runCmd(mountEngineCmd('database', this.backend), false);
  });
  hooks.afterEach(async function () {
    await authPage.login();
    return runCmd(deleteEngineCmd(this.backend), false);
  });

  test('can enable the database secrets engine', async function (assert) {
    const backend = `database-${Date.now()}`;
    await mountSecrets.enable('database', backend);
    await settled();
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'Mounts and redirects to connection list page'
    );
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert
      .dom('.active[data-test-secret-list-tab="Connections"]')
      .exists('Has Connections tab which is active');
    await click('[data-test-tab="overview"]');
    assert.strictEqual(currentURL(), `/vault/secrets/${backend}/overview`, 'Tab links to overview page');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state also exists on overview page');
    assert.dom('[data-test-secret-list-tab="Roles"]').exists('Has Roles tab');
    await visit('/vault/secrets');
    // Cleanup backend
    await runCmd(deleteEngineCmd(backend), false);
  });

  for (const testCase of connectionTests) {
    test(`database connection create and edit: ${testCase.plugin}`, async function (assert) {
      assert.expect(19 + testCase.assertCount);
      const backend = this.backend;
      await connectionPage.visitCreate({ backend });
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/create`, 'Correct creation URL');
      assert
        .dom('[data-test-empty-state-title]')
        .hasText('No plugin selected', 'No plugin is selected by default and empty state shows');
      await connectionPage.dbPlugin(testCase.plugin);
      assert.dom('[data-test-empty-state]').doesNotExist('Empty state goes away after plugin selected');
      await connectionPage.name(testCase.name);
      if (testCase.plugin === 'elasticsearch-database-plugin') {
        await connectionPage.url(testCase.url);
        await connectionPage.username(testCase.elasticUser);
        await connectionPage.password(testCase.elasticPassword);
      } else {
        await connectionPage.connectionUrl(testCase.url);
      }
      testCase.requiredFields(assert, testCase.plugin);
      assert.dom('[data-test-input="verify_connection"]').isChecked('verify is checked');
      await connectionPage.toggleVerify();
      assert.dom('[data-test-input="verify_connection"]').isNotChecked('verify is unchecked');
      assert
        .dom('[data-test-database-oracle-alert]')
        .doesNotExist('does not show oracle alert for non-oracle plugins');
      await connectionPage.save();
      await settled();
      assert
        .dom('.modal.is-active .title')
        .hasText('Rotate your root credentials?', 'Modal appears asking to rotate root credentials');
      assert.dom('[data-test-enable-connection]').exists('Enable button exists');
      await click('[data-test-enable-connection]');
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/show/${testCase.name}`),
        `Saves connection and takes you to show page for ${testCase.name}`
      );
      assert
        .dom(`[data-test-row-value="Password"]`)
        .doesNotExist(`Does not show Password value on show page for ${testCase.name}`);
      await connectionPage.edit();
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/edit/${testCase.name}`),
        `Edit connection button and takes you to edit page for ${testCase.name}`
      );
      assert.dom(`[data-test-input="name"]`).hasAttribute('readonly');
      assert.dom(`[data-test-input="plugin_name"]`).hasAttribute('readonly');
      assert.dom('[data-test-input="password"]').doesNotExist('Password is not displayed on edit form');
      assert.dom('[data-test-toggle-input="show-password"]').exists('Update password toggle exists');

      assert.dom('[data-test-input="verify_connection"]').isNotChecked('verify is still unchecked');
      await connectionPage.save();
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/show/${testCase.name}`);
      // click "Add Role"
      await connectionPage.addRole();
      await settled();
      assert.strictEqual(
        searchSelectComponent.selectedOptions[0].text,
        testCase.name,
        'Database connection is pre-selected on the form'
      );
      await click('[data-test-database-role-cancel]');
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/list`, 'Cancel button links to list view');
    });
  }
  test('database connection create and edit: vault-plugin-database-oracle', async function (assert) {
    assert.expect(11);
    // keep oracle as separate test because it behaves differently than the others
    const testCase = {
      name: 'oracle-connection',
      plugin: 'vault-plugin-database-oracle',
      url: `{{username}}/{{password}}@localhost:1521/OraDoc.localhost`,
      requiredFields: async (assert, name) => {
        assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
        assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
        assert
          .dom('[data-test-input="max_open_connections"]')
          .exists(`Max open connections exists for ${name}`);
        assert
          .dom('[data-test-input="max_idle_connections"]')
          .exists(`Max idle connections exists for ${name}`);
        assert
          .dom('[data-test-input="max_connection_lifetime"]')
          .exists(`Max connection lifetime exists for ${name}`);
        assert
          .dom('[data-test-input="root_rotation_statements"]')
          .exists(`Root rotation statements exists for ${name}`);
        assert
          .dom('[data-test-database-oracle-alert]')
          .hasTextContaining(
            `Warning Please ensure that your Oracle plugin has the default name of vault-plugin-database-oracle. Custom naming is not supported in the UI at this time. If the plugin is already named vault-plugin-database-oracle, disregard this warning.`,
            'warning banner displays for oracle plugin name'
          );
      },
    };
    const backend = this.backend;
    await connectionPage.visitCreate({ backend });
    assert.strictEqual(currentURL(), `/vault/secrets/${backend}/create`, 'Correct creation URL');
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('No plugin selected', 'No plugin is selected by default and empty state shows');
    await connectionPage.dbPlugin(testCase.plugin);
    assert.dom('[data-test-empty-state]').doesNotExist('Empty state goes away after plugin selected');
    assert.dom('[data-test-database-oracle-alert]').exists('shows oracle alert');
    await connectionPage.name(testCase.name);
    await connectionPage.connectionUrl(testCase.url);
    testCase.requiredFields(assert, testCase.plugin);
    // Cannot save without plugin mounted
    // TODO: add fake server response for fuller test coverage
  });

  test('Can create and delete a connection', async function (assert) {
    const backend = this.backend;
    const connectionDetails = {
      plugin: 'mongodb-database-plugin',
      id: 'horses-db',
      fields: [
        { label: 'Connection name', name: 'name', value: 'horses-db' },
        { label: 'Connection URL', name: 'connection_url', value: 'mongodb://127.0.0.1:235/horses' },
        { label: 'Username', name: 'username', value: 'user', hideOnShow: true },
        { label: 'Password', name: 'password', password: 'so-secure', hideOnShow: true },
        { label: 'Write concern', name: 'write_concern' },
      ],
    };
    await visit(`/vault/secrets/${backend}/list`);
    await connectionPage.createLink();
    assert.strictEqual(currentURL(), `/vault/secrets/${backend}/create`, 'Create link goes to create page');
    assert
      .dom('[data-test-empty-state-title]')
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
    await settled();
    assert
      .dom('.modal.is-active .title')
      .hasText('Rotate your root credentials?', 'Modal appears asking to ');
    await connectionPage.enable();
    assert.strictEqual(
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
    await connectionPage.delete();
    assert
      .dom('.modal.is-active .title')
      .hasText('Delete connection?', 'Modal appears asking to confirm delete action');
    await fillIn('[data-test-confirmation-modal-input="Delete connection?"]', connectionDetails.id);
    await click('[data-test-confirm-button]');

    assert.strictEqual(currentURL(), `/vault/secrets/${backend}/list`, 'Redirects to connection list page');
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('No connections in this backend', 'No connections listed because it was deleted');
  });

  test('buttons show up for managing connection', async function (assert) {
    const backend = this.backend;
    const connection = await newConnection(backend);
    const CONNECTION_VIEW_ONLY = `
      path "${backend}/config" {
        capabilities = ["list"]
      }
      path "${backend}/config/*" {
        capabilities = ["read"]
      }
    `;
    const token = await runCmd(tokenWithPolicyCmd('test-policy', CONNECTION_VIEW_ONLY));
    await navToConnection(backend, connection);
    assert
      .dom('[data-test-database-connection-delete]')
      .hasText('Delete connection', 'Delete connection button exists with correct text');
    assert
      .dom('[data-test-database-connection-reset]')
      .hasText('Reset connection', 'Reset button exists with correct text');
    assert.dom('[data-test-secret-create]').hasText('Add role', 'Add role button exists with correct text');
    assert.dom('[data-test-edit-link]').hasText('Edit configuration', 'Edit button exists with correct text');
    await authPage.logout();
    // Check with restricted permissions
    await authPage.login(token);
    await click('[data-test-sidebar-nav-link="Secrets engines"]');
    assert.dom(`[data-test-secrets-backend-link="${backend}"]`).exists('Shows backend on secret list page');
    await navToConnection(backend, connection);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${backend}/show/${connection}`,
      'Allows reading connection'
    );
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
    // confirm get credentials card is an option to select. Regression bug.
    await typeIn('.ember-text-field', 'blah');
    assert.dom('[data-test-get-credentials]').isEnabled();
  });

  test('Role create form', async function (assert) {
    const backend = this.backend;
    // Connection needed for role fields
    await newConnection(backend);
    await rolePage.visitCreate({ backend });
    await rolePage.name('bar');
    assert
      .dom('[data-test-component="empty-state"]')
      .exists({ count: 2 }, 'Two empty states exist before selections made');
    await clickTrigger('#database');
    assert.strictEqual(searchSelectComponent.options.length, 1, 'list shows existing connections so far');
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
    // Real connection (actual running db) required to save role, so we aren't testing that flow yet
  });

  test('root and limited access', async function (assert) {
    const backend = this.backend;
    const NO_ROLES_POLICY = `
      path "${backend}/roles/*" {
        capabilities = ["delete"]
      }
      path "${backend}/static-roles/*" {
        capabilities = ["delete"]
      }
      path "${backend}/config/*" {
        capabilities = ["list", "create", "read", "update"]
      }
      path "${backend}/creds/*" {
        capabilities = ["list", "create", "read", "update"]
      }
    `;
    const token = await runCmd(tokenWithPolicyCmd('test-policy', NO_ROLES_POLICY));

    // test root user flow first
    await visit(`/vault/secrets/${backend}/overview`);

    assert.dom('[data-test-component="empty-state"]').exists('renders empty state');
    assert.dom('[data-test-secret-list-tab="Connections"]').exists('renders connections tab');
    assert.dom('[data-test-secret-list-tab="Roles"]').exists('renders connections tab');

    await click('[data-test-secret-create="connections"]');
    assert.strictEqual(currentURL(), `/vault/secrets/${backend}/create`);

    // Login with restricted policy
    await logout.visit();
    await authPage.login(token);
    await visit(`/vault/secrets/${backend}/overview`);
    assert.dom('[data-test-tab="overview"]').exists('renders overview tab');
    assert.dom('[data-test-secret-list-tab="Connections"]').exists('renders connections tab');
    assert
      .dom('[data-test-secret-list-tab="Roles"]')
      .doesNotExist(`does not show the roles tab because it does not have permissions`);
    assert
      .dom('[data-test-selectable-card="Connections"]')
      .exists({ count: 1 }, 'renders only the connection card');

    await click('[data-test-action-text="Configure new"]');
    assert.strictEqual(currentURL(), `/vault/secrets/${backend}/create?itemType=connection`);
  });
});
