import { module, test } from 'qunit';
import { settled } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { testAliasCRUD, testAliasDeleteFromForm } from '../../_shared-alias-tests';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | /access/identity/entities/aliases/add', function (hooks) {
  // TODO come back and figure out why this is failing.  Seems to be a race condition
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it allows create, list, delete of an entity alias', async function (assert) {
    assert.expect(6);
    const name = `alias-${Date.now()}`;
    await testAliasCRUD(name, 'entities', assert);
    await settled();
  });

  test('it allows delete from the edit form', async function (assert) {
    assert.expect(4);
    const name = `alias-${Date.now()}`;
    await testAliasDeleteFromForm(name, 'entities', assert);
    await settled();
  });
});
