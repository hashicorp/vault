import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, settled } from '@ember/test-helpers';
import { resolve } from 'rsvp';
import { run } from '@ember/runloop';
import Service from '@ember/service';
import hbs from 'htmlbars-inline-precompile';

let capabilities;
const storeService = Service.extend({
  queryRecord() {
    return resolve(capabilities);
  },
});
module('Integration | Component | secret edit', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    capabilities = null;
    this.codeMirror = this.owner.lookup('service:code-mirror');
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeService);
    });
  });

  test('it disables JSON toggle in show mode when is an advanced format', async function(assert) {
    this.set('mode', 'show');
    this.set('model', {
      secretData: {
        int: 2,
        null: null,
        float: 1.234,
      },
    });

    await render(hbs`{{secret-edit mode=mode model=model }}`);
    assert.dom('[data-test-secret-json-toggle]').isDisabled();
  });

  test('it does JSON toggle in show mode when showing string data', async function(assert) {
    this.set('mode', 'show');
    this.set('model', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });

    await render(hbs`{{secret-edit mode=mode model=model }}`);
    assert.dom('[data-test-secret-json-toggle]').isNotDisabled();
  });

  test('it shows an error when creating and data is not an object', async function(assert) {
    this.set('mode', 'create');
    this.set('model', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });

    await render(hbs`{{secret-edit mode=mode model=model preferAdvancedEdit=true }}`);
    let instance = this.codeMirror.instanceFor(find('[data-test-component=json-editor]').id);
    instance.setValue(JSON.stringify([{ foo: 'bar' }]));
    await settled();
    assert.dom('[data-test-error]').includesText('Vault expects data to be formatted as an JSON object');
  });

  test('it allows saving when the model isError', async function(assert) {
    this.set('mode', 'create');
    this.set('model', {
      isError: true,
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });
    await render(hbs`<SecretEdit @mode={{mode}} @model={{model}} />`);
    assert.dom('[data-test-secret-save]').isNotDisabled();
  });

  test('it shows an error when editing and the data is not an object', async function(assert) {
    this.set('mode', 'edit');
    capabilities = {
      canUpdate: true,
    };
    this.set('model', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });

    await render(hbs`{{secret-edit mode=mode model=model preferAdvancedEdit=true }}`);
    let instance = this.codeMirror.instanceFor(find('[data-test-component=json-editor]').id);
    instance.setValue(JSON.stringify([{ foo: 'bar' }]));
    await settled();
    assert.dom('[data-test-error]').includesText('Vault expects data to be formatted as an JSON object');
  });
});
