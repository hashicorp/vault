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
  destructureClientCounts,
  namespaceArrayToObject,
  sortMonthsByTimestamp,
} from 'core/utils/client-count-utils';
import { LICENSE_START } from 'vault/mirage/handlers/clients';
import {
  ACTIVITY_RESPONSE_STUB as RESPONSE,
  VERSION_HISTORY,
  SERIALIZED_ACTIVITY_RESPONSE,
} from 'vault/tests/helpers/clients';

/*
formatByNamespace, formatByMonths, destructureClientCounts are utils 
used to normalize the sys/counters/activity response in the clients/activity 
serializer. these functions are tested individually here, instead of all at once 
in a serializer test for easier debugging
*/
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

    const [formattedNoData, formattedWithActivity] = formatByMonths(RESPONSE.months);

    // instead of asserting the whole expected response, broken up so tests are easier to debug
    // but kept whole above to copy/paste updated response expectations in the future
    const [expectedNoData, expectedWithActivity] = SERIALIZED_ACTIVITY_RESPONSE.by_month;
    const { namespaces, new_clients } = expectedWithActivity;

    assert.propEqual(formattedNoData, expectedNoData, 'it formats months without data');
    assert.propEqual(
      formattedWithActivity.namespaces,
      namespaces,
      'it formats namespaces array for months with data'
    );
    assert.propEqual(
      formattedWithActivity.new_clients,
      new_clients,
      'it formats new_clients block for months with data'
    );
    assert.propEqual(RESPONSE.months, original, 'it does not modify original months array');
    assert.propEqual(formatByMonths([]), [], 'it returns an empty array if the months key is empty');
  });

  test('formatByNamespace: formats namespace array with mounts', async function (assert) {
    assert.expect(3);
    const original = [...RESPONSE.by_namespace];
    const [formattedRoot, formattedNs1] = formatByNamespace(RESPONSE.by_namespace);
    const [root, ns1] = SERIALIZED_ACTIVITY_RESPONSE.by_namespace;

    assert.propEqual(formattedRoot, root, 'it formats root namespace');
    assert.propEqual(formattedNs1, ns1, 'it formats ns1/ namespace');
    assert.propEqual(RESPONSE.by_namespace, original, 'it does not modify original by_namespace array');
  });

  test('formatByNamespace: formats namespace array with no mounts (activity log data < 1.10)', async function (assert) {
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
        mounts: 'no mount accessor (pre-1.10 upgrade?)',
      },
    ];
    const expected = [
      {
        acme_clients: 0,
        clients: 30,
        entity_clients: 10,
        label: 'root',
        mounts: [],
        non_entity_clients: 20,
        secret_syncs: 0,
      },
    ];
    assert.propEqual(formatByNamespace(noMounts), expected, 'it formats namespace without mounts');
    assert.propEqual(formatByNamespace([]), [], 'it returns an empty array if the by_namespace key is empty');
  });

  test('destructureClientCounts: homogenizes key names when both old and new keys exist, or just old key names', async function (assert) {
    assert.expect(2);
    const original = { ...RESPONSE.total };
    const expected = {
      entity_clients: 1816,
      non_entity_clients: 3117,
      secret_syncs: 2672,
      acme_clients: 200,
      clients: 7805,
    };
    assert.propEqual(destructureClientCounts(RESPONSE.total), expected);
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

  test('namespaceArrayToObject: it generates namespaces_by_key without modifying original', async function (assert) {
    assert.expect(5);

    // month at 0-index has no data so use second month in array, empty month format covered by formatByMonths test above
    const original = { ...RESPONSE.months[1] };
    const expectedObject = SERIALIZED_ACTIVITY_RESPONSE.by_month[1].namespaces_by_key;
    const formattedTotal = formatByNamespace(RESPONSE.months[1].namespaces);

    const testObject = namespaceArrayToObject(
      formattedTotal,
      formatByNamespace(RESPONSE.months[1].new_clients.namespaces),
      '9/23',
      '2023-09-01T00:00:00Z'
    );

    const { root } = testObject;
    const { root: expectedRoot } = expectedObject;
    assert.propEqual(root.new_clients, expectedRoot.new_clients, 'it formats namespaces new_clients');
    assert.propEqual(root.mounts_by_key, expectedRoot.mounts_by_key, 'it formats namespaces mounts_by_key');
    assert.propContains(root, expectedRoot, 'namespace has correct keys');

    assert.propEqual(
      namespaceArrayToObject(formattedTotal, formatByNamespace([]), '9/23', '2023-09-01T00:00:00Z'),
      {},
      'returns an empty object when there are no new clients '
    );
    assert.propEqual(RESPONSE.months[1], original, 'it does not modify original month data');
  });
});
