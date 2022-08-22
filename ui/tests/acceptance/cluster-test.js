import { create } from 'ember-cli-page-object';
import { settled, click } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
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

const USER = 'end-user';
const PASSWORD = 'mypassword';

const authAccessor = async function (path) {
  await enablePage.enable('userpass', path);
  await consoleComponent.runCommands([`write auth/${path}/users/end-user password="${PASSWORD}"`]);
};

const setupUser = async function () {
  let authMethodPath = `userpass-${new Date().getTime()}`;
  await authAccessor(authMethodPath);
};

module('Acceptance | cluster', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await logout.visit();
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

    assert.dom('[data-test-navbar-item="policies"]').doesNotExist();
    await logout.visit();
  });

  test('it hides mfa setup if user has not entityId (ex: is a root user)', async function (assert) {
    await setupUser();

    await logout.visit();
    await settled();
    await authPage.loginUsername(USER, PASSWORD);
    await click('.nav-user-button button');
    assert.dom('[data-test-status-link="mfa"]').exists();
    await logout.visit();

    await authPage.login('root');
    await settled();
    await click('.nav-user-button button');
    assert.dom('[data-test-status-link="mfa"]').doesNotExist();
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
    await settled();
    assert.dom('[data-test-navbar-item="policies"]').hasAttribute('href', '/ui/vault/policies/rgp');
    await logout.visit();
  });
});
