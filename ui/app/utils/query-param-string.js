/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isEmptyValue } from 'core/helpers/is-empty-value';

/**
 * queryParamString converts an object to a query param string with URL encoded keys and values.
 * It does not include values that are considered invalid query params (below).
 * @param {object} queryObject with key-value pairs of desired URL params
 * @returns string like ?key=val1&key2=val2
 */

// we can't just rely on falsy values because "false" and "0" are valid query params
const INVALID_QP = [undefined, null, ''];

export default function queryParamString(queryObject) {
  if (
    !queryObject ||
    isEmptyValue(queryObject) ||
    typeof queryObject !== 'object' ||
    Array.isArray(queryObject) ||
    Object.values(queryObject).every((v) => INVALID_QP.includes(v))
  )
    return '';
  return Object.keys(queryObject).reduce((prev, key) => {
    const value = queryObject[key];
    if (INVALID_QP.includes(value)) return prev;
    const keyval = `${encodeURIComponent(key)}=${encodeURIComponent(value)}`;
    if (prev === '?') {
      return `${prev}${keyval}`;
    }
    return `${prev}&${keyval}`;
  }, '?');
}
