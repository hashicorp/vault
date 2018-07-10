import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import { testAliasCRUD, testAliasDeleteFromForm } from '../../_shared-alias-tests';

moduleForAcceptance('Acceptance | /access/identity/groups/aliases/add', {
  beforeEach() {
    return authLogin();
  },
});

test('it allows create, list, delete of an entity alias', function(assert) {
  let name = `alias-${Date.now()}`;
  testAliasCRUD(name, 'groups', assert);
});

test('it allows delete from the edit form', function(assert) {
  let name = `alias-${Date.now()}`;
  testAliasDeleteFromForm(name, 'groups', assert);
});
