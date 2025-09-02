/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// TODO: separate nested into distinct exported consts
export const CLIENT_COUNT = {
  card: (name: string) => `[data-test-card="${name}"]`,
  counts: {
    description: '[data-test-counts-description]',
    configDisabled: '[data-test-counts-disabled]',
    namespaces: '[data-test-counts-namespaces]',
    mountPaths: '[data-test-counts-auth-mounts]',
    startDiscrepancy: '[data-test-counts-start-discrepancy]',
  },
  dateRange: {
    dropdownOption: (idx = 0) => `[data-test-date-range-billing-start="${idx}"]`,
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
  legend: '[data-test-counts-card-legend]',
  legendDot: (nth: number) => `.legend-item:nth-child(${nth}) > span`,

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

export const FILTERS = {
  dropdown: (name: string) => `[data-test-dropdown="${name}"]`,
  dropdownToggle: (name: string) => `[data-test-dropdown="${name}"] button`,
  dropdownItem: (name: string) => `[data-test-dropdown-item="${name}"]`,
  dropdownSearch: (name: string) => `[data-test-dropdown="${name}"] input`,
  tag: (filter?: string, value?: string) =>
    filter && value ? `[data-test-filter-tag="${filter} ${value}"]` : '[data-test-filter-tag]',
  clearTag: (value: string) => `[aria-label="Dismiss ${value}"]`,
};
