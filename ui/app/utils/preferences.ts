/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import localStorage from 'vault/lib/local-storage';

/**
 * Central registry for Vault UI user preferences.
 *
 * This module is the single source of truth for every preference key name,
 * its type, and its default value. Consuming components MUST reference this
 * registry rather than using ad-hoc key string literals.
 *
 * Keys follow the namespaced convention `vault:prefs:<name>`. Reads and writes
 * go through the existing local-storage module (`vault/lib/local-storage`);
 * there is no separate storage abstraction layer. Reading a key that is absent
 * from localStorage returns the documented default defined here.
 */

export const NAMESPACE = 'vault:prefs';

interface PreferenceDefinition {
  key: string;
  type: 'boolean';
  default: boolean;
}

export const PREFERENCES: Record<string, PreferenceDefinition> = {
  // Opt-in consent for anonymous usage telemetry. Fails safe: defaults off
  // until the user explicitly opts in. Per-browser (localStorage only).
  telemetryConsent: {
    key: `${NAMESPACE}:telemetryConsent`,
    type: 'boolean',
    default: false,
  },
};

export type PreferenceName = keyof typeof PREFERENCES;

function resolve(name: PreferenceName): PreferenceDefinition {
  const def = PREFERENCES[name];
  if (!def) {
    throw new Error(`[preferences] Unknown preference "${name}". Register it in app/utils/preferences.ts.`);
  }
  return def;
}

export function getPreference(name: PreferenceName): boolean {
  const def = resolve(name);
  try {
    const stored = localStorage.getItem(def.key);
    return stored === null || stored === undefined ? def.default : stored;
  } catch {
    return def.default;
  }
}

/**
 * Whether a value has been explicitly stored for this preference, as opposed to
 * falling back to its default. Consent gating needs this to tell "user has never
 * decided" (show the banner) apart from "user explicitly declined" (stay quiet).
 * Unreadable storage is treated as "not stored" (safe default).
 */
export function hasPreference(name: PreferenceName): boolean {
  const def = resolve(name);
  try {
    const stored = localStorage.getItem(def.key);
    return stored !== null && stored !== undefined;
  } catch {
    return false;
  }
}

export function setPreference(name: PreferenceName, value: boolean): void {
  const { key } = resolve(name);
  try {
    localStorage.setItem(key, value);
  } catch {
    // Cannot persist (storage unavailable). Hence a no-op
  }
}
