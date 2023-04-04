/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  defaultGroup: '[data-test-details-group="default"]',
  urlsGroup: '[data-test-details-group="Issuer URLs"]',
  groupTitle: '[data-test-group-title]',
  row: '[data-test-component="info-table-row"]',
  rotateRoot: '[data-test-pki-issuer-rotate-root]',
  crossSign: '[data-test-pki-issuer-cross-sign]',
  signIntermediate: '[data-test-pki-issuer-sign-int]',
  download: '[data-test-issuer-download]',
  configure: '[data-test-pki-issuer-configure]',
  rotateModal: '[data-test-modal-background="Rotate this root"]',
  rotateModalGenerate: '[data-test-root-rotate-step-one]',
  valueByName: (name) => `[data-test-value-div="${name}"]`,
};
