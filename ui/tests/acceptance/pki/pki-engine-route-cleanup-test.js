/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { SELECTORS } from 'vault/tests/helpers/pki/workflow';

/**
 * This test module should test that dirty route models are cleaned up when the user leaves the page
 */
module('Acceptance | pki engine route cleanup test', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    await authPage.login();
    // Setup PKI engine
    const mountPath = `pki-workflow-${uuidv4()}`;
    await enablePage.enable('pki', mountPath);
    this.mountPath = mountPath;
    await logout.visit();
  });

  hooks.afterEach(async function () {
    await logout.visit();
    await authPage.login();
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.mountPath}`]);
  });

  module('configuration', function () {
    test('create config', async function (assert) {
      let configs, urls, config;
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.emptyStateLink);
      configs = this.store.peekAll('pki/action');
      urls = this.store.peekRecord('pki/config/urls', this.mountPath);
      config = configs.objectAt(0);
      assert.strictEqual(configs.length, 1, 'One config model present');
      assert.false(urls.hasDirtyAttributes, 'URLs is loaded from endpoint');
      assert.true(config.hasDirtyAttributes, 'Config model is dirty');

      // Cancel button rolls it back
      await click(SELECTORS.configuration.cancelButton);
      configs = this.store.peekAll('pki/action');
      urls = this.store.peekRecord('pki/config/urls', this.mountPath);
      assert.strictEqual(configs.length, 0, 'config model is rolled back on cancel');
      assert.strictEqual(urls.id, this.mountPath, 'Urls still exists on exit');

      await click(SELECTORS.emptyStateLink);
      configs = this.store.peekAll('pki/action');
      urls = this.store.peekRecord('pki/config/urls', this.mountPath);
      config = configs.objectAt(0);
      assert.strictEqual(configs.length, 1, 'One config model present');
      assert.false(urls.hasDirtyAttributes, 'URLs is loaded from endpoint');
      assert.true(config.hasDirtyAttributes, 'Config model is dirty');

      // Exit page via link rolls it back
      await click(SELECTORS.overviewBreadcrumb);
      configs = this.store.peekAll('pki/action');
      urls = this.store.peekRecord('pki/config/urls', this.mountPath);
      assert.strictEqual(configs.length, 0, 'config model is rolled back on cancel');
      assert.strictEqual(urls.id, this.mountPath, 'Urls still exists on exit');
    });
  });

  module('role routes', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // Configure PKI
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.emptyStateLink);
      await click(SELECTORS.configuration.optionByKey('generate-root'));
      await fillIn(SELECTORS.configuration.typeField, 'internal');
      await fillIn(SELECTORS.configuration.inputByName('commonName'), 'my-root-cert');
      await click(SELECTORS.configuration.generateRootSave);
      await logout.visit();
    });

    test('create role exit via cancel', async function (assert) {
      let roles;
      await authPage.login();
      // Create PKI
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.rolesTab);
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'No roles exist yet');
      await click(SELECTORS.createRoleLink);
      roles = this.store.peekAll('pki/role');
      const role = roles.objectAt(0);
      assert.strictEqual(roles.length, 1, 'New role exists');
      assert.true(role.isNew, 'Role is new model');
      await click(SELECTORS.roleForm.roleCancelButton);
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'Role is removed from store');
    });
    test('create role exit via breadcrumb', async function (assert) {
      let roles;
      await authPage.login();
      // Create PKI
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.rolesTab);
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'No roles exist yet');
      await click(SELECTORS.createRoleLink);
      roles = this.store.peekAll('pki/role');
      const role = roles.objectAt(0);
      assert.strictEqual(roles.length, 1, 'New role exists');
      assert.true(role.isNew, 'Role is new model');
      await click(SELECTORS.overviewBreadcrumb);
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'Role is removed from store');
    });
    test('edit role', async function (assert) {
      let roles, role;
      const roleId = 'workflow-edit-role';
      await authPage.login();
      // Create PKI
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.rolesTab);
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'No roles exist yet');
      await click(SELECTORS.createRoleLink);
      await fillIn(SELECTORS.roleForm.roleName, roleId);
      await click(SELECTORS.roleForm.roleCreateButton);
      assert.dom('[data-test-value-div="Role name"]').hasText(roleId, 'Shows correct role after create');
      roles = this.store.peekAll('pki/role');
      role = roles.objectAt(0);
      assert.strictEqual(roles.length, 1, 'Role is created');
      assert.false(role.hasDirtyAttributes, 'Role no longer has dirty attributes');

      // Edit role
      await click(SELECTORS.editRoleLink);
      await click(SELECTORS.roleForm.issuerRefToggle);
      await fillIn(SELECTORS.roleForm.issuerRefSelect, 'foobar');
      role = this.store.peekRecord('pki/role', roleId);
      assert.true(role.hasDirtyAttributes, 'Role has dirty attrs');
      // Exit page via cancel button
      await click(SELECTORS.roleForm.roleCancelButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/${roleId}/details`);
      role = this.store.peekRecord('pki/role', roleId);
      assert.false(role.hasDirtyAttributes, 'Role dirty attrs have been rolled back');

      // Edit again
      await click(SELECTORS.editRoleLink);
      await click(SELECTORS.roleForm.issuerRefToggle);
      await fillIn(SELECTORS.roleForm.issuerRefSelect, 'foobar2');
      role = this.store.peekRecord('pki/role', roleId);
      assert.true(role.hasDirtyAttributes, 'Role has dirty attrs');
      // Exit page via breadcrumbs
      await click(SELECTORS.overviewBreadcrumb);
      role = this.store.peekRecord('pki/role', roleId);
      assert.false(role.hasDirtyAttributes, 'Role dirty attrs have been rolled back');
    });
  });

  module('issuer routes', function () {
    test('import issuer exit via cancel', async function (assert) {
      let issuers;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      issuers = this.store.peekAll('pki/issuer');
      assert.strictEqual(issuers.length, 0, 'No issuer models exist yet');
      await click(SELECTORS.importIssuerLink);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 1, 'Action model created');
      const issuer = issuers.objectAt(0);
      assert.true(issuer.hasDirtyAttributes, 'Action has dirty attrs');
      assert.true(issuer.isNew, 'Action is new');
      // Exit
      await click('[data-test-pki-ca-cert-cancel]');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 0, 'Action is removed from store');
    });
    test('import issuer exit via breadcrumb', async function (assert) {
      let issuers;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      issuers = this.store.peekAll('pki/issuer');
      assert.strictEqual(issuers.length, 0, 'No issuers exist yet');
      await click(SELECTORS.importIssuerLink);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 1, 'Action model created');
      const issuer = issuers.objectAt(0);
      assert.true(issuer.hasDirtyAttributes, 'Action model has dirty attrs');
      assert.true(issuer.isNew, 'Action model is new');
      // Exit
      await click(SELECTORS.overviewBreadcrumb);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 0, 'Issuer is removed from store');
    });
    test('generate root exit via cancel', async function (assert) {
      let actions;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'No actions exist yet');
      await click(SELECTORS.generateIssuerDropdown);
      await click(SELECTORS.generateIssuerRoot);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 1, 'Action model for generate-root created');
      const action = actions.objectAt(0);
      assert.true(action.hasDirtyAttributes, 'Action has dirty attrs');
      assert.true(action.isNew, 'Action is new');
      assert.strictEqual(action.actionType, 'generate-root', 'Action type is correct');
      // Exit
      await click(SELECTORS.configuration.generateRootCancel);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'Action is removed from store');
    });
    test('generate root exit via breadcrumb', async function (assert) {
      let actions;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'No actions exist yet');
      await click(SELECTORS.generateIssuerDropdown);
      await click(SELECTORS.generateIssuerRoot);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 1, 'Action model for generate-root created');
      const action = actions.objectAt(0);
      assert.true(action.hasDirtyAttributes, 'Action has dirty attrs');
      assert.true(action.isNew, 'Action is new');
      assert.strictEqual(action.actionType, 'generate-root');
      // Exit
      await click(SELECTORS.overviewBreadcrumb);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'Action is removed from store');
    });
    test('generate intermediate csr exit via cancel', async function (assert) {
      let actions;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'No actions exist yet');
      await await click(SELECTORS.generateIssuerDropdown);
      await click(SELECTORS.generateIssuerIntermediate);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 1, 'Action model for generate-csr created');
      const action = actions.objectAt(0);
      assert.true(action.hasDirtyAttributes, 'Action has dirty attrs');
      assert.true(action.isNew, 'Action is new');
      assert.strictEqual(action.actionType, 'generate-csr');
      // Exit
      await click('[data-test-cancel]');
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'Action is removed from store');
    });
    test('generate intermediate csr exit via breadcrumb', async function (assert) {
      let actions;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.issuersTab);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'No actions exist yet');
      await click(SELECTORS.generateIssuerDropdown);
      await click(SELECTORS.generateIssuerIntermediate);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 1, 'Action model for generate-csr created');
      const action = actions.objectAt(0);
      assert.true(action.hasDirtyAttributes, 'Action has dirty attrs');
      assert.true(action.isNew, 'Action is new');
      assert.strictEqual(action.actionType, 'generate-csr');
      // Exit
      await click(SELECTORS.overviewBreadcrumb);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'Action is removed from store');
    });
    test('edit issuer exit', async function (assert) {
      let issuers, issuer;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.emptyStateLink);
      await click(SELECTORS.configuration.optionByKey('generate-root'));
      await fillIn(SELECTORS.configuration.typeField, 'internal');
      await fillIn(SELECTORS.configuration.inputByName('commonName'), 'my-root-cert');
      await click(SELECTORS.configuration.generateRootSave);
      // Go to list view so we fetch all the issuers
      await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);

      issuers = this.store.peekAll('pki/issuer');
      const issuerId = issuers.objectAt(0).id;
      assert.strictEqual(issuers.length, 1, 'Issuer exists on model in list');
      await visit(`/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/details`);
      await click(SELECTORS.issuerDetails.configure);
      issuer = this.store.peekRecord('pki/issuer', issuerId);
      assert.false(issuer.hasDirtyAttributes, 'Model not dirty');
      await fillIn('[data-test-input="issuerName"]', 'foobar');
      assert.true(issuer.hasDirtyAttributes, 'Model is dirty');
      await click(SELECTORS.overviewBreadcrumb);
      issuers = this.store.peekAll('pki/issuer');
      assert.strictEqual(issuers.length, 1, 'Issuer exists on model in overview');
      issuer = this.store.peekRecord('pki/issuer', issuerId);
      assert.false(issuer.hasDirtyAttributes, 'Dirty attrs were rolled back');
    });
  });

  module('key routes', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      // Configure PKI -- key creation not allowed unless configured
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.emptyStateLink);
      await click(SELECTORS.configuration.optionByKey('generate-root'));
      await fillIn(SELECTORS.configuration.typeField, 'internal');
      await fillIn(SELECTORS.configuration.inputByName('commonName'), 'my-root-cert');
      await click(SELECTORS.configuration.generateRootSave);
      await logout.visit();
    });
    test('create key exit', async function (assert) {
      let keys, key;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.keysTab);
      keys = this.store.peekAll('pki/key');
      const configKeyId = keys.objectAt(0).id;
      assert.strictEqual(keys.length, 1, 'One key exists from config');
      // Create key
      await click(SELECTORS.keyPages.generateKey);
      keys = this.store.peekAll('pki/key');
      key = keys.objectAt(1);
      assert.strictEqual(keys.length, 2, 'New key exists');
      assert.true(key.isNew, 'Role is new model');
      // Exit
      await click(SELECTORS.keyForm.keyCancelButton);
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'Second key is removed from store');
      assert.strictEqual(keys.objectAt(0).id, configKeyId);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`, 'url is correct');

      // Create again
      await click(SELECTORS.keyPages.generateKey);
      assert.strictEqual(keys.length, 2, 'New key exists');
      keys = this.store.peekAll('pki/key');
      key = keys.objectAt(1);
      assert.true(key.isNew, 'Key is new model');
      // Exit
      await click(SELECTORS.overviewBreadcrumb);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`, 'url is correct');
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'Key is removed from store');
    });
    test('edit key exit', async function (assert) {
      let keys, key;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(SELECTORS.keysTab);
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'One key from config exists');
      assert.dom('.list-item-row').exists({ count: 1 }, 'single row for key');
      await click('.list-item-row');
      // Edit
      await click(SELECTORS.keyPages.keyEditLink);
      await fillIn(SELECTORS.keyForm.keyNameInput, 'foobar');
      keys = this.store.peekAll('pki/key');
      key = keys.objectAt(0);
      assert.true(key.hasDirtyAttributes, 'Key model is dirty');
      // Exit
      await click(SELECTORS.keyForm.keyCancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/${key.id}/details`,
        'url is correct'
      );
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'Key list has 1');
      assert.false(key.hasDirtyAttributes, 'Key dirty attrs have been rolled back');

      // Edit again
      await click(SELECTORS.keyPages.keyEditLink);
      await fillIn(SELECTORS.keyForm.keyNameInput, 'foobar');
      keys = this.store.peekAll('pki/key');
      key = keys.objectAt(0);
      assert.true(key.hasDirtyAttributes, 'Key model is dirty');

      // Exit via breadcrumb
      await click(SELECTORS.overviewBreadcrumb);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`, 'url is correct');
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'Key list has 1');
      assert.false(key.hasDirtyAttributes, 'Key dirty attrs have been rolled back');
    });
  });
});
