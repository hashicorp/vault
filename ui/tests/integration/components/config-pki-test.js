import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import configPki from '../../pages/components/config-pki';

const component = create(configPki);

module('Integration | Component | config pki', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
  });

  hooks.afterEach(function() {
    component.removeContext();
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

  const setupAndRender = async function(context, section = 'tidy') {
    context.set('config', config());
    context.set('section', section);
    await context.render(hbs`{{config-pki section=section config=config}}`);
  };

  test('it renders tidy section', async function(assert) {
    await setupAndRender(this);
    assert.ok(component.text.startsWith('You can tidy up the backend'));
    assert.notOk(component.hasTitle, 'No title for tidy section');
    assert.equal(component.fields.length, 2);
    assert.ok(component.fields.objectAt(0).labelText, 'Tidy cert store');
    assert.ok(component.fields.objectAt(1).labelText, 'Another attr');
  });

  test('it renders crl section', async function(assert) {
    await setupAndRender(this, 'crl');
    assert.ok(component.hasTitle, 'renders the title');
    assert.equal(component.title, 'Certificate Revocation List (CRL) config');
    assert.ok(component.text.startsWith('Set the duration for which the generated CRL'));
    assert.equal(component.fields.length, 1);
    assert.ok(component.fields.objectAt(0).labelText, 'Crl');
  });

  test('it renders urls section', async function(assert) {
    await setupAndRender(this, 'urls');
    assert.notOk(component.hasTitle, 'No title for urls section');
    assert.equal(component.fields.length, 1);
    assert.ok(component.fields.objectAt(0).labelText, 'urls');
  });

  test('it calls save with the correct arguments', async function(assert) {
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
        return resolve();
      })
    );
    this.set('section', section);
    await render(hbs`{{config-pki section=section config=config onRefresh=onRefresh}}`);

    component.submit();
  });
});
