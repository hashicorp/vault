/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { settled, click, visit, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { login, loginMethod, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const tokenWithPolicy = async function (name, policy) {
  return await runCmd([
    `write sys/policies/acl/${name} policy=${btoa(policy)}`,
    `write -field=client_token auth/token/create policies=${name}`,
  ]);
};

module('Acceptance | cluster', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    return login();
  });

  test('hides nav item if user does not have permission', async function (assert) {
    const deny_policies_policy = `
      path "sys/policies/*" {
        capabilities = ["deny"]
      },
    `;

    const userToken = await tokenWithPolicy('hide-policies-nav', deny_policies_policy);
    await login(userToken);
    await visit('/vault/access');

    assert.dom('[data-test-sidebar-nav-link="ACL policies"]').doesNotExist();
    assert.dom('[data-test-sidebar-nav-link="Role governing policies"]').doesNotExist();
    assert.dom('[data-test-sidebar-nav-link="Endpoint governing policies"]').doesNotExist();
  });

  test('it hides mfa setup if user does not have entityId (ex: is a root user)', async function (assert) {
    const username = 'end-user';
    const password = 'mypassword';
    const path = `cluster-userpass-${uuidv4()}`;

    await visit('/vault/settings/auth/enable');
    await mountBackend('userpass', path);
    await runCmd([`write auth/${path}/users/end-user password="${password}"`]);

    await loginMethod({ username, password, path }, { authType: 'userpass', toggleOptions: true });

    await click(GENERAL.button('user-menu-trigger'));
    assert.dom('[data-test-user-menu-item="mfa"]').exists();

    await login();
    await settled();
    await click(GENERAL.button('user-menu-trigger'));
    assert.dom('[data-test-user-menu-item="mfa"]').doesNotExist();
  });

  test('enterprise nav item links to first route that user has access to', async function (assert) {
    const read_rgp_policy = `
      path "sys/policies/rgp" {
        capabilities = ["read"]
      },
    `;

    const userToken = await tokenWithPolicy('show-policies-nav', read_rgp_policy);
    await login(userToken);
    await visit('/vault/access');

    assert.dom('[data-test-sidebar-nav-link="ACL policies"]').doesNotExist();
    assert.dom('[data-test-sidebar-nav-link="Role governing policies"]').exists();
    assert.dom('[data-test-sidebar-nav-link="Endpoint governing policies"]').doesNotExist();
  });

  test('shows error banner if resultant-acl check fails', async function (assert) {
    const version = this.owner.lookup('service:version');
    const login_only = `
      path "auth/token/lookup-self" {
        capabilities = ["read"]
      },
    `;
    // note: the default policy is attached to a user unless you add the no_default_policy=true flag
    // you can confirm this by running `vault token lookup` on the generated token
    const noDefaultPolicyUser = await runCmd([
      `write sys/policies/acl/login-only policy=${btoa(login_only)}`,
      `write -field=client_token auth/token/create no_default_policy=true policies="login-only"`,
    ]);

    assert
      .dom('[data-test-resultant-acl-banner]')
      .doesNotExist('Resultant ACL banner does not show as root user with access to everything');

    await logout();
    assert
      .dom('[data-test-resultant-acl-banner]')
      .doesNotExist('Resultant ACL banner does not show on login page');

    await login(noDefaultPolicyUser);
    const expectedText = version.isEnterprise
      ? "Resultant ACL check failed Links might be shown that you don't have access to. Contact your administrator to update your policy. Log into root namespace"
      : "Resultant ACL check failed Links might be shown that you don't have access to. Contact your administrator to update your policy.";

    assert
      .dom('[data-test-resultant-acl-banner]')
      .includesText(expectedText, 'Resultant ACL banner shows appropriate message for OSS/Enterprise');
  });

  test('redirects to secret-engines from legacy /secrets path', async function (assert) {
    await visit('/vault/secrets');
    assert.strictEqual(
      currentURL(),
      '/vault/secrets-engines',
      'Navigating to /secrets redirects to /secrets-engines'
    );
  });
});
