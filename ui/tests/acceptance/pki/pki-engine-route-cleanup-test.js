/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import {
  PKI_CONFIGURE_CREATE,
  PKI_ISSUER_LIST,
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
    await login();
    // Setup PKI engine
    const mountPath = `pki-workflow-${uuidv4()}`;
    await enablePage.enable('pki', mountPath);
    this.mountPath = mountPath;
  });

  hooks.afterEach(async function () {
    await login();
    // Cleanup engine
    await runCmd([`delete sys/mounts/${this.mountPath}`]);
  });

  module('role routes', function (hooks) {
    hooks.beforeEach(async function () {
      await login();
      // Configure PKI
      await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
      await click(`${GENERAL.emptyStateActions} a`);
      await click(PKI_CONFIGURE_CREATE.optionByKey('generate-root'));
      await fillIn(GENERAL.inputByAttr('type'), 'internal');
      await fillIn(GENERAL.inputByAttr('common_name'), 'my-root-cert');
      await click(GENERAL.submitButton);
    });

    test('create role exit via cancel', async function (assert) {
      let roles;
      await login();
      // Create PKI
      await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
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
      await login();
      // Create PKI
      await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
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
      await login();
      // Create PKI
      await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
      await click(GENERAL.secretTab('Roles'));
      roles = this.store.peekAll('pki/role');
      assert.strictEqual(roles.length, 0, 'No roles exist yet');
      await click(PKI_ROLE_DETAILS.createRoleLink);
      await fillIn(GENERAL.inputByAttr('name'), roleId);
      await click(GENERAL.submitButton);
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
      assert.strictEqual(
        currentURL(),
        `/vault/secrets-engines/${this.mountPath}/pki/roles/${roleId}/details`
      );
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
      await login();
      await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
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
      assert.strictEqual(currentURL(), `/vault/secrets-engines/${this.mountPath}/pki/issuers`);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 0, 'Action is removed from store');
    });
    test('import issuer exit via breadcrumb', async function (assert) {
      let issuers;
      await login();
      await visit(`/vault/secrets-engines/${this.mountPath}/pki/overview`);
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
      assert.strictEqual(currentURL(), `/vault/secrets-engines/${this.mountPath}/pki/overview`);
      issuers = this.store.peekAll('pki/action');
      assert.strictEqual(issuers.length, 0, 'Issuer is removed from store');
    });
  });
});
