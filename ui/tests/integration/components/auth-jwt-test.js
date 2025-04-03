/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { _cancelTimers as cancelTimers } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { fillIn, render, settled, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { resolve } from 'rsvp';
import { create } from 'ember-cli-page-object';
import form from '../../pages/components/auth-jwt';
import { ERROR_JWT_LOGIN } from 'vault/components/auth-jwt';
import { callbackData } from 'vault/tests/helpers/oidc-window-stub';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

const component = create(form);

const renderIt = async (context, { path = 'jwt', type = 'jwt' } = {}) => {
  const handler = (data, e) => {
    if (e && e.preventDefault) e.preventDefault();
    return resolve();
  };

  context.error = '';
  context.handler = sinon.spy(handler);
  context.roleName = '';
  context.selectedAuthPath = path;
  context.selectedAuthType = type;
  await render(hbs`
    <AuthJwt
      @roleName={{this.roleName}}
      @selectedAuthPath={{this.selectedAuthPath}}
      @selectedAuthType={{this.selectedAuthType}}
      @onError={{fn (mut this.error)}}
      @onNamespace={{fn (mut this.namespace)}}
      @onSelectedAuth={{fn (mut this.selectedAuth)}}
      @onSubmit={{this.handler}}
      @onRoleName={{fn (mut this.roleName)}}
    />
    `);
};
module('Integration | Component | auth jwt', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.windowStub = sinon.stub(window, 'open');

    this.owner.lookup('service:router').reopen({
      urlFor() {
        return 'http://example.com';
      },
    });
    this.server.get('/auth/:path/oidc/callback', function () {
      return { auth: { client_token: 'token' } };
    });
    this.server.post('/auth/:path/oidc/auth_url', (_, request) => {
      const { role } = JSON.parse(request.requestBody);
      if (['okta', 'test', 'bar'].includes(role)) {
        const auth_url = role === 'test' ? 'http://example.com' : role === 'okta' ? 'http://okta.com' : '';
        return {
          data: { auth_url },
        };
      }
      const errors = role === 'foo' ? ['role "foo" could not be found'] : [ERROR_JWT_LOGIN];
      return overrideResponse(400, { errors });
    });
  });

  hooks.afterEach(function () {
    this.windowStub.restore();
  });

  test('it renders the yield', async function (assert) {
    await render(hbs`<AuthJwt @onSubmit={{action (mut this.submit)}}>Hello!</AuthJwt>`);
    assert.strictEqual(component.yieldContent, 'Hello!', 'yields properly');
  });

  test('it fetches auth_url when type changes', async function (assert) {
    assert.expect(2);
    await renderIt(this, { path: '', type: 'jwt' });
    // auth_url is requested on initial render so stubbing after rendering the component
    // to test auth_url is called when the type changes
    this.server.post('/auth/:path/oidc/auth_url', (_, request) => {
      assert.true(true, 'request is made to auth_url');
      const { path } = request.params;
      assert.strictEqual(path, 'oidc', `path param is updated type: ${path}`);
      return {
        data: { auth_url: '' },
      };
    });
    this.set('selectedAuthType', 'oidc');
    await settled();
  });

  test('if auth path exists it uses it to build url request instead of type', async function (assert) {
    this.server.post('/auth/:path/oidc/auth_url', (_, request) => {
      const { path } = request.params;
      assert.strictEqual(path, 'custom-jwt', `path param is custom path: ${path}`);
      return {};
    });
    await renderIt(this, { path: 'custom-jwt' });
  });

  test('jwt: it renders and makes auth_url requests', async function (assert) {
    let postCount = 0;
    this.server.post('/auth/:path/oidc/auth_url', (_, request) => {
      postCount++;
      const { path } = request.params;
      const expectedUrl = `/v1/auth/${path}/oidc/auth_url`;
      assert.strictEqual(request.url, expectedUrl);
      return overrideResponse(400, { errors: [ERROR_JWT_LOGIN] });
    });
    await renderIt(this);
    await settled();
    assert.strictEqual(postCount, 1, 'request to the default path is made');
    assert.ok(component.jwtPresent, 'renders jwt field');
    assert.ok(component.rolePresent, 'renders jwt field');

    this.set('selectedAuthPath', 'foo');
    await settled();
    assert.strictEqual(postCount, 2, 'a second request was made');
  });

  test('jwt: it calls passed action on login', async function (assert) {
    await renderIt(this);
    await component.login();
    assert.ok(this.handler.calledOnce);
  });

  test('oidc: test role: it renders', async function (assert) {
    let postCount = 0;
    this.server.post('/auth/:path/oidc/auth_url', (_, request) => {
      postCount++;
      const { role } = JSON.parse(request.requestBody);
      const auth_url = role === 'test' ? 'http://example.com' : role === 'okta' ? 'http://okta.com' : '';
      return {
        data: { auth_url },
      };
    });
    await renderIt(this, { path: 'foo', type: 'oidc' });
    await settled();
    await fillIn(AUTH_FORM.roleInput, 'test');
    assert
      .dom(AUTH_FORM.input('jwt'))
      .doesNotExist('does not show jwt token input if role matches OIDC login url');
    assert.dom(AUTH_FORM.login).hasText('Sign in with OIDC Provider');
    await fillIn(AUTH_FORM.roleInput, 'okta');
    // 1 for initial render, 1 for each time role changed = 3
    assert.strictEqual(postCount, 3, 'fetches the auth_url when the role changes');
    assert.dom(AUTH_FORM.login).hasText('Sign in with Okta', 'recognizes auth methods with certain urls');
  });

  test('oidc: it fetches auth_url when path changes', async function (assert) {
    assert.expect(2);
    await renderIt(this, { path: 'oidc', type: 'oidc' });
    // auth_url is requested on initial render so stubbing after rendering the component
    // to test auth_url is called when the :path changes
    this.server.post('/auth/:path/oidc/auth_url', (_, request) => {
      assert.true(true, 'request is made to auth_url');
      assert.strictEqual(request?.params?.path, 'foo', 'request params are { path: foo }');
      return {
        data: { auth_url: '' },
      };
    });

    this.set('selectedAuthPath', 'foo');
    await settled();
  });

  test('oidc: it calls window.open popup window on login', async function (assert) {
    sinon.replaceGetter(window, 'screen', () => ({ height: 600, width: 500 }));
    await renderIt(this, { path: 'foo', type: 'oidc' });
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.windowStub.calledOnce;
    });

    const call = this.windowStub.lastCall;
    assert.deepEqual(
      call.args,
      ['http://example.com', 'vaultOIDCWindow', 'width=500,height=600,resizable,scrollbars=yes,top=0,left=0'],
      'called with expected args'
    );
    sinon.restore();
  });

  // not the greatest test because this test would also pass if the origin matched
  // because event.isTrusted is always false (another condition checked by the component)
  test('oidc: fails silently when event origin does not match window origin', async function (assert) {
    assert.expect(3);
    // prevent test incorrectly passing because the event isn't triggered at all
    // by also asserting that the message event fires
    const message = { data: callbackData(), origin: 'http://hackerz.com' };
    const assertEvent = (event) => {
      assert.propEqual(event.data, message.data, 'message has expected data');
      assert.strictEqual(event.origin, message.origin, 'message has expected origin');
    };
    window.addEventListener('message', assertEvent);

    await renderIt(this, { path: 'foo', type: 'oidc' });
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.windowStub.calledOnce;
    });

    window.dispatchEvent(new MessageEvent('message', message));
    cancelTimers();
    await settled();
    assert.false(this.handler.called, 'should not call the submit handler');

    // Cleanup
    window.removeEventListener('message', assertEvent);
  });

  test('oidc: fails silently when event is not trusted', async function (assert) {
    assert.expect(2);
    // prevent test incorrectly passing because the event isn't triggered at all
    // by also asserting that the message event fires
    const messageData = callbackData();
    const assertEvent = (event) => {
      assert.propEqual(event.data, messageData, 'message event fires');
    };
    window.addEventListener('message', assertEvent);

    await renderIt(this, { path: 'foo', type: 'oidc' });
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.windowStub.calledOnce;
    });
    // mocking a message event is always untrusted (there is no way to override isTrusted on the window object)
    window.dispatchEvent(new MessageEvent('message', { data: messageData }));

    cancelTimers();
    await settled();
    assert.false(this.handler.called, 'should not call the submit handler');

    // Cleanup
    window.removeEventListener('message', assertEvent);
  });

  test('oidc: it should trigger error callback when role is not found', async function (assert) {
    await renderIt(this, { path: 'oidc', type: 'oidc' });
    await component.role('foo');
    await component.login();
    assert.strictEqual(
      this.error,
      'Invalid role. Please try again.',
      'Error message is returned when role is not found'
    );
  });

  test('oidc: it should trigger error callback when role is returned without auth_url', async function (assert) {
    await renderIt(this, { path: 'oidc', type: 'oidc' });
    await component.role('bar');
    await component.login();
    assert.strictEqual(
      this.error,
      'Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.',
      'Error message is returned when role is returned without auth_url'
    );
  });
});
