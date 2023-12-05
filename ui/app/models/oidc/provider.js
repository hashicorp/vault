/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
    {
      type: 'containsWhiteSpace',
      message: 'Name cannot contain whitespace.',
    },
  ],
};

@withModelValidations(validations)
export default class OidcProviderModel extends Model {
  @attr('string', { editDisabled: true }) name;
  @attr('string', {
    subText:
      'The scheme, host, and optional port for your issuer. This will be used to build the URL that validates ID tokens.',
    placeholderText: 'e.g. https://example.com:8200',
    docLink: '/vault/api-docs/secret/identity/oidc-provider#create-or-update-a-provider',
    helpText: `Optional. This defaults to a URL with Vault's api_addr`,
  })
  issuer;

  @attr('array', {
    label: 'Supported scopes',
    subText: 'Scopes define information about a user and the OIDC service. Optional.',
    editType: 'searchSelect',
    models: ['oidc/scope'],
    fallbackComponent: 'string-list',
    onlyAllowExisting: true,
  })
  scopesSupported;

  @attr('array', { label: 'Allowed applications' }) allowedClientIds; // no editType because does not use form-field component

  // TODO refactor when field-to-attrs is refactored as decorator
  _attributeMeta = null; // cache initial result of expandAttributeMeta in getter and return
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['name', 'issuer', 'scopesSupported']);
    }
    return this._attributeMeta;
  }

  @lazyCapabilities(apiPath`identity/oidc/provider/${'name'}`, 'name') providerPath;
  get canRead() {
    return this.providerPath.get('canRead');
  }
  get canEdit() {
    return this.providerPath.get('canUpdate');
  }
  get canDelete() {
    return this.providerPath.get('canDelete');
  }
}
