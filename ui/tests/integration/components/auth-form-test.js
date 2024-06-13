/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { later, _cancelTimers as cancelTimers } from '@ember/runloop';
import EmberObject from '@ember/object';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { validate } from 'uuid';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

const workingAuthService = Service.extend({
  authenticate() {
    return resolve({});
  },
  handleError() {},
  setLastFetch() {},
});

const routerService = Service.extend({
  transitionTo() {
    return {
      followRedirects() {
        return resolve();
      },
    };
  },
});

module('Integration | Component | auth form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.owner.register('service:router', routerService);
    this.router = this.owner.lookup('service:router');
    this.handleAuth = {
      unlinked() {
        return {
          perform: () => sinon.spy(),
        };
      },
    };
    this.onSuccess = sinon.spy();
    this.renderComponent = async () => {
      return render(hbs`
        <AuthForm
          @wrappedToken={{this.wrappedToken}}
          @cluster={{this.cluster}}
          @selectedAuth={{this.selectedAuth}}
          @handleAuth={{this.handleAuth}}
          @delayIsIdle={{this.delayIsIdle}}
        />`);
    };

    this.renderParent = async () => {
      return render(hbs`
        <Auth::Page
          @wrappedToken={{this.wrappedToken}}
          @cluster={{this.cluster}}
          @namespace={{this.namespaceQueryParam}}
          @selectedAuth={{this.authMethod}}
          @onSuccess={{this.onSuccess}}
        />
        `);
    };
  });

  const CSP_ERR_TEXT = `Error This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`;
  test('it renders error on CSP violation', async function (assert) {
    assert.expect(2);
    this.set('cluster', EmberObject.create({ standby: true }));
    this.set('selectedAuth', 'token');
    await this.renderParent();
    assert.dom(GENERAL.messageError).doesNotExist();
    this.owner.lookup('service:csp-event').handleEvent({ violatedDirective: 'connect-src' });
    await settled();
    assert.dom(GENERAL.messageError).hasText(CSP_ERR_TEXT);
  });

  test('it renders with vault style errors', async function (assert) {
    assert.expect(1);
    this.server.get('/auth/token/lookup-self', () => {
      return new Response(400, { 'Content-Type': 'application/json' }, { errors: ['Not allowed'] });
    });

    this.set('cluster', EmberObject.create({}));
    this.set('selectedAuth', 'token');
    await this.renderParent();
    await click(AUTH_FORM.login);
    assert.dom(GENERAL.messageError).hasText('Error Authentication failed: Not allowed');
  });

  test('it renders AdapterError style errors', async function (assert) {
    assert.expect(1);
    this.server.get('/auth/token/lookup-self', () => {
      return new Response(400, { 'Content-Type': 'application/json' }, { errors: ['API Error here'] });
    });

    this.set('cluster', EmberObject.create({}));
    this.set('selectedAuth', 'token');
    await this.renderParent();
    await click(AUTH_FORM.login);
    assert
      .dom(GENERAL.messageError)
      .hasText('Error Authentication failed: API Error here', 'shows the error from the API');
  });

  test('it renders no tabs when no methods are passed', async function (assert) {
    const methods = {
      'approle/': {
        type: 'approle',
      },
    };
    this.server.get('/sys/internal/ui/mounts', () => {
      return { data: { auth: methods } };
    });
    await this.renderComponent();

    assert.dom(AUTH_FORM.tabs()).doesNotExist();
  });

  test('it renders all the supported methods and Other tab when methods are present', async function (assert) {
    const methods = {
      'foo/': {
        type: 'userpass',
      },
      'approle/': {
        type: 'approle',
      },
    };
    this.server.get('/sys/internal/ui/mounts', () => {
      return { data: { auth: methods } };
    });
    this.set('cluster', EmberObject.create({}));
    await this.renderComponent();

    assert.dom(AUTH_FORM.tabs()).exists({ count: 2 });
    assert.dom(AUTH_FORM.tabs('foo')).exists('tab uses the path in the label');
    assert.dom(AUTH_FORM.tabs('other')).exists('second tab is the Other tab');
  });

  test('it renders the description', async function (assert) {
    const methods = {
      'approle/': {
        type: 'userpass',
        description: 'app description',
      },
    };
    this.server.get('/sys/internal/ui/mounts', () => {
      return { data: { auth: methods } };
    });
    this.set('cluster', EmberObject.create({}));
    await this.renderComponent();

    assert.dom(AUTH_FORM.description).hasText('app description');
  });

  test('it calls authenticate with the correct path', async function (assert) {
    this.owner.unregister('service:auth');
    this.owner.register('service:auth', workingAuthService);
    this.auth = this.owner.lookup('service:auth');
    const authSpy = sinon.spy(this.auth, 'authenticate');
    const methods = {
      'foo/': {
        type: 'userpass',
      },
    };
    this.server.get('/sys/internal/ui/mounts', () => {
      return { data: { auth: methods } };
    });

    this.set('cluster', EmberObject.create({}));
    this.set('selectedAuth', 'foo/');
    await this.renderParent();
    await click(AUTH_FORM.login);
    assert.ok(authSpy.calledOnce, 'a call to authenticate was made');
    const { data } = authSpy.getCall(0).args[0];
    assert.strictEqual(data.path, 'foo', 'uses the id for the path');
  });

  test('it renders no tabs when no supported methods are present in passed methods', async function (assert) {
    const methods = {
      'approle/': {
        type: 'approle',
      },
    };
    this.server.get('/sys/internal/ui/mounts', () => {
      return { data: { auth: methods } };
    });
    this.set('cluster', EmberObject.create({}));
    await this.renderComponent();

    assert.dom(AUTH_FORM.tabs()).doesNotExist();
  });

  test('it makes a request to unwrap if passed a wrappedToken and logs in', async function (assert) {
    assert.expect(3);
    this.owner.register('service:auth', workingAuthService);
    this.auth = this.owner.lookup('service:auth');
    const authSpy = sinon.stub(this.auth, 'authenticate');
    this.server.post('/sys/wrapping/unwrap', (_, req) => {
      assert.strictEqual(req.url, '/v1/sys/wrapping/unwrap', 'makes call to unwrap the token');
      assert.strictEqual(
        req.requestHeaders['X-Vault-Token'],
        wrappedToken,
        'uses passed wrapped token for the unwrap'
      );
      return {
        auth: {
          client_token: '12345',
        },
      };
    });

    const wrappedToken = '54321';
    this.set('wrappedToken', wrappedToken);
    this.set('cluster', EmberObject.create({}));
    await this.renderParent();
    later(() => cancelTimers(), 50);
    await settled();
    assert.ok(authSpy.calledOnce, 'a call to authenticate was made');
  });

  test('it shows an error if unwrap errors', async function (assert) {
    this.server.post('/sys/wrapping/unwrap', () => {
      return new Response(
        400,
        { 'Content-Type': 'application/json' },
        { errors: ['There was an error unwrapping!'] }
      );
    });

    this.set('wrappedToken', '54321');
    await this.renderComponent();
    later(() => cancelTimers(), 50);

    await settled();
    assert.dom(GENERAL.messageError).hasText('Error Token unwrap failed: There was an error unwrapping!');
  });

  test('it should retain oidc role when mount path is changed', async function (assert) {
    assert.expect(2);

    const auth_url = 'http://dev-foo-bar.com';
    this.server.post('/auth/:path/oidc/auth_url', (_, req) => {
      const { role, redirect_uri } = JSON.parse(req.requestBody);
      const goodRequest =
        req.params.path === 'foo-oidc' &&
        role === 'foo' &&
        redirect_uri.includes('/auth/foo-oidc/oidc/callback');

      return new Response(
        goodRequest ? 200 : 400,
        { 'Content-Type': 'application/json' },
        JSON.stringify(
          goodRequest ? { data: { auth_url } } : { errors: [`role "${role}" could not be found`] }
        )
      );
    });
    window.open = (url) => {
      assert.strictEqual(url, auth_url, 'auth_url is returned when required params are passed');
    };

    this.owner.lookup('service:router').reopen({
      urlFor(route, { auth_path }) {
        return `/auth/${auth_path}/oidc/callback`;
      },
    });

    this.set('cluster', EmberObject.create({}));
    await this.renderComponent();

    await fillIn(GENERAL.selectByAttr('auth-method'), 'oidc');
    await fillIn(AUTH_FORM.input('role'), 'foo');
    await click(AUTH_FORM.moreOptions);
    await fillIn(AUTH_FORM.input('role'), 'foo');
    await fillIn(AUTH_FORM.mountPathInput, 'foo-oidc');
    assert.dom(AUTH_FORM.input('role')).hasValue('foo', 'role is retained when mount path is changed');
    await click(AUTH_FORM.login);
  });

  test('it should set nonce value as uuid for okta method type', async function (assert) {
    assert.expect(1);

    this.server.post('/auth/okta/login/foo', (_, req) => {
      const { nonce } = JSON.parse(req.requestBody);
      assert.true(validate(nonce), 'Nonce value passed as uuid for okta login');
      return {
        auth: {
          client_token: '12345',
        },
      };
    });

    this.set('cluster', EmberObject.create({}));
    await this.renderParent();

    await fillIn(GENERAL.selectByAttr('auth-method'), 'okta');
    await fillIn(AUTH_FORM.input('username'), 'foo');
    await fillIn(AUTH_FORM.input('password'), 'bar');
    await click(AUTH_FORM.login);
  });
});
