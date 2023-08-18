/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

const selectors = {
  ttlFormGroup: '[data-test-ttl-inputs]',
  toggle: '[data-test-ttl-toggle]',
  toggleByLabel: (label) => `[data-test-ttl-toggle="${label}"]`,
  label: '[data-test-ttl-form-label]',
  subtext: '[data-test-ttl-form-subtext]',
  tooltipTrigger: `[data-test-tooltip-trigger]`,
  ttlValue: '[data-test-ttl-value]',
  ttlUnit: '[data-test-select="ttl-unit"]',
  valueInputByLabel: (label) => `[data-test-ttl-value="${label}"]`,
  unitInputByLabel: (label) => `[data-test-ttl-unit="${label}"] [data-test-select="ttl-unit"]`,
};

export default selectors;
