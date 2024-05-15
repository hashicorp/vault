/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | permissions', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  test('it calls resultant-acl with the users root namespace', async function (assert) {
    assert.expect(1);
    const adapter = this.owner.lookup('adapter:permissions');
    const nsService = this.owner.lookup('service:namespace');
    const auth = this.owner.lookup('service:auth');
    nsService.setNamespace('admin/foo');
    auth.setCluster('1');
    auth.set('tokens', ['vault-_root_☃1']);
    auth.setTokenData('vault-_root_☃1', { userRootNamespace: 'admin/bar', backend: { mountPath: 'token' } });

    this.server.get('/sys/internal/ui/resultant-acl', (schema, request) => {
      assert.strictEqual(
        request.requestHeaders['X-Vault-Namespace'],
        'admin/bar',
        'Namespace is users root not current path'
      );
      return {
        data: {
          exact_paths: {},
          glob_paths: {},
        },
      };
    });
    await adapter.query();
  });
  test('it calls resultant-acl with the users root namespace when root', async function (assert) {
    assert.expect(1);
    const adapter = this.owner.lookup('adapter:permissions');
    const nsService = this.owner.lookup('service:namespace');
    const auth = this.owner.lookup('service:auth');
    nsService.setNamespace('admin');
    auth.setCluster('1');
    auth.set('tokens', ['vault-_root_☃1']);
    auth.setTokenData('vault-_root_☃1', { userRootNamespace: '', backend: { mountPath: 'token' } });

    this.server.get('/sys/internal/ui/resultant-acl', (schema, request) => {
      assert.false(
        Object.keys(request.requestHeaders).includes('X-Vault-Namespace'),
        'request is called without namespace'
      );

      return {
        data: {
          exact_paths: {},
          glob_paths: {},
        },
      };
    });
    await adapter.query();
  });
});
