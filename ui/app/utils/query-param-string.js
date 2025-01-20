/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isEmptyValue } from 'core/helpers/is-empty-value';

/**
 * queryParamString converts an object to a query param string with URL encoded keys and values.
 * It does not include values that are falsey.
 * @param {object} queryObject with key-value pairs of desired URL params
 * @returns string like ?key=val1&key2=val2
 */
export default function queryParamString(queryObject) {
  if (
    !queryObject ||
    isEmptyValue(queryObject) ||
    typeof queryObject !== 'object' ||
    Array.isArray(queryObject)
  )
    return '';
  return Object.keys(queryObject).reduce((prev, key) => {
    const value = queryObject[key];
    if (!value) return prev;
    const keyval = `${encodeURIComponent(key)}=${encodeURIComponent(value)}`;
    if (prev === '?') {
      return `${prev}${keyval}`;
    }
    return `${prev}&${keyval}`;
  }, '?');
}
