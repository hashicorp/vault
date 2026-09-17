/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { debounce } from '@ember/runloop';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { USER_PREFERENCES_PERSONA_SET } from 'vault/utils/analytic-events';
import { getPreference, getStringPreference, setStringPreference } from 'vault/utils/preferences';

import type AnalyticsService from 'vault/services/analytics';

/**
 * UserPreferences::Persona
 *
 * "Your role" section of the User Preferences page. Renders an HDS Select field
 * that lets users self-identify their role. The selection is persisted in
 * localStorage via the string preferences registry. When the user has already
 * opted into telemetry, the selection is also forwarded to Segment as a
 * `UI Interaction` event (`USER_PREFERENCES_PERSONA_SET`) so it can inform future UX decisions.
 *
 * When "Other" is selected, a free-text field appears. The typed value is
 * stored alongside the selection (e.g. `{ value: 'other', customRole: 'DevOps / SRE' }`)
 * so the raw label is preserved without lossy transformation. When sent to Segment,
 * the value is formatted as `other` or `other-<kebab-role>`.
 */

export interface PersonaPreference {
  value: string;
  customRole?: string;
}

export default class UserPreferencesPersona extends Component {
  @service declare readonly analytics: AnalyticsService;

  // Persona option values — the order matches the dropdown display order.
  personas = [
    { value: 'developer', label: 'Developer' },
    { value: 'platform-engineer', label: 'Platform Engineer / Vault Administrator' },
    { value: 'security-analyst', label: 'Security Analyst' },
    { value: 'other', label: 'Other' },
  ];

  // Parsed initial state from localStorage. Supports backward-compatible string or object structure.
  @tracked selectedValue = this.initialPersona.value;
  @tracked otherRawText = this.initialPersona.customRole ?? '';
  @tracked otherInputError = '';

  // Only letters, digits, and spaces are allowed in the free-text custom role field.
  private readonly customRolePattern = /^[a-zA-Z0-9 ]*$/;

  private get initialPersona(): PersonaPreference {
    const raw = getStringPreference('persona');
    if (!raw) return { value: '' };

    try {
      const parsed = JSON.parse(raw);
      if (typeof parsed === 'object' && parsed !== null && 'value' in parsed) {
        return parsed as PersonaPreference;
      }
    } catch {
      // Stored value is not valid JSON; return the empty default.
    }

    return { value: '' };
  }

  /** True when the selected persona value is 'other'. */
  get isOther(): boolean {
    return this.selectedValue === 'other';
  }

  /** Normalises arbitrary text to lowercase kebab-case for analytics events. */
  private toKebab(value: string): string {
    return value
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }

  @action
  updatePersona(event: Event) {
    const { value } = event.target as HTMLSelectElement;
    this.selectedValue = value;
    this.otherRawText = '';
    const pref: PersonaPreference = { value };
    this.persistAndTrack(pref);
  }

  @action
  updateOtherText(event: Event) {
    this.otherRawText = (event.target as HTMLInputElement).value;
    this.otherInputError = this.customRolePattern.test(this.otherRawText)
      ? ''
      : 'Role description can only contain letters, numbers, and spaces.';

    if (this.otherInputError) return;

    const pref: PersonaPreference = {
      value: 'other',
      customRole: this.otherRawText,
    };
    debounce(this, this.persistAndTrack, pref, 300);
  }

  persistAndTrack(pref: PersonaPreference): void {
    // Omit customRole from the serialised object when it is empty so the stored
    // value stays clean: { value: "other" } rather than { value: "other", customRole: "" }.
    const toStore: PersonaPreference = { value: pref.value };
    if (pref.customRole) toStore.customRole = pref.customRole;
    setStringPreference('persona', JSON.stringify(toStore));

    if (getPreference('telemetryConsent')) {
      let eventValue = pref.value;
      if (pref.value === 'other' && pref.customRole) {
        const kebab = this.toKebab(pref.customRole);
        eventValue = kebab ? `other-${kebab}` : 'other';
      }
      this.analytics.trackEvent(USER_PREFERENCES_PERSONA_SET, {
        namespace: 'user-preferences',
        action: 'persona_set',
        elementId: 'persona-select',
        channel: 'webpage',
        location: 'user-preferences',
        objectType: 'persona',
        object: eventValue,
        resultValue: eventValue,
      });
    }
  }
}
