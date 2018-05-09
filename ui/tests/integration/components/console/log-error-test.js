import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('console/log-error', 'Integration | Component | console/log error', {
  integration: true
});

test('it renders', function(assert) {

  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });
  const errorText = "Error deleting at: sys/foo.\nURL: v1/sys/foo\nCode: 404";
  this.set('content', errorText);

  this.render(hbs`{{console/log-error content=content}}`);

  assert.dom('pre').hasClass('has-text-danger');
  assert.dom('pre').includesText(errorText);

});
