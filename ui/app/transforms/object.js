/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Transform from '@ember-data/serializer/transform';
import { typeOf } from '@ember/utils';
/*
  DS.attr('object')
*/
export default class ObjectTransform extends Transform {
  deserialize(value) {
    if (typeOf(value) !== 'object') {
      return {};
    } else {
      return value;
    }
  }
  serialize(value) {
    if (typeOf(value) !== 'object') {
      return {};
    } else {
      return value;
    }
  }
}
