/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isKnownExternalPlugin } from 'vault/utils/external-plugin-helpers';

/**
 * Utility functions for semantic version handling
 */

/**
 * Clean a version string by removing prefixes and suffixes
 * @param version - The version string to clean (e.g., "v1.2.3+ent")
 * @returns The cleaned version string (e.g., "1.2.3")
 */
export function cleanVersion(version: string): string {
  return version.replace(/^v/, '').split(/[+-]/)[0] || '';
}

/**
 * Parse a version string into numeric parts
 * @param version - The version string to parse
 * @returns Array of numeric version parts
 */
export function parseVersion(version: string): number[] {
  const cleanVer = cleanVersion(version);
  return cleanVer.split('.').map((n) => parseInt(n) || 0);
}

/**
 * Compare two version strings using semantic version rules
 * @param a - First version to compare
 * @param b - Second version to compare
 * @returns Negative if a < b, positive if a > b, 0 if equal
 */
export function compareVersions(a: string, b: string): number {
  const aParts = parseVersion(a);
  const bParts = parseVersion(b);

  const maxLength = Math.max(aParts.length, bParts.length);
  for (let i = 0; i < maxLength; i++) {
    const aPart = aParts[i] || 0;
    const bPart = bParts[i] || 0;

    if (aPart !== bPart) {
      return aPart - bPart;
    }
  }
  return 0;
}

/**
 * Sort an array of version strings in semantic version order
 * @param versions - Array of version strings to sort
 * @param descending - If true, sort highest version first (default: false)
 * @returns New sorted array (does not mutate original)
 */
export function sortVersions(versions: string[], descending = false): string[] {
  const sorted = versions.slice().sort((a, b) => compareVersions(a, b));
  return descending ? sorted.reverse() : sorted;
}

/**
 * Find the highest version from an array of version strings
 * @param versions - Array of version strings
 * @returns The highest version string, or null if array is empty
 */
export function getHighestVersion(versions: string[]): string | null {
  if (versions.length === 0) return null;

  const sorted = sortVersions(versions, true);
  return sorted[0] || null;
}

/**
 * Check if version A is greater than version B
 * @param a - First version
 * @param b - Second version
 * @returns True if a > b
 */
export function isVersionGreater(a: string, b: string): boolean {
  return compareVersions(a, b) > 0;
}

/**
 * Check if two versions are equal
 * @param a - First version
 * @param b - Second version
 * @returns True if versions are equal
 */
export function areVersionsEqual(a: string, b: string): boolean {
  return compareVersions(a, b) === 0;
}

/**
 * Check if a version string is valid and non-empty
 * @param version - The version string to validate
 * @returns True if the version is valid
 */
export function isValidVersion(version: string): boolean {
  if (!version || typeof version !== 'string') return false;

  const trimmed = version.trim();
  if (trimmed === '' || trimmed === 'null') return false;

  // Basic semantic version pattern check (allows prefixes like 'v' and suffixes like '+ent')
  const semverPattern = /^v?\d+(\.\d+)*([+-].+)?$/;
  const cleanVer = cleanVersion(trimmed);

  return semverPattern.test(`v${cleanVer}`);
}

/**
 * Check if a plugin version is required for external plugins
 * @param pluginType - The plugin type (e.g., 'keymgmt' or 'vault-plugin-secrets-keymgmt')
 * @param pluginVersion - The plugin version string
 * @returns True if the plugin version requirement is satisfied
 */
export function isPluginVersionValidForType(pluginType: string, pluginVersion?: string): boolean {
  if (!pluginType) return false;

  if (isKnownExternalPlugin(pluginType)) {
    // External plugins require a valid version
    return isValidVersion(pluginVersion || '');
  } else {
    // Builtin plugins should not have a version specified
    return !pluginVersion || pluginVersion.trim() === '' || pluginVersion.trim() === 'null';
  }
}
