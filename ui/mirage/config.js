const EXPIRY_DATE = '2021-05-12T23:20:50.52Z';

export default function () {
  this.namespace = 'v1';

  this.get('sys/internal/counters/activity', function (db) {
    let data = {};
    const firstRecord = db['clients/activities'].first();
    if (firstRecord) {
      data = firstRecord;
    }
    return {
      data,
      request_id: '0001',
    };
  });

  this.get('sys/internal/counters/config', function (db) {
    return {
      request_id: '00001',
      data: db['clients/configs'].first(),
    };
  });

  this.get('/sys/internal/ui/feature-flags', () => {
    return {
      data: {
        start_time: '2019-11-01T00:00:00Z',
        end_time: '2020-10-31T23:59:59Z',
        total: {
          distinct_entities: 200,
          non_entity_tokens: 100,
          clients: 300,
        },
        by_namespace: [
          {
            _comment: 'by_namespace will remain as it is',
          },
        ],
        months: [
          {
            'jan/2022': {
              counts: {
                distinct_entities: 100,
                non_entity_tokens: 50,
                clients: 150,
              },
              namespaces: [
                {
                  id: 'root',
                  path: '',
                  counts: {
                    distinct_entities: 50,
                    non_entity_tokens: 25,
                    clients: 75,
                  },
                  mounts: [
                    {
                      path: 'auth/aws/login',
                      counts: {
                        distinct_entities: 25,
                        non_entity_tokens: 12,
                        clients: 37,
                      },
                    },
                    {
                      path: 'auth/approle/login',
                      counts: {
                        distinct_entities: 25,
                        non_entity_tokens: 13,
                        clients: 38,
                      },
                    },
                  ],
                },
                {
                  namespace_id: 'ns1',
                  namespace_path: '',
                  counts: {
                    distinct_entities: 50,
                    non_entity_tokens: 25,
                    clients: 75,
                  },
                  mounts: [
                    {
                      mount_path: 'auth/aws/login',
                      counts: {
                        distinct_entities: 20,
                        non_entity_tokens: 10,
                        clients: 30,
                      },
                    },
                    {
                      mount_path: 'auth/approle/login',
                      counts: {
                        distinct_entities: 30,
                        non_entity_tokens: 15,
                        clients: 45,
                      },
                    },
                  ],
                },
              ],
              new: {
                counts: {
                  distinct_entities: 30,
                  non_entity_tokens: 10,
                  clients: 40,
                },
                namespaces: [
                  {
                    namespace_id: 'root',
                    namespace_path: '',
                    counts: {
                      distinct_entities: 15,
                      non_entity_tokens: 5,
                      clients: 20,
                    },
                    mounts: [
                      {
                        mount_path: 'auth/aws/login',
                        counts: {
                          distinct_entities: 5,
                          non_entity_tokens: 2,
                          clients: 7,
                        },
                      },
                      {
                        mount_path: 'auth/aws/login',
                        counts: {
                          distinct_entities: 10,
                          non_entity_tokens: 3,
                          clients: 13,
                        },
                      },
                    ],
                  },
                  {
                    namespace_id: 'ns1',
                    namespace_path: '',
                    counts: {
                      distinct_entities: 15,
                      non_entity_tokens: 5,
                      clients: 20,
                    },
                    mounts: [
                      {
                        mount_path: 'auth/aws/login',
                        counts: {
                          distinct_entities: 10,
                          non_entity_tokens: 1,
                          clients: 11,
                        },
                      },
                      {
                        mount_path: 'auth/aws/login',
                        counts: {
                          distinct_entities: 5,
                          non_entity_tokens: 4,
                          clients: 9,
                        },
                      },
                    ],
                  },
                ],
              },
            },
          },
          {
            'feb/2022': {
              counts: {
                _comment: 'total monthly clients',
                distinct_entities: 100,
                non_entity_tokens: 50,
                clients: 150,
              },
              namespaces: [
                {
                  namespace_id: 'root',
                  namespace_path: '',
                  counts: {
                    distinct_entities: 60,
                    non_entity_tokens: 10,
                    clients: 70,
                  },
                  mounts: [
                    {
                      mount_path: 'auth/aws/login',
                      counts: {
                        distinct_entities: 30,
                        non_entity_tokens: 5,
                        clients: 35,
                      },
                    },
                    {
                      mount_path: 'auth/approle/login',
                      counts: {
                        distinct_entities: 30,
                        non_entity_tokens: 5,
                        clients: 35,
                      },
                    },
                  ],
                },
                {
                  namespace_id: 'ns1',
                  namespace_path: '',
                  counts: {
                    distinct_entities: 40,
                    non_entity_tokens: 40,
                    clients: 80,
                  },
                  mounts: [
                    {
                      mount_path: 'auth/aws/login',
                      counts: {
                        distinct_entities: 20,
                        non_entity_tokens: 20,
                        clients: 40,
                      },
                    },
                    {
                      mount_path: 'auth/approle/login',
                      counts: {
                        distinct_entities: 20,
                        non_entity_tokens: 20,
                        clients: 40,
                      },
                    },
                  ],
                },
              ],
              new: {
                counts: {
                  distinct_entities: 20,
                  non_entity_tokens: 5,
                  clients: 25,
                },
                namespaces: [
                  {
                    namespace_id: 'root',
                    namespace_path: '',
                    counts: {
                      distinct_entities: 10,
                      non_entity_tokens: 3,
                      clients: 13,
                    },
                    mounts: [
                      {
                        mount_path: 'auth/aws/login',
                        counts: {
                          distinct_entities: 5,
                          non_entity_tokens: 1,
                          clients: 6,
                        },
                      },
                      {
                        mount_path: 'auth/aws/login',
                        counts: {
                          distinct_entities: 5,
                          non_entity_tokens: 2,
                          clients: 7,
                        },
                      },
                    ],
                  },
                  {
                    namespace_id: 'ns1',
                    namespace_path: '',
                    counts: {
                      distinct_entities: 10,
                      non_entity_tokens: 2,
                      clients: 12,
                    },
                    mounts: [
                      {
                        mount_path: 'auth/aws/login',
                        counts: {
                          distinct_entities: 5,
                          non_entity_tokens: 2,
                          clients: 7,
                        },
                      },
                      {
                        mount_path: 'auth/aws/login',
                        counts: {
                          distinct_entities: 5,
                          non_entity_tokens: 0,
                          clients: 5,
                        },
                      },
                    ],
                  },
                ],
              },
            },
          },
        ],
      },
    };
  });

  this.get('/sys/internal/counters/activity/monthly', function () {
    return {
      data: {
        by_namespace: [
          {
            namespace_id: 'Z4Rzh',
            namespace_path: 'namespace1/',
            counts: {
              distinct_entities: 867,
              non_entity_tokens: 939,
              clients: 1806,
            },
          },
          {
            namespace_id: 'DcgzU',
            namespace_path: 'namespace17/',
            counts: {
              distinct_entities: 966,
              non_entity_tokens: 550,
              clients: 1516,
            },
          },
          {
            namespace_id: '5SWT8',
            namespace_path: 'namespacelonglonglong4/',
            counts: {
              distinct_entities: 996,
              non_entity_tokens: 417,
              clients: 1413,
            },
          },
          {
            namespace_id: 'XGu7R',
            namespace_path: 'namespace12/',
            counts: {
              distinct_entities: 829,
              non_entity_tokens: 540,
              clients: 1369,
            },
          },
          {
            namespace_id: 'yHcL9',
            namespace_path: 'namespace11/',
            counts: {
              distinct_entities: 563,
              non_entity_tokens: 705,
              clients: 1268,
            },
          },
          {
            namespace_id: 'F0xGm',
            namespace_path: 'namespace10/',
            counts: {
              distinct_entities: 925,
              non_entity_tokens: 255,
              clients: 1180,
            },
          },
          {
            namespace_id: 'aJuQG',
            namespace_path: 'namespace9/',
            counts: {
              distinct_entities: 935,
              non_entity_tokens: 239,
              clients: 1174,
            },
          },
          {
            namespace_id: 'bw5UO',
            namespace_path: 'namespace6/',
            counts: {
              distinct_entities: 810,
              non_entity_tokens: 363,
              clients: 1173,
            },
          },
          {
            namespace_id: 'IeyJp',
            namespace_path: 'namespace14/',
            counts: {
              distinct_entities: 774,
              non_entity_tokens: 392,
              clients: 1166,
            },
          },
          {
            namespace_id: 'Uc0o8',
            namespace_path: 'namespace16/',
            counts: {
              distinct_entities: 408,
              non_entity_tokens: 743,
              clients: 1151,
            },
          },
          {
            namespace_id: 'R6L40',
            namespace_path: 'namespace2/',
            counts: {
              distinct_entities: 292,
              non_entity_tokens: 736,
              clients: 1028,
            },
          },
          {
            namespace_id: 'Rqa3W',
            namespace_path: 'namespace13/',
            counts: {
              distinct_entities: 160,
              non_entity_tokens: 803,
              clients: 963,
            },
          },
          {
            namespace_id: 'MSgZE',
            namespace_path: 'namespace7/',
            counts: {
              distinct_entities: 201,
              non_entity_tokens: 657,
              clients: 858,
            },
          },
          {
            namespace_id: 'kxU4t',
            namespace_path: 'namespacelonglonglong3/',
            counts: {
              distinct_entities: 742,
              non_entity_tokens: 26,
              clients: 768,
            },
          },
          {
            namespace_id: '5xKya',
            namespace_path: 'namespace15/',
            counts: {
              distinct_entities: 663,
              non_entity_tokens: 19,
              clients: 682,
            },
          },
          {
            namespace_id: '5KxXA',
            namespace_path: 'namespace18anotherlong/',
            counts: {
              distinct_entities: 470,
              non_entity_tokens: 196,
              clients: 666,
            },
          },
          {
            namespace_id: 'AAidI',
            namespace_path: 'namespace20/',
            counts: {
              distinct_entities: 429,
              non_entity_tokens: 60,
              clients: 489,
            },
          },
          {
            namespace_id: 'BCl56',
            namespace_path: 'namespace8/',
            counts: {
              distinct_entities: 61,
              non_entity_tokens: 201,
              clients: 262,
            },
          },
          {
            namespace_id: 'yYNw2',
            namespace_path: 'namespace19/',
            counts: {
              distinct_entities: 165,
              non_entity_tokens: 85,
              clients: 250,
            },
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 67,
              non_entity_tokens: 9,
              clients: 76,
            },
          },
        ],
        distinct_entities: 11323,
        non_entity_tokens: 7935,
        clients: 19258,
      },
    };
  });

  this.get('/sys/health', function () {
    return {
      initialized: true,
      sealed: false,
      standby: false,
      license: {
        expiry: '2021-05-12T23:20:50.52Z',
        state: 'stored',
      },
      performance_standby: false,
      replication_performance_mode: 'disabled',
      replication_dr_mode: 'disabled',
      server_time_utc: 1622562585,
      version: '1.9.0+ent',
      cluster_name: 'vault-cluster-e779cd7c',
      cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
      last_wal: 121,
    };
  });

  this.get('/sys/license/status', function () {
    return {
      data: {
        autoloading_used: false,
        stored: {
          expiration_time: EXPIRY_DATE,
          features: ['DR Replication', 'Namespaces', 'Lease Count Quotas', 'Automated Snapshots'],
          license_id: '0eca7ef8-ebc0-f875-315e-3cc94a7870cf',
          performance_standby_count: 0,
          start_time: '2020-04-28T00:00:00Z',
        },
        persisted_autoload: {
          expiration_time: EXPIRY_DATE,
          features: ['DR Replication', 'Namespaces', 'Lease Count Quotas', 'Automated Snapshots'],
          license_id: '0eca7ef8-ebc0-f875-315e-3cc94a7870cf',
          performance_standby_count: 0,
          start_time: '2020-04-28T00:00:00Z',
        },
        autoloaded: {
          expiration_time: EXPIRY_DATE,
          features: ['DR Replication', 'Namespaces', 'Lease Count Quotas', 'Automated Snapshots'],
          license_id: '0eca7ef8-ebc0-f875-315e-3cc94a7870cf',
          performance_standby_count: 0,
          start_time: '2020-04-28T00:00:00Z',
        },
      },
    };
  });

  this.get('sys/namespaces', function () {
    return {
      data: {
        keys: [
          'ns1/',
          'ns2/',
          'ns3/',
          'ns4/',
          'ns5/',
          'ns6/',
          'ns7/',
          'ns8/',
          'ns9/',
          'ns10/',
          'ns11/',
          'ns12/',
          'ns13/',
          'ns14/',
          'ns15/',
          'ns16/',
          'ns17/',
          'ns18/',
        ],
      },
    };
  });

  this.passthrough();
}
