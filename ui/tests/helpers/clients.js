/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';
import { SELECTORS as GENERAL } from 'vault/tests/helpers/general-selectors';
import { click } from '@ember/test-helpers';
import { addMonths, startOfMonth, subMonths } from 'date-fns';

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
    authMounts: '[data-test-counts-auth-mounts]',
    startDiscrepancy: '[data-test-counts-start-discrepancy]',
  },
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
  monthlyNewChart: '[data-test-chart="monthly new"]',
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

export async function dateDropdownSelect(month, year) {
  const { dateDropdown, counts } = SELECTORS;
  await click(counts.startEdit);
  await click(dateDropdown.toggleMonth);
  await click(dateDropdown.selectMonth(month));
  await click(dateDropdown.toggleYear);
  await click(dateDropdown.selectYear(year));
  await click(dateDropdown.submit);
}

export const STATIC_NOW = new Date('2023-01-13T14:15:00');
// for testing, we're in the middle of a license/billing period
export const LICENSE_START = startOfMonth(subMonths(STATIC_NOW, 6)); // 2022-07-01
// upgrade happened 1 month after license start
export const UPGRADE_DATE = addMonths(LICENSE_START, 1); // 2022-08-01
