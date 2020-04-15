/**
 * @module ReplicationSecondaryCard
 * ARG TODO finish
 *
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-secondary-card';

const STATES = {
  streamWals: 'stream-wals',
  idle: 'idle',
  transientFailure: 'transient-failure',
  shutdown: 'shutdown',
};

export default Component.extend({
  layout,
  data: null,
  dr: computed('data', function() {
    let dr = this.data.dr;
    if (!dr) {
      return false;
    }
    return dr;
  }),
  state: computed('dr', function() {
    return this.dr && this.dr.state ? this.dr.state : 'unknown';
  }),
  connection: computed('data', function() {
    return this.data.drStateDisplay ? this.data.drStateDisplay : 'unknown';
  }),
  lastWAL: computed('dr', function() {
    return this.dr && this.dr.lastWAL ? this.dr.lastWAL : 0;
  }),
  lastRemoteWAL: computed('dr', function() {
    return this.dr && this.dr.lastRemoteWAL ? this.dr.lastRemoteWAL : 0;
  }),
  delta: computed('data', function() {
    return Math.abs(this.get('lastWAL') - this.get('lastRemoteWAL'));
  }),
  inSyncState: computed('state', function() {
    if (this.state === STATES.streamWals) {
      return true;
    }
  }),
  hasErrorClass: computed('data', 'title', 'state', 'connection', function() {
    if (this.title === 'States') {
      if (
        this.state === STATES.idle ||
        this.connection === STATES.transientFailure ||
        this.connection === STATES.shutdown
      ) {
        return true;
      }
    }
  }),
});
