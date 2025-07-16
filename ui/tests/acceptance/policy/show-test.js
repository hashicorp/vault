/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/policy/show';
import { login } from 'vault/tests/helpers/auth/auth-helpers';

module('Acceptance | policy/acl/:name', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  test('it redirects to list if navigating to root', async function (assert) {
    await page.visit({ type: 'acl', name: 'root' });
    assert.strictEqual(
      currentURL(),
      '/vault/policies/acl',
      'navigation to root show redirects you to policy list'
    );
  });

  test('it navigates to edit when the toggle is clicked', async function (assert) {
    await page.visit({ type: 'acl', name: 'default' }).toggleEdit();
    assert.strictEqual(currentURL(), '/vault/policy/acl/default/edit', 'toggle navigates to edit page');
  });
});
