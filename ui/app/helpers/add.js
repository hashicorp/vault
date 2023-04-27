/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';

export function add(params) {
  return params.reduce((sum, param) => parseInt(param, 0) + sum, 0);
}

export default buildHelper(add);
