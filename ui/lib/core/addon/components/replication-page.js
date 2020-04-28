import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-page';

const MODE = {
  dr: 'Disaster Recovery',
  performance: 'Performance',
};

export default Component.extend({
  layout,
  isSecondary: computed('model', function() {
    const { model } = this;
    return model.replicationAttrs.isSecondary;
  }),
  clusterMode: computed('model.{replicationAttrs}', function() {
    const { model } = this;
    return model.replicationAttrs.mode;
  }),
  replicationDetails: computed('model.{replicationMode}', function() {
    const { model } = this;
    const replicationMode = model.replicationMode;
    return model[replicationMode];
  }),
  isDisabled: computed('replicationDetails', function() {
    if (this.replicationDetails.mode === 'disabled' || this.replicationDetails.mode === 'primary') {
      return true;
    }
    return false;
  }),
  mode: computed('model.{replicationMode}', function() {
    let mode = this.model.replicationMode;
    return MODE[mode];
  }),
  message: computed('model.{anyReplicationEnabled}', 'mode', function() {
    if (this.model.anyReplicationEnabled) {
      return `This ${this.mode} secondary has not been enabled.  You can do so from the ${this.mode} Primary.`;
    }
    return `This cluster has not been enabled as a ${this.mode} Secondary. You can do so by enabling replication and adding a secondary from the ${this.mode} Primary.`;
  }),
});
