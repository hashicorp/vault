/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { MfaCreateTotpMethodRequest } from '@hashicorp/vault-client-typescript';

type MfaCreateTotpMethodFormData = MfaCreateTotpMethodRequest & {
  name: string;
};

export default class MfaCreateTotpMethodForm extends Form<MfaCreateTotpMethodFormData> {
  type = 'totp';

  formFields = [
    new FormField('issuer', 'string', {
      label: 'Issuer',
      subText: 'The human-readable name of the keys issuing organization.',
    }),
    new FormField('period', undefined, {
      label: 'Period',
      editType: 'ttl',
      helperTextEnabled: 'How long each generated TOTP is valid.',
      hideToggle: true,
    }),
    new FormField('key_size', 'number', {
      label: 'Key size',
      subText: 'The size in bytes of the Vault generated key.',
      helpText: 'Byte size of the generated key.',
    }),
    new FormField('qr_size', 'number', {
      label: 'QR size',
      subText: 'The pixel size of the generated square QR code.',
      helpText: 'Pixel size of the QR code.',
    }),
    new FormField('algorithm', 'string', {
      label: 'Algorithm',
      editType: 'radio',
      possibleValues: ['SHA1', 'SHA256', 'SHA512'],
      subText: 'The hashing algorithm used to generate the TOTP code.',
    }),
    new FormField('digits', 'number', {
      label: 'Digits',
      editType: 'radio',
      possibleValues: [6, 8],
      subText: 'The number digits in the generated TOTP code.',
      helpText: 'TOTP code length.',
    }),
    new FormField('skew', 'number', {
      label: 'Skew',
      editType: 'radio',
      possibleValues: [0, 1],
      subText: 'The number of delay periods allowed when validating a TOTP token.',
    }),
    new FormField('max_validation_attempts', 'number'),
    new FormField('enable_self_enrollment', 'boolean', {
      label: 'Enable self-enrollment',
      editType: 'toggleButton',
      helperTextEnabled:
        'Let end users enroll in this MFA method on their own. You still control which auth mounts, groups, or entities it applies to.',
      helperTextDisabled:
        'Let end users enroll in this MFA method on their own. You still control which auth mounts, groups, or entities it applies to.',
    }),
  ];

  validations: Validations = {
    issuer: [{ type: 'presence', message: 'Issuer is required.' }],
  };
}
