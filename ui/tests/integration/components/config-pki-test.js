import Ember from 'ember';
import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import configPki from '../../pages/components/config-pki';

const component = create(configPki);

moduleForComponent('config-pki', 'Integration | Component | config pki', {
  integration: true,
  beforeEach() {
    component.setContext(this);
    Ember.getOwner(this).lookup('service:flash-messages').registerTypes(['success']);
  },

  afterEach() {
    component.removeContext();
  },
});

const config = function(saveFn) {
  return {
    save: saveFn,
    rollbackAttributes: () => {},
    tidyAttrs: [
      {
        type: 'boolean',
        name: 'tidyCertStore',
      },
      {
        type: 'string',
        name: 'anotherAttr',
      },
    ],
    crlAttrs: [
      {
        type: 'string',
        name: 'crl',
      },
    ],
    urlsAttrs: [
      {
        type: 'string',
        name: 'urls',
      },
    ],
  };
};

const setupAndRender = function(context, section = 'tidy') {
  context.set('config', config());
  context.set('section', section);
  context.render(hbs`{{config-pki section=section config=config}}`);
};

test('it renders tidy section', function(assert) {
  setupAndRender(this);
  assert.ok(component.text.startsWith('You can tidy up the backend'));
  assert.notOk(component.hasTitle, 'No title for tidy section');
  assert.equal(component.fields().count, 2);
  assert.ok(component.fields(0).labelText, 'Tidy cert store');
  assert.ok(component.fields(1).labelText, 'Another attr');
});

test('it renders crl section', function(assert) {
  setupAndRender(this, 'crl');
  assert.ok(component.hasTitle, 'renders the title');
  assert.equal(component.title, 'Certificate Revocation List (CRL) Config');
  assert.ok(component.text.startsWith('Set the duration for which the generated CRL'));
  assert.equal(component.fields().count, 1);
  assert.ok(component.fields(0).labelText, 'Crl');
});

test('it renders urls section', function(assert) {
  setupAndRender(this, 'urls');
  assert.notOk(component.hasTitle, 'No title for urls section');
  assert.equal(component.fields().count, 1);
  assert.ok(component.fields(0).labelText, 'urls');
});

test('it calls save with the correct arguments', function(assert) {
  assert.expect(3);
  const section = 'tidy';
  this.set('onRefresh', () => {
    assert.ok(true, 'refresh called');
  });
  this.set(
    'config',
    config(options => {
      assert.equal(options.adapterOptions.method, section, 'method passed to save');
      assert.deepEqual(
        options.adapterOptions.fields,
        ['tidyCertStore', 'anotherAttr'],
        'fields passed to save'
      );
      return Ember.RSVP.resolve();
    })
  );
  this.set('section', section);
  this.render(hbs`{{config-pki section=section config=config onRefresh=onRefresh}}`);

  component.submit();
});
