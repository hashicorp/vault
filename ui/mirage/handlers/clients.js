export default function (server) {
  // 1.10 API response
  server.get('sys/version-history', function () {
    return {
      keys: ['1.9.0', '1.9.1', '1.9.2'],
      key_info: {
        '1.9.0': {
          previous_version: null,
          timestamp_installed: '2021-11-03T10:23:16Z',
        },
        '1.9.1': {
          previous_version: '1.9.0',
          timestamp_installed: '2021-12-03T10:23:16Z',
        },
        '1.9.2': {
          previous_version: '1.9.1',
          timestamp_installed: '2021-01-03T10:23:16Z',
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

  server.get(
    '/sys/internal/counters/activity',
    function () {
      return {
        request_id: '26be5ab9-dcac-9237-ec12-269a8ca647d5',
        data: {
          by_namespace: [
            {
              namespace_id: 'root',
              namespace_path: '',
              counts: {
                distinct_entities: 10,
                entity_clients: 10,
                non_entity_tokens: 10,
                non_entity_clients: 10,
                clients: 20,
              },
              mounts: [
                {
                  mount_path: 'auth/up1/',
                  counts: {
                    distinct_entities: 0,
                    entity_clients: 0,
                    non_entity_tokens: 0,
                    non_entity_clients: 10,
                    clients: 10,
                  },
                },
                {
                  mount_path: 'auth/up2/',
                  counts: {
                    distinct_entities: 0,
                    entity_clients: 10,
                    non_entity_tokens: 0,
                    non_entity_clients: 0,
                    clients: 10,
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
              namespace_id: 'oImjk',
              namespace_path: 'ns2/',
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
          ],
          end_time: '2022-05-31T23:59:59Z',
          months: [
            {
              timestamp: '2021-05-01T00:00:00Z',
              counts: {
                distinct_entities: 0,
                entity_clients: 13,
                non_entity_tokens: 0,
                non_entity_clients: 12,
                clients: 25,
              },
              namespaces: [
                {
                  namespace_id: 'root',
                  namespace_path: '',
                  counts: {
                    distinct_entities: 0,
                    entity_clients: 8,
                    non_entity_tokens: 0,
                    non_entity_clients: 7,
                    clients: 15,
                  },
                  mounts: [
                    {
                      mount_path: 'auth/up2/',
                      counts: {
                        distinct_entities: 0,
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
                        non_entity_tokens: 0,
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
              ],
              new_clients: {
                counts: {
                  distinct_entities: 0,
                  entity_clients: 3,
                  non_entity_tokens: 0,
                  non_entity_clients: 2,
                  clients: 5,
                },
                namespaces: [
                  {
                    namespace_id: 'root',
                    namespace_path: '',
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
            {
              timestamp: '2021-04-01T00:00:00Z',
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
            {
              timestamp: '2021-03-01T00:00:00Z',
              counts: {
                distinct_entities: 0,
                entity_clients: 7,
                non_entity_tokens: 0,
                non_entity_clients: 8,
                clients: 15,
              },
              namespaces: [
                {
                  namespace_id: 'root',
                  namespace_path: '',
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
                  namespace_id: 's07UR',
                  namespace_path: 'ns1/',
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
              ],
              new_clients: {
                counts: {
                  distinct_entities: 0,
                  entity_clients: 2,
                  non_entity_tokens: 0,
                  non_entity_clients: 3,
                  clients: 5,
                },
                namespaces: [
                  {
                    namespace_id: 's07UR',
                    namespace_path: 'ns1/',
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
                ],
              },
            },
            {
              timestamp: '2021-02-01T00:00:00Z',
              counts: {
                distinct_entities: 0,
                entity_clients: 5,
                non_entity_tokens: 0,
                non_entity_clients: 5,
                clients: 10,
              },
              namespaces: [
                {
                  namespace_id: 'root',
                  namespace_path: '',
                  counts: {
                    distinct_entities: 0,
                    entity_clients: 5,
                    non_entity_tokens: 0,
                    non_entity_clients: 5,
                    clients: 10,
                  },
                  mounts: [
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
                  ],
                },
              ],
              new_clients: {
                counts: {
                  distinct_entities: 0,
                  entity_clients: 3,
                  non_entity_tokens: 0,
                  non_entity_clients: 2,
                  clients: 5,
                },
                namespaces: [
                  {
                    namespace_id: 'root',
                    namespace_path: '',
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
            {
              timestamp: '2021-01-01T00:00:00Z',
              counts: {
                distinct_entities: 0,
                entity_clients: 2,
                non_entity_tokens: 0,
                non_entity_clients: 3,
                clients: 5,
              },
              namespaces: [
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
              ],
              new_clients: {
                counts: {
                  distinct_entities: 0,
                  entity_clients: 2,
                  non_entity_tokens: 0,
                  non_entity_clients: 3,
                  clients: 5,
                },
                namespaces: [
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
                ],
              },
            },
          ],
          start_time: '2021-05-01T00:00:00Z',
          total: {
            distinct_entities: 20,
            entity_clients: 20,
            non_entity_tokens: 20,
            non_entity_clients: 20,
            clients: 40,
          },
        },
      };
    },
    { timing: 3000 }
  );

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
                path: 'auth/method/uMGBU',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/woiej',
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
                path: 'auth/method/ABCD1',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/ABCD2',
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
                path: 'auth/method/XYZZ2',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/XYZZ1',
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
