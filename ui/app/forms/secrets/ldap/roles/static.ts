/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Form, { type FormOptions } from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { LdapWriteStaticRoleRequest } from '@hashicorp/vault-client-typescript';

type LdapStaticRoleFormData = LdapWriteStaticRoleRequest & {
  name: string;
};

type LdapStaticRoleFormOptions = FormOptions & {
  isSelfManaged?: boolean;
};

const USERNAME_SUB_TEXT = {
  selfManaged: 'The username of the existing LDAP entry to manage.',
  default:
    "The name of the user to be used when logging in. This is useful when DN isn't used for login purposes.",
};

export default class LdapStaticRoleForm extends Form<LdapStaticRoleFormData> {
  // Self-managed mounts require dn and password fields, which root-managed mounts do not. Held on the
  // instance so the field indicators and the validators below read the same value.
  isSelfManaged: boolean;

  constructor(data: Partial<LdapStaticRoleFormData> = {}, options: LdapStaticRoleFormOptions = {}) {
    super(data, options);
    this.isSelfManaged = options.isSelfManaged ?? false;

    if (this.isSelfManaged) {
      // Only override the username subText for self-managed configs; the field
      // already carries USERNAME_SUB_TEXT.default as its initial value.
      const usernameField = this.formFields.find((f) => f.name === 'username');
      if (usernameField) {
        usernameField.options.subText = USERNAME_SUB_TEXT.selfManaged;
      }
      const dnField = this.formFields.find((f) => f.name === 'dn');
      if (dnField) {
        dnField.options.isRequired = true;
      }
      const passwordField = this.formFields.find((f) => f.name === 'password');
      if (passwordField) {
        // Only on create: an existing role already has a stored password, and leaving it untouched is valid.
        passwordField.options.isRequired = this.isNew;
      }
    }
  }

  formFields = [
    new FormField('name', 'string', {
      label: 'Role name',
      subText: 'The name of the role that will be used in Vault.',
      editDisabled: true,
      isRequired: true,
    }),
    new FormField('dn', 'string', {
      label: 'Distinguished name',
      subText: 'Distinguished name (DN) of entry Vault should manage.',
      editDisabled: true,
    }),
    new FormField('username', 'string', {
      label: 'Username',
      subText: USERNAME_SUB_TEXT.default,
      isRequired: true,
    }),
    new FormField('password', 'string', {
      editType: 'password',
      label: 'Password',
      subText:
        'Enter the account’s current LDAP password. Vault will use it to manage password rotations going forward.',
    }),
    new FormField('rotation_period', 'number', {
      editType: 'ttl',
      label: 'Rotation period',
      helperTextEnabled:
        'Specifies how often Vault rotates the password. Enter 0 to disable automatic password rotation.',
      hideToggle: true,
      emptyMeansZero: true,
    }),
  ];

  validations: Validations = {
    name: [
      { type: 'presence', message: 'Role name is required' },
      {
        validator: ({ name }: LdapStaticRoleFormData) => {
          // Presence reports an empty name, so skip here to avoid two messages on one field.
          if (!name) return true;
          // Allow alphanumeric, hyphens, underscores, periods, and forward slashes
          const validPattern = /^[a-z0-9\-_./]+$/;
          return validPattern.test(name);
        },
        message:
          'Name must be lowercase and can only contain alphanumeric characters, hyphens, underscores, periods, and forward slashes.',
      },
    ],
    username: [{ type: 'presence', message: 'Enter the username of the LDAP account' }],
    // rotation_period is intentionally unvalidated: 0 is a legitimate value meaning no automated
    // rotation, so an empty field must never block submission.
    // dn is required for self-managed mounts; password only when creating one.
    dn: [
      {
        validator: (data: LdapStaticRoleFormData) => !this.isSelfManaged || Boolean(data.dn?.trim()),
        message: 'Enter the distinguished name (DN) of the LDAP account',
      },
    ],
    password: [
      {
        validator: (data: LdapStaticRoleFormData) =>
          !this.isSelfManaged || !this.isNew || Boolean(data.password?.trim()),
        message: 'Enter the current password for this LDAP account',
      },
    ],
  };
}
