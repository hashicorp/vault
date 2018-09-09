import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/index';

module('Acceptance | settings/configure/secrets/pki', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it redirects to the cert section', function(assert) {
    const path = `pki-${new Date().getTime()}`;
    mountSupportedSecretBackend(assert, 'pki', path);
    page.visit({ backend: path });
    assert.equal(
      currentRouteName(),
      'vault.cluster.settings.configure-secret-backend.section',
      'redirects from the index'
    );
  });
});
