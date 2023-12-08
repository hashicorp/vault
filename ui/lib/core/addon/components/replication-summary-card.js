/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { get } from '@ember/object';
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
  get key() {
    return this.args.title === 'Performance' ? 'performance' : 'dr';
  }
  get lastWAL() {
    return get(this.args.replicationDetails, `${this.key}.lastWAL`) || 0;
  }
  get merkleRoot() {
    return get(this.args.replicationDetails, `${this.key}.merkleRoot`) || 'no hash found';
  }
  get knownSecondariesCount() {
    return get(this.args.replicationDetails, `${this.key}.knownSecondaries.length`) || 0;
  }
}
