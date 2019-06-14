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

    assert.equal(
      find('.page-header .title').textContent.trim(),
      'Vault Enterprise',
      'renders default page title'
    );
    assert.equal(
      find('[data-test-empty-state-title]').textContent.trim(),
      'Upgrade to use this feature',
      'renders default title'
    );
    assert.equal(
      find('[data-test-empty-state-message]').textContent.trim(),
      'You will need Vault Enterprise with this feature included to use this feature.',
      'renders default message'
    );
    assert.equal(findAll('[data-test-upgrade-link]').length, 1, 'renders upgrade link');
  });

  test('it renders with custom attributes', async function(assert) {
    await render(hbs`
      {{upgrade-page title="Test Feature Title" featureName="Specific Feature Name" minimumEdition="Vault Enterprise Premium"}}
      <div id="modal-wormhole"></div>
    `);

    assert.equal(
      find('.page-header .title').textContent.trim(),
      'Test Feature Title',
      'renders custom page title'
    );
    assert.equal(
      find('[data-test-empty-state-title]').textContent.trim(),
      'Upgrade to use Specific Feature Name',
      'renders custom title'
    );
    assert.equal(
      find('[data-test-empty-state-message]').textContent.trim(),
      'You will need Vault Enterprise Premium with Specific Feature Name included to use this feature.',
      'renders custom message'
    );
  });
});
