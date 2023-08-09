/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  name: [
    { type: 'presence', message: 'Name is required.' },
    {
      type: 'containsWhiteSpace',
      message: 'Name cannot contain whitespace.',
    },
  ],
  key: [{ type: 'presence', message: 'Key is required.' }],
};

@withModelValidations(validations)
export default class OidcClientModel extends Model {
  @attr('string', { label: 'Application name', editDisabled: true }) name;
  @attr('string', {
    label: 'Type',
    subText:
      'Specify whether the application type is confidential or public. The public type must use PKCE. This cannot be edited later.',
    editType: 'radio',
    editDisabled: true,
    defaultValue: 'confidential',
    possibleValues: ['confidential', 'public'],
  })
  clientType;

  @attr('array', {
    label: 'Redirect URIs',
    subText:
      'One of these values must exactly match the redirect_uri parameter value used in each authentication request.',
    editType: 'stringArray',
  })
  redirectUris;

  // >> MORE OPTIONS TOGGLE <<

  @attr('string', {
    label: 'Signing key',
    subText: 'Add a key to sign and verify the JSON web tokens (JWT). This cannot be edited later.',
    editType: 'searchSelect',
    editDisabled: true,
    onlyAllowExisting: true,
    defaultValue() {
      return ['default'];
    },
    fallbackComponent: 'input-search',
    selectLimit: 1,
    models: ['oidc/key'],
  })
  key;
  @attr({
    label: 'Access Token TTL',
    editType: 'ttl',
    defaultValue: '24h',
  })
  accessTokenTtl;

  @attr({
    label: 'ID Token TTL',
    editType: 'ttl',
    defaultValue: '24h',
  })
  idTokenTtl;

  // >> END MORE OPTIONS TOGGLE <<

  @attr('array', { label: 'Assign access' }) assignments; // no editType because does not use form-field component
  @attr('string', { label: 'Client ID' }) clientId;
  @attr('string') clientSecret;

  // TODO refactor when field-to-attrs util is refactored as decorator
  _attributeMeta = null; // cache initial result of expandAttributeMeta in getter and return
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['name', 'clientType', 'redirectUris']);
    }
    return this._attributeMeta;
  }

  _fieldToAttrsGroups = null;
  // more options fields
  get fieldGroups() {
    if (!this._fieldToAttrsGroups) {
      this._fieldToAttrsGroups = fieldToAttrs(this, [
        { 'More options': ['key', 'idTokenTtl', 'accessTokenTtl'] },
      ]);
    }
    return this._fieldToAttrsGroups;
  }

  // CAPABILITIES //
  @lazyCapabilities(apiPath`identity/oidc/client/${'name'}`, 'name') clientPath;
  get canRead() {
    return this.clientPath.get('canRead');
  }
  get canEdit() {
    return this.clientPath.get('canUpdate');
  }
  get canDelete() {
    return this.clientPath.get('canDelete');
  }
}
