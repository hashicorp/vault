import { currentURL, currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/index';

module('Acceptance | settings/auth/configure', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it redirects to section options when there are no other sections', function(assert) {
    const path = `approle-${new Date().getTime()}`;
    const type = 'approle';
    enablePage.visit();
    enablePage.form.mount(type, path);
    page.visit({ path });
    assert.equal(currentRouteName(), 'vault.cluster.settings.auth.configure.section');
    assert.equal(currentURL(), `/vault/settings/auth/configure/${path}/options`, 'loads the options route');
  });

  test('it redirects to the first section', function(assert) {
    const path = `aws-${new Date().getTime()}`;
    const type = 'aws';
    enablePage.visit();
    enablePage.form.mount(type, path);
    page.visit({ path });
    assert.equal(currentRouteName(), 'vault.cluster.settings.auth.configure.section');
    assert.equal(
      currentURL(),
      `/vault/settings/auth/configure/${path}/client`,
      'loads the first section for the type of auth method'
    );
  });
});
