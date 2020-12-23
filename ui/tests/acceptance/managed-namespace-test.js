import { module, test } from 'qunit';
import { currentURL, visit, fillIn } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import setupMirage from 'ember-cli-mirage/test-support/setup-mirage';

module('Acceptance | Enterprise | Managed namespace root', function(hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function() {
    server.logging = true;
    /**
     * Since the features are fetched on the application load,
     * we have to populate them on the beforeEach hook because
     * the fetch won't trigger again within the tests
     */
    server.create('feature', { feature_flags: ['VAULT_CLOUD_ADMIN_NAMESPACE'] });
  });

  test('it shows the managed namespace toolbar when feature flag exists', async function(assert) {
    await visit('/vault/auth');
    assert.equal(currentURL(), '/vault/auth?namespace=admin&with=token', 'Redirected to base namespace');

    assert.dom('[data-test-namespace-toolbar]').doesNotExist('Normal namespace toolbar does not exist');
    assert.dom('[data-test-managed-namespace-toolbar]').exists('Managed namespace toolbar exists');
    assert.dom('[data-test-managed-namespace-root]').hasText('/admin', 'Shows /admin namespace prefix');
    assert.dom('input#namespace').hasAttribute('placeholder', '/ (Default)');
    await fillIn('input#namespace', '/foo');
    let encodedNamespace = encodeURIComponent('admin/foo');
    assert.equal(
      currentURL(),
      `/vault/auth?namespace=${encodedNamespace}&with=token`,
      'Correctly prepends root to namespace'
    );
  });
});
