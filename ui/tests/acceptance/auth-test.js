import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';

moduleForAcceptance('Acceptance | auth', {
  afterEach() {
    return authLogout();
  },
});

test('auth query params', function(assert) {
  const backends = supportedAuthBackends();
  visit('/vault/auth');
  andThen(function() {
    assert.equal(currentURL(), '/vault/auth');
  });
  backends.reverse().forEach(backend => {
    click(`[data-test-auth-method-link="${backend.type}"]`);
    andThen(function() {
      assert.equal(
        currentURL(),
        `/vault/auth?with=${backend.type}`,
        `has the correct URL for ${backend.type}`
      );
    });
  });
});
