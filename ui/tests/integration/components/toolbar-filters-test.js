import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar-filters', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ToolbarFilters>These are the toolbar filters</ToolbarFilters>`);

    assert.equal(this.element.textContent.trim(), 'These are the toolbar filters');
    assert.ok(isPresent('.toolbar-filters'));
  });
});
