/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RESTSerializer from '@ember-data/serializer/rest';
import { decamelize } from '@ember/string';

export default RESTSerializer.extend({
  keyForAttribute: function (attr) {
    return decamelize(attr);
  },
});
