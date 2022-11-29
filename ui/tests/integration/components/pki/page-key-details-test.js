import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS, KEY_SELECTORS } from 'vault/tests/helpers/pki/page-details';

module('Integration | Component | pki key details page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.owner.lookup('service:flash-messages').registerTypes(['success', 'danger']);
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/key', {
      keyId: '724862ff-6438-bad0-b598-77a6c7f4e934',
      keyType: 'ec',
      keyName: 'exported-key-2',
    });
  });

  test('it renders the page component', async function (assert) {
    assert.expect(8);
    await render(
      hbs`
        <Page::PkiKeyDetails @key={{this.model}} />
      `,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.breadcrumbContainer).exists({ count: 1 }, 'breadcrumb containers exist');
    assert.dom(SELECTORS.breadcrumbs).exists({ count: 4 }, 'Shows 4 breadcrumbs');
    assert.dom(KEY_SELECTORS.title).containsText('View key', 'title renders');
    assert.dom(KEY_SELECTORS.keyIdValue).hasText(' 724862ff-6438-bad0-b598-77a6c7f4e934', 'key id renders');
    assert.dom(KEY_SELECTORS.keyNameValue).hasText('exported-key-2', 'key name renders');
    assert.dom(KEY_SELECTORS.keyTypeValue).hasText('ec', 'key type renders');
    assert.dom(KEY_SELECTORS.keyBitsValue).doesNotExist('does not render empty value');
    assert.dom(KEY_SELECTORS.keyDeleteButton).exists('renders delete button');
  });
});
