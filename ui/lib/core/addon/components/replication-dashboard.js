import Component from '@ember/component';
import { computed } from '@ember/object';
import { clusterStates } from 'core/helpers/cluster-states';
import layout from '../templates/components/replication-dashboard';

export default Component.extend({
  layout,
  data: null,
  mode: computed('data', function() {
    const { data } = this;
    return data.replicationMode;
  }),
  isSecondary: computed('data', function() {
    const { data } = this;
    return data.replicationAttrs.isSecondary;
  }),
  dr: computed('data', function() {
    let dr = this.data.dr;
    if (!dr) {
      return false;
    }
    return dr;
  }),
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
