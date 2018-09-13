import { currentURL, find, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import backendListPage from 'vault/tests/pages/secrets/backends';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import withFlash from 'vault/tests/helpers/with-flash';

module('Acceptance | settings', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  hooks.afterEach(function() {
    return logout.visit();
  });

  test('settings', async function(assert) {
    const now = new Date().getTime();
    const type = 'consul';
    const path = `path-${now}`;

    // mount unsupported backend
    await visit('/vault/settings/mount-secret-backend');
    assert.equal(currentURL(), '/vault/settings/mount-secret-backend');

    await mountSecrets.selectType(type);
    await withFlash(
      mountSecrets
        .next()
        .path(path)
        .toggleOptions()
        .defaultTTLVal(100)
        .defaultTTLUnit('s')
        .submit(),
      () => {
        assert.ok(
          find('[data-test-flash-message]').textContent.trim(),
          `Successfully mounted '${type}' at '${path}'!`
        );
      }
    );
    assert.equal(currentURL(), `/vault/secrets`, 'redirects to secrets page');
    let row = backendListPage.rows.filterBy('path', path + '/')[0];
    await row.menu();
    await backendListPage.configLink();
    assert.ok(currentURL(), '/vault/secrets/${path}/configuration', 'navigates to the config page');
  });
});
