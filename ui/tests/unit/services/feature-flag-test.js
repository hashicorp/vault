/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | feature-flag', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const service = this.owner.lookup('service:feature-flag');
    assert.ok(service);
  });

  test('it returns the namespace root when flag is present', function (assert) {
    const service = this.owner.lookup('service:feature-flag');
    assert.strictEqual(service.managedNamespaceRoot, null, 'Managed namespace root is null by default');
    service.setFeatureFlags(['VAULT_CLOUD_ADMIN_NAMESPACE']);
    assert.strictEqual(service.managedNamespaceRoot, 'admin', 'Managed namespace is admin when flag present');
    service.setFeatureFlags(['SOMETHING_ELSE']);
    assert.strictEqual(
      service.managedNamespaceRoot,
      null,
      'Flags were overwritten and root namespace is null again'
    );
  });
});
