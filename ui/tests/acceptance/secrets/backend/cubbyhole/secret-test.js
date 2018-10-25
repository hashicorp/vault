import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | secrets/cubbyhole/create', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub({ usePassthrough: true });
    return authPage.login();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('it creates and can view a secret with the cubbyhole backend', async function(assert) {
    const kvPath = `cubbyhole-kv-${new Date().getTime()}`;
    await listPage.visitRoot({ backend: 'cubbyhole' });
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');

    await listPage.create();
    await editPage.createSecret(kvPath, 'foo', 'bar');
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });
});
