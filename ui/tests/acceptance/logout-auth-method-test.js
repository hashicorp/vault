import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, visit, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { fakeWindow, buildMessage } from '../helpers/oidc-window-stub';
import sinon from 'sinon';
import { later, _cancelTimers as cancelTimers } from '@ember/runloop';

module('Acceptance | logout auth method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.openStub = sinon.stub(window, 'open').callsFake(() => fakeWindow.create());
  });
  hooks.afterEach(function () {
    this.openStub.restore();
  });

  // coverage for bug where token was selected as auth method for oidc and jwt
  test('it should populate oidc auth method on logout', async function (assert) {
    this.server.post('/auth/oidc/oidc/auth_url', () => ({
      data: { auth_url: 'http://example.com' },
    }));
    this.server.get('/auth/foo/oidc/callback', () => ({
      auth: { client_token: 'root' },
    }));
    // ensure clean state
    sessionStorage.removeItem('selectedAuth');
    await visit('/vault/auth');
    await fillIn('[data-test-select="auth-method"]', 'oidc');
    later(() => {
      window.postMessage(buildMessage().data, window.origin);
      cancelTimers();
    }, 50);
    await click('[data-test-auth-submit]');
    await click('.nav-user-button button');
    await click('#logout');
    assert
      .dom('[data-test-select="auth-method"]')
      .hasValue('oidc', 'Previous auth method selected on logout');
  });
});
