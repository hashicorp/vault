/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';

export function multiLineJoin([arr]) {
  return arr.join('\n');
}

export default buildHelper(multiLineJoin);
