import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/access/identity/create';
import { testCRUD, testDeleteFromForm } from '../_shared-tests';

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
  testCRUD(name, 'entities', assert);
});

test('it can be deleted from the edit form', function(assert) {
  let name = `entity-${Date.now()}`;
  testDeleteFromForm(name, 'entities', assert);
});
