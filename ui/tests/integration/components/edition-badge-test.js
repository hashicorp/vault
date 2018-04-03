import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('edition-badge', 'Integration | Component | edition badge', {
  integration: true,
});

test('it renders', function(assert) {
  this.render(hbs`
    {{edition-badge edition="Custom"}}
  `);

  assert.equal(this.$('.edition-badge').text().trim(), 'Custom', 'contains edition');

  this.render(hbs`
    {{edition-badge edition="Enterprise"}}
  `);

  assert.equal(this.$('.edition-badge').text().trim(), 'Ent', 'abbreviates Enterprise');
});
