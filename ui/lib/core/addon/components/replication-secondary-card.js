/**
 * @module ReplicationSecondaryCard
 * ARG TODO finish
 *
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-secondary-card';
import { clusterStates } from 'core/helpers/cluster-states';

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
    return this.dr.connection_state ? this.dr.connection_state : 'unknown';
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
    // if our definition of what is considered 'synced' changes,
    // we should use the clusterStates helper instead
    return this.state === 'stream-wals';
  }),

  hasErrorClass: computed('data', 'title', 'state', 'connection', function() {
    const { title, state, connection } = this;

    // only show errors on the state card
    if (title === 'States') {
      const currentClusterisOk = clusterStates([state]).isOk;
      const primaryIsOk = clusterStates([connection]).isOk;
      return !(currentClusterisOk && primaryIsOk);
    }
    return false;
  }),
});
