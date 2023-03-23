/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper as buildHelper } from '@ember/component/helper';

export function filterWildcard([string, array]) {
  if (!string || !array) {
    return;
  }
  if (!string.id && string) {
    string = { id: string };
  }
  const stringId = string.id;
  const filterBy = (stringId) =>
    array.filter((item) => new RegExp('^' + stringId.replace(/\*/g, '.*') + '$').test(item));
  return filterBy(stringId).length;
}

export default buildHelper(filterWildcard);
