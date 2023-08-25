/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const SELECTORS = {
  cardTitle: '[data-test-dashboard-secrets-engines-header]',
  secretEnginesTableRows: '[data-test-dashboard-secrets-engines-table] tr',
  getSecretEngineAccessor: (engineId) => `[data-test-secrets-engines-row=${engineId}] [data-test-accessor]`,
  getSecretEngineDescription: (engineId) =>
    `[data-test-secrets-engines-row=${engineId}] [data-test-description]`,
};

export default SELECTORS;
