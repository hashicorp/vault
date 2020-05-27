import Component from '@ember/component';
import { computed } from '@ember/object';
import { clusterStates } from 'core/helpers/cluster-states';
import { capitalize } from '@ember/string';
import layout from '../templates/components/replication-dashboard';

export default Component.extend({
  layout,
  data: null,
  replicationDetails: null,
  isSecondary: null,
  isSyncing: computed('replicationDetails.{state}', 'isSecondary', function() {
    const { state } = this.replicationDetails;
    const isSecondary = this.isSecondary;
    return isSecondary && state && clusterStates([state]).isSyncing;
  }),
  isReindexing: computed('replicationDetails.{reindex_in_progress}', function() {
    const { replicationDetails } = this;
    return !!replicationDetails.reindex_in_progress;
  }),
  reindexingStage: computed('replicationDetails.{reindex_stage}', function() {
    const { replicationDetails } = this;
    const stage = replicationDetails.reindex_stage;
    // specify the stage if we have one
    if (stage) {
      return `: ${capitalize(stage)}`;
    }
    return '';
  }),
  progressBar: computed('replicationDetails.{reindex_building_progress,reindex_building_total}', function() {
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
  summaryState: computed(
    'replicationDetailsSummary.dr.{state}',
    'replicationDetailsSummary.performance.{state}',
    function() {
      const { replicationDetailsSummary } = this;
      const drState = replicationDetailsSummary.dr.state;
      const performanceState = replicationDetailsSummary.performance.state;

      if (drState !== performanceState) {
        // when DR and Performance is enabled on the same cluster,
        // the states should always be the same
        // we are leaving this console log statement to be sure
        console.log('DR State: ', drState, 'Performance State: ', performanceState);
      }

      return drState;
    }
  ),
});
