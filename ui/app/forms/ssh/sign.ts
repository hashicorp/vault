/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { Validations } from 'vault/vault/app-types';

interface SshSignData {
  public_key: string;
  key_id?: string;
  valid_principals?: string;
  cert_type?: string;
  critical_options?: Record<string, unknown>;
  extensions?: Record<string, unknown>;
  ttl?: string;
}

export default class SshSignForm extends Form<SshSignData> {
  get formFieldGroups() {
    return [
      new FormFieldGroup('default', [
        new FormField('public_key', 'string', { label: 'Public key' }),
        new FormField('valid_principals', 'string', {
          label: 'Valid principals',
          helpText:
            'Specifies valid principals, either usernames or hostnames, that the certificate should be signed for. Required unless the role has specified allow_empty_principals.',
        }),
      ]),
      new FormFieldGroup('More options', [
        new FormField('key_id', 'string', { label: 'Key ID' }),
        new FormField('cert_type', 'string', {
          label: 'Certificate Type',
          possibleValues: ['user', 'host'],
          defaultValue: 'user',
        }),
        new FormField('critical_options', 'object', { label: 'Critical Options' }),
        new FormField('extensions', 'object', { label: 'Extensions' }),
        new FormField('ttl', 'string', { label: 'TTL', editType: 'ttl' }),
      ]),
    ];
  }

  validations: Validations = {
    public_key: [{ type: 'presence', message: 'Public Key is required' }],
  };
}
