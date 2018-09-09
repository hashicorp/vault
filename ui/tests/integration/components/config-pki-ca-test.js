import { resolve } from 'rsvp';
import EmberObject from '@ember/object';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import configPki from '../../pages/components/config-pki-ca';

const component = create(configPki);

const storeStub = Service.extend({
  createRecord(type, args) {
    return EmberObject.create(args, {
      save() {
        return resolve(this);
      },
      destroyRecord() {},
      send() {},
      unloadRecord() {},
    });
  },
});

module('Integration | Component | config pki ca', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
    this.owner.register('service:store', storeStub);
    this.storeService = this.owner.lookup('service:store');
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  const config = function(pem) {
    return EmberObject.create({
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

  test('it renders, with pem', async function(assert) {
    const c = config('pem');
    this.set('config', c);
    await render(hbs`{{config-pki-ca config=config}}`);
    assert.notOk(component.hasTitle, 'no title in the default state');
    assert.equal(component.replaceCAText, 'Replace CA');
    assert.equal(component.downloadLinks().count, 3, 'shows download links');
  });
});
