import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | upgrade page', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders with defaults', async function (assert) {
    await render(hbs`
      {{upgrade-page}}
      <div id="modal-wormhole"></div>
    `);

    assert.dom('.page-header .title').hasText('Vault Enterprise', 'renders default page title');
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('Upgrade to use this feature', 'renders default title');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'You will need Vault Enterprise with this feature included to use this feature.',
        'renders default message'
      );
    assert.dom('[data-test-upgrade-link]').exists({ count: 1 }, 'renders upgrade link');
  });

  test('it renders with custom attributes', async function (assert) {
    await render(hbs`
      {{upgrade-page title="Test Feature Title" minimumEdition="Vault Enterprise Premium"}}
      <div id="modal-wormhole"></div>
    `);

    assert.dom('.page-header .title').hasText('Test Feature Title', 'renders custom page title');
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('Upgrade to use Test Feature Title', 'renders custom title');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'You will need Vault Enterprise Premium with Test Feature Title included to use this feature.',
        'renders custom message'
      );
  });
});
