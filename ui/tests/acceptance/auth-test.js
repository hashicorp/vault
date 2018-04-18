import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import authForm from '../pages/components/auth-form';
import { create } from 'ember-cli-page-object';

const component = create(authForm);

moduleForAcceptance('Acceptance | auth', {
  beforeEach() {
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

test('it clears token when changing selected auth method', function(assert) {
  visit('/vault/auth');
  andThen(function() {
    assert.equal(currentURL(), '/vault/auth');
  });
  component.token('token').tabs.filterBy('name', 'GitHub')[0].link();
  component.tabs.filterBy('name', 'Token')[0].link();
  andThen(function() {
    assert.equal(component.tokenValue, '', 'it clears the token value when toggling methods');
  });
});
