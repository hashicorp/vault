/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Ember Data model type mapping utilities for secret engines.
 *
 * This file contains functions that determine the appropriate Ember model type
 * based on engine type and context. These utilities are specifically related to
 * Ember Data model management.
 *
 * TODO: Migrate to API service instead of Ember Data model types.
 * When routes are converted to TypeScript, the string-based model type approach
 * becomes problematic for type safety. Direct API service calls would provide
 * better type safety and eliminate the need for model type mapping utilities.
 * This would align with the eventual migration away from Ember Data.
 */

/**
 * Helper function to determine the model type from a secret path for transform engine.
 * @param secret - The secret path to analyze
 * @returns The model type based on the secret path prefix, or 'transform' if no recognized prefix
 */
function getTransformModelTypeFromSecretPath(secret: string): string {
  switch (true) {
    case secret.startsWith('role/'):
      return 'transform/role';
    case secret.startsWith('template/'):
      return 'transform/template';
    case secret.startsWith('alphabet/'):
      return 'transform/alphabet';
    default:
      return 'transform';
  }
}

/**
 * Helper function to determine the model type from query parameters for transform engine.
 * @param transformType - The transform type from context (transformType or tab)
 * @returns The model type based on the transform type, or 'transform' if no match
 */
function getTransformModelTypeFromParams(transformType?: string): string {
  const validTypes = ['role', 'template', 'alphabet'];
  if (transformType && validTypes.includes(transformType)) {
    return `transform/${transformType}`;
  }
  return 'transform';
}

/**
 * Main helper function to determine the transform model type based on context.
 * @param context - Context object containing secret path, transformType, or tab
 * @returns The appropriate transform model type
 */
function getTransformModelType(context: { transformType?: string; tab?: string; secret?: string }): string {
  // Check secret name prefix first (for existing secrets)
  if (context.secret) {
    const secretBasedType = getTransformModelTypeFromSecretPath(context.secret);
    // If secret has a recognized prefix, use it. Otherwise, fall back to tab/transformType
    if (secretBasedType !== 'transform') {
      return secretBasedType;
    }
  }

  // Fall back to query parameters (for new secrets or navigation, or when secret has no recognized prefix)
  const transformType = context.transformType || context.tab;
  return getTransformModelTypeFromParams(transformType);
}

/**
 * Engine type to Ember model type mapping for secrets engines.
 * Used by routes to determine the correct Ember model type for a given engine.
 */
const ENGINE_TYPE_TO_MODEL_TYPE_MAP = {
  database: (context: { isRole?: boolean; tab?: string; secret?: string }) => {
    if (context.isRole || context.tab === 'role' || context.secret?.startsWith('role/')) {
      return 'database/role';
    }
    return 'database/connection';
  },
  transit: () => 'transit-key',
  ssh: () => 'role-ssh',
  aws: () => 'role-aws',
  cubbyhole: () => 'secret',
  kv: () => 'secret',
  keymgmt: (context: { tab?: string; itemType?: string }) =>
    `keymgmt/${context.itemType || context.tab || 'key'}`,
  transform: getTransformModelType,
  generic: () => 'secret',
  totp: () => 'totp-key',
} as const;

/**
 * Get the appropriate Ember model type for a given effective engine type and context.
 *
 * @param effectiveEngineType - The effective engine type (after external plugin mapping)
 * @param context - Context object with additional parameters needed for some engines
 * @returns The Ember model type string
 */
export function getModelTypeForEngine(
  effectiveEngineType: string,
  context: {
    tab?: string;
    itemType?: string;
    secret?: string;
    isRole?: boolean;
    transformType?: string;
  } = {}
): string {
  const modelTypeFn =
    ENGINE_TYPE_TO_MODEL_TYPE_MAP[effectiveEngineType as keyof typeof ENGINE_TYPE_TO_MODEL_TYPE_MAP];
  return modelTypeFn ? modelTypeFn(context) : 'secret';
}
