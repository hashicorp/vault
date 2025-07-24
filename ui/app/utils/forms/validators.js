/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent } from '@ember/utils';
import { capitalize } from '@ember/string';

/*
 * Validators
 *
 * These functions are most often used by the validation decorator to check if a field is valid.
 * Each validator returns:
 *   - `true`  → the value PASSES validation (is valid) and no message is shown in the UI.
 *   - `false` → the value FAILS validation (is invalid), triggering the associated error message in the UI.
 *
 * Validators are referenced by name in a validations object:
 *
 *   const validations = {
 *     path: [
 *       { type: 'presence', message: `Path can't be blank.` },
 *       { type: 'noEndingSlash', message: `Path can't end with '/'` },
 *       { type: 'noWhitespace', message: WHITESPACE_WARNING('path'), level: 'warn' },
 *     ],
 *   };
 *
 * Examples of return values:
 *   presence('abc')                     → true
 *   length('abc', { min: 5 })           → false
 *   number('12')                        → true
 *   noWhitespace('foo bar')             → false
 *   containsWhitespace('foo bar')       → true
 *   noTrailingWhitespace('foo ')        → false
 *   noEndingSlash('foo/')               → false
 */

export const presence = (value) => isPresent(value);

export const length = (value, { nullable = false, min, max } = {}) => {
  if (!min && !max) return true; // nothing to validate against
  // value could be an integer if the attr has a default numeric value
  const valueLength = value?.toString().length;
  if (valueLength) {
    const underMin = min && valueLength < min;
    const overMax = max && valueLength > max;
    return !(underMin || overMax);
  }
  return nullable;
};

export const number = (value, { nullable = false } = {}) => {
  // since 0 is falsy, !value returns true even though 0 is a valid number
  if (!value && value !== 0) return nullable;
  return !isNaN(value);
};

/**
 * Returns true if ANY whitespace exists.
 * Not something to use in the validators decorator because we want to return false if whitespace exists.
 */
export const containsWhitespace = (value) => /\s/.test(value);

/**
 * Returns true if the value contains NO whitespace characters.
 * This is used in the validators decorator because we want to return false if whitespace exists and show a warning.
 */
export const noWhitespace = (value) => !/\s/.test(value);

/**
 * Returns false if the value ends in a slash.
 */
export const noEndingSlash = (value) => !/\/$/.test(value);

// -------------------------
// General validators/utilities
// -------------------------

// HTML inputs coerce to strings; this checks if the value can be parsed as a non-string JSON type.
export const canParseToNonString = (value) => {
  try {
    JSON.parse(value); // could be object, array, number, null, true, false
    return true;
  } catch (e) {
    return false;
  }
};

// -------------------------
// Messages
// -------------------------

export const WHITESPACE_WARNING = (item) =>
  `${capitalize(
    item
  )} contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.`;

export const NON_STRING_WARNING =
  'This value will be saved as a string. If you need to save a non-string value, please use the JSON editor.';

// -------------------------
// Default export
// -------------------------

export default {
  presence,
  length,
  number,
  containsWhitespace,
  noWhitespace,
  noEndingSlash,
  canParseToNonString,
  WHITESPACE_WARNING,
  NON_STRING_WARNING,
};
