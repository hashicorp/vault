/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// TODO: separate nested into distinct exported consts
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
  statText: (label: string) => `[data-test-stat-text="${label}"]`,
  chartContainer: (title: string) => `[data-test-chart-container="${title}"]`,
  charts: {
    chart: (title: string) => `[data-test-chart="${title}"]`, // newer lineal charts
    statTextValue: (label: string) =>
      label ? `[data-test-stat-text="${label}"] .stat-value` : '[data-test-stat-text]',
    legend: '[data-test-chart-container-legend]',
    legendLabel: (nth: number) => `.legend-label:nth-child(${nth * 2})`, // nth * 2 accounts for dots in legend
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
  usageStats: (title: string) => `[data-test-usage-stats="${title}"]`,
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
    selectMonth: (month: string) => `[data-test-dropdown-month="${month}"]`,
    selectYear: (year: string) => `[data-test-dropdown-year="${year}"]`,
    submit: '[data-test-date-dropdown-submit]',
  },
  calendarWidget: {
    trigger: '[data-test-calendar-widget-trigger]',
    currentMonth: '[data-test-current-month]',
    currentBillingPeriod: '[data-test-current-billing-period]',
    customEndMonth: '[data-test-show-calendar]',
    previousYear: '[data-test-previous-year]',
    nextYear: '[data-test-next-year]',
    calendarMonth: (month: string) => `[data-test-calendar-month="${month}"]`,
  },
  selectedAuthMount: 'div#mounts-search-select [data-test-selected-option] div',
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
