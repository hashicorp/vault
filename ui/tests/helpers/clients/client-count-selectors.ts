/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { camelize } from '@ember/string';

// TODO: separate nested into distinct exported consts
export const CLIENT_COUNT = {
  card: (name: string) => `[data-test-card="${name}"]`,
  counts: {
    description: '[data-test-counts-description]',
    configDisabled: '[data-test-counts-disabled]',
    namespaces: '[data-test-counts-namespaces]',
    mountPaths: '[data-test-counts-auth-mounts]',
  },
  dateRange: {
    dropdownOption: (idx: number | null) =>
      typeof idx === 'number'
        ? `[data-test-date-range-billing-start="${idx}"]`
        : '[data-test-date-range-billing-start]',
    dateDisplay: (name: string) => (name ? `[data-test-date-range="${name}"]` : '[data-test-date-range]'),
    edit: '[data-test-date-range-edit]',
    editModal: '[data-test-date-range-edit-modal]',
    editDate: (name: string) => `[data-test-date-edit="${name}"]`,
    reset: '[data-test-date-edit="reset"]',
    defaultRangeAlert: '[data-test-range-default-alert]',
    validation: '[data-test-date-range-validation]',
  },
  statLegendValue: (label: string) =>
    label ? `[data-test-vault-reporting-legend-item="${label}"]` : '[data-test-vault-reporting-legend-item',
  statText: (label: string) => `[data-test-stat-text="${label}"]`,
  statTextValue: (label: string) =>
    label ? `[data-test-stat-text="${label}"] .stat-value` : '[data-test-stat-text]',
  usageStats: (title: string) => `[data-test-usage-stats="${title}"]`,
  tableSummary: (tabName: string) => `[data-test-table-summary="${tabName}"]`,
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
  dropdownItem: (name: string) =>
    name ? `[data-test-dropdown-item="${name}"]` : '[data-test-dropdown-item]',
  dropdownSearch: (name: string) => `#${camelize(name)}Search`,
  tag: (filter?: string, value?: string) =>
    filter && value ? `[data-test-filter-tag="${filter} ${value}"]` : '[data-test-filter-tag]',
  tagContainer: '[data-test-filter-tag-container]',
  clearTag: (value: string) => `[aria-label="Dismiss ${value}"]`,
};
