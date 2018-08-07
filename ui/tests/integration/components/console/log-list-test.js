import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('console/log-list', 'Integration | Component | console/log list', {
  integration: true,
});

test('it renders', function(assert) {
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });
  const listContent = { keys: ['one', 'two'] };
  const expectedText = 'Keys\none\ntwo';

  this.set('content', listContent);

  this.render(hbs`{{console/log-list content=content}}`);

  assert.dom('pre').includesText(`${expectedText}`);
});
