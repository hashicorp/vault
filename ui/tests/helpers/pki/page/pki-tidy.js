/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import { SELECTORS as TIDY_FORM } from './pki-tidy-form';

export const SELECTORS = {
  hdsAlertTitle: '[data-test-tidy-status-alert-title]',
  hdsAlertDescription: '[data-test-tidy-status-alert-description]',
  alertUpdatedAt: '[data-test-tidy-status-alert-updated-at]',
  cancelTidyAction: '[data-test-cancel-tidy-action]',
  hdsAlertButtonText: '[data-test-cancel-tidy-action] .hds-button__text',
  timeStartedRow: '[data-test-value-div="Time started"]',
  timeFinishedRow: '[data-test-value-div="Time finished"]',
  cancelTidyModalBackground: '#pki-cancel-tidy-modal',
  tidyEmptyStateConfigure: '[data-test-tidy-empty-state-configure]',
  manualTidyToolbar: '[data-test-pki-manual-tidy-config]',
  autoTidyToolbar: '[data-test-pki-auto-tidy-config]',
  tidyConfigureModal: {
    configureTidyModal: '#pki-tidy-modal',
    tidyModalAutoButton: '[data-test-tidy-modal-auto-button]',
    tidyModalManualButton: '[data-test-tidy-modal-manual-button]',
    tidyModalCancelButton: '[data-test-tidy-modal-cancel-button]',
    tidyOptionsModal: '[data-test-pki-tidy-options-modal]',
  },
  tidyEmptyState: '[data-test-component="empty-state"]',
  tidyForm: {
    ...TIDY_FORM,
  },
};
