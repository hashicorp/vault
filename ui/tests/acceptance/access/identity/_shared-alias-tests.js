import { currentRouteName } from '@ember/test-helpers';
import page from 'vault/tests/pages/access/identity/aliases/add';
import aliasIndexPage from 'vault/tests/pages/access/identity/aliases/index';
import aliasShowPage from 'vault/tests/pages/access/identity/aliases/show';
import createItemPage from 'vault/tests/pages/access/identity/create';
import showItemPage from 'vault/tests/pages/access/identity/show';
import withFlash from 'vault/tests/helpers/with-flash';

export const testAliasCRUD = async function(name, itemType, assert) {
  let itemID;
  let aliasID;
  let idRow;
  if (itemType === 'groups') {
    createItemPage.createItem(itemType, 'external');
  } else {
    createItemPage.createItem(itemType);
  }
  await withFlash();
  idRow = showItemPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
  itemID = idRow.rowValue;
  await page.visit({ item_type: itemType, id: itemID });

  await withFlash(page.editForm.name(name).submit(), () => {
    assert.ok(
      aliasShowPage.flashMessage.latestMessage.startsWith(
        'Successfully saved',
        `${itemType}: shows a flash message`
      )
    );
  });

  idRow = aliasShowPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
  aliasID = idRow.rowValue;
  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.aliases.show',
    'navigates to the correct route'
  );
  assert.ok(aliasShowPage.nameContains(name), `${itemType}: renders the name on the show page`);

  await aliasIndexPage.visit({ item_type: itemType });
  assert.equal(
    aliasIndexPage.items.filterBy('name', name).length,
    1,
    `${itemType}: lists the entity in the entity list`
  );

  let item = aliasIndexPage.items.filterBy('name', name)[0];
  await item.menu();

  await aliasIndexPage.delete();
  await withFlash(aliasIndexPage.confirmDelete(), () => {
    aliasIndexPage.flashMessage.latestMessage.startsWith(
      'Successfully deleted',
      `${itemType}: shows flash message`
    );
  });

  assert.equal(aliasIndexPage.items.filterBy('id', aliasID).length, 0, `${itemType}: the row is deleted`);
};

export const testAliasDeleteFromForm = async function(name, itemType, assert) {
  let itemID;
  let aliasID;
  let idRow;
  if (itemType === 'groups') {
    createItemPage.createItem(itemType, 'external');
  } else {
    createItemPage.createItem(itemType);
  }
  await withFlash();
  idRow = showItemPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
  itemID = idRow.rowValue;
  await page.visit({ item_type: itemType, id: itemID });

  await withFlash(page.editForm.name(name).submit());

  idRow = aliasShowPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
  aliasID = idRow.rowValue;
  await aliasShowPage.edit();

  assert.equal(
    currentRouteName(),
    'vault.cluster.access.identity.aliases.edit',
    `${itemType}: navigates to edit on create`
  );
  await page.editForm.delete();

  await withFlash(page.editForm.confirmDelete(), () => {
    aliasIndexPage.flashMessage.latestMessage.startsWith(
      'Successfully deleted',
      `${itemType}: shows flash message`
    );
  });

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
