/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  versionDropdown: '[data-test-version-dropdown]',
  version: (number) => `[data-test-version="${number}"]`,
  editMetadata: '[data-test-edit-metadata]',
};
