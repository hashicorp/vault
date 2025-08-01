/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { visit, click, currentURL } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { login } from 'vault/tests/helpers/auth/auth-helpers';

// integration test handle general display assertions, acceptance handles nav + filtering
module.skip('Acceptance | clients | counts | client list', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    return visit('/vault');
  });

  test('it navigates to client list tab', async function (assert) {
    assert.expect(3);
    await click(GENERAL.navLink('Client Count'));
    await click(GENERAL.tab('client list'));
    assert.strictEqual(currentURL(), '/vault/clients/counts/client-list', 'it navigates to client list tab');
    assert.dom(GENERAL.tab('client list')).hasClass('active');
    await click(GENERAL.navLink('Back to main navigation'));
    assert.strictEqual(currentURL(), '/vault/dashboard', 'it navigates back to dashboard');
  });
});
