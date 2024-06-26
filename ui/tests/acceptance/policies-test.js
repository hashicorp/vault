/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentURL, currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import { GENERAL } from '../helpers/general-selectors';

module('Acceptance | policies', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('it redirects from plural policies route with unknown policy type', async function (assert) {
    await visit('/vault/policies/foo');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');
  });

  test('it redirects from singular policy route with unknown policy type', async function (assert) {
    await visit('/vault/policy/foo/default');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');

    await visit('/vault/policy/foo/default/edit');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');
  });

  test('it redirects to policies/acl on index navigation', async function (assert) {
    await visit('/vault/policy/acl');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');
  });

  test('it navigates to and from policy show page from sidebar', async function (assert) {
    await visit('/vault/dashboard');
    await click(GENERAL.navLink('Policies'));
    assert.strictEqual(currentURL(), '/vault/policies/acl');
    await click('[data-test-policy-link="default"]');
    assert.strictEqual(currentURL(), '/vault/policy/acl/default');
    await click(GENERAL.navLink('Back to main navigation'));
    assert.strictEqual(currentURL(), '/vault/dashboard');
  });
});
