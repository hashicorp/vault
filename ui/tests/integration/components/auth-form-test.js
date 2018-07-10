import { moduleForComponent, test } from 'ember-qunit';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import Ember from 'ember';
import wait from 'ember-test-helpers/wait';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import Pretender from 'pretender';
import { create } from 'ember-cli-page-object';
import authForm from '../../pages/components/auth-form';

const component = create(authForm);
const BACKENDS = supportedAuthBackends();

const authService = Ember.Service.extend({
  authenticate() {
    return Ember.$.getJSON('http://localhost:2000');
  },
});

const workingAuthService = Ember.Service.extend({
  authenticate() {
    return Ember.RSVP.resolve({});
  },
  setLastFetch() {},
});

const routerService = Ember.Service.extend({
  transitionTo() {
    return Ember.RSVP.resolve();
  },
  replaceWith() {
    return Ember.RSVP.resolve();
  },
});
moduleForComponent('auth-form', 'Integration | Component | auth form', {
  integration: true,
  beforeEach() {
    Ember.getOwner(this).lookup('service:csp-event').attach();
    component.setContext(this);
    this.register('service:router', routerService);
    this.inject.service('router');
  },

  afterEach() {
    Ember.getOwner(this).lookup('service:csp-event').remove();
    component.removeContext();
  },
});

const CSP_ERR_TEXT = `Error This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`;
test('it renders error on CSP violation', function(assert) {
  this.register('service:auth', authService);
  this.inject.service('auth');
  this.set('cluster', Ember.Object.create({ standby: true }));
  this.set('selectedAuth', 'token');
  this.render(hbs`{{auth-form cluster=cluster selectedAuth=selectedAuth}}`);
  assert.equal(component.errorText, '');
  component.login();
  // because this is an ember-concurrency backed service,
  // we have to manually force settling the run queue
  Ember.run.later(() => Ember.run.cancelTimers(), 50);
  return wait().then(() => {
    assert.equal(component.errorText, CSP_ERR_TEXT);
  });
});

test('it renders with vault style errors', function(assert) {
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

  this.set('cluster', Ember.Object.create({}));
  this.set('selectedAuth', 'token');
  this.render(hbs`{{auth-form cluster=cluster selectedAuth=selectedAuth}}`);
  return component.login().then(() => {
    assert.equal(component.errorText, 'Error Authentication failed: Not allowed');
    server.shutdown();
  });
});

test('it renders AdapterError style errors', function(assert) {
  let server = new Pretender(function() {
    this.get('/v1/auth/**', () => {
      return [400, { 'Content-Type': 'application/json' }];
    });
  });

  this.set('cluster', Ember.Object.create({}));
  this.set('selectedAuth', 'token');
  this.render(hbs`{{auth-form cluster=cluster selectedAuth=selectedAuth}}`);
  return component.login().then(() => {
    assert.equal(component.errorText, 'Error Authentication failed: Bad Request');
    server.shutdown();
  });
});

test('it renders all the supported tabs when no methods are passed', function(assert) {
  this.render(hbs`{{auth-form cluster=cluster}}`);
  assert.equal(component.tabs.length, BACKENDS.length, 'renders a tab for every backend');
});

test('it renders all the supported methods and Other tab when methods are present', function(assert) {
  let methods = [
    {
      type: 'userpass',
      id: 'foo',
      path: 'foo/',
    },
    {
      type: 'approle',
      id: 'approle',
      path: 'approle/',
    },
  ];
  this.set('methods', methods);

  this.render(hbs`{{auth-form cluster=cluster methods=methods}}`);
  assert.equal(component.tabs.length, 2, 'renders a tab for userpass and Other');
  assert.equal(component.tabs.objectAt(0).name, 'foo', 'uses the path in the label');
  assert.equal(component.tabs.objectAt(1).name, 'Other', 'second tab is the Other tab');
});

test('it renders all the supported methods when no supported methods are present in passed methods', function(
  assert
) {
  let methods = [
    {
      type: 'approle',
      id: 'approle',
      path: 'approle/',
    },
  ];
  this.set('methods', methods);
  this.render(hbs`{{auth-form cluster=cluster methods=methods}}`);
  assert.equal(component.tabs.length, BACKENDS.length, 'renders a tab for every backend');
});

test('it makes a request to unwrap if passed a wrappedToken and logs in', function(assert) {
  this.register('service:auth', workingAuthService);
  this.inject.service('auth');
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
  this.render(hbs`{{auth-form cluster=cluster wrappedToken=wrappedToken}}`);
  Ember.run.later(() => Ember.run.cancelTimers(), 50);
  return wait().then(() => {
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
});

test('it shows an error if unwrap errors', function(assert) {
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
  this.render(hbs`{{auth-form cluster=cluster wrappedToken=wrappedToken}}`);
  Ember.run.later(() => Ember.run.cancelTimers(), 50);

  return wait().then(() => {
    assert.equal(
      component.errorText,
      'Error Token unwrap failed: There was an error unwrapping!',
      'shows the error'
    );
    server.shutdown();
  });
});
