import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('upgrade-page', 'Integration | Component | upgrade page', {
  integration: true,
});

test('it renders with defaults', function(assert) {
  this.render(hbs`
    {{upgrade-page}}
    <div id="modal-wormhole"></div>
  `);

  assert.equal(this.$('.page-header .title').text().trim(), 'Vault Enterprise', 'renders default title');
  assert.equal(
    this.$('[data-test-upgrade-feature-description]').text().trim(),
    'This is a Vault Enterprise feature.',
    'renders default description'
  );
  assert.equal(this.$('[data-test-upgrade-link]').length, 1, 'renders upgrade link');
});

test('it renders with custom attributes', function(assert) {
  this.render(hbs`
    {{upgrade-page title="Test Feature Title" featureName="Specific Feature Name" minimumEdition="Premium"}}
    <div id="modal-wormhole"></div>
  `);

  assert.equal(this.$('.page-header .title').text().trim(), 'Test Feature Title', 'renders default title');
  assert.equal(
    this.$('[data-test-upgrade-feature-description]').text().trim(),
    'Specific Feature Name is a Premium feature.',
    'renders default description'
  );
});
