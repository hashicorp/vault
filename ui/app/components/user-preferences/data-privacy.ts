/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { getPreference } from 'vault/utils/preferences';

import type AnalyticsService from 'vault/services/analytics';

/**
 * UserPreferences::DataPrivacy
 *
 * "Data & Privacy" section of the User Preferences page. Renders the "Share
 * usage metrics" telemetry-consent control and routes changes through the
 * analytics service's recordConsent, which persists the choice and applies it
 * to the current session (starts on accept, stops on decline) without a reload.
 */
export default class UserPreferencesDataPrivacy extends Component {
  @service declare readonly analytics: AnalyticsService;

  // Initialize from storage; absent key resolves to the registry default (off).
  @tracked telemetryConsent = getPreference('telemetryConsent');

  // Items the anonymous telemetry would include / never include. Presentational.
  included = ['Feature use', 'Navigation and page visits', 'Interaction events (clicks, configurations)'];

  excluded = [
    'No secrets, certificates or keys',
    'No namespace or secret paths',
    'No auth tokens or identity data',
  ];

  @action
  updateConsent(event: Event) {
    const { checked } = event.target as HTMLInputElement;
    this.telemetryConsent = checked;
    // Persist and apply to the current session (start/stop) — no Save/Cancel.
    this.analytics.recordConsent(checked);
  }
}
