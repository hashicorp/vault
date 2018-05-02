import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/access/identity/create';
import showPage from 'vault/tests/pages/access/identity/show';
import indexPage from 'vault/tests/pages/access/identity/index';

moduleForAcceptance('Acceptance | /access/identity/entities/create', {
  beforeEach() {
    return authLogin();
  },
});

test('it visits the correct page', function(assert) {
  page.visit({ item_type: 'entities' });
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.create',
      'navigates to the correct route'
    );
  });
});

test('it allows create, list, delete of an entity', function(assert) {
  let name = `entity-${Date.now()}`;
  let id;
  page.visit({ item_type: 'entities' });
  page.editForm.name(name).submit();
  andThen(() => {
    let idRow = showPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0];
    id = idRow.rowValue;
    assert.equal(currentRouteName(), 'vault.cluster.access.identity.show', 'navigates to show on create');
    assert.ok(
      showPage.flashMessage.latestMessage.startsWith('Successfully saved Entity', 'shows a flash message')
    );
    assert.ok(showPage.nameContains(name), 'renders the name on the show page');
  });

  indexPage.visit({ item_type: 'entities' });
  andThen(() => {
    assert.equal(indexPage.items.filterBy('id', id).length, 1, 'lists the entity in the entity list');
    indexPage.items.filterBy('id', id)[0].menu();
  });
  indexPage.delete().confirmDelete();

  andThen(() => {
    assert.equal(indexPage.items.filterBy('id', id).length, 0, 'the row is deleted');
    indexPage.flashMessage.latestMessage.startsWith('Successfully deleted', 'shows flash message');
  });
});
test('it can be deleted from the edit form', function(assert) {
  let name = `entity-${Date.now()}`;
  let id;
  page.visit({ item_type: 'entities' });
  page.editForm.name(name).submit();
  andThen(() => {
    id = showPage.rows.filterBy('hasLabel').filterBy('rowLabel', 'ID')[0].rowValue
  });
  showPage.edit();
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.access.identity.edit', 'navigates to edit on create');
  });
  page.editForm.delete().confirmDelete();
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.access.identity.index', 'navigates to list page on delete');
    assert.equal(indexPage.items.filterBy('id', id).length, 0, 'the row does not show in the list');
    indexPage.flashMessage.latestMessage.startsWith('Successfully deleted', 'shows flash message');
  });
});


test('it visits the correct page', function(assert) {
  page.visit({ item_type: 'groups' });
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.create',
      'navigates to the correct route'
    );
  });
});
