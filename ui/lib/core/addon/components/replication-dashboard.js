/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { clusterStates } from 'core/helpers/cluster-states';
import { capitalize } from '@ember/string';

/**
 * @module ReplicationDashboard
 * The `ReplicationDashboard` component is a contextual component of the replication-page component.
 * It organizes cluster data specific to mode (dr or performance) and also the type (primary or secondary).
 * It is the parent contextual component of the replication-<name>-card components.
 *
 * @example
 * ```js
 * <ReplicationDashboard
    @componentToRender='replication-primary-card'
    @isSecondary=false
    @isSummaryDashboard=false
    @replicationDetailsSummary={}
    @replicationDetails={{replicationDetails}}
    @clusterMode=primary
    @reindexingDetails={{reindexingDetails}}
    />
 * ```
 * @param {String} [componentToRender=''] - A string that determines which card component is displayed.  There are three options, replication-primary-card, replication-secondary-card, replication-summary-card.
 * @param {Boolean} [isSecondary=false] - Used to determine the title and display logic.
 * @param {Boolean} [isSummaryDashboard=false] -  Only true when the cluster is both a dr and performance primary. If true, replicationDetailsSummary is populated and used to pass through the cluster details.
 * @param {Object} replicationDetailsSummary=null - An Ember data object computed off the Ember Model.  It combines the Model.dr and Model.performance objects into one and contains details specific to the mode replication.
 * @param {Object} replicationDetails=null - An Ember data object pulled from the Ember Model. It contains details specific to the whether the replication is dr or performance.
 * @param {String} clusterMode=null - The cluster mode passed through to a table component.
 * @param {Object} reindexingDetails=null - An Ember data object used to show a reindexing progress bar.
 */

export default class ReplicationDashboard extends Component {
  get isSyncing() {
    const { state } = this.args.replicationDetails;
    const isSecondary = this.args.isSecondary;
    return isSecondary && state && clusterStates([state]).isSyncing;
  }
  get isReindexing() {
    return !!this.args.replicationDetails.reindex_in_progress;
  }
  get reindexingStage() {
    const stage = this.args.replicationDetails.reindex_stage;
    // specify the stage if we have one
    if (stage) {
      return `: ${capitalize(stage)}`;
    }
    return '';
  }
  get progressBar() {
    const { reindex_building_progress, reindex_building_total } = this.args.replicationDetails;
    let progressBar = null;

    if (reindex_building_progress && reindex_building_total) {
      progressBar = {
        value: reindex_building_progress,
        max: reindex_building_total,
      };
    }

    return progressBar;
  }
  get summaryState() {
    const { replicationDetailsSummary } = this.args;
    const drState = replicationDetailsSummary.dr.state;
    const performanceState = replicationDetailsSummary.performance.state;

    if (drState !== performanceState) {
      // when DR and Performance is enabled on the same cluster,
      // the states should always be the same
      // we are leaving this console log statement to be sure
      console.log('DR State: ', drState, 'Performance State: ', performanceState); // eslint-disable-line
    }

    return drState;
  }
  get reindexMessage() {
    if (!this.args.isSecondary) {
      return 'This can cause a delay depending on the size of the data store. You can <b>not</b> use Vault during this time.';
    }
    return 'This can cause a delay depending on the size of the data store. You can use Vault during this time.';
  }
}
