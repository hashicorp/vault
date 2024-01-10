/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  configure: '[data-test-pki-issuer-configure]',
  copyButtonByName: (name) => `[data-test-value-div="${name}"] [data-test-copy-button]`,
  crossSign: '[data-test-pki-issuer-cross-sign]',
  defaultGroup: '[data-test-details-group="default"]',
  download: '[data-test-issuer-download]',
  groupTitle: '[data-test-group-title]',
  parsingAlertBanner: '[data-test-parsing-error-alert-banner]',
  rotateModal: '#pki-rotate-root-modal',
  rotateModalGenerate: '[data-test-root-rotate-step-one]',
  rotateRoot: '[data-test-pki-issuer-rotate-root]',
  row: '[data-test-component="info-table-row"]',
  signIntermediate: '[data-test-pki-issuer-sign-int]',
  urlsGroup: '[data-test-details-group="Issuer URLs"]',
  valueByName: (name) => `[data-test-value-div="${name}"]`,
};
