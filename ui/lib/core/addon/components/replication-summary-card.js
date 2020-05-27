import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-summary-card';
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
  stateDr: computed('replicationDetails.dr.{state}', function() {
    return this.replicationDetails.dr && this.replicationDetails.dr.state
      ? this.replicationDetails.dr.state
      : 'unknown';
  }),
  statePerformance: computed('replicationDetails.performance.{state}', function() {
    return this.replicationDetails.performance && this.replicationDetails.performance.state
      ? this.replicationDetails.performance.state
      : 'unknown';
  }),
  connection: computed(
    'replicationDetails.dr.{connection_state}',
    'replicationDetails.performance.{connection_state}',
    function() {
      // ARG TODO figure this out
      return 'figure-me-out';
    }
  ),
  lastDrWAL: computed('replicationDetails.dr.{lastWAL}', function() {
    return this.replicationDetails.dr && this.replicationDetails.dr.lastWAL
      ? this.replicationDetails.dr.lastWAL
      : 0;
  }),
  lastPerformanceWAL: computed('replicationDetails.performance.{lastWAL}', function() {
    return this.replicationDetails.performance && this.replicationDetails.performance.lastWAL
      ? this.replicationDetails.performance.lastWAL
      : 0;
  }),
  inSyncState: computed('stateDr', 'statePerformance', function() {
    // if our definition of what is considered 'synced' changes,
    // we should use the clusterStates helper instead
    // ARG FIGURE OUT
    return this.state === 'stream-wals';
  }),
  hasErrorClass: computed(
    'replicationDetails',
    'title',
    'stateDr',
    'statePerformance',
    'connection',
    function() {
      const { title, stateDr, statePerformance, connection } = this;
      // ARG figure out
      // only show errors on the state card
      // if (title === 'States') {
      //   const currentClusterisOk = clusterStates([state]).isOk;
      //   const primaryIsOk = clusterStates([connection]).isOk;
      //   return !(currentClusterisOk && primaryIsOk);
      // }
      return false;
    }
  ),
  primaryClusterAddr: computed('replicationDetails.dr.{primaryClusterAddr}', function() {
    return 'meep';
    // return this.replicationDetails.primaryClusterAddr;
  }),
});
