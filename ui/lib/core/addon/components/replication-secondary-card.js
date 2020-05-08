import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-secondary-card';
import { clusterStates } from 'core/helpers/cluster-states';

/**
 * @module ReplicationSecondaryCard
 * ReplicationSecondaryCard components
 *
 * @example
 * ```js
 * <ReplicationSecondaryCard
    @title='States'
    @replicationDetails=replicationDetails
    />
 * ```
 * @param {string} [title=null] - The title to be displayed on the top left corner of the card.
 * @param replicationDetails=null{DS.Model.replicationDetails} - An Ember data object off the Ember data model.  It is computed at the parent component and passed through to this component.
 */

export default Component.extend({
  layout,
  title: null,
  replicationDetails: null,
  state: computed('replicationDetails.{state}', function() {
    return this.replicationDetails && this.replicationDetails.state
      ? this.replicationDetails.state
      : 'unknown';
  }),
  connection: computed('replicationDetails.{connection_state}', function() {
    return this.replicationDetails.connection_state ? this.replicationDetails.connection_state : 'unknown';
  }),
  lastRemoteWAL: computed('replicationDetails.{lastRemoteWAL}', function() {
    return this.replicationDetails && this.replicationDetails.lastRemoteWAL
      ? this.replicationDetails.lastRemoteWAL
      : 0;
  }),
  inSyncState: computed('state', function() {
    // if our definition of what is considered 'synced' changes,
    // we should use the clusterStates helper instead
    return this.state === 'stream-wals';
  }),
  hasErrorClass: computed('replicationDetails', 'title', 'state', 'connection', function() {
    const { title, state, connection } = this;

    // only show errors on the state card
    if (title === 'States') {
      const currentClusterisOk = clusterStates([state]).isOk;
      const primaryIsOk = clusterStates([connection]).isOk;
      return !(currentClusterisOk && primaryIsOk);
    }
    return false;
  }),
  primaryClusterAddr: computed('replicationDetails.{primaryClusterAddr}', function() {
    return this.replicationDetails.primaryClusterAddr;
  }),
});
