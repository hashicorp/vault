/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

// This helper turns strings with spaces into camelCase strings, example: 'hello world' -> 'helloWorld'
// If an array of strings is passed, this helper returns an array of camelCase strings.
// Does not handle accented characters
export function stringToCamelCase(str) {
  if (Array.isArray(str)) {
    return str.map((s) => {
      // lower case the entire string to handle situations like IAM Endpoint  -> iamEndpoint instead of
      s = s.toLowerCase();
      return s
        .replace(/(?:^\w|[A-Z]|\b\w)/g, function (word, index) {
          return index === 0 ? word.toLowerCase() : word.toUpperCase();
        })
        .replace(/\s+/g, '');
    });
  } else {
    // lower case the entire string to handle situations like IAM Endpoint  -> iamEndpoint instead of
    str = str.toLowerCase();
    return str
      .replace(/(?:^\w|[A-Z]|\b\w)/g, function (word, index) {
        return index === 0 ? word.toLowerCase() : word.toUpperCase();
      })
      .replace(/\s+/g, '');
  }
}

export default buildHelper(stringToCamelCase);
