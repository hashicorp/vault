/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isAddonEngine } from 'vault/utils/all-engines-metadata';
import engineDisplayData from 'vault/helpers/engines-display-data';

/**
 * Get the appropriate route for exiting configuration based on engine type and version.
 * This handles the logic for determining whether to use the backends route for
 * isOnlyMountable engines, or the engine-specific routes for other engines.
 *
 * @param engineType - The type of the engine
 * @param version - The version of the engine (relevant for KV engines)
 * @returns The full route path for the exit configuration button
 */
function getExitConfigurationRoute(engineType: string, version?: number): string {
  const engineData = engineDisplayData(engineType);

  if (engineData.isOnlyMountable) {
    return 'vault.cluster.secrets.backends';
  }

  const baseRoute = 'vault.cluster.secrets.backend';
  const shouldUseEngineRoute = isAddonEngine(engineType, version || 1);

  if (shouldUseEngineRoute && engineData.engineRoute) {
    return `${baseRoute}.${engineData.engineRoute}`;
  }

  return `${baseRoute}.list-root`;
}

/**
 * Handlebars helper to get the appropriate exit configuration route for a secrets engine.
 * This helper handles all the logic for determining the correct route based on the engine type and version.
 *
 * Usage:
 * @route={{exit-configuration-route engineType version}}
 *
 * @param engineType - The type of the secrets engine
 * @param version - The version of the engine (optional, defaults to 1)
 * @returns The full route path for the exit configuration button
 */
export default function exitConfigurationRoute(engineType: string, version?: number): string {
  return getExitConfigurationRoute(engineType, version);
}
