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
    // last_wal
    let lastWAL = this.data.dr.lastWAL ? this.data.dr.lastWAL : 0;
    // last_remote_wal
    let lastRemoteWAL = this.data.dr.lastRemoteWAL ? this.data.dr.lastRemoteWAL : 0;
    return Math.abs(lastWAL - lastRemoteWAL);
  }),
});
