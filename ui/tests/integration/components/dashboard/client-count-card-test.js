/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import timestamp from 'core/utils/timestamp';
import { parseAPITimestamp } from 'core/utils/date-formatters';

module('Integration | Component | dashboard/client-count-card', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.license = {
      startTime: '2018-04-03T14:15:30',
    };
  });

  test('it should display client count information', async function (assert) {
    this.server.get('sys/internal/counters/activity', () => {
      return {
        request_id: 'some-activity-id',
        data: {
          months: [
            {
              timestamp: '2023-08-01T00:00:00-07:00',
              counts: {},
              namespaces: [
                {
                  namespace_id: 'root',
                  namespace_path: '',
                  counts: {},
                  mounts: [{ mount_path: 'auth/up2/', counts: {} }],
                },
              ],
              new_clients: {
                counts: {
                  clients: 12,
                },
                namespaces: [
                  {
                    namespace_id: 'root',
                    namespace_path: '',
                    counts: {
                      clients: 12,
                    },
                    mounts: [{ mount_path: 'auth/up2/', counts: {} }],
                  },
                ],
              },
            },
          ],
          total: {
            clients: 300417,
            entity_clients: 73150,
            non_entity_clients: 227267,
          },
        },
      };
    });

    await render(hbs`<Dashboard::ClientCountCard @license={{this.license}} />`);
    assert.dom('[data-test-client-count-title]').hasText('Client count');
    assert.dom('[data-test-stat-text="total-clients"] .stat-label').hasText('Total');
    assert
      .dom('[data-test-stat-text="total-clients"] .stat-text')
      .hasText(
        `The number of clients in this billing period (Apr 2018 - ${parseAPITimestamp(
          timestamp.now().toISOString(),
          'MMM yyyy'
        )}).`
      );
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('300,417');
    assert.dom('[data-test-stat-text="new-clients"] .stat-label').hasText('New');
    assert
      .dom('[data-test-stat-text="new-clients"] .stat-text')
      .hasText('The number of clients new to Vault in the current month.');
    assert.dom('[data-test-stat-text="new-clients"] .stat-value').hasText('12');
    this.server.get('sys/internal/counters/activity', () => {
      return {
        request_id: 'some-activity-id',
        data: {
          months: [
            {
              timestamp: '2023-09-01T00:00:00-07:00',
              counts: {},
              namespaces: [
                {
                  namespace_id: 'root',
                  namespace_path: '',
                  counts: {},
                  mounts: [{ mount_path: 'auth/up2/', counts: {} }],
                },
              ],
              new_clients: {
                counts: {
                  clients: 5,
                },
                namespaces: [
                  {
                    namespace_id: 'root',
                    namespace_path: '',
                    counts: {
                      clients: 12,
                    },
                    mounts: [{ mount_path: 'auth/up2/', counts: {} }],
                  },
                ],
              },
            },
          ],
          total: {
            clients: 120,
            entity_clients: 100,
            non_entity_clients: 100,
          },
        },
      };
    });
    await click('[data-test-refresh]');
    assert.dom('[data-test-stat-text="total-clients"] .stat-label').hasText('Total');
    assert
      .dom('[data-test-stat-text="total-clients"] .stat-text')
      .hasText(
        `The number of clients in this billing period (Apr 2018 - ${parseAPITimestamp(
          timestamp.now().toISOString(),
          'MMM yyyy'
        )}).`
      );
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('120');
    assert.dom('[data-test-stat-text="new-clients"] .stat-label').hasText('New');
    assert
      .dom('[data-test-stat-text="new-clients"] .stat-text')
      .hasText('The number of clients new to Vault in the current month.');
    assert.dom('[data-test-stat-text="new-clients"] .stat-value').hasText('5');
  });
});
