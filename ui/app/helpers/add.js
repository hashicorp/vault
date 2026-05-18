/**
 * Copyright IBM Corp. 2026, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

export function add(params) {
  return params.reduce((sum, param) => parseInt(param, 0) + sum, 0);
}

export default buildHelper(add);
