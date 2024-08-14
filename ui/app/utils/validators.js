/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent } from '@ember/utils';
import { capitalize } from '@ember/string';

/*
* Model Validators
these return false when the condition fails because false means "invalid"
for example containsWhiteSpace returns "false" when a value HAS whitespace 
because that is an invalid value
*/
export const presence = (value) => isPresent(value);

export const length = (value, { nullable = false, min, max } = {}) => {
  if (!min && !max) return;
  // value could be an integer if the attr has a default value of some number
  const valueLength = value?.toString().length;
  if (valueLength) {
    const underMin = min && valueLength < min;
    const overMax = max && valueLength > max;
    return underMin || overMax ? false : true;
  }
  return nullable;
};

export const number = (value, { nullable = false } = {}) => {
  // since 0 is falsy, !value returns true even though 0 is a valid number
  if (!value && value !== 0) return nullable;
  return !isNaN(value);
};

export const containsWhiteSpace = (value) => {
  return !hasWhitespace(value);
};

export const endsInSlash = (value) => {
  const validation = new RegExp('/$');
  return !validation.test(value);
};

/*
* General Validators
these utils return true or false relative to the function name
*/

export const hasWhitespace = (value) => {
  const validation = new RegExp('\\s', 'g'); // search for whitespace
  return validation.test(value);
};

// HTML form inputs transform values to a string type
// this returns if the value can be evaluated as non-string, i.e. "null"
export const isNonString = (value) => {
  try {
    // if parsable the value could be an object, array, number, null, true or false
    JSON.parse(value);
    return true;
  } catch (e) {
    return false;
  }
};

export const WHITESPACE_WARNING = (item) =>
  `${capitalize(
    item
  )} contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.`;

export const NON_STRING_WARNING =
  'This value will be saved as a string. If you need to save a non-string value, please use the JSON editor.';

export default {
  presence,
  length,
  number,
  containsWhiteSpace,
  endsInSlash,
  isNonString,
  hasWhitespace,
  WHITESPACE_WARNING,
  NON_STRING_WARNING,
};
