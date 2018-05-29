import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/settings/mount-secret-backend';
import configPage from 'vault/tests/pages/secrets/backend/configuration';

moduleForAcceptance('Acceptance | settings/mount-secret-backend', {
  beforeEach() {
    return authLogin();
  },
});

test('it sets the ttl corrects when mounting', function(assert) {
  page.visit();
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
  });
  // always force the new mount to the top of the list
  const path = `kv-${new Date().getTime()}`;
  const defaultTTLHours = 100;
  const maxTTLHours = 300;
  const defaultTTLSeconds = defaultTTLHours * 60 * 60;
  const maxTTLSeconds = maxTTLHours * 60 * 60;
  page
    .type('kv')
    .path(path)
    .toggleOptions()
    .defaultTTLVal(defaultTTLHours)
    .defaultTTLUnit('h')
    .maxTTLVal(maxTTLHours)
    .maxTTLUnit('h')
    .submit();

  configPage.visit({ backend: path });
  andThen(() => {
    assert.equal(configPage.defaultTTL, defaultTTLSeconds, 'shows the proper TTL');
    assert.equal(configPage.maxTTL, maxTTLSeconds, 'shows the proper max TTL');
  });
});
