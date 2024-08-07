/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

export function isEmptyValue(value, hasDefault = false) {
  if (hasDefault) {
    value = hasDefault;
  }
  if (typeof value === 'object' && value !== null) {
    return Object.keys(value).length === 0;
  }
  return value == null || value === '';
}

export default helper(function ([value], { hasDefault = false }) {
  return isEmptyValue(value, hasDefault);
});
