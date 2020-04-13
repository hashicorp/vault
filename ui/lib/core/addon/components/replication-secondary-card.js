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
  state: computed('data', function() {
    let dr = this.data.dr;
    return dr.state ? dr.state : 'unknown';
  }),
  connection: computed('data', function() {
    return this.data.drStateDisplay ? this.data.drStateDisplay : 'unknown';
  }),
  lastWAL: computed('data', function() {
    let dr = this.data.dr;
    return dr.lastWAL ? dr.lastWAL : '0';
  }),
  lastRemoteWAL: computed('data', function() {
    let dr = this.data.dr;
    return dr.lastRemoteWAL ? dr.lastRemoteWAL : '0';
  }),
  delta: computed('data', function() {
    let dr = this.data.dr;
    let lastWAL = dr && dr.lastWAL ? dr.lastWAL : 0;
    let lastRemoteWAL = dr && dr.lastRemoteWAL ? dr.lastRemoteWAL : 0;

    return Math.abs(lastWAL - lastRemoteWAL);
  }),
});
