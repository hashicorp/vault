/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { EngineDisplayData } from './all-engines-metadata';

/**
 * Constants for plugin catalog functionality
 */
const DEFAULT_EXTERNAL_PLUGIN_GLYPH = '';

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
    database?: string[];
  };
}

/**
 * Enhanced engine display data with plugin catalog information
 */
export interface EnhancedEngineDisplayData extends EngineDisplayData {
  version?: string;
  builtin?: boolean;
  deprecationStatus?: string;
  isAvailable?: boolean;
  pluginData?: PluginCatalogPlugin;
}

/**
 * Function to merge plugin versions with static engine data and add external plugin discovery
 * Enhanced to support plugins from catalog that don't exist in static metadata
 *
 * @param staticEngines - Array of static engine metadata
 * @param secretEnginesDetailed - Array of detailed secret engine info from catalog
 * @param databasePluginsDetailed - Array of detailed database plugin info from catalog
 * @returns Enhanced engines with version information and dynamically discovered external plugins
 */
export function addVersionsToEngines(
  allEngines: EngineDisplayData[],
  secretEnginesDetailed: PluginCatalogPlugin[] = [],
  databasePluginsDetailed: PluginCatalogPlugin[] = []
): EnhancedEngineDisplayData[] {
  if (!secretEnginesDetailed || !Array.isArray(secretEnginesDetailed)) {
    return allEngines;
  }

  // First, enhance existing static engines with catalog data
  const enhancedEngines = allEngines.map((engine) => {
    // Special handling for Database engine
    if (engine.type === 'database') {
      // Database engine is available if there are any database plugins
      const isDatabaseAvailable = databasePluginsDetailed.length > 0;

      if (isDatabaseAvailable) {
        // Find a representative database plugin for version info (prefer first available)
        const representativePlugin = databasePluginsDetailed[0];

        return {
          ...engine,
          builtin: representativePlugin?.builtin,
          deprecationStatus: representativePlugin?.deprecation_status,
          version: representativePlugin?.version,
          isAvailable: true,
          pluginData: representativePlugin,
        };
      } else {
        return {
          ...engine,
          isAvailable: false,
        };
      }
    }

    const pluginData = secretEnginesDetailed.find((plugin) => plugin.name === engine.type);

    if (pluginData) {
      return {
        ...engine,
        builtin: pluginData.builtin,
        deprecationStatus: pluginData.deprecation_status,
        version: pluginData.version,
        isAvailable: true,
        pluginData,
      };
    }

    return {
      ...engine,
      isAvailable: false,
    };
  });

  // Find plugins in catalog that don't exist in static metadata
  const staticEngineTypes = new Set(allEngines.map((engine) => engine.type));
  const externalPlugins: EnhancedEngineDisplayData[] = [];

  // Process secret engines from the detailed array
  secretEnginesDetailed.forEach((plugin) => {
    // Skip if this plugin already exists in static metadata
    if (staticEngineTypes.has(plugin.name)) {
      return;
    }

    // Look for a static engine type that this plugin name contains or matches
    // This handles cases like "my-custom-aws-plugin" matching the "aws" static engine
    const matchingStaticEngine = allEngines.find((engine) => {
      // Direct type match (e.g., plugin.name = "aws" matches engine.type = "aws")
      if (plugin.name === engine.type) {
        return true;
      }
      // Pattern match (e.g., plugin.name = "my-custom-aws-plugin" contains "aws")
      return plugin.name.includes(engine.type) || plugin.name.includes(engine.type.replace('-', ''));
    });

    // Create external engine metadata with defaults
    const externalEngine: EnhancedEngineDisplayData = {
      type: plugin.name,
      displayName: plugin.name
        .split('-')
        .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' '), // Convert kebab-case to Title Case
      mountCategory: ['secret'],
      pluginCategory: 'external', // Mark as external since it's not in static metadata
      glyph: matchingStaticEngine?.glyph || DEFAULT_EXTERNAL_PLUGIN_GLYPH, // Use glyph from matching type or default
      isAvailable: true,
      builtin: plugin.builtin,
      deprecationStatus: plugin.deprecation_status,
      version: plugin.version,
      pluginData: plugin,
    };

    externalPlugins.push(externalEngine);
  });

  return [...enhancedEngines, ...externalPlugins];
}

/**
 * Validates plugin catalog response structure to ensure it contains the expected data format.
 * Checks for the presence of the required `data.detailed` array which contains plugin information.
 *
 * @param response - Raw response from plugin catalog API endpoint
 * @returns boolean indicating if response has valid structure for plugin catalog data
 *
 * @example
 * ```typescript
 * const response = await api.getPluginCatalog();
 * if (isValidPluginCatalogResponse(response)) {
 *   // Safe to access response.data.detailed, response.data.secret, etc.
 *   const plugins = response.data.detailed;
 * }
 * ```
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
 * Categorizes engines by their availability status for display purposes.
 * Used in Phase 4 to separate enabled and disabled plugins in the UI.
 * Engines with isAvailable === false are considered disabled.
 *
 * @param engines - Array of enhanced engine data with availability information
 * @returns Object containing separate arrays for enabled and disabled engines
 *
 * @example
 * ```typescript
 * const engines = addVersionsToEngines(staticEngines, catalogData);
 * const { enabled, disabled } = categorizeEnginesByStatus(engines);
 *
 * // Render enabled plugins first, then disabled plugins with different styling
 * enabled.forEach(engine => renderEnabledPlugin(engine));
 * disabled.forEach(engine => renderDisabledPlugin(engine));
 * ```
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
