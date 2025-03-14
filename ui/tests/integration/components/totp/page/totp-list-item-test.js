import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | totp-list-item', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    const item = EmberObject.create({
      id: 'foo',
      canRead: true,
      canDelete: true,
    });

    this.set('itemPath', 'key/foo');
    this.set('itemType', 'key');
    this.set('item', item);
    await render(hbs`<SecretList::TotpListItem
          @item={{this.item}}
        />`);
    assert.dom(`[data-test-secret-link=${item.id}`).exists('has correct link');
  });

  test('it has details and delete menu item', async function (assert) {
    const item = EmberObject.create({
      id: 'foo',
      canRead: true,
      canDelete: true,
    });

    this.set('itemPath', 'key/foo');
    this.set('itemType', 'key');
    this.set('item', item);

    await render(hbs`<SecretList::TotpListItem
      @item={{this.item}}
    />`);
    await click('[data-test-popup-menu-trigger]');

    assert.dom('.hds-dropdown li').exists({ count: 2 }, 'has both options');
    assert.dom('.hds-dropdown li:nth-of-type(1)').hasText('Details', 'first list item is "details"');
    assert.dom('.hds-dropdown li:nth-of-type(2)').hasText('Delete', 'second list item is "delete"');
  });
});
