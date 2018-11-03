import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | empty-state', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`{{empty-state}}`);

    assert.equal(this.element.textContent.trim(), '');

    // Template block usage:
    await render(hbs`
      {{#empty-state
        title="Empty State Title"
        message="This is the empty state message"
      }}
        Actions Link
      {{/empty-state}}
    `);

    assert.equal(
      find('.empty-state-title').textContent.trim(),
      'Empty State Title',
      'renders empty state title'
    );
    assert.equal(
      find('.empty-state-message').textContent.trim(),
      'This is the empty state message',
      'renders empty state message'
    );
    assert.equal(
      find('.empty-state-actions').textContent.trim(),
      'Actions Link',
      'renders empty state actions'
    );
  });
});
