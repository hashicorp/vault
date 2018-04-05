import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/index';

moduleForAcceptance('Acceptance | settings/configure/secrets/pki', {
  beforeEach() {
    return authLogin();
  },
});

test('it redirects to the cert section', function(assert) {
  const path = `pki-${new Date().getTime()}`;
  mountSupportedSecretBackend(assert, 'pki', path);
  page.visit({ backend: path });
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.settings.configure-secret-backend.section',
      'redirects from the index'
    );
  });
});
