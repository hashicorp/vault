import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';

module('Acceptance | settings/configure/secrets/pki/urls', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it saves urls config', function(assert) {
    const path = `pki-${new Date().getTime()}`;
    mountSupportedSecretBackend(assert, 'pki', path);
    page.visit({ backend: path, section: 'urls' });
    assert.equal(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');

    page.form
      .fields(0)
      .input('foo')
      .change();
    page.form.submit();

    assert.ok(page.form.hasError, 'shows error on invalid input');

    page.form
      .fields(0)
      .input('foo.example.com')
      .change();
    page.form.submit();

    assert.equal(page.lastMessage, 'The urls config for this backend has been updated.');
  });
});
