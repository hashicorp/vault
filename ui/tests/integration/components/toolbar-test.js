import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<Toolbar>This is the toolbar content</Toolbar>`);

    assert.dom(this.element).hasText('This is the toolbar content');
    assert.ok(isPresent('.toolbar'));
    assert.ok(isPresent('.toolbar-scroller'));
  });
});
