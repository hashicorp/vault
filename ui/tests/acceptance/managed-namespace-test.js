/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { currentURL, visit, fillIn } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { getManagedNamespace } from 'vault/routes/vault/cluster';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

module('Acceptance | Enterprise | Managed namespace root', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.get('/sys/internal/ui/feature-flags', () => {
      return { feature_flags: ['VAULT_CLOUD_ADMIN_NAMESPACE'] };
    });
  });

  test('it shows the managed namespace toolbar when feature flag exists', async function (assert) {
    await visit('/vault/auth');
    assert.ok(currentURL().startsWith('/vault/auth'), 'Redirected to auth');
    assert.ok(currentURL().includes('?namespace=admin'), 'with base namespace');
    assert.dom(AUTH_FORM.managedNsRoot).hasValue('/admin', 'Shows /admin namespace prefix');
    assert.dom(AUTH_FORM.managedNsRoot).hasAttribute('readonly', '', 'HVD managed prefix input is readonly');
    assert.dom('input[name="namespace"]').hasAttribute('placeholder', '/ (default)');
    await fillIn('input[name="namespace"]', '/foo');
    const encodedNamespace = encodeURIComponent('admin/foo');
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${encodedNamespace}`,
      'Correctly prepends root to namespace when input starts with /'
    );
    await fillIn('input[name="namespace"]', 'foo');
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${encodedNamespace}`,
      'Correctly prepends root to namespace when input does not start with /'
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

    assert.dom(AUTH_FORM.managedNsRoot).hasValue('/admin', 'Shows /admin namespace prefix');
    assert.dom('input[name="namespace"]').hasValue('/admindev', 'Input has /dev value');
  });
});
