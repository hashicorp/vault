import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

moduleForAcceptance('Acceptance | secrets/cubbyhole/create', {
  beforeEach() {
    this.server = apiStub({ usePassthrough: true });
    return authLogin();
  },
  afterEach() {
    this.server.shutdown();
  },
});

test('it creates and can view a secret with the cubbyhole backend', function(assert) {
  const kvPath = `cubbyhole-kv-${new Date().getTime()}`;
  listPage.visitRoot({ backend: 'cubbyhole' });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');
  });

  listPage.create();
  editPage.createSecret(kvPath, 'foo', 'bar');
  andThen(() => {
    let capabilitiesReq = this.server.passthroughRequests.findBy('url', '/v1/sys/capabilities-self');
    assert.equal(
      JSON.parse(capabilitiesReq.requestBody).paths,
      `cubbyhole/${kvPath}`,
      'calls capabilites with the correct path'
    );
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });
});
