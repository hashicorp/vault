import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/identity/create';
import { testCRUD, testDeleteFromForm } from '../_shared-tests';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | /access/identity/groups/create', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it visits the correct page', async function(assert) {
    await page.visit({ item_type: 'groups' });
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.identity.create',
      'navigates to the correct route'
    );
  });

  test('it allows create, list, delete of an group', async function(assert) {
    let name = `group-${Date.now()}`;
    await testCRUD(name, 'groups', assert);
  });

  test('it can be deleted from the group edit form', async function(assert) {
    let name = `group-${Date.now()}`;
    await testDeleteFromForm(name, 'groups', assert);
  });
});
