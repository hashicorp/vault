/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import type Owner from '@ember/owner';
import Ember from 'ember';
import { getStringPreference, setStringPreference } from 'vault/utils/preferences';

const STORAGE_KEY = 'theme';
// The full localStorage key written by setStringPreference / local-storage.ts.
const STORAGE_FULL_KEY = `vault:prefs:${STORAGE_KEY}`;
export type ThemeChoice = 'dark' | 'light' | 'system';

export default class ThemeService extends Service {
  @tracked theme: ThemeChoice = this._resolveInitialTheme();

  // Held so we can remove the listener when the theme changes away from 'system'
  // or when the service is destroyed.
  private _mediaQuery: MediaQueryList | null = null;
  private _systemListener = () => this._applyTheme(true);

  // Bound storage listener for cross-tab theme sync. Stored so we can remove it
  // in willDestroy without creating a new function reference each time.
  private _storageListener = (event: StorageEvent) => {
    if (event.storageArea !== window.localStorage) return;
    if (event.key !== STORAGE_FULL_KEY) return;

    // newValue is JSON-encoded by local-storage.ts (e.g. `"\"dark\""`).
    // A null newValue means the key was removed; fall back to the default.
    let incoming: string;
    try {
      incoming = event.newValue !== null ? (JSON.parse(event.newValue) as string) : 'system';
    } catch {
      return; // unparseable — ignore
    }

    const isValidTheme = (v: string): v is ThemeChoice => v === 'dark' || v === 'light' || v === 'system';
    if (!isValidTheme(incoming)) return;

    this.theme = incoming;
    this._applyTheme(true);
    this._syncSystemListener();
  };

  constructor(owner: Owner) {
    super(owner);
    // Leave unanimated as there is no previous state to cross-fade from on boot.
    this._applyTheme();
    this._syncSystemListener();
    // Listen for theme changes written by other tabs on the same origin.
    window.addEventListener('storage', this._storageListener);
  }

  /** True when the effective (resolved) theme is dark. */
  get isDarkMode(): boolean {
    if (this.theme === 'system') {
      return window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false;
    }
    return this.theme === 'dark';
  }

  @action
  setTheme(choice: ThemeChoice): void {
    this.theme = choice;
    setStringPreference(STORAGE_KEY, choice);
    this._applyTheme(true);
    // Re-register because the 'system' branch may have changed.
    this._syncSystemListener();
  }

  willDestroy(): void {
    this._mediaQuery?.removeEventListener('change', this._systemListener);
    window.removeEventListener('storage', this._storageListener);
    super.willDestroy();
  }

  private _resolveInitialTheme(): ThemeChoice {
    return getStringPreference(STORAGE_KEY) as ThemeChoice;
  }

  /**
   * Writes `data-theme="dark"` when dark mode is active, removes it otherwise.
   *
   * When `animate` is true the write is wrapped in `document.startViewTransition`
   * so the browser cross-dissolves the before/after page snapshots (see
   * styles/theme/theme-fade.scss for timing).
   */
  private _applyTheme(animate = false): void {
    const write = () => {
      // Resolve the effective theme: for 'system', defer to the OS preference.
      if (this.isDarkMode) {
        document.documentElement.setAttribute('data-theme', 'dark');
      } else {
        // 'light' is the default; no attribute needed.
        document.documentElement.removeAttribute('data-theme');
      }
    };

    if (animate && !Ember.testing && typeof document.startViewTransition === 'function') {
      document.startViewTransition(write);
    } else {
      write();
    }
  }

  /**
   * Keeps the OS-level `prefers-color-scheme` listener in sync with the current
   * theme choice. When the theme is 'system' we subscribe so that OS dark/light
   * toggles update the page immediately without a reload. For explicit 'dark' or
   * 'light' choices the listener is removed — OS changes should have no effect.
   */
  private _syncSystemListener(): void {
    this._mediaQuery?.removeEventListener('change', this._systemListener);
    this._mediaQuery = null;

    if (this.theme === 'system') {
      this._mediaQuery = window.matchMedia?.('(prefers-color-scheme: dark)') ?? null;
      this._mediaQuery?.addEventListener('change', this._systemListener);
    }
  }
}
