/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn } from '@ember/test-helpers';
import { SELECTORS as GENERAL } from 'vault/tests/helpers/general-selectors';

export const PAGE = {
  ...GENERAL,
  cta: {
    summary: '[data-test-cta-container] p',
    button: '[data-test-cta-button]',
  },
  create: {
    selectType: (type) => `[data-test-select-destination="${type}"]`,
  },
  form: {
    fillInByAttr: async (attr, value) => {
      // for handling more complex form input elements by attr name
      switch (attr) {
        case 'credentials':
          await click('[data-test-text-toggle]');
          return fillIn('[data-test-text-file-textarea]', value);
        case 'deploymentEnvironments':
          await click('[data-test-input="deploymentEnvironments"] input#deployment');
          await click('[data-test-input="deploymentEnvironments"] input#preview');
          return await click('[data-test-input="deploymentEnvironments"] input#production');
        default:
          return fillIn(`[data-test-input="${attr}"]`, value);
      }
    },

    cancelButton: '[data-test-cancel]',
    saveButton: '[data-test-save]',
    credentials: async () => {
      await click('[data-test-text-toggle]');
      await fillIn('[data-test-text-file-textarea]');
    },
  },
};
