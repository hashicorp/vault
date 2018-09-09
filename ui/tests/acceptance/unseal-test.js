import { click, fillIn, findAll, currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';

const { unseal } = VAULT_KEYS;

module('Acceptance | unseal', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  hooks.afterEach(function() {
    return authLogout();
  });

  test('seal then unseal', async function(assert) {
    await visit('/vault/settings/seal');
    assert.equal(currentURL(), '/vault/settings/seal');

    // seal
    await click('[data-test-seal] button');
    await click('[data-test-confirm-button]');
    pollCluster();
    assert.equal(currentURL(), '/vault/unseal', 'vault is on the unseal page');

    // unseal
    await fillIn('[data-test-shamir-input]', unseal);
    await click('button[type="submit"]');
    pollCluster();
    assert.equal(findAll('[data-test-cluster-status]').length, 0, 'ui does not show sealed warning');
    assert.ok(currentURL().match(/\/vault\/auth/), 'vault is ready to authenticate');
  });
});
