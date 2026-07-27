/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { getPreference, setPreference } from 'vault/utils/preferences';

/**
 * UserPreferences::DataPrivacy
 *
 * Presentational "Data & Privacy" section of the User Preferences page. Renders
 * the "Share usage metrics" telemetry-consent control, persisting the user's
 * choice to localStorage via the central preference registry
 * (`vault:prefs:telemetryConsent`).
 *
 * Consent fails safe: when no value is stored the toggle reads off (opt-in).
 * This round is presentation + persistence only — no telemetry SDK is wired to
 * the stored value.
 */
export default class UserPreferencesDataPrivacy extends Component {
  // Initialize from storage; absent key resolves to the registry default (off).
  @tracked telemetryConsent = getPreference('telemetryConsent');

  // Items the anonymous telemetry would include / never include. Presentational.
  included = [
    'Feature usage patterns',
    'Navigation flows and page visits',
    'UI interaction events (clicks, form interactions)',
  ];

  excluded = [
    'Secret values or credentials',
    'Namespace paths or secret keys',
    'Auth tokens or identity data',
  ];

  @action
  updateConsent(event: Event) {
    const { checked } = event.target as HTMLInputElement;
    this.telemetryConsent = checked;
    // Persist on change through the registry — no Save/Cancel.
    setPreference('telemetryConsent', checked);
  }
}
