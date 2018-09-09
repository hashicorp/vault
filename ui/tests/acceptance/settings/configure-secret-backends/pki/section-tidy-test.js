import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';

module('Acceptance | settings/configure/secrets/pki/tidy', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it saves tidy config', function(assert) {
    const path = `pki-${new Date().getTime()}`;
    mountSupportedSecretBackend(assert, 'pki', path);
    page.visit({ backend: path, section: 'tidy' });
    assert.equal(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');
    page.form.fields();

    page.form.fields(0).clickLabel();
    page.form.submit();

    assert.equal(page.lastMessage, 'The tidy config for this backend has been updated.');
  });
});
