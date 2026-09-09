/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { Validations } from 'vault/app-types';
import type { OidcWriteClientRequest } from '@hashicorp/vault-client-typescript';

type OidcClientFormData = OidcWriteClientRequest & {
  name: string;
};

export default class OidcClientForm extends Form<OidcClientFormData> {
  formFieldGroups = [
    new FormFieldGroup('default', [
      new FormField('name', 'string', {
        label: 'Application name',
        editDisabled: true,
      }),
      new FormField('client_type', 'string', {
        label: 'Type',
        subText:
          'Specify whether the application type is confidential or public. The public type must use PKCE. This cannot be edited later.',
        editType: 'radio',
        editDisabled: true,
        possibleValues: ['confidential', 'public'],
      }),
      new FormField('redirect_uris', 'array', {
        label: 'Redirect URIs',
        subText:
          'One of these values must exactly match the redirect_uri parameter value used in each authentication request.',
        editType: 'stringArray',
      }),
    ]),
    new FormFieldGroup('More options', [
      // SearchSelect within the FormField component works in conjunction with Ember Data Models
      // we can still use the component since it supports passing in an array of objects as options for the select
      // yield out the field so keys can be fetched in the route and passed directly to SearchSelect
      new FormField('key', undefined, {
        label: 'Signing key',
        subText: 'Add a key to sign and verify the JSON web tokens (JWT). This cannot be edited later.',
        editType: 'yield',
        editDisabled: true,
      }),
      new FormField('id_token_ttl', undefined, {
        label: 'ID Token TTL',
        editType: 'ttl',
        defaultValue: '24h',
      }),
      new FormField('access_token_ttl', undefined, {
        label: 'Access Token TTL',
        editType: 'ttl',
        defaultValue: '24h',
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
    key: [{ type: 'presence', message: 'Key is required.' }],
  };
}
