import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('secret-edit', 'Integration | Component | secret edit', {
  integration: true,
  beforeEach() {
    this.inject.service('code-mirror', { as: 'codeMirror' });
  },
});

test('it disables JSON toggle in show mode when is an advanced format', function(assert) {
  this.set('mode', 'show');
  this.set('key', {
    secretData: {
      int: 2,
      null: null,
      float: 1.234,
    },
  });

  this.render(hbs`{{secret-edit mode=mode key=key }}`);
  assert.dom('[data-test-secret-json-toggle]').isDisabled();
});

test('it does JSON toggle in show mode when showing string data', function(assert) {
  this.set('mode', 'show');
  this.set('key', {
    secretData: {
      int: '2',
      null: 'null',
      float: '1.234',
    },
  });

  this.render(hbs`{{secret-edit mode=mode key=key }}`);
  assert.dom('[data-test-secret-json-toggle]').isNotDisabled();
});

test('it shows an error when creating and data is not an object', function(assert) {
  this.set('mode', 'create');
  this.set('key', {
    secretData: {
      int: '2',
      null: 'null',
      float: '1.234',
    },
  });

  this.render(hbs`{{secret-edit mode=mode key=key preferAdvancedEdit=true }}`);
  let instance = this.codeMirror.instanceFor(this.$('[data-test-component=json-editor]').attr('id'));
  instance.setValue(JSON.stringify([{ foo: 'bar' }]));
  assert.dom('[data-test-error]').includesText('Vault expects data to be formatted as an JSON object');
});

test('it shows an error when editing and the data is not an object', function(assert) {
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

  this.render(hbs`{{secret-edit capabilities=capabilities mode=mode key=key preferAdvancedEdit=true }}`);
  let instance = this.codeMirror.instanceFor(this.$('[data-test-component=json-editor]').attr('id'));
  instance.setValue(JSON.stringify([{ foo: 'bar' }]));
  assert.dom('[data-test-error]').includesText('Vault expects data to be formatted as an JSON object');
});
