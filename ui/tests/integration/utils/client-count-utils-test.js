/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import {
  filterVersionHistory,
  flattenDataset,
  formatByMonths,
  formatByNamespace,
  destructureCounts,
  namespaceArrayToObject,
  sortMonthsByTimestamp,
} from 'core/utils/client-count-utils';
import { LICENSE_START } from 'vault/mirage/handlers/clients';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { addMonths, isAfter, isBefore } from 'date-fns';

const MONTHS = [
  {
    timestamp: '2021-05-01T00:00:00Z',
    counts: {
      distinct_entities: 25,
      non_entity_tokens: 25,
      clients: 50,
    },
    namespaces: [
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 13,
          non_entity_tokens: 7,
          clients: 20,
        },
        mounts: [
          {
            mount_path: 'auth/up2/',
            counts: {
              distinct_entities: 8,
              non_entity_tokens: 0,
              clients: 8,
            },
          },
          {
            mount_path: 'auth/up1/',
            counts: {
              distinct_entities: 0,
              non_entity_tokens: 7,
              clients: 7,
            },
          },
        ],
      },
      {
        namespace_id: 's07UR',
        namespace_path: 'ns1/',
        counts: {
          distinct_entities: 5,
          non_entity_tokens: 5,
          clients: 10,
        },
        mounts: [
          {
            mount_path: 'auth/up1/',
            counts: {
              distinct_entities: 0,
              non_entity_tokens: 5,
              clients: 5,
            },
          },
          {
            mount_path: 'auth/up2/',
            counts: {
              distinct_entities: 5,
              non_entity_tokens: 0,
              clients: 5,
            },
          },
        ],
      },
    ],
    new_clients: {
      counts: {
        distinct_entities: 3,
        non_entity_tokens: 2,
        clients: 5,
      },
      namespaces: [
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 3,
            non_entity_tokens: 2,
            clients: 5,
          },
          mounts: [
            {
              mount_path: 'auth/up2/',
              counts: {
                distinct_entities: 3,
                non_entity_tokens: 0,
                clients: 3,
              },
            },
            {
              mount_path: 'auth/up1/',
              counts: {
                distinct_entities: 0,
                non_entity_tokens: 2,
                clients: 2,
              },
            },
          ],
        },
      ],
    },
  },
  {
    timestamp: '2021-10-01T00:00:00Z',
    counts: {
      distinct_entities: 20,
      entity_clients: 20,
      non_entity_tokens: 20,
      non_entity_clients: 20,
      clients: 40,
    },
    namespaces: [
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 8,
          entity_clients: 8,
          non_entity_tokens: 7,
          non_entity_clients: 7,
          clients: 15,
        },
        mounts: [
          {
            mount_path: 'auth/up2/',
            counts: {
              distinct_entities: 8,
              entity_clients: 8,
              non_entity_tokens: 0,
              non_entity_clients: 0,
              clients: 8,
            },
          },
          {
            mount_path: 'auth/up1/',
            counts: {
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_tokens: 7,
              non_entity_clients: 7,
              clients: 7,
            },
          },
        ],
      },
      {
        namespace_id: 's07UR',
        namespace_path: 'ns1/',
        counts: {
          distinct_entities: 5,
          entity_clients: 5,
          non_entity_tokens: 5,
          non_entity_clients: 5,
          clients: 10,
        },
        mounts: [
          {
            mount_path: 'auth/up1/',
            counts: {
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_tokens: 5,
              non_entity_clients: 5,
              clients: 5,
            },
          },
          {
            mount_path: 'auth/up2/',
            counts: {
              distinct_entities: 5,
              entity_clients: 5,
              non_entity_tokens: 0,
              non_entity_clients: 0,
              clients: 5,
            },
          },
        ],
      },
    ],
    new_clients: {
      counts: {
        distinct_entities: 3,
        entity_clients: 3,
        non_entity_tokens: 2,
        non_entity_clients: 2,
        clients: 5,
      },
      namespaces: [
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 3,
            entity_clients: 3,
            non_entity_tokens: 2,
            non_entity_clients: 2,
            clients: 5,
          },
          mounts: [
            {
              mount_path: 'auth/up2/',
              counts: {
                distinct_entities: 3,
                entity_clients: 3,
                non_entity_tokens: 0,
                non_entity_clients: 0,
                clients: 3,
              },
            },
            {
              mount_path: 'auth/up1/',
              counts: {
                distinct_entities: 0,
                entity_clients: 0,
                non_entity_tokens: 2,
                non_entity_clients: 2,
                clients: 2,
              },
            },
          ],
        },
      ],
    },
  },
  {
    timestamp: '2021-09-01T00:00:00Z',
    counts: {
      distinct_entities: 0,
      entity_clients: 17,
      non_entity_tokens: 0,
      non_entity_clients: 18,
      clients: 35,
    },
    namespaces: [
      {
        namespace_id: 'oImjk',
        namespace_path: 'ns2/',
        counts: {
          distinct_entities: 0,
          entity_clients: 5,
          non_entity_tokens: 0,
          non_entity_clients: 5,
          clients: 10,
        },
        mounts: [
          {
            mount_path: 'auth/up1/',
            counts: {
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_tokens: 0,
              non_entity_clients: 5,
              clients: 5,
            },
          },
          {
            mount_path: 'auth/up2/',
            counts: {
              distinct_entities: 0,
              entity_clients: 5,
              non_entity_tokens: 0,
              non_entity_clients: 0,
              clients: 5,
            },
          },
        ],
      },
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 0,
          entity_clients: 2,
          non_entity_tokens: 0,
          non_entity_clients: 3,
          clients: 5,
        },
        mounts: [
          {
            mount_path: 'auth/up1/',
            counts: {
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_tokens: 0,
              non_entity_clients: 3,
              clients: 3,
            },
          },
          {
            mount_path: 'auth/up2/',
            counts: {
              distinct_entities: 0,
              entity_clients: 2,
              non_entity_tokens: 0,
              non_entity_clients: 0,
              clients: 2,
            },
          },
        ],
      },
      {
        namespace_id: 's07UR',
        namespace_path: 'ns1/',
        counts: {
          distinct_entities: 0,
          entity_clients: 3,
          non_entity_tokens: 0,
          non_entity_clients: 2,
          clients: 5,
        },
        mounts: [
          {
            mount_path: 'auth/up2/',
            counts: {
              distinct_entities: 0,
              entity_clients: 3,
              non_entity_tokens: 0,
              non_entity_clients: 0,
              clients: 3,
            },
          },
          {
            mount_path: 'auth/up1/',
            counts: {
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_tokens: 0,
              non_entity_clients: 2,
              clients: 2,
            },
          },
        ],
      },
    ],
    new_clients: {
      counts: {
        distinct_entities: 0,
        entity_clients: 10,
        non_entity_tokens: 0,
        non_entity_clients: 10,
        clients: 20,
      },
      namespaces: [
        {
          namespace_id: 'oImjk',
          namespace_path: 'ns2/',
          counts: {
            distinct_entities: 0,
            entity_clients: 5,
            non_entity_tokens: 0,
            non_entity_clients: 5,
            clients: 10,
          },
          mounts: [
            {
              mount_path: 'auth/up1/',
              counts: {
                distinct_entities: 0,
                entity_clients: 0,
                non_entity_tokens: 0,
                non_entity_clients: 5,
                clients: 5,
              },
            },
            {
              mount_path: 'auth/up2/',
              counts: {
                distinct_entities: 0,
                entity_clients: 5,
                non_entity_tokens: 0,
                non_entity_clients: 0,
                clients: 5,
              },
            },
          ],
        },
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 0,
            entity_clients: 2,
            non_entity_tokens: 0,
            non_entity_clients: 3,
            clients: 5,
          },
          mounts: [
            {
              mount_path: 'auth/up1/',
              counts: {
                distinct_entities: 0,
                entity_clients: 0,
                non_entity_tokens: 0,
                non_entity_clients: 3,
                clients: 3,
              },
            },
            {
              mount_path: 'auth/up2/',
              counts: {
                distinct_entities: 0,
                entity_clients: 2,
                non_entity_tokens: 0,
                non_entity_clients: 0,
                clients: 2,
              },
            },
          ],
        },
        {
          namespace_id: 's07UR',
          namespace_path: 'ns1/',
          counts: {
            distinct_entities: 0,
            entity_clients: 3,
            non_entity_tokens: 0,
            non_entity_clients: 2,
            clients: 5,
          },
          mounts: [
            {
              mount_path: 'auth/up2/',
              counts: {
                distinct_entities: 0,
                entity_clients: 3,
                non_entity_tokens: 0,
                non_entity_clients: 0,
                clients: 3,
              },
            },
            {
              mount_path: 'auth/up1/',
              counts: {
                distinct_entities: 0,
                entity_clients: 0,
                non_entity_tokens: 0,
                non_entity_clients: 2,
                clients: 2,
              },
            },
          ],
        },
      ],
    },
  },
];

const BY_NAMESPACE = [
  {
    namespace_id: '96OwG',
    namespace_path: 'test-ns/',
    counts: {
      distinct_entities: 18290,
      entity_clients: 18290,
      non_entity_tokens: 18738,
      non_entity_clients: 18738,
      clients: 37028,
    },
    mounts: [
      {
        mount_path: 'path-1',
        counts: {
          distinct_entities: 6403,
          entity_clients: 6403,
          non_entity_tokens: 6300,
          non_entity_clients: 6300,
          clients: 12703,
        },
      },
      {
        mount_path: 'path-2',
        counts: {
          distinct_entities: 5699,
          entity_clients: 5699,
          non_entity_tokens: 6777,
          non_entity_clients: 6777,
          clients: 12476,
        },
      },
      {
        mount_path: 'path-3',
        counts: {
          distinct_entities: 6188,
          entity_clients: 6188,
          non_entity_tokens: 5661,
          non_entity_clients: 5661,
          clients: 11849,
        },
      },
    ],
  },
  {
    namespace_id: 'root',
    namespace_path: '',
    counts: {
      distinct_entities: 19099,
      entity_clients: 19099,
      non_entity_tokens: 17781,
      non_entity_clients: 17781,
      clients: 36880,
    },
    mounts: [
      {
        mount_path: 'path-3',
        counts: {
          distinct_entities: 6863,
          entity_clients: 6863,
          non_entity_tokens: 6801,
          non_entity_clients: 6801,
          clients: 13664,
        },
      },
      {
        mount_path: 'path-2',
        counts: {
          distinct_entities: 6047,
          entity_clients: 6047,
          non_entity_tokens: 5957,
          non_entity_clients: 5957,
          clients: 12004,
        },
      },
      {
        mount_path: 'path-1',
        counts: {
          distinct_entities: 6189,
          entity_clients: 6189,
          non_entity_tokens: 5023,
          non_entity_clients: 5023,
          clients: 11212,
        },
      },
      {
        mount_path: 'auth/up2/',
        counts: {
          distinct_entities: 0,
          entity_clients: 50,
          non_entity_tokens: 0,
          non_entity_clients: 23,
          clients: 73,
        },
      },
      {
        mount_path: 'auth/up1/',
        counts: {
          distinct_entities: 0,
          entity_clients: 25,
          non_entity_tokens: 0,
          non_entity_clients: 15,
          clients: 40,
        },
      },
    ],
  },
];

const EMPTY_MONTHS = [
  {
    timestamp: '2021-06-01T00:00:00Z',
    counts: null,
    namespaces: null,
    new_clients: null,
  },
  {
    timestamp: '2021-07-01T00:00:00Z',
    counts: null,
    namespaces: null,
    new_clients: null,
  },
];

const SOME_OBJECT = { foo: 'bar' };

module('Integration | Util | client count utils', function (hooks) {
  setupTest(hooks);

  test('filterVersionHistory: returns version data that occurred during activity date range', async function (assert) {
    assert.expect(1);
    // LICENSE_START: '2023-07-02T00:00:00Z'
    const versionHistory = [
      {
        version: '1.9.0',
        previousVersion: null,
        timestampInstalled: LICENSE_START.toISOString(),
      },
      {
        version: '1.9.1',
        previousVersion: '1.9.0',
        timestampInstalled: addMonths(LICENSE_START, 1).toISOString(),
      },
      {
        version: '1.10.1',
        previousVersion: '1.9.1',
        timestampInstalled: addMonths(LICENSE_START, 2).toISOString(),
      },
      {
        version: '1.14.4',
        previousVersion: '1.10.1',
        timestampInstalled: addMonths(LICENSE_START, 3).toISOString(),
      },
      {
        version: '1.16.0',
        previousVersion: '1.14.4',
        timestampInstalled: addMonths(LICENSE_START, 4).toISOString(),
      },
    ];
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
    const activity = {
      startTime: '2023-07-02T00:00:00Z', // same as license start to catch same day edge cases
      endTime: '2024-03-04T16:14:21.000Z',
    };
    assert.propEqual(
      filterVersionHistory(versionHistory, activity.startTime, activity.endTime),
      expected,
      'it returns only upgrades that happened between given start and end times.'
    );
  });

  test('formatByMonths: formats the months array', async function (assert) {
    assert.expect(103);
    const keyNameAssertions = (object, objectName) => {
      const objectKeys = Object.keys(object);
      assert.false(objectKeys.includes('counts'), `${objectName} doesn't include 'counts' key`);
      assert.true(objectKeys.includes('clients'), `${objectName} includes 'clients' key`);
      assert.true(objectKeys.includes('entity_clients'), `${objectName} includes 'entity_clients' key`);
      assert.true(
        objectKeys.includes('non_entity_clients'),
        `${objectName} includes 'non_entity_clients' key`
      );
    };
    const assertClientCounts = (object, originalObject) => {
      const newObjectKeys = ['clients', 'entity_clients', 'non_entity_clients'];
      const originalKeys = Object.keys(originalObject.counts).includes('entity_clients')
        ? newObjectKeys
        : ['clients', 'distinct_entities', 'non_entity_tokens'];

      newObjectKeys.forEach((key, i) => {
        assert.strictEqual(
          object[key],
          originalObject.counts[originalKeys[i]],
          `${object.month} ${key} equal original counts`
        );
      });
    };

    const formattedMonths = formatByMonths(MONTHS);
    assert.notEqual(formattedMonths, MONTHS, 'does not modify original array');

    formattedMonths.forEach((month) => {
      const originalMonth = MONTHS.find((m) => month.month === parseAPITimestamp(m.timestamp, 'M/yy'));
      // if originalMonth is found (not undefined) then the formatted month has an accurate, parsed timestamp
      assert.ok(originalMonth, `month has parsed timestamp of ${month.month}`);
      assert.ok(month.namespaces_by_key, `month includes 'namespaces_by_key' key`);

      keyNameAssertions(month, 'formatted month');
      assertClientCounts(month, originalMonth);

      assert.ok(month.new_clients.month, 'new clients key has a month key');
      keyNameAssertions(month.new_clients, 'formatted month new_clients');
      assertClientCounts(month.new_clients, originalMonth.new_clients);

      month.namespaces.forEach((namespace) => keyNameAssertions(namespace, 'namespace within month'));
      month.new_clients.namespaces.forEach((namespace) =>
        keyNameAssertions(namespace, 'new client namespaces within month')
      );
    });

    // method fails gracefully
    const expected = [
      {
        counts: null,
        month: '6/21',
        namespaces: [],
        namespaces_by_key: {},
        new_clients: {
          month: '6/21',
          namespaces: [],
          timestamp: '2021-06-01T00:00:00Z',
        },
        timestamp: '2021-06-01T00:00:00Z',
      },
      {
        counts: null,
        month: '7/21',
        namespaces: [],
        namespaces_by_key: {},
        new_clients: {
          month: '7/21',
          namespaces: [],
          timestamp: '2021-07-01T00:00:00Z',
        },
        timestamp: '2021-07-01T00:00:00Z',
      },
    ];
    assert.strictEqual(formatByMonths(SOME_OBJECT), SOME_OBJECT, 'it returns if arg is not an array');
    assert.propEqual(formatByMonths(EMPTY_MONTHS), expected, 'it does not error with null months');
    assert.ok(formatByMonths([...EMPTY_MONTHS, ...MONTHS]), 'it does not error with combined data');
  });

  test('formatByNamespace: formats namespace arrays with and without mounts', async function (assert) {
    assert.expect(102);
    const keyNameAssertions = (object, objectName) => {
      const objectKeys = Object.keys(object);
      assert.false(objectKeys.includes('counts'), `${objectName} doesn't include 'counts' key`);
      assert.true(objectKeys.includes('label'), `${objectName} includes 'label' key`);
      assert.true(objectKeys.includes('clients'), `${objectName} includes 'clients' key`);
      assert.true(objectKeys.includes('entity_clients'), `${objectName} includes 'entity_clients' key`);
      assert.true(
        objectKeys.includes('non_entity_clients'),
        `${objectName} includes 'non_entity_clients' key`
      );
    };
    const keyValueAssertions = (object, pathName, originalObject) => {
      const keysToAssert = ['clients', 'entity_clients', 'non_entity_clients'];
      assert.strictEqual(object.label, originalObject[pathName], `${pathName} matches label`);

      keysToAssert.forEach((key) => {
        assert.strictEqual(object[key], originalObject.counts[key], `number of ${key} equal original`);
      });
    };

    const formattedNamespaces = formatByNamespace(BY_NAMESPACE);
    assert.notEqual(formattedNamespaces, MONTHS, 'does not modify original array');

    formattedNamespaces.forEach((namespace) => {
      const origNamespace = BY_NAMESPACE.find((ns) => ns.namespace_path === namespace.label);
      keyNameAssertions(namespace, 'formatted namespace');
      keyValueAssertions(namespace, 'namespace_path', origNamespace);

      namespace.mounts.forEach((mount) => {
        const origMount = origNamespace.mounts.find((m) => m.mount_path === mount.label);
        keyNameAssertions(mount, 'formatted mount');
        keyValueAssertions(mount, 'mount_path', origMount);
      });
    });

    const nsWithoutMounts = {
      namespace_id: '96OwG',
      namespace_path: 'no-mounts-ns/',
      counts: {
        distinct_entities: 18290,
        entity_clients: 18290,
        non_entity_tokens: 18738,
        non_entity_clients: 18738,
        clients: 37028,
      },
      mounts: [],
    };

    const formattedNsWithoutMounts = formatByNamespace([nsWithoutMounts])[0];
    keyNameAssertions(formattedNsWithoutMounts, 'namespace without mounts');
    keyValueAssertions(formattedNsWithoutMounts, 'namespace_path', nsWithoutMounts);
    assert.strictEqual(formattedNsWithoutMounts.mounts.length, 0, 'formatted namespace has no mounts');

    assert.strictEqual(formatByNamespace(SOME_OBJECT), SOME_OBJECT, 'it returns if arg is not an array');
  });

  test('destructureCounts: homogenizes key names when both old and new keys exist, or just old key names', async function (assert) {
    assert.expect(2);
    const original = {
      distinct_entities: 3,
      entity_clients: 3,
      non_entity_tokens: 5,
      non_entity_clients: 5,
      secret_syncs: 10,
      acme_clients: 4,
      clients: 22,
    };
    const expected = {
      entity_clients: 3,
      non_entity_clients: 5,
      secret_syncs: 10,
      acme_clients: 4,
      clients: 22,
    };
    assert.propEqual(destructureCounts(original), expected);
    assert.propEqual(
      original,
      {
        distinct_entities: 3,
        entity_clients: 3,
        non_entity_tokens: 5,
        non_entity_clients: 5,
        secret_syncs: 10,
        acme_clients: 4,
        clients: 22,
      },
      'original array is not modified'
    );
  });

  test('flattenDataset: removes the counts key and flattens the dataset', async function (assert) {
    assert.expect(22);
    const flattenedNamespace = flattenDataset(BY_NAMESPACE[0]);
    const flattenedMount = flattenDataset(BY_NAMESPACE[0].mounts[0]);
    const flattenedMonth = flattenDataset(MONTHS[0]);
    const flattenedNewMonthClients = flattenDataset(MONTHS[0].new_clients);
    const objectNullCounts = { counts: null, foo: 'bar' };

    const keyNameAssertions = (object, objectName) => {
      const objectKeys = Object.keys(object);
      assert.false(objectKeys.includes('counts'), `${objectName} doesn't include 'counts' key`);
      assert.true(objectKeys.includes('clients'), `${objectName} includes 'clients' key`);
      assert.true(objectKeys.includes('entity_clients'), `${objectName} includes 'entity_clients' key`);
      assert.true(
        objectKeys.includes('non_entity_clients'),
        `${objectName} includes 'non_entity_clients' key`
      );
    };

    keyNameAssertions(flattenedNamespace, 'namespace object');
    keyNameAssertions(flattenedMount, 'mount object');
    keyNameAssertions(flattenedMonth, 'month object');
    keyNameAssertions(flattenedNewMonthClients, 'month new_clients object');

    assert.strictEqual(
      flattenDataset(SOME_OBJECT),
      SOME_OBJECT,
      "it returns original object if counts key doesn't exist"
    );

    assert.strictEqual(
      flattenDataset(objectNullCounts),
      objectNullCounts,
      'it returns original object if counts are null'
    );

    assert.propEqual(
      flattenDataset(['some array']),
      ['some array'],
      'it fails gracefully if an array is passed in'
    );
    assert.strictEqual(flattenDataset(null), null, 'it fails gracefully if null is passed in');
    assert.strictEqual(
      flattenDataset('some string'),
      'some string',
      'it fails gracefully if a string is passed in'
    );
    assert.propEqual(
      flattenDataset(new Object()),
      new Object(),
      'it fails gracefully if an empty object is passed in'
    );
  });

  test('sortMonthsByTimestamp: sorts timestamps chronologically, oldest to most recent', async function (assert) {
    assert.expect(4);
    const sortedMonths = sortMonthsByTimestamp(MONTHS);
    assert.ok(
      isBefore(parseAPITimestamp(sortedMonths[0].timestamp), parseAPITimestamp(sortedMonths[1].timestamp)),
      'first timestamp date is earlier than second'
    );
    assert.ok(
      isAfter(parseAPITimestamp(sortedMonths[2].timestamp), parseAPITimestamp(sortedMonths[1].timestamp)),
      'third timestamp date is later second'
    );
    assert.notEqual(sortedMonths[1], MONTHS[1], 'it does not modify original array');
    assert.strictEqual(sortedMonths[0], MONTHS[0], 'it does not modify original array');
  });

  test('namespaceArrayToObject: transforms data without modifying original', async function (assert) {
    assert.expect(30);

    const assertClientCounts = (object, originalObject) => {
      const valuesToCheck = ['clients', 'entity_clients', 'non_entity_clients'];

      valuesToCheck.forEach((key) => {
        assert.strictEqual(object[key], originalObject[key], `${key} equal original counts`);
      });
    };
    const totalClientsByNamespace = formatByNamespace(MONTHS[1].namespaces);
    const newClientsByNamespace = formatByNamespace(MONTHS[1].new_clients.namespaces);

    const byNamespaceKeyObject = namespaceArrayToObject(
      totalClientsByNamespace,
      newClientsByNamespace,
      '10/21',
      '2021-10-01T00:00:00Z'
    );

    assert.propEqual(
      formatByNamespace(MONTHS[1].namespaces),
      totalClientsByNamespace,
      'it does not modify original array'
    );
    assert.propEqual(
      formatByNamespace(MONTHS[1].new_clients.namespaces),
      newClientsByNamespace,
      'it does not modify original array'
    );

    const namespaceKeys = Object.keys(byNamespaceKeyObject);
    namespaceKeys.forEach((nsKey) => {
      const newNsObject = byNamespaceKeyObject[nsKey];
      const originalNsData = totalClientsByNamespace.find((ns) => ns.label === nsKey);
      assertClientCounts(newNsObject, originalNsData);
      const mountKeys = Object.keys(newNsObject.mounts_by_key);
      mountKeys.forEach((mKey) => {
        const mountData = originalNsData.mounts.find((m) => m.label === mKey);
        assertClientCounts(newNsObject.mounts_by_key[mKey], mountData);
      });
    });

    namespaceKeys.forEach((nsKey) => {
      const newNsObject = byNamespaceKeyObject[nsKey];
      const originalNsData = newClientsByNamespace.find((ns) => ns.label === nsKey);
      if (!originalNsData) return;
      assertClientCounts(newNsObject.new_clients, originalNsData);
      const mountKeys = Object.keys(newNsObject.mounts_by_key);

      mountKeys.forEach((mKey) => {
        const mountData = originalNsData.mounts.find((m) => m.label === mKey);
        assertClientCounts(newNsObject.mounts_by_key[mKey].new_clients, mountData);
      });
    });

    assert.propEqual(
      namespaceArrayToObject(null, null, '10/21', 'timestamp-here'),
      {},
      'returns an empty object when totalClientsByNamespace = null'
    );
  });
});
