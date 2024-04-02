/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';
import { click } from '@ember/test-helpers';

import { LICENSE_START } from 'vault/mirage/handlers/clients';
import { addMonths } from 'date-fns';

/** Scenarios
  Config off, no data
  Config on, no data
  Config on, with data
  Filtering (data with mounts)
  Filtering (data without mounts)
  Filtering (data without mounts)
  * -- HISTORY ONLY --
  Filtering different date ranges (hist only)
  Upgrade warning
  No permissions for license
  Version
  queries available
  queries unavailable
  License start date this month
*/
export const CLIENT_COUNT = {
  counts: {
    startLabel: '[data-test-counts-start-label]',
    description: '[data-test-counts-description]',
    startMonth: '[data-test-counts-start-month]',
    startEdit: '[data-test-counts-start-edit]',
    startDropdown: '[data-test-counts-start-dropdown]',
    configDisabled: '[data-test-counts-disabled]',
    namespaces: '[data-test-counts-namespaces]',
    mountPaths: '[data-test-counts-auth-mounts]',
    startDiscrepancy: '[data-test-counts-start-discrepancy]',
  },
  tokenTab: {
    entity: '[data-test-monthly-new-entity]',
    nonentity: '[data-test-monthly-new-nonentity]',
    legend: '[data-test-monthly-new-legend]',
  },
  syncTab: {
    total: '[data-test-total-sync-clients]',
    average: '[data-test-average-sync-clients]',
  },
  charts: {
    chart: (title) => `[data-test-chart="${title}"]`, // newer lineal charts
    statTextValue: (label) =>
      label ? `[data-test-stat-text-container="${label}"] .stat-value` : '[data-test-stat-text-container]',
    legend: '[data-test-chart-container-legend]',
    legendLabel: (nth) => `.legend-label:nth-child(${nth * 2})`, // nth * 2 accounts for dots in legend
    timestamp: '[data-test-chart-container-timestamp]',
    dataBar: '[data-test-vertical-bar]',
    xAxisLabel: '[data-test-x-axis] text',
    // selectors for old d3 charts
    verticalBar: '[data-test-vertical-bar-chart]',
    lineChart: '[data-test-line-chart]',
    bar: {
      xAxisLabel: '[data-test-vertical-chart="x-axis-labels"] text',
      dataBar: '[data-test-vertical-chart="data-bar"]',
    },
    line: {
      xAxisLabel: '[data-test-line-chart] [data-test-x-axis] text',
      plotPoint: '[data-test-line-chart="plot-point"]',
    },
  },
  usageStats: '[data-test-usage-stats]',
  dateDisplay: '[data-test-date-display]',
  attributionBlock: '[data-test-clients-attribution]',
  filterBar: '[data-test-clients-filter-bar]',
  rangeDropdown: '[data-test-calendar-widget-trigger]',
  monthDropdown: '[data-test-toggle-month]',
  yearDropdown: '[data-test-toggle-year]',
  currentBillingPeriod: '[data-test-current-billing-period]',
  dateDropdown: {
    toggleMonth: '[data-test-toggle-month]',
    toggleYear: '[data-test-toggle-year]',
    selectMonth: (month) => `[data-test-dropdown-month="${month}"]`,
    selectYear: (year) => `[data-test-dropdown-year="${year}"]`,
    submit: '[data-test-date-dropdown-submit]',
  },
  calendarWidget: {
    trigger: '[data-test-calendar-widget-trigger]',
    currentMonth: '[data-test-current-month]',
    currentBillingPeriod: '[data-test-current-billing-period]',
    customEndMonth: '[data-test-show-calendar]',
    previousYear: '[data-test-previous-year]',
    nextYear: '[data-test-next-year]',
    calendarMonth: (month) => `[data-test-calendar-month="${month}"]`,
  },
  runningTotalMonthStats: '[data-test-running-total="single-month-stats"]',
  runningTotalMonthlyCharts: '[data-test-running-total="monthly-charts"]',
  selectedAuthMount: 'div#auth-method-search-select [data-test-selected-option] div',
  selectedNs: 'div#namespace-search-select [data-test-selected-option] div',
  upgradeWarning: '[data-test-clients-upgrade-warning]',
};

export const CHART_ELEMENTS = {
  entityClientDataBars: '[data-test-group="entity_clients"]',
  nonEntityDataBars: '[data-test-group="non_entity_clients"]',
  yLabels: '[data-test-group="y-labels"]',
  actionBars: '[data-test-group="action-bars"]',
  labelActionBars: '[data-test-group="label-action-bars"]',
  totalValues: '[data-test-group="total-values"]',
};

export function sendResponse(data, httpStatus = 200) {
  if (httpStatus === 403) {
    return [
      httpStatus,
      { 'Content-Type': 'application/json' },
      JSON.stringify({ errors: ['permission denied'] }),
    ];
  }
  if (httpStatus === 204) {
    // /activity endpoint returns 204 when no data, while
    // /activity/monthly returns 200 with zero values on data
    return [httpStatus, { 'Content-Type': 'application/json' }];
  }
  return [httpStatus, { 'Content-Type': 'application/json' }, JSON.stringify(data)];
}

export function overrideResponse(httpStatus, data) {
  if (httpStatus === 403) {
    return new Response(
      403,
      { 'Content-Type': 'application/json' },
      JSON.stringify({ errors: ['permission denied'] })
    );
  }
  // /activity endpoint returns 204 when no data, while
  // /activity/monthly returns 200 with zero values on data
  if (httpStatus === 204) {
    return new Response(204, { 'Content-Type': 'application/json' });
  }
  return new Response(200, { 'Content-Type': 'application/json' }, JSON.stringify(data));
}

export async function dateDropdownSelect(month, year) {
  const { dateDropdown, counts } = CLIENT_COUNT;
  await click(counts.startEdit);
  await click(dateDropdown.toggleMonth);
  await click(dateDropdown.selectMonth(month));
  await click(dateDropdown.toggleYear);
  await click(dateDropdown.selectYear(year));
  await click(dateDropdown.submit);
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
      namespace_path: 'ns/1',
      counts: {
        distinct_entities: 783,
        entity_clients: 783,
        non_entity_tokens: 1193,
        non_entity_clients: 1193,
        secret_syncs: 275,
        acme_clients: 125,
        clients: 2376,
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
          namespace_path: 'ns/1',
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
            namespace_path: 'ns/1',
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
      label: 'ns/1',
      clients: 2376,
      entity_clients: 783,
      non_entity_clients: 1193,
      secret_syncs: 275,
      acme_clients: 125,
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
          acme_clients: 125,
          clients: 125,
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
          label: 'ns/1',
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
        'ns/1': {
          month: '9/23',
          timestamp: '2023-09-01T00:00:00Z',
          clients: 3085,
          entity_clients: 50,
          non_entity_clients: 140,
          secret_syncs: 2770,
          acme_clients: 125,
          new_clients: {
            month: '9/23',
            label: 'ns/1',
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
            label: 'ns/1',
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
