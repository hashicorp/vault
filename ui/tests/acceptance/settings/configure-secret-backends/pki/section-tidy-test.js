import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';

moduleForAcceptance('Acceptance | settings/configure/secrets/pki/tidy', {
  beforeEach() {
    return authLogin();
  },
});

test('it saves tidy config', function(assert) {
  const path = `pki-${new Date().getTime()}`;
  mountSupportedSecretBackend(assert, 'pki', path);
  page.visit({ backend: path, section: 'tidy' });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');
    page.form.fields();
  });

  page.form.fields(0).clickLabel();
  page.form.submit();

  andThen(() => {
    assert.equal(page.lastMessage, 'The tidy config for this backend has been updated.');
  });
});
