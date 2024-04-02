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
import { ACTIVITY_RESPONSE as RESPONSE, VERSION_HISTORY, EXPECTED_FORMAT } from 'vault/tests/helpers/clients';

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
    assert.expect(4);
    const original = [...RESPONSE.months];

    const [formattedNoData, formattedWithActivity] = formatByMonths(RESPONSE.months);

    // instead of asserting the whole expected response, broken up so tests are easier to debug
    // but kept whole above to copy/paste updated response expectations in the future
    const [expectedNoData, expectedWithActivity] = EXPECTED_FORMAT;
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
  });

  test('formatByNamespace: formats namespace array with mounts', async function (assert) {
    assert.expect(3);
    const original = [...RESPONSE.by_namespace];
    const expected = [
      {
        label: 'root',
        clients: 5429,
        entity_clients: 1033,
        non_entity_clients: 1924,
        secret_syncs: 2397,
        acme_clients: 75,
        mounts: [
          {
            acme_clients: 0,
            clients: 2957,
            entity_clients: 1033,
            label: 'auth/authid0',
            non_entity_clients: 1924,
            secret_syncs: 0,
          },
          {
            acme_clients: 0,
            clients: 2397,
            entity_clients: 0,
            label: 'kvv2-engine-0',
            non_entity_clients: 0,
            secret_syncs: 2397,
          },
          {
            acme_clients: 75,
            clients: 75,
            entity_clients: 0,
            label: 'pki-engine-0',
            non_entity_clients: 0,
            secret_syncs: 0,
          },
        ],
      },
      {
        label: 'ns/1',
        clients: 2376,
        entity_clients: 783,
        non_entity_clients: 1193,
        secret_syncs: 275,
        acme_clients: 125,
        mounts: [
          {
            label: 'auth/authid0',
            clients: 1976,
            entity_clients: 783,
            non_entity_clients: 1193,
            secret_syncs: 0,
            acme_clients: 0,
          },
          {
            label: 'kvv2-engine-0',
            clients: 275,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 275,
            acme_clients: 0,
          },
          {
            label: 'pki-engine-0',
            acme_clients: 125,
            clients: 125,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 0,
          },
        ],
      },
    ];
    const [formattedRoot, formattedNs1] = formatByNamespace(RESPONSE.by_namespace);
    const [root, ns1] = expected;

    assert.propEqual(formattedRoot, root, 'it formats root namespace');
    assert.propEqual(formattedNs1, ns1, 'it formats ns1/ namespace');
    assert.propEqual(RESPONSE.by_namespace, original, 'it does not modify original by_namespace array');
  });

  test('formatByNamespace: formats namespace array with no mounts (activity log data < 1.10)', async function (assert) {
    assert.expect(1);
    // TODO waiting to hear from backend whether the mounts key will actually exist or not
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
  });

  test('homogenizeClientNaming: homogenizes key names when both old and new keys exist, or just old key names', async function (assert) {
    assert.expect(2);
    const original = { ...RESPONSE.total };
    const expected = {
      entity_clients: 1816,
      non_entity_clients: 3117,
      secret_syncs: 2672,
      acme_clients: 200,
      clients: 7805,
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

  test('namespaceArrayToObject: it generates namespaces_by_key without modifying original', async function (assert) {
    assert.expect(3);
    const { namespaces_by_key: expected } = EXPECTED_FORMAT[1];

    const { namespaces, new_clients } = RESPONSE.months[1];
    const original = { ...RESPONSE.months[1] };
    const byNamespaceKeyObject = namespaceArrayToObject(
      formatByNamespace(namespaces),
      formatByNamespace(new_clients.namespaces),
      '9/23',
      '2023-09-01T00:00:00Z'
    );

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
    assert.propEqual(RESPONSE.months[1], original, 'it does not modify original month data');
  });
});
