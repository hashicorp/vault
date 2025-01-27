/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const LDAP_SELECTORS = {
  roleItem: (type: string, name: string) => `[data-test-role="${type} ${name}"]`,
  libraryItem: (name: string) => `[data-test-library="${name}"]`,
  roleMenu: (type: string, name: string) => `[data-test-popup-menu-trigger="${type} ${name}"]`,
  libraryMenu: (name: string) => `[data-test-popup-menu-trigger="${name}"]`,
  action: (action: string) => `[data-test-${action}]`,
};
