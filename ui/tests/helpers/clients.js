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
  return new Response(200, { 'Content-Type': 'application/json' }, JSON.stringify(data));
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
