/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ALL_ENGINES, type EngineDisplayData } from 'core/utils/all-engines-metadata';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';

/**
 * Default metadata for unknown engine plugins
 */
export const unknownEngineMetadata = (methodType?: string): EngineDisplayData => ({
  type: methodType || 'unknown',
  displayName: methodType || 'Unknown plugin',
  glyph: 'lock',
  mountCategory: ['secret', 'auth'],
});

/**
 * Helper function to retrieve engine metadata for a given `methodType`.
 * It searches the `ALL_ENGINES` array for an engine with a matching type and returns its metadata object.
 * The `ALL_ENGINES` array includes secret and auth engines, including those supported only in enterprise.
 * These details (such as mount type and enterprise licensing) are included in the returned engine object.
 *
 * For external plugins that have a builtin mapping (e.g., "vault-plugin-secrets-keymgmt" -> "keymgmt"),
 * this function returns the metadata for the corresponding builtin engine, preserving the original
 * external plugin name in the type field.
 *
 * Example usage:
 * const engineMetadata = engineDisplayData('kmip');
 * if (engineMetadata?.requiresEnterprise) {
 *   console.log(`This mount: ${engineMetadata.engineType} requires an enterprise license`);
 * }
 *
 * @param {string} methodType - The engine type (sometimes called backend) to look up (e.g., "aws", "azure", "vault-plugin-secrets-keymgmt").
 * @returns {Object} - The engine metadata, which includes information about its mount type (e.g., secret or auth)
 *   and whether it requires an enterprise license. For unknown engines, returns a default unknown plugin object.
 */
export default function engineDisplayData(methodType: string): EngineDisplayData {
  // First try to find an exact match
  const builtinEngine = ALL_ENGINES?.find((t) => t.type === methodType);
  if (builtinEngine) {
    return builtinEngine;
  }

  // If no direct match, check if this is a known external plugin and use its builtin mapping
  const effectiveType = getEffectiveEngineType(methodType);
  if (effectiveType !== methodType) {
    // This is a known external plugin with a builtin mapping
    const mappedEngine = ALL_ENGINES?.find((t) => t.type === effectiveType);
    if (mappedEngine) {
      // Return the mapped engine metadata but preserve the original external plugin type
      return {
        ...mappedEngine,
        type: methodType, // Keep the original external plugin name for identification
      };
    }
  }

  // Return default unknown plugin metadata
  return unknownEngineMetadata(methodType);
}
