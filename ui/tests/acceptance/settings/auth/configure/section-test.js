import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/section';
import indexPage from 'vault/tests/pages/settings/auth/configure/index';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import consolePanel from 'vault/tests/pages/components/console/ui-panel';
import authPage from 'vault/tests/pages/auth';
import withFlash from 'vault/tests/helpers/with-flash';

const cli = create(consolePanel);

module('Acceptance | settings/auth/configure/section', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub({ usePassthrough: true });
    return authPage.login();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('it can save options', async function(assert) {
    const path = `approle-${new Date().getTime()}`;
    const type = 'approle';
    const section = 'options';
    await enablePage.enable(type, path);
    await page.visit({ path, section });
    await page.fillInTextarea('description', 'This is AppRole!');
    await withFlash(page.save(), () => {
      assert.equal(
        page.flash.latestMessage,
        `The configuration was saved successfully.`,
        'success flash shows'
      );
    });
    let tuneRequest = this.server.passthroughRequests.filterBy('url', `/v1/sys/mounts/auth/${path}/tune`)[0];
    let keys = Object.keys(JSON.parse(tuneRequest.requestBody));
    assert.ok(keys.includes('default_lease_ttl'), 'passes default_lease_ttl on tune');
    assert.ok(keys.includes('max_lease_ttl'), 'passes max_lease_ttl on tune');
  });

  for (let type of ['aws', 'azure', 'gcp', 'github', 'kubernetes', 'ldap', 'okta', 'radius']) {
    test(`it shows tabs for auth method: ${type}`, async assert => {
      let path = `${type}-${Date.now()}`;
      await cli.consoleInput(`write sys/auth/${path} type=${type}`);
      await cli.enter();
      await indexPage.visit({ path });
      // aws has 4 tabs, the others will have 'Configuration' and 'Method Options' tabs
      let numTabs = type === 'aws' ? 4 : 2;
      assert.equal(page.tabs.length, numTabs, 'shows correct number of tabs');
    });
  }
});
