import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-page';

const MODE = {
  dr: 'Disaster Recovery',
  performance: 'Performance',
};

export default Component.extend({
  layout,
  mode: computed('model', function() {
    let mode = this.model.rm.mode;
    return MODE[mode];
  }),
  dr: computed('model', function() {
    let dr = this.model.dr;
    if (!dr) {
      return false;
    }
    return dr;
  }),
  isDisabled: computed('dr', function() {
    // this conditional only applies to DR secondaries.
    if (this.dr.mode === 'disabled' || this.dr.mode === 'primary') {
      return true;
    }
    return false;
  }),
  message: computed('model', function() {
    if (this.model.anyReplicationEnabled) {
      return `This ${this.mode} secondary has not been enabled.  You can do so from the Disaster Recovery Primary.`;
    }
    return `This cluster has not been enabled as a ${this.mode} Secondary. You can do so by enabling replication and adding a secondary from the ${this.mode} Primary.`;
  }),
});
