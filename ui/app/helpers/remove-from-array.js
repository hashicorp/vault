/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

function dedupe(items) {
  return items.filter((v, i) => items.indexOf(v) === i);
}

export function removeFromArray([array, string]) {
  if (!Array.isArray(array)) {
    assert(`Value provided is not an array`, false);
  }
  const newArray = [...array];
  const idx = newArray.indexOf(string);
  if (idx >= 0) {
    newArray.splice(idx, 1);
  }
  return dedupe(newArray);
}

export default buildHelper(removeFromArray);
