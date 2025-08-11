/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visit, click, currentURL } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { ClientFilters } from 'core/utils/client-count-utils';
import { FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';

// integration test handle general display assertions, acceptance handles nav + filtering
module('Acceptance | clients | counts | client list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.get('sys/internal/counters/activity', () => {
      return {
        request_id: 'some-activity-id',
        data: ACTIVITY_RESPONSE_STUB,
      };
    });
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

  test('filters update query params', async function (assert) {
    assert.expect(3);
    const ns = 'ns1';
    const mPath = 'auth/userpass-0';
    const mType = 'userpass';
    const url = '/vault/clients/counts/client-list';
    await visit(url);
    assert.strictEqual(currentURL(), url, 'URL does not contain query params');
    // select namespace
    await click(FILTERS.dropdownToggle(ClientFilters.NAMESPACE));
    await click(FILTERS.dropdownItem(ns));
    // select mount path
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_PATH));
    await click(FILTERS.dropdownItem(mPath));
    // select mount type
    await click(FILTERS.dropdownToggle(ClientFilters.MOUNT_TYPE));
    await click(FILTERS.dropdownItem(mType));
    await click(GENERAL.button('Apply filters'));
    assert.strictEqual(
      currentURL(),
      `${url}?mountPath=${encodeURIComponent(mPath)}&mountType=${mType}&nsLabel=${ns}`,
      'url query params match filters'
    );
    await click(GENERAL.button('Clear filters'));
    assert.strictEqual(currentURL(), url, '"Clear filters" resets URL query params');
  });
});
