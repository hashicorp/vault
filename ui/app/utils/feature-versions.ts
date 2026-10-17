/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Registry of version-gated features.
 *
 * Each entry maps a stable feature key to the minimum Vault binary version
 * required for that feature to be available in the UI.
 *
 * CONVENTIONS
 * -----------
 * - `key`     Stable kebab-case identifier. Never rename a key — existing call
 *             sites reference it by string. Add a new entry instead.
 * - `name`    Human-readable display name used in error messages and UI copy,
 *             e.g. "This feature requires Vault ≥ X.Y.Z".
 * - `version` Minimum Vault version string (no prefix/suffix), e.g. "2.1.0".
 *             Update this field here if a feature is backported — all call sites
 *             pick up the change automatically.
 *
 * ADDING A NEW VERSION-GATED FEATURE
 * -----------------------------------
 * 1. Add an entry to FEATURE_VERSIONS below.
 * 2. Gate the feature in your component/route with `this.version.hasFeature('your-key')`.
 * 3. Inside the gate, call new API methods via `this.api.versioned().sys.newMethod()`.
 *
 * See ui/docs/versioned-api-calls.md for the full guide.
 */

export interface VersionedFeature {
  /** Stable kebab-case identifier used at call sites. */
  key: string;
  /** Human-readable name for UI copy and error messages. */
  name: string;
  /** Minimum Vault version required (semver, no v prefix or +ent suffix). */
  version: string;
}

/**
 * The canonical registry of version-gated features.
 * Add one entry here per feature that requires a minimum Vault version.
 */
export const FEATURE_VERSIONS: VersionedFeature[] = [
  // Example — replace with real entries as version-gated features are added:
  // { key: 'agents', name: 'Agentic Security', version: '2.1.0' },
];

/**
 * Look up a feature entry by its key.
 * Returns `undefined` if the key is not registered in FEATURE_VERSIONS.
 */
export function getFeatureVersion(key: string): VersionedFeature | undefined {
  return FEATURE_VERSIONS.find((f) => f.key === key);
}
