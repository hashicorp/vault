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
