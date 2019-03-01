import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/auth/enable';
import listPage from 'vault/tests/pages/access/methods';
import authPage from 'vault/tests/pages/auth';
import withFlash from 'vault/tests/helpers/with-flash';

module('Acceptance | settings/auth/enable', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it mounts and redirects', async function(assert) {
    // always force the new mount to the top of the list
    const path = `approle-${new Date().getTime()}`;
    const type = 'approle';
    await page.visit();
    assert.equal(currentRouteName(), 'vault.cluster.settings.auth.enable');
    await withFlash(page.enable(type, path), () => {
      assert.equal(
        page.flash.latestMessage,
        `Successfully mounted the ${type} auth method at ${path}.`,
        'success flash shows'
      );
    });
    assert.equal(
      currentRouteName(),
      'vault.cluster.settings.auth.configure.section',
      'redirects to the auth config page'
    );

    await listPage.visit();
    assert.ok(listPage.findLinkById(path), 'mount is present in the list');
  });
});
