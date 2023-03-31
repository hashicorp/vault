/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';

function hasWhitespace([string]) {
  const whitespace = /\s/;
  return whitespace.test(string);
}

export default helper(hasWhitespace);
