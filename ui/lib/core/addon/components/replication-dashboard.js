import Component from '@ember/component';
import { computed } from '@ember/object';
import { clusterStates } from 'core/helpers/cluster-states';
import layout from '../templates/components/replication-dashboard';

export default Component.extend({
  layout,
  data: null,
  mode: null,
  isSecondary: null,
  dr: null,
  isSyncing: computed('dr', 'isSecondary', function() {
    const { state } = this.dr;
    const isSecondary = this.isSecondary;
    return isSecondary && state && clusterStates([state]).isSyncing;
  }),
  isReindexing: computed('data', function() {
    // TODO: make this a real value
    return true;
  }),
});
