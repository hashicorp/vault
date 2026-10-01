/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

import AnalyticsService from 'vault/services/analytics';
import { ThemeChoice } from 'vault/services/theme';
import type ThemeService from 'vault/services/theme';

import { PREFERENCES_THEME_SET } from 'vault/utils/analytic-events';
export default class Appearance extends Component {
  @service declare readonly analytics: AnalyticsService;
  @service declare readonly theme: ThemeService;

  get currentTheme() {
    return this.theme.theme;
  }

  @action
  selectTheme(choice: ThemeChoice) {
    this.theme.setTheme(choice);
    this.analytics.trackEvent(PREFERENCES_THEME_SET, {
      namespace: 'preferences',
      action: 'theme_set',
      elementId: 'theme-select',
      channel: 'webpage',
      location: 'Appearance',
      objectType: 'theme',
      object: choice,
    });
  }
}
