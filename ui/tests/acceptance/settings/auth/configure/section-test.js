import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/section';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

moduleForAcceptance('Acceptance | settings/auth/configure/section', {
  beforeEach() {
    this.server = apiStub({ usePassthrough: true });
    return authLogin();
  },
  afterEach() {
    this.server.shutdown();
  },
});

test('it can save options', function(assert) {
  const path = `approle-${new Date().getTime()}`;
  const type = 'approle';
  const section = 'options';
  enablePage.visit().enableAuth(type, path);
  page.visit({ path, section });
  andThen(() => {
    page.fields().findByName('description').textarea('This is AppRole!');
    page.save();
  });
  andThen(() => {
    let tuneRequest = this.server.passthroughRequests.filterBy('url', `/v1/sys/mounts/auth/${path}/tune`)[0];
    let keys = Object.keys(JSON.parse(tuneRequest.requestBody));
    assert.ok(keys.includes('default_lease_ttl'), 'passes default_lease_ttl on tune');
    assert.ok(keys.includes('max_lease_ttl'), 'passes max_lease_ttl on tune');

    assert.equal(
      page.flash.latestMessage,
      `The configuration options were saved successfully.`,
      'success flash shows'
    );
  });
});
