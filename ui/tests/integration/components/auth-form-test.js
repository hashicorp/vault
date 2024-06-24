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
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import authForm from '../../pages/components/auth-form';
import { validate } from 'uuid';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';

const component = create(authForm);

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
    this.onSuccess = sinon.spy();
  });

  const CSP_ERR_TEXT = `Error This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`;
  test('it renders error on CSP violation', async function (assert) {
    assert.expect(2);
    this.set('cluster', EmberObject.create({ standby: true }));
    this.set('selectedAuth', 'token');
    await render(hbs`<AuthForm @cluster={{this.cluster}} @selectedAuth={{this.selectedAuth}} />`);
    assert.false(component.errorMessagePresent, false);
    this.owner.lookup('service:csp-event').handleEvent({ violatedDirective: 'connect-src' });
    await settled();
    assert.strictEqual(component.errorText, CSP_ERR_TEXT);
  });

  test('it renders with vault style errors', async function (assert) {
    assert.expect(1);
    this.server.get('/auth/token/lookup-self', () => {
      return new Response(400, { 'Content-Type': 'application/json' }, { errors: ['Not allowed'] });
    });

    this.set('cluster', EmberObject.create({}));
    this.set('selectedAuth', 'token');
    await render(hbs`<AuthForm @cluster={{this.cluster}} @selectedAuth={{this.selectedAuth}} />`);
    await component.login();
    assert.strictEqual(component.errorText, 'Error Authentication failed: Not allowed');
  });

  test('it renders AdapterError style errors', async function (assert) {
    assert.expect(1);
    this.server.get('/auth/token/lookup-self', () => {
      return new Response(400, { 'Content-Type': 'application/json' }, { errors: ['API Error here'] });
    });

    this.set('cluster', EmberObject.create({}));
    this.set('selectedAuth', 'token');
    await render(hbs`<AuthForm @cluster={{this.cluster}} @selectedAuth={{this.selectedAuth}} />`);
    return component.login().then(() => {
      assert.strictEqual(
        component.errorText,
        'Error Authentication failed: API Error here',
        'shows the error from the API'
      );
    });
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
    await render(hbs`<AuthForm @cluster={{this.cluster}} />`);

    assert.strictEqual(component.tabs.length, 0, 'renders a tab for every backend');
    server.shutdown();
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
    await render(hbs`<AuthForm @cluster={{this.cluster}} />`);

    assert.strictEqual(component.tabs.length, 2, 'renders a tab for userpass and Other');
    assert.strictEqual(component.tabs.objectAt(0).name, 'foo', 'uses the path in the label');
    assert.strictEqual(component.tabs.objectAt(1).name, 'Other', 'second tab is the Other tab');
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
    await render(hbs`<AuthForm @cluster={{this.cluster}} />`);

    assert.strictEqual(
      component.descriptionText,
      'app description',
      'renders a description for auth methods'
    );
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
    await render(hbs`<AuthForm @cluster={{this.cluster}} @selectedAuth={{this.selectedAuth}} />`);
    await component.login();

    await settled();
    assert.ok(authSpy.calledOnce, 'a call to authenticate was made');
    const { data } = authSpy.getCall(0).args[0];
    assert.strictEqual(data.path, 'foo', 'uses the id for the path');
    authSpy.restore();
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
    await render(hbs`<AuthForm @cluster={{this.cluster}} />`);

    server.shutdown();
    assert.strictEqual(component.tabs.length, 0, 'renders a tab for every backend');
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
    await render(
      hbs`<AuthForm @cluster={{this.cluster}} @wrappedToken={{this.wrappedToken}} @onSuccess={{this.onSuccess}} />`
    );
    later(() => cancelTimers(), 50);
    await settled();
    assert.ok(authSpy.calledOnce, 'a call to authenticate was made');
    authSpy.restore();
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
    await render(hbs`<AuthForm @cluster={{this.cluster}} @wrappedToken={{this.wrappedToken}} />`);
    later(() => cancelTimers(), 50);
    await settled();
    assert.strictEqual(
      component.errorText,
      'Error Token unwrap failed: There was an error unwrapping!',
      'shows the error'
    );
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
    await render(hbs`<AuthForm @cluster={{this.cluster}} />`);

    await component.selectMethod('oidc');
    await component.oidcRole('foo');
    await component.oidcMoreOptions();
    await component.oidcMountPath('foo-oidc');
    assert.dom('[data-test-role]').hasValue('foo', 'role is retained when mount path is changed');
    await component.login();
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
    await render(hbs`<AuthForm @cluster={{this.cluster}} />`);

    await component.selectMethod('okta');
    await component.username('foo');
    await component.password('bar');
    await component.login();
  });
});
