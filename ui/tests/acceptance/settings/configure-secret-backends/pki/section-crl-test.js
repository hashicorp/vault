import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';

moduleForAcceptance('Acceptance | settings/configure/secrets/pki/crl', {
  beforeEach() {
    return authLogin();
  },
});

test('it saves crl config', function(assert) {
  const path = `pki-${new Date().getTime()}`;
  mountSupportedSecretBackend(assert, 'pki', path);
  page.visit({ backend: path, section: 'crl' });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');
  });

  page.form.fillInField('time', 3);
  page.form.fillInField('unit', 'h');
  page.form.submit();

  andThen(() => {
    assert.equal(page.lastMessage, 'The crl config for this backend has been updated.');
  });
});
