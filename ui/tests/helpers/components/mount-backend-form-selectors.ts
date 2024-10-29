/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const MOUNT_BACKEND_FORM = {
  header: '[data-test-mount-form-header]',
  mountType: (name: string) => `[data-test-mount-type="${name}"]`,
};
