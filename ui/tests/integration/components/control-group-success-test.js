import { moduleForComponent, test } from 'ember-qunit';
import Ember from 'ember';
import wait from 'ember-test-helpers/wait';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import controlGroupSuccess from '../../pages/components/control-group-success';

const component = create(controlGroupSuccess);

const controlGroupService = Ember.Service.extend({
  init() {
    this.set('wrapInfo', null);
  },
  wrapInfoForAccessor() {
    return this.get('wrapInfo');
  },
  deleteControlGroupToken: sinon.stub(),
});

const routerService = Ember.Service.extend({
  transitionTo() {
    return Ember.RSVP.resolve();
  },
});

const storeService = Ember.Service.extend({
  adapterFor() {
    return {
      toolAction() {
        return Ember.RSVP.resolve({ data: { foo: 'bar' } });
      },
    };
  },
});

moduleForComponent('control-group-success', 'Integration | Component | control group success', {
  integration: true,
  beforeEach() {
    component.setContext(this);
    this.register('service:control-group', controlGroupService);
    this.inject.service('controlGroup');
    this.register('service:router', routerService);
    this.inject.service('router');

    this.register('service:store', storeService);
    this.inject.service('store');
  },

  afterEach() {
    component.removeContext();
  },
});

const setup = (modelData = {}, controlGroupData = {}) => {
  let modelDefaults = {
    approved: false,
    requestPath: 'foo/bar',
    id: 'accessor',
    requestEntity: { id: 'requestor', name: 'entity8509' },
    reload: sinon.stub(),
  };
  let controlGroupDefaults = { entity_id: 'requestor' };

  return {
    model: {
      ...modelDefaults,
      ...modelData,
    },
    authData: {
      ...authDataDefaults,
      ...authData,
    },
  };
};

test('render with saved token', function(assert) {
  let { model, authData } = setup();
  this.set('model', model);
  this.set('auth.authData', authData);
  this.render(hbs`{{control-group model=model}}`);
  assert.ok(component.showsAccessorCallout, 'shows accessor callout');
  assert.equal(component.bannerPrefix, 'Locked');
  assert.equal(component.bannerText, 'The path you requested is locked by a control group');
  assert.equal(component.requestorText, `You are requesting access to ${model.requestPath}`);
  assert.equal(component.showsTokenText, false, 'does not show token message when there is no token');
  assert.ok(component.showsRefresh, 'shows refresh button');
  assert.ok(component.authorizationText, 'Awaiting authorization.');
});

test('render without token', function(assert) {
  let { model, authData } = setup();
  this.set('model', model);
  this.set('auth.authData', authData);
  this.set('controlGroup.wrapInfo', { token: 'token' });
  this.render(hbs`{{control-group model=model}}`);
  assert.equal(component.showsTokenText, true, 'shows token message');
  assert.equal(component.token, 'token', 'shows token value');
});

test('render without token: unwrapped', function(assert) {});
