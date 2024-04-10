/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const TTL_PICKER = {
  ttlFormGroup: '[data-test-ttl-inputs]',
  toggle: '[data-test-ttl-toggle]',
  toggleByLabel: (label: string) => `[data-test-ttl-toggle="${label}"]`,
  label: '[data-test-ttl-form-label]',
  subtext: '[data-test-ttl-form-subtext]',
  tooltipTrigger: `[data-test-tooltip-trigger]`,
  ttlValue: '[data-test-ttl-value]',
  ttlUnit: '[data-test-select="ttl-unit"]',
  valueInputByLabel: (label: string) => `[data-test-ttl-value="${label}"]`,
  unitInputByLabel: (label: string) => `[data-test-ttl-unit="${label}"] [data-test-select="ttl-unit"]`,
};
