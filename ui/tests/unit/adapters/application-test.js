import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

module('Unit | Adapter | application', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('ajax list call - reroutes to filter-list', async function(assert) {
    let adapter = this.owner.lookup('adapter:application');
    let url = '/v1/secret/metadata';
    let expectedRerouteUrl = '/v1/sys/internal/ui/filtered-path/secret/metadata?list=true';
    await adapter.ajax(url, 'GET', { data: { list: true } });
    assert.equal(
      this.server.handledRequests[0].url,
      expectedRerouteUrl,
      'redirects the request to the new url'
    );
  });

  test('ajax list call - retries when the endpoint 403s', async function(assert) {
    this.server.get('/v1/sys/internal/ui/filtered-path/secret/metadata', () => {
      return [403, { 'Content-Type': 'application/json' }, JSON.stringify({})];
    });

    let adapter = this.owner.lookup('adapter:application');
    let url = '/v1/secret/metadata';
    let expectedRerouteUrl = '/v1/sys/internal/ui/filtered-path/secret/metadata?list=true';
    let expectedRetryUrl = '/v1/secret/metadata?list=true';
    await adapter.ajax(url, 'GET', { data: { list: true } });
    assert.equal(
      this.server.handledRequests[0].url,
      expectedRerouteUrl,
      'redirects the request to the new url'
    );
    assert.equal(this.server.handledRequests[1].url, expectedRetryUrl, 'retries on the original url');
  });
});
