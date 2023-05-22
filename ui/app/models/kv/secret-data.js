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
  // ARG TODO add must have path name, no spaces, etc.
};
const formFieldProps = ['path', 'data'];

// ARG TODO: so far this is the data endpoint
@withModelValidations(validations)
@withFormFields(formFieldProps)
export default class KVSecretDataModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', {
    label: 'Path',
    subText: 'The path for this secret.',
  })
  path;
  @attr('object', {
    editType: 'kv',
  })
  data;

  @attr('string') createdTime;
  @attr('object') customMetadata;
  @attr('string') deletionTime;
  @attr('boolean') destroyed;
  @attr('number') version;
}
