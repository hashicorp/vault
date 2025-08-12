/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { EngineDisplayData } from './all-engines-metadata';

/**
 * Plugin catalog response structure from Vault API
 */
export interface PluginCatalogPlugin {
  name: string;
  type: string;
  builtin: boolean;
  version: string;
  deprecation_status?: string;
  oci_image?: string;
  runtime?: string;
}

export interface PluginCatalogResponse {
  data: {
    detailed: PluginCatalogPlugin[];
    secret?: string[];
    auth?: string[];
  };
}

/**
 * Enhanced engine display data with plugin catalog information
 */
export interface EnhancedEngineDisplayData extends EngineDisplayData {
  version?: string;
  builtin?: boolean;
  deprecation_status?: string;
  isAvailable?: boolean;
  pluginData?: PluginCatalogPlugin;
}

/**
 * Simple function to merge plugin versions with static engine data
 * This is the minimal Phase 1 implementation - just adds version info
 *
 * @param staticEngines - Array of static engine metadata
 * @param pluginCatalogData - Array of plugin data from catalog API
 * @returns Enhanced engines with version information
 */
export function addVersionsToEngines(
  staticEngines: EngineDisplayData[],
  pluginCatalogData: PluginCatalogPlugin[]
): EnhancedEngineDisplayData[] {
  if (!pluginCatalogData || !Array.isArray(pluginCatalogData)) {
    return staticEngines;
  }

  return staticEngines.map((engine) => {
    const pluginData = pluginCatalogData.find((plugin) => plugin.name === engine.type);

    if (pluginData) {
      return {
        ...engine,
        version: formatPluginVersion(pluginData.version),
        builtin: pluginData.builtin,
        deprecation_status: pluginData.deprecation_status,
        isAvailable: true,
        pluginData,
      };
    }

    return {
      ...engine,
      isAvailable: false,
    };
  });
}

/**
 * Formats a plugin version string for display
 * Removes builtin suffixes ('+builtin', '+builtin.vault') for cleaner presentation
 *
 * @param version - Raw version string from plugin catalog
 * @returns Formatted version string or undefined if no version
 */
export function formatPluginVersion(version?: string): string | undefined {
  if (!version) return undefined;

  // Remove any '+builtin' suffix variations for cleaner display
  // Handles: +builtin.vault, +builtin, and any other +builtin.* patterns
  return version.replace(/\+builtin(\.[^+]*)?$/, '');
}

/**
 * Validates plugin catalog response structure
 *
 * @param response - Raw response from plugin catalog API
 * @returns boolean indicating if response is valid
 */
export function isValidPluginCatalogResponse(response: unknown): response is PluginCatalogResponse {
  if (!response || typeof response !== 'object') {
    return false;
  }

  const typedResponse = response as Record<string, unknown>;

  return Boolean(
    typedResponse['data'] &&
      typeof typedResponse['data'] === 'object' &&
      Array.isArray((typedResponse['data'] as Record<string, unknown>)['detailed'])
  );
}
