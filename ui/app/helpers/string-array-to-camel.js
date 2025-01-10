/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

// This helper is similar to the Ember string camelize helper but it does some additional handling:
// it allows you to pass in an array of strings
// it lowercases the entire string before converting to camelCase preventing situations like IAM Endpoint  -> iamEndpoint instead of iAMEndpoint
// it does not handle accented characters so try not use for user inputted strings.
export function stringArrayToCamelCase(str) {
  if (!str) return;
  if (Array.isArray(str)) {
    return str.map((s) => {
      assert(`must pass in a string or array of strings`, typeof s === 'string');
      // lower case the entire string to handle situations like IAM Endpoint  -> iamEndpoint instead of iAMEndpoint
      s = s.toLowerCase();
      return s
        .replace(/(?:^\w|[A-Z]|\b\w)/g, function (word, index) {
          return index === 0 ? word.toLowerCase() : word.toUpperCase();
        })
        .replace(/\s+/g, '');
    });
  } else {
    // lower case the entire string to handle situations like IAM Endpoint  -> iamEndpoint instead of iAMEndpoint
    assert(`must pass in a string or array of strings`, typeof str === 'string');
    str = str.toLowerCase();
    return str
      .replace(/(?:^\w|[A-Z]|\b\w)/g, function (word, index) {
        return index === 0 ? word.toLowerCase() : word.toUpperCase();
      })
      .replace(/\s+/g, '');
  }
}

export default buildHelper(stringArrayToCamelCase);
