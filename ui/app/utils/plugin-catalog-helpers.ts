/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isEmpty } from '@ember/utils';
import type { PluginCatalogPlugin } from 'vault/services/plugin-catalog';
import type { EngineDisplayData } from './all-engines-metadata';
import { getBuiltinTypeFromExternalPlugin, isKnownExternalPlugin } from './external-plugin-helpers';

/**
 * Constants for plugin catalog functionality
 */
const DEFAULT_EXTERNAL_PLUGIN_GLYPH = '';

/**
 * Plugin categories used throughout the application
 */
export const PLUGIN_CATEGORIES = {
  GENERIC: 'generic',
  CLOUD: 'cloud',
  INFRA: 'infra',
  EXTERNAL: 'external',
} as const;

/**
 * Mount categories for different engine types
 */
export const MOUNT_CATEGORIES = {
  SECRET: 'secret',
  AUTH: 'auth',
  DATABASE: 'database',
} as const;

/**
 * Plugin types used in catalog responses
 */
export const PLUGIN_TYPES = {
  SECRET: 'secret',
  AUTH: 'auth',
  DATABASE: 'database',
} as const;

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
 * Enhances static engine data with plugin catalog information including availability status,
 * deprecation status, and discovery of external plugins not in static metadata.
 *
 * @param allEngines - Array of static engine metadata
 * @param secretEnginesDetailed - Array of detailed secret engine info from catalog
 * @param databasePluginsDetailed - Array of detailed database plugin info from catalog
 * @returns Enhanced engines with catalog data and dynamically discovered external plugins
 */
export function enhanceEnginesWithCatalogData(
  allEngines: EngineDisplayData[],
  secretEnginesDetailed: PluginCatalogPlugin[] = [],
  databasePluginsDetailed: PluginCatalogPlugin[] = []
): EnhancedEngineDisplayData[] {
  if (isEmpty(secretEnginesDetailed) && isEmpty(databasePluginsDetailed)) {
    return allEngines;
  }

  // First, enhance existing static engines with catalog data
  const enhancedEngines = allEngines.map((engine) => {
    if (engine.type === MOUNT_CATEGORIES.DATABASE) {
      const isDatabaseAvailable = databasePluginsDetailed.length > 0;

      if (isDatabaseAvailable) {
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
    // Skip if this plugin already exists in static metadata or is a builtin plugin
    if (staticEngineTypes.has(plugin.name) || plugin.builtin) {
      return;
    }

    // Skip plugins that have known builtin mappings - these should appear in their
    // respective categories (e.g., KV, AWS) rather than in the "External" category
    if (isKnownExternalPlugin(plugin.name)) {
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

    // Only create external engines for custom external plugins (external plugins without mappings to builtin Vault plugins)
    const externalEngine: EnhancedEngineDisplayData = {
      type: plugin.name,
      displayName: plugin.name
        .split('-')
        .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' '), // Convert kebab-case to Title Case
      mountCategory: [MOUNT_CATEGORIES.SECRET],
      pluginCategory: PLUGIN_CATEGORIES.EXTERNAL, // Mark as external since it's not in static metadata
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
 * Categorizes engines by their availability status for display purposes.
 * Separate enabled and disabled plugins in the UI.
 * Engines with isAvailable === false are considered disabled.
 *
 * @param engines - Array of enhanced engine data with availability information
 * @returns Object containing separate arrays for enabled and disabled engines
 *
 * @example
 * ```typescript
 * const engines = enhanceEnginesWithCatalogData(staticEngines, catalogData);
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

export function getPluginVersionsFromEngineType(list: PluginCatalogPlugin[] | undefined, name: string) {
  if (!list) return [];

  return list.reduce((acc: string[], item: PluginCatalogPlugin) => {
    if (item.name === name) acc.push(item.version);
    return acc;
  }, []);
}

/**
 * Version information for a specific plugin engine
 */
export interface EngineVersionInfo {
  version: string;
  pluginName: string;
  isBuiltin: boolean;
}

/**
 * Result containing version information and unversioned plugin detection
 */
export interface EngineVersionResult {
  versions: EngineVersionInfo[];
  hasUnversionedPlugins: boolean;
}

/**
 * Retrieves all available plugin versions for a specific engine type from the catalog.
 * This enables users to choose between builtin and external plugin variants when mounting
 * secrets engines, supporting both standard Vault engines and custom external plugins.
 *
 * The function handles the mapping between external plugin names (e.g., "vault-plugin-secrets-kv")
 * and their corresponding engine types (e.g., "kv") to provide a unified version selection experience.
 *
 * @param secretEnginesDetailed - Array of detailed secret engine info from catalog API
 * @param engineType - The engine type to get versions for (e.g., 'kv', 'aws')
 * @param pluginType - Optional plugin type filter ('secret', 'auth', 'database')
 * @returns Object containing version information array and flag for unversioned plugins
 */
export function getAllVersionsForEngineType(
  secretEnginesDetailed: PluginCatalogPlugin[] | undefined,
  engineType: string,
  pluginType = 'secret'
): EngineVersionResult {
  if (
    !engineType ||
    !secretEnginesDetailed ||
    typeof engineType !== 'string' ||
    !Array.isArray(secretEnginesDetailed)
  ) {
    return { versions: [], hasUnversionedPlugins: false };
  }

  let hasUnversionedPlugins = false;
  const filteredVersions: EngineVersionInfo[] = [];

  secretEnginesDetailed.forEach((plugin) => {
    // Basic validation
    if (!plugin?.name || typeof plugin?.builtin !== 'boolean' || typeof plugin?.version !== 'string') {
      return;
    }

    // Filter by plugin type (secret, auth, database)
    if (plugin.type !== pluginType) {
      return;
    }

    // Check if this plugin matches the engine type
    const isDirectMatch = plugin.name === engineType;
    const builtin = getBuiltinTypeFromExternalPlugin(plugin.name);
    const isExternalMatch = builtin === engineType;

    if (!isDirectMatch && !isExternalMatch) {
      return;
    }

    // Check for unversioned plugins (empty version strings)
    if (plugin.version === '') {
      hasUnversionedPlugins = true;
      return; // Don't include in versions array
    }

    // Include versioned plugins
    filteredVersions.push({
      version: plugin.version,
      pluginName: plugin.name,
      isBuiltin: plugin.builtin,
    });
  });

  return { versions: filteredVersions, hasUnversionedPlugins };
}
