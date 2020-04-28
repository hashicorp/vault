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
    let mode = this.model.replicationMode;
    return MODE[mode];
  }),
  isSecondary: computed('model', function() {
    const { model } = this;
    return model.replicationAttrs.isSecondary;
  }),
  replicationDetails: computed('model', function() {
    const { model } = this;
    const replicationMode = this.model.replicationMode;
    return model[replicationMode];
  }),
  isDisabled: computed('replicationDetails', function() {
    if (this.replicationDetails.mode === 'disabled' || this.replicationDetails.mode === 'primary') {
      return true;
    }
    return false;
  }),
  message: computed('model', function() {
    if (this.model.anyReplicationEnabled) {
      return `This ${this.mode} secondary has not been enabled.  You can do so from the ${this.mode} Primary.`;
    }
    return `This cluster has not been enabled as a ${this.mode} Secondary. You can do so by enabling replication and adding a secondary from the ${this.mode} Primary.`;
  }),
});
