/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  cardName: (name) => `[data-test-card="${name}"]`,
  emptyState: (name) => `[data-test-empty-state="${name}"]`,
};
