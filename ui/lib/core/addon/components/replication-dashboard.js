import Component from '@ember/component';
import { computed } from '@ember/object';
import { clusterStates } from 'core/helpers/cluster-states';
import layout from '../templates/components/replication-dashboard';

export default Component.extend({
  layout,
  data: null,
  mode: computed('data', function() {}),
  dr: computed('data', function() {
    let dr = this.data.dr;
    if (!dr) {
      return false;
    }
    return dr;
  }),
  isSyncing: computed('dr', function() {
    const { state } = this.dr;
    return state && clusterStates([state]).isSyncing;
  }),
  isReindexing: computed('data', function() {
    return true;
  }),
});
