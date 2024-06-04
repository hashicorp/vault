/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';
import { SELECTORS as GENERAL } from 'vault/tests/helpers/general-selectors';
import { click } from '@ember/test-helpers';

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
export const SELECTORS = {
  ...GENERAL,
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
      hoverCircle: (month) => `[data-test-hover-circle="${month}"]`,
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
  dateDropdownSubmit: '[data-test-date-dropdown-submit]',
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
  return new Response(httpStatus, { 'Content-Type': 'application/json' }, JSON.stringify(data));
}

export async function dateDropdownSelect(month, year) {
  const { dateDropdown, counts } = SELECTORS;
  await click(counts.startEdit);
  await click(dateDropdown.toggleMonth);
  await click(dateDropdown.selectMonth(month));
  await click(dateDropdown.toggleYear);
  await click(dateDropdown.selectYear(year));
  await click(dateDropdown.submit);
}

export const ACTIVITY_RESPONSE_STUB = {
  start_time: '2023-08-01T00:00:00Z',
  end_time: '2023-10-31T23:59:59Z', // is always the last day and hour of the month queried
  by_namespace: [
    {
      namespace_id: 'e67m31',
      namespace_path: 'ns1',
      counts: {
        clients: 13204,
        entity_clients: 4256,
        non_entity_clients: 4138,
        secret_syncs: 4810,
        distinct_entities: 4256,
        non_entity_tokens: 4138,
      },
      mounts: [
        {
          mount_path: 'auth/authid/0',
          counts: {
            clients: 8394,
            entity_clients: 4256,
            non_entity_clients: 4138,
            secret_syncs: 0,
            distinct_entities: 4256,
            non_entity_tokens: 4138,
          },
        },
        {
          mount_path: 'kvv2-engine-0',
          counts: {
            clients: 4810,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 4810,
            distinct_entities: 0,
            non_entity_tokens: 0,
          },
        },
      ],
    },
    {
      namespace_id: 'root',
      namespace_path: '',
      counts: {
        clients: 12381,
        entity_clients: 4002,
        non_entity_clients: 4089,
        secret_syncs: 4290,
        distinct_entities: 4002,
        non_entity_tokens: 4089,
      },
      mounts: [
        {
          mount_path: 'auth/authid/0',
          counts: {
            clients: 8091,
            entity_clients: 4002,
            non_entity_clients: 4089,
            secret_syncs: 0,
            distinct_entities: 4002,
            non_entity_tokens: 4089,
          },
        },
        {
          mount_path: 'kvv2-engine-0',
          counts: {
            clients: 4290,
            entity_clients: 0,
            non_entity_clients: 0,
            secret_syncs: 4290,
            distinct_entities: 0,
            non_entity_tokens: 0,
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
        clients: 3928,
        entity_clients: 832,
        non_entity_clients: 930,
        secret_syncs: 238,
        distinct_entities: 832,
        non_entity_tokens: 930,
      },
      namespaces: [
        {
          namespace_id: 'e67m31',
          namespace_path: 'ns1',
          counts: {
            clients: 1047,
            entity_clients: 708,
            non_entity_clients: 182,
            secret_syncs: 157,
            distinct_entities: 708,
            non_entity_tokens: 182,
          },
          mounts: [
            {
              mount_path: 'auth/authid/0',
              counts: {
                clients: 890,
                entity_clients: 708,
                non_entity_clients: 182,
                secret_syncs: 0,
                distinct_entities: 0,
                non_entity_tokens: 0,
              },
            },
            {
              mount_path: 'kvv2-engine-0',
              counts: {
                clients: 157,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 157,
                distinct_entities: 0,
                non_entity_tokens: 0,
              },
            },
          ],
        },
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            clients: 1947,
            entity_clients: 124,
            non_entity_clients: 748,
            secret_syncs: 81,
            distinct_entities: 124,
            non_entity_tokens: 748,
          },
          mounts: [
            {
              mount_path: 'auth/authid/0',
              counts: {
                clients: 872,
                entity_clients: 124,
                non_entity_clients: 748,
                secret_syncs: 0,
                distinct_entities: 0,
                non_entity_tokens: 0,
              },
            },
            {
              mount_path: 'kvv2-engine-0',
              counts: {
                clients: 81,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 81,
                distinct_entities: 0,
                non_entity_tokens: 0,
              },
            },
          ],
        },
      ],
      new_clients: {
        counts: {
          clients: 364,
          entity_clients: 59,
          non_entity_clients: 112,
          secret_syncs: 49,
          distinct_entities: 59,
          non_entity_tokens: 112,
        },
        namespaces: [
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              clients: 191,
              entity_clients: 25,
              non_entity_clients: 50,
              secret_syncs: 25,
              distinct_entities: 25,
              non_entity_tokens: 50,
            },
            mounts: [
              {
                mount_path: 'auth/authid/0',
                counts: {
                  clients: 75,
                  entity_clients: 25,
                  non_entity_clients: 50,
                  secret_syncs: 0,
                  distinct_entities: 0,
                  non_entity_tokens: 0,
                },
              },
              {
                mount_path: 'kvv2-engine-0',
                counts: {
                  clients: 25,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 25,
                  distinct_entities: 0,
                  non_entity_tokens: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'e67m31',
            namespace_path: 'ns1',
            counts: {
              clients: 173,
              entity_clients: 34,
              non_entity_clients: 62,
              secret_syncs: 24,
              distinct_entities: 34,
              non_entity_tokens: 62,
            },
            mounts: [
              {
                mount_path: 'auth/authid/0',
                counts: {
                  clients: 96,
                  entity_clients: 34,
                  non_entity_clients: 62,
                  secret_syncs: 0,
                  distinct_entities: 0,
                  non_entity_tokens: 0,
                },
              },
              {
                mount_path: 'kvv2-engine-0',
                counts: {
                  clients: 24,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  secret_syncs: 24,
                  distinct_entities: 0,
                  non_entity_tokens: 0,
                },
              },
            ],
          },
        ],
      },
    },
    {
      timestamp: '2023-10-01T00:00:00Z',
      counts: {
        clients: 3928,
        entity_clients: 832,
        non_entity_clients: 930,
        secret_syncs: 238,
        distinct_entities: 832,
        non_entity_tokens: 930,
      },
      namespaces: [
        {
          namespace_id: 'e67m31',
          namespace_path: 'ns1',
          counts: {
            clients: 1981,
            entity_clients: 708,
            non_entity_clients: 182,
            secret_syncs: 157,
            distinct_entities: 708,
            non_entity_tokens: 182,
          },
          mounts: [
            {
              mount_path: 'auth/authid/0',
              counts: {
                clients: 890,
                entity_clients: 708,
                non_entity_clients: 182,
                secret_syncs: 0,
                distinct_entities: 0,
                non_entity_tokens: 0,
              },
            },
            {
              mount_path: 'kvv2-engine-0',
              counts: {
                clients: 157,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 157,
                distinct_entities: 0,
                non_entity_tokens: 0,
              },
            },
          ],
        },
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            clients: 1947,
            entity_clients: 124,
            non_entity_clients: 748,
            secret_syncs: 81,
            distinct_entities: 124,
            non_entity_tokens: 748,
          },
          mounts: [
            {
              mount_path: 'auth/authid/0',
              counts: {
                clients: 872,
                entity_clients: 124,
                non_entity_clients: 748,
                secret_syncs: 0,
                distinct_entities: 0,
                non_entity_tokens: 0,
              },
            },
            {
              mount_path: 'kvv2-engine-0',
              counts: {
                clients: 81,
                entity_clients: 0,
                non_entity_clients: 0,
                secret_syncs: 81,
                distinct_entities: 0,
                non_entity_tokens: 0,
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
  ],
  total: {
    clients: 35287,
    entity_clients: 8258,
    non_entity_clients: 8227,
    secret_syncs: 9100,
    distinct_entities: 8258,
    non_entity_tokens: 8227,
  },
};
