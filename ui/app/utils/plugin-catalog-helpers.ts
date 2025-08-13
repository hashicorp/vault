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
    },
    {
      displayName: 'Example Cloud Service',
      type: 'example-cloud',
      glyph: 'cloud',
      pluginCategory: 'cloud',
      mountCategory: ['secret'],
      isAvailable: false,
    },
    {
      displayName: 'Test Infrastructure Tool',
      type: 'test-infra',
      glyph: 'server',
      pluginCategory: 'infra',
      mountCategory: ['secret'],
      isAvailable: false,
    },
  ];

  // First, enhance existing static engines with catalog data
  const enhancedEngines = allEngines.map((engine) => {
    // Special handling for Database engine
    if (engine.type === 'database') {
      // Database engine is available if there are any database plugins
      const isDatabaseAvailable = databasePluginsList.length > 0 || databasePluginsDetailed.length > 0;

      if (isDatabaseAvailable) {
        // Find a representative database plugin for version info (prefer first available)
        const representativePlugin = databasePluginsDetailed[0];
        const cleanedVersion =
          representativePlugin?.version?.replace(/\+builtin.*$/, '') || representativePlugin?.version;

        return {
          ...engine,
          builtin: representativePlugin?.builtin ?? true, // Database engine is typically builtin
          deprecation_status: representativePlugin?.deprecation_status || 'supported',
          version: cleanedVersion,
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
      // Clean version string by removing +builtin suffixes
      const cleanedVersion = pluginData.version?.replace(/\+builtin.*$/, '') || pluginData.version;

      return {
        ...engine,
        builtin: pluginData.builtin,
        deprecation_status: pluginData.deprecation_status,
        version: cleanedVersion,
        isAvailable: true,
        isExternalPlugin: false, // Static engines are not external
        pluginData,
      };
    }

    // DEMO: Mark certain plugins as disabled for testing ONLY if they're not in the real catalog
    if (demoDisabledPlugins.includes(engine.type)) {
      return {
        ...engine,
        isAvailable: false,
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

      // Look for a static engine type that this plugin name contains or matches
      // This handles cases like "my-custom-aws-plugin" matching the "aws" static engine
      const matchingStaticEngine = allEngines.find((engine) => {
        // Direct type match (e.g., secretEngineName = "aws" matches engine.type = "aws")
        if (secretEngineName === engine.type) {
          return true;
        }
        // Pattern match (e.g., secretEngineName = "my-custom-aws-plugin" contains "aws")
        return (
          secretEngineName.includes(engine.type) || secretEngineName.includes(engine.type.replace('-', ''))
        );
      });

      // Create dynamic engine metadata with defaults
      const dynamicEngine: EnhancedEngineDisplayData = {
        type: secretEngineName,
        displayName: secretEngineName
          .split('-')
          .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(' '), // Convert kebab-case to Title Case
        mountCategory: ['secret'],
        pluginCategory: 'external', // Mark as external since it's not in static metadata
        glyph: matchingStaticEngine?.glyph || 'file-text', // Use glyph from matching type or default
        isAvailable: true,
        // Use detailed info if available, otherwise create minimal plugin data
        builtin: detailedInfo?.builtin ?? false,
        deprecation_status: detailedInfo?.deprecation_status,
        version: detailedInfo?.version?.replace(/\+builtin.*$/, '') || detailedInfo?.version,
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
        glyph: matchingStaticEngine?.glyph || 'file-text', // Use glyph from matching type or default
        isAvailable: true,
        builtin: plugin.builtin,
        deprecation_status: plugin.deprecation_status,
        version: plugin.version?.replace(/\+builtin.*$/, '') || plugin.version,
        pluginData: plugin,
        isExternalPlugin: true,
      };

      dynamicPlugins.push(dynamicEngine);
    });
  }

  // Only add demo plugins if we have engines with data (not in tests with empty arrays)
  const shouldAddDemoPlugins =
    secretEnginesList.length > 0 ||
    secretEnginesDetailed.length > 0 ||
    databasePluginsList.length > 0 ||
    databasePluginsDetailed.length > 0;

  if (shouldAddDemoPlugins) {
    return [...enhancedEngines, ...dynamicPlugins, ...demoFakePlugins];
  }

  return [...enhancedEngines, ...dynamicPlugins];
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
