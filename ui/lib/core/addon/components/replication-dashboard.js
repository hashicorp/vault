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

  didReceiveAttrs() {
    this._super(arguments);
    console.log('the dashboard component received new attrs!');
  },

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
});
