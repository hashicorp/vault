import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | secret edit', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.codeMirror = this.owner.lookup('service:code-mirror');
  });

  test('it disables JSON toggle in show mode when is an advanced format', async function(assert) {
    this.set('mode', 'show');
    this.set('key', {
      secretData: {
        int: 2,
        null: null,
        float: 1.234,
      },
    });

    await render(hbs`{{secret-edit mode=mode key=key }}`);
    assert.dom('[data-test-secret-json-toggle]').isDisabled();
  });

  test('it does JSON toggle in show mode when showing string data', async function(assert) {
    this.set('mode', 'show');
    this.set('key', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });

    await render(hbs`{{secret-edit mode=mode key=key }}`);
    assert.dom('[data-test-secret-json-toggle]').isNotDisabled();
  });

  test('it shows an error when creating and data is not an object', async function(assert) {
    this.set('mode', 'create');
    this.set('key', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });

    await render(hbs`{{secret-edit mode=mode key=key preferAdvancedEdit=true }}`);
    let instance = this.codeMirror.instanceFor(find('[data-test-component=json-editor]').id);
    instance.setValue(JSON.stringify([{ foo: 'bar' }]));
    await settled();
    assert.dom('[data-test-error]').includesText('Vault expects data to be formatted as an JSON object');
  });

  test('it shows an error when editing and the data is not an object', async function(assert) {
    this.set('mode', 'edit');
    this.set('capabilities', {
      canUpdate: true,
    });
    this.set('key', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });

    await render(hbs`{{secret-edit capabilities=capabilities mode=mode key=key preferAdvancedEdit=true }}`);
    let instance = this.codeMirror.instanceFor(find('[data-test-component=json-editor]').id);
    instance.setValue(JSON.stringify([{ foo: 'bar' }]));
    await settled();
    assert.dom('[data-test-error]').includesText('Vault expects data to be formatted as an JSON object');
  });
});
