import { currentURL, currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import authPage from 'vault/tests/pages/auth';
import withFlash from 'vault/tests/helpers/with-flash';

module('Acceptance | secrets/secret/create', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub({ usePassthrough: true });
    return authPage.login();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('it creates a secret and redirects', async function(assert) {
    const path = `kv-path-${new Date().getTime()}`;
    await listPage.visitRoot({ backend: 'secret' });
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');

    await listPage.create();
    await editPage.createSecret(path, 'foo', 'bar');

    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });

  test('version 1 performs the correct capabilities lookup', async function(assert) {
    let enginePath = `kv-${new Date().getTime()}`;
    let secretPath = 'foo/bar';
    // mount version 1 engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await withFlash(
      mountSecrets
        .next()
        .path(enginePath)
        .version(1)
        .submit()
    );

    await listPage.create();
    await editPage.createSecret(secretPath, 'foo', 'bar');

    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });

  test('it redirects to the path ending in / for list pages', async function(assert) {
    const path = `foo/bar/kv-path-${new Date().getTime()}`;
    await listPage.visitRoot({ backend: 'secret' });
    await listPage.create();
    await editPage.createSecret(path, 'foo', 'bar');
    await listPage.visit({ backend: 'secret', id: 'foo/bar' });
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list');
    assert.ok(currentURL().endsWith('/'), 'redirects to the path ending in a slash');
  });
});
