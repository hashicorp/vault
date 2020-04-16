import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-page';

export default Component.extend({
  layout,
  dr: computed('model', function() {
    let dr = this.model.dr;
    if (!dr) {
      return false;
    }
    return dr;
  }),
  isDisabled: computed('dr', function() {
    if (this.dr.mode === 'disabled') {
      return true;
    }
    return false;
  }),
});
