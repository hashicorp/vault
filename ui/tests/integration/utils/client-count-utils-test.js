/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import {
  filterVersionHistory,
  formatByMonths,
  formatByNamespace,
  destructureClientCounts,
  namespaceArrayToObject,
  sortMonthsByTimestamp,
  setStartTimeQuery,
} from 'core/utils/client-count-utils';
import clientsHandler from 'vault/mirage/handlers/clients';
import {
  ACTIVITY_RESPONSE_STUB as RESPONSE,
  MIXED_ACTIVITY_RESPONSE_STUB as MIXED_RESPONSE,
  SERIALIZED_ACTIVITY_RESPONSE,
} from 'vault/tests/helpers/clients/client-count-helpers';

/*
formatByNamespace, formatByMonths, destructureClientCounts are utils
used to normalize the sys/counters/activity response in the clients/activity
serializer. these functions are tested individually here, instead of all at once
in a serializer test for easier debugging
*/

// TODO refactor tests into a module for each util method, then make each assertion its separate tests

module('Integration | Util | client count utils', function (hooks) {
  setupTest(hooks);

  module('filterVersionHistory', function (hooks) {
    setupMirage(hooks);

    hooks.beforeEach(async function () {
      clientsHandler(this.server);
      const store = this.owner.lookup('service:store');
      // format returned by model hook in routes/vault/cluster/clients.ts
      this.versionHistory = await store.findAll('clients/version-history').then((resp) => {
        return resp.map(({ version, previousVersion, timestampInstalled }) => {
          return {
            // order of keys needs to match expected order
            previousVersion,
            timestampInstalled,
            version,
          };
        });
      });
    });

    test('it returns version data for upgrade to notable versions: 1.9, 1.10, 1.17', async function (assert) {
      assert.expect(3);
      const original = [...this.versionHistory];
      const expected = [
        {
          previousVersion: '1.9.0',
          timestampInstalled: '2023-08-02T00:00:00.000Z',
          version: '1.9.1',
        },
        {
          previousVersion: '1.9.1',
          timestampInstalled: '2023-09-02T00:00:00.000Z',
          version: '1.10.1',
        },
        {
          previousVersion: '1.16.0',
          timestampInstalled: '2023-12-02T00:00:00.000Z',
          version: '1.17.0',
        },
      ];
      // set start/end times longer than version history to test all relevant upgrades return
      const startTime = '2023-06-02T00:00:00Z'; // first upgrade installed '2023-07-02T00:00:00Z'
      const endTime = '2024-03-04T16:14:21.000Z'; // latest upgrade installed '2023-12-02T01:00:00.000Z'
      const filteredHistory = filterVersionHistory(this.versionHistory, startTime, endTime);
      assert.deepEqual(
        JSON.stringify(filteredHistory),
        JSON.stringify(expected),
        'it returns all notable upgrades'
      );
      assert.notPropContains(
        filteredHistory,
        {
          version: '1.9.0',
          previousVersion: null,
          timestampInstalled: '2023-07-02T00:00:00Z',
        },
        'does not include version history if previous_version is null'
      );
      assert.propEqual(this.versionHistory, original, 'it does not modify original array');
    });

    test('it only returns version data for initial upgrades between given date range', async function (assert) {
      assert.expect(2);
      const expected = [
        {
          previousVersion: '1.9.0',
          timestampInstalled: '2023-08-02T00:00:00.000Z',
          version: '1.9.1',
        },
        {
          previousVersion: '1.9.1',
          timestampInstalled: '2023-09-02T00:00:00.000Z',
          version: '1.10.1',
        },
      ];
      const startTime = '2023-08-02T00:00:00.000Z'; // same date as 1.9.1 install date to catch same day edge cases
      const endTime = '2023-11-02T00:00:00.000Z';
      const filteredHistory = filterVersionHistory(this.versionHistory, startTime, endTime);
      assert.deepEqual(
        JSON.stringify(filteredHistory),
        JSON.stringify(expected),
        'it only returns upgrades during date range'
      );
      assert.notPropContains(
        filteredHistory,
        {
          version: '1.10.3',
          previousVersion: '1.10.1',
          timestampInstalled: '2023-09-23T00:00:00.000Z',
        },
        'it does not return subsequent patch versions of the same notable upgrade version'
      );
    });
  });

  test('formatByMonths: it formats the months array', async function (assert) {
    assert.expect(9);
    const original = [...RESPONSE.months];

    const [formattedNoData, formattedWithActivity, formattedNoNew] = formatByMonths(RESPONSE.months);

    // instead of asserting the whole expected response, broken up so tests are easier to debug
    // but kept whole above to copy/paste updated response expectations in the future
    const [expectedNoData, expectedWithActivity, expectedNoNew] = SERIALIZED_ACTIVITY_RESPONSE.by_month;

    assert.propEqual(formattedNoData, expectedNoData, 'it formats months without data');
    ['namespaces', 'new_clients', 'namespaces_by_key'].forEach((key) => {
      assert.propEqual(
        formattedWithActivity[key],
        expectedWithActivity[key],
        `it formats ${key} array for months with data`
      );
      assert.propEqual(
        formattedNoNew[key],
        expectedNoNew[key],
        `it formats the ${key} array for months with no new clients`
      );
    });

    assert.propEqual(RESPONSE.months, original, 'it does not modify original months array');
    assert.propEqual(formatByMonths([]), [], 'it returns an empty array if the months key is empty');
  });

  test('formatByNamespace: it formats namespace array with mounts', async function (assert) {
    const original = [...RESPONSE.by_namespace];
    const expectedNs1 = SERIALIZED_ACTIVITY_RESPONSE.by_namespace.find((ns) => ns.label === 'ns1');
    const formattedNs1 = formatByNamespace(RESPONSE.by_namespace).find((ns) => ns.label === 'ns1');
    assert.expect(2 + expectedNs1.mounts.length * 2);

    assert.propEqual(formattedNs1, expectedNs1, 'it formats ns1/ namespace');
    assert.propEqual(RESPONSE.by_namespace, original, 'it does not modify original by_namespace array');

    formattedNs1.mounts.forEach((mount) => {
      const expectedMount = expectedNs1.mounts.find((m) => m.label === mount.label);
      assert.propEqual(Object.keys(mount), Object.keys(expectedMount), `${mount} as expected keys`);
      assert.propEqual(Object.values(mount), Object.values(expectedMount), `${mount} as expected values`);
    });
  });

  test('destructureClientCounts: it returns relevant key names when both old and new keys exist', async function (assert) {
    assert.expect(2);
    const original = { ...RESPONSE.total };
    const expected = {
      acme_clients: 9702,
      clients: 35287,
      entity_clients: 8258,
      non_entity_clients: 8227,
      secret_syncs: 9100,
    };
    assert.propEqual(destructureClientCounts(RESPONSE.total), expected);
    assert.propEqual(RESPONSE.total, original, 'it does not modify original object');
  });

  test('sortMonthsByTimestamp: sorts timestamps chronologically, oldest to most recent', async function (assert) {
    assert.expect(2);
    // API returns them in order so this test is extra extra
    const unOrdered = [RESPONSE.months[1], RESPONSE.months[0], RESPONSE.months[3], RESPONSE.months[2]]; // mixup order
    const original = [...RESPONSE.months];
    const expected = RESPONSE.months;
    assert.propEqual(sortMonthsByTimestamp(unOrdered), expected);
    assert.propEqual(RESPONSE.months, original, 'it does not modify original array');
  });

  test('namespaceArrayToObject: it returns namespaces_by_key and mounts_by_key', async function (assert) {
    // namespaceArrayToObject only called when there are counts, so skip month 0 which has no counts
    for (let i = 1; i < RESPONSE.months.length; i++) {
      const original = { ...RESPONSE.months[i] };
      const expectedObject = SERIALIZED_ACTIVITY_RESPONSE.by_month[i].namespaces_by_key;
      const formattedTotal = formatByNamespace(RESPONSE.months[i].namespaces);
      const testObject = namespaceArrayToObject(
        formattedTotal,
        formatByNamespace(RESPONSE.months[i].new_clients.namespaces),
        `${i + 6}/23`,
        original.timestamp
      );
      const { root } = testObject;
      const { root: expectedRoot } = expectedObject;

      assert.propEqual(
        root?.new_clients,
        expectedRoot?.new_clients,
        `it formats namespaces new_clients for ${original.timestamp}`
      );
      assert.propEqual(root.mounts_by_key, expectedRoot.mounts_by_key, 'it formats namespaces mounts_by_key');
      assert.propContains(root, expectedRoot, 'namespace has correct keys');

      assert.propEqual(RESPONSE.months[i], original, 'it does not modify original month data');
    }
  });

  // TESTS FOR COMBINED ACTIVITY DATA - no mount attribution < 1.10
  test('it formats the namespaces array with no mount attribution (activity log data < 1.10)', async function (assert) {
    assert.expect(2);
    const noMounts = [
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 10,
          entity_clients: 10,
          non_entity_tokens: 20,
          non_entity_clients: 20,
          secret_syncs: 0,
          acme_clients: 0,
          clients: 30,
        },
        mounts: [
          {
            counts: {
              distinct_entities: 10,
              entity_clients: 10,
              non_entity_tokens: 20,
              non_entity_clients: 20,
              secret_syncs: 0,
              acme_clients: 0,
              clients: 30,
            },
            mount_path: 'no mount accessor (pre-1.10 upgrade?)',
          },
        ],
      },
    ];
    const expected = [
      {
        acme_clients: 0,
        clients: 30,
        entity_clients: 10,
        label: 'root',
        mounts: [
          {
            acme_clients: 0,
            clients: 30,
            entity_clients: 10,
            label: 'no mount accessor (pre-1.10 upgrade?)',
            non_entity_clients: 20,
            secret_syncs: 0,
          },
        ],
        non_entity_clients: 20,
        secret_syncs: 0,
      },
    ];
    assert.propEqual(formatByNamespace(noMounts), expected, 'it formats namespace without mounts');
    assert.propEqual(formatByNamespace([]), [], 'it returns an empty array if the by_namespace key is empty');
  });

  test('it formats the months array with mixed activity data', async function (assert) {
    assert.expect(3);

    const [, formattedWithActivity] = formatByMonths(MIXED_RESPONSE.months);
    // mirage isn't set up to generate mixed data, so hardcoding the expected responses here
    assert.propEqual(
      formattedWithActivity.namespaces,
      [
        {
          acme_clients: 0,
          clients: 3,
          entity_clients: 3,
          label: 'root',
          mounts: [
            {
              acme_clients: 0,
              clients: 2,
              entity_clients: 2,
              label: 'no mount accessor (pre-1.10 upgrade?)',
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            {
              acme_clients: 0,
              clients: 1,
              entity_clients: 1,
              label: 'auth/u/',
              non_entity_clients: 0,
              secret_syncs: 0,
            },
          ],
          non_entity_clients: 0,
          secret_syncs: 0,
        },
      ],
      'it formats combined data for monthly namespaces spanning upgrade to 1.10'
    );
    assert.propEqual(
      formattedWithActivity.new_clients,
      {
        acme_clients: 0,
        clients: 3,
        entity_clients: 3,
        month: '4/24',
        namespaces: [
          {
            acme_clients: 0,
            clients: 3,
            entity_clients: 3,
            label: 'root',
            mounts: [
              {
                acme_clients: 0,
                clients: 2,
                entity_clients: 2,
                label: 'no mount accessor (pre-1.10 upgrade?)',
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                acme_clients: 0,
                clients: 1,
                entity_clients: 1,
                label: 'auth/u/',
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            ],
            non_entity_clients: 0,
            secret_syncs: 0,
          },
        ],
        non_entity_clients: 0,
        secret_syncs: 0,
        timestamp: '2024-04-01T00:00:00Z',
      },
      'it formats combined data for monthly new_clients spanning upgrade to 1.10'
    );
    assert.propEqual(
      formattedWithActivity.namespaces_by_key,
      {
        root: {
          acme_clients: 0,
          clients: 3,
          entity_clients: 3,
          month: '4/24',
          mounts_by_key: {
            'auth/u/': {
              acme_clients: 0,
              clients: 1,
              entity_clients: 1,
              label: 'auth/u/',
              month: '4/24',
              new_clients: {
                acme_clients: 0,
                clients: 1,
                entity_clients: 1,
                label: 'auth/u/',
                month: '4/24',
                timestamp: '2024-04-01T00:00:00Z',
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              non_entity_clients: 0,
              secret_syncs: 0,
              timestamp: '2024-04-01T00:00:00Z',
            },
            'no mount accessor (pre-1.10 upgrade?)': {
              acme_clients: 0,
              clients: 2,
              entity_clients: 2,
              label: 'no mount accessor (pre-1.10 upgrade?)',
              month: '4/24',
              new_clients: {
                acme_clients: 0,
                clients: 2,
                entity_clients: 2,
                label: 'no mount accessor (pre-1.10 upgrade?)',
                month: '4/24',
                timestamp: '2024-04-01T00:00:00Z',
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              non_entity_clients: 0,
              secret_syncs: 0,
              timestamp: '2024-04-01T00:00:00Z',
            },
          },
          new_clients: {
            acme_clients: 0,
            clients: 3,
            entity_clients: 3,
            label: 'root',
            month: '4/24',
            timestamp: '2024-04-01T00:00:00Z',
            mounts: [
              {
                acme_clients: 0,
                clients: 2,
                entity_clients: 2,
                label: 'no mount accessor (pre-1.10 upgrade?)',
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                acme_clients: 0,
                clients: 1,
                entity_clients: 1,
                label: 'auth/u/',
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            ],
            non_entity_clients: 0,
            secret_syncs: 0,
          },
          non_entity_clients: 0,
          secret_syncs: 0,
          timestamp: '2024-04-01T00:00:00Z',
        },
      },
      'it formats combined data for monthly namespaces_by_key spanning upgrade to 1.10'
    );
  });

  test('setStartTimeQuery: it returns start time query for activity log', async function (assert) {
    assert.expect(6);
    const apiPath = 'sys/internal/counters/config';
    assert.strictEqual(setStartTimeQuery(true, {}), null, `it returns null if no permission to ${apiPath}`);
    assert.strictEqual(
      setStartTimeQuery(false, {}),
      null,
      `it returns null for community edition and no permission to ${apiPath}`
    );
    assert.strictEqual(
      setStartTimeQuery(true, { billingStartTimestamp: new Date('2022-06-08T00:00:00Z') }),
      1654646400,
      'it returns unix time if enterprise and billing_start_timestamp exists'
    );
    assert.strictEqual(
      setStartTimeQuery(false, { billingStartTimestamp: new Date('0001-01-01T00:00:00Z') }),
      null,
      'it returns null time for community edition even if billing_start_timestamp exists'
    );
    assert.strictEqual(
      setStartTimeQuery(false, { foo: 'bar' }),
      null,
      'it returns null if billing_start_timestamp key does not exist'
    );
    assert.strictEqual(
      setStartTimeQuery(false, undefined),
      null,
      'fails gracefully if no config model is passed'
    );
  });
});
