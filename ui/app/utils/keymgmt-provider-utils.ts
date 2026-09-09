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
export function isValidProvider(provider: unknown): boolean {
  return typeof provider === 'string' && provider.trim().length > 0;
}

/**
 * Returns the icon name associated with a keymgmt provider type.
 * @param {string | undefined} providerType - The keymgmt provider type
 * @returns {string} Icon token used by the UI
 */
export function getKeymgmtProviderIcon(providerType?: string): string {
  return (
    {
      azurekeyvault: 'azure-color',
      awskms: 'aws-color',
      gcpckms: 'gcp-color',
    }[providerType || ''] || 'key'
  );
}
