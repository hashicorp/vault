import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';

const { unseal } = VAULT_KEYS;

moduleForAcceptance('Acceptance | unseal', {
  beforeEach() {
    return authLogin();
  },
  afterEach() {
    return authLogout();
  },
});

test('seal then unseal', function(assert) {
  visit('/vault/settings/seal');
  andThen(function() {
    assert.equal(currentURL(), '/vault/settings/seal');
  });

  // seal
  click('[data-test-seal] button');
  click('[data-test-confirm-button]');
  andThen(() => {
    pollCluster();
  });
  andThen(function() {
    assert.equal(currentURL(), '/vault/unseal', 'vault is on the unseal page');
  });

  // unseal
  fillIn('[data-test-shamir-input]', unseal);
  click('button[type="submit"]');
  andThen(() => {
    pollCluster();
  });
  andThen(() => {
    assert.equal(find('[data-test-cluster-status]').length, 0, 'ui does not show sealed warning');
    assert.ok(currentURL().match(/\/vault\/auth/), 'vault is ready to authenticate');
  });
});
