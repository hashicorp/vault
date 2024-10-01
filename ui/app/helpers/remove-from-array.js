/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

function dedupe(items) {
  return items.filter((v, i) => items.indexOf(v) === i);
}

export function removeManyFromArray(array, toRemove) {
  assert(`Both values must be an array`, Array.isArray(array) && Array.isArray(toRemove));
  const a = [...(array || [])];
  return a.filter((v) => !toRemove.includes(v));
}

export function removeFromArray(array, string) {
  assert(`Value provided is not an array`, Array.isArray(array));
  const newArray = [...array];
  const idx = newArray.indexOf(string);
  if (idx >= 0) {
    newArray.splice(idx, 1);
  }
  return dedupe(newArray);
}

export default buildHelper(function ([array, string]) {
  if (Array.isArray(string)) {
    return removeManyFromArray(array, string);
  }
  return removeFromArray(array, string);
});
