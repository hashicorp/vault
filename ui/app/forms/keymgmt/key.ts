/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import {
  KeyManagementUpdateKeyRequestTypeEnum,
  KeyManagementUpdateKeyRequest,
} from '@hashicorp/vault-client-typescript';

const KEY_TYPES = Object.values(KeyManagementUpdateKeyRequestTypeEnum);

type KeyFormData = KeyManagementUpdateKeyRequest & {
  name: string;
};

export default class KeymgmtKeyForm extends Form<KeyFormData> {
  icon = 'key';

  get formFieldGroups() {
    // Edit mode: only show min_enabled_version and deletion_allowed
    if (!this.isNew) {
      return [
        new FormFieldGroup('default', [
          new FormField('min_enabled_version', 'number', {
            label: 'Minimum enabled version',
            defaultValue: 0,
            defaultShown: 'All versions enabled',
          }),
          new FormField('deletion_allowed', 'boolean', {
            label: 'Allow deletion',
            defaultValue: false,
          }),
        ]),
      ];
    }

    // Create mode: show all fields
    return [
      new FormFieldGroup('default', [
        new FormField('name', 'string', {
          label: 'Key name',
          subText: 'This is the name of the key that shows in Vault.',
        }),
        new FormField('type', 'string', {
          label: 'Type',
          subText: 'The type of cryptographic key that will be created.',
          possibleValues: KEY_TYPES,
          defaultValue: KeyManagementUpdateKeyRequestTypeEnum.RSA_2048,
        }),
        new FormField('deletion_allowed', 'boolean', {
          label: 'Allow deletion',
          defaultValue: false,
        }),
      ]),
    ];
  }
}
