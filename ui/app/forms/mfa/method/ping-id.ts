/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { MfaCreatePingIdMethodRequest } from '@hashicorp/vault-client-typescript';

type MfaCreatePingIdMethodFormData = MfaCreatePingIdMethodRequest & {
  name: string;
};

export default class MfaCreatePingIdMethodForm extends Form<MfaCreatePingIdMethodFormData> {
  type = 'pingid';

  formFields = [
    new FormField('username_format', 'string', {
      label: 'Username format',
      subText: 'How to map identity names to MFA method names. ',
    }),
    new FormField('settings_file_base64', 'string', {
      label: 'Settings file',
      subText: 'A base-64 encoded third party setting file retrieved from the PingIDs configuration page.',
    }),
    new FormField('use_signature', 'boolean'),
    new FormField('idp_url', 'string'),
    new FormField('admin_url', 'string'),
    new FormField('authenticator_url', 'string'),
    new FormField('org_alias', 'string'),
  ];

  validations: Validations = {
    settings_file_base64: [{ type: 'presence', message: 'Settings file base64 is required.' }],
  };
}
