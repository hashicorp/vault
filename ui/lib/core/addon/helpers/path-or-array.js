/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

export function pathOrArray([maybeArray, target]) {
  if (Array.isArray(maybeArray)) {
    return maybeArray;
  }
  return target[maybeArray];
}

export default buildHelper(pathOrArray);
