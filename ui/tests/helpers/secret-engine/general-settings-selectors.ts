/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  manageDropdown: '[data-test-manage-dropdown]',
  manageDropdownItem: (name: string) => `[data-test-manage-dropdown-item="${name}"]`,
  label: (name: string) => `[data-test-label="${name}"]`,
  helperText: (name: string) => `[data-test-helper-text="${name}"]`,
  ttlPickerV2: '[data-test-ttl-picker-v2]',
  versionCard: {
    engineType: '[data-test-engine-type]',
    currentVersion: '[data-test-engine-current-version]',
    versionsDropdown: `[data-test-versions-dropdown]`,
  },
};
