/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | flags', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const service = this.owner.lookup('service:flags');
    assert.ok(service);
  });

  test('it returns the namespace root when flag is present', function (assert) {
    const service = this.owner.lookup('service:flags');
    assert.strictEqual(service.managedNamespaceRoot, null, 'Managed namespace root is null by default');
    service.setFlags(['VAULT_CLOUD_ADMIN_NAMESPACE']);
    assert.strictEqual(service.managedNamespaceRoot, 'admin', 'Managed namespace is admin when flag present');
    service.setFlags(['SOMETHING_ELSE']);
    assert.strictEqual(
      service.managedNamespaceRoot,
      null,
      'Flags were overwritten and root namespace is null again'
    );
  });
});
