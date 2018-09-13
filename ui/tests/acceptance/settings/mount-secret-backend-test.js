import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/mount-secret-backend';
import configPage from 'vault/tests/pages/secrets/backend/configuration';
import authPage from 'vault/tests/pages/auth';
import withFlash from 'vault/tests/helpers/with-flash';

module('Acceptance | settings/mount-secret-backend', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it sets the ttl corrects when mounting', async function(assert) {
    // always force the new mount to the top of the list
    const path = `kv-${new Date().getTime()}`;
    const defaultTTLHours = 100;
    const maxTTLHours = 300;
    const defaultTTLSeconds = defaultTTLHours * 60 * 60;
    const maxTTLSeconds = maxTTLHours * 60 * 60;

    await page.visit();
    assert.equal(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await page.selectType('kv');
    await withFlash(
      page
        .next()
        .path(path)
        .toggleOptions()
        .defaultTTLVal(defaultTTLHours)
        .defaultTTLUnit('h')
        .maxTTLVal(maxTTLHours)
        .maxTTLUnit('h')
        .submit()
    );
    await configPage.visit({ backend: path });
    assert.equal(configPage.defaultTTL, defaultTTLSeconds, 'shows the proper TTL');
    assert.equal(configPage.maxTTL, maxTTLSeconds, 'shows the proper max TTL');
  });
});
