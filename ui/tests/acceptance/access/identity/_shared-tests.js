import { settled, currentRouteName, click, waitUntil, find } from '@ember/test-helpers';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';
import page from 'vault/tests/pages/access/identity/create';
import showPage from 'vault/tests/pages/access/identity/show';
import indexPage from 'vault/tests/pages/access/identity/index';

export const testCRUD = async (name, itemType, assert) => {
  await page.visit({ item_type: itemType });
  await settled();
  await page.editForm.name(name).submit();
  await settled();
  assert.ok(
    showPage.flashMessage.latestMessage.startsWith('Successfully saved'),
    `${itemType}: shows a flash message`
  );
  assert.strictEqual(
    currentRouteName(),
    'vault.cluster.access.identity.show',
    `${itemType}: navigates to show on create`
  );
  assert.ok(showPage.nameContains(name), `${itemType}: renders the name on the show page`);

  await indexPage.visit({ item_type: itemType });
  await settled();
  assert.strictEqual(
    indexPage.items.filterBy('name', name).length,
    1,
    `${itemType}: lists the entity in the entity list`
  );
  await indexPage.items.filterBy('name', name)[0].menu();
  await waitUntil(() => find('[data-test-item-delete]'));
  await indexPage.delete();
  await settled();
  await indexPage.confirmDelete();
  await settled();
  assert.ok(
    indexPage.flashMessage.latestMessage.startsWith('Successfully deleted'),
    `${itemType}: shows flash message`
  );
  assert.strictEqual(indexPage.items.filterBy('name', name).length, 0, `${itemType}: the row is deleted`);
};

export const testDeleteFromForm = async (name, itemType, assert) => {
  await page.visit({ item_type: itemType });
  await settled();
  await page.editForm.name(name);
  await page.editForm.metadataKey('hello');
  await page.editForm.metadataValue('goodbye');
  await clickTrigger('#policies');
  // first option should be "default"
  await selectChoose('#policies', '.ember-power-select-option', 0);
  await page.editForm.submit();
  await click('[data-test-tab-subnav="policies"]');
  assert.dom('.list-item-row').exists({ count: 1 }, 'One item is under policies');
  await click('[data-test-tab-subnav="metadata"]');
  assert.dom('.info-table-row').hasText('hello goodbye', 'Metadata shows on tab');
  await showPage.edit();
  assert.strictEqual(
    currentRouteName(),
    'vault.cluster.access.identity.edit',
    `${itemType}: navigates to edit on create`
  );
  await settled();
  await page.editForm.delete();
  await settled();
  await page.editForm.confirmDelete();
  await settled();
  assert.ok(
    indexPage.flashMessage.latestMessage.startsWith('Successfully deleted'),
    `${itemType}: shows flash message`
  );
  assert.strictEqual(
    currentRouteName(),
    'vault.cluster.access.identity.index',
    `${itemType}: navigates to list page on delete`
  );
  assert.strictEqual(
    indexPage.items.filterBy('name', name).length,
    0,
    `${itemType}: the row does not show in the list`
  );
};
