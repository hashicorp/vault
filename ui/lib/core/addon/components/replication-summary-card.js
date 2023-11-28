/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

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

export default class ReplicationSummaryCard extends Component {
  get lastDrWAL() {
    return this.args.replicationDetails.dr?.lastWAL || 0;
  }
  get lastPerformanceWAL() {
    return this.args.replicationDetails.performance?.lastWAL || 0;
  }
  get merkleRootDr() {
    return this.args.replicationDetails.dr?.merkleRoot || '';
  }
  get merkleRootPerformance() {
    return this.args.replicationDetails.performance?.merkleRoot || '';
  }
  get knownSecondariesDr() {
    const knownSecondaries = this.args.replicationDetails.dr.knownSecondaries;
    return knownSecondaries?.length ?? 0;
  }
  get knownSecondariesPerformance() {
    const knownSecondaries = this.args.replicationDetails.performance.knownSecondaries;
    return knownSecondaries?.length ?? 0;
  }
}
