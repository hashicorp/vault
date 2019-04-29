import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar-link', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ToolbarLink @params="/secrets">Link</ToolbarLink>`);

    assert.equal(this.element.textContent.trim(), 'Link');
    assert.ok(isPresent('.toolbar-link'));
    assert.ok(isPresent('.icon'));
  });
});
