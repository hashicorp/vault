/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, visit, currentRouteName } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { runCmd } from 'vault/tests/helpers/commands';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';

const SELECTORS = {
  rule: (name) => (name ? `[data-test-rule="${name}"]` : '[data-test-rule]'),
  popupMenu: (name) => `[data-test-rule="${name}"] ${GENERAL.menuTrigger}`,
};
// read view for custom login settings
module('Acceptance | Enterprise | config-ui/login-settings', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    return await login();
  });

  test('it renders empty state if no login settings exist', async function (assert) {
    await visit('vault/config-ui/login-settings');

    assert.dom(GENERAL.emptyStateTitle).hasText('No UI login rules yet');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Login rules can be used to select default and back up login methods and customize which methods display in the web UI login form. Available to be created via the CLI or HTTP API.'
      );
  });

  test('it falls back error template if no permission', async function (assert) {
    this.server.get('/sys/config/ui/login/default-auth', () => overrideResponse(403));
    await visit('vault/config-ui/login-settings');
    assert.dom(GENERAL.pageError.error).hasText('Error permission denied');
  });

  module('list, read and delete', function (hooks) {
    hooks.beforeEach(async function () {
      await login();

      // create login rules
      await runCmd([
        `write sys/config/ui/login/default-auth/testRule backup_auth_types=userpass default_auth_type=okta disable_inheritance=false namespace=ns1`,
        'write sys/config/ui/login/default-auth/testRule2 backup_auth_types=oidc default_auth_type=ldap disable_inheritance=true namespace=ns2',
      ]);
    });

    hooks.afterEach(async function () {
      await login();

      // cleanup login rules
      await runCmd([
        'delete sys/config/ui/login/default-auth/testRule',
        'delete sys/config/ui/login/default-auth/testRule2',
      ]);
    });

    test('fetched login rule list renders', async function (assert) {
      // Visit the login settings list index page
      await visit('vault/config-ui/login-settings');

      // verify fetched rules are rendered in list
      assert.dom(SELECTORS.rule()).exists({ count: 2 });
      assert.dom(SELECTORS.rule('testRule')).hasText('testRule ns1/ Inheritance enabled');
      assert.dom(SELECTORS.rule('testRule2')).hasText('testRule2 ns2/ Inheritance disabled');
    });

    test('delete rule from list view', async function (assert) {
      // Visit the login settings list index page
      await visit('vault/config-ui/login-settings');

      assert.dom(SELECTORS.rule()).exists({ count: 2 });
      await click(SELECTORS.popupMenu('testRule'));
      await click(GENERAL.menuItem('delete-rule'));

      assert.dom(GENERAL.confirmationModal).exists();
      await click(GENERAL.confirmButton);

      // verify success message from deletion
      assert.dom(GENERAL.latestFlashContent).includesText('Successfully deleted rule testRule.');
      assert.dom(SELECTORS.rule('testRule')).doesNotExist();
      assert.dom(SELECTORS.rule()).exists({ count: 1 });
    });

    test('navigate to rule details page and renders rule data', async function (assert) {
      // visit individual rule page
      await visit('vault/config-ui/login-settings');

      await click(SELECTORS.popupMenu('testRule'));
      await click(GENERAL.menuItem('view-rule'));

      // verify that user is redirected to the rule details page
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.login-settings.rule.details',
        'goes to rule details page'
      );

      // verify fetched rule data is rendered
      assert.dom(GENERAL.infoRowValue('Default method')).hasText('okta');
      assert.dom(GENERAL.infoRowValue('Namespace')).hasText('ns1/');
      assert.dom(GENERAL.infoRowValue('Backup methods')).hasText('userpass');
      assert.dom(GENERAL.infoRowValue('Inheritance enabled')).hasText('Yes');
    });

    test('it navigates to rule details from linked block', async function (assert) {
      await visit('vault/config-ui/login-settings');
      await click(SELECTORS.rule('testRule2'));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.login-settings.rule.details',
        'goes to rule details page'
      );

      assert.dom(GENERAL.infoRowValue('Default method')).hasText('ldap');
      assert.dom(GENERAL.infoRowValue('Namespace')).hasText('ns2/');
      assert.dom(GENERAL.infoRowValue('Backup methods')).hasText('oidc');
      assert.dom(GENERAL.infoRowValue('Inheritance enabled')).hasText('No');
    });
  });
});
