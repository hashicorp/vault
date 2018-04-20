import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import authForm from '../pages/components/auth-form';
import { create } from 'ember-cli-page-object';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

const component = create(authForm);

moduleForAcceptance('Acceptance | auth', {
  beforeEach() {
    this.server = apiStub({ usePassthrough: true });
    return authLogout();
  },
  afterEach() {
    this.server.shutdown();
  },
});

test('auth query params', function(assert) {
  const backends = supportedAuthBackends();
  visit('/vault/auth');
  andThen(() => {
    assert.equal(currentURL(), '/vault/auth');
  });
  backends.reverse().forEach(backend => {
    click(`[data-test-auth-method-link="${backend.type}"]`);
    andThen(() => {
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
  andThen(() => {
    assert.equal(currentURL(), '/vault/auth');
  });
  component.token('token').tabs.filterBy('name', 'GitHub')[0].link();
  component.tabs.filterBy('name', 'Token')[0].link();
  andThen(() => {
    assert.equal(component.tokenValue, '', 'it clears the token value when toggling methods');
  });
});

test('it sends the right attributes when authenticating', function(assert) {
  let backends = supportedAuthBackends();
  visit('/vault/auth');
  backends.reverse().forEach(backend => {
    click(`[data-test-auth-method-link="${backend.type}"]`);
    if (backend.type === 'GitHub') {
      component.token('token');
    }
    component.login();
    andThen(() => {
      let lastRequest = this.server.passthroughRequests[this.server.passthroughRequests.length - 1];
      let body = JSON.parse(lastRequest.requestBody);
      if (backend.type === 'token') {
        assert.ok(
          Object.keys(lastRequest.requestHeaders).includes('X-Vault-Token'),
          'token uses vault token header'
        );
      } else if (backend.type === 'GitHub') {
        assert.ok(Object.keys(body).includes('token'), 'GitHub includes token');
      } else {
        assert.ok(Object.keys(body).includes('password'), `${backend.type} includes password`);
      }
    });
  });
});
