/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const KV_LIST = {
  // Page::List in KV
  createSecret: '[data-test-toolbar-create-secret]',
  item: (secret: string | null) => (!secret ? '[data-test-list-item]' : `[data-test-list-item="${secret}"]`),
  filter: `[data-test-kv-list-filter]`,
  listMenuDelete: `[data-test-popup-metadata-delete]`,
  listMenuCreate: `[data-test-popup-create-new-version]`,
  overviewCard: '[data-test-overview-card-container="View secret"]',
  overviewInput: '[data-test-view-secret] input',
  overviewButton: '[data-test-get-secret-detail]',
};
