import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-keys';

module('Integration | Component | pki key details page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:flash-messages').registerTypes(['success', 'danger']);
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
    this.store.pushPayload('pki/key', {
      modelName: 'pki/key',
      key_id: '724862ff-6438-bad0-b598-77a6c7f4e934',
      key_type: 'ec',
      key_name: 'test-key',
    });
    this.model = this.store.peekRecord('pki/key', '724862ff-6438-bad0-b598-77a6c7f4e934');
  });

  test('it renders the page component and deletes a key', async function (assert) {
    assert.expect(7);
    this.server.delete(`${this.backend}/key/${this.model.keyId}`, () => {
      assert.ok(true, 'confirming delete fires off destroyRecord()');
    });

    await render(
      hbs`
        <Page::PkiKeyDetails
          @key={{this.model}} 
          @canDelete={{true}}
          @canEdit={{true}} 
        />
      `,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.keyIdValue).hasText(' 724862ff-6438-bad0-b598-77a6c7f4e934', 'key id renders');
    assert.dom(SELECTORS.keyNameValue).hasText('test-key', 'key name renders');
    assert.dom(SELECTORS.keyTypeValue).hasText('ec', 'key type renders');
    assert.dom(SELECTORS.keyBitsValue).doesNotExist('does not render empty value');
    assert.dom(SELECTORS.keyEditLink).exists('renders edit link');
    assert.dom(SELECTORS.keyDeleteButton).exists('renders delete button');
    await click(SELECTORS.keyDeleteButton);
    await click(SELECTORS.confirmDelete);
  });

  test('it does not render actions when capabilities are false', async function (assert) {
    assert.expect(2);

    await render(
      hbs`
        <Page::PkiKeyDetails
          @key={{this.model}} 
          @canDelete={{false}}
          @canEdit={{false}} 
        />
      `,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.keyDeleteButton).doesNotExist('does not render delete button if no permission');
    assert.dom(SELECTORS.keyEditLink).doesNotExist('does not render edit button if no permission');
  });
});
