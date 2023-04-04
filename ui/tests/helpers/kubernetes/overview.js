/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  rolesCardTitle: '[data-test-selectable-card="Roles"] .title',
  rolesCardSubTitle: '[data-test-selectable-card-container="Roles"] p',
  rolesCardLink: '[data-test-selectable-card="Roles"] a',
  rolesCardNumRoles: '[data-test-roles-card-overview-num]',
  generateCredentialsCardTitle: '[data-test-selectable-card="Generate credentials"] .title',
  generateCredentialsCardSubTitle: '[data-test-selectable-card-container="Generate credentials"] p',
  generateCredentialsCardButton: '[data-test-generate-credential-button]',
  emptyStateTitle: '.empty-state .empty-state-title',
  emptyStateMessage: '.empty-state .empty-state-message',
  emptyStateActionText: '.empty-state .empty-state-actions',
};
