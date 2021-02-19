import { currentRouteName } from '@ember/test-helpers';
import { module, skip } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/identity/create';
import { testCRUD, testDeleteFromForm } from '../_shared-tests';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | /access/identity/entities/create', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  skip('it visits the correct page', async function(assert) {
    await page.visit({ item_type: 'entities' });
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.create',
      'navigates to the correct route'
    );
  });

  skip('it allows create, list, delete of an entity', async function(assert) {
    let name = `entity-${Date.now()}`;
    await testCRUD(name, 'entities', assert);
  });

  skip('it can be deleted from the edit form', async function(assert) {
    let name = `entity-${Date.now()}`;
    await testDeleteFromForm(name, 'entities', assert);
  });
});
