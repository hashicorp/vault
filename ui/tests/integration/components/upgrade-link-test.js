import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | upgrade link', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders with overlay', async function(assert) {
    await render(hbs`
       <div id="modal-wormhole"></div>
       <div class="upgrade-link-container">
         {{#upgrade-link data-test-link}}upgrade{{/upgrade-link}}
       </div>
    `);

    assert.equal(
      find('.upgrade-link-container button').textContent.trim(),
      'upgrade',
      'renders link content'
    );
    assert.equal(
      find('#modal-wormhole .upgrade-overlay-title').textContent.trim(),
      'Try Vault Enterprise free for 30 days',
      'contains overlay content'
    );
    assert.equal(
      findAll('#modal-wormhole a[href^="https://hashicorp.com/products/vault/trial?source=vaultui"]').length,
      1,
      'contains info link'
    );
  });

  test('it adds custom classes', async function(assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <div class="upgrade-link-container">
        {{#upgrade-link linkClass="button upgrade-button"}}upgrade{{/upgrade-link}}
      </div>
    `);

    assert.equal(
      find('.upgrade-link-container button').getAttribute('class'),
      'link button upgrade-button',
      'adds classes to link'
    );
  });
});
