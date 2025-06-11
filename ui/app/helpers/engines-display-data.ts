/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

/**
 * Helper function to retrieve engine metadata for a given `methodType`.
 * It searches the `ALL_ENGINES` array for an engine with a matching type and returns its metadata object.
 * The `ALL_ENGINES` array includes secret and auth engines, including those supported only in enterprise.
 * These details (such as mount type and enterprise licensing) are included in the returned engine object.
 *
 * Example usage:
 * const engineMetadata = engineDisplayData('kmip');
 * if (engineMetadata?.requiresEnterprise) {
 *   console.log(`This mount: ${engineMetadata.engineType} requires an enterprise license`);
 * }
 *
 * @param {string} methodType - The engine type (sometimes called backend) to look up (e.g., "aws", "azure").
 * @returns {Object|undefined} - The engine metadata, which includes information about its mount type (e.g., secret or auth)
 *   and whether it requires an enterprise license. Returns undefined if no match is found.
 */
export default function engineDisplayData(methodType: string) {
  const engine = ALL_ENGINES?.find((t) => t.type === methodType);
  return engine;
}
