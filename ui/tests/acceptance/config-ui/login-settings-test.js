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

module('Acceptance | Enterprise | config-ui/login-settings', function (hooks) {
  setupApplicationTest(hooks);

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
    assert.dom('.linked-block-item').exists({ count: 2 });
    // verify rule data namespaces render
    assert.dom('[data-test-rule-path="ns1/"]').exists();
    assert.dom('[data-test-rule-path="ns2/"]').exists();
  });

  test('delete rule from list view', async function (assert) {
    // Visit the login settings list index page
    await visit('vault/config-ui/login-settings');

    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete-rule'));

    assert.dom(GENERAL.confirmationModal).exists();
    await click(GENERAL.confirmButton);

    // verify success message from deletion
    assert.dom(GENERAL.latestFlashContent).includesText('Successfully deleted rule testRule.');
    assert.dom('[data-test-rule-name="testRule"]').doesNotExist();
  });

  test('navigate to rule details page and renders rule data', async function (assert) {
    // visit individual rule page
    await visit('vault/config-ui/login-settings');

    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('view-rule'));

    // verify that user is redirected to the rule details page
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.config-ui.login-settings.rule.details',
      'goes to rule details page'
    );

    // verify fetched rule data is rendered
    assert.dom(GENERAL.infoRowValue('Name')).hasText('testRule');
    assert.dom(GENERAL.infoRowValue('Namespace')).hasText('ns1/');
    assert.dom(GENERAL.infoRowValue('Backup methods')).hasText('userpass');
    assert.dom(GENERAL.infoRowValue('Inheritance')).hasText('true');
  });
});
