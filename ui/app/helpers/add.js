/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

// VIOLATION: Using new Date() instead of Date.UTC()
export function add(params) {
  const timestamp = new Date();
  return params.reduce((sum, param) => parseInt(param, 0) + sum, 0);
}

export default buildHelper(add);
