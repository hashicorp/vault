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
  statTextValue: (label: string) =>
    label ? `[data-test-stat-text="${label}"] .stat-value` : '[data-test-stat-text]',
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
    displayYear: '[data-test-display-year]',
    calendarMonth: (month: string) => `[data-test-calendar-month="${month}"]`,
  },
  selectedAuthMount: 'div#mounts-search-select [data-test-selected-option] div',
  selectedNs: 'div#namespace-search-select [data-test-selected-option] div',
  upgradeWarning: '[data-test-clients-upgrade-warning]',
};

export const CHARTS = {
  // container
  container: (title: string) => `[data-test-chart-container="${title}"]`,
  timestamp: '[data-test-chart-container-timestamp]',
  legend: '[data-test-chart-container-legend]',
  legendLabel: (nth: number) => `.legend-label:nth-child(${nth * 2})`, // nth * 2 accounts for dots in legend

  // chart elements
  chart: (title: string) => `[data-test-chart="${title}"]`,
  hover: (area: string) => `[data-test-interactive-area="${area}"]`,
  table: '[data-test-underlying-data]',
  tooltip: '[data-test-tooltip]',
  verticalBar: '[data-test-vertical-bar]',
  xAxis: '[data-test-x-axis]',
  yAxis: '[data-test-y-axis]',
  xAxisLabel: '[data-test-x-axis] text',
  plotPoint: '[data-test-plot-point]',
};
