/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

function dedupe(items) {
  return items.filter((v, i) => items.indexOf(v) === i);
}

export function addManyToArray(array, otherArray) {
  assert(`Both values must be an array`, Array.isArray(array) && Array.isArray(otherArray));
  const newArray = [...array].concat(otherArray);
  return dedupe(newArray);
}

export function addToArray(array, string) {
  if (!Array.isArray(array)) {
    assert(`Value provided is not an array`, false);
  }
  const newArray = [...array];
  newArray.push(string);
  return dedupe(newArray);
}

export default buildHelper(function ([array, string]) {
  if (Array.isArray(string)) {
    return addManyToArray(array, string);
  }
  return addToArray(array, string);
});
