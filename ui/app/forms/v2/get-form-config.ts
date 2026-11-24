/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import GENERATED_CONFIGS from './generated/index';
import OVERRIDE_CONFIGS from './overrides/index';
import type { FormConfig } from './form-config';

export type FormConfigKey = keyof typeof GENERATED_CONFIGS | keyof typeof OVERRIDE_CONFIGS;

type AllConfigs = typeof GENERATED_CONFIGS & typeof OVERRIDE_CONFIGS;

/**
 * Extract the payload type for a given form config key
 */
export type ExtractPayload<K extends FormConfigKey> = AllConfigs[K] extends FormConfig<
  infer TPayload,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  infer _TResponse
>
  ? TPayload
  : never;

/**
 * Extract the response type for a given form config key
 */
export type ExtractResponse<K extends FormConfigKey> = AllConfigs[K] extends FormConfig<
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  infer _TPayload,
  infer TResponse
>
  ? TResponse
  : never;

/**
 * Retrieves a V2 form configuration, preferring overrides when they exist.
 * @param configName - The camelCase config name (e.g., 'aliCloudDeleteAuthRole', 'azureConfigureAuth')
 * @returns V2FormConfig instance with the exact type for the given config
 */
export function getFormConfig<K extends FormConfigKey>(
  configName: K
): FormConfig<ExtractPayload<K>, ExtractResponse<K>> {
  // Check overrides first
  const override = OVERRIDE_CONFIGS[configName as keyof typeof OVERRIDE_CONFIGS];
  if (override) {
    return override as unknown as FormConfig<ExtractPayload<K>, ExtractResponse<K>>;
  }

  // Fall back to generated configs
  const config = GENERATED_CONFIGS[configName as keyof typeof GENERATED_CONFIGS];

  if (!config) {
    throw new Error(`Form configuration not found for: ${configName}`);
  }

  return config as unknown as FormConfig<ExtractPayload<K>, ExtractResponse<K>>;
}
