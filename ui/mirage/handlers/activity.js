export default function (server) {
  // 1.10 API response
  server.get('/sys/internal/counters/activity', function () {
    return {
      data: {
        start_time: '2019-11-01T00:00:00Z',
        end_time: '2020-10-31T23:59:59Z',
        total: {
          _comment1: 'total client counts',
          clients: 300,
          _comment2: 'following 2 fields are deprecated',
          distinct_entities: 200,
          non_entity_tokens: 100,
        },
        by_namespace: [
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 110,
              non_entity_tokens: 35,
              clients: 145,
            },
            mounts: [
              {
                path: 'auth/aws/login',
                counts: {
                  clients: 44,
                  entity_clients: 30,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/approle/login',
                counts: {
                  clients: 51,
                  entity_clients: 35,
                  non_entity_clients: 16,
                },
              },
            ],
          },
          {
            namespace_id: 'DochC',
            namespace_path: 'ns1',
            counts: {
              distinct_entities: 90,
              non_entity_tokens: 65,
              clients: 155,
            },
            mounts: [
              {
                path: 'auth/aws/login',
                counts: {
                  clients: 65,
                  entity_clients: 50,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/approle/login',
                counts: {
                  clients: 80,
                  entity_clients: 60,
                  non_entity_clients: 20,
                },
              },
            ],
          },
          {
            namespace_id: 'DochC',
            namespace_path: 'ns1',
            counts: {
              distinct_entities: 90,
              non_entity_tokens: 65,
              clients: 155,
            },
            mounts: [
              {
                path: 'auth/aws/login',
                counts: {
                  clients: 65,
                  entity_clients: 50,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/approle/login',
                counts: {
                  clients: 80,
                  entity_clients: 60,
                  non_entity_clients: 20,
                },
              },
            ],
          },
          {
            namespace_id: 'DochC',
            namespace_path: 'ns1',
            counts: {
              distinct_entities: 90,
              non_entity_tokens: 65,
              clients: 155,
            },
            mounts: [
              {
                path: 'auth/aws/login',
                counts: {
                  clients: 65,
                  entity_clients: 50,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/approle/login',
                counts: {
                  clients: 80,
                  entity_clients: 60,
                  non_entity_clients: 20,
                },
              },
            ],
          },
          {
            namespace_id: 'DochC',
            namespace_path: 'ns1',
            counts: {
              distinct_entities: 90,
              non_entity_tokens: 65,
              clients: 155,
            },
            mounts: [
              {
                path: 'auth/aws/login',
                counts: {
                  clients: 65,
                  entity_clients: 50,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/approle/login',
                counts: {
                  clients: 80,
                  entity_clients: 60,
                  non_entity_clients: 20,
                },
              },
            ],
          },
          {
            namespace_id: 'DochC',
            namespace_path: 'ns1',
            counts: {
              distinct_entities: 90,
              non_entity_tokens: 65,
              clients: 155,
            },
            mounts: [
              {
                path: 'auth/aws/login',
                counts: {
                  clients: 65,
                  entity_clients: 50,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/approle/login',
                counts: {
                  clients: 80,
                  entity_clients: 60,
                  non_entity_clients: 20,
                },
              },
            ],
          },
          {
            namespace_id: 'DochC',
            namespace_path: 'ns1',
            counts: {
              distinct_entities: 90,
              non_entity_tokens: 65,
              clients: 155,
            },
            mounts: [
              {
                path: 'auth/aws/login',
                counts: {
                  clients: 65,
                  entity_clients: 50,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/approle/login',
                counts: {
                  clients: 80,
                  entity_clients: 60,
                  non_entity_clients: 20,
                },
              },
            ],
          },
          {
            namespace_id: 'DochC',
            namespace_path: 'ns1',
            counts: {
              distinct_entities: 90,
              non_entity_tokens: 65,
              clients: 155,
            },
            mounts: [
              {
                path: 'auth/aws/login',
                counts: {
                  clients: 65,
                  entity_clients: 50,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/approle/login',
                counts: {
                  clients: 80,
                  entity_clients: 60,
                  non_entity_clients: 20,
                },
              },
            ],
          },
        ],
        months: [
          {
            month_year: 'January_2022',
            counts: {
              clients: 150,
              entity_clients: 100,
              non_entity_clients: 50,
            },
            namespaces: [
              {
                id: 'root',
                path: '',
                counts: {
                  clients: 75,
                  entity_clients: 50,
                  non_entity_clients: 25,
                },
                mounts: [
                  {
                    path: 'auth/aws/login',
                    counts: {
                      clients: 37,
                      entity_clients: 25,
                      non_entity_clients: 12,
                    },
                  },
                  {
                    path: 'auth/approle/login',
                    counts: {
                      clients: 38,
                      entity_clients: 25,
                      non_entity_clients: 13,
                    },
                  },
                ],
              },
              {
                id: 'ns1',
                path: '',
                counts: {
                  clients: 75,
                  entity_clients: 50,
                  non_entity_clients: 25,
                },
                mounts: [
                  {
                    path: 'auth/aws/login',
                    counts: {
                      clients: 30,
                      entity_clients: 20,
                      non_entity_clients: 10,
                    },
                  },
                  {
                    path: 'auth/approle/login',
                    counts: {
                      clients: 45,
                      entity_clients: 30,
                      non_entity_clients: 15,
                    },
                  },
                ],
              },
            ],
            new_clients: {
              counts: {
                clients: 40,
                entity_clients: 30,
                non_entity_clients: 10,
              },
              namespaces: [
                {
                  id: 'root',
                  path: '',
                  counts: {
                    clients: 20,
                    entity_clients: 15,
                    non_entity_clients: 5,
                  },
                  mounts: [
                    {
                      path: 'auth/aws/login',
                      counts: {
                        clients: 7,
                        entity_clients: 5,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/aws/login',
                      counts: {
                        clients: 13,
                        entity_clients: 10,
                        non_entity_clients: 3,
                      },
                    },
                  ],
                },
                {
                  id: 'ns1',
                  path: '',
                  counts: {
                    clients: 20,
                    entity_clients: 15,
                    non_entity_clients: 5,
                  },
                  mounts: [
                    {
                      path: 'auth/aws/login',
                      counts: {
                        clients: 11,
                        entity_clients: 10,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/aws/login',
                      counts: {
                        clients: 9,
                        entity_clients: 5,
                        non_entity_clients: 4,
                      },
                    },
                  ],
                },
              ],
            },
          },
          {
            month_year: 'February_2022',
            counts: {
              _comment: 'total monthly clients',
              clients: 150,
              entity_clients: 100,
              non_entity_clients: 50,
            },
            namespaces: [
              {
                id: 'root',
                path: '',
                counts: {
                  clients: 70,
                  entity_clients: 60,
                  non_entity_clients: 10,
                },
                mounts: [
                  {
                    path: 'auth/aws/login',
                    counts: {
                      clients: 35,
                      entity_clients: 30,
                      non_entity_clients: 5,
                    },
                  },
                  {
                    path: 'auth/approle/login',
                    counts: {
                      clients: 35,
                      entity_clients: 30,
                      non_entity_clients: 5,
                    },
                  },
                ],
              },
              {
                id: 'ns1',
                path: '',
                counts: {
                  clients: 80,
                  entity_clients: 40,
                  non_entity_clients: 40,
                },
                mounts: [
                  {
                    path: 'auth/aws/login',
                    counts: {
                      clients: 40,
                      entity_clients: 20,
                      non_entity_clients: 20,
                    },
                  },
                  {
                    path: 'auth/approle/login',
                    counts: {
                      clients: 40,
                      entity_clients: 20,
                      non_entity_clients: 20,
                    },
                  },
                ],
              },
            ],
            new_clients: {
              counts: {
                clients: 25,
                entity_clients: 20,
                non_entity_clients: 5,
              },
              namespaces: [
                {
                  id: 'root',
                  path: '',
                  counts: {
                    clients: 13,
                    entity_clients: 10,
                    non_entity_clients: 3,
                  },
                  mounts: [
                    {
                      path: 'auth/aws/login',
                      counts: {
                        clients: 6,
                        entity_clients: 5,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/aws/login',
                      counts: {
                        clients: 7,
                        entity_clients: 5,
                        non_entity_clients: 2,
                      },
                    },
                  ],
                },
                {
                  id: 'ns1',
                  path: '',
                  counts: {
                    clients: 12,
                    entity_clients: 10,
                    non_entity_clients: 2,
                  },
                  mounts: [
                    {
                      path: 'auth/aws/login',
                      counts: {
                        clients: 7,
                        entity_clients: 5,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/aws/login',
                      counts: {
                        clients: 5,
                        entity_clients: 5,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                },
              ],
            },
          },
        ],
      },
    };
  });

  // current (<= 1.9) API response
  server.get('/sys/internal/counters/activity/monthly', function () {
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
}
