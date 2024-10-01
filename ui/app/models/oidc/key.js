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
export default class OidcKeyModel extends Model {
  @attr('string', { editDisabled: true }) name;
  @attr('string', {
    defaultValue: 'RS256',
    possibleValues: ['RS256', 'RS384', 'RS512', 'ES256', 'ES384', 'ES512', 'EdDSA'],
  })
  algorithm;

  @attr({ editType: 'ttl', defaultValue: '24h' }) rotationPeriod;
  @attr({ label: 'Verification TTL', editType: 'ttl', defaultValue: '24h' }) verificationTtl;
  @attr('array', { label: 'Allowed applications' }) allowedClientIds; // no editType because does not use form-field component

  // TODO refactor when field-to-attrs is refactored as decorator
  _attributeMeta = null; // cache initial result of expandAttributeMeta in getter and return
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, [
        'name',
        'algorithm',
        'rotationPeriod',
        'verificationTtl',
      ]);
    }
    return this._attributeMeta;
  }

  @lazyCapabilities(apiPath`identity/oidc/key/${'name'}`, 'name') keyPath;
  @lazyCapabilities(apiPath`identity/oidc/key/${'name'}/rotate`, 'name') rotatePath;
  get canRead() {
    return this.keyPath.get('canRead');
  }
  get canEdit() {
    return this.keyPath.get('canUpdate');
  }
  get canRotate() {
    return this.rotatePath.get('canUpdate');
  }
  get canDelete() {
    return this.keyPath.get('canDelete');
  }
}
