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

export const containsDataOctet = (value) => {
  return !hasDataOctet(value);
};

export const containsForwardSlash = (value) => {
  return !hasForwardSlash(value);
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

export const hasDataOctet = (value) => {
  // A percent-encoded data octet is a character triplet that represents a byte's numeric value in a Uniform Resource Identifier (URI):
  // Format: A percent sign (%) followed by two hexadecimal digits
  // Example: The percent-encoding for / is %2f
  // In KVv2 we want to warn users that their secret path includes a percent-encoded data octet and that we will not transform it
  const regex = /%([0-9A-Fa-f]{2})/g;
  return !!value.match(regex);
};

export const hasForwardSlash = (value) => {
  // only show if forward slash is not the last value. If it's the last value the endsInSlash validator will catch it.
  const notLastChar = value.slice(0, -1);
  const regex = /\//g;
  return regex.test(notLastChar);
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

export const DATA_OCTET_WARNING = (item) =>
  `${capitalize(item)} contains a percent encoded data octet. The UI will not decode this.`;

export const FORWARD_SLASH_WARNING = (item) =>
  `${capitalize(
    item
  )} contains a forward slash. The UI will interpret this as the name of a directory. Example: foo/bar where foo will be the directory name and foo the secret path.`;

export const NON_STRING_WARNING =
  'This value will be saved as a string. If you need to save a non-string value, please use the JSON editor.';

export default {
  presence,
  length,
  number,
  containsWhiteSpace,
  containsDataOctet,
  containsForwardSlash,
  endsInSlash,
  isNonString,
  hasWhitespace,
  WHITESPACE_WARNING,
  NON_STRING_WARNING,
};
