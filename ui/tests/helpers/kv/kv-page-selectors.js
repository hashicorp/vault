/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  details: {
    versionDropdown: '[data-test-version-dropdown]',
    version: (number) => `[data-test-version="${number}"]`,
    editMetadataBtn: '[data-test-edit-metadata]',
  },
  edit: {
    kvRow: '[data-test-kv-row]',
    inputByAttr: (attr) => `[data-test-input="${attr}"]`,
    automateSecretDeletion: '[data-test-ttl-value="Automate secret deletion"]',
    metadataCancel: '[data-test-kv-metadata-cancel]',
  },
};
