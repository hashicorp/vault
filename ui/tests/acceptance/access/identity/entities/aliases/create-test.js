import { module, test, skip } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { testAliasCRUD, testAliasDeleteFromForm } from '../../_shared-alias-tests';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | /access/identity/entities/aliases/add', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  skip('it allows create, list, delete of an entity alias', function(assert) {
    let name = `alias-${Date.now()}`;
    return testAliasCRUD(name, 'entities', assert);
  });

  skip('it allows delete from the edit form', function(assert) {
    let name = `alias-${Date.now()}`;
    return testAliasDeleteFromForm(name, 'entities', assert);
  });
});
