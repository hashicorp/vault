/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { toLabel } from 'core/helpers/to-label';

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
        ];
      case 'okta':
        return ['username_format', 'mount_accessor', 'org_name', 'api_token', 'base_url', 'primary_email'];
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
        ];
      default:
        return [];
    }
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
