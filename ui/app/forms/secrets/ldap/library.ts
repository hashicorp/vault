/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { LdapLibraryConfigureRequest } from '@hashicorp/vault-client-typescript';

// omit 'disable_check_in_enforcement' here to transform it from boolean to string for display
type LdapLibraryFormData = Omit<LdapLibraryConfigureRequest, 'disable_check_in_enforcement'> & {
  disable_check_in_enforcement: string;
  name: string;
};

export default class LdapLibraryForm extends Form<LdapLibraryFormData> {
  formFields = [
    new FormField('name', 'string', {
      label: 'Library name',
      editDisabled: true,
    }),
    new FormField('service_account_names', 'string', {
      editType: 'stringArray',
      label: 'Accounts',
      subText:
        'The names of all the accounts that can be checked out from this set. These accounts must only be used by Vault, and may only be in one set.',
    }),
    new FormField('ttl', 'string', {
      editType: 'ttl',
      label: 'Default lease TTL',
      helperTextDisabled: 'Vault will use the default lease duration.',
      defaultValue: '24h',
      defaultShown: 'Engine default',
    }),
    new FormField('max_ttl', 'string', {
      editType: 'ttl',
      label: 'Max lease TTL',
      helperTextDisabled: 'Vault will use the default lease duration.',
      defaultValue: '24h',
      defaultShown: 'Engine default',
    }),
    // this is a boolean from the server but is transformed in the serializer to display as Disabled or Enabled
    new FormField('disable_check_in_enforcement', 'string', {
      editType: 'radio',
      label: 'Check-in enforcement',
      subText:
        'When enabled, accounts must be checked in by the entity or client token that checked them out. If disabled, anyone with the right permission can check the account back in.',
      possibleValues: ['Disabled', 'Enabled'],
    }),
  ];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Library name is required.' }],
    service_account_names: [{ type: 'presence', message: 'At least one service account is required.' }],
  };
}
