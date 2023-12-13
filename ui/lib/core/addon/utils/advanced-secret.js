/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * Method to check whether the secret value is a nested object (returns true)
 * All other values return false
 * @param value string or stringified JSON
 * @returns boolean
 */
export function isAdvancedSecret(value) {
  try {
    const json = JSON.parse(value);
    if (Array.isArray(json)) return false;
    return Object.values(json).some((value) => typeof value !== 'string');
  } catch (e) {
    return false;
  }
}
