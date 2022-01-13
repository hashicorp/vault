export default function (server) {
  server.get('/sys/internal/counters/activity/monthly', function () {
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
}
