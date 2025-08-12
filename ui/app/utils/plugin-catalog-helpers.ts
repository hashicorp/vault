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
  demoDisabled?: boolean; // Flag for demo disabled plugins
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

  // DEMO: Force some plugins to be disabled for testing purposes
  // You can delete this section later - it's just for demonstration
  const demoDisabledPlugins = ['consul', 'nomad', 'terraform', 'alicloud', 'mongodbatlas', 'venafi'];

  // DEMO: Add some fake plugins for testing disabled plugin UI
  const demoFakePlugins: EnhancedEngineDisplayData[] = [
    {
      displayName: 'Demo Plugin Alpha',
      type: 'demo-alpha',
      glyph: 'file-text',
      pluginCategory: 'generic',
      mountCategory: ['secret'],
      isAvailable: false,
      demoDisabled: true,
    },
    {
      displayName: 'Example Cloud Service',
      type: 'example-cloud',
      glyph: 'cloud',
      pluginCategory: 'cloud',
      mountCategory: ['secret'],
      isAvailable: false,
      demoDisabled: true,
    },
    {
      displayName: 'Test Infrastructure Tool',
      type: 'test-infra',
      glyph: 'server',
      pluginCategory: 'infra',
      mountCategory: ['secret'],
      isAvailable: false,
      demoDisabled: true,
    },
  ];

  const enhancedEngines = staticEngines.map((engine) => {
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

    // DEMO: Mark certain plugins as disabled for testing ONLY if they're not in the real catalog
    if (demoDisabledPlugins.includes(engine.type)) {
      return {
        ...engine,
        isAvailable: false,
        demoDisabled: true, // Flag to indicate this is a demo disabled plugin
      };
    }

    return {
      ...engine,
      isAvailable: false,
    };
  });

  // DEMO: Add fake plugins to the list
  return [...enhancedEngines, ...demoFakePlugins];
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

/**
 * Categorizes engines by their availability status
 * Used in Phase 4 to separate enabled and disabled plugins
 *
 * @param engines - Array of enhanced engine data
 * @returns Object with enabled and disabled engine arrays
 */
export interface CategorizedEngines {
  enabled: EnhancedEngineDisplayData[];
  disabled: EnhancedEngineDisplayData[];
}

export function categorizeEnginesByStatus(engines: EnhancedEngineDisplayData[]): CategorizedEngines {
  const enabled: EnhancedEngineDisplayData[] = [];
  const disabled: EnhancedEngineDisplayData[] = [];

  engines.forEach((engine) => {
    if (engine.isAvailable !== false) {
      enabled.push(engine);
    } else {
      disabled.push(engine);
    }
  });

  return { enabled, disabled };
}
