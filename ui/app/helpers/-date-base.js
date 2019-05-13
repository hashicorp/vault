import { run } from '@ember/runloop';
import Helper from '@ember/component/helper';
import { get } from '@ember/object';

export default Helper.extend({
  disableInterval: false,

  compute(value, { interval }) {
    if (get(this, 'disableInterval')) {
      return;
    }

    this.clearTimer();

    if (interval) {
      /*
             * NOTE: intentionally a setTimeout so tests do not block on it
             * as the run loop queue is never clear so tests will stay locked waiting
             * for queue to clear.
             */
      this.intervalTimer = setTimeout(() => {
        run(() => this.recompute());
      }, parseInt(interval, 10));
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
