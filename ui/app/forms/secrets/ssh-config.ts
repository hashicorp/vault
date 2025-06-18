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
    generateSigningKey: [
      {
        validator(data: SshConfigForm['data']) {
          const { publicKey, privateKey, generateSigningKey } = data;
          // if generateSigningKey is false, both public and private keys are required
          if (!generateSigningKey && (!publicKey || !privateKey)) {
            return false;
          }
          return true;
        },
        message: 'Provide a Public and Private key or set "Generate Signing Key" to true.',
      },
    ],
    publicKey: [
      {
        validator(data: SshConfigForm['data']) {
          const { publicKey, privateKey } = data;
          // regardless of generateSigningKey, if one key is set they both need to be set.
          return publicKey || privateKey ? !!(publicKey && privateKey) : true;
        },
        message: 'You must provide a Public and Private keys or leave both unset.',
      },
    ],
  };

  formFields = [
    new FormField('privateKey', 'string', { sensitive: true }),
    new FormField('publicKey', 'string', { sensitive: true }),
    new FormField('generateSigningKey', 'boolean'),
  ];
}
