import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/page-role-details';

module('Integration | Component | pki/role detail page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role', {
      name: 'Foobar',
      noStore: false,
      keyUsage: [],
      extKeyUsage: ['bar', 'baz'],
    });
    this.model.backend = 'pki';
  });

  test('it should render the page component', async function (assert) {
    assert.expect(7);
    await render(
      hbs`
      <Page::Roles::Role::DetailsPage @role={{this.model}} />
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.breadcrumbContainer).exists({ count: 1 }, 'breadcrumb containers exist');
    assert.dom(SELECTORS.breadcrumbs).exists({ count: 4 }, 'Shows 4 breadcrumbs');
    assert.dom(SELECTORS.title).containsText('PKI Role Foobar', 'Title includes type and name of role');
    // Attribute-specific checks
    assert.dom(SELECTORS.issuerLabel).hasText('Issuer', 'Label is');
    assert.dom(SELECTORS.keyUsageValue).hasText('None', 'Key usage shows none when array is empty');
    assert
      .dom(SELECTORS.extKeyUsageValue)
      .containsText('bar, baz', 'Key usage shows comma-joined values when array has items');
    assert.dom(SELECTORS.noStoreValue).containsText('Yes', 'noStore shows opposite of what the value is');
  });
});
