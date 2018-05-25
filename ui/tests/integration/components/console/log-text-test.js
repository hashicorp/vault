import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('console/log-text', 'Integration | Component | console/log text', {
  integration: true,
});

test('it renders', function(assert) {
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });
  const text = 'Success! You did a thing!';
  this.set('content', text);

  this.render(hbs`{{console/log-text content=content}}`);

  assert.dom('pre').includesText(text);
});
