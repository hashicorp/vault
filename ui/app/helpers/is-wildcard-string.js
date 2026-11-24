/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

export function isWildcardString([string]) {
  if (!string) {
    return false;
  }
  // string is actually an object which is what comes in in the searchSelect component
  if (typeof string === 'object') {
    string = Object.values(string).toString();
  }

  return string.includes('*');
}

export default buildHelper(isWildcardString);
