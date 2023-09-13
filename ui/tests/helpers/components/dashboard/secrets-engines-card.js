/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const SELECTORS = {
  secretEnginesTableRows: '[data-test-dashboard-table="${name}"] tr',
  getSecretEngineAccessor: (engineId) => `[data-test-secrets-engines-row=${engineId}] [data-test-accessor]`,
  getSecretEngineDescription: (engineId) =>
    `[data-test-secrets-engines-row=${engineId}] [data-test-description]`,
};

export default SELECTORS;
