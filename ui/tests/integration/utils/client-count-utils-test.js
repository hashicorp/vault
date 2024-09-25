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
  sortMonthsByTimestamp,
  filterByMonthDataForMount,
  filteredTotalForMount,
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
          timestampInstalled: '2023-08-02T00:00:00Z',
          version: '1.9.1',
        },
        {
          previousVersion: '1.9.1',
          timestampInstalled: '2023-09-02T00:00:00Z',
          version: '1.10.1',
        },
        {
          previousVersion: '1.16.0',
          timestampInstalled: '2023-12-02T00:00:00Z',
          version: '1.17.0',
        },
      ];
      // set start/end times longer than version history to test all relevant upgrades return
      const startTime = '2023-06-02T00:00:00Z'; // first upgrade installed '2023-07-02T00:00:00Z'
      const endTime = '2024-03-04T16:14:21Z'; // latest upgrade installed '2023-12-02T00:00:00Z'
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
          timestampInstalled: '2023-08-02T00:00:00Z',
          version: '1.9.1',
        },
        {
          previousVersion: '1.9.1',
          timestampInstalled: '2023-09-02T00:00:00Z',
          version: '1.10.1',
        },
      ];
      const startTime = '2023-08-02T00:00:00Z'; // same date as 1.9.1 install date to catch same day edge cases
      const endTime = '2023-11-02T00:00:00Z';
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
          timestampInstalled: '2023-09-23T00:00:00Z',
        },
        'it does not return subsequent patch versions of the same notable upgrade version'
      );
    });
  });

  test('formatByMonths: it formats the months array', async function (assert) {
    assert.expect(7);
    const original = [...RESPONSE.months];

    const [formattedNoData, formattedWithActivity, formattedNoNew] = formatByMonths(RESPONSE.months);

    // instead of asserting the whole expected response, broken up so tests are easier to debug
    // but kept whole above to copy/paste updated response expectations in the future
    const [expectedNoData, expectedWithActivity, expectedNoNew] = SERIALIZED_ACTIVITY_RESPONSE.by_month;

    assert.propEqual(formattedNoData, expectedNoData, 'it formats months without data');
    ['namespaces', 'new_clients'].forEach((key) => {
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

  // TESTS FOR COMBINED ACTIVITY DATA - no mount attribution < 1.10
  test('it formats the namespaces array with no mount attribution (activity log data < 1.10)', async function (assert) {
    assert.expect(2);
    const noMounts = [
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          entity_clients: 10,
          non_entity_clients: 20,
          secret_syncs: 0,
          acme_clients: 0,
          clients: 30,
        },
        mounts: [
          {
            counts: {
              entity_clients: 10,
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
    assert.expect(2);

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
  });

  module('filterByMonthDataForMount', function (hooks) {
    hooks.beforeEach(function () {
      this.getExpected = (label, count = 0, newCount = 0) => {
        return {
          month: '6/23',
          namespaces: [],
          label,
          timestamp: '2023-06-01T00:00:00Z',
          acme_clients: count,
          clients: count,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 0,
          new_clients: {
            month: '6/23',
            timestamp: '2023-06-01T00:00:00Z',
            namespaces: [],
            label,
            acme_clients: newCount,
            clients: newCount,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 0,
          },
        };
      };
    });

    test('it works when month has no data', async function (assert) {
      const months = [
        {
          month: '6/23',
          timestamp: '2023-06-01T00:00:00Z',
          namespaces: [],
          new_clients: {
            month: '6/23',
            timestamp: '2023-06-01T00:00:00Z',
            namespaces: [],
          },
        },
      ];
      const result = filterByMonthDataForMount(months, 'root', 'some-mount');
      // no data is different than zero, it implies no data was being saved at that time
      // so we don't fill in missing data with zeros to differentiate those two states
      assert.deepEqual(result[0], months[0], 'does not change month when no data');
    });

    test('it works when month has no new clients', async function (assert) {
      const months = [
        {
          month: '6/23',
          timestamp: '2023-06-01T00:00:00Z',
          acme_clients: 11,
          clients: 11,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 0,
          namespaces: [
            {
              label: 'root',
              acme_clients: 11,
              clients: 11,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
              mounts: [
                {
                  label: 'some-mount',
                  acme_clients: 11,
                  clients: 11,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
              ],
            },
          ],
          new_clients: {
            month: '6/23',
            timestamp: '2023-06-01T00:00:00Z',
            namespaces: [],
          },
        },
      ];

      let result = filterByMonthDataForMount(months, 'root', 'some-mount');
      assert.propEqual(result[0], this.getExpected('some-mount', 11), 'works when mount is found');
      result = filterByMonthDataForMount(months, 'root', 'another-mount');
      assert.deepEqual(result[0], this.getExpected('another-mount', 0), 'works when mount is not found');
      result = filterByMonthDataForMount(months, 'unknown-child', 'some-mount');
      assert.deepEqual(result[0], this.getExpected('some-mount', 0), 'works when namespace is not found');
    });

    test('it works when month has new clients', async function (assert) {
      const months = [
        {
          month: '6/23',
          timestamp: '2023-06-01T00:00:00Z',
          acme_clients: 22,
          clients: 22,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 0,
          namespaces: [
            {
              label: 'root',
              acme_clients: 22,
              clients: 22,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
              mounts: [
                {
                  label: 'some-mount',
                  acme_clients: 22,
                  clients: 22,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
              ],
            },
          ],
          new_clients: {
            month: '6/23',
            timestamp: '2023-06-01T00:00:00Z',
            namespaces: [
              {
                label: 'root',
                acme_clients: 11,
                clients: 11,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
                mounts: [
                  {
                    label: 'some-mount',
                    acme_clients: 11,
                    clients: 11,
                    entity_clients: 0,
                    non_entity_clients: 0,
                    secret_syncs: 0,
                  },
                ],
              },
            ],
          },
        },
      ];
      let result = filterByMonthDataForMount(months, 'root', 'some-mount');
      assert.propEqual(result[0], this.getExpected('some-mount', 22, 11), 'works when mount is found');
      result = filterByMonthDataForMount(months, 'root', 'another-mount');
      assert.deepEqual(result[0], this.getExpected('another-mount', 0), 'works when mount is not found');
      result = filterByMonthDataForMount(months, 'unknown-child', 'some-mount');
      assert.deepEqual(result[0], this.getExpected('some-mount', 0), 'works when namespace is not found');
    });
  });

  module('filteredTotalForMount', function (hooks) {
    hooks.beforeEach(function () {
      this.byNamespace = SERIALIZED_ACTIVITY_RESPONSE.by_namespace;
    });

    const emptyCounts = {
      acme_clients: 0,
      clients: 0,
      entity_clients: 0,
      non_entity_clients: 0,
      secret_syncs: 0,
    };

    [
      {
        when: 'no namespace filter passed',
        result: 'it returns empty counts',
        ns: '',
        mount: 'auth/authid/0',
        expected: emptyCounts,
      },
      {
        when: 'no mount filter passed',
        result: 'it returns empty counts',
        ns: 'ns1',
        mount: '',
        expected: emptyCounts,
      },
      {
        when: 'no matching ns/mount exists',
        result: 'it returns empty counts',
        ns: 'ns1',
        mount: 'auth/authid/1',
        expected: emptyCounts,
      },
      {
        when: 'mount and label have extra slashes',
        result: 'it returns the data sanitized',
        ns: 'ns1/',
        mount: 'auth/authid/0/',
        expected: {
          label: 'auth/authid/0',
          acme_clients: 0,
          clients: 8394,
          entity_clients: 4256,
          non_entity_clients: 4138,
          secret_syncs: 0,
        },
      },
      {
        when: 'mount within root',
        result: 'it returns the data',
        ns: 'root',
        mount: 'kvv2-engine-0',
        expected: {
          label: 'kvv2-engine-0',
          acme_clients: 0,
          clients: 4290,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 4290,
        },
      },
    ].forEach((testCase) => {
      test(`it returns correct values when ${testCase.when}`, async function (assert) {
        const actual = filteredTotalForMount(this.byNamespace, testCase.ns, testCase.mount);
        assert.deepEqual(actual, testCase.expected);
      });
    });
  });
});
