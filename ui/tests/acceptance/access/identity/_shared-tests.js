import page from 'vault/tests/pages/access/identity/create';
import showPage from 'vault/tests/pages/access/identity/show';
import indexPage from 'vault/tests/pages/access/identity/index';

export const testCRUD = (name, itemType, assert) => {
  page.visit({ item_type: itemType });
  page.editForm.name(name).submit();
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.show',
      `${itemType}: navigates to show on create`
    );
    assert.ok(
      showPage.flashMessage.latestMessage.startsWith(
        'Successfully saved',
        `${itemType}: shows a flash message`
      )
    );
    assert.ok(showPage.nameContains(name), `${itemType}: renders the name on the show page`);
  });

  indexPage.visit({ item_type: itemType });
  andThen(() => {
    assert.equal(
      indexPage.items.filterBy('name', name).length,
      1,
      `${itemType}: lists the entity in the entity list`
    );
    indexPage.items.filterBy('name', name)[0].menu();
  });
  indexPage.delete().confirmDelete();

  andThen(() => {
    assert.equal(indexPage.items.filterBy('name', name).length, 0, `${itemType}: the row is deleted`);
    indexPage.flashMessage.latestMessage.startsWith(
      'Successfully deleted',
      `${itemType}: shows flash message`
    );
  });
};

export const testDeleteFromForm = (name, itemType, assert) => {
  page.visit({ item_type: itemType });
  page.editForm.name(name).submit();
  showPage.edit();
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.edit',
      `${itemType}: navigates to edit on create`
    );
  });
  page.editForm.delete().confirmDelete();
  andThen(() => {
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
    indexPage.flashMessage.latestMessage.startsWith(
      'Successfully deleted',
      `${itemType}: shows flash message`
    );
  });
};
