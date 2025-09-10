/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { findAll } from '@ember/test-helpers';
import { CHARTS } from './client-count-selectors';

export function assertBarChart(assert, chartName, byMonthData, isStacked = false) {
  // assertion count is byMonthData.length, plus 2
  const chart = CHARTS.chart(chartName);
  const dataBars = findAll(`${chart} ${CHARTS.verticalBar}`).filter(
    (b) => b.hasAttribute('height') && b.getAttribute('height') !== '0'
  );
  const xAxisLabels = findAll(`${chart} ${CHARTS.xAxisLabel}`);

  let count = byMonthData.filter((m) => m.clients).length;
  if (isStacked) count = count * 2;

  assert.strictEqual(dataBars.length, count, `${chartName}: it renders bars for each non-zero month`);
  assert.strictEqual(
    xAxisLabels.length,
    byMonthData.length,
    `${chartName}: it renders a label for each month`
  );

  xAxisLabels.forEach((e, i) => {
    assert.dom(e).hasText(`${byMonthData[i].month}`, `renders x-axis label: ${byMonthData[i].month}`);
  });
}

export const ACTIVITY_RESPONSE_STUB = {
  start_time: '2023-06-01T00:00:00Z',
  end_time: '2023-09-30T23:59:59Z', // is always the last day and hour of the month queried
  by_namespace: [
    {
      namespace_id: 'e67m31',
      namespace_path: 'ns1/',
      counts: {
        acme_clients: 5699,
        clients: 18903,
        entity_clients: 4256,
        non_entity_clients: 4138,
        secret_syncs: 4810,
      },
      mounts: [
        {
          mount_path: 'auth/userpass/0/',
          mount_type: 'userpass/',
          counts: {
            acme_clients: 0,
            clients: 8394,
            entity_clients: 4256,
            non_entity_clients: 4138,
            secret_syncs: 0,
          },
        },
        {
          mount_path: 'acme/pki/0/',
          mount_type: 'pki/',
          counts: {
            acme_clients: 5699,
            clients: 5699,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 0,
          },
        },
        {
          mount_path: 'secrets/kv/0/',
          mount_type: 'kv/',
          counts: {
            acme_clients: 0,
            clients: 4810,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 4810,
          },
        },
      ],
    },
    {
      namespace_id: 'root',
      namespace_path: '',
      counts: {
        acme_clients: 4003,
        clients: 16384,
        entity_clients: 4002,
        non_entity_clients: 4089,
        secret_syncs: 4290,
      },
      mounts: [
        {
          mount_path: 'auth/userpass/0/',
          mount_type: 'userpass/',
          counts: {
            acme_clients: 0,
            clients: 8091,
            entity_clients: 4002,
            non_entity_clients: 4089,
            secret_syncs: 0,
          },
        },
        {
          mount_path: 'secrets/kv/0/',
          mount_type: 'kv/',
          counts: {
            acme_clients: 0,
            clients: 4290,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 4290,
          },
        },
        {
          mount_path: 'acme/pki/0/',
          mount_type: 'pki/',
          counts: {
            acme_clients: 4003,
            clients: 4003,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 0,
          },
        },
      ],
    },
  ],
  months: [
    {
      timestamp: '2023-06-01T00:00:00Z',
      counts: null,
      namespaces: null,
      new_clients: null,
    },
    {
      timestamp: '2023-07-01T00:00:00Z',
      counts: {
        acme_clients: 100,
        clients: 400,
        entity_clients: 100,
        non_entity_clients: 100,
        secret_syncs: 100,
      },
      namespaces: [
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            acme_clients: 100,
            clients: 400,
            entity_clients: 100,
            non_entity_clients: 100,
            secret_syncs: 100,
          },
          mounts: [
            {
              mount_path: 'auth/userpass/0/',
              mount_type: 'userpass/',
              counts: {
                acme_clients: 0,
                clients: 200,
                entity_clients: 100,
                non_entity_clients: 100,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'acme/pki/0/',
              mount_type: 'pki/',
              counts: {
                acme_clients: 100,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'secrets/kv/0/',
              mount_type: 'kv/',
              counts: {
                acme_clients: 0,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 100,
              },
            },
          ],
        },
      ],
      new_clients: {
        counts: {
          acme_clients: 100,
          clients: 400,
          entity_clients: 100,
          non_entity_clients: 100,
          secret_syncs: 100,
        },
        namespaces: [
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              acme_clients: 100,
              clients: 400,
              entity_clients: 100,
              non_entity_clients: 100,
              secret_syncs: 100,
            },
            mounts: [
              {
                mount_path: 'auth/userpass/0/',
                mount_type: 'userpass/',
                counts: {
                  acme_clients: 0,
                  clients: 200,
                  entity_clients: 100,
                  non_entity_clients: 100,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'acme/pki/0/',
                mount_type: 'pki/',
                counts: {
                  acme_clients: 100,
                  clients: 100,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'secrets/kv/0/',
                mount_type: 'kv/',
                counts: {
                  acme_clients: 0,
                  clients: 100,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 100,
                },
              },
            ],
          },
        ],
      },
    },
    {
      timestamp: '2023-08-01T00:00:00Z',
      counts: {
        acme_clients: 100,
        clients: 400,
        entity_clients: 100,
        non_entity_clients: 100,
        secret_syncs: 100,
      },
      namespaces: [
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            acme_clients: 100,
            clients: 400,
            entity_clients: 100,
            non_entity_clients: 100,
            secret_syncs: 100,
          },
          mounts: [
            {
              mount_path: 'auth/userpass/0/',
              mount_type: 'userpass/',
              counts: {
                acme_clients: 0,
                clients: 200,
                entity_clients: 100,
                non_entity_clients: 100,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'acme/pki/0/',
              mount_type: 'pki/',
              counts: {
                acme_clients: 100,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'secrets/kv/0/',
              mount_type: 'kv/',
              counts: {
                acme_clients: 0,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 100,
              },
            },
          ],
        },
      ],
      new_clients: {
        counts: null,
        namespaces: null,
      },
    },
    {
      timestamp: '2023-09-01T00:00:00Z',
      counts: {
        acme_clients: 1928,
        clients: 3928,
        entity_clients: 832,
        non_entity_clients: 930,
        secret_syncs: 238,
      },
      namespaces: [
        {
          namespace_id: 'e67m31',
          namespace_path: 'ns1/',
          counts: {
            acme_clients: 934,
            clients: 1981,
            entity_clients: 708,
            non_entity_clients: 182,
            secret_syncs: 157,
          },
          mounts: [
            {
              mount_path: 'acme/pki/0/',
              counts: {
                acme_clients: 934,
                clients: 934,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'auth/userpass/0/',
              counts: {
                acme_clients: 0,
                clients: 890,
                entity_clients: 708,
                non_entity_clients: 182,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'secrets/kv/0/',
              counts: {
                acme_clients: 0,
                clients: 157,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 157,
              },
            },
          ],
        },
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            acme_clients: 994,
            clients: 1947,
            entity_clients: 124,
            non_entity_clients: 748,
            secret_syncs: 81,
          },
          mounts: [
            {
              mount_path: 'acme/pki/0/',
              counts: {
                acme_clients: 994,
                clients: 994,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'auth/userpass/0/',
              counts: {
                acme_clients: 0,
                clients: 872,
                entity_clients: 124,
                non_entity_clients: 748,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'secrets/kv/0/',
              counts: {
                acme_clients: 0,
                clients: 81,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 81,
              },
            },
          ],
        },
      ],
      new_clients: {
        counts: {
          acme_clients: 144,
          clients: 364,
          entity_clients: 59,
          non_entity_clients: 112,
          secret_syncs: 49,
        },
        namespaces: [
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              acme_clients: 91,
              clients: 191,
              entity_clients: 25,
              non_entity_clients: 50,
              secret_syncs: 25,
            },
            mounts: [
              {
                mount_path: 'acme/pki/0/',
                mount_type: 'pki/',
                counts: {
                  acme_clients: 91,
                  clients: 91,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'auth/userpass/0/',
                mount_type: 'userpass/',
                counts: {
                  acme_clients: 0,
                  clients: 75,
                  entity_clients: 25,
                  non_entity_clients: 50,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'secrets/kv/0/',
                mount_type: 'kv/',
                counts: {
                  acme_clients: 0,
                  clients: 25,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 25,
                },
              },
            ],
          },
          {
            namespace_id: 'e67m31',
            namespace_path: 'ns1/',
            counts: {
              acme_clients: 53,
              clients: 173,
              entity_clients: 34,
              non_entity_clients: 62,
              secret_syncs: 24,
            },
            mounts: [
              {
                mount_path: 'auth/userpass/0/',
                mount_type: 'userpass/',
                counts: {
                  acme_clients: 0,
                  clients: 96,
                  entity_clients: 34,
                  non_entity_clients: 62,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'acme/pki/0/',
                mount_type: 'pki/',
                counts: {
                  acme_clients: 53,
                  clients: 53,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'secrets/kv/0/',
                mount_type: 'kv/',
                counts: {
                  acme_clients: 0,
                  clients: 24,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 24,
                },
              },
            ],
          },
        ],
      },
    },
  ],
  total: {
    acme_clients: 9702,
    clients: 35287,
    entity_clients: 8258,
    non_entity_clients: 8227,
    secret_syncs: 9100,
  },
};

// combined activity data before and after 1.10 upgrade when Vault added mount attribution
export const MIXED_ACTIVITY_RESPONSE_STUB = {
  start_time: '2024-03-01T00:00:00Z',
  end_time: '2024-04-30T23:59:59Z',
  total: {
    acme_clients: 0,
    clients: 3,
    entity_clients: 3,
    non_entity_clients: 0,
    secret_syncs: 0,
  },
  by_namespace: [
    {
      counts: {
        acme_clients: 0,
        clients: 3,
        entity_clients: 3,
        non_entity_clients: 0,
        secret_syncs: 0,
      },
      mounts: [
        {
          counts: {
            acme_clients: 0,
            clients: 2,
            entity_clients: 2,
            non_entity_clients: 0,
            secret_syncs: 0,
          },
          mount_path: 'no mount accessor (pre-1.10 upgrade?)',
          mount_type: 'no mount path (pre-1.10 upgrade?)',
        },
        {
          counts: {
            acme_clients: 0,
            clients: 1,
            entity_clients: 1,
            non_entity_clients: 0,
            secret_syncs: 0,
          },
          mount_path: 'auth/userpass/0/',
          mount_type: 'userpass',
        },
      ],
      namespace_id: 'root',
      namespace_path: '',
    },
  ],
  months: [
    {
      counts: null,
      namespaces: null,
      new_clients: null,
      timestamp: '2024-03-01T00:00:00Z',
    },
    {
      counts: {
        acme_clients: 0,
        clients: 3,
        entity_clients: 3,
        non_entity_clients: 0,
        secret_syncs: 0,
      },
      namespaces: [
        {
          counts: {
            acme_clients: 0,
            clients: 3,
            entity_clients: 3,
            non_entity_clients: 0,
            secret_syncs: 0,
          },
          mounts: [
            {
              counts: {
                acme_clients: 0,
                clients: 2,
                entity_clients: 2,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              mount_path: 'no mount accessor (pre-1.10 upgrade?)',
              mount_type: 'no mount path (pre-1.10 upgrade?)',
            },
            {
              counts: {
                acme_clients: 0,
                clients: 1,
                entity_clients: 1,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              mount_path: 'auth/userpass/0/',
              mount_type: 'userpass',
            },
          ],
          namespace_id: 'root',
          namespace_path: '',
        },
      ],
      new_clients: {
        counts: {
          acme_clients: 0,
          clients: 3,
          entity_clients: 3,
          non_entity_clients: 0,
          secret_syncs: 0,
        },
        namespaces: [
          {
            counts: {
              acme_clients: 0,
              clients: 3,
              entity_clients: 3,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            mounts: [
              {
                counts: {
                  acme_clients: 0,
                  clients: 2,
                  entity_clients: 2,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
                mount_path: 'no mount accessor (pre-1.10 upgrade?)',
                mount_type: 'no mount path (pre-1.10 upgrade?)',
              },
              {
                counts: {
                  acme_clients: 0,
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
                mount_path: 'auth/userpass/0/',
                mount_type: 'userpass',
              },
            ],
            namespace_id: 'root',
            namespace_path: '',
          },
        ],
      },
      timestamp: '2024-04-01T00:00:00Z',
    },
  ],
};

// order of this array matters because index 0 is a month without data
export const SERIALIZED_ACTIVITY_RESPONSE = {
  total: {
    acme_clients: 9702,
    clients: 35287,
    entity_clients: 8258,
    non_entity_clients: 8227,
    secret_syncs: 9100,
  },
  by_namespace: [
    {
      label: 'ns1/',
      acme_clients: 5699,
      clients: 18903,
      entity_clients: 4256,
      non_entity_clients: 4138,
      secret_syncs: 4810,
      mounts: [
        {
          label: 'auth/userpass/0/',
          mount_path: 'auth/userpass/0/',
          mount_type: 'userpass',
          namespace_path: 'ns1/',
          acme_clients: 0,
          clients: 8394,
          entity_clients: 4256,
          non_entity_clients: 4138,
          secret_syncs: 0,
        },
        {
          label: 'acme/pki/0/',
          mount_path: 'acme/pki/0/',
          mount_type: 'pki',
          namespace_path: 'ns1/',
          acme_clients: 5699,
          clients: 5699,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 0,
        },
        {
          label: 'secrets/kv/0/',
          mount_path: 'secrets/kv/0/',
          mount_type: 'kv',
          namespace_path: 'ns1/',
          acme_clients: 0,
          clients: 4810,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 4810,
        },
      ],
    },
    {
      label: 'root',
      acme_clients: 4003,
      clients: 16384,
      entity_clients: 4002,
      non_entity_clients: 4089,
      secret_syncs: 4290,
      mounts: [
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
        {
          label: 'secrets/kv/0/',
          mount_path: 'secrets/kv/0/',
          mount_type: 'kv',
          namespace_path: 'root',
          acme_clients: 0,
          clients: 4290,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 4290,
        },
        {
          label: 'acme/pki/0/',
          mount_path: 'acme/pki/0/',
          mount_type: 'pki',
          namespace_path: 'root',
          acme_clients: 4003,
          clients: 4003,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 0,
        },
      ],
    },
  ],
  by_month: [
    {
      timestamp: '2023-06-01T00:00:00Z',
      namespaces: [],
      new_clients: {
        timestamp: '2023-06-01T00:00:00Z',
        namespaces: [],
      },
    },
    {
      timestamp: '2023-07-01T00:00:00Z',
      acme_clients: 100,
      clients: 400,
      entity_clients: 100,
      non_entity_clients: 100,
      secret_syncs: 100,
      namespaces: [
        {
          label: 'root',
          acme_clients: 100,
          clients: 400,
          entity_clients: 100,
          non_entity_clients: 100,
          secret_syncs: 100,
          mounts: [
            {
              label: 'auth/userpass/0/',
              namespace_path: 'root',
              mount_path: 'auth/userpass/0/',
              mount_type: 'userpass',
              acme_clients: 0,
              clients: 200,
              entity_clients: 100,
              non_entity_clients: 100,
              secret_syncs: 0,
            },
            {
              label: 'acme/pki/0/',
              namespace_path: 'root',
              mount_path: 'acme/pki/0/',
              mount_type: 'pki',
              acme_clients: 100,
              clients: 100,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            {
              label: 'secrets/kv/0/',
              namespace_path: 'root',
              mount_path: 'secrets/kv/0/',
              mount_type: 'kv',
              acme_clients: 0,
              clients: 100,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 100,
            },
          ],
        },
      ],
      new_clients: {
        timestamp: '2023-07-01T00:00:00Z',
        acme_clients: 100,
        clients: 400,
        entity_clients: 100,
        non_entity_clients: 100,
        secret_syncs: 100,
        namespaces: [
          {
            label: 'root',
            acme_clients: 100,
            clients: 400,
            entity_clients: 100,
            non_entity_clients: 100,
            secret_syncs: 100,
            mounts: [
              {
                label: 'auth/userpass/0/',
                mount_path: 'auth/userpass/0/',
                mount_type: 'userpass',
                namespace_path: 'root',
                acme_clients: 0,
                clients: 200,
                entity_clients: 100,
                non_entity_clients: 100,
                secret_syncs: 0,
              },
              {
                label: 'acme/pki/0/',
                mount_path: 'acme/pki/0/',
                namespace_path: 'root',
                mount_type: 'pki',
                acme_clients: 100,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                label: 'secrets/kv/0/',
                mount_path: 'secrets/kv/0/',
                mount_type: 'kv',
                namespace_path: 'root',
                acme_clients: 0,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 100,
              },
            ],
          },
        ],
      },
    },
    {
      timestamp: '2023-08-01T00:00:00Z',
      acme_clients: 100,
      clients: 400,
      entity_clients: 100,
      non_entity_clients: 100,
      secret_syncs: 100,
      namespaces: [
        {
          label: 'root',
          acme_clients: 100,
          clients: 400,
          entity_clients: 100,
          non_entity_clients: 100,
          secret_syncs: 100,
          mounts: [
            {
              label: 'auth/userpass/0/',
              mount_path: 'auth/userpass/0/',
              namespace_path: 'root',
              mount_type: 'userpass',
              acme_clients: 0,
              clients: 200,
              entity_clients: 100,
              non_entity_clients: 100,
              secret_syncs: 0,
            },
            {
              label: 'acme/pki/0/',
              mount_path: 'acme/pki/0/',
              namespace_path: 'root',
              mount_type: 'pki',
              acme_clients: 100,
              clients: 100,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },

            {
              label: 'secrets/kv/0/',
              mount_path: 'secrets/kv/0/',
              namespace_path: 'root',
              mount_type: 'kv',
              acme_clients: 0,
              clients: 100,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 100,
            },
          ],
        },
      ],
      new_clients: {
        timestamp: '2023-08-01T00:00:00Z',
        namespaces: [],
      },
    },
    {
      timestamp: '2023-09-01T00:00:00Z',
      acme_clients: 1928,
      clients: 3928,
      entity_clients: 832,
      non_entity_clients: 930,
      secret_syncs: 238,
      namespaces: [
        {
          label: 'ns1/',
          acme_clients: 934,
          clients: 1981,
          entity_clients: 708,
          non_entity_clients: 182,
          secret_syncs: 157,
          mounts: [
            {
              label: 'acme/pki/0/',
              mount_path: 'acme/pki/0/',
              acme_clients: 934,
              clients: 934,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            {
              label: 'auth/userpass/0/',
              mount_path: 'auth/userpass/0/',
              acme_clients: 0,
              clients: 890,
              entity_clients: 708,
              non_entity_clients: 182,
              secret_syncs: 0,
            },
            {
              label: 'secrets/kv/0/',
              mount_path: 'secrets/kv/0/',
              acme_clients: 0,
              clients: 157,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 157,
            },
          ],
        },
        {
          label: 'root',
          acme_clients: 994,
          clients: 1947,
          entity_clients: 124,
          non_entity_clients: 748,
          secret_syncs: 81,
          mounts: [
            {
              label: 'acme/pki/0/',
              mount_path: 'acme/pki/0/',
              acme_clients: 994,
              clients: 994,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            {
              label: 'auth/userpass/0/',
              mount_path: 'auth/userpass/0/',
              acme_clients: 0,
              clients: 872,
              entity_clients: 124,
              non_entity_clients: 748,
              secret_syncs: 0,
            },
            {
              label: 'secrets/kv/0/',
              mount_path: 'secrets/kv/0/',
              acme_clients: 0,
              clients: 81,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 81,
            },
          ],
        },
      ],
      new_clients: {
        timestamp: '2023-09-01T00:00:00Z',
        acme_clients: 144,
        clients: 364,
        entity_clients: 59,
        non_entity_clients: 112,
        secret_syncs: 49,
        namespaces: [
          {
            label: 'root',
            acme_clients: 91,
            clients: 191,
            entity_clients: 25,
            non_entity_clients: 50,
            secret_syncs: 25,
            mounts: [
              {
                label: 'acme/pki/0/',
                mount_path: 'acme/pki/0/',
                mount_type: 'pki',
                acme_clients: 91,
                clients: 91,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                label: 'auth/userpass/0/',
                mount_path: 'auth/userpass/0/',
                mount_type: 'userpass',
                acme_clients: 0,
                clients: 75,
                entity_clients: 25,
                non_entity_clients: 50,
                secret_syncs: 0,
              },
              {
                label: 'secrets/kv/0/',
                mount_path: 'secrets/kv/0/',
                mount_type: 'kv',
                acme_clients: 0,
                clients: 25,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 25,
              },
            ],
          },
          {
            label: 'ns1/',
            acme_clients: 53,
            clients: 173,
            entity_clients: 34,
            non_entity_clients: 62,
            secret_syncs: 24,
            mounts: [
              {
                label: 'auth/userpass/0/',
                mount_path: 'auth/userpass/0/',
                mount_type: 'userpass',
                acme_clients: 0,
                clients: 96,
                entity_clients: 34,
                non_entity_clients: 62,
                secret_syncs: 0,
              },
              {
                label: 'acme/pki/0/',
                mount_path: 'acme/pki/0/',
                mount_type: 'pki',
                acme_clients: 53,
                clients: 53,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                label: 'secrets/kv/0/',
                mount_path: 'secrets/kv/0/',
                mount_type: 'kv',
                acme_clients: 0,
                clients: 24,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 24,
              },
            ],
          },
        ],
      },
    },
  ],
};

export const ENTITY_EXPORT = `{"entity_name":"entity_b3e2a7ff","entity_alias_name":"bob","local_entity_alias":false,"client_id":"5692c6ef-c871-128e-fb06-df2be7bfc0db","client_type":"entity","namespace_id":"vK5Bt","namespace_path":"ns1/","mount_accessor":"auth_userpass_f47ad0b4","mount_type":"userpass","mount_path":"auth/userpass/0/","token_creation_time":"2022-09-15T23:48:09Z","client_first_used_time":"2023-09-15T23:48:09Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":["7537e6b7-3b06-65c2-1fb2-c83116eb5e6f"]}
{"entity_name":"entity_b3e2a7ff","entity_alias_name":"bob","local_entity_alias":false,"client_id":"daf8420c-0b6b-34e6-ff38-ee1ed093bea9","client_type":"entity","namespace_id":"root","namespace_path":"","mount_accessor":"auth_userpass_f47ad0b4","mount_type":"userpass","mount_path":"auth/userpass/","token_creation_time":"2020-08-15T23:48:09Z","client_first_used_time":"2025-07-15T23:48:09Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":["7537e6b7-3b06-65c2-1fb2-c83116eb5e6f"]}
{"entity_name":"bob-smith","entity_alias_name":"bob","local_entity_alias":false,"client_id":"23a04911-5d72-ba98-11d3-527f2fcf3a81","client_type":"entity","namespace_id":"root","namespace_path":"","mount_accessor":"auth_userpass_de28062c","mount_type":"userpass","mount_path":"auth/userpass-test/","token_creation_time":"2020-08-15T23:52:38Z","client_first_used_time":"2025-08-15T23:53:19Z","policies":["base"],"entity_metadata":{"organization":"ACME Inc.","team":"QA"},"entity_alias_metadata":{},"entity_alias_custom_metadata":{"account":"Tester Account"},"entity_group_ids":["7537e6b7-3b06-65c2-1fb2-c83116eb5e6f"]}
{"entity_name":"alice-johnson","entity_alias_name":"alice","local_entity_alias":false,"client_id":"a7c8d912-4f61-23b5-88e4-627a3dcf2b92","client_type":"entity","namespace_id":"root","namespace_path":"","mount_accessor":"auth_userpass_f47ad0b4","mount_type":"userpass","mount_path":"auth/userpass/","token_creation_time":"2020-08-16T09:15:42Z","client_first_used_time":"2025-09-16T09:16:03Z","policies":["admin","audit"],"entity_metadata":{"organization":"TechCorp","team":"DevOps","location":"San Francisco"},"entity_alias_metadata":{"department":"Engineering"},"entity_alias_custom_metadata":{"role":"Senior Engineer"},"entity_group_ids":["7537e6b7-3b06-65c2-1fb2-c83116eb5e6f","a1b2c3d4-5e6f-7g8h-9i0j-k1l2m3n4o5p6"]}
{"entity_name":"charlie-brown","entity_alias_name":"charlie","local_entity_alias":true,"client_id":"b9e5f824-7c92-34d6-a1f8-738b4ecf5d73","client_type":"entity","namespace_id":"whUNi","namespace_path":"ns2/","mount_accessor":"auth_ldap_8a3b9c2d","mount_type":"ldap","mount_path":"auth/ldap/","token_creation_time":"2020-08-16T14:22:17Z","client_first_used_time":"2025-10-16T14:22:45Z","policies":["developer","read-only"],"entity_metadata":{"organization":"StartupXYZ","team":"Backend"},"entity_alias_metadata":{"cn":"charlie.brown","ou":"development"},"entity_alias_custom_metadata":{"project":"microservices"},"entity_group_ids":["c7d8e9f0-1a2b-3c4d-5e6f-789012345678"]}
{"entity_name":"diana-prince","entity_alias_name":"diana","local_entity_alias":false,"client_id":"e4f7a935-2b68-47c9-b3e6-849c5dfb7a84","client_type":"entity","namespace_id":"aT9S5","namespace_path":"ns1/","mount_accessor":"auth_oidc_1f2e3d4c","mount_type":"oidc","mount_path":"auth/oidc/","token_creation_time":"2020-08-17T11:08:33Z","client_first_used_time":"2025-11-17T11:09:01Z","policies":["security","compliance"],"entity_metadata":{"organization":"SecureTech","team":"Security","clearance":"high"},"entity_alias_metadata":{"email":"diana.prince@securetech.com"},"entity_alias_custom_metadata":{"access_level":"L4"},"entity_group_ids":["f8e7d6c5-4b3a-2918-7654-321098765432"]}
{"entity_name":"frank-castle","entity_alias_name":"frank","local_entity_alias":false,"client_id":"c6b9d248-5a71-39e4-c7f2-951d8eaf6b95","client_type":"entity","namespace_id":"root","namespace_path":"","mount_accessor":"auth_jwt_9d8c7b6a","mount_type":"jwt","mount_path":"auth/jwt/","token_creation_time":"2020-08-17T16:43:28Z","client_first_used_time":"2025-12-17T16:44:12Z","policies":["operations","monitoring"],"entity_metadata":{"organization":"CloudOps","team":"SRE","region":"us-east-1"},"entity_alias_metadata":{"sub":"frank.castle@cloudops.io","iss":"https://auth.cloudops.io"},"entity_alias_custom_metadata":{"on_call":"true","expertise":"kubernetes"},"entity_group_ids":["9a8b7c6d-5e4f-3210-9876-543210fedcba"]}
{"entity_name":"grace-hopper","entity_alias_name":"grace","local_entity_alias":true,"client_id":"d8a3e517-6f94-42b7-d5c8-062f9bce4a73","client_type":"entity","namespace_id":"YMjS8","namespace_path":"ns5/","mount_accessor":"auth_userpass_3e2d1c0b","mount_type":"userpass","mount_path":"auth/userpass-legacy/","token_creation_time":"2020-08-18T08:17:55Z","client_first_used_time":"2025-06-18T08:18:23Z","policies":["legacy-admin","data-access"],"entity_metadata":{"organization":"LegacySystems","team":"Platform","tenure":"senior"},"entity_alias_metadata":{"legacy_id":"grace.hopper.001"},"entity_alias_custom_metadata":{"system_access":"mainframe","certification":"vault-admin"},"entity_group_ids":["1f2e3d4c-5b6a-7980-1234-567890abcdef"]}
`;

const NON_ENTITY_EXPORT = `{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"46dcOXXH+P1VEQiKTQjtWXEtBlbHdMOWwz+svXf3xuU=","client_type":"non-entity-token","namespace_id":"whUNi","namespace_path":"ns2/","mount_accessor":"auth_ns_token_3b2bf405","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:21Z","client_first_used_time":"2025-05-15T16:19:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"VKAJVITyTwyqF1GUzwYHwkaK6bbnL1zN8ZJ7viKR8no=","client_type":"non-entity-token","namespace_id":"omjn8","namespace_path":"ns8/","mount_accessor":"auth_ns_token_07b90be7","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:22Z","client_first_used_time":"2025-05-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"ww4L5n9WE32lPNh3UBgT3JxTDZb1a+m/3jqUffp04tQ=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:23Z","client_first_used_time":"2025-05-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"cBLb9erIROCw7cczXpfkXTOdnZoVwfWF4EAPD9k61lU=","client_type":"non-entity-token","namespace_id":"aT9S5","namespace_path":"ns1/","mount_accessor":"auth_ns_token_62a4e52a","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:21Z","client_first_used_time":"2025-06-15T16:19:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"KMHoH3Kvr6nnW2ZIs+i37pYvyVtnuaL3DmyVxUL6boI=","client_type":"non-entity-token","namespace_id":"YMjS8","namespace_path":"ns5/","mount_accessor":"auth_ns_token_45cbc810","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:22Z","client_first_used_time":"2025-06-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"hcMH4P4IGAN13cJqkwIJLXYoPLTodtOj/wPTZKS0x4U=","client_type":"non-entity-token","namespace_id":"ZNdL5","namespace_path":"ns7/","mount_accessor":"auth_ns_token_8bbd9440","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:22Z","client_first_used_time":"2025-06-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Oby0ABLmfhqYdfqGfljGHHhAA5zX+BwsGmFu4QGJZd0=","client_type":"non-entity-token","namespace_id":"bJIgY","namespace_path":"ns9/","mount_accessor":"auth_ns_token_8d188479","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:22Z","client_first_used_time":"2025-06-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Z6MjZuH/VD7HU11efiKoM/hfoxssSbeu4c6DhC7zUZ4=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:23Z","client_first_used_time":"2025-07-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"1UxaPHJUOPWrf0ivMgBURK6WHzbfXGkcn/C/xI3AeHQ=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:24Z","client_first_used_time":"2025-07-15T16:19:24Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"hfFbwhMucs/f84p2QTOiBLT72i0WLVkIgCGV7RIuWlo=","client_type":"non-entity-token","namespace_id":"x6sKN","namespace_path":"ns4/","mount_accessor":"auth_ns_token_2aaebdc2","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:21Z","client_first_used_time":"2025-07-15T16:19:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"sOdIr+zoNqOUa4hq6Jv4LCGVr0sTLGbvcRPVGAtUA7g=","client_type":"non-entity-token","namespace_id":"Rsvk5","namespace_path":"ns6/","mount_accessor":"auth_ns_token_f603fd8d","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:22Z","client_first_used_time":"2025-07-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"vOIAwNhe6P6HFdJQgUIU/8K6Z5e+oxyVP5x3KtTKS6U=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"ZOkJY3P7IzOqulsnEI0JAQQXwTPnXmpGUh9otqNUclc=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Lsha/HH+xLZq92XG4GYZVlwVQCiqPCUIuoego4aCybU=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Tsl/u7CDTYSXA9HRwlNTW7K/yyEe5PDkLOVTvTWy3q0=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:23Z","client_first_used_time":"2025-09-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"vnq6JntpiGV4FN6GDICLECe2in31aanLA6Q1UWqBmL0=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:24Z","client_first_used_time":"2025-09-15T16:19:24Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"MRMrywfPPL3QnKFMBGfRjjmaefBRH1VKpQVIfrd0Xb4=","client_type":"non-entity-token","namespace_id":"6aDiU","namespace_path":"ns3/","mount_accessor":"auth_ns_token_ef771c23","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:21Z","client_first_used_time":"2025-09-15T16:19:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Rce6fjHs15+hDl5XdXbWmzGNYrTcQsJuaoqfs9Vrhvw=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2020-08-15T16:19:24Z","client_first_used_time":"2025-09-15T16:19:24Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
`;

const ACME_EXPORT = `{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.3U8nSB_yMBvrdu7PvAVykKurDiaH_vQGaEdAUsp-Cew","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:47:54Z","client_first_used_time":"2025-06-21T18:47:54Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.77tKDzxw0i81Nr4XLliTP9xRsztXLTuS16nN32B9jHA","client_type":"pki-acme","namespace_id":"whUNi","namespace_path":"ns2/","mount_accessor":"pki_06dad7b8","mount_type":"ns_pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:48:17Z","client_first_used_time":"2025-07-21T18:48:17Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.RoN77EahLU0wfem4z--ZqJSaqOZ7RvBWR3OkPHM_xaw","client_type":"pki-acme","namespace_id":"omjn8","namespace_path":"ns8/","mount_accessor":"pki_06dad7b8","mount_type":"ns_pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:49:26Z","client_first_used_time":"2025-08-21T18:49:26Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.5S7vaaJIrXormQSLHv4YkBhVfu6Ug0GERhTVTCrq-Fk","client_type":"pki-acme","namespace_id":"aT9S5","namespace_path":"ns1/","mount_accessor":"pki_06dad7b8","mount_type":"ns_pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:45:12Z","client_first_used_time":"2025-08-21T18:45:12Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.IeJQMwtkReJHVNL6fZLmqiu8-Re4JdKCQixXkfcaSRE","client_type":"pki-acme","namespace_id":"YMjS8","namespace_path":"ns5/","mount_accessor":"pki_06dad7b8","mount_type":"ns_pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:45:41Z","client_first_used_time":"2025-08-21T18:45:41Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.vTm3SCZom90qy3SuyIacpVsQgGLx7ASf3SeGpqn5XBA","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:47:19Z","client_first_used_time":"2025-08-21T18:47:19Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.64jWs15k6roUH6MiQ2u80K08Bmqw8IQOpqTpDZgZ1f4","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:47:25Z","client_first_used_time":"2025-08-21T18:47:25Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.RkjnwyIIn6bnc4LDdKQ9HNfnhuVXT7vQONXgGHJl4CE","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:49:21Z","client_first_used_time":"2025-08-21T18:49:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.uozIMLVXDMU7Fc2TFFwq0-uE1GFSui5rbTI1XyNAYBY","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:44:44Z","client_first_used_time":"2025-08-21T18:44:44Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.WiLdlzq93WtVmObB__CC2SPX6sI7EVLTTzxOIRHHN3o","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:44:49Z","client_first_used_time":"2025-08-21T18:44:49Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.P65jgamzwLYbKyxTlJFD5DL3sIUbusbXcQhYaysgzlU","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:45:59Z","client_first_used_time":"2025-08-21T18:45:59Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.2REWUkDLXAG2UB0ZJQcjPnHc4H39aq8fG3LMaHSHKow","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:46:05Z","client_first_used_time":"2025-08-21T18:46:05Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.Eeyq9-EfWv-iE9Aj3DzCU4r9P8V1Maewx51vcxMN-jA","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:46:10Z","client_first_used_time":"2025-08-21T18:46:10Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.vaeb2KR58sRuMUdUlv2TsbaOkSICTAxmJxhkuOs8ZiM","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:46:22Z","client_first_used_time":"2025-08-21T18:46:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.xEPG0eNfrAfRgXg6AKjsCrFPMs0IbLTCfUsCie_rfzY","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:46:51Z","client_first_used_time":"2025-08-21T18:46:51Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"pki-acme.Bkg4862LEoFXJUDWlfFtJHU9a69KRJPiEdw5XCbkkAI","client_type":"pki-acme","namespace_id":"root","namespace_path":"","mount_accessor":"pki_06dad7b8","mount_type":"pki","mount_path":"pki_int/","token_creation_time":"2020-08-21T18:47:42Z","client_first_used_time":"2025-08-21T18:47:42Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
`;

const SECRET_SYNC_EXPORT = `{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"secret-sync.3U8nSB_yMBvrdu7PvAVykKurDiaH_vQGaEdAUsp-Cew","client_type":"secret-sync","namespace_id":"root","namespace_path":"","mount_accessor":"kv_06dad7b8","mount_type":"kv","mount_path":"secrets/kv/0/","token_creation_time":"2020-08-21T18:47:54Z","client_first_used_time":"2025-05-21T18:47:54Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"secret-sync.77tKDzxw0i81Nr4XLliTP9xRsztXLTuS16nN32B9jHA","client_type":"secret-sync","namespace_id":"ZNdL5","namespace_path":"ns7/","mount_accessor":"kv_06dad7b8","mount_type":"ns_kv","mount_path":"secrets/kv/0/","token_creation_time":"2020-08-21T18:48:17Z","client_first_used_time":"2025-05-21T18:48:17Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"secret-sync.RoN77EahLU0wfem4z--ZqJSaqOZ7RvBWR3OkPHM_xaw","client_type":"secret-sync","namespace_id":"bJIgY","namespace_path":"ns9/","mount_accessor":"kv_12abc3d4","mount_type":"ns_kv","mount_path":"secrets/kv/1","token_creation_time":"2020-08-21T18:49:26Z","client_first_used_time":"2025-06-21T18:49:26Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"secret-sync.5S7vaaJIrXormQSLHv4YkBhVfu6Ug0GERhTVTCrq-Fk","client_type":"secret-sync","namespace_id":"x6sKN","namespace_path":"ns4/","mount_accessor":"kv_06dad7b8","mount_type":"ns_kv","mount_path":"secrets/kv/0/","token_creation_time":"2020-08-21T18:45:12Z","client_first_used_time":"2025-06-21T18:45:12Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"secret-sync.IeJQMwtkReJHVNL6fZLmqiu8-Re4JdKCQixXkfcaSRE","client_type":"secret-sync","namespace_id":"Rsvk5","namespace_path":"ns6/","mount_accessor":"kv_12abc3d4","mount_type":"ns_kv","mount_path":"secrets/kv/1","token_creation_time":"2020-08-21T18:45:41Z","client_first_used_time":"2025-07-21T18:45:41Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"secret-sync.vTm3SCZom90qy3SuyIacpVsQgGLx7ASf3SeGpqn5XBA","client_type":"secret-sync","namespace_id":"root","namespace_path":"","mount_accessor":"kv_06dad7b8","mount_type":"kv","mount_path":"secrets/kv/0/","token_creation_time":"2020-08-21T18:47:19Z","client_first_used_time":"2025-08-21T18:47:19Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"secret-sync.64jWs15k6roUH6MiQ2u80K08Bmqw8IQOpqTpDZgZ1f4","client_type":"secret-sync","namespace_id":"6aDiU","namespace_path":"ns3/","mount_accessor":"kv_12abc3d4","mount_type":"kv","mount_path":"secrets/kv/1","token_creation_time":"2020-08-21T18:47:25Z","client_first_used_time":"2025-08-21T18:47:25Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"secret-sync.RkjnwyIIn6bnc4LDdKQ9HNfnhuVXT7vQONXgGHJl4CE","client_type":"secret-sync","namespace_id":"root","namespace_path":"","mount_accessor":"kv_06dad7b8","mount_type":"kv","mount_path":"secrets/kv/0/","token_creation_time":"2020-08-21T18:49:21Z","client_first_used_time":"2025-08-21T18:49:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
`;

export const ACTIVITY_EXPORT_STUB = ENTITY_EXPORT + NON_ENTITY_EXPORT + ACME_EXPORT + SECRET_SYNC_EXPORT;
