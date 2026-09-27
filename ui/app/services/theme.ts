/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import type Owner from '@ember/owner';
import Ember from 'ember';

const STORAGE_KEY = 'vault:theme';

export type ThemeChoice = 'dark' | 'light' | 'system';

export default class ThemeService extends Service {
  @tracked theme: ThemeChoice = this._resolveInitialTheme();

  constructor(owner: Owner) {
    super(owner);
    // Leave unanimated as there is no previous state to cross-fade from on boot.
    this._applyTheme();
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
    window.localStorage.setItem(STORAGE_KEY, choice);
    this._applyTheme(true);
  }

  private _resolveInitialTheme(): ThemeChoice {
    const stored = window.localStorage.getItem(STORAGE_KEY);
    if (stored === 'dark' || stored === 'light' || stored === 'system') {
      return stored;
    }
    return 'system';
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
}
