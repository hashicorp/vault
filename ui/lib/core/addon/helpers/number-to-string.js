/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';

// useful for select dropdowns when the API param is an integer but options are strings
// because after selecting an option, selected values are returned as type string
export function numberToString([number]) {
  if (typeof number === 'number') {
    return number.toString();
  }
  return number;
}

export default helper(numberToString);
