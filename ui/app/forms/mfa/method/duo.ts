/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { MfaCreateDuoMethodRequest } from '@hashicorp/vault-client-typescript';

type MfaCreateDuoMethodFormData = MfaCreateDuoMethodRequest & {
  name: string;
};

export default class MfaCreateDuoMethodForm extends Form<MfaCreateDuoMethodFormData> {
  type = 'duo';

  formFields = [
    new FormField('username_format', 'string', {
      label: 'Username format',
      subText: 'How to map identity names to MFA method names. ',
    }),
    new FormField('secret_key', 'string', {
      label: 'Duo secret key',
      sensitive: true,
    }),
    new FormField('integration_key', 'string', {
      label: 'Duo integration key',
      sensitive: true,
    }),
    new FormField('api_hostname', 'string', {
      label: 'Duo API hostname',
    }),
    new FormField('push_info', 'string', {
      label: 'Duo push information',
      subText: 'Additional information displayed to the user when the push is presented to them.',
    }),
    new FormField('use_passcode', 'boolean', {
      label: 'Passcode reminder',
      subText: 'If this is turned on, the user is reminded to use the passcode upon MFA validation.',
    }),
  ];

  validations: Validations = {
    secret_key: [{ type: 'presence', message: 'Secret key is required.' }],
    integration_key: [{ type: 'presence', message: 'Integration key is required.' }],
    api_hostname: [{ type: 'presence', message: 'API hostname is required.' }],
  };
}
