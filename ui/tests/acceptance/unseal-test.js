import { click, currentRouteName, currentURL, fillIn, settled, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';

const { unsealKeys } = VAULT_KEYS;

module('Acceptance | unseal', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  hooks.afterEach(function () {
    return logout.visit();
  });

  test('seal then unseal', async function (assert) {
    await visit('/vault/settings/seal');

    assert.equal(currentURL(), '/vault/settings/seal');

    // seal
    await click('[data-test-seal] button');

    await click('[data-test-confirm-button]');

    await pollCluster(this.owner);
    await settled();
    assert.equal(currentURL(), '/vault/unseal', 'vault is on the unseal page');

    // unseal
    for (const key of unsealKeys) {
      await fillIn('[data-test-shamir-input]', key);

      await click('button[type="submit"]');

      await pollCluster(this.owner);
      await settled();
    }

    assert.dom('[data-test-cluster-status]').doesNotExist('ui does not show sealed warning');
    assert.equal(currentRouteName(), 'vault.cluster.auth', 'vault is ready to authenticate');
  });
});
