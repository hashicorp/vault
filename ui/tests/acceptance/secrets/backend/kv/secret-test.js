import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

moduleForAcceptance('Acceptance | secrets/secret/create', {
  beforeEach() {
    this.server = apiStub({ usePassthrough: true });
    return authLogin();
  },
  afterEach() {
    this.server.shutdown();
  },
});

test('it creates a secret and redirects', function(assert) {
  const path = `kv-path-${new Date().getTime()}`;
  listPage.visitRoot({ backend: 'secret' });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');
  });

  listPage.create();
  editPage.createSecret(path, 'foo', 'bar');
  andThen(() => {
    let capabilitiesReq = this.server.passthroughRequests.findBy('url', '/v1/sys/capabilities-self');
    assert.equal(
      JSON.parse(capabilitiesReq.requestBody).paths,
      `secret/data/${path}`,
      'calls capabilites with the correct path'
    );
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });
});

test('version 1 performs the correct capabilities lookup', function(assert) {
  let enginePath = `kv-${new Date().getTime()}`;
  let secretPath = 'foo/bar';
  // mount version 1 engine
  mountSecrets.visit().path(enginePath).type('kv').version(1).submit();

  listPage.create();
  editPage.createSecret(secretPath, 'foo', 'bar');
  andThen(() => {
    let capabilitiesReq = this.server.passthroughRequests.findBy('url', '/v1/sys/capabilities-self');
    assert.equal(
      JSON.parse(capabilitiesReq.requestBody).paths,
      `${enginePath}/${secretPath}`,
      'calls capabilites with the correct path'
    );
  });
});
