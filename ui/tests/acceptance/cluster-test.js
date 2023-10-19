/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create } from 'ember-cli-page-object';
import { settled, click, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';

const consoleComponent = create(consoleClass);

const tokenWithPolicy = async function (name, policy) {
  await consoleComponent.runCommands([
    `write sys/policies/acl/${name} policy=${btoa(policy)}`,
    `write -field=client_token auth/token/create policies=${name}`,
  ]);

  return consoleComponent.lastLogOutput;
};

module('Acceptance | cluster', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    return authPage.login();
  });

  test('hides nav item if user does not have permission', async function (assert) {
    const deny_policies_policy = `
      path "sys/policies/*" {
        capabilities = ["deny"]
      },
    `;

    const userToken = await tokenWithPolicy('hide-policies-nav', deny_policies_policy);
    await logout.visit();
    await authPage.login(userToken);
    await visit('/vault/access');

    assert.dom('[data-test-sidebar-nav-link="Policies"]').doesNotExist();
  });

  test('it hides mfa setup if user has not entityId (ex: is a root user)', async function (assert) {
    const user = 'end-user';
    const password = 'mypassword';
    const path = `cluster-userpass-${uuidv4()}`;

    await enablePage.enable('userpass', path);
    await consoleComponent.runCommands([`write auth/${path}/users/end-user password="${password}"`]);

    await logout.visit();
    await settled();
    await authPage.loginUsername(user, password, path);
    await click('[data-test-user-menu-trigger]');
    assert.dom('[data-test-user-menu-item="mfa"]').exists();
    await logout.visit();

    await authPage.login('root');
    await settled();
    await click('[data-test-user-menu-trigger]');
    assert.dom('[data-test-user-menu-item="mfa"]').doesNotExist();
  });

  test('enterprise nav item links to first route that user has access to', async function (assert) {
    const read_rgp_policy = `
      path "sys/policies/rgp" {
        capabilities = ["read"]
      },
    `;

    const userToken = await tokenWithPolicy('show-policies-nav', read_rgp_policy);
    await logout.visit();
    await authPage.login(userToken);
    await visit('/vault/access');

    assert.dom('[data-test-sidebar-nav-link="Policies"]').hasAttribute('href', '/ui/vault/policies/rgp');
  });

  test('shows error banner if resultant-acl check fails', async function (assert) {
    const login_only = `
      path "auth/token/lookup-self" {
        capabilities = ["read"]
      },
    `;
    await consoleComponent.runCommands([
      `write sys/policies/acl/login-only policy=${btoa(login_only)}`,
      `write -field=client_token auth/token/create no_default_policy=true policies="login-only"`,
    ]);
    const noDefaultPolicyUser = consoleComponent.lastLogOutput;
    assert.dom('[data-test-resultant-acl-banner]').doesNotExist('Resultant ACL banner does not show as root');
    await logout.visit();
    assert.dom('[data-test-resultant-acl-banner]').doesNotExist('Does not show on login page');
    await authPage.login(noDefaultPolicyUser);
    assert.dom('[data-test-resultant-acl-banner]').includesText('Resultant ACL check failed');
  });
});
