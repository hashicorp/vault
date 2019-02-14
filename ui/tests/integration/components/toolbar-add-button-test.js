import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar-add-button', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ToolbarAddButton @params={{array "/"}}>Link</ToolbarAddButton>`);

    assert.equal(this.element.textContent.trim(), 'Link');
    assert.ok(isPresent('.toolbar-button'));
    assert.ok(isPresent('.icon'));
  });
});
