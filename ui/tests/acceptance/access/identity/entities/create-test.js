import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/identity/create';
import { testCRUD, testDeleteFromForm } from '../_shared-tests';

module('Acceptance | /access/identity/entities/create', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it visits the correct page', function(assert) {
    page.visit({ item_type: 'entities' });
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.create',
      'navigates to the correct route'
    );
  });

  test('it allows create, list, delete of an entity', function(assert) {
    let name = `entity-${Date.now()}`;
    testCRUD(name, 'entities', assert);
  });

  test('it can be deleted from the edit form', function(assert) {
    let name = `entity-${Date.now()}`;
    testDeleteFromForm(name, 'entities', assert);
  });
});
