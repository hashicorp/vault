/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';
import {regions} from "vault/helpers/aws-regions";
import {withFormFields} from "vault/decorators/model-form-fields";

const formFields = ['scope'];

@withFormFields(formFields)
export default class Config extends Model {
  @attr('string', {
    editType: 'radio',
    label: 'Storage scope',
    subText:
      'Defines the behavior of the physical storage used for cubbyhole\'s secrets. If per token, cubbyhole is destroyed after the token\'s expiration. If per identity, secrets are persisted and are linked to user\'s identity lifetime',
    possibleValues: ['per-token', 'per-identity'],
    defaultValue: 'per-token',
  })
  scope;
  get attrs() {
    const keys = ['scope'];
    return expandAttributeMeta(this, keys);
  }
}
