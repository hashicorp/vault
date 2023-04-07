/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';

@withExpandedAttributes()
export default class NamespaceModel extends Model {
  @attr('string', {
    validationAttr: 'pathIsValid',
    invalidMessage: 'You have entered and invalid path please only include letters, numbers, -, ., and _.',
  })
  path;

  get pathIsValid() {
    return this.path && this.path.match(/^[\w\d-.]+$/g);
  }

  get fields() {
    return ['path'].map((f) => this.allByKey[f]);
  }
}
