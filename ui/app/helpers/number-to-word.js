/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

export function numberToWord(number, capitalize) {
  const word =
    {
      0: 'zero',
      1: 'one',
      2: 'two',
      3: 'three',
      4: 'four',
      5: 'five',
      6: 'six',
      7: 'seven',
      8: 'eight',
      9: 'nine',
    }[number] || number;
  return capitalize && typeof word === 'string' ? `${word.charAt(0).toUpperCase()}${word.slice(1)}` : word;
}

export default helper(function ([number], { capitalize }) {
  return numberToWord(number, capitalize);
});
