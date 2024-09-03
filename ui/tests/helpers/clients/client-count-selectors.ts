/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// TODO: separate nested into distinct exported consts
export const CLIENT_COUNT = {
  counts: {
    description: '[data-test-counts-description]',
    configDisabled: '[data-test-counts-disabled]',
    namespaces: '[data-test-counts-namespaces]',
    mountPaths: '[data-test-counts-auth-mounts]',
    startDiscrepancy: '[data-test-counts-start-discrepancy]',
  },
  dateRange: {
    dateDisplay: (name: string) => (name ? `[data-test-date-range="${name}"]` : '[data-test-date-range]'),
    edit: '[data-test-date-range-edit]',
    editModal: '[data-test-date-range-edit-modal]',
    editDate: (name: string) => `[data-test-date-edit="${name}"]`,
    reset: '[data-test-date-edit="reset"]',
    defaultRangeAlert: '[data-test-range-default-alert]',
    validation: '[data-test-date-range-validation]',
  },
  statText: (label: string) => `[data-test-stat-text="${label}"]`,
  statTextValue: (label: string) =>
    label ? `[data-test-stat-text="${label}"] .stat-value` : '[data-test-stat-text]',
  usageStats: (title: string) => `[data-test-usage-stats="${title}"]`,
  attributionBlock: (type: string) =>
    type ? `[data-test-clients-attribution="${type}"]` : '[data-test-clients-attribution]',
  filterBar: '[data-test-clients-filter-bar]',
  nsFilter: '#namespace-search-select',
  mountFilter: '#mounts-search-select',
  selectedAuthMount: 'div#mounts-search-select [data-test-selected-option] div',
  selectedNs: 'div#namespace-search-select [data-test-selected-option] div',
  upgradeWarning: '[data-test-clients-upgrade-warning]',
  exportButton: '[data-test-export-button]',
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
