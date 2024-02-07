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
    nsService.setNamespace('admin/foo');
    nsService.reopen({
      userRootNamespace: 'admin/bar',
    });
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
});
