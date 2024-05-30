/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { run } from '@ember/runloop';
import Helper from '@ember/component/helper';
import Ember from 'ember';

export default Helper.extend({
  disableInterval: false,

  compute(value, { interval }) {
    if (Ember.testing) {
      // issues with flaky test, suspect it has to the do with the run loop not being cleared as intended farther down.
      return;
    }
    if (this.disableInterval) {
      return;
    }

    this.clearTimer();

    if (interval) {
      /*
       * NOTE: intentionally a setTimeout so tests do not block on it
       * as the run loop queue is never clear so tests will stay locked waiting
       * for queue to clear.
       */
      this.intervalTimer = setTimeout(
        () => {
          run(() => this.recompute());
        },
        parseInt(interval, 10)
      );
    }
  },

  clearTimer() {
    clearTimeout(this.intervalTimer);
  },

  destroy() {
    this.clearTimer();
    this._super(...arguments);
  },
});
