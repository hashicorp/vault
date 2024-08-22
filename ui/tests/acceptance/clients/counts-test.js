/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import clientsHandler, { STATIC_NOW } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { visit, click, currentURL, fillIn } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import timestamp from 'core/utils/timestamp';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Acceptance | clients | counts', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    clientsHandler(this.server);
    this.store = this.owner.lookup('service:store');
    return authPage.login();
  });

  test('it should prompt user to query start time for community version', async function (assert) {
    assert.expect(2);
    this.owner.lookup('service:version').type = 'community';
    await visit('/vault/clients/counts/overview');
    assert.dom(GENERAL.emptyStateTitle).hasText('No start date found');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'In order to get the most from this data, please enter a start month above. Vault will calculate new clients starting from that month.'
      );
  });

  test('it should redirect to counts overview route for transitions to parent', async function (assert) {
    await visit('/vault/clients');
    assert.strictEqual(currentURL(), '/vault/clients/counts/overview', 'Redirects to counts overview route');
  });

  test('it should persist filter query params between child routes', async function (assert) {
    await visit('/vault/clients/counts/overview');
    await click(CLIENT_COUNT.dateRange.edit);
    await fillIn(CLIENT_COUNT.dateRange.editDate('start'), '2023-03');
    await fillIn(CLIENT_COUNT.dateRange.editDate('end'), '2023-10');
    await click(GENERAL.saveButton);
    assert.strictEqual(
      currentURL(),
      '/vault/clients/counts/overview?end_time=1698710400&start_time=1677628800',
      'Start and end times added as query params'
    );

    await click(GENERAL.tab('token'));
    assert.strictEqual(
      currentURL(),
      '/vault/clients/counts/token?end_time=1698710400&start_time=1677628800',
      'Start and end times persist through child route change'
    );

    await click(GENERAL.navLink('Dashboard'));
    await click(GENERAL.navLink('Client Count'));
    assert.strictEqual(
      currentURL(),
      '/vault/clients/counts/overview',
      'Query params are reset when exiting route'
    );
  });

  test('it should render empty state if no permission to query activity data', async function (assert) {
    assert.expect(2);
    server.get('/sys/internal/counters/activity', () => {
      return overrideResponse(403);
    });
    await visit('/vault/clients/counts/overview');
    assert.dom(GENERAL.emptyStateTitle).hasText('You are not authorized');
    assert
      .dom(GENERAL.emptyStateActions)
      .hasText(
        'You must be granted permissions to view this page. Ask your administrator if you think you should have access to the /v1/sys/internal/counters/activity endpoint.'
      );
  });

  test('it should use the first month timestamp from default response rather than response start_time', async function (assert) {
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
          start_time: '2023-04-01T00:00:00Z', // reflects the first month with data
          end_time: '2023-04-30T00:00:00Z',
          by_namespace: [],
          months: [
            {
              timestamp: '2023-02-01T00:00:00Z',
              counts: null,
              namespaces: null,
              new_clients: null,
            },
            {
              timestamp: '2023-03-01T00:00:00Z',
              counts: null,
              namespaces: null,
              new_clients: null,
            },
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
                      mount_path: 'auth/authid/0',
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
                        mount_path: 'auth/authid/0',
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
    assert.dom(CLIENT_COUNT.dateRange.dateDisplay('start')).hasText('February 2023');
    assert.dom(CLIENT_COUNT.dateRange.dateDisplay('end')).hasText('April 2023');
    assert.dom(CLIENT_COUNT.counts.startDiscrepancy).exists();
  });
});
