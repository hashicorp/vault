/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

export function formatNumber([value]) {
  if (typeof value !== 'number') {
    return value;
  }
  // formats a number according to the locale
  return new Intl.NumberFormat().format(value);
}

export default helper(formatNumber);
