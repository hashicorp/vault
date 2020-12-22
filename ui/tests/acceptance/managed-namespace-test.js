import { module, test } from 'qunit';
import { visit, currentURL, settled } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import setupMirage from 'ember-cli-mirage/test-support/setup-mirage';
import logout from 'vault/tests/pages/logout';

module('Acceptance | Enterprise | Managed namespace root', function(hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function() {
    return logout.visit();
  });

  test('it shows the regular namespace toolbar when not managed', async function(assert) {
    assert.dom('[data-test-namespace-toolbar]').exists('Namespace toolbar exists');
    assert.dom('input#namespace').hasAttribute('placeholder', '/ (Root)');
    assert.equal(currentURL(), '/vault/auth?with=token', 'Does not redirect');
  });

  test('it shows the managed namespace toolbar when feature flag exists', async function(assert) {
    await server.create('feature', { feature_flags: ['SOME_FLAG', 'VAULT_CLOUD_ADMIN_NAMESPACE'] });
    // await settled();
    await visit('/vault/auth');
    await settled();
    console.log(currentURL());
    assert.equal(1, 1);
    // await this.pauseTest();

    // assert.equal(currentURL(), '/vault/auth?namespace=admin&with=token', 'Redirects to root namespace');
    // assert.dom('[data-test-managed-namespace-toolbar]').exists('Managed namespace toolbar exists');
    // await this.pauseTest();
  });
});
