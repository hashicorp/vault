import { currentURL, currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/index';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | settings/auth/configure', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.timestamp = new Date().getTime();
    return authPage.login();
  });

  test('it redirects to section options when there are no other sections', async function (assert) {
    const path = `approle-${this.timestamp}`;
    const type = 'approle';
    await enablePage.enable(type, path);
    await page.visit({ path });
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.auth.configure.section');
    assert.strictEqual(
      currentURL(),
      `/vault/settings/auth/configure/${path}/options`,
      'loads the options route'
    );
  });

  test('it redirects to the first section', async function (assert) {
    const path = `aws-${this.timestamp}`;
    const type = 'aws';
    await enablePage.enable(type, path);
    await page.visit({ path });
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.auth.configure.section');
    assert.strictEqual(
      currentURL(),
      `/vault/settings/auth/configure/${path}/client`,
      'loads the first section for the type of auth method'
    );
  });
});
