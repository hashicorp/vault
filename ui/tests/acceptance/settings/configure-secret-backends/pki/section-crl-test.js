import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';

module('Acceptance | settings/configure/secrets/pki/crl', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it saves crl config', function(assert) {
    const path = `pki-${new Date().getTime()}`;
    mountSupportedSecretBackend(assert, 'pki', path);
    page.visit({ backend: path, section: 'crl' });
    assert.equal(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');

    page.form.fillInField('time', 3);
    page.form.fillInField('unit', 'h');
    page.form.submit();

    assert.equal(page.lastMessage, 'The crl config for this backend has been updated.');
  });
});
