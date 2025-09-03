/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visit, click, currentURL } from '@ember/test-helpers';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { ClientFilters } from 'core/utils/client-count-utils';
import { CLIENT_COUNT, FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { ACTIVITY_EXPORT_STUB } from 'vault/tests/helpers/clients/client-count-helpers';

// integration test handle general display assertions, acceptance handles nav + filtering
module('Acceptance | clients | counts | client list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    // This tab is hidden on community so the version is stubbed for consistent test running on either version
    this.version = this.owner.lookup('service:version');
    this.version.type === 'enterprise';

    // The activity export endpoint returns a ReadableStream of json lines, this is not easily mocked using mirage.
    // Stubbing the adapter method return instead.
    const mockResponse = {
      status: 200,
      ok: true,
      text: () => Promise.resolve(ACTIVITY_EXPORT_STUB.trim()),
    };
    const store = this.owner.lookup('service:store');
    const adapter = store.adapterFor('clients/activity');
    this.exportDataStub = sinon.stub(adapter, 'exportData');
    this.exportDataStub.resolves(mockResponse);
    await login();
    return visit('/vault');
  });

  hooks.afterEach(async function () {
    this.exportDataStub.restore();
  });

  test('it hides client list tab on community', async function (assert) {
    this.version.type === 'community';
    assert.dom(GENERAL.tab('client list')).doesNotExist();
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

  test('filters are preset if URL includes query params', async function (assert) {
    assert.expect(4);
    const ns = 'ns2/';
    const mPath = 'auth/userpass/';
    const mType = 'userpass';
    await visit(
      `vault/clients/counts/client-list?namespace_path=${ns}&mount_path=${mPath}&mount_type=${mType}`
    );
    assert.dom(FILTERS.tag()).exists({ count: 3 }, '3 filter tags render');
    assert.dom(FILTERS.tag(ClientFilters.NAMESPACE, ns)).exists();
    assert.dom(FILTERS.tag(ClientFilters.MOUNT_PATH, mPath)).exists();
    assert.dom(FILTERS.tag(ClientFilters.MOUNT_TYPE, mType)).exists();
  });

  test('selecting filters update URL query params', async function (assert) {
    assert.expect(3);
    const ns = 'ns2/';
    const mPath = 'auth/userpass/';
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
    assert.strictEqual(
      currentURL(),
      `${url}?mount_path=${encodeURIComponent(mPath)}&mount_type=${mType}&namespace_path=${encodeURIComponent(
        ns
      )}`,
      'url query params match filters'
    );
    await click(GENERAL.button('Clear filters'));
    assert.strictEqual(currentURL(), url, '"Clear filters" resets URL query params');
  });

  test('it renders error message if export fails', async function (assert) {
    this.exportDataStub.throws(new Error('No data to export in provided time range.'));
    await visit('/vault/clients/counts/client-list');
    await click(CLIENT_COUNT.dateRange.edit);
    await click(CLIENT_COUNT.dateRange.dropdownOption(4));
    assert.dom(GENERAL.emptyStateTitle).hasText('Error');
    assert.dom(GENERAL.emptyStateActions).hasText('No data to export in provided time range.');
    // Assert the empty state message renders below the page header so user can query other dates
    assert.dom(GENERAL.tab('overview')).exists('Overview tab still renders');
    assert.dom(GENERAL.tab('client list')).exists('Client list tab still renders');
  });
});
