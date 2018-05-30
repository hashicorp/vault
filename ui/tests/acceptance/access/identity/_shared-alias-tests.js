import page from 'vault/tests/pages/access/identity/aliases/add';
import aliasIndexPage from 'vault/tests/pages/access/identity/aliases/index';
import aliasShowPage from 'vault/tests/pages/access/identity/aliases/show';
import createItemPage from 'vault/tests/pages/access/identity/create';
import showItemPage from 'vault/tests/pages/access/identity/show';

export const testAliasCRUD = (name, itemType, assert) => {
  let itemID;
  let aliasID;
  if (itemType === 'groups') {
    createItemPage.createItem(itemType, 'external');
  } else {
    createItemPage.createItem(itemType);
  }
  andThen(() => {
    let idRow = showItemPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
    itemID = idRow.rowValue;
    page.visit({ item_type: itemType, id: itemID });
  });
  page.editForm.name(name).submit();
  andThen(() => {
    let idRow = aliasShowPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
    aliasID = idRow.rowValue;
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.aliases.show',
      'navigates to the correct route'
    );
    assert.ok(
      aliasShowPage.flashMessage.latestMessage.startsWith(
        'Successfully saved',
        `${itemType}: shows a flash message`
      )
    );
    assert.ok(aliasShowPage.nameContains(name), `${itemType}: renders the name on the show page`);
  });

  aliasIndexPage.visit({ item_type: itemType });
  andThen(() => {
    assert.equal(
      aliasIndexPage.items.filterBy('name', name).length,
      1,
      `${itemType}: lists the entity in the entity list`
    );
    aliasIndexPage.items.filterBy('name', name)[0].menu();
  });
  aliasIndexPage.delete().confirmDelete();

  andThen(() => {
    assert.equal(aliasIndexPage.items.filterBy('id', aliasID).length, 0, `${itemType}: the row is deleted`);
    aliasIndexPage.flashMessage.latestMessage.startsWith(
      'Successfully deleted',
      `${itemType}: shows flash message`
    );
  });
};

export const testAliasDeleteFromForm = (name, itemType, assert) => {
  let itemID;
  let aliasID;
  if (itemType === 'groups') {
    createItemPage.createItem(itemType, 'external');
  } else {
    createItemPage.createItem(itemType);
  }
  andThen(() => {
    let idRow = showItemPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
    itemID = idRow.rowValue;
    page.visit({ item_type: itemType, id: itemID });
  });
  page.editForm.name(name).submit();
  andThen(() => {
    let idRow = aliasShowPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
    aliasID = idRow.rowValue;
  });
  aliasShowPage.edit();

  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.aliases.edit',
      `${itemType}: navigates to edit on create`
    );
  });
  page.editForm.delete().confirmDelete();
  andThen(() => {
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
    aliasIndexPage.flashMessage.latestMessage.startsWith(
      'Successfully deleted',
      `${itemType}: shows flash message`
    );
  });
};
