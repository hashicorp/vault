/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * External plugin utilities for managing external plugin mappings and metadata.
 *
 * This file handles the mapping between external plugin names and their builtin equivalents,
 * providing utilities to determine effective engine types for display and routing purposes.
 */

/**
 * Map of external plugin names to their builtin counterparts.
 * This mapping allows external plugins to use the same UI experience as their builtin equivalents.
 *
 * Future: When the backend provides unique plugin IDs, this mapping can serve as a fallback
 * for external plugins that don't have unique IDs available.
 */
export const EXTERNAL_PLUGIN_TO_BUILTIN_MAP: Record<string, string> = {
  'vault-plugin-secrets-ad': 'ad',
  'vault-plugin-secrets-alicloud': 'alicloud',
  'vault-plugin-secrets-azure': 'azure',
  'vault-plugin-secrets-gcp': 'gcp',
  'vault-plugin-secrets-gcpkms': 'gcpkms',
  'vault-plugin-secrets-keymgmt': 'keymgmt',
  'vault-plugin-secrets-kubernetes': 'kubernetes',
  'vault-plugin-secrets-kv': 'kv',
  'vault-plugin-secrets-mongodbatlas': 'mongodbatlas',
  'vault-plugin-secrets-openldap': 'openldap',
  'vault-plugin-secrets-terraform': 'terraform',
} as const;

/**
 * Get the builtin engine type for a given external plugin name.
 * This function checks the external plugin mapping to find the corresponding builtin type.
 *
 * @param externalPluginName - The name of the external plugin (e.g., "vault-plugin-secrets-keymgmt")
 * @returns The builtin engine type if a mapping exists, otherwise undefined
 */
export function getBuiltinTypeFromExternalPlugin(externalPluginName: string): string | undefined {
  return EXTERNAL_PLUGIN_TO_BUILTIN_MAP[externalPluginName];
}

/**
 * Check if a plugin name is a known external plugin that maps to a builtin.
 *
 * @param pluginName - The plugin name to check
 * @returns True if the plugin name is in the external plugin mapping
 */
export function isKnownExternalPlugin(pluginName: string): boolean {
  return pluginName in EXTERNAL_PLUGIN_TO_BUILTIN_MAP;
}

/**
 * Get the effective engine type for display purposes.
 * For external plugins that have a builtin mapping, returns the builtin type.
 * For other plugins, returns the original type.
 *
 * @param pluginType - The original plugin type
 * @returns The effective type to use for engine metadata lookup
 */
export function getEffectiveEngineType(pluginType: string): string {
  return getBuiltinTypeFromExternalPlugin(pluginType) || pluginType;
}

/**
 * Get the external plugin name for a given builtin engine type.
 * This function performs a reverse lookup on the external plugin mapping.
 *
 * @param builtinType - The builtin engine type (e.g., "keymgmt")
 * @returns The external plugin name if a mapping exists, otherwise null
 */
export function getExternalPluginNameFromBuiltin(builtinType: string): string | null {
  // Find the external plugin name that maps to this builtin type
  for (const [externalName, mappedBuiltin] of Object.entries(EXTERNAL_PLUGIN_TO_BUILTIN_MAP)) {
    if (mappedBuiltin === builtinType) {
      return externalName;
    }
  }
  return null;
}
