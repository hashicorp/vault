import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { testAliasCRUD, testAliasDeleteFromForm } from '../../_shared-alias-tests';

module('Acceptance | /access/identity/groups/aliases/add', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it allows create, list, delete of an entity alias', function(assert) {
    let name = `alias-${Date.now()}`;
    testAliasCRUD(name, 'groups', assert);
  });

  test('it allows delete from the edit form', function(assert) {
    let name = `alias-${Date.now()}`;
    testAliasDeleteFromForm(name, 'groups', assert);
  });
});
