/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { filterVersionHistory, flattenMounts, filterTableData } from 'core/utils/client-counts/helpers';
import clientsHandler from 'vault/mirage/handlers/clients';
import {
  SERIALIZED_ACTIVITY_RESPONSE,
  ENTITY_EXPORT,
} from 'vault/tests/helpers/clients/client-count-helpers';

module('Unit | Util | client counts | helpers', function (hooks) {
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

  test('flattenMounts: it flattens mount data', async function (assert) {
    assert.expect(2);
    const original = [...SERIALIZED_ACTIVITY_RESPONSE.by_namespace];
    const expected = [
      ...SERIALIZED_ACTIVITY_RESPONSE.by_namespace[0].mounts,
      ...SERIALIZED_ACTIVITY_RESPONSE.by_namespace[1].mounts,
    ];
    const actual = flattenMounts(SERIALIZED_ACTIVITY_RESPONSE.by_namespace);

    assert.propEqual(actual, expected, 'it returns mounts from each namespace object into a single array');
    assert.propEqual(
      SERIALIZED_ACTIVITY_RESPONSE.by_namespace,
      original,
      'it does not modify original by_namespace array'
    );
  });

  module('filterTableData', function (hooks) {
    hooks.beforeEach(async function () {
      const activityByMount = flattenMounts(SERIALIZED_ACTIVITY_RESPONSE.by_namespace);
      this.mockMountData = [...activityByMount];
      // copy mock data before using the filterTableData function to assert filtering doesn't modify the original array
      const original = [...this.mockMountData];
      this.assertOriginal = (assert) =>
        assert.propEqual(this.mockMountData, original, 'filtering does not mutate dataset');
    });

    test('it returns original data if no filters are passed', async function (assert) {
      const emptyObject = filterTableData(this.mockMountData, {});
      assert.propEqual(emptyObject, this.mockMountData, 'when filters arg is an empty object');
      this.assertOriginal(assert);

      const emptyValues = filterTableData(this.mockMountData, {
        namespace_path: '',
        mount_path: '',
        mount_type: '',
      });
      assert.propEqual(emptyValues, this.mockMountData, 'when filters have are empty strings');
      this.assertOriginal(assert);

      const nullFilter = filterTableData(this.mockMountData, null);
      assert.propEqual(nullFilter, this.mockMountData, 'returns all data when no filters are null');
      this.assertOriginal(assert);
    });

    test('it filters data for a single filter', async function (assert) {
      const namespaceFilter = filterTableData(this.mockMountData, {
        namespace_path: 'root',
        mount_path: '',
        mount_type: '',
      });
      const expectedNamespaceFilter = this.mockMountData.filter((m) => m.namespace_path === 'root');
      assert.propEqual(namespaceFilter, expectedNamespaceFilter, 'it filters by namespace_path');
      this.assertOriginal(assert);

      const mountPathFilter = filterTableData(this.mockMountData, {
        namespace_path: '',
        mount_path: 'acme/pki/0/',
        mount_type: '',
      });
      const expectedMountPathFilter = this.mockMountData.filter((m) => m.mount_path === 'acme/pki/0/');
      assert.propEqual(mountPathFilter, expectedMountPathFilter, 'it filters by mount_path');
      this.assertOriginal(assert);

      const mountTypeFilter = filterTableData(this.mockMountData, {
        namespace_path: '',
        mount_path: '',
        mount_type: 'userpass',
      });
      const expectedMountTypeFilter = this.mockMountData.filter((m) => m.mount_type === 'userpass');
      assert.propEqual(mountTypeFilter, expectedMountTypeFilter, 'it filters by mount_type');
      this.assertOriginal(assert);
    });

    test('it filters data for a multiple filters', async function (assert) {
      const twoFilters = filterTableData(this.mockMountData, {
        namespace_path: 'root',
        mount_path: '',
        mount_type: 'userpass',
      });
      const expectedTwoFilters = this.mockMountData.filter(
        (m) => m.namespace_path === 'root' && m.mount_type === 'userpass'
      );
      assert.propEqual(twoFilters, expectedTwoFilters, 'it filters by namespace_path and mount_type');
      this.assertOriginal(assert);

      const allFilters = filterTableData(this.mockMountData, {
        namespace_path: 'root',
        mount_path: 'auth/userpass/0/',
        mount_type: 'userpass',
      });
      const expectedAllFilters = [
        {
          label: 'auth/userpass/0/',
          mount_path: 'auth/userpass/0/',
          mount_type: 'userpass',
          namespace_path: 'root',
          acme_clients: 0,
          clients: 8091,
          entity_clients: 4002,
          non_entity_clients: 4089,
          secret_syncs: 0,
        },
      ];
      assert.propEqual(allFilters, expectedAllFilters, 'it filters by all filters');
      this.assertOriginal(assert);
    });

    test('it returns an empty array when there are no matches', async function (assert) {
      const noMatches = filterTableData(this.mockMountData, {
        namespace_path: 'does not exist',
        mount_path: '',
        mount_type: '',
      });
      assert.propEqual(noMatches, [], 'returns an empty array when no data matches filters');
      this.assertOriginal(assert);
    });

    test('it returns an empty array when filter includes keys the dataset does not contain', async function (assert) {
      const noMatches = filterTableData(this.mockMountData, { foo: 'root', bar: '' });
      assert.propEqual(noMatches, [], 'returns an empty array when no keys match dataset');
      this.assertOriginal(assert);
    });

    test('it matches on empty strings or "root" for the root namespace', async function (assert) {
      const mockExportData = ENTITY_EXPORT.trim()
        .split('\n')
        .map((line) => JSON.parse(line));
      const combinedData = [...this.mockMountData, ...mockExportData];
      const filteredData = filterTableData(combinedData, { namespace_path: 'root' });
      const expected = [
        {
          acme_clients: 0,
          clients: 8091,
          entity_clients: 4002,
          label: 'auth/userpass/0/',
          mount_path: 'auth/userpass/0/',
          mount_type: 'userpass',
          namespace_path: 'root',
          non_entity_clients: 4089,
          secret_syncs: 0,
        },
        {
          acme_clients: 0,
          clients: 4290,
          entity_clients: 0,
          label: 'secrets/kv/0/',
          mount_path: 'secrets/kv/0/',
          mount_type: 'kv',
          namespace_path: 'root',
          non_entity_clients: 0,
          secret_syncs: 4290,
        },
        {
          acme_clients: 4003,
          clients: 4003,
          entity_clients: 0,
          label: 'acme/pki/0/',
          mount_path: 'acme/pki/0/',
          mount_type: 'pki',
          namespace_path: 'root',
          non_entity_clients: 0,
          secret_syncs: 0,
        },
        {
          client_first_used_time: '2025-07-15T23:48:09Z',
          client_id: 'daf8420c-0b6b-34e6-ff38-ee1ed093bea9',
          client_type: 'entity',
          entity_alias_custom_metadata: {},
          entity_alias_metadata: {},
          entity_alias_name: 'bob',
          entity_group_ids: ['7537e6b7-3b06-65c2-1fb2-c83116eb5e6f'],
          entity_metadata: {},
          entity_name: 'entity_b3e2a7ff',
          local_entity_alias: false,
          mount_accessor: 'auth_userpass_f47ad0b4',
          mount_path: 'auth/userpass/',
          mount_type: 'userpass',
          namespace_id: 'root',
          namespace_path: '',
          policies: [],
          token_creation_time: '2020-08-15T23:48:09Z',
        },
        {
          client_first_used_time: '2025-08-15T23:53:19Z',
          client_id: '23a04911-5d72-ba98-11d3-527f2fcf3a81',
          client_type: 'entity',
          entity_alias_custom_metadata: {
            account: 'Tester Account',
          },
          entity_alias_metadata: {},
          entity_alias_name: 'bob',
          entity_group_ids: ['7537e6b7-3b06-65c2-1fb2-c83116eb5e6f'],
          entity_metadata: {
            organization: 'ACME Inc.',
            team: 'QA',
          },
          entity_name: 'bob-smith',
          local_entity_alias: false,
          mount_accessor: 'auth_userpass_de28062c',
          mount_path: 'auth/userpass-test/',
          mount_type: 'userpass',
          namespace_id: 'root',
          namespace_path: '',
          policies: ['base'],
          token_creation_time: '2020-08-15T23:52:38Z',
        },
        {
          client_first_used_time: '2025-09-16T09:16:03Z',
          client_id: 'a7c8d912-4f61-23b5-88e4-627a3dcf2b92',
          client_type: 'entity',
          entity_alias_custom_metadata: {
            role: 'Senior Engineer',
          },
          entity_alias_metadata: {
            department: 'Engineering',
          },
          entity_alias_name: 'alice',
          entity_group_ids: ['7537e6b7-3b06-65c2-1fb2-c83116eb5e6f', 'a1b2c3d4-5e6f-7g8h-9i0j-k1l2m3n4o5p6'],
          entity_metadata: {
            location: 'San Francisco',
            organization: 'TechCorp',
            team: 'DevOps',
          },
          entity_name: 'alice-johnson',
          local_entity_alias: false,
          mount_accessor: 'auth_userpass_f47ad0b4',
          mount_path: 'auth/userpass/',
          mount_type: 'userpass',
          namespace_id: 'root',
          namespace_path: '',
          policies: ['admin', 'audit'],
          token_creation_time: '2020-08-16T09:15:42Z',
        },
        {
          client_first_used_time: '2025-12-17T16:44:12Z',
          client_id: 'c6b9d248-5a71-39e4-c7f2-951d8eaf6b95',
          client_type: 'entity',
          entity_alias_custom_metadata: {
            expertise: 'kubernetes',
            on_call: 'true',
          },
          entity_alias_metadata: {
            iss: 'https://auth.cloudops.io',
            sub: 'frank.castle@cloudops.io',
          },
          entity_alias_name: 'frank',
          entity_group_ids: ['9a8b7c6d-5e4f-3210-9876-543210fedcba'],
          entity_metadata: {
            organization: 'CloudOps',
            region: 'us-east-1',
            team: 'SRE',
          },
          entity_name: 'frank-castle',
          local_entity_alias: false,
          mount_accessor: 'auth_jwt_9d8c7b6a',
          mount_path: 'auth/jwt/',
          mount_type: 'jwt',
          namespace_id: 'root',
          namespace_path: '',
          policies: ['operations', 'monitoring'],
          token_creation_time: '2020-08-17T16:43:28Z',
        },
      ];

      filteredData.forEach((d, idx) => {
        const identifier = idx < 3 ? `label: ${d.label}` : `client_id: ${d.client_id}`;
        assert.propEqual(d, expected[idx], `filtered data contains ${identifier}`);
      });
      this.assertOriginal(assert);
    });
  });
});
