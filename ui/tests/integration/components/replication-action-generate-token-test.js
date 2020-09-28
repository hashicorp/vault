import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | replication-action-generate-token', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders with the expected elements', async function(assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      {{replication-action-generate-token}}
    `);
    assert.dom('h4.title').hasText('Generate operation token', 'renders default title');
    assert.dom('[data-test-replication-action-trigger]').hasText('Generate token', 'renders default CTA');
  });
});
