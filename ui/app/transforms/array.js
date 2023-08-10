/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Transform from '@ember-data/serializer/transform';
import { isArray, A } from '@ember/array';
/*
  This should go inside a globally available place for all apps

  DS.attr('array')
*/
export default Transform.extend({
  deserialize(value) {
    if (isArray(value)) {
      return A(value);
    } else {
      return A();
    }
  },
  serialize(value) {
    if (isArray(value)) {
      return A(value);
    } else {
      return A();
    }
  },
});
