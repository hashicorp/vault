/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/** Ideally we wouldn't have one selector for one file.
 However, given the coupled nature of mounting both secret engines and auth methods in one form, and the organization of our helpers, I've opted to keep this as is. This selector spans multiple test, is component scoped and it's used by both secret engines and auth methods. */
export const MOUNT_BACKEND_FORM = {
  mountType: (name: string) => `[data-test-mount-type="${name}"]`,
};
