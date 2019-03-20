import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import withFlash from 'vault/tests/helpers/with-flash';

module('Acceptance | settings/configure/secrets/pki/urls', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it saves urls config', async function(assert) {
    const path = `pki-${new Date().getTime()}`;
    await enablePage.enable('pki', path);
    await page.visit({ backend: path, section: 'urls' });
    assert.equal(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');

    await page.form.fields
      .objectAt(0)
      .input('foo')
      .change();
    await page.form.submit();

    assert.ok(page.form.hasError, 'shows error on invalid input');

    await page.form.fields
      .objectAt(0)
      .input('foo.example.com')
      .change();
    await withFlash(page.form.submit(), () => {
      assert.equal(page.lastMessage, 'The urls config for this backend has been updated.');
    });
  });
});
