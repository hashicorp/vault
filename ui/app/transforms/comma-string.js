/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Transform from '@ember-data/serializer/transform';

/**
 * transforms array types from the server to a comma separated string
 * useful when using changedAttributes() in the serializer to track attribute changes for PATCH requests
 * because arrays are not trackable and strings are!
 */
export default class CommaString extends Transform {
  deserialize(serialized) {
    if (Array.isArray(serialized)) {
      return serialized.join(',');
    }
    return serialized;
  }

  serialize(deserialized) {
    if (typeof deserialized === 'string') {
      return deserialized.split(',');
    }
    return deserialized;
  }
}
