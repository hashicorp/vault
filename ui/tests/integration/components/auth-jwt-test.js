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
import { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS, ERROR_JWT_LOGIN } from 'vault/components/auth-jwt';
import { fakeWindow, buildMessage } from 'vault/tests/helpers/oidc-window-stub';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

const component = create(form);
const windows = [];

fakeWindow.reopen({
  init() {
    this._super(...arguments);
    windows.push(this);
  },
  open() {
    return fakeWindow.create();
  },
  close() {
    windows.forEach((w) => w.trigger('close'));
  },
});

const renderIt = async (context, path = 'jwt') => {
  const handler = (data, e) => {
    if (e && e.preventDefault) e.preventDefault();
    return resolve();
  };
  const fake = fakeWindow.create();
  context.set('window', fake);
  context.set('handler', sinon.spy(handler));
  context.set('roleName', '');
  context.set('selectedAuthPath', path);
  await render(hbs`
    <AuthJwt
      @window={{this.window}}
      @roleName={{this.roleName}}
      @selectedAuthPath={{this.selectedAuthPath}}
      @onError={{action (mut this.error)}}
      @onLoading={{action (mut this.isLoading)}}
      @onNamespace={{action (mut this.namespace)}}
      @onSelectedAuth={{action (mut this.selectedAuth)}}
      @onSubmit={{action this.handler}}
      @onRoleName={{action (mut this.roleName)}}
    />
    `);
};
module('Integration | Component | auth jwt', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.openSpy = sinon.spy(fakeWindow.proto(), 'open');
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
    this.openSpy.restore();
    this.server.shutdown();
  });

  test('it renders the yield', async function (assert) {
    await render(hbs`<AuthJwt @onSubmit={{action (mut this.submit)}}>Hello!</AuthJwt>`);
    assert.strictEqual(component.yieldContent, 'Hello!', 'yields properly');
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
    // setting the path also fires off a request to auth_url but this happens inconsistently in tests
    // setting here so it doesn't affect the postCount because it's not relevant to what's being tested
    this.set('selectedAuthPath', 'foo');
    let postCount = 0;
    this.server.post('/auth/:path/oidc/auth_url', (_, request) => {
      postCount++;
      const { role } = JSON.parse(request.requestBody);
      const auth_url = role === 'test' ? 'http://example.com' : role === 'okta' ? 'http://okta.com' : '';
      return {
        data: { auth_url },
      };
    });
    await renderIt(this);
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
    this.set('selectedAuthPath', 'foo');
    await renderIt(this);
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
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.openSpy.calledOnce;
    });

    cancelTimers();
    await settled();

    const call = this.openSpy.getCall(0);
    assert.deepEqual(
      call.args,
      ['http://example.com', 'vaultOIDCWindow', 'width=500,height=600,resizable,scrollbars=yes,top=0,left=0'],
      'called with expected args'
    );
  });

  test('oidc: it calls error handler when popup is closed', async function (assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.openSpy.calledOnce;
    });
    this.window.close();
    await settled();
    assert.strictEqual(this.error, ERROR_WINDOW_CLOSED, 'calls onError with error string');
  });

  test('oidc: shows error when message posted with state key, wrong params', async function (assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.openSpy.calledOnce;
    });
    this.window.trigger(
      'message',
      buildMessage({ data: { source: 'oidc-callback', state: 'state', foo: 'bar' } })
    );
    cancelTimers();
    await settled();

    assert.strictEqual(this.error, ERROR_MISSING_PARAMS, 'calls onError with params missing error');
  });

  test('oidc: storage event fires with state key, correct params', async function (assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.openSpy.calledOnce;
    });
    this.window.trigger('message', buildMessage());
    await settled();
    const [callbackData, , token] = this.handler.lastCall.args;
    assert.propEqual(
      callbackData,
      { mfa_requirement: undefined },
      'mfa_requirement is undefined if not returned by response'
    );
    assert.strictEqual(token, 'token', 'calls the onSubmit handler with token');
  });

  test('oidc: fails silently when event origin does not match window origin', async function (assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.openSpy.calledOnce;
    });
    this.window.trigger('message', buildMessage({ origin: 'http://hackerz.com' }));

    cancelTimers();
    await settled();

    assert.false(this.handler.called, 'should not call the submit handler');
  });

  test('oidc: fails silently when event is not trusted', async function (assert) {
    await renderIt(this);
    this.set('selectedAuthPath', 'foo');
    await component.role('test');
    component.login();
    await waitUntil(() => {
      return this.openSpy.calledOnce;
    });
    this.window.trigger('message', buildMessage({ isTrusted: false }));
    cancelTimers();
    await settled();

    assert.false(this.handler.called, 'should not call the submit handler');
  });

  test('oidc: it should trigger error callback when role is not found', async function (assert) {
    await renderIt(this, 'oidc');
    await component.role('foo');
    await component.login();
    assert.strictEqual(
      this.error,
      'Invalid role. Please try again.',
      'Error message is returned when role is not found'
    );
  });

  test('oidc: it should trigger error callback when role is returned without auth_url', async function (assert) {
    await renderIt(this, 'oidc');
    await component.role('bar');
    await component.login();
    assert.strictEqual(
      this.error,
      'Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.',
      'Error message is returned when role is returned without auth_url'
    );
  });
});
