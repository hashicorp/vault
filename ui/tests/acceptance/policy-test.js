/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentURL, currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | policies', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it redirects to acls with unknown policy type', async function (assert) {
    await visit('/vault/policies/foo');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');
  });
});
