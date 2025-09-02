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
      namespace_path: 'ns1',
      counts: {
        acme_clients: 5699,
        clients: 18903,
        entity_clients: 4256,
        non_entity_clients: 4138,
        secret_syncs: 4810,
      },
      mounts: [
        {
          mount_path: 'auth/userpass/0',
          mount_type: 'userpass',
          counts: {
            acme_clients: 0,
            clients: 8394,
            entity_clients: 4256,
            non_entity_clients: 4138,
            secret_syncs: 0,
          },
        },
        {
          mount_path: 'acme/pki/0',
          mount_type: 'pki',
          counts: {
            acme_clients: 5699,
            clients: 5699,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 0,
          },
        },
        {
          mount_path: 'secrets/kv/0',
          mount_type: 'kv',
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
          mount_path: 'auth/userpass/0',
          mount_type: 'userpass',
          counts: {
            acme_clients: 0,
            clients: 8091,
            entity_clients: 4002,
            non_entity_clients: 4089,
            secret_syncs: 0,
          },
        },
        {
          mount_path: 'secrets/kv/0',
          mount_type: 'kv',
          counts: {
            acme_clients: 0,
            clients: 4290,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 4290,
          },
        },
        {
          mount_path: 'acme/pki/0',
          mount_type: 'pki',
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
              mount_path: 'auth/userpass/0',
              mount_type: 'userpass',
              counts: {
                acme_clients: 0,
                clients: 200,
                entity_clients: 100,
                non_entity_clients: 100,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'acme/pki/0',
              mount_type: 'pki',
              counts: {
                acme_clients: 100,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'secrets/kv/0',
              mount_type: 'kv',
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
                mount_path: 'auth/userpass/0',
                mount_type: 'userpass',
                counts: {
                  acme_clients: 0,
                  clients: 200,
                  entity_clients: 100,
                  non_entity_clients: 100,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'acme/pki/0',
                mount_type: 'pki',
                counts: {
                  acme_clients: 100,
                  clients: 100,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'secrets/kv/0',
                mount_type: 'kv',
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
              mount_path: 'auth/userpass/0',
              mount_type: 'userpass',
              counts: {
                acme_clients: 0,
                clients: 200,
                entity_clients: 100,
                non_entity_clients: 100,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'acme/pki/0',
              mount_type: 'pki',
              counts: {
                acme_clients: 100,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'secrets/kv/0',
              mount_type: 'kv',
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
          namespace_path: 'ns1',
          counts: {
            acme_clients: 934,
            clients: 1981,
            entity_clients: 708,
            non_entity_clients: 182,
            secret_syncs: 157,
          },
          mounts: [
            {
              mount_path: 'acme/pki/0',
              counts: {
                acme_clients: 934,
                clients: 934,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'auth/userpass/0',
              counts: {
                acme_clients: 0,
                clients: 890,
                entity_clients: 708,
                non_entity_clients: 182,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'secrets/kv/0',
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
              mount_path: 'acme/pki/0',
              counts: {
                acme_clients: 994,
                clients: 994,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'auth/userpass/0',
              counts: {
                acme_clients: 0,
                clients: 872,
                entity_clients: 124,
                non_entity_clients: 748,
                secret_syncs: 0,
              },
            },
            {
              mount_path: 'secrets/kv/0',
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
                mount_path: 'acme/pki/0',
                mount_type: 'pki',
                counts: {
                  acme_clients: 91,
                  clients: 91,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'auth/userpass/0',
                mount_type: 'userpass',
                counts: {
                  acme_clients: 0,
                  clients: 75,
                  entity_clients: 25,
                  non_entity_clients: 50,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'secrets/kv/0',
                mount_type: 'kv',
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
            namespace_path: 'ns1',
            counts: {
              acme_clients: 53,
              clients: 173,
              entity_clients: 34,
              non_entity_clients: 62,
              secret_syncs: 24,
            },
            mounts: [
              {
                mount_path: 'auth/userpass/0',
                mount_type: 'userpass',
                counts: {
                  acme_clients: 0,
                  clients: 96,
                  entity_clients: 34,
                  non_entity_clients: 62,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'acme/pki/0',
                mount_type: 'pki',
                counts: {
                  acme_clients: 53,
                  clients: 53,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 0,
                },
              },
              {
                mount_path: 'secrets/kv/0',
                mount_type: 'kv',
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
          mount_path: 'auth/userpass/0',
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
              mount_path: 'auth/userpass/0',
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
                mount_path: 'auth/userpass/0',
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
      label: 'ns1',
      acme_clients: 5699,
      clients: 18903,
      entity_clients: 4256,
      non_entity_clients: 4138,
      secret_syncs: 4810,
      mounts: [
        {
          label: 'auth/userpass/0',
          mount_path: 'auth/userpass/0',
          mount_type: 'userpass',
          namespace_path: 'ns1',
          acme_clients: 0,
          clients: 8394,
          entity_clients: 4256,
          non_entity_clients: 4138,
          secret_syncs: 0,
        },
        {
          label: 'acme/pki/0',
          mount_path: 'acme/pki/0',
          mount_type: 'pki',
          namespace_path: 'ns1',
          acme_clients: 5699,
          clients: 5699,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 0,
        },
        {
          label: 'secrets/kv/0',
          mount_path: 'secrets/kv/0',
          mount_type: 'kv',
          namespace_path: 'ns1',
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
          label: 'auth/userpass/0',
          mount_path: 'auth/userpass/0',
          mount_type: 'userpass',
          namespace_path: 'root',
          acme_clients: 0,
          clients: 8091,
          entity_clients: 4002,
          non_entity_clients: 4089,
          secret_syncs: 0,
        },
        {
          label: 'secrets/kv/0',
          mount_path: 'secrets/kv/0',
          mount_type: 'kv',
          namespace_path: 'root',
          acme_clients: 0,
          clients: 4290,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 4290,
        },
        {
          label: 'acme/pki/0',
          mount_path: 'acme/pki/0',
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
              label: 'auth/userpass/0',
              namespace_path: 'root',
              mount_path: 'auth/userpass/0',
              mount_type: 'userpass',
              acme_clients: 0,
              clients: 200,
              entity_clients: 100,
              non_entity_clients: 100,
              secret_syncs: 0,
            },
            {
              label: 'acme/pki/0',
              namespace_path: 'root',
              mount_path: 'acme/pki/0',
              mount_type: 'pki',
              acme_clients: 100,
              clients: 100,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            {
              label: 'secrets/kv/0',
              namespace_path: 'root',
              mount_path: 'secrets/kv/0',
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
                label: 'auth/userpass/0',
                mount_path: 'auth/userpass/0',
                mount_type: 'userpass',
                namespace_path: 'root',
                acme_clients: 0,
                clients: 200,
                entity_clients: 100,
                non_entity_clients: 100,
                secret_syncs: 0,
              },
              {
                label: 'acme/pki/0',
                mount_path: 'acme/pki/0',
                namespace_path: 'root',
                mount_type: 'pki',
                acme_clients: 100,
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                label: 'secrets/kv/0',
                mount_path: 'secrets/kv/0',
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
              label: 'auth/userpass/0',
              mount_path: 'auth/userpass/0',
              namespace_path: 'root',
              mount_type: 'userpass',
              acme_clients: 0,
              clients: 200,
              entity_clients: 100,
              non_entity_clients: 100,
              secret_syncs: 0,
            },
            {
              label: 'acme/pki/0',
              mount_path: 'acme/pki/0',
              namespace_path: 'root',
              mount_type: 'pki',
              acme_clients: 100,
              clients: 100,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },

            {
              label: 'secrets/kv/0',
              mount_path: 'secrets/kv/0',
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
          label: 'ns1',
          acme_clients: 934,
          clients: 1981,
          entity_clients: 708,
          non_entity_clients: 182,
          secret_syncs: 157,
          mounts: [
            {
              label: 'acme/pki/0',
              mount_path: 'acme/pki/0',
              acme_clients: 934,
              clients: 934,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            {
              label: 'auth/userpass/0',
              mount_path: 'auth/userpass/0',
              acme_clients: 0,
              clients: 890,
              entity_clients: 708,
              non_entity_clients: 182,
              secret_syncs: 0,
            },
            {
              label: 'secrets/kv/0',
              mount_path: 'secrets/kv/0',
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
              label: 'acme/pki/0',
              mount_path: 'acme/pki/0',
              acme_clients: 994,
              clients: 994,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
            {
              label: 'auth/userpass/0',
              mount_path: 'auth/userpass/0',
              acme_clients: 0,
              clients: 872,
              entity_clients: 124,
              non_entity_clients: 748,
              secret_syncs: 0,
            },
            {
              label: 'secrets/kv/0',
              mount_path: 'secrets/kv/0',
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
                label: 'acme/pki/0',
                mount_path: 'acme/pki/0',
                mount_type: 'pki',
                acme_clients: 91,
                clients: 91,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                label: 'auth/userpass/0',
                mount_path: 'auth/userpass/0',
                mount_type: 'userpass',
                acme_clients: 0,
                clients: 75,
                entity_clients: 25,
                non_entity_clients: 50,
                secret_syncs: 0,
              },
              {
                label: 'secrets/kv/0',
                mount_path: 'secrets/kv/0',
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
            label: 'ns1',
            acme_clients: 53,
            clients: 173,
            entity_clients: 34,
            non_entity_clients: 62,
            secret_syncs: 24,
            mounts: [
              {
                label: 'auth/userpass/0',
                mount_path: 'auth/userpass/0',
                mount_type: 'userpass',
                acme_clients: 0,
                clients: 96,
                entity_clients: 34,
                non_entity_clients: 62,
                secret_syncs: 0,
              },
              {
                label: 'acme/pki/0',
                mount_path: 'acme/pki/0',
                mount_type: 'pki',
                acme_clients: 53,
                clients: 53,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
              {
                label: 'secrets/kv/0',
                mount_path: 'secrets/kv/0',
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

export const ACTIVITY_EXPORT_STUB = `
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"46dcOXXH+P1VEQiKTQjtWXEtBlbHdMOWwz+svXf3xuU=","client_type":"non-entity-token","namespace_id":"whUNi","namespace_path":"test-ns-2/","mount_accessor":"auth_ns_token_3b2bf405","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:21Z","client_first_used_time":"2025-08-15T16:19:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"VKAJVITyTwyqF1GUzwYHwkaK6bbnL1zN8ZJ7viKR8no=","client_type":"non-entity-token","namespace_id":"omjn8","namespace_path":"test-ns-8/","mount_accessor":"auth_ns_token_07b90be7","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:22Z","client_first_used_time":"2025-08-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"ww4L5n9WE32lPNh3UBgT3JxTDZb1a+m/3jqUffp04tQ=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"cBLb9erIROCw7cczXpfkXTOdnZoVwfWF4EAPD9k61lU=","client_type":"non-entity-token","namespace_id":"aT9S5","namespace_path":"test-ns-1/","mount_accessor":"auth_ns_token_62a4e52a","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:21Z","client_first_used_time":"2025-08-15T16:19:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"KMHoH3Kvr6nnW2ZIs+i37pYvyVtnuaL3DmyVxUL6boI=","client_type":"non-entity-token","namespace_id":"YMjS8","namespace_path":"test-ns-5/","mount_accessor":"auth_ns_token_45cbc810","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:22Z","client_first_used_time":"2025-08-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"hcMH4P4IGAN13cJqkwIJLXYoPLTodtOj/wPTZKS0x4U=","client_type":"non-entity-token","namespace_id":"ZNdL5","namespace_path":"test-ns-7/","mount_accessor":"auth_ns_token_8bbd9440","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:22Z","client_first_used_time":"2025-08-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Oby0ABLmfhqYdfqGfljGHHhAA5zX+BwsGmFu4QGJZd0=","client_type":"non-entity-token","namespace_id":"bJIgY","namespace_path":"test-ns-9/","mount_accessor":"auth_ns_token_8d188479","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:22Z","client_first_used_time":"2025-08-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Z6MjZuH/VD7HU11efiKoM/hfoxssSbeu4c6DhC7zUZ4=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"1UxaPHJUOPWrf0ivMgBURK6WHzbfXGkcn/C/xI3AeHQ=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:24Z","client_first_used_time":"2025-08-15T16:19:24Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"hfFbwhMucs/f84p2QTOiBLT72i0WLVkIgCGV7RIuWlo=","client_type":"non-entity-token","namespace_id":"x6sKN","namespace_path":"test-ns-4/","mount_accessor":"auth_ns_token_2aaebdc2","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:21Z","client_first_used_time":"2025-08-15T16:19:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"sOdIr+zoNqOUa4hq6Jv4LCGVr0sTLGbvcRPVGAtUA7g=","client_type":"non-entity-token","namespace_id":"Rsvk5","namespace_path":"test-ns-6/","mount_accessor":"auth_ns_token_f603fd8d","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:22Z","client_first_used_time":"2025-08-15T16:19:22Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"vOIAwNhe6P6HFdJQgUIU/8K6Z5e+oxyVP5x3KtTKS6U=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"ZOkJY3P7IzOqulsnEI0JAQQXwTPnXmpGUh9otqNUclc=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Lsha/HH+xLZq92XG4GYZVlwVQCiqPCUIuoego4aCybU=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Tsl/u7CDTYSXA9HRwlNTW7K/yyEe5PDkLOVTvTWy3q0=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:23Z","client_first_used_time":"2025-08-15T16:19:23Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"vnq6JntpiGV4FN6GDICLECe2in31aanLA6Q1UWqBmL0=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:24Z","client_first_used_time":"2025-08-15T16:19:24Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"MRMrywfPPL3QnKFMBGfRjjmaefBRH1VKpQVIfrd0Xb4=","client_type":"non-entity-token","namespace_id":"6aDiU","namespace_path":"test-ns-3/","mount_accessor":"auth_ns_token_ef771c23","mount_type":"ns_token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:21Z","client_first_used_time":"2025-08-15T16:19:21Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"","entity_alias_name":"","local_entity_alias":false,"client_id":"Rce6fjHs15+hDl5XdXbWmzGNYrTcQsJuaoqfs9Vrhvw=","client_type":"non-entity-token","namespace_id":"root","namespace_path":"","mount_accessor":"auth_token_360f591b","mount_type":"token","mount_path":"auth/token/","token_creation_time":"2025-08-15T16:19:24Z","client_first_used_time":"2025-08-15T16:19:24Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":[]}
{"entity_name":"entity_b3e2a7ff","entity_alias_name":"bob","local_entity_alias":false,"client_id":"5692c6ef-c871-128e-fb06-df2be7bfc0db","client_type":"entity","namespace_id":"root","namespace_path":"","mount_accessor":"auth_userpass_f47ad0b4","mount_type":"userpass","mount_path":"auth/userpass/","token_creation_time":"2025-08-15T23:48:09Z","client_first_used_time":"2025-08-15T23:48:09Z","policies":[],"entity_metadata":{},"entity_alias_metadata":{},"entity_alias_custom_metadata":{},"entity_group_ids":["7537e6b7-3b06-65c2-1fb2-c83116eb5e6f"]}
{"entity_name":"bob-smith","entity_alias_name":"bob","local_entity_alias":false,"client_id":"23a04911-5d72-ba98-11d3-527f2fcf3a81","client_type":"entity","namespace_id":"root","namespace_path":"","mount_accessor":"auth_userpass_de28062c","mount_type":"userpass","mount_path":"auth/userpass-test/","token_creation_time":"2025-08-15T23:52:38Z","client_first_used_time":"2025-08-15T23:53:19Z","policies":["base"],"entity_metadata":{"organization":"ACME Inc.","team":"QA"},"entity_alias_metadata":{},"entity_alias_custom_metadata":{"account":"Tester Account"},"entity_group_ids":["7537e6b7-3b06-65c2-1fb2-c83116eb5e6f"]}
`;
