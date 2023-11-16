/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

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
  dashboardActiveTab: '.active[data-test-dashboard]',
  emptyStateTitle: '[data-test-empty-state-title]',
  usageStats: '[data-test-usage-stats]',
  dateDisplay: '[data-test-date-display]',
  attributionBlock: '[data-test-clients-attribution]',
  filterBar: '[data-test-clients-filter-bar]',
  rangeDropdown: '[data-test-calendar-widget-trigger]',
  monthDropdown: '[data-test-toggle-month]',
  yearDropdown: '[data-test-toggle-year]',
  dateDropdownSubmit: '[data-test-date-dropdown-submit]',
  runningTotalMonthStats: '[data-test-running-total="single-month-stats"]',
  runningTotalMonthlyCharts: '[data-test-running-total="monthly-charts"]',
  monthlyUsageBlock: '[data-test-monthly-usage]',
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
