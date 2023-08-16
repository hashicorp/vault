import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import {
  flattenDataset,
  formatByMonths,
  formatByNamespace,
  homogenizeClientNaming,
  sortMonthsByTimestamp,
  namespaceArrayToObject,
} from 'core/utils/client-count-utils';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import isBefore from 'date-fns/isBefore';
import isAfter from 'date-fns/isAfter';

// import { setupMirage } from 'ember-cli-mirage/test-support';
// import ENV from 'vault/config/environment';
// import { formatRFC3339 } from 'date-fns';

module('Integration | Util | client count utils', function (hooks) {
  setupTest(hooks);
  // setupMirage(hooks);

  // TODO: wire up to stubbed API/mirage?
  // hooks.before(function () {
  //   ENV['ember-cli-mirage'].handler = 'clients';
  // });
  // hooks.after(function () {
  //   ENV['ember-cli-mirage'].handler = null;
  // });

  /* MONTHS array contains: (update when backend work done on months )
  - one month with only old client naming
  */

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
      let newObjectKeys = ['clients', 'entity_clients', 'non_entity_clients'];
      let originalKeys = Object.keys(originalObject.counts).includes('entity_clients')
        ? newObjectKeys
        : ['clients', 'distinct_entities', 'non_entity_tokens'];

      newObjectKeys.forEach((key, i) => {
        assert.equal(
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
    let expected = [
      {
        counts: null,
        month: '6/21',
        namespaces: [],
        namespaces_by_key: {},
        new_clients: {
          month: '6/21',
          namespaces: [],
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
        },
        timestamp: '2021-07-01T00:00:00Z',
      },
    ];
    assert.equal(formatByMonths(SOME_OBJECT), SOME_OBJECT, 'it returns if arg is not an array');
    assert.propEqual(expected, formatByMonths(EMPTY_MONTHS), 'it does not error with null months');
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
      assert.equal(object.label, originalObject[pathName], `${pathName} matches label`);

      keysToAssert.forEach((key) => {
        assert.equal(object[key], originalObject.counts[key], `number of ${key} equal original`);
      });
    };

    const formattedNamespaces = formatByNamespace(BY_NAMESPACE);
    assert.notEqual(formattedNamespaces, MONTHS, 'does not modify original array');

    formattedNamespaces.forEach((namespace) => {
      let origNamespace = BY_NAMESPACE.find((ns) => ns.namespace_path === namespace.label);
      keyNameAssertions(namespace, 'formatted namespace');
      keyValueAssertions(namespace, 'namespace_path', origNamespace);

      namespace.mounts.forEach((mount) => {
        let origMount = origNamespace.mounts.find((m) => m.mount_path === mount.label);
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

    let formattedNsWithoutMounts = formatByNamespace([nsWithoutMounts])[0];
    keyNameAssertions(formattedNsWithoutMounts, 'namespace without mounts');
    keyValueAssertions(formattedNsWithoutMounts, 'namespace_path', nsWithoutMounts);
    assert.equal(formattedNsWithoutMounts.mounts.length, 0, 'formatted namespace has no mounts');

    assert.equal(formatByNamespace(SOME_OBJECT), SOME_OBJECT, 'it returns if arg is not an array');
  });

  test('homogenizeClientNaming: homogenizes key names when both old and new keys exist, or just old key names', async function (assert) {
    assert.expect(168);
    const keyNameAssertions = (object, objectName) => {
      const objectKeys = Object.keys(object);
      assert.false(
        objectKeys.includes('distinct_entities'),
        `${objectName} doesn't include 'distinct_entities' key`
      );
      assert.false(
        objectKeys.includes('non_entity_tokens'),
        `${objectName} doesn't include 'non_entity_tokens' key`
      );
      assert.true(objectKeys.includes('entity_clients'), `${objectName} includes 'entity_clients' key`);
      assert.true(
        objectKeys.includes('non_entity_clients'),
        `${objectName} includes 'non_entity_clients' key`
      );
    };

    let transformedMonths = [...MONTHS];
    transformedMonths.forEach((month) => {
      month.counts = homogenizeClientNaming(month.counts);
      keyNameAssertions(month.counts, 'month counts');

      month.new_clients.counts = homogenizeClientNaming(month.new_clients.counts);
      keyNameAssertions(month.new_clients.counts, 'month new counts');

      month.namespaces.forEach((ns) => {
        ns.counts = homogenizeClientNaming(ns.counts);
        keyNameAssertions(ns.counts, 'namespace counts');

        ns.mounts.forEach((mount) => {
          mount.counts = homogenizeClientNaming(mount.counts);
          keyNameAssertions(mount.counts, 'mount counts');
        });
      });

      month.new_clients.namespaces.forEach((ns) => {
        ns.counts = homogenizeClientNaming(ns.counts);
        keyNameAssertions(ns.counts, 'namespace new counts');

        ns.mounts.forEach((mount) => {
          mount.counts = homogenizeClientNaming(mount.counts);
          keyNameAssertions(mount.counts, 'mount new counts');
        });
      });
    });
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

    assert.equal(
      flattenDataset(SOME_OBJECT),
      SOME_OBJECT,
      "it returns original object if counts key doesn't exist"
    );

    assert.equal(
      flattenDataset(objectNullCounts),
      objectNullCounts,
      'it returns original object if counts are null'
    );

    assert.propEqual(
      ['some array'],
      flattenDataset(['some array']),
      'it fails gracefully if an array is passed in'
    );
    assert.equal(flattenDataset(null), null, 'it fails gracefully if null is passed in');
    assert.equal(
      flattenDataset('some string'),
      'some string',
      'it fails gracefully if a string is passed in'
    );
    assert.propEqual(
      new Object(),
      flattenDataset(new Object()),
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
    assert.equal(sortedMonths[0], MONTHS[0], 'it does not modify original array');
  });

  test('namespaceArrayToObject: transforms data without modifying original', async function (assert) {
    assert.expect(30);

    const assertClientCounts = (object, originalObject) => {
      let valuesToCheck = ['clients', 'entity_clients', 'non_entity_clients'];

      valuesToCheck.forEach((key) => {
        assert.equal(object[key], originalObject[key], `${key} equal original counts`);
      });
    };
    const totalClientsByNamespace = formatByNamespace(MONTHS[1].namespaces);
    const newClientsByNamespace = formatByNamespace(MONTHS[1].new_clients.namespaces);

    const byNamespaceKeyObject = namespaceArrayToObject(
      totalClientsByNamespace,
      newClientsByNamespace,
      '10/21'
    );

    assert.propEqual(
      totalClientsByNamespace,
      formatByNamespace(MONTHS[1].namespaces),
      'it does not modify original array'
    );
    assert.propEqual(
      newClientsByNamespace,
      formatByNamespace(MONTHS[1].new_clients.namespaces),
      'it does not modify original array'
    );

    let namespaceKeys = Object.keys(byNamespaceKeyObject);
    namespaceKeys.forEach((nsKey) => {
      const newNsObject = byNamespaceKeyObject[nsKey];
      let originalNsData = totalClientsByNamespace.find((ns) => ns.label === nsKey);
      assertClientCounts(newNsObject, originalNsData);
      let mountKeys = Object.keys(newNsObject.mounts_by_key);
      mountKeys.forEach((mKey) => {
        let mountData = originalNsData.mounts.find((m) => m.label === mKey);
        assertClientCounts(newNsObject.mounts_by_key[mKey], mountData);
      });
    });

    namespaceKeys.forEach((nsKey) => {
      const newNsObject = byNamespaceKeyObject[nsKey];
      let originalNsData = newClientsByNamespace.find((ns) => ns.label === nsKey);
      if (!originalNsData) return;
      assertClientCounts(newNsObject.new_clients, originalNsData);
      let mountKeys = Object.keys(newNsObject.mounts_by_key);

      mountKeys.forEach((mKey) => {
        let mountData = originalNsData.mounts.find((m) => m.label === mKey);
        assertClientCounts(newNsObject.mounts_by_key[mKey].new_clients, mountData);
      });
    });

    assert.propEqual(
      {},
      namespaceArrayToObject(null, null, '10/21'),
      'returns an empty object when totalClientsByNamespace = null'
    );
  });
});
