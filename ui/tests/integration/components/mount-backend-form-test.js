import Ember from 'ember';
import { moduleForComponent, test } from 'ember-qunit';
import wait from 'ember-test-helpers/wait';
import hbs from 'htmlbars-inline-precompile';

import { create } from 'ember-cli-page-object';
import mountBackendForm from '../../pages/components/mount-backend-form';

import { startMirage } from 'vault/initializers/ember-cli-mirage';
import sinon from 'sinon';

const component = create(mountBackendForm);

moduleForComponent('mount-backend-form', 'Integration | Component | mount backend form', {
  integration: true,
  beforeEach() {
    component.setContext(this);
    Ember.getOwner(this).lookup('service:flash-messages').registerTypes(['success', 'danger']);
    this.server = startMirage();
  },

  afterEach() {
    component.removeContext();
    this.server.shutdown();
  },
});

test('it renders', function(assert) {
  this.render(hbs`{{mount-backend-form}}`);
  assert.equal(component.header, 'Enable an authentication method', 'renders auth header in default state');
  assert.equal(component.fields().count, 2, 'renders 2 fields');
});

test('it changes path when type is changed', function(assert) {
  this.render(hbs`{{mount-backend-form}}`);
  assert.equal(component.pathValue, 'approle', 'defaults to approle (first in the list)');
  component.type('aws');
  assert.equal(component.pathValue, 'aws', 'updates to the value of the type');
});

test('it keeps path value if the user has changed it', function(assert) {
  this.render(hbs`{{mount-backend-form}}`);
  assert.equal(component.pathValue, 'approle', 'defaults to approle (first in the list)');
  component.path('newpath');
  component.type('aws');
  assert.equal(component.pathValue, 'newpath', 'updates to the value of the type');
});

test('it calls mount success', function(assert) {
  const spy = sinon.spy();
  this.set('onMountSuccess', spy);
  this.render(hbs`{{mount-backend-form onMountSuccess=onMountSuccess}}`);

  component.path('foo').type('approle').submit();
  return wait().then(() => {
    assert.equal(this.server.db.authMethods.length, 1, 'it enables an auth method');
    assert.ok(spy.calledOnce, 'calls the passed success method');
  });
});

test('it calls mount mount config error', function(assert) {
  const spy = sinon.spy();
  const spy2 = sinon.spy();
  this.set('onMountSuccess', spy);
  this.set('onConfigError', spy2);
  this.render(hbs`{{mount-backend-form onMountSuccess=onMountSuccess onConfigError=onConfigError}}`);

  component.path('bar').type('kubernetes');
  // kubernetes requires a host + a cert / pem, so only filling the host will error
  component.fields().fillIn('kubernetesHost', 'host');
  component.submit();
  return wait().then(() => {
    assert.equal(this.server.db.authMethods.length, 1, 'it still enables an auth method');
    assert.equal(spy.callCount, 0, 'does not call the success method');
    assert.ok(spy2.calledOnce, 'calls the passed error method');
  });
});
