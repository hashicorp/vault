import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/access/identity/create';
import { testCRUD, testDeleteFromForm } from '../_shared-tests';

moduleForAcceptance('Acceptance | /access/identity/groups/create', {
  beforeEach() {
    return authLogin();
  },
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

test('it allows create, list, delete of an group', function(assert) {
  let name = `group-${Date.now()}`;
  testCRUD(name, 'groups', assert);
});

test('it can be deleted from the group edit form', function(assert) {
  let name = `group-${Date.now()}`;
  testDeleteFromForm(name, 'groups', assert);
});
