import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-dashboard';

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
    if (this.dr.state === 'merkle-sync' || this.dr.state === 'merkle-diff') {
      return true;
    }
  }),
});
