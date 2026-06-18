/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { Validations } from 'vault/app-types';
import type { OidcWriteProviderRequest } from '@hashicorp/vault-client-typescript';

type OidcProviderFormData = OidcWriteProviderRequest & {
  name: string;
};

export default class OidcProviderForm extends Form<OidcProviderFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('name', 'string', { editDisabled: true }),
      new FormField('issuer', 'string', {
        subText:
          "The scheme, host, and optional port for your issuer. This will be used to build the URL that validates ID tokens. Defaults to a URL with Vault's api_addr.",
        placeholder: 'e.g. https://example.com:8200',
        docLink: '/vault/api-docs/secret/identity/oidc-provider#create-or-update-a-provider',
      }),
      // SearchSelect within the FormField component works in conjunction with Ember Data Models
      // we can still use the component since it supports passing in an array of objects as options for the select
      // yield out the field so scopes can be fetched in the route and passed directly to SearchSelect
      new FormField('supported_scopes', undefined, {
        label: 'Supported scopes',
        subText: 'Scopes define information about a user and the OIDC service. Optional.',
        editType: 'yield',
      }),
    ]),
  ];

  validations: Validations = {
    name: [
      { type: 'presence', message: 'Name is required.' },
      {
        type: 'containsWhiteSpace',
        message: 'Name cannot contain whitespace.',
      },
    ],
  };
}
