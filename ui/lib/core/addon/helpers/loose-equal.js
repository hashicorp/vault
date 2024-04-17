/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

/*
 * use sparingly *
  ex: logic for an HTML element's selected boolean because <select> values are strings
  strict equal (===) will fail if the API param is a number
  <option selected={{loose-equal model.someAttr someOption)}} value={{someOption}}>
*/
export function looseEqual([a, b]) {
  // loose equal 0 == '' returns true, we don't want that
  if ((a === 0 && b === '') || (a === '' && b === 0)) {
    return false;
  }
  return a == b;
}

export default helper(looseEqual);
