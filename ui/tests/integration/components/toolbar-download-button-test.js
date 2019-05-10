import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar-download-button', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ToolbarDownloadButton @actionText="Link" />`);

    assert.equal(this.element.textContent.trim(), 'Link');
    assert.ok(isPresent('.toolbar-link'));
    assert.ok(isPresent('.icon'));
  });
});
