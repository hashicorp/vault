import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | console/log command', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    const commandText = 'list this/path';
    this.set('content', commandText);

    await render(hbs`{{console/log-command content=content}}`);

    assert.dom('pre').includesText(commandText);
  });
});
