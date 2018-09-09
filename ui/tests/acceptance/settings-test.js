import { currentURL, find, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import backendListPage from 'vault/tests/pages/secrets/backends';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | settings', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  hooks.afterEach(function() {
    return authLogout();
  });

  test('settings', async function(assert) {
    const now = new Date().getTime();
    const type = 'consul';
    const path = `path-${now}`;

    // mount unsupported backend
    await visit('/vault/settings/mount-secret-backend');
    assert.equal(currentURL(), '/vault/settings/mount-secret-backend');

    mountSecrets
      .selectType(type)
      .next()
      .path(path)
      .toggleOptions()
      .defaultTTLVal(100)
      .defaultTTLUnit('s')
      .submit();
    assert.equal(currentURL(), `/vault/secrets`, 'redirects to secrets page');
    assert.ok(
      find('[data-test-flash-message]').textContent.trim(),
      `Successfully mounted '${type}' at '${path}'!`
    );
    let row = backendListPage.rows().findByPath(path);
    row.menu();
    backendListPage.configLink();
    assert.ok(currentURL(), '/vault/secrets/${path}/configuration', 'navigates to the config page');
  });
});
