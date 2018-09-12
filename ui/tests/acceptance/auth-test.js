import { click, currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import authForm from '../pages/components/auth-form';
import { create } from 'ember-cli-page-object';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import logout from 'vault/tests/pages/logout';

const component = create(authForm);

module('Acceptance | auth', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub({ usePassthrough: true });
    return logout.visit();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('auth query params', async function(assert) {
    let backends = supportedAuthBackends();
    await visit('/vault/auth');
    assert.equal(currentURL(), '/vault/auth?with=token');
    for (let backend of backends.reverse()) {
      await click(`[data-test-auth-method-link="${backend.type}"]`);
      assert.equal(
        currentURL(),
        `/vault/auth?with=${backend.type}`,
        `has the correct URL for ${backend.type}`
      );
    }
  });

  test('it clears token when changing selected auth method', async function(assert) {
    await visit('/vault/auth');
    assert.equal(currentURL(), '/vault/auth?with=token');
    await component
      .token('token')
      .tabs.filterBy('name', 'GitHub')[0]
      .link();
    await component.tabs.filterBy('name', 'Token')[0].link();
    assert.equal(component.tokenValue, '', 'it clears the token value when toggling methods');
  });

  test('it sends the right attributes when authenticating', async function(assert) {
    let backends = supportedAuthBackends();
    await visit('/vault/auth');
    for (let backend of backends.reverse()) {
      await click(`[data-test-auth-method-link="${backend.type}"]`);
      if (backend.type === 'github') {
        await component.token('token');
      }
      await component.login();
      let lastRequest = this.server.passthroughRequests[this.server.passthroughRequests.length - 1];
      let body = JSON.parse(lastRequest.requestBody);
      if (backend.type === 'token') {
        assert.ok(
          Object.keys(lastRequest.requestHeaders).includes('X-Vault-Token'),
          'token uses vault token header'
        );
      } else if (backend.type === 'github') {
        assert.ok(Object.keys(body).includes('token'), 'GitHub includes token');
      } else {
        assert.ok(Object.keys(body).includes('password'), `${backend.type} includes password`);
      }
    }
  });
});
