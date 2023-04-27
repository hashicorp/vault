/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-summary-card';

/**
 * @module ReplicationSummaryCard
 * The `ReplicationSummaryCard` is a card-like component.  It displays cluster mode details for both DR and Performance
 *
 * @example
 * ```js
 * <ReplicationSummaryCard
    @title='States'
    @replicationDetails={DS.Model.replicationDetailsSummary}
    />
 * ```
 * @param {String} [title=null] - The title to be displayed on the top left corner of the card.
 * @param {Object} replicationDetails=null - An Ember data object computed off the Ember Model.  It combines the Model.dr and Model.performance objects into one and contains details specific to the mode replication.
 */

export default Component.extend({
  layout,
  title: null,
  replicationDetails: null,
  lastDrWAL: computed('replicationDetails.dr.lastWAL', function () {
    return this.replicationDetails.dr.lastWAL || 0;
  }),
  lastPerformanceWAL: computed('replicationDetails.performance.lastWAL', function () {
    return this.replicationDetails.performance.lastWAL || 0;
  }),
  merkleRootDr: computed('replicationDetails.dr.merkleRoot', function () {
    return this.replicationDetails.dr.merkleRoot || '';
  }),
  merkleRootPerformance: computed('replicationDetails.performance.merkleRoot', function () {
    return this.replicationDetails.performance.merkleRoot || '';
  }),
  knownSecondariesDr: computed('replicationDetails.dr.knownSecondaries', function () {
    const knownSecondaries = this.replicationDetails.dr.knownSecondaries;
    return knownSecondaries.length;
  }),
  knownSecondariesPerformance: computed('replicationDetails.performance.knownSecondaries', function () {
    const knownSecondaries = this.replicationDetails.performance.knownSecondaries;
    return knownSecondaries.length;
  }),
});
