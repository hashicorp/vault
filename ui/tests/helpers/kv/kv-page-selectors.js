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
};
