/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { SshConfigureCaRequest } from '@hashicorp/vault-client-typescript';

export default class SshConfigForm extends Form<SshConfigureCaRequest> {
  validations: Validations = {
    generate_signing_key: [
      {
        validator(data: SshConfigForm['data']) {
          const { public_key, private_key, generate_signing_key } = data;
          // if generateSigningKey is false, both public and private keys are required
          if (!generate_signing_key && (!public_key || !private_key)) {
            return false;
          }
          return true;
        },
        message: 'Provide a Public and Private key or set "Generate Signing Key" to true.',
      },
    ],
    public_key: [
      {
        validator(data: SshConfigForm['data']) {
          const { public_key, private_key } = data;
          // regardless of generateSigningKey, if one key is set they both need to be set.
          return public_key || private_key ? !!(public_key && private_key) : true;
        },
        message: 'You must provide a Public and Private keys or leave both unset.',
      },
    ],
  };

  formFields = [
    new FormField('private_key', 'string', { sensitive: true }),
    new FormField('public_key', 'string', { sensitive: true }),
    new FormField('generate_signing_key', 'boolean'),
  ];
}
