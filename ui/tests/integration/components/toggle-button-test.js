import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toggle button', function(hooks) {
  setupRenderingTest(hooks);

  test('toggle functionality', async function(assert) {
    this.set('toggleTarget', {});

    await render(hbs`{{toggle-button toggleTarget=toggleTarget toggleAttr="toggled"}}`);

    assert.equal(find('button').textContent.trim(), 'More options', 'renders default closedLabel');

    await click('button');
    assert.equal(this.get('toggleTarget.toggled'), true, 'it toggles the attr on the target');
    assert.equal(find('button').textContent.trim(), 'Hide options', 'renders default openLabel');
    await click('button');
    assert.equal(this.get('toggleTarget.toggled'), false, 'it toggles the attr on the target');

    this.set('closedLabel', 'Open the options!');
    this.set('openLabel', 'Close the options!');
    await render(
      hbs`{{toggle-button toggleTarget=toggleTarget toggleAttr="toggled" closedLabel=closedLabel openLabel=openLabel}}`
    );

    assert.equal(find('button').textContent.trim(), 'Open the options!', 'renders passed closedLabel');
    await click('button');
    assert.equal(find('button').textContent.trim(), 'Close the options!', 'renders passed openLabel');
  });
});
