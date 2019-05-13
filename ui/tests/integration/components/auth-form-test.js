import { later, run } from '@ember/runloop';
import EmberObject from '@ember/object';
import { resolve } from 'rsvp';
import $ from 'jquery';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import Pretender from 'pretender';
import { create } from 'ember-cli-page-object';
import authForm from '../../pages/components/auth-form';

const component = create(authForm);

const authService = Service.extend({
  authenticate() {
    return $.getJSON('http://localhost:2000');
  },
  setLastFetch() {},
});

const workingAuthService = Service.extend({
  authenticate() {
    return resolve({});
  },
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

module('Integration | Component | auth form', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.lookup('service:csp-event').attach();
    this.owner.register('service:router', routerService);
    this.router = this.owner.lookup('service:router');
  });

  hooks.afterEach(function() {
    this.owner.lookup('service:csp-event').remove();
  });

  const CSP_ERR_TEXT = `Error This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`;
  test('it renders error on CSP violation', async function(assert) {
    this.owner.unregister('service:auth');
    this.owner.register('service:auth', authService);
    this.auth = this.owner.lookup('service:auth');
    this.set('cluster', EmberObject.create({ standby: true }));
    this.set('selectedAuth', 'token');
    await render(hbs`{{auth-form cluster=cluster selectedAuth=selectedAuth}}`);
    assert.equal(component.errorText, '');
    component.login();
    // because this is an ember-concurrency backed service,
    // we have to manually force settling the run queue
    later(() => run.cancelTimers(), 50);
    return settled().then(() => {
      assert.equal(component.errorText, CSP_ERR_TEXT);
    });
  });

  test('it renders with vault style errors', async function(assert) {
    let server = new Pretender(function() {
      this.get('/v1/auth/**', () => {
        return [
          400,
          { 'Content-Type': 'application/json' },
          JSON.stringify({
            errors: ['Not allowed'],
          }),
        ];
      });
    });

    this.set('cluster', EmberObject.create({}));
    this.set('selectedAuth', 'token');
    await render(hbs`{{auth-form cluster=cluster selectedAuth=selectedAuth}}`);
    return component.login().then(() => {
      assert.equal(component.errorText, 'Error Authentication failed: Not allowed');
      server.shutdown();
    });
  });

  test('it renders AdapterError style errors', async function(assert) {
    let server = new Pretender(function() {
      this.get('/v1/auth/**', () => {
        return [400, { 'Content-Type': 'application/json' }];
      });
    });

    this.set('cluster', EmberObject.create({}));
    this.set('selectedAuth', 'token');
    await render(hbs`{{auth-form cluster=cluster selectedAuth=selectedAuth}}`);
    return component.login().then(() => {
      assert.equal(component.errorText, 'Error Authentication failed: Bad Request');
      server.shutdown();
    });
  });

  test('it renders no tabs when no methods are passed', async function(assert) {
    let methods = {
      'approle/': {
        type: 'approle',
      },
    };
    let server = new Pretender(function() {
      this.get('/v1/sys/internal/ui/mounts', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify({ data: { auth: methods } })];
      });
    });
    await render(hbs`<AuthForm @cluster={{cluster}} />`);

    await settled();
    assert.equal(component.tabs.length, 0, 'renders a tab for every backend');
    server.shutdown();
  });

  test('it renders all the supported methods and Other tab when methods are present', async function(assert) {
    let methods = {
      'foo/': {
        type: 'userpass',
      },
      'approle/': {
        type: 'approle',
      },
    };
    let server = new Pretender(function() {
      this.get('/v1/sys/internal/ui/mounts', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify({ data: { auth: methods } })];
      });
    });

    this.set('cluster', EmberObject.create({}));
    await render(hbs`{{auth-form cluster=cluster }}`);
    await settled();
    assert.equal(component.tabs.length, 2, 'renders a tab for userpass and Other');
    assert.equal(component.tabs.objectAt(0).name, 'foo', 'uses the path in the label');
    assert.equal(component.tabs.objectAt(1).name, 'Other', 'second tab is the Other tab');
    server.shutdown();
  });

  test('it calls authenticate with the correct path', async function(assert) {
    this.owner.unregister('service:auth');
    this.owner.register('service:auth', workingAuthService);
    this.auth = this.owner.lookup('service:auth');
    let authSpy = sinon.spy(this.get('auth'), 'authenticate');
    let methods = {
      'foo/': {
        type: 'userpass',
      },
    };
    let server = new Pretender(function() {
      this.get('/v1/sys/internal/ui/mounts', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify({ data: { auth: methods } })];
      });
    });

    this.set('cluster', EmberObject.create({}));
    this.set('selectedAuth', 'foo/');
    await render(hbs`{{auth-form cluster=cluster selectedAuth=selectedAuth}}`);
    await component.login();

    await settled();
    assert.ok(authSpy.calledOnce, 'a call to authenticate was made');
    let { data } = authSpy.getCall(0).args[0];
    assert.equal(data.path, 'foo', 'uses the id for the path');
    authSpy.restore();
    server.shutdown();
  });

  test('it renders no tabs when no supported methods are present in passed methods', async function(assert) {
    let methods = {
      'approle/': {
        type: 'approle',
      },
    };
    let server = new Pretender(function() {
      this.get('/v1/sys/internal/ui/mounts', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify({ data: { auth: methods } })];
      });
    });
    this.set('cluster', EmberObject.create({}));
    await render(hbs`<AuthForm @cluster={{cluster}} />`);
    await settled();
    server.shutdown();
    assert.equal(component.tabs.length, 0, 'renders a tab for every backend');
  });

  test('it makes a request to unwrap if passed a wrappedToken and logs in', async function(assert) {
    this.owner.register('service:auth', workingAuthService);
    this.auth = this.owner.lookup('service:auth');
    let authSpy = sinon.spy(this.get('auth'), 'authenticate');
    let server = new Pretender(function() {
      this.post('/v1/sys/wrapping/unwrap', () => {
        return [
          200,
          { 'Content-Type': 'application/json' },
          JSON.stringify({
            auth: {
              client_token: '12345',
            },
          }),
        ];
      });
    });

    let wrappedToken = '54321';
    this.set('wrappedToken', wrappedToken);
    this.set('cluster', EmberObject.create({}));
    await render(hbs`<AuthForm @cluster={{cluster}} @wrappedToken={{wrappedToken}} />`);
    later(() => run.cancelTimers(), 50);
    await settled();
    assert.equal(server.handledRequests[0].url, '/v1/sys/wrapping/unwrap', 'makes call to unwrap the token');
    assert.equal(
      server.handledRequests[0].requestHeaders['X-Vault-Token'],
      wrappedToken,
      'uses passed wrapped token for the unwrap'
    );
    assert.ok(authSpy.calledOnce, 'a call to authenticate was made');
    server.shutdown();
    authSpy.restore();
  });

  test('it shows an error if unwrap errors', async function(assert) {
    let server = new Pretender(function() {
      this.post('/v1/sys/wrapping/unwrap', () => {
        return [
          400,
          { 'Content-Type': 'application/json' },
          JSON.stringify({
            errors: ['There was an error unwrapping!'],
          }),
        ];
      });
    });

    this.set('wrappedToken', '54321');
    await render(hbs`{{auth-form cluster=cluster wrappedToken=wrappedToken}}`);
    later(() => run.cancelTimers(), 50);

    await settled();
    assert.equal(
      component.errorText,
      'Error Token unwrap failed: There was an error unwrapping!',
      'shows the error'
    );
    server.shutdown();
  });
});
