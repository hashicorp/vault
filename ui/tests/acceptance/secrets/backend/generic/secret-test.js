import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import editPage from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import showPage from 'vault/tests/pages/secrets/backend/kv/show';
import listPage from 'vault/tests/pages/secrets/backend/list';
import consolePanel from 'vault/tests/pages/components/console/ui-panel';

import { create } from 'ember-cli-page-object';

import apiStub from 'vault/tests/helpers/noop-all-api-requests';

const cli = create(consolePanel);

moduleForAcceptance('Acceptance | secrets/generic/create', {
  beforeEach() {
    this.server = apiStub({ usePassthrough: true });
    return authLogin();
  },
  afterEach() {
    this.server.shutdown();
  },
});

test('it creates and can view a secret with the generic backend', function(assert) {
  const path = `generic-${new Date().getTime()}`;
  const kvPath = `generic-kv-${new Date().getTime()}`;
  cli.consoleInput(`write sys/mounts/${path} type=generic`);
  cli.enter();
  cli.consoleInput(`write ${path}/foo bar=baz`);
  cli.enter();
  listPage.visitRoot({ backend: path });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');
    assert.equal(listPage.secrets.length, 1, 'lists one secret in the backend');
  });

  listPage.create();
  editPage.createSecret(kvPath, 'foo', 'bar');
  andThen(() => {
    let capabilitiesReq = this.server.passthroughRequests.findBy('url', '/v1/sys/capabilities-self');
    assert.equal(
      JSON.parse(capabilitiesReq.requestBody).paths,
      `${path}/${kvPath}`,
      'calls capabilites with the correct path'
    );
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.show', 'redirects to the show page');
    assert.ok(showPage.editIsPresent, 'shows the edit button');
  });
});

test('upgrading generic to version 2 lists all existing secrets, and CRUD continues to work', function(
  assert
) {
  const path = `generic-${new Date().getTime()}`;
  const kvPath = `generic-kv-${new Date().getTime()}`;
  cli.consoleInput(`write sys/mounts/${path} type=generic`);
  cli.enter();
  cli.consoleInput(`write ${path}/foo bar=baz`);
  cli.enter();
  // upgrade to version 2 generic mount
  cli.consoleInput(`write sys/mounts/${path}/tune options=version=2`);
  cli.enter();
  listPage.visitRoot({ backend: path });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');
    assert.equal(listPage.secrets.length, 1, 'lists the old secret in the backend');
  });

  listPage.create();
  editPage.createSecret(kvPath, 'foo', 'bar');
  andThen(() => {
    let capabilitiesReq = this.server.passthroughRequests.findBy('url', '/v1/sys/capabilities-self');
    assert.equal(
      JSON.parse(capabilitiesReq.requestBody).paths,
      `${path}/data/${kvPath}`,
      'calls capabilites with the correct path'
    );
  });
  listPage.visitRoot({ backend: path });
  andThen(() => {
    assert.equal(listPage.secrets.length, 2, 'lists two secrets in the backend');
  });
});
