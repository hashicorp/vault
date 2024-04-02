/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import {
  filterVersionHistory,
  formatByMonths,
  formatByNamespace,
  homogenizeClientNaming,
  namespaceArrayToObject,
  sortMonthsByTimestamp,
} from 'core/utils/client-count-utils';
import { LICENSE_START } from 'vault/mirage/handlers/clients';
import { ACTIVITY_RESPONSE as RESPONSE, VERSION_HISTORY } from 'vault/tests/helpers/clients';

module('Integration | Util | client count utils', function (hooks) {
  setupTest(hooks);

  test('filterVersionHistory: returns version data for relevant upgrades that occurred during date range', async function (assert) {
    assert.expect(2);
    // LICENSE_START is '2023-07-02T00:00:00Z'
    const original = [...VERSION_HISTORY];
    const expected = [
      {
        previousVersion: null,
        timestampInstalled: '2023-07-02T00:00:00.000Z',
        version: '1.9.0',
      },
      {
        previousVersion: '1.9.1',
        timestampInstalled: '2023-09-02T00:00:00.000Z',
        version: '1.10.1',
      },
    ];

    const startTime = LICENSE_START.toISOString(); // same as license start to catch same day edge cases
    const endTime = '2024-03-04T16:14:21.000Z';
    assert.propEqual(
      filterVersionHistory(VERSION_HISTORY, startTime, endTime),
      expected,
      'it only returns upgrades between given start and end times'
    );
    assert.propEqual(VERSION_HISTORY, original, 'it does not modify original array');
  });

  test('formatByMonths: formats the months array', async function (assert) {
    assert.expect(5);
    const original = [...RESPONSE.months];
    const expected = [
      {
        month: '8/23',
        timestamp: '2023-08-01T00:00:00-07:00',
        namespaces: [],
        new_clients: {
          month: '8/23',
          timestamp: '2023-08-01T00:00:00-07:00',
          namespaces: [],
        },
        namespaces_by_key: {},
      },
      {
        month: '9/23',
        timestamp: '2023-09-01T00:00:00-07:00',
        clients: 8592,
        entity_clients: 1329,
        non_entity_clients: 1738,
        secret_syncs: 5525,
        namespaces: [
          {
            label: 'root',
            clients: 5632,
            entity_clients: 1279,
            non_entity_clients: 1598,
            secret_syncs: 2755,
            mounts: [
              {
                label: 'auth/authid0',
                clients: 2877,
                entity_clients: 1279,
                non_entity_clients: 1598,
                secret_syncs: 0,
              },
              {
                label: 'kvv2-engine-0',
                clients: 2755,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 2755,
              },
            ],
          },
          {
            label: 'ns/1',
            clients: 2960,
            entity_clients: 50,
            non_entity_clients: 140,
            secret_syncs: 2770,
            mounts: [
              {
                label: 'kvv2-engine-0',
                clients: 2770,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 2770,
              },
              {
                label: 'auth/authid0',
                clients: 190,
                entity_clients: 50,
                non_entity_clients: 140,
                secret_syncs: 0,
              },
            ],
          },
        ],
        namespaces_by_key: {
          root: {
            month: '9/23',
            timestamp: '2023-09-01T00:00:00-07:00',
            clients: 5632,
            entity_clients: 1279,
            non_entity_clients: 1598,
            secret_syncs: 2755,
            new_clients: {
              month: '9/23',
              label: 'root',
              clients: 94,
              entity_clients: 9,
              non_entity_clients: 19,
              secret_syncs: 66,
              mounts: [
                {
                  label: 'kvv2-engine-0',
                  clients: 66,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 66,
                },
                {
                  label: 'auth/authid0',
                  clients: 28,
                  entity_clients: 9,
                  non_entity_clients: 19,
                  secret_syncs: 0,
                },
              ],
            },
            mounts_by_key: {
              'auth/authid0': {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00-07:00',
                label: 'auth/authid0',
                clients: 2877,
                entity_clients: 1279,
                non_entity_clients: 1598,
                secret_syncs: 0,
                new_clients: {
                  month: '9/23',
                  label: 'auth/authid0',
                  clients: 28,
                  entity_clients: 9,
                  non_entity_clients: 19,
                  secret_syncs: 0,
                },
              },
              'kvv2-engine-0': {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00-07:00',
                label: 'kvv2-engine-0',
                clients: 2755,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 2755,
                new_clients: {
                  month: '9/23',
                  label: 'kvv2-engine-0',
                  clients: 66,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 66,
                },
              },
            },
          },
          'ns/1': {
            month: '9/23',
            timestamp: '2023-09-01T00:00:00-07:00',
            clients: 2960,
            entity_clients: 50,
            non_entity_clients: 140,
            secret_syncs: 2770,
            new_clients: {
              month: '9/23',
              label: 'ns/1',
              clients: 192,
              entity_clients: 30,
              non_entity_clients: 62,
              secret_syncs: 100,
              mounts: [
                {
                  label: 'kvv2-engine-0',
                  clients: 100,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 100,
                },
                {
                  label: 'auth/authid0',
                  clients: 92,
                  entity_clients: 30,
                  non_entity_clients: 62,
                  secret_syncs: 0,
                },
              ],
            },
            mounts_by_key: {
              'kvv2-engine-0': {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00-07:00',
                label: 'kvv2-engine-0',
                clients: 2770,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 2770,
                new_clients: {
                  month: '9/23',
                  label: 'kvv2-engine-0',
                  clients: 100,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 100,
                },
              },
              'auth/authid0': {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00-07:00',
                label: 'auth/authid0',
                clients: 190,
                entity_clients: 50,
                non_entity_clients: 140,
                secret_syncs: 0,
                new_clients: {
                  month: '9/23',
                  label: 'auth/authid0',
                  clients: 92,
                  entity_clients: 30,
                  non_entity_clients: 62,
                  secret_syncs: 0,
                },
              },
            },
          },
        },
        new_clients: {
          month: '9/23',
          timestamp: '2023-09-01T00:00:00-07:00',
          clients: 286,
          entity_clients: 39,
          non_entity_clients: 81,
          secret_syncs: 166,
          namespaces: [
            {
              label: 'ns/1',
              clients: 192,
              entity_clients: 30,
              non_entity_clients: 62,
              secret_syncs: 100,
              mounts: [
                {
                  label: 'kvv2-engine-0',
                  clients: 100,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 100,
                },
                {
                  label: 'auth/authid0',
                  clients: 92,
                  entity_clients: 30,
                  non_entity_clients: 62,
                  secret_syncs: 0,
                },
              ],
            },
            {
              label: 'root',
              clients: 94,
              entity_clients: 9,
              non_entity_clients: 19,
              secret_syncs: 66,
              mounts: [
                {
                  label: 'kvv2-engine-0',
                  clients: 66,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 66,
                },
                {
                  label: 'auth/authid0',
                  clients: 28,
                  entity_clients: 9,
                  non_entity_clients: 19,
                  secret_syncs: 0,
                },
              ],
            },
          ],
        },
      },
    ];

    const [formattedNoData, formattedWithActivity] = formatByMonths(RESPONSE.months);

    // instead of asserting the whole expected response, broken up so tests are easier to debug
    // but kept whole above to copy/paste updated response expectations in the future
    const [expectedNoData, expectedWithActivity] = expected;
    const { namespaces, namespaces_by_key, new_clients } = expectedWithActivity;

    assert.propEqual(formattedNoData, expectedNoData, 'it formats months without data');
    assert.propEqual(
      formattedWithActivity.namespaces,
      namespaces,
      'it formats namespaces array for months with data'
    );
    assert.propEqual(
      formattedWithActivity.namespaces_by_key,
      namespaces_by_key,
      'it formats namespaces_by_key for months with data'
    );
    assert.propEqual(
      formattedWithActivity.new_clients,
      new_clients,
      'it formats new_clients block for months with data'
    );
    assert.propEqual(RESPONSE.months, original, 'it does not modify original months array');
  });

  test('formatByNamespace: formats namespace arrays with and without mounts', async function (assert) {
    assert.expect(2);
    const original = [...RESPONSE.by_namespace];
    const expected = [
      {
        clients: 5354,
        entity_clients: 1033,
        label: 'root',
        mounts: [
          {
            clients: 2957,
            entity_clients: 1033,
            label: 'auth/authid0',
            non_entity_clients: 1924,
            secret_syncs: 0,
          },
          {
            clients: 2397,
            entity_clients: 0,
            label: 'kvv2-engine-0',
            non_entity_clients: 0,
            secret_syncs: 2397,
          },
        ],
        non_entity_clients: 1924,
        secret_syncs: 2397,
      },
      {
        clients: 2251,
        entity_clients: 783,
        label: 'ns/1',
        mounts: [
          {
            clients: 1976,
            entity_clients: 783,
            label: 'auth/authid0',
            non_entity_clients: 1193,
            secret_syncs: 0,
          },
          {
            clients: 275,
            entity_clients: 0,
            label: 'kvv2-engine-0',
            non_entity_clients: 0,
            secret_syncs: 275,
          },
        ],
        non_entity_clients: 1193,
        secret_syncs: 275,
      },
    ];
    assert.propEqual(formatByNamespace(RESPONSE.by_namespace), expected);
    assert.propEqual(RESPONSE.by_namespace, original, 'it does not modify original by_namespace array');
  });

  test('homogenizeClientNaming: homogenizes key names when both old and new keys exist, or just old key names', async function (assert) {
    assert.expect(2);
    const original = { ...RESPONSE.total };
    const expected = {
      entity_clients: 1816,
      non_entity_clients: 3117,
      secret_syncs: 2672,
      clients: 7605,
    };
    assert.propEqual(homogenizeClientNaming(RESPONSE.total), expected);
    assert.propEqual(RESPONSE.total, original, 'it does not modify original object');
  });

  test('sortMonthsByTimestamp: sorts timestamps chronologically, oldest to most recent', async function (assert) {
    assert.expect(2);
    // API returns them in order so this test is extra extra
    const unOrdered = [RESPONSE.months[1], RESPONSE.months[0]]; // mixup order
    const original = [...RESPONSE.months];
    const expected = RESPONSE.months;
    assert.propEqual(sortMonthsByTimestamp(unOrdered), expected);
    assert.propEqual(RESPONSE.months, original, 'it does not modify original array');
  });

  test('namespaceArrayToObject: transforms data without modifying original', async function (assert) {
    assert.expect(2);
    const { namespaces, new_clients } = RESPONSE.months[1];
    const monthNamespaces = formatByNamespace(namespaces);
    const newClients = formatByNamespace(new_clients.namespaces);
    const byNamespaceKeyObject = namespaceArrayToObject(
      monthNamespaces,
      newClients,
      '9/23',
      '2023-9-01T00:00:00Z'
    );
    const expected = {
      root: {
        month: '9/23',
        timestamp: '2023-9-01T00:00:00Z',
        clients: 5632,
        entity_clients: 1279,
        non_entity_clients: 1598,
        secret_syncs: 2755,
        new_clients: {
          month: '9/23',
          label: 'root',
          clients: 94,
          entity_clients: 9,
          non_entity_clients: 19,
          secret_syncs: 66,
          mounts: [
            {
              label: 'kvv2-engine-0',
              clients: 66,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 66,
            },
            {
              label: 'auth/authid0',
              clients: 28,
              entity_clients: 9,
              non_entity_clients: 19,
              secret_syncs: 0,
            },
          ],
        },
        mounts_by_key: {
          'auth/authid0': {
            month: '9/23',
            timestamp: '2023-9-01T00:00:00Z',
            label: 'auth/authid0',
            clients: 2877,
            entity_clients: 1279,
            non_entity_clients: 1598,
            secret_syncs: 0,
            new_clients: {
              month: '9/23',
              label: 'auth/authid0',
              clients: 28,
              entity_clients: 9,
              non_entity_clients: 19,
              secret_syncs: 0,
            },
          },
          'kvv2-engine-0': {
            month: '9/23',
            timestamp: '2023-9-01T00:00:00Z',
            label: 'kvv2-engine-0',
            clients: 2755,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 2755,
            new_clients: {
              month: '9/23',
              label: 'kvv2-engine-0',
              clients: 66,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 66,
            },
          },
        },
      },
      'ns/1': {
        month: '9/23',
        timestamp: '2023-9-01T00:00:00Z',
        clients: 2960,
        entity_clients: 50,
        non_entity_clients: 140,
        secret_syncs: 2770,
        new_clients: {
          month: '9/23',
          label: 'ns/1',
          clients: 192,
          entity_clients: 30,
          non_entity_clients: 62,
          secret_syncs: 100,
          mounts: [
            {
              label: 'kvv2-engine-0',
              clients: 100,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 100,
            },
            {
              label: 'auth/authid0',
              clients: 92,
              entity_clients: 30,
              non_entity_clients: 62,
              secret_syncs: 0,
            },
          ],
        },
        mounts_by_key: {
          'kvv2-engine-0': {
            month: '9/23',
            timestamp: '2023-9-01T00:00:00Z',
            label: 'kvv2-engine-0',
            clients: 2770,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 2770,
            new_clients: {
              month: '9/23',
              label: 'kvv2-engine-0',
              clients: 100,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 100,
            },
          },
          'auth/authid0': {
            month: '9/23',
            timestamp: '2023-9-01T00:00:00Z',
            label: 'auth/authid0',
            clients: 190,
            entity_clients: 50,
            non_entity_clients: 140,
            secret_syncs: 0,
            new_clients: {
              month: '9/23',
              label: 'auth/authid0',
              clients: 92,
              entity_clients: 30,
              non_entity_clients: 62,
              secret_syncs: 0,
            },
          },
        },
      },
    };
    assert.propEqual(
      byNamespaceKeyObject,
      expected,
      'it returns object with namespaces by key and includes mounts_by_key'
    );
    assert.propEqual(
      namespaceArrayToObject(null, null, '10/21', 'timestamp-here'),
      {},
      'returns an empty object when monthByNamespace = null'
    );
  });
});
