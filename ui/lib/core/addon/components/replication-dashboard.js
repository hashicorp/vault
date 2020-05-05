import Component from '@ember/component';
import { computed } from '@ember/object';
import { clusterStates } from 'core/helpers/cluster-states';
import layout from '../templates/components/replication-dashboard';

export default Component.extend({
  layout,
  data: null,
  replicationDetails: null,
  isSecondary: null,
  isSyncing: computed('replicationDetails', 'isSecondary', function() {
    const { state } = this.replicationDetails;
    const isSecondary = this.isSecondary;
    return isSecondary && state && clusterStates([state]).isSyncing;
  }),
  isReindexing: computed('replicationDetails', function() {
    const { replicationDetails } = this;
    return !!replicationDetails.reindex_in_progress;
  }),
  reindexingStage: computed('replicationDetails', function() {
    const { replicationDetails } = this;
    const stage = replicationDetails.reindex_stage;
    // specify the stage if we have one
    if (stage) {
      return `: ${stage}`;
    }
    return '';
  }),
  reindexingProgress: computed('replicationDetails', function() {
    // TODO: use this value to display a progress bar
    const { replicationDetails } = this;
    const { reindex_building_progress, reindex_building_total } = replicationDetails;
    let progress = '';

    if (reindex_building_progress && reindex_building_total) {
      // convert progress to a percentage
      progress = (reindex_building_progress / reindex_building_total) * 100;
    }

    return progress;
  }),
});
