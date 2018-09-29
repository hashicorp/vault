import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | console/log list', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.on('myAction', function(val) { ... });
    const listContent = { keys: ['one', 'two'] };
    const expectedText = 'Keys\none\ntwo';

    this.set('content', listContent);

    await render(hbs`{{console/log-list content=content}}`);

    assert.dom('pre').includesText(`${expectedText}`);
  });
});
