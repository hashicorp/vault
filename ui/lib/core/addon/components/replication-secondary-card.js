/**
 * @module ReplicationSecondaryCard
 * ARG TODO finish
 *
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-secondary-card';

export default Component.extend({
  layout,
  delta: computed('data', function() {
    let dr = this.data.dr;
    let lastWAL = dr && dr.lastWAL ? dr.lastWAL : 0;
    let lastRemoteWAL = dr && dr.lastRemoteWAL ? dr.lastRemoteWAL : 0;

    return Math.abs(lastWAL - lastRemoteWAL);
  }),
  errorClass: computed('data', 'title', 'state', 'connection', function() {
    let dr = this.data.dr;

    if (!dr) {
      return false;
    }

    if (this.title === 'States') {
      if (this.state === 'idle' || this.connection === 'transient-failure') {
        return true;
      }
    }
  }),
});
