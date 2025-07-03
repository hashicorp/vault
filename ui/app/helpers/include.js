/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

export default helper(function includes([array, item]) {
  return Array.isArray(array) && array.includes(item);
});
