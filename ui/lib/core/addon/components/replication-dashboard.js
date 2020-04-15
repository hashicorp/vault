import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-dashboard';

const MERKLE_STATES = { sync: 'merkle-sync', diff: 'merkle-diff' };

export default Component.extend({
  layout,
  data: null,
  dr: computed('data', function() {
    let dr = this.data.dr;
    if (!dr) {
      return false;
    }
    return dr;
  }),
  isSyncing: computed('data', function() {
    if (this.dr.state === MERKLE_STATES.sync || this.dr.state === MERKLE_STATES.diff) {
      return true;
    }
  }),
});
