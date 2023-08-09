/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';

export default helper(function isEmptyValue([value], { hasDefault = false }) {
  if (hasDefault) {
    value = hasDefault;
  }
  if (typeof value === 'object' && value !== null) {
    return Object.keys(value).length === 0;
  }
  return value == null || value === '';
});
