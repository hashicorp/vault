/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const NAMESPACE_PICKER_SELECTORS = {
  link: (link) => (link ? `[data-test-namespace-link="${link}"]` : '[data-test-namespace-link]'),
  searchInput: 'input[type="search"]',
};
