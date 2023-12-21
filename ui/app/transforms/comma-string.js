/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Transform from '@ember-data/serializer/transform';

/**
 * transforms API arrays to a comma separated string so changedAttributes() tracks changes for PATCH requests
 */
export default class CommaString extends Transform {
  deserialize(serialized) {
    return serialized.join(',');
  }

  serialize(deserialized) {
    return deserialized.split(',');
  }
}
