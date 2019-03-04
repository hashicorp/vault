import { waitFor, settled, currentRouteName } from '@ember/test-helpers';
import page from 'vault/tests/pages/access/identity/create';
import showPage from 'vault/tests/pages/access/identity/show';
import indexPage from 'vault/tests/pages/access/identity/index';

export const testCRUD = async (name, itemType, assert) => {
  await page.visit({ item_type: itemType });
  await page.editForm.name(name).submit();
  assert.ok(
    showPage.flashMessage.latestMessage.startsWith('Successfully saved'),
    `${itemType}: shows a flash message`
  );
  // wait for redirect;
  await settled();
  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.show',
    `${itemType}: navigates to show on create`
  );
  assert.ok(showPage.nameContains(name), `${itemType}: renders the name on the show page`);

  await indexPage.visit({ item_type: itemType });
  assert.equal(
    indexPage.items.filterBy('name', name).length,
    1,
    `${itemType}: lists the entity in the entity list`
  );
  await indexPage.items.filterBy('name', name)[0].menu();
  await waitFor('[data-test-confirm-action-trigger]');

  let btnIndex = itemType === 'groups' ? 0 : 1;
  await indexPage.buttons[btnIndex].delete();
  await indexPage.confirmDelete();
  assert.ok(
    indexPage.flashMessage.latestMessage.startsWith('Successfully deleted'),
    `${itemType}: shows flash message`
  );
  assert.equal(indexPage.items.filterBy('name', name).length, 0, `${itemType}: the row is deleted`);
};

export const testDeleteFromForm = async (name, itemType, assert) => {
  await page.visit({ item_type: itemType });

  await page.editForm.name(name).submit();
  await showPage.edit();
  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.edit',
    `${itemType}: navigates to edit on create`
  );
  await page.editForm.delete();
  await page.editForm.confirmDelete();
  assert.ok(
    indexPage.flashMessage.latestMessage.startsWith('Successfully deleted'),
    `${itemType}: shows flash message`
  );
  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.index',
    `${itemType}: navigates to list page on delete`
  );
  assert.equal(
    indexPage.items.filterBy('name', name).length,
    0,
    `${itemType}: the row does not show in the list`
  );
};
