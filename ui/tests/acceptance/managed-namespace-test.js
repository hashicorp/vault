/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { currentURL, visit, fillIn } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import Pretender from 'pretender';
import { getManagedNamespace } from 'vault/routes/vault/cluster';

const FEATURE_FLAGS_RESPONSE = {
  feature_flags: ['VAULT_CLOUD_ADMIN_NAMESPACE'],
};

module('Acceptance | Enterprise | Managed namespace root', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    /**
     * Since the features are fetched on the application load,
     * we have to populate them on the beforeEach hook because
     * the fetch won't trigger again within the tests
     */
    this.server = new Pretender(function () {
      this.get('/v1/sys/internal/ui/feature-flags', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify(FEATURE_FLAGS_RESPONSE)];
      });
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.get('/v1/sys/license/features', this.passthrough);
      this.get('/v1/sys/internal/ui/mounts', this.passthrough);
    });
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('it shows the managed namespace toolbar when feature flag exists', async function (assert) {
    await visit('/vault/auth');
    assert.ok(currentURL().startsWith('/vault/auth'), 'Redirected to auth');
    assert.ok(currentURL().includes('?namespace=admin'), 'with base namespace');
    assert.dom('[data-test-namespace-toolbar]').doesNotExist('Normal namespace toolbar does not exist');
    assert.dom('[data-test-managed-namespace-toolbar]').exists('Managed namespace toolbar exists');
    assert.dom('[data-test-managed-namespace-root]').hasText('/admin', 'Shows /admin namespace prefix');
    assert.dom('input#namespace').hasAttribute('placeholder', '/ (Default)');
    await fillIn('input#namespace', '/foo');
    const encodedNamespace = encodeURIComponent('admin/foo');
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${encodedNamespace}&with=token`,
      'Correctly prepends root to namespace'
    );
  });

  test('getManagedNamespace helper works as expected', function (assert) {
    let managedNs = getManagedNamespace(null, 'admin');
    assert.strictEqual(managedNs, 'admin', 'returns root ns when no namespace present');
    managedNs = getManagedNamespace('admin/', 'admin');
    assert.strictEqual(managedNs, 'admin', 'returns root ns when matches passed ns');
    managedNs = getManagedNamespace('adminfoo/', 'admin');
    assert.strictEqual(
      managedNs,
      'admin/adminfoo/',
      'appends passed namespace to root even if it matches without slashes'
    );
    managedNs = getManagedNamespace('admin/foo/', 'admin');
    assert.strictEqual(managedNs, 'admin/foo/', 'returns passed namespace if it starts with root and /');
  });

  test('it redirects to root prefixed ns when non-root passed', async function (assert) {
    await visit('/vault/auth?namespace=admindev');
    assert.ok(currentURL().startsWith('/vault/auth'), 'Redirected to auth');
    assert.ok(
      currentURL().includes(`?namespace=${encodeURIComponent('admin/admindev')}`),
      'with appended namespace'
    );

    assert.dom('[data-test-managed-namespace-root]').hasText('/admin', 'Shows /admin namespace prefix');
    assert.dom('input#namespace').hasValue('/admindev', 'Input has /dev value');
  });
});
