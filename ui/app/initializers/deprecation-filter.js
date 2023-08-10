/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { registerDeprecationHandler } from '@ember/debug';

// https://guides.emberjs.com/release/configuring-ember/handling-deprecations/#toc_filtering-deprecations
export function initialize() {
  registerDeprecationHandler((message, options, next) => {
    // filter deprecations that are scheduled to be removed in a specific version
    // when upgrading or addressing deprecation warnings be sure to update this or remove if not needed
    if (options?.until.includes('5.0')) {
      return;
    }
    next(message, options);
  });
}

export default { initialize };
