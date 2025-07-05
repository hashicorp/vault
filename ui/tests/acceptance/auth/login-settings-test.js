/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, fillIn, typeIn, visit, waitFor } from '@ember/test-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import { login, logout, rootToken } from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// Auth form login settings
// This feature has thorough integration test coverage so only testing a few scenarios and direct link functionality
// Tests for read/list views are in ui/tests/acceptance/config-ui/login-settings-test.js
module('Acceptance | Enterprise | auth form custom login settings', function (hooks) {
  setupApplicationTest(hooks);
  hooks.beforeEach(async function () {
    await login();
    await runCmd([
      `write sys/namespaces/test-ns -force`,
      `write test-ns/sys/namespaces/child -force`,
      `write sys/config/ui/login/default-auth/root-rule backup_auth_types=token default_auth_type=okta disable_inheritance=false namespace_path=""`,
      `write sys/config/ui/login/default-auth/ns-rule default_auth_type=ldap disable_inheritance=true namespace_path=test-ns`,
      `write sys/auth/my_oidc type=oidc`,
      `write sys/auth/my_oidc/tune listing_visibility="unauth"`,
    ]);
    return await logout();
  });

  hooks.afterEach(async function () {
    // cleanup login rules
    await visit('/vault/auth?with=token');
    await fillIn(GENERAL.inputByAttr('token'), rootToken);
    await click(GENERAL.submitButton);
    await runCmd([
      'delete sys/config/ui/login/default-auth/root-rule',
      'delete sys/config/ui/login/default-auth/ns-rule',
      'delete sys/auth/my_oidc',
      'delete test-ns/sys/namespaces/child -f',
      'delete sys/namespaces/test-ns -f',
    ]);
  });

  test('it renders login settings for root namespace', async function (assert) {
    await visit('/vault/auth');
    await waitFor(AUTH_FORM.tabBtn('okta'));
    assert.dom(AUTH_FORM.tabBtn('okta')).hasAttribute('aria-selected', 'true');
    assert.dom(AUTH_FORM.authForm('okta')).exists('it renders default method');
    assert.dom(AUTH_FORM.advancedSettings).exists();

    await click(GENERAL.button('Sign in with other methods'));
    assert.dom(AUTH_FORM.authForm('token')).exists('it renders backup method');
  });

  test('it renders login settings for namespaces', async function (assert) {
    await visit('/vault/auth');
    await fillIn(GENERAL.inputByAttr('namespace'), 'test-ns');
    await waitFor(AUTH_FORM.authForm('ldap'));
    assert.dom(AUTH_FORM.authForm('ldap')).exists('it renders default method');
    assert.dom(AUTH_FORM.advancedSettings).exists();
    assert
      .dom(GENERAL.button('Sign in with other methods'))
      .doesNotExist('it does not render alternate view');

    // type in so that the namespace is "test-ns/child"
    await typeIn(GENERAL.inputByAttr('namespace'), '/child');
    await waitFor(AUTH_FORM.authForm('okta'));
    assert
      .dom(AUTH_FORM.authForm('okta'))
      .exists('it inherits view from root namespace because "test-ns" settings are not inheritable');
  });

  test('it ignores login settings if query param references a visible mount path', async function (assert) {
    await visit('/vault/auth?with=my_oidc%2F');
    await waitFor(AUTH_FORM.tabBtn('oidc'));
    assert
      .dom(AUTH_FORM.tabBtn('oidc'))
      .hasAttribute('aria-selected', 'true', 'it selects tab matching query param');
    assert.dom(AUTH_FORM.authForm('oidc')).exists();
    assert.dom(AUTH_FORM.advancedSettings).doesNotExist();
    await click(GENERAL.button('Sign in with other methods'));
    assert.dom(GENERAL.selectByAttr('auth type')).exists('dropdown renders as fallback view');
  });

  test('it ignores login settings if query param references a valid type', async function (assert) {
    await visit('/vault/auth?with=userpass');
    assert.dom(GENERAL.selectByAttr('auth type')).hasValue('userpass', 'dropdown selects userpass');
    await click(GENERAL.backButton);
    assert.dom(AUTH_FORM.tabBtn('oidc')).exists('it renders tabs on "Back" because visible mounts exist');
  });
});
