/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';
import { withModelValidations } from 'vault/decorators/model-validations';

@withExpandedAttributes()
@withModelValidations({
  path: [
    { type: 'presence', message: `Path can't be blank.` },
    { type: 'endsInSlash', message: `Path can't end in forward slash '/'.` },
    {
      type: 'containsWhiteSpace',
      message: "Path can't contain whitespace.",
    },
  ],
})
export default class NamespaceModel extends Model {
  @attr('string')
  path;

  get fields() {
    return ['path'].map((f) => this.allByKey[f]);
  }
}
