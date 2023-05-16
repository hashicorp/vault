/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object'; // eslint-disable-line
import { alias } from '@ember/object/computed'; // eslint-disable-line
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

const validations = {
  maxVersions: [
    { type: 'number', message: 'Maximum versions must be a number.' },
    { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
  ],
};
const formFieldProps = ['path'];

// ARG TODO this can be different... maybe ?
@withModelValidations(validations)
@withFormFields(formFieldProps)
export default class KVSecretModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', {
    label: 'Path',
    subText: 'The path for this secret.',
  })
  path;
  @attr('number') version;
  @attr('string') deletionTime;
  @attr('string') createdTime;
  @attr('boolean') destroyed;
  @attr('number') currentVersion;
}
