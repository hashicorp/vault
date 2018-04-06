import { moduleForComponent, test } from 'ember-qunit';
import Ember from 'ember';
import wait from 'ember-test-helpers/wait';
import hbs from 'htmlbars-inline-precompile';

import Pretender from 'pretender';
import { create } from 'ember-cli-page-object';
import authForm from '../../pages/components/auth-form';

const component = create(authForm);

const authService = Ember.Service.extend({
  authenticate() {
    return Ember.$.getJSON('http://localhost:2000');
  },
});

moduleForComponent('auth-form', 'Integration | Component | auth form', {
  integration: true,
  beforeEach() {
    Ember.getOwner(this).lookup('service:csp-event').attach();
    component.setContext(this);
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
  this.render(hbs`{{auth-form cluster=cluster}}`);
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
  this.render(hbs`{{auth-form cluster=cluster}}`);
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
  this.render(hbs`{{auth-form cluster=cluster}}`);
  return component.login().then(() => {
    assert.equal(component.errorText, 'Error Authentication failed: Bad Request');
    server.shutdown();
  });
});
