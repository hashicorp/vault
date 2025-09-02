/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { Response } from 'miragejs';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { create } from 'ember-cli-page-object';

import databaseHandlers from 'vault/mirage/handlers/database';
import { setupApplicationTest } from 'vault/tests/helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import flashMessage from 'vault/tests/pages/components/flash-message';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

const flash = create(flashMessage);

const PAGE = {
  // GENERIC
  emptyStateTitle: '[data-test-empty-state-title]',
  infoRow: '[data-test-component="info-table-row"]',
  // CONNECTIONS
  rotateModal: '[data-test-db-connection-modal-title]',
  confirmRotate: '[data-test-enable-rotate-connection]',
  skipRotate: '[data-test-enable-connection]',
  // ROLES
  addRole: '[data-test-add-role]',
  roleSettingsSection: '[data-test-role-settings-section]',
  statementsSection: '[data-test-statements-section]',
  editRole: '[data-test-edit-link]',
};

const FORM = {
  creationStatement: (idx = 0) =>
    `[data-test-input="creation_statements"] [data-test-string-list-input="${idx}"]`,
};

async function fillOutConnection(name) {
  await fillIn(GENERAL.inputByAttr('name'), name);
  await fillIn(GENERAL.inputByAttr('plugin_name'), 'mysql-database-plugin');
  await fillIn(GENERAL.inputByAttr('connection_url'), '{{username}}:{{password}}@tcp(127.0.0.1:33060)/');
  await fillIn(GENERAL.inputByAttr('username'), 'admin');
  await fillIn(GENERAL.inputByAttr('password'), 'very-secure');
}

/**
 * This test set is for testing the flow for database secrets engine.
 */
module('Acceptance | database workflow', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    databaseHandlers(this.server);
    this.backend = `db-workflow-${uuidv4()}`;
    this.store = this.owner.lookup('service:store');
    await login();
    await runCmd(mountEngineCmd('database', this.backend), false);
  });

  hooks.afterEach(async function () {
    await login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  module('connections', function (hooks) {
    hooks.beforeEach(function () {
      this.expectedRows = [
        { label: 'Database plugin', value: 'mysql-database-plugin' },
        { label: 'Connection name', value: `connect-${this.backend}` },
        { label: 'Use custom password policy', value: 'Default' },
        { label: 'Connection URL', value: '{{username}}:{{password}}@tcp(127.0.0.1:33060)/' },
        { label: 'Max open connections', value: '4' },
        { label: 'Max idle connections', value: '0' },
        { label: 'Max connection lifetime', value: '0s' },
        { label: 'Username template', value: 'Default' },
        {
          label: 'Root rotation statements',
          value: `Default`,
        },
        { label: 'Rotate static roles immediately', value: 'Yes' },
      ];
    });
    test('create with rotate', async function (assert) {
      assert.expect(26);
      this.server.post('/:backend/rotate-root/:name', () => {
        assert.ok(true, 'rotate root called');
        new Response(204);
      });
      await visit(`/vault/secrets/${this.backend}/overview`);
      assert.dom(PAGE.emptyStateTitle).hasText('Connect a database', 'empty state title is correct');
      await click(SES.createSecretLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/create?itemType=connection`,
        'Takes you to create page'
      );

      // fill in connection details
      await fillOutConnection(`connect-${this.backend}`);
      await click(GENERAL.submitButton);

      assert.dom(PAGE.rotateModal).hasText('Rotate your root credentials?', 'rotate modal is shown');
      await click(PAGE.confirmRotate);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/connect-${this.backend}`,
        'Takes you to details page for connection'
      );
      assert.dom(PAGE.infoRow).exists({ count: this.expectedRows.length }, 'correct number of rows');
      this.expectedRows.forEach(({ label, value }) => {
        assert.dom(GENERAL.infoRowLabel(label)).hasText(label, `Label for ${label} is correct`);
        assert.dom(GENERAL.infoRowValue(label)).hasText(value, `Value for ${label} is correct`);
      });
    });
    test('create without rotate', async function (assert) {
      assert.expect(25);
      this.server.post('/:backend/rotate-root/:name', () => {
        assert.notOk(true, 'rotate root called when it should not have been');
        new Response(204);
      });
      await visit(`/vault/secrets/${this.backend}/overview`);
      assert.dom(PAGE.emptyStateTitle).hasText('Connect a database', 'empty state title is correct');
      await click(SES.createSecretLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/create?itemType=connection`,
        'Takes you to create page'
      );

      // fill in connection details
      await fillOutConnection(`connect-${this.backend}`);
      await click(GENERAL.submitButton);

      assert.dom(PAGE.rotateModal).hasText('Rotate your root credentials?', 'rotate modal is shown');
      await click(PAGE.skipRotate);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/connect-${this.backend}`,
        'Takes you to details page for connection'
      );
      assert.dom(PAGE.infoRow).exists({ count: this.expectedRows.length }, 'correct number of rows');
      this.expectedRows.forEach(({ label, value }) => {
        assert.dom(GENERAL.infoRowLabel(label)).hasText(label, `Label for ${label} is correct`);
        assert.dom(GENERAL.infoRowValue(label)).hasText(value, `Value for ${label} is correct`);
      });
    });
    test('create failure', async function (assert) {
      assert.expect(27);
      this.server.post('/:backend/rotate-root/:name', (schema, req) => {
        const okay = req.params.name !== 'bad-connection';
        assert.ok(okay, 'rotate root called but not for bad-connection');
        new Response(204);
      });
      await visit(`/vault/secrets/${this.backend}/overview`);
      assert.dom(PAGE.emptyStateTitle).hasText('Connect a database', 'empty state title is correct');
      await click(SES.createSecretLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/create?itemType=connection`,
        'Takes you to create page'
      );

      // fill in connection details
      await fillOutConnection(`bad-connection`);
      await click(GENERAL.submitButton);
      assert.strictEqual(
        flash.latestMessage,
        `error creating database object: error verifying - ping: Error 1045 (28000): Access denied for user 'admin'@'192.168.65.1' (using password: YES)`,
        'shows the error message from API'
      );
      await fillIn(GENERAL.inputByAttr('name'), `connect-${this.backend}`);
      await click(GENERAL.submitButton);
      assert.dom(PAGE.rotateModal).hasText('Rotate your root credentials?', 'rotate modal is shown');
      await click(PAGE.confirmRotate);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/connect-${this.backend}`,
        'Takes you to details page for connection'
      );
      assert.dom(PAGE.infoRow).exists({ count: this.expectedRows.length }, 'correct number of rows');
      this.expectedRows.forEach(({ label, value }) => {
        assert.dom(GENERAL.infoRowLabel(label)).hasText(label, `Label for ${label} is correct`);
        assert.dom(GENERAL.infoRowValue(label)).hasText(value, `Value for ${label} is correct`);
      });
    });

    test('create connection with rotate failure', async function (assert) {
      await visit(`/vault/secrets/${this.backend}/overview`);
      assert.dom(PAGE.emptyStateTitle).hasText('Connect a database', 'empty state title is correct');
      await click(SES.createSecretLink);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/create?itemType=connection`,
        'Takes you to create page'
      );

      // fill in connection details
      await fillOutConnection(`fail-rotate`);
      await click(GENERAL.submitButton);
      assert.dom(PAGE.rotateModal).hasText('Rotate your root credentials?', 'rotate modal is shown');
      await click(PAGE.confirmRotate);

      assert.strictEqual(
        flash.latestMessage,
        `Error rotating root credentials: 1 error occurred: * failed to update user: failed to change password: Error 1045 (28000): Access denied for user 'admin'@'%' (using password: YES)`,
        'shows the error message from API'
      );
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/fail-rotate`,
        'Takes you to details page for connection'
      );
    });
  });
  module('dynamic roles', function (hooks) {
    hooks.beforeEach(async function () {
      this.connection = `connect-${this.backend}`;
      await visit(`/vault/secrets/${this.backend}/create`);
      await fillOutConnection(this.connection);
      await click(GENERAL.submitButton);
      await visit(`/vault/secrets/${this.backend}/show/${this.connection}`);
    });

    test('it creates a dynamic role attached to the current connection', async function (assert) {
      const roleName = 'dynamic-role';
      await click(PAGE.addRole);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/create?initialKey=${this.connection}&itemType=role`,
        'Takes you to create role page'
      );

      assert
        .dom(`${PAGE.roleSettingsSection} ${PAGE.emptyStateTitle}`)
        .hasText('No role type selected', 'roles section shows empty state before selecting role type');
      assert
        .dom(`${PAGE.statementsSection} ${PAGE.emptyStateTitle}`)
        .hasText('No role type selected', 'statements section shows empty state before selecting role type');

      await fillIn(GENERAL.inputByAttr('name'), roleName);
      assert.dom('[data-test-selected-option]').hasText(this.connection, 'Connection is selected by default');

      await fillIn(GENERAL.inputByAttr('type'), 'dynamic');
      assert
        .dom(`${PAGE.roleSettingsSection} ${PAGE.emptyStateTitle}`)
        .doesNotExist('roles section no longer has empty state');
      assert
        .dom(`${PAGE.statementsSection} ${PAGE.emptyStateTitle}`)
        .doesNotExist('statements section no longer has empty state');

      // Fill in multiple creation statements
      await fillIn(FORM.creationStatement(), `GRANT SELECT ON *.* TO '{{name}}'@'%'`);
      await click(`[data-test-string-list-row="0"] [data-test-string-list-button="add"]`);
      await fillIn(FORM.creationStatement(1), `GRANT CREATE ON *.* TO '{{name}}'@'%'`);
      await click(GENERAL.submitButton);
      // DETAILS
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/role/${roleName}`,
        'Takes you to details page for role after save'
      );
      assert.dom(PAGE.infoRow).exists({ count: 7 }, 'correct number of info rows displayed');
      [
        { label: 'Role name', value: roleName },
        { label: 'Connection name', value: this.connection },
        { label: 'Type of role', value: 'dynamic' },
        { label: 'Generated credentials’s Time-to-Live (TTL)', value: '1 hour' },
        { label: 'Generated credentials’s maximum Time-to-Live (Max TTL)', value: '1 day' },
        {
          label: 'Creation statements',
          value: `GRANT SELECT ON *.* TO '{{name}}'@'%',GRANT CREATE ON *.* TO '{{name}}'@'%'`,
        },
        { label: 'Revocation statements', value: 'Default' },
      ].forEach(({ label, value }) => {
        assert.dom(GENERAL.infoRowLabel(label)).hasText(label, `Label for ${label} is correct`);
        assert.dom(GENERAL.infoRowValue(label)).hasText(value, `Value for ${label} is correct`);
      });
      // EDIT
      await click(PAGE.editRole);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/edit/role/${roleName}?itemType=role`,
        'Takes you to edit page for role'
      );
      // TODO: these should be readonly not disabled
      assert.dom(GENERAL.inputByAttr('name')).isDisabled('Name is read-only');
      assert.dom(GENERAL.inputByAttr('database')).isDisabled('Database is read-only');
      assert.dom(GENERAL.inputByAttr('type')).isDisabled('Type is read-only');
      await fillIn('[data-test-ttl-value="Generated credentials’s Time-to-Live (TTL)"]', '2');
      await click(GENERAL.submitButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/role/${roleName}`,
        'Takes you to details page for role after save'
      );
      assert
        .dom(GENERAL.infoRowValue('Generated credentials’s Time-to-Live (TTL)'))
        .hasText('2 hours', 'Shows updated TTL');

      // CREDENTIALS
      await click(GENERAL.button('dynamic'));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/credentials/${roleName}?roleType=dynamic`,
        'Takes you to credentials page for role'
      );
      assert
        .dom('[data-test-credentials-warning]')
        .exists('shows warning about credentials only being available once');
      assert
        .dom(`${GENERAL.infoRowValue('Username')} [data-test-masked-input]`)
        .hasText('***********', 'Username is masked');

      await click(`${GENERAL.infoRowValue('Username')} ${GENERAL.button('toggle-masked')}`);
      assert
        .dom(`${GENERAL.infoRowValue('Username')} [data-test-masked-input]`)
        .hasText('generated-username', 'Username is generated');

      assert
        .dom(`${GENERAL.infoRowValue('Password')} [data-test-masked-input]`)
        .hasText('***********', 'Password is masked');

      await click(`${GENERAL.infoRowValue('Password')} ${GENERAL.button('toggle-masked')}`);
      assert
        .dom(`${GENERAL.infoRowValue('Password')} [data-test-masked-input]`)
        .hasText('generated-password', 'Password is generated');
      assert
        .dom(GENERAL.infoRowValue('Lease Duration'))
        .hasText('3600', 'shows lease duration from response');
      assert
        .dom(GENERAL.infoRowValue('Lease ID'))
        .hasText(`database/creds/${roleName}/abcd`, 'shows lease ID from response');
    });
  });

  module('static roles', function (hooks) {
    hooks.beforeEach(async function () {
      this.setup = async ({ toggleRotateOff = false }) => {
        this.connection = `connect-${this.backend}`;
        await visit(`/vault/secrets/${this.backend}/create`);
        await fillOutConnection(this.connection);
        if (toggleRotateOff) {
          await click(GENERAL.toggleInput('toggle-skip_static_role_rotation_import'));
        }
        await click(GENERAL.submitButton);
        await visit(`/vault/secrets/${this.backend}/show/${this.connection}`);
      };
    });

    test('set parent db to rotate static roles immediately, verify static role reflects that default', async function (assert) {
      await this.setup({ toggleRotateOff: false });

      const roleName = 'static-role';
      await click(PAGE.addRole);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/create?initialKey=${this.connection}&itemType=role`,
        'Takes you to create role page'
      );

      await fillIn(GENERAL.inputByAttr('name'), roleName);

      await fillIn(GENERAL.inputByAttr('type'), 'static');

      assert
        .dom('[data-test-toggle-subtext]')
        .containsText(`Vault will rotate the password for this static role on creation.`);
    });

    test('set parent db to not rotate static roles immediately, verify static role reflects that default', async function (assert) {
      await this.setup({ toggleRotateOff: true });

      const roleName = 'static-role';
      await click(PAGE.addRole);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/create?initialKey=${this.connection}&itemType=role`,
        'Takes you to create role page'
      );

      await fillIn(GENERAL.inputByAttr('name'), roleName);

      await fillIn(GENERAL.inputByAttr('type'), 'static');

      assert
        .dom('[data-test-toggle-subtext]')
        .containsText(`Vault will not rotate this role's password on creation.`);
    });
  });
});
