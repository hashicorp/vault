/**
 * Copyright IBM Corp. 2026, 2025
 * SPDX-License-Identifier: BUSL-1.1
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
