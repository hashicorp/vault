import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import Pretender from 'pretender';

moduleForAcceptance('Acceptance | secrets/secret/create', {
  beforeEach() {
    this.server = new Pretender(function() {
      this.post('/v1/**', this.passthrough);
      this.put('/v1/**', this.passthrough);
      this.get('/v1/**', this.passthrough);
      this.delete('/v1/**', this.passthrough);
    });
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
      JSON.parse(capabilitiesReq.requestBody).path,
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
  mountSecrets.visit().path(enginePath).type('kv').version(1).submit();

  listPage.create();
  editPage.createSecret(secretPath, 'foo', 'bar');
  andThen(() => {
    let capabilitiesReq = this.server.passthroughRequests.findBy('url', '/v1/sys/capabilities-self');
    assert.equal(
      JSON.parse(capabilitiesReq.requestBody).path,
      `${enginePath}/${secretPath}`,
      'calls capabilites with the correct path'
    );
  });
});
