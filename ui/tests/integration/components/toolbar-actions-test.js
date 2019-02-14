import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar-actions', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ToolbarActions>These are the toolbar actions</ToolbarActions>`);

    assert.equal(this.element.textContent.trim(), 'These are the toolbar actions');
    assert.ok(isPresent('.toolbar-actions'));
  });
});
