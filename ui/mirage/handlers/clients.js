import { formatISO, isBefore, sub } from 'date-fns';

export default function (server) {
  // 1.10 API response
  server.get('sys/version-history', function () {
    return {
      data: {
        keys: ['1.9.0', '1.9.1', '1.9.2', '1.10.1'],
        key_info: {
          '1.9.0': {
            previous_version: null,
            timestamp_installed: '2021-07-03T10:23:16Z',
          },
          '1.9.1': {
            previous_version: '1.9.0',
            timestamp_installed: '2021-08-03T10:23:16Z',
          },
          '1.9.2': {
            previous_version: '1.9.1',
            timestamp_installed: '2021-09-03T10:23:16Z',
          },
          '1.10.1': {
            previous_version: '1.9.2',
            timestamp_installed: '2021-10-03T10:23:16Z',
          },
        },
      },
    };
  });

  server.get('sys/license/status', function () {
    return {
      data: {
        autoloading_used: true,
        autoloaded: {
          expiration_time: '2022-05-17T23:59:59.999Z',
          features: [
            'HSM',
            'Performance Replication',
            'DR Replication',
            'MFA',
            'Sentinel',
            'Seal Wrapping',
            'Control Groups',
            'Performance Standby',
            'Namespaces',
            'KMIP',
            'Entropy Augmentation',
            'Transform Secrets Engine',
            'Lease Count Quotas',
            'Key Management Secrets Engine',
            'Automated Snapshots',
          ],
          license_id: '060d7820-fa59-f95c-832b-395db0aeb9ba',
          performance_standby_count: 9999,
          start_time: '2021-05-17T00:00:00Z',
        },
        persisted_autoload: {
          expiration_time: '2022-05-17T23:59:59.999Z',
          features: [
            'HSM',
            'Performance Replication',
            'DR Replication',
            'MFA',
            'Sentinel',
            'Seal Wrapping',
            'Control Groups',
            'Performance Standby',
            'Namespaces',
            'KMIP',
            'Entropy Augmentation',
            'Transform Secrets Engine',
            'Lease Count Quotas',
            'Key Management Secrets Engine',
            'Automated Snapshots',
          ],
          license_id: '060d7820-fa59-f95c-832b-395db0aeb9ba',
          performance_standby_count: 9999,
          start_time: '2021-05-17T00:00:00Z',
        },
      },
    };
  });

  server.get('sys/internal/counters/config', function () {
    return {
      request_id: '00001',
      data: {
        default_report_months: 12,
        enabled: 'default-enable',
        queries_available: true,
        retention_months: 24,
      },
    };
  });

  server.get('/sys/internal/counters/activity', (schema, req) => {
    const { start_time, end_time } = req.queryParams;
    // fake client counting start date so warning shows if user queries earlier start date
    const counts_start = '2020-10-17T00:00:00Z';
    return {
      request_id: '25f55fbb-f253-9c46-c6f0-3cdd3ada91ab',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        by_namespace: [
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
            ],
          },
        ],
        end_time: end_time || formatISO(sub(new Date(), { months: 1 })),
        months: [],
        start_time: isBefore(new Date(start_time), new Date(counts_start)) ? counts_start : start_time,
        total: {
          distinct_entities: 37389,
          entity_clients: 37389,
          non_entity_tokens: 36519,
          non_entity_clients: 36519,
          clients: 73908,
        },
      },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });

  server.get('/sys/internal/counters/activity/monthly', function () {
    return {
      request_id: '26be5ab9-dcac-9237-ec12-269a8ca64742',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        by_namespace: [
          {
            namespace_id: '0lHBL',
            namespace_path: 'ns1/',
            counts: {
              distinct_entities: 85,
              non_entity_tokens: 15,
              clients: 100,
            },
            mounts: [
              {
                mount_path: 'auth/method/uMGBU',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                mount_path: 'auth/method/woiej',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
            ],
          },
          {
            namespace_id: 'RxD81',
            namespace_path: 'ns2/',
            counts: {
              distinct_entities: 35,
              non_entity_tokens: 20,
              clients: 55,
            },
            mounts: [
              {
                mount_path: 'auth/method/ABCD1',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                mount_path: 'auth/method/ABCD2',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
            ],
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 12,
              non_entity_tokens: 8,
              clients: 20,
            },
            mounts: [
              {
                mount_path: 'auth/method/XYZZ2',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                mount_path: 'auth/method/XYZZ1',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
            ],
          },
        ],
        distinct_entities: 132,
        non_entity_tokens: 43,
        clients: 175,
      },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });
}
