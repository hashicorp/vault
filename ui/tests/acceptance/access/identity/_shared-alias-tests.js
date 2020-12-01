import { currentRouteName, settled } from '@ember/test-helpers';
import page from 'vault/tests/pages/access/identity/aliases/add';
import aliasIndexPage from 'vault/tests/pages/access/identity/aliases/index';
import aliasShowPage from 'vault/tests/pages/access/identity/aliases/show';
import createItemPage from 'vault/tests/pages/access/identity/create';
import showItemPage from 'vault/tests/pages/access/identity/show';

export const testAliasCRUD = async function(name, itemType, assert) {
  let itemID;
  let aliasID;
  let idRow;
  if (itemType === 'groups') {
    await createItemPage.createItem(itemType, 'external');
    await settled();
  } else {
    await createItemPage.createItem(itemType);
    await settled();
  }
  idRow = showItemPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
  itemID = idRow.rowValue;
  await page.visit({ item_type: itemType, id: itemID });
  await settled();
  await page.editForm.name(name).submit();
  await settled();
  assert.ok(
    aliasShowPage.flashMessage.latestMessage.startsWith('Successfully saved'),
    `${itemType}: shows a flash message`
  );

  idRow = aliasShowPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
  aliasID = idRow.rowValue;
  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.aliases.show',
    'navigates to the correct route'
  );
  assert.ok(aliasShowPage.nameContains(name), `${itemType}: renders the name on the show page`);

  await aliasIndexPage.visit({ item_type: itemType });
  await settled();
  assert.equal(
    aliasIndexPage.items.filterBy('name', name).length,
    1,
    `${itemType}: lists the entity in the entity list`
  );

  let item = aliasIndexPage.items.filterBy('name', name)[0];
  await item.menu();
  await settled();
  await aliasIndexPage.delete();
  await settled();
  await aliasIndexPage.confirmDelete();
  await settled();
  assert.ok(
    aliasIndexPage.flashMessage.latestMessage.startsWith('Successfully deleted'),
    `${itemType}: shows flash message`
  );

  assert.equal(aliasIndexPage.items.filterBy('id', aliasID).length, 0, `${itemType}: the row is deleted`);
};

export const testAliasDeleteFromForm = async function(name, itemType, assert) {
  let itemID;
  let aliasID;
  let idRow;
  if (itemType === 'groups') {
    await createItemPage.createItem(itemType, 'external');
    await settled();
  } else {
    await createItemPage.createItem(itemType);
    await settled();
  }
  idRow = showItemPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
  itemID = idRow.rowValue;
  await page.visit({ item_type: itemType, id: itemID });
  await settled();
  await page.editForm.name(name).submit();
  await settled();
  idRow = aliasShowPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
  aliasID = idRow.rowValue;
  await aliasShowPage.edit();
  await settled();
  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.aliases.edit',
    `${itemType}: navigates to edit on create`
  );
  await page.editForm.delete();
  await settled();
  await page.editForm.confirmDelete();
  await settled();
  assert.ok(
    aliasIndexPage.flashMessage.latestMessage.startsWith('Successfully deleted'),
    `${itemType}: shows flash message`
  );

  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.aliases.index',
    `${itemType}: navigates to list page on delete`
  );
  assert.equal(
    aliasIndexPage.items.filterBy('id', aliasID).length,
    0,
    `${itemType}: the row does not show in the list`
  );
};
