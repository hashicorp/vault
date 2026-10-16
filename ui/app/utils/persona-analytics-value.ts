/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Formats a persona preference for analytics. Predefined personas are sent as-is;
 * "other" is sent as `other-<kebab-role>` when a custom role is provided, or
 * `other` when it is empty or normalises to nothing.
 *
 * Examples:
 *   { value: 'developer' }                          → "developer"
 *   { value: 'other' }                              → "other"
 *   { value: 'other', customRole: 'DevOps SRE' }    → "other-devops-sre"
 */
export default function personaAnalyticsValue(value: string, customRole?: string): string {
  if (value !== 'other' || !customRole) return value;

  const kebab = customRole
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');

  return kebab ? `other-${kebab}` : 'other';
}
