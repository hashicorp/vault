import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import listPage from 'vault/tests/pages/secrets/backend/list';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import Pretender from 'pretender';

moduleForAcceptance('Acceptance | secrets/secret/', {
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

test('it performs the capabilities lookup properly', function(assert) {
  listPage.visitRoot({ backend: 'secret' });

  andThen(() => {
    let capabilitiesReq = this.server.passthroughRequests.findBy('url', '/v1/sys/capabilities-self');
    assert.equal(
      JSON.parse(capabilitiesReq.requestBody).path,
      'secret/data/',
      'calls capabilites with the correct path'
    );
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');
    assert.ok(listPage.createIsPresent, 'shows the create button');
  });
});

test('it hides create if create is not in the capabilities', function(assert) {
  this.server.post('/v1/sys/capabilities-self', () => {
    return [
      200,
      { 'Content-Type': 'application/json' },
      JSON.stringify({
        capabilities: ['read'],
      }),
    ];
  });
  listPage.visitRoot({ backend: 'secret' });

  andThen(() => {
    assert.notOk(listPage.createIsPresent, 'does not show the create button');
  });
});

test('version 1 performs the correct capabilities lookup', function(assert) {
  const path = `kv-${new Date().getTime()}`;
  mountSecrets.visit().path(path).type('kv').version(1).submit();

  andThen(() => {
    let capabilitiesReq = this.server.passthroughRequests.findBy('url', '/v1/sys/capabilities-self');
    assert.equal(
      JSON.parse(capabilitiesReq.requestBody).path,
      path + '/',
      'calls capabilites with the correct path'
    );
    assert.equal(currentRouteName(), 'vault.cluster.secrets.backend.list-root', 'navigates to the list page');
    assert.ok(listPage.createIsPresent, 'shows the create button');
  });
});
