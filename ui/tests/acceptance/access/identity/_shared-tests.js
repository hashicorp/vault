import { currentRouteName } from '@ember/test-helpers';
import page from 'vault/tests/pages/access/identity/create';
import showPage from 'vault/tests/pages/access/identity/show';
import indexPage from 'vault/tests/pages/access/identity/index';

export const testCRUD = async (name, itemType, assert) => {
  await page.visit({ item_type: itemType });
  let save = page.editForm.name(name).submit();

  await showPage.flashMessage.waitForFlash();
  assert.ok(
    showPage.flashMessage.latestMessage.startsWith('Successfully saved', `${itemType}: shows a flash message`)
  );
  await showPage.flashMessage.clickLast();
  await save;
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
  await indexPage.delete();
  let deleted = indexPage.confirmDelete();

  await indexPage.flashMessage.waitForFlash();
  indexPage.flashMessage.latestMessage.startsWith('Successfully deleted', `${itemType}: shows flash message`);
  await indexPage.flashMessage.clickLast();
  await deleted;
  assert.equal(indexPage.items.filterBy('name', name).length, 0, `${itemType}: the row is deleted`);
};

export const testDeleteFromForm = async (name, itemType, assert) => {
  await page.visit({ item_type: itemType });

  let save = page.editForm.name(name).submit();

  await showPage.flashMessage.waitForFlash();
  await showPage.flashMessage.clickLast();
  await save;
  await showPage.edit();
  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.edit',
    `${itemType}: navigates to edit on create`
  );
  await page.editForm.delete();
  let deleted = page.editForm.confirmDelete();
  await indexPage.flashMessage.waitForFlash();
  indexPage.flashMessage.latestMessage.startsWith('Successfully deleted', `${itemType}: shows flash message`);
  await indexPage.flashMessage.clickLast();
  await deleted;
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
