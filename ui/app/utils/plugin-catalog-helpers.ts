/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isEmpty } from '@ember/utils';
import type { EngineDisplayData } from './all-engines-metadata';
import type { PluginCatalogPlugin } from 'vault/services/plugin-catalog';

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
