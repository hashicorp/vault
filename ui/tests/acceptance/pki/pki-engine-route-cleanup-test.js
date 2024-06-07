/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import {
  PKI_CONFIGURE_CREATE,
  PKI_ISSUER_DETAILS,
  PKI_ISSUER_LIST,
  PKI_KEYS,
  PKI_ROLE_DETAILS,
} from 'vault/tests/helpers/pki/pki-selectors';

const OVERVIEW_BREADCRUMB = '[data-test-breadcrumbs] li:nth-of-type(2) > a';
/**
 * This test module should test that dirty route models are cleaned up when
 * the user leaves the page via cancel or breadcrumb navigation
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
    await runCmd([`delete sys/mounts/${this.mountPath}`]);
  });

  module('configuration', function () {
    test('create config', async function (assert) {
      let configs, urls, config;
      await authPage.login(this.pkiAdminToken);
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(`${GENERAL.emptyStateActions} a`);
      configs = this.store.peekAll('pki/action');
      urls = this.store.peekRecord('pki/config/urls', this.mountPath);
      config = configs.at(0);
      assert.strictEqual(configs.length, 1, 'One config model present');
      assert.false(urls.hasDirtyAttributes, 'URLs is loaded from endpoint');
      assert.true(config.hasDirtyAttributes, 'Config model is dirty');

      // Cancel button rolls it back
      await click(GENERAL.cancelButton);
      configs = this.store.peekAll('pki/action');
      urls = this.store.peekRecord('pki/config/urls', this.mountPath);
      assert.strictEqual(configs.length, 0, 'config model is rolled back on cancel');
      assert.strictEqual(urls.id, this.mountPath, 'Urls still exists on exit');

      await click(`${GENERAL.emptyStateActions} a`);
      configs = this.store.peekAll('pki/action');
      urls = this.store.peekRecord('pki/config/urls', this.mountPath);
      config = configs.at(0);
      assert.strictEqual(configs.length, 1, 'One config model present');
      assert.false(urls.hasDirtyAttributes, 'URLs is loaded from endpoint');
      assert.true(config.hasDirtyAttributes, 'Config model is dirty');

      // Exit page via link rolls it back
      await click(OVERVIEW_BREADCRUMB);
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
      await click(`${GENERAL.emptyStateActions} a`);
      await click(PKI_CONFIGURE_CREATE.optionByKey('generate-root'));
      await fillIn(GENERAL.inputByAttr('type'), 'internal');
      await fillIn(GENERAL.inputByAttr('commonName'), 'my-root-cert');
      await click(GENERAL.saveButton);
      await logout.visit();
    });

    test('create role exit via cancel', async function (assert) {
      let roles;
      await authPage.login();
      // Create PKI
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Roles'));
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'No roles exist yet');
      await click(PKI_ROLE_DETAILS.createRoleLink);
      roles = this.store.peekAll('pki/role');
      const role = roles.at(0);
      assert.strictEqual(roles.length, 1, 'New role exists');
      assert.true(role.isNew, 'Role is new model');
      await click(GENERAL.cancelButton);
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'Role is removed from store');
    });
    test('create role exit via breadcrumb', async function (assert) {
      let roles;
      await authPage.login();
      // Create PKI
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Roles'));
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'No roles exist yet');
      await click(PKI_ROLE_DETAILS.createRoleLink);
      roles = this.store.peekAll('pki/role');
      const role = roles.at(0);
      assert.strictEqual(roles.length, 1, 'New role exists');
      assert.true(role.isNew, 'Role is new model');
      await click(OVERVIEW_BREADCRUMB);
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'Role is removed from store');
    });
    test('edit role', async function (assert) {
      let roles, role;
      const roleId = 'workflow-edit-role';
      await authPage.login();
      // Create PKI
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Roles'));
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'No roles exist yet');
      await click(PKI_ROLE_DETAILS.createRoleLink);
      await fillIn(GENERAL.inputByAttr('name'), roleId);
      await click(GENERAL.saveButton);
      assert.dom(GENERAL.infoRowValue('Role name')).hasText(roleId, 'Shows correct role after create');
      roles = this.store.peekAll('pki/role');
      role = roles.at(0);
      assert.strictEqual(roles.length, 1, 'Role is created');
      assert.false(role.hasDirtyAttributes, 'Role no longer has dirty attributes');

      // Edit role
      await click(PKI_ROLE_DETAILS.editRoleLink);
      await click(GENERAL.ttl.toggle('issuerRef-toggle'));
      await fillIn(GENERAL.selectByAttr('issuerRef'), 'foobar');
      role = this.store.peekRecord('pki/role', roleId);
      assert.true(role.hasDirtyAttributes, 'Role has dirty attrs');
      // Exit page via cancel button
      await click(GENERAL.cancelButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/roles/${roleId}/details`);
      role = this.store.peekRecord('pki/role', roleId);
      assert.false(role.hasDirtyAttributes, 'Role dirty attrs have been rolled back');

      // Edit again
      await click(PKI_ROLE_DETAILS.editRoleLink);
      await click(GENERAL.ttl.toggle('issuerRef-toggle'));
      await fillIn(GENERAL.selectByAttr('issuerRef'), 'foobar2');
      role = this.store.peekRecord('pki/role', roleId);
      assert.true(role.hasDirtyAttributes, 'Role has dirty attrs');
      // Exit page via breadcrumbs
      await click(OVERVIEW_BREADCRUMB);
      role = this.store.peekRecord('pki/role', roleId);
      assert.false(role.hasDirtyAttributes, 'Role dirty attrs have been rolled back');
    });
  });

  module('issuer routes', function () {
    test('import issuer exit via cancel', async function (assert) {
      let issuers;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Issuers'));
      issuers = this.store.peekAll('pki/issuer');
      assert.strictEqual(issuers.length, 0, 'No issuer models exist yet');
      await click(PKI_ISSUER_LIST.importIssuerLink);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 1, 'Action model created');
      const issuer = issuers.at(0);
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
      await click(GENERAL.secretTab('Issuers'));
      issuers = this.store.peekAll('pki/issuer');
      assert.strictEqual(issuers.length, 0, 'No issuers exist yet');
      await click(PKI_ISSUER_LIST.importIssuerLink);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 1, 'Action model created');
      const issuer = issuers.at(0);
      assert.true(issuer.hasDirtyAttributes, 'Action model has dirty attrs');
      assert.true(issuer.isNew, 'Action model is new');
      // Exit
      await click(OVERVIEW_BREADCRUMB);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 0, 'Issuer is removed from store');
    });
    test('generate root exit via cancel', async function (assert) {
      let actions;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Issuers'));
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'No actions exist yet');
      await click(PKI_ISSUER_LIST.generateIssuerDropdown);
      await click(PKI_ISSUER_LIST.generateIssuerRoot);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 1, 'Action model for generate-root created');
      const action = actions.at(0);
      assert.true(action.hasDirtyAttributes, 'Action has dirty attrs');
      assert.true(action.isNew, 'Action is new');
      assert.strictEqual(action.actionType, 'generate-root', 'Action type is correct');
      // Exit
      await click(GENERAL.cancelButton);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/issuers`);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'Action is removed from store');
    });
    test('generate root exit via breadcrumb', async function (assert) {
      let actions;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Issuers'));
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'No actions exist yet');
      await click(PKI_ISSUER_LIST.generateIssuerDropdown);
      await click(PKI_ISSUER_LIST.generateIssuerRoot);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 1, 'Action model for generate-root created');
      const action = actions.at(0);
      assert.true(action.hasDirtyAttributes, 'Action has dirty attrs');
      assert.true(action.isNew, 'Action is new');
      assert.strictEqual(action.actionType, 'generate-root');
      // Exit
      await click(OVERVIEW_BREADCRUMB);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'Action is removed from store');
    });
    test('generate intermediate csr exit via cancel', async function (assert) {
      let actions;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Issuers'));
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'No actions exist yet');
      await await click(PKI_ISSUER_LIST.generateIssuerDropdown);
      await click(PKI_ISSUER_LIST.generateIssuerIntermediate);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 1, 'Action model for generate-csr created');
      const action = actions.at(0);
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
      await click(GENERAL.secretTab('Issuers'));
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'No actions exist yet');
      await click(PKI_ISSUER_LIST.generateIssuerDropdown);
      await click(PKI_ISSUER_LIST.generateIssuerIntermediate);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 1, 'Action model for generate-csr created');
      const action = actions.at(0);
      assert.true(action.hasDirtyAttributes, 'Action has dirty attrs');
      assert.true(action.isNew, 'Action is new');
      assert.strictEqual(action.actionType, 'generate-csr');
      // Exit
      await click(OVERVIEW_BREADCRUMB);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`);
      actions = this.store.peekAll('pki/action');
      assert.strictEqual(actions.length, 0, 'Action is removed from store');
    });
    test('edit issuer exit', async function (assert) {
      let issuers, issuer;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(`${GENERAL.emptyStateActions} a`);
      await click(PKI_CONFIGURE_CREATE.optionByKey('generate-root'));
      await fillIn(GENERAL.inputByAttr('type'), 'internal');
      await fillIn(GENERAL.inputByAttr('commonName'), 'my-root-cert');
      await click(GENERAL.saveButton);
      // Go to list view so we fetch all the issuers
      await visit(`/vault/secrets/${this.mountPath}/pki/issuers`);

      issuers = this.store.peekAll('pki/issuer');
      const issuerId = issuers.at(0).id;
      assert.strictEqual(issuers.length, 1, 'Issuer exists on model in list');
      await visit(`/vault/secrets/${this.mountPath}/pki/issuers/${issuerId}/details`);
      await click(PKI_ISSUER_DETAILS.configure);
      issuer = this.store.peekRecord('pki/issuer', issuerId);
      assert.false(issuer.hasDirtyAttributes, 'Model not dirty');
      await fillIn('[data-test-input="issuerName"]', 'foobar');
      assert.true(issuer.hasDirtyAttributes, 'Model is dirty');
      await click(OVERVIEW_BREADCRUMB);
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
      await click(`${GENERAL.emptyStateActions} a`);
      await click(PKI_CONFIGURE_CREATE.optionByKey('generate-root'));
      await fillIn(GENERAL.inputByAttr('type'), 'internal');
      await fillIn(GENERAL.inputByAttr('commonName'), 'my-root-cert');
      await click(GENERAL.saveButton);
      await logout.visit();
    });
    test('create key exit', async function (assert) {
      let keys, key;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Keys'));
      keys = this.store.peekAll('pki/key');
      const configKeyId = keys.at(0).id;
      assert.strictEqual(keys.length, 1, 'One key exists from config');
      // Create key
      await click(PKI_KEYS.generateKey);
      keys = this.store.peekAll('pki/key');
      key = keys.at(1);
      assert.strictEqual(keys.length, 2, 'New key exists');
      assert.true(key.isNew, 'Role is new model');
      // Exit
      await click(GENERAL.cancelButton);
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'Second key is removed from store');
      assert.strictEqual(keys.at(0).id, configKeyId);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/keys`, 'url is correct');

      // Create again
      await click(PKI_KEYS.generateKey);
      assert.strictEqual(keys.length, 2, 'New key exists');
      keys = this.store.peekAll('pki/key');
      key = keys.at(1);
      assert.true(key.isNew, 'Key is new model');
      // Exit
      await click(OVERVIEW_BREADCRUMB);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`, 'url is correct');
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'Key is removed from store');
    });
    test('edit key exit', async function (assert) {
      let keys, key;
      await authPage.login();
      await visit(`/vault/secrets/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Keys'));
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'One key from config exists');
      assert.dom('.list-item-row').exists({ count: 1 }, 'single row for key');
      await click('.list-item-row');
      // Edit
      await click(PKI_KEYS.keyEditLink);
      await fillIn(GENERAL.inputByAttr('keyName'), 'foobar');
      keys = this.store.peekAll('pki/key');
      key = keys.at(0);
      assert.true(key.hasDirtyAttributes, 'Key model is dirty');
      // Exit
      await click(GENERAL.cancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.mountPath}/pki/keys/${key.id}/details`,
        'url is correct'
      );
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'Key list has 1');
      assert.false(key.hasDirtyAttributes, 'Key dirty attrs have been rolled back');

      // Edit again
      await click(PKI_KEYS.keyEditLink);
      await fillIn(GENERAL.inputByAttr('keyName'), 'foobar');
      keys = this.store.peekAll('pki/key');
      key = keys.at(0);
      assert.true(key.hasDirtyAttributes, 'Key model is dirty');

      // Exit via breadcrumb
      await click(OVERVIEW_BREADCRUMB);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.mountPath}/pki/overview`, 'url is correct');
      keys = this.store.peekAll('pki/key');
      assert.strictEqual(keys.length, 1, 'Key list has 1');
      assert.false(key.hasDirtyAttributes, 'Key dirty attrs have been rolled back');
    });
  });
});
