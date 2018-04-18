import Ember from 'ember';
import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import configPki from '../../pages/components/config-pki-ca';

const component = create(configPki);

const storeStub = Ember.Service.extend({
  createRecord(type, args) {
    return Ember.Object.create(args, {
      save() {
        return Ember.RSVP.resolve(this);
      },
      destroyRecord() {},
      send() {},
      unloadRecord() {},
    });
  },
});

moduleForComponent('config-pki-ca', 'Integration | Component | config pki ca', {
  integration: true,
  beforeEach() {
    component.setContext(this);
    Ember.getOwner(this).lookup('service:flash-messages').registerTypes(['success']);
    this.register('service:store', storeStub);
    this.inject.service('store', { as: 'storeService' });
  },

  afterEach() {
    component.removeContext();
  },
});

const config = function(pem) {
  return Ember.Object.create({
    pem: pem,
    backend: 'pki',
    caChain: 'caChain',
    der: new File(['der'], { type: 'text/plain' }),
  });
};

const setupAndRender = function(context, onRefresh) {
  const refreshFn = onRefresh || function() {};
  context.set('config', config());
  context.set('onRefresh', refreshFn);
  context.render(hbs`{{config-pki-ca onRefresh=onRefresh config=config}}`);
};

test('it renders, no pem', function(assert) {
  setupAndRender(this);

  assert.notOk(component.hasTitle, 'no title in the default state');
  assert.equal(component.replaceCAText, 'Configure CA');
  assert.equal(component.downloadLinks().count, 0, 'there are no download links');

  component.replaceCA();
  assert.equal(component.title, 'Configure CA Certificate');
  component.back();

  component.setSignedIntermediateBtn();
  assert.equal(component.title, 'Set signed intermediate');
  component.back();
});

test('it renders, with pem', function(assert) {
  const c = config('pem');
  this.set('config', c);
  this.render(hbs`{{config-pki-ca config=config}}`);
  assert.notOk(component.hasTitle, 'no title in the default state');
  assert.equal(component.replaceCAText, 'Replace CA');
  assert.equal(component.downloadLinks().count, 3, 'shows download links');
});
