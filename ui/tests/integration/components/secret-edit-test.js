import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('secret-edit', 'Integration | Component | secret edit', {
  integration: true,
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

test('it does JSON toggle in show mode when is', function(assert) {
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
