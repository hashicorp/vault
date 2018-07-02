import { test } from 'qunit';
import { create } from 'ember-cli-page-object';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/section';
import indexPage from 'vault/tests/pages/settings/auth/configure/index';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import consolePanel from 'vault/tests/pages/components/console/ui-panel';

const cli = create(consolePanel);

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

['aws', 'azure', 'gcp', 'github', 'kubernetes', 'ldap', 'okta', 'radius'].forEach(function(type) {
  test(`it shows tabs for auth method: ${type}`, assert => {
    let path = `${type}-${Date.now()}`;
    cli.consoleInput(`write sys/auth/${path} type=${type}`);
    cli.enter();
    indexPage.visit({ path });
    andThen(() => {
      // aws has 4 tabs, the others will have 'Configuration' and 'Method Options' tabs
      let numTabs = type === 'aws' ? 4 : 2;
      assert.equal(page.tabs.length, numTabs, 'shows correct number of tabs');
    });
  });
});
