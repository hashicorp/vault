/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { LdapWriteStaticRoleRequest } from '@hashicorp/vault-client-typescript';

type LdapDynamicRoleFormData = LdapWriteStaticRoleRequest & {
  name: string;
};

export default class LdapStaticRoleForm extends Form<LdapDynamicRoleFormData> {
  formFields = [
    new FormField('name', 'string', {
      label: 'Role name',
      subText: 'The name of the role that will be used in Vault.',
      editDisabled: true,
    }),
    new FormField('dn', 'string', {
      label: 'Distinguished name',
      subText: 'Distinguished name (DN) of entry Vault should manage.',
    }),
    new FormField('username', 'string', {
      label: 'Username',
      subText:
        "The name of the user to be used when logging in. This is useful when DN isn't used for login purposes.",
    }),
    new FormField('rotation_period', 'number', {
      editType: 'ttl',
      label: 'Rotation period',
      helperTextEnabled:
        'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds.',
      hideToggle: true,
    }),
  ];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required' }],
    username: [{ type: 'presence', message: 'Username is required.' }],
    rotation_period: [{ type: 'presence', message: 'Rotation Period is required.' }],
  };
}
