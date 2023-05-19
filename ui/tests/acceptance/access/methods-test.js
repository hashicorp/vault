/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/methods';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | /access/', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it navigates', async function (assert) {
    await page.visit();
    assert.strictEqual(currentRouteName(), 'vault.cluster.access.methods', 'navigates to the correct route');
    assert.ok(page.methodsLink.isActive, 'the first link is active');
    assert.strictEqual(page.methodsLink.text, 'Authentication methods');
  });
});
