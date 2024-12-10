/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

const lifetimeValues = ['session', 'persistent']
const validations = {
  lifetime: [
    {
      validator(model) {
        const { lifetime } = model;
        return lifetimeValues.includes(lifetime);
      },
      message: 'You can choose either session or persistent',
    },
  ],
};
// there are more options available on the API, but the UI does not support them yet.
@withModelValidations(validations)
export default class Config extends Model {
  @attr('string', {
    editType: 'radio',
    label: 'Storage scope',
    subText:
      "Defines the behavior of the storage used for cubbyhole's secrets. Per default, cubbyhole data is destroyed after the token's expiration. If set persistent, secrets are persisted and are linked to user's entity",
    possibleValues: lifetimeValues,
    defaultValue: 'per-token',
  })
  scope;
  get attrs() {
    const keys = ['scope'];
    return expandAttributeMeta(this, keys);
  }
  get formFields() {
    const keys = ['scope'];
    return expandAttributeMeta(this, keys);
  }
}
