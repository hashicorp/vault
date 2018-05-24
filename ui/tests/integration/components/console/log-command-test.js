import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('console/log-command', 'Integration | Component | console/log command', {
  integration: true,
});

test('it renders', function(assert) {
  const commandText = 'list this/path';
  this.set('content', commandText);

  this.render(hbs`{{console/log-command content=content}}`);

  assert.dom('pre').includesText(commandText);
});
