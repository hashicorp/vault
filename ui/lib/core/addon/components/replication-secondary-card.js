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
    let lastRemoteWAL = dr.lastRemoteWAL ? dr.lastRemoteWAL : 0;

    return Math.abs(lastWAL - lastRemoteWAL);
  }),
});
