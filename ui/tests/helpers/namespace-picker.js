/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const NAMESPACE_PICKER_SELECTORS = {
  link: (link) => (link ? `[data-test-namespace-link="${link}"]` : '[data-test-namespace-link]'),
  refreshList: '[data-test-refresh-namespaces]',
  toggle: '[data-test-toggle-input="namespace-id"]',
  searchInput: 'input[type="search"]',
  manageButton: '[data-test-manage-namespaces]',
};
