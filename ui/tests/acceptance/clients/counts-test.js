/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import clientsHandler, { STATIC_NOW } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { visit, currentURL, click } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import timestamp from 'core/utils/timestamp';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { format } from 'date-fns';

module('Acceptance | clients | counts', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.timestampStub = sinon.stub(timestamp, 'now');
    this.timestampStub.returns(STATIC_NOW);
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    return login();
  });

  test('it should prompt user to query start time for community version', async function (assert) {
    assert.expect(2);
    this.owner.lookup('service:version').type = 'community';
    await visit('/vault/clients/counts/overview');
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('Input the start and end dates to view client attribution by path.');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Only historical data may be queried. No data is available for the current month.');
  });

  test('it does not make a request to the export api on community versions', async function (assert) {
    assert.expect(1);
    this.owner.lookup('service:version').type = 'community';
    server.get('/sys/internal/counters/activity/export', () => {
      // passing "false" because a request should NOT be made, so if this assertion is hit we want it to fail
      assert.true(false, 'it does not make request to export API on community versions ');
    });
    await visit('/vault/clients/counts/overview');
    assert.dom(GENERAL.tab('client list')).doesNotExist();
  });

  test('it should redirect to counts overview route for transitions to parent', async function (assert) {
    await visit('/vault/clients');
    assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'Redirects to counts overview route');
  });

  test('it should render empty state if no permission to query activity data', async function (assert) {
    assert.expect(2);
    server.get('/sys/internal/counters/activity', () => {
      return overrideResponse(403);
    });
    await visit('/vault/clients/counts/overview');
    assert.dom(GENERAL.emptyStateTitle).hasText('ERROR 403 You are not authorized');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'You must be granted permissions to view this page. Ask your administrator if you think you should have access to the /v1/sys/internal/counters/activity endpoint.'
      );
  });

  test('it should use the response start_time as the timestamp', async function (assert) {
    const getCounts = () => {
      return {
        acme_clients: 0,
        clients: 0,
        entity_clients: 0,
        non_entity_clients: 0,
        secret_syncs: 0,
        distinct_entities: 0,
        non_entity_tokens: 1,
      };
    };
    // set to enterprise because when community the initial activity call is skipped
    this.owner.lookup('service:version').type = 'enterprise';
    this.server.get('/sys/internal/counters/activity', function () {
      return {
        request_id: 'some-activity-id',
        data: {
          start_time: '2023-04-01T00:00:00Z', // API returns complete billing cycles, so we use this date as the source of truth
          end_time: '2023-04-30T00:00:00Z',
          by_namespace: [],
          months: [
            {
              timestamp: '2023-04-01T00:00:00Z',
              counts: getCounts(),
              namespaces: [
                {
                  namespace_id: 'root',
                  namespace_path: '',
                  counts: getCounts(),
                  mounts: [
                    {
                      mount_path: 'auth/userpass-0',
                      counts: getCounts(),
                    },
                  ],
                },
              ],
              new_clients: {
                counts: getCounts(),
                namespaces: [
                  {
                    namespace_id: 'root',
                    namespace_path: '',
                    counts: getCounts(),
                    mounts: [
                      {
                        mount_path: 'auth/userpass-0',
                        counts: getCounts(),
                      },
                    ],
                  },
                ],
              },
            },
          ],
          total: getCounts(),
        },
      };
    });
    await visit('/vault/clients/counts/overview');
    assert.dom(CLIENT_COUNT.dateRange.dateDisplay('start')).hasText('April 2023');
    assert.dom(CLIENT_COUNT.dateRange.dateDisplay('end')).hasText('April 2023');
  });

  module('manual refresh', function (hooks) {
    hooks.beforeEach(async function () {
      const router = this.owner.lookup('service:router');
      this.refreshSpy = sinon.spy(router, 'refresh');
      return login();
    });

    // Date querying is different in CE vs Enterprise, but the refresh behaves the same.
    // For simplicity, just test this action on enterprise versions.
    test('enterprise: it refreshes the overview route and preserves query params', async function (assert) {
      await visit('/vault/clients/counts/overview');
      assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'current view is the overview page');
      // Change date to add query params
      await click(CLIENT_COUNT.dateRange.edit);
      await click(CLIENT_COUNT.dateRange.dropdownOption(1));
      assert
        .dom(GENERAL.hdsPageHeaderSubtitle)
        .hasTextContaining(`Dashboard last updated: ${format(STATIC_NOW, 'MMM d yyyy')}`);
      // Save URL with query params before clicking refresh
      const url = currentURL();
      // re-stub with a completely different year/month/day before clicking refresh
      // to mock the timestamp updating when page reloads
      const fakeUpdatedNow = new Date('2025-07-02T23:25:13Z');
      this.timestampStub.returns(fakeUpdatedNow);
      await click(GENERAL.button('Refresh page'));
      assert.true(this.refreshSpy.calledOnce, 'router.refresh() is called once');
      assert.strictEqual(currentURL(), url, 'url is the same after clicking refresh');
      assert
        .dom(GENERAL.hdsPageHeaderSubtitle)
        .hasTextContaining(`Dashboard last updated: ${format(fakeUpdatedNow, 'MMM d yyyy')}`);
    });

    test('enterprise: it refreshes the client-list route and preserves query params', async function (assert) {
      await visit('/vault/clients/counts/client-list');
      assert.strictEqual(
        currentURL(),
        '/vault/clients/counts/client-list',
        'current view is the client-list page'
      );
      // Change date to add query params
      await click(CLIENT_COUNT.dateRange.edit);
      await click(CLIENT_COUNT.dateRange.dropdownOption(1));
      assert
        .dom(GENERAL.hdsPageHeaderSubtitle)
        .hasTextContaining(`Dashboard last updated: ${format(STATIC_NOW, 'MMM d yyyy')}`);
      // Save URL with query params before clicking refresh
      const url = currentURL();
      // re-stub with a completely different year/month/day before clicking refresh
      // to mock the timestamp updating when page reloads
      const fakeUpdatedNow = new Date('2025-07-02T23:25:13Z');
      this.timestampStub.returns(fakeUpdatedNow);
      await click(GENERAL.button('Refresh page'));
      assert.true(this.refreshSpy.calledOnce, 'router.refresh() is called once');
      assert.strictEqual(currentURL(), url, 'url is the same after clicking refresh');
      assert
        .dom(GENERAL.hdsPageHeaderSubtitle)
        .hasTextContaining(`Dashboard last updated: ${format(fakeUpdatedNow, 'MMM d yyyy')}`);
    });
  });
});
