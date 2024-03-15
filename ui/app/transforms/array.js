/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Transform from '@ember-data/serializer/transform';
import { isArray } from '@ember/array';
/*
  This should go inside a globally available place for all apps

  DS.attr('array')
*/
export default class ArrayTransform extends Transform {
  deserialize(value) {
    if (isArray(value)) {
      return [...value];
    } else {
      return [];
    }
  }
  serialize(value) {
    if (isArray(value)) {
      return [...value];
    } else {
      return [];
    }
  }
}
