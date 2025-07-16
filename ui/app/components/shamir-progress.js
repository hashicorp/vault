/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  threshold: null,
  progress: null,
  classNames: ['shamir-progress'],
  progressDecimal: computed('threshold', 'progress', function () {
    const { threshold, progress } = this;
    if (threshold && progress) {
      return progress / threshold;
    }
    return 0;
  }),
});
