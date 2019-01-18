import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import backendsPage from 'vault/tests/pages/secrets/backends';
import authPage from 'vault/tests/pages/auth';
import withFlash from 'vault/tests/helpers/with-flash';

module('Acceptance | alicloud/enable', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('enable alicloud', async function(assert) {
    let enginePath = `alicloud-${new Date().getTime()}`;
    await mountSecrets.visit();
    await mountSecrets.selectType('alicloud');
    await withFlash(
      mountSecrets
        .next()
        .path(enginePath)
        .submit()
    );

    assert.equal(currentRouteName(), 'vault.cluster.secrets.backends', 'redirects to the backends page');

    assert.ok(backendsPage.rows.filterBy('path', `${enginePath}/`)[0], 'shows the alicloud engine');
  });
});
