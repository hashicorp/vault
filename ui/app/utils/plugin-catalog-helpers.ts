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
  deprecation_status?: string;
  isAvailable?: boolean;
  pluginData?: PluginCatalogPlugin;
  isExternalPlugin?: boolean; // Flag to indicate this engine was discovered from catalog, not static metadata
}

/**
 * Function to merge plugin versions with static engine data and add dynamic plugin discovery
 * Enhanced to support plugins from catalog that don't exist in static metadata
 *
 * @param staticEngines - Array of static engine metadata
 * @param secretEnginesList - Array of secret engine names from catalog
 * @param secretEnginesDetailed - Array of detailed secret engine info from catalog
 * @param databasePluginsList - Array of database plugin names from catalog
 * @param databasePluginsDetailed - Array of detailed database plugin info from catalog
 * @returns Enhanced engines with version information and dynamically discovered plugins
 */
export function addVersionsToEngines(
  allEngines: EngineDisplayData[],
  secretEnginesList: string[] = [],
  secretEnginesDetailed: PluginCatalogPlugin[] = [],
  databasePluginsList: string[] = [],
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
      const isDatabaseAvailable = databasePluginsList.length > 0 || databasePluginsDetailed.length > 0;

      if (isDatabaseAvailable) {
        // Find a representative database plugin for version info (prefer first available)
        const representativePlugin = databasePluginsDetailed[0];

        return {
          ...engine,
          builtin: representativePlugin?.builtin ?? true, // Database engine is typically builtin
          deprecation_status: representativePlugin?.deprecation_status || 'supported',
          version: representativePlugin?.version,
          isAvailable: true,
          isExternalPlugin: false,
          pluginData: representativePlugin,
        };
      } else {
        return {
          ...engine,
          isAvailable: false,
          isExternalPlugin: false,
        };
      }
    }

    const pluginData = secretEnginesDetailed.find((plugin) => plugin.name === engine.type);

    if (pluginData) {
      return {
        ...engine,
        builtin: pluginData.builtin,
        deprecation_status: pluginData.deprecation_status,
        version: pluginData.version,
        isAvailable: true,
        isExternalPlugin: false, // Static engines are not external
        pluginData,
      };
    }

    return {
      ...engine,
      isAvailable: false,
      isExternalPlugin: false, // Static engines are not external even if unavailable
    };
  });

  // Find plugins in catalog that don't exist in static metadata
  const staticEngineTypes = new Set(allEngines.map((engine) => engine.type));
  const dynamicPlugins: EnhancedEngineDisplayData[] = [];

  // Create a map of detailed plugin information for quick lookup
  const detailedPluginMap = new Map<string, PluginCatalogPlugin>();
  secretEnginesDetailed.forEach((plugin) => {
    detailedPluginMap.set(plugin.name, plugin);
  });

  // Process secret engines from the categorized list (if available)
  if (secretEnginesList.length > 0) {
    secretEnginesList.forEach((secretEngineName) => {
      // Skip if this plugin already exists in static metadata
      if (staticEngineTypes.has(secretEngineName)) {
        return;
      }

      // Look up detailed information for this secret engine
      const detailedInfo = detailedPluginMap.get(secretEngineName);

      // Create dynamic engine metadata with defaults
      const dynamicEngine: EnhancedEngineDisplayData = {
        type: secretEngineName,
        displayName: secretEngineName
          .split('-')
          .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(' '), // Convert kebab-case to Title Case
        mountCategory: ['secret'],
        pluginCategory: 'external', // Mark as external since it's not in static metadata
        glyph: DEFAULT_EXTERNAL_PLUGIN_GLYPH,
        isAvailable: true,
        // Use detailed info if available, otherwise create minimal plugin data
        builtin: detailedInfo?.builtin ?? false,
        deprecation_status: detailedInfo?.deprecation_status,
        version: detailedInfo?.version,
        pluginData: detailedInfo || {
          name: secretEngineName,
          type: 'secret',
          builtin: false,
          version: 'unknown',
        },
        isExternalPlugin: true,
      };

      dynamicPlugins.push(dynamicEngine);
    });
  } else {
    // Fallback: process plugins from detailed list if no categorized list available
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
      }); // Create dynamic engine metadata with defaults
      const dynamicEngine: EnhancedEngineDisplayData = {
        type: plugin.name,
        displayName: plugin.name
          .split('-')
          .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(' '), // Convert kebab-case to Title Case
        mountCategory: ['secret'],
        pluginCategory: 'external', // Mark as external since it's not in static metadata
        glyph: matchingStaticEngine?.glyph || DEFAULT_EXTERNAL_PLUGIN_GLYPH, // Use glyph from matching type or default
        isAvailable: true,
        builtin: plugin.builtin,
        deprecation_status: plugin.deprecation_status,
        version: plugin.version,
        pluginData: plugin,
        isExternalPlugin: true,
      };

      dynamicPlugins.push(dynamicEngine);
    });
  }

  return [...enhancedEngines, ...dynamicPlugins];
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
