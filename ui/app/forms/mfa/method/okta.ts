/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { MfaCreateOktaMethodRequest } from '@hashicorp/vault-client-typescript';

type MfaCreateOktaMethodFormData = MfaCreateOktaMethodRequest & {
  name: string;
};

export default class MfaCreateOktaMethodForm extends Form<MfaCreateOktaMethodFormData> {
  type = 'okta';

  formFields = [
    new FormField('username_format', 'string', {
      label: 'Username format',
      subText: 'How to map identity names to MFA method names. ',
    }),
    new FormField('mount_accessor', 'string'),
    new FormField('org_name', 'string', {
      label: 'Organization name',
      subText: 'Name of the organization to be used in the Okta API.',
    }),
    new FormField('api_token', 'string', {
      label: 'Okta API key',
    }),
    new FormField('base_url', 'string', {
      label: 'Base URL',
      subText:
        'If set, will be used as the base domain for API requests. Example are okta.com, oktapreview.com and okta-emea.com.',
    }),
    new FormField('primary_email', 'boolean'),
  ];

  validations: Validations = {
    org_name: [{ type: 'presence', message: 'Org name is required.' }],
    api_token: [{ type: 'presence', message: 'API token is required.' }],
  };
}
