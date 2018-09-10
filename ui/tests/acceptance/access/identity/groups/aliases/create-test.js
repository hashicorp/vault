import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { testAliasCRUD, testAliasDeleteFromForm } from '../../_shared-alias-tests';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | /access/identity/groups/aliases/add', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it allows create, list, delete of an entity alias', async function(assert) {
    let name = `alias-${Date.now()}`;
    await testAliasCRUD(name, 'groups', assert);
  });

  test('it allows delete from the edit form', async function(assert) {
    let name = `alias-${Date.now()}`;
    await testAliasDeleteFromForm(name, 'groups', assert);
  });
});
