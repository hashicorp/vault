/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { toLabel } from 'core/helpers/to-label';
import { isPresent } from '@ember/utils';

// Restricted-tier write-only credential names that must never appear in the
// UI read model or the DOM. The corresponding attribute definitions on the
// MfaMethod model are intentionally preserved so the create/update write
// form continues to function correctly.
const WRITE_ONLY_FIELD_NAMES = new Set([
  'secret_key',
  'integration_key',
  'api_token',
  'settings_file_base64',
]);

// Maps each MFA type to the subset of WRITE_ONLY_FIELD_NAMES that were
// formerly listed in displayFields for that type. Used to build per-type
// indicator rows (label + configured boolean) so the detail template can
// show presence status without the secret string.
const WRITE_ONLY_FIELDS_BY_TYPE = {
  duo: ['secret_key', 'integration_key'],
  okta: ['api_token'],
  pingid: ['settings_file_base64'],
  totp: [],
};

/**
 * Returns true only when a credential value is present AND non-blank.
 * isPresent from @ember/utils returns true for whitespace-only strings;
 * the additional trim() check ensures those are also treated as not
 * configured, matching the edge case requirement.
 */
function isConfigured(value) {
  return isPresent(value) && (typeof value !== 'string' || value.trim().length > 0);
}

export default class MfaMethodController extends Controller {
  @service router;
  @service flashMessages;
  @service api;

  queryParams = ['tab'];
  tab = 'config';

  get displayFields() {
    const type = this.model.method.type;
    switch (type) {
      case 'duo':
        return [
          'username_format',
          'secret_key',
          'integration_key',
          'api_hostname',
          'push_info',
          'use_passcode',
        ].filter((f) => !WRITE_ONLY_FIELD_NAMES.has(f));
      case 'okta':
        return [
          'username_format',
          'mount_accessor',
          'org_name',
          'api_token',
          'base_url',
          'primary_email',
        ].filter((f) => !WRITE_ONLY_FIELD_NAMES.has(f));
      case 'totp':
        return [
          'issuer',
          'period',
          'key_size',
          'qr_size',
          'algorithm',
          'digits',
          'skew',
          'max_validation_attempts',
          'enable_self_enrollment',
        ];
      case 'pingid':
        return [
          'username_format',
          'settings_file_base64',
          'use_signature',
          'idp_url',
          'admin_url',
          'authenticator_url',
          'org_alias',
        ].filter((f) => !WRITE_ONLY_FIELD_NAMES.has(f));
      default:
        return [];
    }
  }

  /**
   * An in-memory copy of the method read payload for use by the detail
   * template. Restricted write-only credential strings are absent from
   * this object — each write-only field is represented by a presence
   * indicator `{ configured: boolean }` with no secret property, so
   * those strings cannot be bound into the DOM.
   *
   * The original ember-data record is NOT mutated. The write form retains
   * full access to all attributes via this.model.method.
   */
  get sanitizedModel() {
    const method = this.model?.method;
    if (!method) return null;

    const result = {};

    // Copy non-secret display fields from the source record.
    for (const field of this.displayFields) {
      result[field] = method[field];
    }

    // Add configured-only presence flags for write-only fields relevant to
    // this method type. TOTP yields an empty array — no flags are added.
    // A TOTP record that unexpectedly carries secret_key or
    // settings_file_base64 will therefore not expose those values here.
    const writeOnlyForType = WRITE_ONLY_FIELDS_BY_TYPE[method.type] ?? [];
    for (const field of writeOnlyForType) {
      result[field] = { configured: isConfigured(method[field]) };
    }

    return result;
  }

  /**
   * Array of `{ label, configured }` indicator rows for the detail template
   * to display in place of each Restricted credential field. Contains no
   * secret property and no secret string — only a human-readable label and
   * a boolean presence flag.
   */
  get writeOnlyIndicators() {
    const method = this.model?.method;
    if (!method) return [];

    const writeOnlyForType = WRITE_ONLY_FIELDS_BY_TYPE[method.type] ?? [];
    return writeOnlyForType.map((field) => ({
      label: this.label(field),
      configured: isConfigured(method[field]),
    }));
  }

  label(field) {
    const key = field.replace('config.', '');
    const label = toLabel([key]);
    // map specific fields to custom labels
    return (
      {
        push_info: 'Duo push information',
        use_passcode: 'Passcode reminder',
        org_name: 'Organization name',
        api_token: 'Okta API key',
        base_url: 'Base URL',
        qr_size: 'QR size',
        enable_self_enrollment: 'Enable self-enrollment',
        api_hostname: 'API hostname',
      }[key] || label
    );
  }

  @action
  async deleteMethod() {
    try {
      const { type, id } = this.model.method;
      switch (type) {
        case 'totp':
          await this.api.identity.mfaDeleteTotpMethod(id);
          break;
        case 'pingid':
          await this.api.identity.mfaDeletePingIdMethod(id);
          break;
        case 'duo':
          await this.api.identity.mfaDeleteDuoMethod(id);
          break;
        case 'okta':
          await this.api.identity.mfaDeleteOktaMethod(id);
          break;
        default:
          throw new Error(`Unknown MFA method type: ${type}`);
      }
      this.flashMessages.success('MFA method deleted successfully.');
      this.router.transitionTo('vault.cluster.access.mfa.methods');
    } catch (error) {
      this.flashMessages.danger('There was an error deleting this MFA method.');
    }
  }
}
