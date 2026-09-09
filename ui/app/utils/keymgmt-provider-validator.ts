/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Validates if a provider value is a valid string that can be used for API calls.
 * Provider must be a non-empty string to be valid.
 * This prevents passing objects like { permissionsError: true } to API calls.
 * @param {*} provider - The provider value to validate
 * @returns {boolean} true if provider is a valid non-empty string, false otherwise
 */
export function isValidProvider(provider: any): boolean {
  return provider && typeof provider === 'string' && provider.trim().length > 0;
}
