/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import { clusterStates } from 'core/helpers/cluster-states';
import { capitalize } from '@ember/string';
import layout from '../templates/components/replication-dashboard';

/**
 * @module ReplicationDashboard
 * The `ReplicationDashboard` component is a contextual component of the replication-page component.
 * It organizes cluster data specific to mode (dr or performance) and also the type (primary or secondary).
 * It is the parent contextual component of the replication-<name>-card components.
 *
 * @example
 * ```js
 * <ReplicationDashboard
    @data={{model}}
    @componentToRender='replication-primary-card'
    @isSecondary=false
    @isSummaryDashboard=false
    @replicationDetailsSummary={}
    @replicationDetails={{replicationDetails}}
    @clusterMode=primary
    @reindexingDetails={{reindexingDetails}}
    />
 * ```
 * @param {Object} data=null - An Ember data object that is pulled from the Ember Cluster Model.
 * @param {String} [componentToRender=''] - A string that determines which card component is displayed.  There are three options, replication-primary-card, replication-secondary-card, replication-summary-card.
 * @param {Boolean} [isSecondary=false] - Used to determine the title and display logic.
 * @param {Boolean} [isSummaryDashboard=false] -  Only true when the cluster is both a dr and performance primary. If true, replicationDetailsSummary is populated and used to pass through the cluster details.
 * @param {Object} replicationDetailsSummary=null - An Ember data object computed off the Ember Model.  It combines the Model.dr and Model.performance objects into one and contains details specific to the mode replication.
 * @param {Object} replicationDetails=null - An Ember data object pulled from the Ember Model. It contains details specific to the whether the replication is dr or performance.
 * @param {String} clusterMode=null - The cluster mode passed through to a table component.
 * @param {Object} reindexingDetails=null - An Ember data object used to show a reindexing progress bar.
 */

export default Component.extend({
  layout,
  componentToRender: '',
  data: null,
  isSecondary: false,
  isSummaryDashboard: false,
  replicationDetails: null,
  replicationDetailsSummary: null,
  isSyncing: computed('replicationDetails.state', 'isSecondary', function () {
    const { state } = this.replicationDetails;
    const isSecondary = this.isSecondary;
    return isSecondary && state && clusterStates([state]).isSyncing;
  }),
  isReindexing: computed('replicationDetails.reindex_in_progress', function () {
    const { replicationDetails } = this;
    return !!replicationDetails.reindex_in_progress;
  }),
  reindexingStage: computed('replicationDetails.reindex_stage', function () {
    const { replicationDetails } = this;
    const stage = replicationDetails.reindex_stage;
    // specify the stage if we have one
    if (stage) {
      return `: ${capitalize(stage)}`;
    }
    return '';
  }),
  progressBar: computed('replicationDetails.{reindex_building_progress,reindex_building_total}', function () {
    const { reindex_building_progress, reindex_building_total } = this.replicationDetails;
    let progressBar = null;

    if (reindex_building_progress && reindex_building_total) {
      progressBar = {
        value: reindex_building_progress,
        max: reindex_building_total,
      };
    }

    return progressBar;
  }),
  summaryState: computed('replicationDetailsSummary.{dr.state,performance.state}', function () {
    const { replicationDetailsSummary } = this;
    const drState = replicationDetailsSummary.dr.state;
    const performanceState = replicationDetailsSummary.performance.state;

    if (drState !== performanceState) {
      // when DR and Performance is enabled on the same cluster,
      // the states should always be the same
      // we are leaving this console log statement to be sure
      console.log('DR State: ', drState, 'Performance State: ', performanceState); // eslint-disable-line
    }

    return drState;
  }),
  reindexMessage: computed('isSecondary', 'progressBar', function () {
    if (!this.isSecondary) {
      return 'This can cause a delay depending on the size of the data store. You can <b>not</b> use Vault during this time.';
    }
    return 'This can cause a delay depending on the size of the data store. You can use Vault during this time.';
  }),
});
