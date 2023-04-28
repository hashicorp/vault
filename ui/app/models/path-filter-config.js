/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';

export default Model.extend({
  mode: attr('string'),
  paths: attr('array', {
    defaultValue: function () {
      return [];
    },
  }),
});
