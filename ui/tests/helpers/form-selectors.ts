/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export const FORM = {
  header: (title?: string) => (title ? `[data-test-form-header="${title}"]` : '[data-test-form-header]'),
  description: (title?: string) =>
    title ? `[data-test-form-description="${title}"]` : '[data-test-form-description]',
};
