/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visit, click, currentURL } from '@ember/test-helpers';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { ClientFilters } from 'core/utils/client-counts/helpers';
import { CLIENT_COUNT, FILTERS } from 'vault/tests/helpers/clients/client-count-selectors';
import { ACTIVITY_EXPORT_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import timestamp from 'core/utils/timestamp';
import clientsHandler, { STATIC_NOW } from 'vault/mirage/handlers/clients';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

// integration test handle general display assertions, acceptance handles nav + filtering
module('Acceptance | clients | counts | client list', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    //* Community version setup
    // This tab is hidden on community this setup is stubbed for consistent test running on either version
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
    // Return a consistent billing start timestamp for community versions
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    clientsHandler(this.server);
    //* End CE setup

    // The activity export endpoint returns a ReadableStream of json lines, this is not easily mocked using mirage.
    // Stubbing the api service method instead.
    const mockResponse = {
      raw: new Response(ACTIVITY_EXPORT_STUB, {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    };
    const api = this.owner.lookup('service:api');
    this.exportDataStub = sinon.stub(api.sys, 'internalClientActivityExportRaw');
    this.exportDataStub.resolves(mockResponse);

    await login();
    return visit('/vault');
  });

  hooks.afterEach(async function () {
    this.exportDataStub.restore();
  });

  test('it hides client list tab on community', async function (assert) {
    this.version.type = 'community';
    assert.dom(GENERAL.tab('client list')).doesNotExist();

    // Navigate directly to URL to test redirect
    await visit('/vault/clients/counts/client-list');
    assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'it redirects to overview');
  });

  // skip this test on CE test runs because GET sys/license/features an enterprise only endpoint
  test('enterprise: it hides client list tab on HVD managed clusters', async function (assert) {
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    assert.dom(GENERAL.tab('client list')).doesNotExist();

    // Navigate directly to URL to test redirect
    await visit('/vault/clients/counts/client-list');
    assert.strictEqual(
      currentURL(),
      '/vault/clients/counts/overview?namespace=admin',
      'it redirects to overview'
    );
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
    await click(GENERAL.dropdownToggle(ClientFilters.NAMESPACE));
    await click(FILTERS.dropdownItem(ns));
    // select mount path
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_PATH));
    await click(FILTERS.dropdownItem(mPath));
    // select mount type
    await click(GENERAL.dropdownToggle(ClientFilters.MOUNT_TYPE));
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

  test('it renders error message if export has no data', async function (assert) {
    const emptyResponse = {
      raw: new Response({}, { status: 204, headers: { 'Content-Type': 'application/json' } }),
    };
    this.exportDataStub.resolves(emptyResponse);
    await visit('/vault/clients/counts/client-list');
    await click(CLIENT_COUNT.dateRange.edit);
    await click(CLIENT_COUNT.dateRange.dropdownOption(4));
    assert.dom(GENERAL.emptyStateTitle).hasText('No data found');
    assert.dom(GENERAL.emptyStateMessage).hasText('No data to export in provided time range.');
    // Assert the empty state message renders below the page header so user can query other dates
    assert.dom(GENERAL.tab('overview')).exists('Overview tab still renders');
    assert.dom(GENERAL.tab('client list')).exists('Client list tab still renders');
  });

  test('it renders error message for permission denied', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.exportDataStub.rejects(getErrorResponse(error, 403));
    await visit('/vault/clients/counts/client-list');
    assert.dom(GENERAL.emptyStateTitle).hasText('You are not authorized');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Viewing export data requires sudo permissions to /sys/internal/counters/activity/export.');
    assert.dom(GENERAL.emptyStateActions).hasText('Client Export Documentation');
  });

  // since permissions errors are specially handled, test that a non 403 is handled correctly
  test('it renders error message for a server error', async function (assert) {
    const error = { errors: ['uh oh'] };
    this.exportDataStub.rejects(getErrorResponse(error, 500));
    await visit('/vault/clients/counts/client-list');
    assert.dom(GENERAL.pageError.errorTitle(500)).hasText('Error');
    assert.dom(GENERAL.pageError.errorDetails).hasText('uh oh');
  });
});
