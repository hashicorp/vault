/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import OpenApiForm from 'vault/forms/open-api';
import FormField from 'vault/utils/forms/field';

import type { LdapConfigureRequest } from '@hashicorp/vault-client-typescript';
import type Form from 'vault/forms/form';
import type { Validations } from 'vault/app-types';
import FormFieldGroup from 'vault/utils/forms/field-group';

export default class LdapConfigForm extends OpenApiForm<LdapConfigureRequest> {
  constructor(...args: ConstructorParameters<typeof Form>) {
    super('LdapConfigureRequest', ...args);

    this.formFields.forEach((field) => {
      // password_policy field has special handling
      if (field.name === 'password_policy') {
        field.options = {
          editType: 'optionalText',
          label: 'Use custom password policy',
          subText: 'Specify the name of an existing password policy.',
          defaultSubText: 'Unless a custom policy is specified, Vault will use a default.',
          defaultShown: 'Default',
          docLink: '/vault/docs/concepts/password-policies',
        };
      } else if (field.name === 'binddn') {
        // binddn and bindpass subText mentions that they are optional but the docs have them marked as required
        // update text to avoid confusion
        field.options = {
          ...field.options,
          label: 'Administrator distinguished name',
          subText:
            'Distinguished name of the administrator to bind (Bind DN) when performing user and group search. Example: cn=vault,ou=Users,dc=example,dc=com.',
        };
      } else if (field.name === 'bindpass') {
        field.options = {
          ...field.options,
          label: 'Administrator password',
          subText: 'Password to use along with Bind DN when performing user search.',
        };
      }
    });

    // set up formFieldGroups
    const groupsMap = [
      { name: 'default', keys: ['binddn', 'bindpass', 'url', 'password_policy'] },
      {
        name: 'TLS options',
        keys: ['starttls', 'insecure_tls', 'certificate', 'client_tls_cert', 'client_tls_key'],
      },
      {
        name: 'More options',
        keys: ['userdn', 'userattr', 'upndomain', 'connection_timeout', 'request_timeout'],
      },
    ];

    this.formFieldGroups = groupsMap.map(({ name, keys }) => {
      const fields = keys.map((key) => this.formFields.find((field) => field.name === key) as FormField);
      return new FormFieldGroup(name, fields);
    });
  }

  validations: Validations = {
    binddn: [{ type: 'presence', message: 'Administrator distinguished name is required.' }],
    bindpass: [{ type: 'presence', message: 'Administrator password is required.' }],
  };
}
