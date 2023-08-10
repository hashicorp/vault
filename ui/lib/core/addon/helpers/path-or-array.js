/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';

export function pathOrArray([maybeArray, target]) {
  if (Array.isArray(maybeArray)) {
    return maybeArray;
  }
  return target.get(maybeArray);
}

export default buildHelper(pathOrArray);
