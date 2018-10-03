import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | edition badge', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      {{edition-badge edition="Custom"}}
    `);

    assert.equal(find('.edition-badge').textContent.trim(), 'Custom', 'contains edition');

    await render(hbs`
      {{edition-badge edition="Enterprise"}}
    `);

    assert.equal(find('.edition-badge').textContent.trim(), 'Enterprise', 'renders edition');
  });
});
