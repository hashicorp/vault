/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, findAll } from '@ember/test-helpers';

import { LICENSE_START } from 'vault/mirage/handlers/clients';
import { addMonths } from 'date-fns';
import { CLIENT_COUNT } from './client-count-selectors';

export async function dateDropdownSelect(month, year) {
  const { dateDropdown, counts } = CLIENT_COUNT;
  await click(counts.startEdit);
  await click(dateDropdown.toggleMonth);
  await click(dateDropdown.selectMonth(month));
  await click(dateDropdown.toggleYear);
  await click(dateDropdown.selectYear(year));
  await click(dateDropdown.submit);
}

export function assertChart(assert, chartName, byMonthData) {
  // assertion count is byMonthData.length + 2
  const chart = CLIENT_COUNT.charts.chart(chartName);
  const dataBars = findAll(`${chart} ${CLIENT_COUNT.charts.dataBar}`).filter((b) => b.hasAttribute('height'));
  const xAxisLabels = findAll(`${chart} ${CLIENT_COUNT.charts.xAxisLabel}`);

  assert.strictEqual(
    dataBars.length,
    byMonthData.filter((m) => m.clients).length,
    `${chartName}: it renders bars for each non-zero month`
  );

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
  start_time: '2023-08-01T00:00:00Z',
  end_time: '2023-09-30T23:59:59Z', // is always the last day and hour of the month queried
  by_namespace: [
    {
      namespace_id: 'root',
      namespace_path: '',
      counts: {
        distinct_entities: 1033,
        entity_clients: 1033,
        non_entity_tokens: 1924,
        non_entity_clients: 1924,
        secret_syncs: 2397,
        acme_clients: 75,
        clients: 5429,
      },
      mounts: [
        {
          mount_path: 'auth/authid0',
          counts: {
            clients: 2957,
            entity_clients: 1033,
            non_entity_clients: 1924,
            distinct_entities: 1033,
            non_entity_tokens: 1924,
            secret_syncs: 0,
            acme_clients: 0,
          },
        },
        {
          mount_path: 'kvv2-engine-0',
          counts: {
            clients: 2397,
            entity_clients: 0,
            non_entity_clients: 0,
            distinct_entities: 0,
            non_entity_tokens: 0,
            secret_syncs: 2397,
            acme_clients: 0,
          },
        },
        {
          mount_path: 'pki-engine-0',
          counts: {
            clients: 75,
            entity_clients: 0,
            non_entity_clients: 0,
            distinct_entities: 0,
            non_entity_tokens: 0,
            secret_syncs: 0,
            acme_clients: 75,
          },
        },
      ],
    },
    {
      namespace_id: '81ry61',
      namespace_path: 'ns1',
      counts: {
        distinct_entities: 783,
        entity_clients: 783,
        non_entity_tokens: 1193,
        non_entity_clients: 1193,
        secret_syncs: 275,
        acme_clients: 200,
        clients: 2451,
      },
      mounts: [
        {
          mount_path: 'auth/authid0',
          counts: {
            clients: 1976,
            entity_clients: 783,
            non_entity_clients: 1193,
            distinct_entities: 783,
            non_entity_tokens: 1193,
            secret_syncs: 0,
            acme_clients: 0,
          },
        },
        {
          mount_path: 'kvv2-engine-0',
          counts: {
            clients: 275,
            entity_clients: 0,
            non_entity_clients: 0,
            distinct_entities: 0,
            non_entity_tokens: 0,
            secret_syncs: 275,
            acme_clients: 0,
          },
        },
        {
          mount_path: 'pki-engine-0',
          counts: {
            clients: 200,
            entity_clients: 0,
            non_entity_clients: 0,
            distinct_entities: 0,
            non_entity_tokens: 0,
            secret_syncs: 0,
            acme_clients: 200,
          },
        },
      ],
    },
  ],
  months: [
    {
      timestamp: '2023-08-01T00:00:00Z',
      counts: null,
      namespaces: null,
      new_clients: null,
    },
    {
      timestamp: '2023-09-01T00:00:00Z',
      counts: {
        distinct_entities: 1329,
        entity_clients: 1329,
        non_entity_tokens: 1738,
        non_entity_clients: 1738,
        secret_syncs: 5525,
        acme_clients: 200,
        clients: 8792,
      },
      namespaces: [
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 1279,
            entity_clients: 1279,
            non_entity_tokens: 1598,
            non_entity_clients: 1598,
            secret_syncs: 2755,
            acme_clients: 75,
            clients: 5707,
          },
          mounts: [
            {
              mount_path: 'auth/authid0',
              counts: {
                clients: 2877,
                entity_clients: 1279,
                non_entity_clients: 1598,
                distinct_entities: 1279,
                non_entity_tokens: 1598,
                secret_syncs: 0,
                acme_clients: 0,
              },
            },
            {
              mount_path: 'kvv2-engine-0',
              counts: {
                clients: 2755,
                entity_clients: 0,
                non_entity_clients: 0,
                distinct_entities: 0,
                non_entity_tokens: 0,
                secret_syncs: 2755,
                acme_clients: 0,
              },
            },
            {
              mount_path: 'pki-engine-0',
              counts: {
                clients: 75,
                entity_clients: 0,
                non_entity_clients: 0,
                distinct_entities: 0,
                non_entity_tokens: 0,
                secret_syncs: 0,
                acme_clients: 75,
              },
            },
          ],
        },
        {
          namespace_id: '81ry61',
          namespace_path: 'ns1',
          counts: {
            distinct_entities: 50,
            entity_clients: 50,
            non_entity_tokens: 140,
            non_entity_clients: 140,
            secret_syncs: 2770,
            acme_clients: 125,
            clients: 3085,
          },
          mounts: [
            {
              mount_path: 'kvv2-engine-0',
              counts: {
                clients: 2770,
                entity_clients: 0,
                non_entity_clients: 0,
                distinct_entities: 0,
                non_entity_tokens: 0,
                secret_syncs: 2770,
                acme_clients: 0,
              },
            },
            {
              mount_path: 'auth/authid0',
              counts: {
                clients: 190,
                entity_clients: 50,
                non_entity_clients: 140,
                distinct_entities: 50,
                non_entity_tokens: 140,
                secret_syncs: 0,
                acme_clients: 0,
              },
            },
            {
              mount_path: 'pki-engine-0',
              counts: {
                clients: 125,
                entity_clients: 0,
                non_entity_clients: 0,
                distinct_entities: 0,
                non_entity_tokens: 0,
                secret_syncs: 0,
                acme_clients: 125,
              },
            },
          ],
        },
      ],
      new_clients: {
        counts: {
          distinct_entities: 39,
          entity_clients: 39,
          non_entity_tokens: 81,
          non_entity_clients: 81,
          secret_syncs: 166,
          acme_clients: 50,
          clients: 336,
        },
        namespaces: [
          {
            namespace_id: '81ry61',
            namespace_path: 'ns1',
            counts: {
              distinct_entities: 30,
              entity_clients: 30,
              non_entity_tokens: 62,
              non_entity_clients: 62,
              secret_syncs: 100,
              acme_clients: 30,
              clients: 222,
            },
            mounts: [
              {
                mount_path: 'kvv2-engine-0',
                counts: {
                  clients: 100,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  distinct_entities: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 100,
                  acme_clients: 0,
                },
              },
              {
                mount_path: 'auth/authid0',
                counts: {
                  clients: 92,
                  entity_clients: 30,
                  non_entity_clients: 62,
                  distinct_entities: 30,
                  non_entity_tokens: 62,
                  secret_syncs: 0,
                  acme_clients: 0,
                },
              },
              {
                mount_path: 'pki-engine-0',
                counts: {
                  clients: 30,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  distinct_entities: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 0,
                  acme_clients: 30,
                },
              },
            ],
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 9,
              entity_clients: 9,
              non_entity_tokens: 19,
              non_entity_clients: 19,
              secret_syncs: 66,
              acme_clients: 20,
              clients: 114,
            },
            mounts: [
              {
                mount_path: 'kvv2-engine-0',
                counts: {
                  clients: 66,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  distinct_entities: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 66,
                  acme_clients: 0,
                },
              },
              {
                mount_path: 'auth/authid0',
                counts: {
                  clients: 28,
                  entity_clients: 9,
                  non_entity_clients: 19,
                  distinct_entities: 9,
                  non_entity_tokens: 19,
                  secret_syncs: 0,
                  acme_clients: 0,
                },
              },
              {
                mount_path: 'pki-engine-0',
                counts: {
                  clients: 20,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  distinct_entities: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 0,
                  acme_clients: 20,
                },
              },
            ],
          },
        ],
      },
    },
  ],
  total: {
    distinct_entities: 1816,
    entity_clients: 1816,
    non_entity_tokens: 3117,
    non_entity_clients: 3117,
    secret_syncs: 2672,
    acme_clients: 200,
    clients: 7805,
  },
};

// combined activity data before and after 1.10 upgrade when Vault added mount attribution
export const MIXED_ACTIVITY_RESPONSE_STUB = {
  start_time: '2024-03-01T00:00:00Z',
  end_time: '2024-04-30T23:59:59Z',
  total: {
    acme_clients: 0,
    clients: 3,
    distinct_entities: 3,
    entity_clients: 3,
    non_entity_clients: 0,
    non_entity_tokens: 0,
    secret_syncs: 0,
  },
  by_namespace: [
    {
      counts: {
        acme_clients: 0,
        clients: 3,
        distinct_entities: 3,
        entity_clients: 3,
        non_entity_clients: 0,
        non_entity_tokens: 0,
        secret_syncs: 0,
      },
      mounts: [
        {
          counts: {
            acme_clients: 0,
            clients: 2,
            distinct_entities: 2,
            entity_clients: 2,
            non_entity_clients: 0,
            non_entity_tokens: 0,
            secret_syncs: 0,
          },
          mount_path: 'no mount accessor (pre-1.10 upgrade?)',
        },
        {
          counts: {
            acme_clients: 0,
            clients: 1,
            distinct_entities: 1,
            entity_clients: 1,
            non_entity_clients: 0,
            non_entity_tokens: 0,
            secret_syncs: 0,
          },
          mount_path: 'auth/u/',
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
        distinct_entities: 0,
        entity_clients: 3,
        non_entity_clients: 0,
        non_entity_tokens: 0,
        secret_syncs: 0,
      },
      namespaces: [
        {
          counts: {
            acme_clients: 0,
            clients: 3,
            distinct_entities: 0,
            entity_clients: 3,
            non_entity_clients: 0,
            non_entity_tokens: 0,
            secret_syncs: 0,
          },
          mounts: [
            {
              counts: {
                acme_clients: 0,
                clients: 2,
                distinct_entities: 0,
                entity_clients: 2,
                non_entity_clients: 0,
                non_entity_tokens: 0,
                secret_syncs: 0,
              },
              mount_path: 'no mount accessor (pre-1.10 upgrade?)',
            },
            {
              counts: {
                acme_clients: 0,
                clients: 1,
                distinct_entities: 0,
                entity_clients: 1,
                non_entity_clients: 0,
                non_entity_tokens: 0,
                secret_syncs: 0,
              },
              mount_path: 'auth/u/',
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
          distinct_entities: 0,
          entity_clients: 3,
          non_entity_clients: 0,
          non_entity_tokens: 0,
          secret_syncs: 0,
        },
        namespaces: [
          {
            counts: {
              acme_clients: 0,
              clients: 3,
              distinct_entities: 0,
              entity_clients: 3,
              non_entity_clients: 0,
              non_entity_tokens: 0,
              secret_syncs: 0,
            },
            mounts: [
              {
                counts: {
                  acme_clients: 0,
                  clients: 2,
                  distinct_entities: 0,
                  entity_clients: 2,
                  non_entity_clients: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 0,
                },
                mount_path: 'no mount accessor (pre-1.10 upgrade?)',
              },
              {
                counts: {
                  acme_clients: 0,
                  clients: 1,
                  distinct_entities: 0,
                  entity_clients: 1,
                  non_entity_clients: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 0,
                },
                mount_path: 'auth/u/',
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
// format returned by model hook in routes/vault/cluster/clients.ts
export const VERSION_HISTORY = [
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

// order of this array matters because index 0 is a month without data
export const SERIALIZED_ACTIVITY_RESPONSE = {
  by_namespace: [
    {
      label: 'root',
      clients: 5429,
      entity_clients: 1033,
      non_entity_clients: 1924,
      secret_syncs: 2397,
      acme_clients: 75,
      mounts: [
        {
          acme_clients: 0,
          clients: 2957,
          entity_clients: 1033,
          label: 'auth/authid0',
          non_entity_clients: 1924,
          secret_syncs: 0,
        },
        {
          acme_clients: 0,
          clients: 2397,
          entity_clients: 0,
          label: 'kvv2-engine-0',
          non_entity_clients: 0,
          secret_syncs: 2397,
        },
        {
          acme_clients: 75,
          clients: 75,
          entity_clients: 0,
          label: 'pki-engine-0',
          non_entity_clients: 0,
          secret_syncs: 0,
        },
      ],
    },
    {
      label: 'ns1',
      clients: 2451,
      entity_clients: 783,
      non_entity_clients: 1193,
      secret_syncs: 275,
      acme_clients: 200,
      mounts: [
        {
          label: 'auth/authid0',
          clients: 1976,
          entity_clients: 783,
          non_entity_clients: 1193,
          secret_syncs: 0,
          acme_clients: 0,
        },
        {
          label: 'kvv2-engine-0',
          clients: 275,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 275,
          acme_clients: 0,
        },
        {
          label: 'pki-engine-0',
          acme_clients: 200,
          clients: 200,
          entity_clients: 0,
          non_entity_clients: 0,
          secret_syncs: 0,
        },
      ],
    },
  ],
  by_month: [
    {
      month: '8/23',
      timestamp: '2023-08-01T00:00:00Z',
      namespaces: [],
      new_clients: {
        month: '8/23',
        timestamp: '2023-08-01T00:00:00Z',
        namespaces: [],
      },
      namespaces_by_key: {},
    },
    {
      month: '9/23',
      timestamp: '2023-09-01T00:00:00Z',
      clients: 8592,
      entity_clients: 1329,
      non_entity_clients: 1738,
      secret_syncs: 5525,
      acme_clients: 200,
      namespaces: [
        {
          label: 'root',
          clients: 5707,
          entity_clients: 1279,
          non_entity_clients: 1598,
          secret_syncs: 2755,
          acme_clients: 75,
          mounts: [
            {
              label: 'auth/authid0',
              clients: 2877,
              entity_clients: 1279,
              non_entity_clients: 1598,
              secret_syncs: 0,
              acme_clients: 0,
            },
            {
              label: 'kvv2-engine-0',
              clients: 2755,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 2755,
              acme_clients: 0,
            },
            {
              label: 'pki-engine-0',
              acme_clients: 75,
              clients: 75,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
          ],
        },
        {
          label: 'ns1',
          clients: 3085,
          entity_clients: 50,
          non_entity_clients: 140,
          secret_syncs: 2770,
          acme_clients: 125,
          mounts: [
            {
              label: 'kvv2-engine-0',
              clients: 2770,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 2770,
              acme_clients: 0,
            },
            {
              label: 'auth/authid0',
              clients: 190,
              entity_clients: 50,
              non_entity_clients: 140,
              secret_syncs: 0,
              acme_clients: 0,
            },
            {
              label: 'pki-engine-0',
              acme_clients: 125,
              clients: 125,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
            },
          ],
        },
      ],
      namespaces_by_key: {
        root: {
          month: '9/23',
          timestamp: '2023-09-01T00:00:00Z',
          clients: 5707,
          entity_clients: 1279,
          non_entity_clients: 1598,
          secret_syncs: 2755,
          acme_clients: 75,
          new_clients: {
            month: '9/23',
            timestamp: '2023-09-01T00:00:00Z',
            label: 'root',
            clients: 114,
            entity_clients: 9,
            non_entity_clients: 19,
            secret_syncs: 66,
            acme_clients: 20,
            mounts: [
              {
                label: 'kvv2-engine-0',
                clients: 66,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 66,
                acme_clients: 0,
              },
              {
                label: 'auth/authid0',
                clients: 28,
                entity_clients: 9,
                non_entity_clients: 19,
                secret_syncs: 0,
                acme_clients: 0,
              },
              {
                label: 'pki-engine-0',
                clients: 20,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
                acme_clients: 20,
              },
            ],
          },
          mounts_by_key: {
            'auth/authid0': {
              month: '9/23',
              timestamp: '2023-09-01T00:00:00Z',
              label: 'auth/authid0',
              clients: 2877,
              entity_clients: 1279,
              non_entity_clients: 1598,
              secret_syncs: 0,
              acme_clients: 0,
              new_clients: {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00Z',
                label: 'auth/authid0',
                clients: 28,
                entity_clients: 9,
                non_entity_clients: 19,
                secret_syncs: 0,
                acme_clients: 0,
              },
            },
            'kvv2-engine-0': {
              month: '9/23',
              timestamp: '2023-09-01T00:00:00Z',
              label: 'kvv2-engine-0',
              clients: 2755,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 2755,
              acme_clients: 0,
              new_clients: {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00Z',
                label: 'kvv2-engine-0',
                clients: 66,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 66,
                acme_clients: 0,
              },
            },
            'pki-engine-0': {
              month: '9/23',
              timestamp: '2023-09-01T00:00:00Z',
              label: 'pki-engine-0',
              clients: 75,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 0,
              acme_clients: 75,
              new_clients: {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00Z',
                label: 'pki-engine-0',
                acme_clients: 20,
                clients: 20,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
          },
        },
        ns1: {
          month: '9/23',
          timestamp: '2023-09-01T00:00:00Z',
          clients: 3085,
          entity_clients: 50,
          non_entity_clients: 140,
          secret_syncs: 2770,
          acme_clients: 125,
          new_clients: {
            month: '9/23',
            timestamp: '2023-09-01T00:00:00Z',
            label: 'ns1',
            clients: 222,
            entity_clients: 30,
            non_entity_clients: 62,
            secret_syncs: 100,
            acme_clients: 30,
            mounts: [
              {
                label: 'kvv2-engine-0',
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 100,
                acme_clients: 0,
              },
              {
                label: 'auth/authid0',
                clients: 92,
                entity_clients: 30,
                non_entity_clients: 62,
                secret_syncs: 0,
                acme_clients: 0,
              },
              {
                label: 'pki-engine-0',
                acme_clients: 30,
                clients: 30,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            ],
          },
          mounts_by_key: {
            'kvv2-engine-0': {
              month: '9/23',
              timestamp: '2023-09-01T00:00:00Z',
              label: 'kvv2-engine-0',
              clients: 2770,
              entity_clients: 0,
              non_entity_clients: 0,
              secret_syncs: 2770,
              acme_clients: 0,
              new_clients: {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00Z',
                label: 'kvv2-engine-0',
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 100,
                acme_clients: 0,
              },
            },
            'auth/authid0': {
              month: '9/23',
              timestamp: '2023-09-01T00:00:00Z',
              label: 'auth/authid0',
              clients: 190,
              entity_clients: 50,
              non_entity_clients: 140,
              secret_syncs: 0,
              acme_clients: 0,
              new_clients: {
                month: '9/23',
                timestamp: '2023-09-01T00:00:00Z',
                label: 'auth/authid0',
                clients: 92,
                entity_clients: 30,
                non_entity_clients: 62,
                secret_syncs: 0,
                acme_clients: 0,
              },
            },
            'pki-engine-0': {
              month: '9/23',
              timestamp: '2023-09-01T00:00:00Z',
              clients: 125,
              acme_clients: 125,
              entity_clients: 0,
              label: 'pki-engine-0',
              non_entity_clients: 0,
              secret_syncs: 0,
              new_clients: {
                acme_clients: 30,
                clients: 30,
                entity_clients: 0,
                label: 'pki-engine-0',
                month: '9/23',
                timestamp: '2023-09-01T00:00:00Z',
                non_entity_clients: 0,
                secret_syncs: 0,
              },
            },
          },
        },
      },
      new_clients: {
        month: '9/23',
        timestamp: '2023-09-01T00:00:00Z',
        clients: 336,
        entity_clients: 39,
        non_entity_clients: 81,
        secret_syncs: 166,
        acme_clients: 50,
        namespaces: [
          {
            label: 'ns1',
            clients: 222,
            entity_clients: 30,
            non_entity_clients: 62,
            secret_syncs: 100,
            acme_clients: 30,
            mounts: [
              {
                label: 'kvv2-engine-0',
                clients: 100,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 100,
                acme_clients: 0,
              },
              {
                label: 'auth/authid0',
                clients: 92,
                entity_clients: 30,
                non_entity_clients: 62,
                secret_syncs: 0,
                acme_clients: 0,
              },
              {
                label: 'pki-engine-0',
                clients: 30,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
                acme_clients: 30,
              },
            ],
          },
          {
            label: 'root',
            clients: 114,
            entity_clients: 9,
            non_entity_clients: 19,
            secret_syncs: 66,
            acme_clients: 20,
            mounts: [
              {
                label: 'kvv2-engine-0',
                clients: 66,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 66,
                acme_clients: 0,
              },
              {
                label: 'auth/authid0',
                clients: 28,
                entity_clients: 9,
                non_entity_clients: 19,
                secret_syncs: 0,
                acme_clients: 0,
              },
              {
                label: 'pki-engine-0',
                clients: 20,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 0,
                acme_clients: 20,
              },
            ],
          },
        ],
      },
    },
  ],
};
