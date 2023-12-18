/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Transform from '@ember-data/serializer/transform';
import { isArray, A } from '@ember/array';
/*
  This should go inside a globally available place for all apps

  DS.attr('array')
*/
export default class ArrayTransform extends Transform {
  deserialize(value) {
    if (isArray(value)) {
      return A(value);
    } else {
      return A();
    }
  }
  serialize(value) {
    if (isArray(value)) {
      return A(value);
    } else {
      return A();
    }
  }
}
