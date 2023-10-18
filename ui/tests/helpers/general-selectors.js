/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx) => `[data-test-crumb="${idx}"] a`,
  breadcrumbs: '[data-test-breadcrumbs]',
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  title: '[data-test-page-title]',
  headerContainer: 'header.page-header',
};
