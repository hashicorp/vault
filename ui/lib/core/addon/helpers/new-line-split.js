/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

export function newLineSplit([lines]) {
  return lines.split('\n');
}

export default helper(newLineSplit);
