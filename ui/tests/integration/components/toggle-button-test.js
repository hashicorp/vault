import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toggle button', function(hooks) {
  setupRenderingTest(hooks);

  test('toggle functionality', async function(assert) {
    this.set('toggleTarget', {});

    await render(hbs`{{toggle-button toggleTarget=toggleTarget toggleAttr="toggled"}}`);

    assert.dom('button').hasText('More options', 'renders default closedLabel');

    await click('button');
    assert.equal(this.toggleTarget.toggled, true, 'it toggles the attr on the target');
    assert.dom('button').hasText('Hide options', 'renders default openLabel');
    await click('button');
    assert.equal(this.toggleTarget.toggled, false, 'it toggles the attr on the target');

    this.set('closedLabel', 'Open the options!');
    this.set('openLabel', 'Close the options!');
    await render(
      hbs`{{toggle-button toggleTarget=toggleTarget toggleAttr="toggled" closedLabel=closedLabel openLabel=openLabel}}`
    );

    assert.dom('button').hasText('Open the options!', 'renders passed closedLabel');
    await click('button');
    assert.dom('button').hasText('Close the options!', 'renders passed openLabel');
  });
});
