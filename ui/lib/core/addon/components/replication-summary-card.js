import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-summary-card';

/**
 * @module ReplicationSummaryCard
 * ReplicationSummaryCard components
 *
 * @example
 * ```js
 * <ReplicationSummaryCard
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
  state: computed('replicationDetails.dr.{state}', 'replicationDetails.performance.{state}', function() {
    // ARG TODO return for top part of dashboard
    // return this.replicationDetails.dr && this.replicationDetails.dr.state
    //   ? this.replicationDetails.dr.state
    //   : 'unknown';
  }),
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
  merkleRootDr: computed('replicationDetails.dr.{merkleRoot}', function() {
    return this.replicationDetails.dr && this.replicationDetails.dr.merkleRoot
      ? this.replicationDetails.dr.merkleRoot
      : '';
  }),
  merkleRootPerformance: computed('replicationDetails.performance.{merkleRoot}', function() {
    return this.replicationDetails.performance && this.replicationDetails.performance.merkleRoot
      ? this.replicationDetails.performance.merkleRoot
      : '';
  }),
  knownSecondariesDr: computed('replicationDetails.dr.{knownSecondaries}', function() {
    const knownSecondaries = this.replicationDetails.dr.knownSecondaries;
    return knownSecondaries.length;
  }),
  knownSecondariesPerformance: computed('replicationDetails.performance.{knownSecondaries}', function() {
    const knownSecondaries = this.replicationDetails.performance.knownSecondaries;
    return knownSecondaries.length;
  }),
});
