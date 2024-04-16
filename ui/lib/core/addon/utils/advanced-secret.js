/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Method to check whether the secret value is a nested object (returns true)
 * All other values return false
 * @param value string or stringified JSON
 * @returns boolean
 */
export function isAdvancedSecret(value) {
  try {
    const obj = typeof value === 'string' ? JSON.parse(value) : value;
    if (Array.isArray(obj)) return false;
    return Object.values(obj).any((value) => typeof value !== 'string');
  } catch (e) {
    return false;
  }
}

/**
 * Method to obfuscate all values in a map, including nested values and arrays
 * @param obj object
 * @returns object
 */
export function obfuscateData(obj) {
  if (typeof obj !== 'object' || Array.isArray(obj)) return obj;
  const newObj = {};
  for (const key of Object.keys(obj)) {
    if (Array.isArray(obj[key])) {
      newObj[key] = obj[key].map(() => '********');
    } else if (typeof obj[key] === 'object') {
      newObj[key] = obfuscateData(obj[key]);
    } else {
      newObj[key] = '********';
    }
  }
  return newObj;
}
