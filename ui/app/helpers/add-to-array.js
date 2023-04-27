/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

function dedupe(items) {
  return items.filter((v, i) => items.indexOf(v) === i);
}

export function addToArray([array, string]) {
  if (!Array.isArray(array)) {
    assert(`Value provided is not an array`, false);
  }
  const newArray = [...array];
  newArray.push(string);
  return dedupe(newArray);
}

export default buildHelper(addToArray);
