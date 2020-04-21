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
    if (this.dr.mode === 'disabled' || this.dr.mode === 'primary') {
      return true;
    }
    return false;
  }),
  title: computed('model', function() {
    let mode = this.model.rm.mode;
    if (mode === 'dr') {
      return 'Disaster Recovery';
    } else if (mode === 'performance') {
      return 'Performance';
    }
    return 'unknown';
  }),
  message: computed('model', function() {
    return 'This Disaster Recovery secondary has not been enabled.  You can do so from the Disaster Recovery Primary.';
  }),
});
