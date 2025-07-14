/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, fillIn, typeIn, visit, waitFor } from '@ember/test-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';
import customLoginHandler from 'vault/mirage/handlers/custom-login';
import customLoginScenario from 'vault/mirage/scenarios/custom-login';

// Auth form login settings
// This feature has thorough integration test coverage so only testing a few scenarios and direct link functionality
// Tests for read/list views are in ui/tests/acceptance/config-ui/login-settings-test.js
module('Acceptance | Enterprise | auth form custom login settings', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    customLoginHandler(this.server);
    customLoginScenario(this.server);
    // mirage scenario sets:
    // root namespace with 'okta' as default and 'token' as backup
    // 'test-ns' with 'ldap' as default and no backups
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
    assert.dom(AUTH_FORM.authForm('okta')).exists('it updates to render child namespace settings');
  });

  module('listing visibility', function (hooks) {
    hooks.beforeEach(function () {
      this.server.get('/sys/internal/ui/mounts', () => {
        // Stub a visible mount that does NOT match a type in the login settings
        return { data: { auth: { 'my_oidc/': { description: '', options: {}, type: 'oidc' } } } };
      });
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
});
