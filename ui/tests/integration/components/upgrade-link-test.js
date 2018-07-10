import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('upgrade-link', 'Integration | Component | upgrade link', {
  integration: true,
});

test('it renders with overlay', function(assert) {
  this.render(hbs`
     <div id="modal-wormhole"></div>
     <div class="upgrade-link-container">
       {{#upgrade-link data-test-link}}upgrade{{/upgrade-link}}
     </div>
  `);

  assert.equal(this.$('.upgrade-link-container button').text().trim(), 'upgrade', 'renders link content');
  assert.equal(
    this.$('#modal-wormhole .upgrade-overlay-title').text().trim(),
    'Try Vault Enterprise Free for 30 Days',
    'contains overlay content'
  );
  assert.equal(
    this.$('#modal-wormhole a[href^="https://hashicorp.com/products/vault/trial?source=vaultui"]').length,
    1,
    'contains info link'
  );
});

test('it adds custom classes', function(assert) {
  this.render(hbs`
    <div id="modal-wormhole"></div>
    <div class="upgrade-link-container">
      {{#upgrade-link linkClass="button upgrade-button"}}upgrade{{/upgrade-link}}
    </div>
  `);

  assert.equal(
    this.$('.upgrade-link-container button').attr('class'),
    'link button upgrade-button',
    'adds classes to link'
  );
});
