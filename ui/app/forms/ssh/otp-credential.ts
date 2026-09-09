/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import { Validations } from 'vault/vault/app-types';

interface SshOtpCredentialData {
  username: string;
  ip: string;
}

export default class SshOtpCredentialForm extends Form<SshOtpCredentialData> {
  formFields = [
    new FormField('username', 'string', { label: 'Username' }),
    new FormField('ip', 'string', { label: 'IP address' }),
  ];

  validations: Validations = {
    ip: [{ type: 'presence', message: 'IP address is required' }],
  };
}
