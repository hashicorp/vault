/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import type { TransitCreateKeyRequest } from '@hashicorp/vault-client-typescript';
import { Validations } from 'vault/vault/app-types';
import { durationToSeconds } from 'core/utils/duration-utils';

type TransitKeyData = TransitCreateKeyRequest & {
  backend?: string;
  id?: string;
  name?: string;
};

export default class TransitKeyForm extends Form<TransitKeyData> {
  // these set funcs just exist on edit
  set convergentEncryptionValue(val: boolean | undefined) {
    if (val === true) {
      this.data.derived = val;
    }
    this.data.convergent_encryption = val;
  }

  set derivedValue(val: boolean | undefined) {
    if (val === false) {
      this.data.convergent_encryption = val;
    }
    this.data.derived = val;
  }

  get formFields() {
    let fields = [];
    if (this.isNew) {
      fields.push(
        ...[
          new FormField('name', 'string', {
            editDisabled: !this.isNew,
          }),
          new FormField('auto_rotate_period', undefined, {
            editType: 'yield',
            label: ' ',
          }),
          new FormField('type', 'string', {
            possibleValues: [
              'aes128-gcm96',
              'aes256-gcm96',
              'chacha20-poly1305',
              'ecdsa-p256',
              'ecdsa-p384',
              'ecdsa-p521',
              'ed25519',
              'rsa-2048',
              'rsa-3072',
              'rsa-4096',
            ],
          }),
          new FormField('exportable', 'boolean', {
            label: 'Exportable',
            editType: 'checkbox',
          }),
          new FormField('derived', 'boolean', {
            label: 'Derived',
            editType: 'checkbox',
          }),
          new FormField('convergent_encryption', 'boolean', {
            label: 'Enable convergent encryption',
            editType: 'checkbox',
          }),
        ]
      );
      if (this.data.type?.startsWith('ecdsa') || this.data.type?.startsWith('rsa')) {
        fields = fields.filter((field) => !['convergent_encryption', 'derived'].includes(field.name));
      } else if (this.data.type === 'ed25519') {
        fields = fields.filter((field) => field.name !== 'convergent_encryption');
      }
    } else {
      fields.push(
        ...[
          new FormField('deletion_allowed', 'boolean', {
            label: 'Allow deletion',
            editType: 'checkbox',
          }),
          new FormField('auto_rotate_period', undefined, {
            label: ' ',
            editType: 'yield',
          }),
          new FormField('min_decryption_version', undefined, {
            label: 'Minimum decryption version',
            editType: 'yield',
          }),
          new FormField('min_encryption_version', undefined, {
            label: 'Minimum encryption version',
            editType: 'yield',
          }),
        ]
      );
    }
    return fields;
  }

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required' }],
    auto_rotate_period: [
      {
        validator(data: TransitKeyData) {
          const { auto_rotate_period } = data;
          if (auto_rotate_period === undefined) {
            return true;
          } else {
            const duration = durationToSeconds(auto_rotate_period);
            // regardless of generateSigningKey, if one key is set they both need to be set.
            return duration === 0 || duration >= 3600;
          }
        },
        message: 'Duration must be longer than 1 hour or set to 0 to disable auto-rotation.',
      },
    ],
  };

  toJSON() {
    const { isValid, state, invalidFormMessage } = super.toJSON();
    const data = { ...this.data } as Record<string, unknown>;
    // set auto_rotate_period to 0 if it's false to avoid api validation error since the form field is a string and the api expects a number
    if (data['auto_rotate_period'] === false) {
      data['auto_rotate_period'] = 0;
    }

    return { isValid, state, invalidFormMessage, data };
  }
}
