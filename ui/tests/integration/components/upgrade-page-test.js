import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, findAll, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | upgrade page', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders with defaults', async function(assert) {
    await render(hbs`
      {{upgrade-page}}
      <div id="modal-wormhole"></div>
    `);

    assert.equal(find('.page-header .title').textContent.trim(), 'Vault Enterprise', 'renders default title');
    assert.equal(
      find('[data-test-upgrade-feature-description]').textContent.trim(),
      'This is a Vault Enterprise feature.',
      'renders default description'
    );
    assert.equal(findAll('[data-test-upgrade-link]').length, 1, 'renders upgrade link');
  });

  test('it renders with custom attributes', async function(assert) {
    await render(hbs`
      {{upgrade-page title="Test Feature Title" featureName="Specific Feature Name" minimumEdition="Premium"}}
      <div id="modal-wormhole"></div>
    `);

    assert.equal(
      find('.page-header .title').textContent.trim(),
      'Test Feature Title',
      'renders default title'
    );
    assert.equal(
      find('[data-test-upgrade-feature-description]').textContent.trim(),
      'Specific Feature Name is a Premium feature.',
      'renders default description'
    );
  });
});
