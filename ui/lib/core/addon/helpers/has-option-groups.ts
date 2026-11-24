/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';
/**
 * Helper to detect if possibleValues contains grouped options (with optgroups).
 * Returns true if any item in the array has a 'group' property, false otherwise.
 *
 * @param {Array} possibleValues - Array that may contain flat options or grouped options
 * @returns {boolean} - True if grouped options detected, false for flat arrays
 */
export function hasOptionGroups(possibleValues: unknown[] | null | undefined): boolean {
  if (!possibleValues || !Array.isArray(possibleValues) || possibleValues.length === 0) {
    return false;
  }
  return possibleValues.some((item) => {
    if (!item || typeof item !== 'object') {
      return false;
    }
    const candidate = item as { group?: unknown; options?: unknown };
    return typeof candidate.group === 'string' && Array.isArray(candidate.options);
  });
}
export default helper(function hasOptionGroupsHelper([possibleValues]: [
  unknown[] | null | undefined,
]): boolean {
  return hasOptionGroups(possibleValues);
});
