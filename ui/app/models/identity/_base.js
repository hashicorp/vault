/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model from '@ember-data/model';
import { assert } from '@ember/debug';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default Model.extend({
  formFields: computed(function () {
    return assert('formFields should be overridden', false);
  }),

  fields: computed('formFields', 'formFields.[]', function () {
    return expandAttributeMeta(this, this.formFields);
  }),

  identityType: computed('constructor.modelName', function () {
    const modelType = this.constructor.modelName.split('/')[1];
    return modelType;
  }),
});
