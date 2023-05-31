/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';

// used when select dropdowns receive an integer from the API but options are strings
// because after selecting an option, selected values are converted to strings
export function numberToString([value], { options }) {
  // only convert value type if options are a string (assumes all options the same type)
  if (typeof value === 'number' && (!options || typeof options.firstObject === 'string')) {
    return value.toString();
  }
  return value;
}

export default helper(numberToString);
