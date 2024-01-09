/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// will trim a given set of endings from the end of a string
// if isExtension is true, the first char of that string will be escaped
// in the regex
export default function (str, endings = [], isExtension = true) {
  const prefix = isExtension ? '\\' : '';
  const trimRegex = new RegExp(endings.map((ext) => `${prefix}${ext}$`).join('|'));
  return str.replace(trimRegex, '');
}
