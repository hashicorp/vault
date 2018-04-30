import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/access/identity/create';
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

test('it creates an entity', function(assert) {
  let name = `entity-${Date.now()}`;
  page.visit({ item_type: 'entities' });
  page.editForm.name(name).submit();
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.access.identity.show', 'navigates to show on create');
    assert.ok(
      indexPage.flashMessage.latestMessage.startsWith('Successfully saved Entity', 'shows a flash message')
    );
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
