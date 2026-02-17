/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import {
  destructureClientCounts,
  formatByMonths,
  formatByNamespace,
  formatQueryParams,
  sortMonthsByTimestamp,
} from 'core/utils/client-counts/serializers';
import {
  ACTIVITY_RESPONSE_STUB as RESPONSE,
  MIXED_ACTIVITY_RESPONSE_STUB as MIXED_RESPONSE,
  SERIALIZED_ACTIVITY_RESPONSE,
} from 'vault/tests/helpers/clients/client-count-helpers';

module('Unit | Util | client counts | serializers', function (hooks) {
  setupTest(hooks);

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
    const expectedNs1 = SERIALIZED_ACTIVITY_RESPONSE.by_namespace.find((ns) => ns.label === 'ns1/');
    const formattedNs1 = formatByNamespace(RESPONSE.by_namespace).find((ns) => ns.label === 'ns1/');
    assert.expect(2 + formattedNs1.mounts.length);

    assert.propEqual(formattedNs1, expectedNs1, 'it formats ns1/ namespace');
    assert.propEqual(RESPONSE.by_namespace, original, 'it does not modify original by_namespace array');

    formattedNs1.mounts.forEach((mount) => {
      const expectedMount = expectedNs1.mounts.find((m) => m.label === mount.label);
      assert.propEqual(mount, expectedMount, `${mount.label} has expected key/value pairs`);
    });
  });

  module('formatQueryParams', function (hooks) {
    hooks.beforeEach(function () {
      this.assertQuery = ({ query, expected }, assert) => {
        const result = formatQueryParams(query);

        assert.propEqual(
          result,
          expected,
          `returned params: ${JSON.stringify(result)} matches expected: ${JSON.stringify(expected)}`
        );
        assert.strictEqual(result?.start_time, expected?.start_time, 'query has expected start_time');
        assert.strictEqual(result?.end_time, expected?.end_time, 'query has expected end_time');
      };
    });

    test('formatQueryParams: it returns formatted query params with valid ISO date strings', function (assert) {
      const query = { start_time: '2023-01-01T00:00:00.000Z', end_time: '2023-12-31T23:59:59.999Z' };
      this.assertQuery({ query, expected: query }, assert);
    });

    test('it returns undefined for invalid date strings', function (assert) {
      const query = { start_time: 'invalid-date', end_time: 'not-a-date' };
      const expected = {};
      this.assertQuery({ query, expected }, assert);
    });

    test('it handles mixed valid and invalid dates', function (assert) {
      const query = { start_time: '2023-01-01T00:00:00.000Z', end_time: 'invalid-date' };
      const expected = { start_time: '2023-01-01T00:00:00.000Z' };
      this.assertQuery({ query, expected }, assert);
    });

    test('it handles empty strings', function (assert) {
      let query = { start_time: '', end_time: '' };
      let expected = {};
      this.assertQuery({ query, expected }, assert);

      query = { start_time: '2023-01-01T00:00:00.000Z', end_time: '' };
      expected = { start_time: '2023-01-01T00:00:00.000Z' };
      this.assertQuery({ query, expected }, assert);

      query = { start_time: '', end_time: '2023-12-31T23:59:59.999Z' };
      expected = { end_time: '2023-12-31T23:59:59.999Z' };
      this.assertQuery({ query, expected }, assert);
    });

    test('it handles undefined values', function (assert) {
      let query = { start_time: undefined, end_time: undefined };
      let expected = {};
      this.assertQuery({ query, expected }, assert);

      query = { start_time: '2023-01-01T00:00:00.000Z', end_time: undefined };
      expected = { start_time: '2023-01-01T00:00:00.000Z' };
      this.assertQuery({ query, expected }, assert);

      query = { start_time: undefined, end_time: '2023-12-31T23:59:59.999Z' };
      expected = { end_time: '2023-12-31T23:59:59.999Z' };
      this.assertQuery({ query, expected }, assert);
    });

    test('it handles null values', function (assert) {
      let query = { start_time: null, end_time: null };
      let expected = {};
      this.assertQuery({ query, expected }, assert);

      query = { start_time: '2023-01-01T00:00:00.000Z', end_time: null };
      expected = { start_time: '2023-01-01T00:00:00.000Z' };
      this.assertQuery({ query, expected }, assert);

      query = { start_time: null, end_time: '2023-12-31T23:59:59.999Z' };
      expected = { end_time: '2023-12-31T23:59:59.999Z' };
      this.assertQuery({ query, expected }, assert);
    });

    test('it handles missing properties', function (assert) {
      const query = {};
      const expected = {};
      this.assertQuery({ query, expected }, assert);
    });

    test('it does not accept date strings that are not ISO formatted', function (assert) {
      const query = { start_time: '2023-06-15', end_time: '2023-06-16T14:30:00Z' };
      const expected = { end_time: '2023-06-16T14:30:00Z' };
      this.assertQuery({ query, expected }, assert);
    });
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
            mount_type: '',
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
            mount_path: 'no mount accessor (pre-1.10 upgrade?)',
            mount_type: '',
            namespace_path: 'root',
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
              mount_path: 'no mount accessor (pre-1.10 upgrade?)',
              mount_type: 'no mount path (pre-1.10 upgrade?)',
              namespace_path: 'root',
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            {
              acme_clients: 0,
              clients: 1,
              entity_clients: 1,
              label: 'auth/userpass/0/',
              mount_path: 'auth/userpass/0/',
              mount_type: 'userpass',
              namespace_path: 'root',
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
                mount_path: 'no mount accessor (pre-1.10 upgrade?)',
                mount_type: 'no mount path (pre-1.10 upgrade?)',
                namespace_path: 'root',
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                acme_clients: 0,
                clients: 1,
                entity_clients: 1,
                label: 'auth/userpass/0/',
                mount_path: 'auth/userpass/0/',
                mount_type: 'userpass',
                namespace_path: 'root',
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
});
